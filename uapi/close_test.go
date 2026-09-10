package uapi

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
	var first api.Runtime = Wrap(native)
	var second api.Runtime = Wrap(native)
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

	if err := first.Close(); err != nil {
		t.Fatal(err)
	}

	if out, err := first.Run(t.Context(), api.NewAnonymousSource("RETURN 45")); err != nil || out == nil || string(out.Content) != "45" {
		t.Fatalf("adapter after Close: output=%+v err=%v", out, err)
	}

	for _, compile := range []func(context.Context, api.Source, ...api.PlanOption) (api.Plan, error){first.Compile, first.CompileDebug} {
		compiled, err := compile(t.Context(), api.NewAnonymousSource("RETURN 1"))
		if err != nil {
			t.Fatalf("compile after adapter Close: %v", err)
		}

		if err := compiled.Close(); err != nil {
			t.Fatal(err)
		}
	}

	child, err := p.NewSession(t.Context())
	if err != nil {
		t.Fatalf("runtime close invalidated published plan: %v", err)
	}

	t.Cleanup(func() { _ = child.Close() })

	if out, err := child.Run(t.Context()); err != nil || out == nil || string(out.Content) != "42" {
		t.Fatalf("published child: output=%+v err=%v", out, err)
	}

	if out, err := second.Run(t.Context(), api.NewAnonymousSource("RETURN 43")); err != nil || out == nil || string(out.Content) != "43" {
		t.Fatalf("independent adapter: output=%+v err=%v", out, err)
	}

	if out, err := native.Run(t.Context(), source.NewAnonymous("RETURN 44")); err != nil || out == nil || string(out.Content) != "44" {
		t.Fatalf("native engine: output=%+v err=%v", out, err)
	}
}

