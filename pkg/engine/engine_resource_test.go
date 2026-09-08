package engine

import (
	"errors"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/engine/internal/resource"
	ferretfs "github.com/MontFerret/ferret/v2/pkg/fs"
	"github.com/MontFerret/ferret/v2/pkg/logging"
	"github.com/MontFerret/ferret/v2/pkg/module"
	ferretnet "github.com/MontFerret/ferret/v2/pkg/net"
	"github.com/MontFerret/ferret/v2/pkg/stdlib"
)

func TestNewRollsBackResourcesAtEveryConstructionStage(t *testing.T) {
	failure := errors.New("construction failed")

	for _, stage := range []struct {
		option    Option
		name      string
		wantHooks bool
	}{
		{name: "option", option: func(*config) error { return failure }},
		{name: "stdlib", option: WithStdlib(stdlib.Only(stdlib.Group("unknown")))},
		{name: "compiler validation", option: func(opts *config) error {
			opts.optimizationLevel = OptimizationLevel(-1)

			return nil
		}},
		{name: "bootstrap logger", option: func(opts *config) error {
			opts.logger = append(opts.logger, logging.WithLevel(logging.LogLevel(100)))

			return nil
		}},
		{name: "bootstrap filesystem", option: WithFSRoot(filepath.Join(t.TempDir(), "missing"))},
		{name: "module registration", wantHooks: true, option: WithModules(testModule{
			registerFn: func(module.Bootstrap) error { return failure },
		})},
		{name: "host build", wantHooks: true, option: WithModules(testModule{
			registerFn: func(boot module.Bootstrap) error {
				boot.Host().Library().Function().A0().Add("RESOURCE_DUPLICATE", testFn0)
				boot.Host().Library().Function().A0().Add("RESOURCE_DUPLICATE", testFn0)

				return nil
			},
		})},
		{name: "init hooks", wantHooks: true, option: WithEngineInitHook(func() error { return failure })},
	} {
		for _, borrowed := range []bool{false, true} {
			ownership := map[bool]string{false: "owned", true: "borrowed"}[borrowed]
			t.Run(stage.name+"/"+ownership, func(t *testing.T) {
				client := &recordingHTTPClient{}
				networkOption := WithNetworkOptions(ferretnet.WithHTTPClient(client))
				wantCloses := 1

				if borrowed {
					networkOption = WithNetwork(mustNewTestNetwork(t, ferretnet.WithHTTPClient(client)))
					wantCloses = 0
				}

				var filesystem ferretfs.FileSystem
				hookCalled := false
				engine, err := New(
					networkOption,
					WithFSRoot(t.TempDir()),
					WithModules(testModule{registerFn: func(boot module.Bootstrap) error {
						filesystem = boot.Host().FileSystem()

						return nil
					}}),
					WithEngineCloseHook(func() error {
						hookCalled = true

						if filesystem != nil {
							if _, err := filesystem.Stat("."); err != nil {
								t.Errorf("filesystem closed before hooks: %v", err)
							}
						}

						if got := client.idleCloseCount(); got != 0 {
							t.Errorf("network closed before hooks: %d", got)
						}

						return nil
					}),
					stage.option,
				)
				if engine != nil {
					_ = engine.Close()

					t.Fatal("failed construction returned an engine")
				}

				if err == nil {
					t.Fatal("expected construction failure")
				}

				if hookCalled != stage.wantHooks {
					t.Fatalf("close hook called = %t, want %t", hookCalled, stage.wantHooks)
				}

				if got := client.idleCloseCount(); got != wantCloses {
					t.Fatalf("network closes = %d, want %d", got, wantCloses)
				}

				if filesystem != nil {
					if _, err := filesystem.Stat("."); err == nil {
						t.Fatal("construction rollback left the filesystem open")
					}
				}
			})
		}
	}
}

