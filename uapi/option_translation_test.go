package uapi

import (
	"context"
	"errors"
	"io"
	"reflect"
	"testing"
	"testing/synctest"
	"time"

	"github.com/MontFerret/api"

	"github.com/MontFerret/ferret/v2/pkg/engine"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

func TestOptionTranslationIsIndependentOfOperationContext(t *testing.T) {
	for _, operation := range []string{"compile", "compile debug", "run", "session", "debug session"} {
		for _, tc := range []struct {
			returnedContextError error
			wantContextError     error
			name                 string
			contextMode          string
			failCallbacks        bool
		}{
			{name: "nil", contextMode: "nil", wantContextError: runtime.ErrInvalidArgument},
			{name: "canceled", contextMode: "canceled", wantContextError: context.Canceled},
			{name: "deadline", contextMode: "deadline", wantContextError: context.DeadlineExceeded},
			{name: "cancel during translation", contextMode: "cancel during", wantContextError: context.Canceled},
			{name: "callback failures", failCallbacks: true},
			{name: "failures with nil context", contextMode: "nil", failCallbacks: true},
			{name: "failures after cancellation", contextMode: "canceled", failCallbacks: true},
			{name: "failures after deadline", contextMode: "deadline", failCallbacks: true},
			{name: "failures during cancellation", contextMode: "cancel during", failCallbacks: true},
			{name: "callback returns cancellation", failCallbacks: true, returnedContextError: context.Canceled, wantContextError: context.Canceled},
			{name: "callback returns deadline", failCallbacks: true, returnedContextError: context.DeadlineExceeded, wantContextError: context.DeadlineExceeded},
		} {
			t.Run(operation+"/"+tc.name, func(t *testing.T) {
				var compileCalls int
				r := newTestRuntime(t, engine.WithBeforeCompileHook(func(context.Context) error {
					compileCalls++

					return nil
				}))
				p, err := r.CompileDebug(t.Context(), api.NewAnonymousSource("RETURN 1"))
				if err != nil {
					t.Fatal(err)
				}

				t.Cleanup(func() { _ = p.Close() })
				compileCalls = 0
				ctx, cancel := context.WithCancel(t.Context())
				defer cancel()
				switch tc.contextMode {
				case "nil":
					ctx = nil
				case "canceled":
					cancel()
				case "deadline":
					var stop context.CancelFunc
					ctx, stop = context.WithDeadline(t.Context(), time.Now().Add(-time.Second))
					defer stop()
				}

				first, second := errors.New("first callback"), errors.New("second callback")
				var order []int
				callbacks := []func() error{func() error {
					order = append(order, 1)
					if tc.contextMode == "cancel during" {
						cancel()
					}

					if tc.failCallbacks {
						return errors.Join(first, tc.returnedContextError)
					}

					return nil
				}, nil, func() error {
					order = append(order, 2)
					if tc.failCallbacks {
						return second
					}

					return nil
				}}

				err = callRejectedOptionOperation(t, r, p, operation, ctx, callbacks)
				if err == nil || !reflect.DeepEqual(order, []int{1, 2}) || compileCalls != 0 {
					t.Fatalf("error=%v callbacks=%v Native compile hooks=%d", err, order, compileCalls)
				}

				for _, cause := range []error{first, second} {
					if errors.Is(err, cause) != tc.failCallbacks {
						t.Fatalf("callback cause %v: error=%v", cause, err)
					}
				}

				if tc.failCallbacks {
					want := errors.Join(first, tc.returnedContextError, second)
					if err.Error() != want.Error() {
						t.Fatalf("callback error order: got=%v want=%v", err, want)
					}
				}

				for _, cause := range []error{runtime.ErrInvalidArgument, context.Canceled, context.DeadlineExceeded} {
					if errors.Is(err, cause) != errors.Is(tc.wantContextError, cause) {
						t.Fatalf("context cause %v: error=%v, want context error=%v", cause, err, tc.wantContextError)
					}
				}
			})
		}
	}
}

func TestCallerCancellationWakesNativeCapacityWaiters(t *testing.T) {
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
				ctx, cancel := context.WithCancel(t.Context())
				defer cancel()
				var child io.Closer
				var createErr error
				finished := false
				go func() {
					if debug {
						child, createErr = p.NewDebugSession(ctx)
					} else {
						child, createErr = p.NewSession(ctx)
					}

					finished = true
				}()
				synctest.Wait()
				if finished {
					t.Fatal("constructor did not wait for capacity")
				}

				cancel()
				synctest.Wait()
				if !finished || child != nil || !errors.Is(createErr, context.Canceled) {
					t.Fatalf("canceled waiter: finished=%t session=%v error=%v", finished, child, createErr)
				}

				if out, err := holder.Run(t.Context()); err != nil || out == nil || string(out.Content) != "42" {
					t.Fatalf("independent session: output=%+v err=%v", out, err)
				}

				if err := holder.Close(); err != nil {
					t.Fatal(err)
				}

				next, err := p.NewSession(t.Context())
				if err != nil {
					t.Fatal(err)
				}

				if err := next.Close(); err != nil {
					t.Fatal(err)
				}
			})
		})
	}
}

func callRejectedOptionOperation(t *testing.T, r api.Runtime, p api.Plan, operation string, ctx context.Context, callbacks []func() error) error {
	t.Helper()

	planOptions := make([]api.PlanOption, len(callbacks))
	sessionOptions := make([]api.SessionOption, len(callbacks))
	for i, callback := range callbacks {
		if callback != nil {
			planOptions[i] = func(api.PlanOptions) error { return callback() }
			sessionOptions[i] = func(api.SessionOptions) error { return callback() }
		}
	}

	src := api.NewAnonymousSource("RETURN 1")
	var child io.Closer
	var err error
	switch operation {
	case "compile":
		child, err = r.Compile(ctx, src, planOptions...)
	case "compile debug":
		child, err = r.CompileDebug(ctx, src, planOptions...)
	case "session":
		child, err = p.NewSession(ctx, sessionOptions...)
	case "debug session":
		child, err = p.NewDebugSession(ctx, sessionOptions...)
	case "run":
		var output *api.Output
		output, err = r.Run(ctx, src, sessionOptions...)
		if output != nil {
			t.Fatalf("rejected Run returned output=%+v error=%v", output, err)
		}
	default:
		t.Fatalf("unknown operation %q", operation)
	}

	if child != nil {
		_ = child.Close()

		t.Fatalf("rejected operation published a resource: error=%v", err)
	}

	return err
}
