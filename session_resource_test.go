package ferret

import (
	"errors"
	"io"
	"strings"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/MontFerret/ferret/v2/internal/resource"
	ferretfs "github.com/MontFerret/ferret/v2/pkg/fs"
	ferretnet "github.com/MontFerret/ferret/v2/pkg/net"
	"github.com/MontFerret/ferret/v2/pkg/vm"
)

func TestSessionHostResourcesAreIsolated(t *testing.T) {
	for _, debug := range []bool{false, true} {
		for _, owned := range []bool{false, true} {
			name := map[bool]string{false: "ordinary", true: "debug"}[debug] + "/" + map[bool]string{false: "borrowed", true: "owned"}[owned]
			t.Run(name, func(t *testing.T) {
				client := &recordingHTTPClient{}
				engine := mustNewEngine(t, WithFSRoot(t.TempDir()), WithNetworkOptions(ferretnet.WithHTTPClient(client)), WithMaxActiveSessions(2))
				t.Cleanup(func() { _ = engine.Close() })
				plan, err := engine.CompileDebug(t.Context(), NewAnonymousSource("RETURN 1"))
				if err != nil {
					t.Fatal(err)
				}

				t.Cleanup(func() { _ = plan.Close() })
				sibling := mustNewSession(t, plan, WithSessionFSRoot(t.TempDir()))
				t.Cleanup(func() { _ = sibling.Close() })
				siblingFS := sibling.fs
				var setters []SessionOption

				if owned {
					setters = []SessionOption{WithSessionFSRoot(t.TempDir())}
				}

				var filesystem ferretfs.FileSystem
				session, err := newPlanSession(plan, t.Context(), setters, planSessionSetup{requiresDebugInfo: debug},
					func(dependencies planSessionDependencies) (io.Closer, error) {
						filesystem = dependencies.filesystem

						if debug {
							return buildDebugSession(dependencies)
						}

						return buildSession(dependencies)
					})
				if err != nil {
					t.Fatal(err)
				}

				t.Cleanup(func() { _ = session.Close() })

				for range 2 {
					if err := session.Close(); err != nil {
						t.Fatal(err)
					}
				}

				if _, err := filesystem.Stat("."); (err != nil) != owned {
					t.Fatalf("session filesystem stat = %v, want closed=%t", err, owned)
				}

				if _, err := engine.host.fs.Stat("."); err != nil {
					t.Fatalf("session closed the Engine filesystem: %v", err)
				}

				if _, err := siblingFS.Stat("."); err != nil {
					t.Fatalf("session closed its sibling's filesystem: %v", err)
				}

				if client.idleCloseCount() != 0 {
					t.Fatal("session closed the Engine network")
				}

				if got := len(engine.limiter.ch); got != 1 {
					t.Fatalf("active permits = %d, want only the sibling's permit", got)
				}

				if err := sibling.Close(); err != nil {
					t.Fatal(err)
				}

				if err := plan.Close(); err != nil {
					t.Fatal(err)
				}

				if err := engine.Close(); err != nil {
					t.Fatal(err)
				}

				if client.idleCloseCount() != 1 {
					t.Fatal("Engine did not retain network ownership")
				}
			})
		}
	}
}