func TestRuntimeCloseDoesNotWaitOrCancelNativeCalls(t *testing.T) {
	for _, owned := range []bool{false, true} {
		for _, mode := range []string{"compile", "debug", "run"} {
			t.Run(map[bool]string{false: "borrowed", true: "owned"}[owned]+"/"+mode, func(t *testing.T) {
				synctest.Test(t, func(t *testing.T) {
					entered, release := make(chan struct{}), make(chan struct{})
					releaseHook := sync.OnceFunc(func() { close(release) })
					defer releaseHook()
					caller := t.Context()
					block := func(ctx context.Context) error {
						if ctx != caller {
							t.Error("adapter replaced the caller context")
						}

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

					var r api.Runtime
					if owned {
						created, err := New(hook)
						if err != nil {
							t.Fatal(err)
						}

						r = created
						t.Cleanup(func() { _ = r.Close() })
					} else {
						r = newTestRuntime(t, hook)
					}

					var p api.Plan
					var output *api.Output
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
					if !closed {
						t.Fatal("Close waited for Native work")
					}

					releaseHook()
					synctest.Wait()
					if !closed || callErr != nil {
						t.Fatalf("admitted call was canceled or close stalled: closed=%t err=%v", closed, callErr)
					}

					if mode == "run" {
						if output == nil || string(output.Content) != "42" {
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

func TestClosePreservesNativeErrorsAndCleansExactlyOnce(t *testing.T) {
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
				if !errors.Is(result, cause) || !errors.As(result, &projected) || len(projected) != 1 {
					t.Fatalf("unstable or incomplete close result: %v", result)
				}
			}

			if calls.Load() != 1 {
				t.Fatalf("cleanup calls=%d", calls.Load())
			}

			if s, ok := target.(api.Session); ok {
				out, err := s.Run(t.Context())
				if !errors.Is(err, runtime.ErrInvalidOperation) || out != nil {
					t.Fatalf("closed Native session: output=%+v err=%v", out, err)
				}
			}
		})
	}
}

func TestRuntimeObservesBorrowedEngineClosure(t *testing.T) {
	native := newTestEngine(t)
	r := Wrap(native)
	if err := native.Close(); err != nil {
		t.Fatal(err)
	}

	for _, compile := range []func(context.Context, api.Source, ...api.PlanOption) (api.Plan, error){r.Compile, r.CompileDebug} {
		p, err := compile(t.Context(), api.NewAnonymousSource("RETURN 1"))
		if p != nil || !errors.Is(err, runtime.ErrInvalidOperation) {
			t.Fatalf("closed Native engine: plan=%v err=%v", p, err)
		}
	}

	if out, err := r.Run(t.Context(), api.NewAnonymousSource("RETURN 1")); !errors.Is(err, runtime.ErrInvalidOperation) || out != nil {
		t.Fatalf("closed Native engine: output=%+v err=%v", out, err)
	}
}

func TestPlanCloseWakesNativeCapacityWaiters(t *testing.T) {
	for _, debug := range []bool{false, true} {
		t.Run(map[bool]string{false: "ordinary", true: "debug"}[debug], func(t *testing.T) {
			synctest.Test(t, func(t *testing.T) {
				r := newTestRuntime(t, engine.WithMaxActiveSessions(1))
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
				var child io.Closer
				var createErr error
				finished := false
				go func() {
					if debug {
						child, createErr = p.NewDebugSession(t.Context())
					} else {
						child, createErr = p.NewSession(t.Context())
					}

					finished = true
				}()
				synctest.Wait()
				if finished {
					t.Fatal("constructor did not wait for capacity")
				}

				if err := p.Close(); err != nil {
					t.Fatal(err)
				}

				synctest.Wait()
				if !finished || child != nil || !errors.Is(createErr, runtime.ErrInvalidOperation) || errors.Is(createErr, context.Canceled) {
					t.Fatalf("finished=%t child=%v Native closure error=%v", finished, child, createErr)
				}

				if s, err := p.NewSession(t.Context()); s != nil || !errors.Is(err, runtime.ErrInvalidOperation) {
					t.Fatalf("closed plan: session=%v err=%v", s, err)
				}

				if s, err := p.NewDebugSession(t.Context()); s != nil || !errors.Is(err, runtime.ErrInvalidOperation) {
					t.Fatalf("closed plan: debug session=%v err=%v", s, err)
				}

				if out, err := holder.Run(t.Context()); err != nil || out == nil || string(out.Content) != "42" {
					t.Fatalf("borrowed Native VM: output=%+v err=%v", out, err)
				}
			})
		})
	}
}

func TestPlanCloseDoesNotWaitForPortableOptionCallbacks(t *testing.T) {
	for _, debug := range []bool{false, true} {
		t.Run(map[bool]string{false: "ordinary", true: "debug"}[debug], func(t *testing.T) {
			synctest.Test(t, func(t *testing.T) {
				var closes atomic.Int32
				r := newTestRuntime(t, engine.WithPlanCloseHook(func() error {
					closes.Add(1)

					return nil
				}))
				p, err := r.CompileDebug(t.Context(), api.NewAnonymousSource("RETURN 1"))
				if err != nil {
					t.Fatal(err)
				}

				t.Cleanup(func() { _ = p.Close() })
				entered, release := make(chan struct{}), make(chan struct{})
				releaseOption := sync.OnceFunc(func() { close(release) })
				defer releaseOption()
				option := func(api.SessionOptions) error {
					close(entered)
					<-release

					return nil
				}

				var child io.Closer
				var createErr error
				go func() {
					if debug {
						child, createErr = p.NewDebugSession(t.Context(), option)
					} else {
						child, createErr = p.NewSession(t.Context(), option)
					}
				}()
				<-entered
				closed := false
				go func() { _ = p.Close(); closed = true }()
				synctest.Wait()
				if !closed || closes.Load() != 1 {
					t.Fatal("adapter delayed Native plan cleanup")
				}

				releaseOption()
				synctest.Wait()
				if child != nil || !errors.Is(createErr, runtime.ErrInvalidOperation) || errors.Is(createErr, context.Canceled) {
					t.Fatalf("child=%v Native closure error=%v", child, createErr)
				}
			})
		})
	}
}
