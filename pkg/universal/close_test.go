package universal

import (
	"context"
	"errors"
	"io"
	"sync"
	"sync/atomic"
	"testing"
	"testing/synctest"

	"github.com/MontFerret/api"
	apidebugger "github.com/MontFerret/api/debugger"
	apidiagnostics "github.com/MontFerret/api/diagnostics"

	"github.com/MontFerret/ferret/v2/pkg/diagnostics"
	"github.com/MontFerret/ferret/v2/pkg/engine"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/source"
)

func TestRuntimeBorrowsEngineAndLeavesChildrenIndependent(t *testing.T) {
	var closes atomic.Int32
	native := newTestEngine(t, engine.WithEngineCloseHook(func() error {
		closes.Add(1)

		return nil
	}))
	var first api.Runtime = New(native)
	var second api.Runtime = New(native)
	t.Cleanup(func() { _ = second.Close() })
	p, err := first.Compile(t.Context(), api.NewAnonymousSource("RETURN 42"))
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() { _ = p.Close() })

	if err := first.Close(); err != nil {
		t.Fatal(err)
	}

	if closes.Load() != 0 {
		t.Fatal("adapter closed borrowed engine")
	}

	child, err := p.NewSession(t.Context())
	if err != nil {
		t.Fatalf("runtime close invalidated published plan: %v", err)
	}

	t.Cleanup(func() { _ = child.Close() })

	if out, err := child.Run(t.Context()); err != nil || string(out.Content) != "42" {
		t.Fatalf("published child: output=%+v err=%v", out, err)
	}

	if out, err := second.Run(t.Context(), api.NewAnonymousSource("RETURN 43")); err != nil || string(out.Content) != "43" {
		t.Fatalf("independent adapter: output=%+v err=%v", out, err)
	}

	if out, err := native.Run(t.Context(), source.NewAnonymous("RETURN 44")); err != nil || out == nil || string(out.Content) != "44" {
		t.Fatalf("native engine: output=%+v err=%v", out, err)
	}
}

func TestRuntimeCloseWaitsWithoutCancelingAdmittedCalls(t *testing.T) {
	for _, mode := range []string{"compile", "debug", "run"} {
		t.Run(mode, func(t *testing.T) {
			synctest.Test(t, func(t *testing.T) {
				entered, release := make(chan struct{}), make(chan struct{})
				block := func(ctx context.Context) error {
					close(entered)
					<-release

					return ctx.Err()
				}

				hook := engine.WithBeforeCompileHook(block)
				if mode == "run" {
					hook = engine.WithBeforeRunHook(func(ctx context.Context) (context.Context, error) {
						return ctx, block(ctx)
					})
				}

				r := newTestRuntime(t, hook)
				var p api.Plan
				var output api.Output
				var callErr error
				go func() {
					src := api.NewAnonymousSource("RETURN 42")
					switch mode {
					case "compile":
						p, callErr = r.Compile(t.Context(), src)
					case "debug":
						p, callErr = r.CompileDebug(t.Context(), src)
					case "run":
						output, callErr = r.Run(t.Context(), src)
					}
				}()
				<-entered
				closed := false
				go func() { _ = r.Close(); closed = true }()
				synctest.Wait()
				if closed {
					t.Fatal("Close returned before admitted work settled")
				}

				calls := 0
				if rejected, err := r.Compile(t.Context(), api.NewAnonymousSource("RETURN 1"), func(api.PlanOptions) error {
					calls++

					return nil
				}); rejected != nil || !errors.Is(err, runtime.ErrInvalidOperation) || calls != 0 {
					t.Fatalf("closed admission: plan=%v err=%v callbacks=%d", rejected, err, calls)
				}

				close(release)
				synctest.Wait()
				if !closed || callErr != nil {
					t.Fatalf("admitted call was canceled or close stalled: closed=%t err=%v", closed, callErr)
				}

				if mode == "run" {
					if string(output.Content) != "42" {
						t.Fatalf("output=%+v", output)
					}
				} else {
					if p == nil {
						t.Fatal("admitted compile lost its plan")
					}

					if err := p.Close(); err != nil {
						t.Fatal(err)
					}
				}
			})
		})
	}
}