func TestSessionManagerCleanupIsConcurrentAndRetainsErrors(t *testing.T) {
	for _, debug := range []bool{false, true} {
		t.Run(map[bool]string{false: "ordinary", true: "debug"}[debug], func(t *testing.T) {
			hookErr, closeErr := errors.New("hook failed"), errors.New("service failed")
			var hooks, closes atomic.Int32
			engine := mustNewEngine(t, WithMaxActiveSessions(1), WithMaxIdleVMsPerPlan(1), WithSessionCloseHook(func() error {
				hooks.Add(1)

				return hookErr
			}))
			t.Cleanup(func() { _ = engine.Close() })
			plan, err := engine.CompileDebug(t.Context(), NewAnonymousSource("RETURN 1"))
			if err != nil {
				t.Fatal(err)
			}

			t.Cleanup(func() { _ = plan.Close() })
			var filesystem ferretfs.FileSystem
			var instance *vm.VM

			session, err := newPlanSession(plan, t.Context(), []SessionOption{WithSessionFSRoot(t.TempDir())}, planSessionSetup{requiresDebugInfo: debug},
				func(dependencies planSessionDependencies) (io.Closer, error) {
					filesystem = dependencies.filesystem
					if err := dependencies.resources.Own("service", func() error {
						closes.Add(1)

						if hooks.Load() != 1 || len(engine.limiter.ch) != 1 {
							t.Error("resource cleanup did not run between hooks and permit release")
						}

						return closeErr
					}); err != nil {
						return nil, err
					}

					if debug {
						return buildDebugSession(dependencies)
					}

					child, err := buildSession(dependencies)
					if err == nil {
						instance = child.vm
					}

					return child, err
				})
			if err != nil {
				t.Fatal(err)
			}

			t.Cleanup(func() { _ = session.Close() })
			results := make([]error, 8)
			var wg sync.WaitGroup

			for i := range results {
				wg.Go(func() { results[i] = session.Close() })
			}

			wg.Wait()

			for _, err := range results {
				if !errors.Is(err, hookErr) || !errors.Is(err, closeErr) || err != results[0] {
					t.Fatalf("concurrent closure lost cleanup errors: %v", err)
				}
			}

			if err := session.Close(); err != results[0] {
				t.Fatalf("repeated closure changed the result: %v", err)
			}

			if hooks.Load() != 1 || closes.Load() != 1 || len(engine.limiter.ch) != 0 {
				t.Fatalf("hooks/closes/permits = %d/%d/%d, want 1/1/0", hooks.Load(), closes.Load(), len(engine.limiter.ch))
			}

			if _, err := filesystem.Stat("."); err == nil {
				t.Fatal("resource cleanup failure left the filesystem open")
			}

			next := mustNewSession(t, plan)
			t.Cleanup(func() { _ = next.Close() })

			if !debug && next.vm != instance {
				t.Fatal("resource cleanup failure prevented returning the VM to its pool")
			}
		})
	}
}

func TestSessionManagerRollsBackBuilderFailures(t *testing.T) {
	for _, mode := range []string{"error", "builder panic", "cleanup panic"} {
		t.Run(mode, func(t *testing.T) {
			engine := mustNewEngine(t, WithFSRoot(t.TempDir()), WithMaxActiveSessions(1))
			t.Cleanup(func() { _ = engine.Close() })
			plan := mustCompilePlan(t, engine, "RETURN 1")
			t.Cleanup(func() { _ = plan.Close() })
			buildErr, closeErr := errors.New("builder failed"), errors.New("filesystem failed")
			filesystem := newCountingCloseFileSystem(t, t.TempDir(), closeErr)
			var recovered any
			var resultErr error

			func() {
				defer func() { recovered = recover() }()

				_, resultErr = newPlanSession(plan, t.Context(), nil, planSessionSetup{},
					func(dependencies planSessionDependencies) (*Session, error) {
						if err := dependencies.resources.Own(resource.FileSystem, func() error {
							err := filesystem.Close()
							if mode == "cleanup panic" {
								panic(closeErr)
							}

							return err
						}); err != nil {
							return nil, err
						}

						if mode == "builder panic" {
							panic(buildErr)
						}

						return nil, buildErr
					})
			}()

			switch mode {
			case "error":
				if recovered != nil || !errors.Is(resultErr, buildErr) || !errors.Is(resultErr, closeErr) || !strings.Contains(resultErr.Error(), "close filesystem") {
					t.Fatalf("rollback panic/error = %v/%v, want joined builder/filesystem errors", recovered, resultErr)
				}
			case "builder panic":
				if recovered != buildErr {
					t.Fatalf("builder panic = %v, want %v", recovered, buildErr)
				}
			case "cleanup panic":
				if recovered != closeErr {
					t.Fatalf("cleanup panic = %v, want %v", recovered, closeErr)
				}
			}

			if filesystem.closeCalls.Load() != 1 || len(engine.limiter.ch) != 0 {
				t.Fatalf("filesystem closes/permits = %d/%d, want 1/0", filesystem.closeCalls.Load(), len(engine.limiter.ch))
			}

			if _, err := filesystem.Stat("."); err == nil {
				t.Fatal("builder failure left the filesystem open")
			}

			if _, err := engine.host.fs.Stat("."); err != nil {
				t.Fatalf("rollback closed the borrowed Engine filesystem: %v", err)
			}

			next := mustNewSession(t, plan)
			if err := next.Close(); err != nil {
				t.Fatal(err)
			}
		})
	}
}
