package engine

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

func TestRunHooksSettleAbortedAdmission(t *testing.T) {
	for _, debug := range []bool{false, true} {
		for _, mode := range []string{"caller", "derived", "deadline", "nil", "before_failure", "pre_canceled"} {
			t.Run(map[bool]string{false: "normal/", true: "debug/"}[debug]+mode, func(t *testing.T) {
				type contextKey struct{}
				caller := context.WithValue(t.Context(), contextKey{}, "caller")
				ctx, cancel := context.WithCancel(caller)
				t.Cleanup(cancel)
				beforeErr, afterErr := errors.New("before failed"), errors.New("after failed")
				wantErr := error(context.Canceled)
				if mode == "deadline" {
					wantErr = context.DeadlineExceeded
				} else if mode == "nil" {
					wantErr = runtime.ErrInvalidArgument
				} else if mode == "before_failure" {
					wantErr = beforeErr
				}

				var calls []string
				var afterContexts []any
				var afterErrors []error
				var executions int
				abort := true
				engine := mustNewEngine(t,
					WithFunctionsRegistrar(func(ns runtime.Namespace) {
						ns.Function().A0().Add("RUN_PROBE", func(context.Context) (runtime.Value, error) {
							executions++

							return runtime.NewInt(42), nil
						})
					}),
					WithBeforeRunHook(func(c context.Context) (context.Context, error) {
						calls = append(calls, "before")
						derived := context.WithValue(c, contextKey{}, "hook")
						if !abort {
							return derived, nil
						}

						switch mode {
						case "caller":
							cancel()
						case "derived":
							canceled, stop := context.WithCancel(derived)
							stop()

							return canceled, nil
						case "deadline":
							expired, stop := context.WithDeadline(derived, time.Unix(1, 0))
							defer stop()

							return expired, nil
						case "nil":
							return nil, nil
						case "before_failure":
							return derived, beforeErr
						}

						return derived, nil
					}),
					WithAfterRunHook(func(c context.Context, err error) error {
						calls = append(calls, "after-first")
						afterContexts = append(afterContexts, c.Value(contextKey{}))
						afterErrors = append(afterErrors, err)

						return nil
					}),
					WithAfterRunHook(func(c context.Context, err error) error {
						calls = append(calls, "after-second")
						afterContexts = append(afterContexts, c.Value(contextKey{}))
						afterErrors = append(afterErrors, err)
						if abort {
							return afterErr
						}

						return nil
					}),
				)
				t.Cleanup(func() { _ = engine.Close() })
				plan, err := engine.CompileDebug(t.Context(), NewAnonymousSource("RETURN RUN_PROBE()"))
				if err != nil {
					t.Fatal(err)
				}

				t.Cleanup(func() { _ = plan.Close() })
				var run func(context.Context) error
				var closeSession func() error
				if debug {
					session, createErr := plan.NewDebugSession(t.Context())
					if createErr != nil {
						t.Fatal(createErr)
					}

					closeSession = session.Close
					run = func(c context.Context) error {
						event, runErr := session.Start(c)
						if runErr != nil {
							if event != nil {
								t.Errorf("aborted start returned event: %+v", event)
							}

							return runErr
						}

						event, runErr = session.Continue(c)
						if runErr == nil && (event.Reason != DebugReasonCompleted || string(event.Output.Content) != "42") {
							t.Fatalf("retry did not complete: %+v", event)
						}

						return runErr
					}
				} else {
					session, createErr := plan.NewSession(t.Context())
					if createErr != nil {
						t.Fatal(createErr)
					}

					closeSession = session.Close
					run = func(c context.Context) error {
						_, runErr := session.Run(c)

						return runErr
					}
				}

				t.Cleanup(func() { _ = closeSession() })
				if mode == "pre_canceled" {
					cancel()
				}

				err = run(ctx)
				if !errors.Is(err, wantErr) || executions != 0 {
					t.Fatalf("abort error=%v executions=%d", err, executions)
				}

				wantCalls := []string{"before", "after-second", "after-first"}
				if mode == "before_failure" {
					wantCalls = []string{"before"}
				} else if mode == "pre_canceled" {
					wantCalls = nil
				}

				if !reflect.DeepEqual(calls, wantCalls) {
					t.Fatalf("hook calls=%v, want %v", calls, wantCalls)
				}

				if len(afterErrors) > 0 {
					wantValue := "hook"
					if mode == "nil" {
						wantValue = "caller"
					}

					for i, seen := range afterErrors {
						if !errors.Is(seen, wantErr) || errors.Is(seen, afterErr) || afterContexts[i] != wantValue {
							t.Fatalf("after hook %d: error=%v context=%v", i, seen, afterContexts[i])
						}
					}

					if !errors.Is(err, afterErr) {
						t.Fatalf("lost after-hook failure: %v", err)
					}
				} else if errors.Is(err, afterErr) {
					t.Fatalf("unexpected after-hook error: %v", err)
				}

				abort = false
				calls = nil
				if err := run(caller); err != nil {
					t.Fatalf("retry failed: %v", err)
				}

				if executions != 1 || !reflect.DeepEqual(calls, []string{"before", "after-second", "after-first"}) {
					t.Fatalf("retry executions=%d hooks=%v", executions, calls)
				}

				if err := closeSession(); err != nil {
					t.Fatal(err)
				}

				if len(calls) != 3 {
					t.Fatalf("close repeated hooks: %v", calls)
				}
			})
		}
	}
}