func TestEngineResourceCleanupIsOnceOnly(t *testing.T) {
	for _, failing := range []bool{false, true} {
		t.Run(map[bool]string{false: "success", true: "cleanup errors"}[failing], func(t *testing.T) {
			filesystem, err := ferretfs.New(ferretfs.WithRoot(t.TempDir()))
			if err != nil {
				t.Fatal(err)
			}

			var hookErr, filesystemErr error

			if failing {
				hookErr, filesystemErr = errors.New("hook failed"), errors.New("filesystem failed")
			}

			replacement := &failingCloseFileSystem{FileSystem: filesystem, closeErr: filesystemErr}
			client := &recordingHTTPClient{}
			var hooks atomic.Int32

			engine := mustNewEngine(t,
				WithNetworkOptions(ferretnet.WithHTTPClient(client)),
				WithEngineCloseHook(func() error {
					hooks.Add(1)

					if _, err := filesystem.Stat("."); err != nil {
						t.Errorf("filesystem closed before hook: %v", err)
					}

					if client.idleCloseCount() != 0 {
						t.Error("network closed before hook")
					}

					return hookErr
				}),
			)
			if err := engine.resources.Own(resource.FileSystem, replacement.Close); err != nil {
				t.Fatal(err)
			}

			engine.host.fs = replacement
			results := make([]error, 8)
			var wg sync.WaitGroup

			for i := range results {
				wg.Go(func() { results[i] = engine.Close() })
			}

			wg.Wait()

			for _, err := range results {
				if !errors.Is(err, hookErr) || !errors.Is(err, filesystemErr) {
					t.Fatalf("lost cleanup cause: %v", err)
				}

				if err != results[0] {
					t.Fatalf("concurrent close did not retain the same result: %v", err)
				}
			}

			if err := engine.Close(); err != results[0] {
				t.Fatalf("repeated close did not retain the same result: %v", err)
			}

			if hooks.Load() != 1 || replacement.closes.Load() != 1 || client.idleCloseCount() != 1 {
				t.Fatalf("hook/filesystem/network closes = %d/%d/%d, want 1/1/1", hooks.Load(), replacement.closes.Load(), client.idleCloseCount())
			}

			if _, err := filesystem.Stat("."); err == nil {
				t.Fatal("filesystem remains open")
			}
		})
	}
}

func TestNetworkReplacementRetiresBeforeNextOption(t *testing.T) {
	first, second := &recordingHTTPClient{}, &recordingHTTPClient{}
	optionErr := errors.New("option failed")
	engine, err := New(
		WithNetworkOptions(ferretnet.WithHTTPClient(first)),
		func(*config) error { return optionErr },
		WithNetworkOptions(ferretnet.WithHTTPClient(second)),
		func(opts *config) error {
			if first.idleCloseCount() != 1 || second.idleCloseCount() != 0 {
				t.Errorf("replacement closes = %d/%d, want 1/0", first.idleCloseCount(), second.idleCloseCount())
			}

			if opts.network.HTTP() != second {
				t.Error("later network option was not applied after failure")
			}

			return nil
		},
	)
	if engine != nil || !errors.Is(err, optionErr) {
		t.Fatalf("construction result = %v/%v, want nil/option error", engine, err)
	}

	if first.idleCloseCount() != 1 || second.idleCloseCount() != 1 {
		t.Fatalf("rollback closes = %d/%d, want 1/1", first.idleCloseCount(), second.idleCloseCount())
	}
}

func TestNewConfigJoinsOptionAndResourceCleanupErrors(t *testing.T) {
	optionErr, closeErr := errors.New("option failed"), errors.New("resource failed")
	closes := 0
	_, err := newConfig([]Option{
		func(opts *config) error {
			return opts.resources.Own("future", func() error {
				closes++

				return closeErr
			})
		},
		func(*config) error { return optionErr },
	})
	if !errors.Is(err, optionErr) || !errors.Is(err, closeErr) {
		t.Fatalf("lost construction/cleanup cause: %v", err)
	}

	if closes != 1 || !strings.Contains(err.Error(), "close future") {
		t.Fatalf("cleanup = %d/%v, want one labeled cleanup failure", closes, err)
	}
}