func TestPlanCloseSettlesPendingConstructorsBeforeCleanup(t *testing.T) {
	for _, debug := range []bool{false, true} {
		for _, blockedOption := range []bool{false, true} {
			name := map[bool]string{false: "ordinary", true: "debug"}[debug] + "/" + map[bool]string{false: "capacity", true: "option"}[blockedOption]
			t.Run(name, func(t *testing.T) {
				synctest.Test(t, func(t *testing.T) {
					var planCloses atomic.Int32
					r := newTestRuntime(t, engine.WithMaxActiveSessions(1), engine.WithPlanCloseHook(func() error {
						planCloses.Add(1)

						return nil
					}))
					p, err := r.CompileDebug(t.Context(), api.NewAnonymousSource("RETURN 42"))
					if err != nil {
						t.Fatal(err)
					}

					t.Cleanup(func() { _ = p.Close() })
					holder, err := p.NewSession(t.Context())
					if err != nil {
						t.Fatal(err)
					}

					t.Cleanup(func() { _ = holder.Close() })
					entered, release := make(chan struct{}), make(chan struct{})
					failure := errors.New("option failure")
					option := func(api.SessionOptions) error {
						close(entered)
						if blockedOption {
							<-release

							return failure
						}

						return nil
					}

					var child io.Closer
					var createErr error
					finished := false
					go func() {
						if debug {
							child, createErr = p.NewDebugSession(t.Context(), option)
						} else {
							child, createErr = p.NewSession(t.Context(), option)
						}

						finished = true
					}()
					<-entered
					synctest.Wait()
					if finished {
						t.Fatal("constructor did not block")
					}

					closed := false
					go func() { _ = p.Close(); closed = true }()
					synctest.Wait()
					if blockedOption {
						if closed || planCloses.Load() != 0 {
							t.Fatal("native cleanup ran while an option was active")
						}

						close(release)
						synctest.Wait()
					}

					if !closed || !finished || child != nil || !errors.Is(createErr, context.Canceled) || planCloses.Load() != 1 {
						t.Fatalf("close=%t finished=%t child=%v err=%v cleanup=%d", closed, finished, child, createErr, planCloses.Load())
					}

					if blockedOption && !errors.Is(createErr, failure) {
						t.Fatalf("lost option failure: %v", createErr)
					}

					// The successfully published holder owns its borrowed VM, even
					// after the plan pool has closed its idle resources.
					if out, err := holder.Run(t.Context()); err != nil || string(out.Content) != "42" {
						t.Fatalf("published session revoked: output=%+v err=%v", out, err)
					}
				})
			})
		}
	}
}

func TestPublishedDebugSessionDoesNotInheritConstructorCancellation(t *testing.T) {
	r := newTestRuntime(t)
	p, err := r.CompileDebug(t.Context(), api.NewAnonymousSource("RETURN @value"))
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() { _ = p.Close() })
	ctx, cancel := context.WithCancel(t.Context())
	defer cancel()
	s, err := p.NewDebugSession(ctx, api.WithParam("value", 42))
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() { _ = s.Close() })
	cancel()
	if err := p.Close(); err != nil {
		t.Fatal(err)
	}

	if _, err := s.Start(t.Context()); err != nil {
		t.Fatal(err)
	}

	event, err := s.Continue(t.Context())
	if err != nil || event == nil || event.Reason != apidebugger.ReasonCompleted || event.Output == nil || string(event.Output.Content) != "42" {
		t.Fatalf("published debugger: event=%+v err=%v", event, err)
	}
}

func TestCloseCachesProjectedErrorsAndCleansExactlyOnce(t *testing.T) {
	for _, kind := range []string{"plan", "session", "debug"} {
		t.Run(kind, func(t *testing.T) {
			cause := diagnostics.NewUnexpectedError(source.Source{}, "close diagnostic")
			var calls atomic.Int32
			hook := func() error {
				calls.Add(1)

				return cause
			}

			option := engine.WithSessionCloseHook(hook)
			if kind == "plan" {
				option = engine.WithPlanCloseHook(hook)
			}

			r := newTestRuntime(t, option)
			p, err := r.CompileDebug(t.Context(), api.NewAnonymousSource("RETURN 1"))
			if err != nil {
				t.Fatal(err)
			}

			t.Cleanup(func() { _ = p.Close() })
			var target io.Closer = p
			switch kind {
			case "session":
				target, err = p.NewSession(t.Context())
			case "debug":
				target, err = p.NewDebugSession(t.Context())
			}

			if err != nil {
				t.Fatal(err)
			}

			results := make([]error, 8)
			var wg sync.WaitGroup
			for i := range results {
				wg.Go(func() { results[i] = target.Close() })
			}

			wg.Wait()
			for _, result := range results {
				var projected apidiagnostics.Diagnostics
				if result != results[0] || !errors.Is(result, cause) || !errors.As(result, &projected) || len(projected) != 1 {
					t.Fatalf("unstable or incomplete close result: %v", result)
				}
			}

			if calls.Load() != 1 {
				t.Fatalf("cleanup calls=%d", calls.Load())
			}
		})
	}
}

func TestPanickingOptionsReleaseAdmission(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		r := newTestRuntime(t)
		p, err := r.Compile(t.Context(), api.NewAnonymousSource("RETURN 1"))
		if err != nil {
			t.Fatal(err)
		}

		for _, call := range []func(){
			func() {
				_, _ = r.Compile(t.Context(), api.NewAnonymousSource("RETURN 1"), func(api.PlanOptions) error { panic("option") })
			},
			func() { _, _ = p.NewSession(t.Context(), func(api.SessionOptions) error { panic("option") }) },
		} {
			func() {
				defer func() {
					if recover() == nil {
						t.Error("option panic was swallowed")
					}
				}()
				call()
			}()
		}

		if err := p.Close(); err != nil {
			t.Fatal(err)
		}

		if err := r.Close(); err != nil {
			t.Fatal(err)
		}
	})
}
