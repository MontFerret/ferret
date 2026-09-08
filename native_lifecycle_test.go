package ferret_test

import (
	"context"
	"errors"
	"io"
	"sync"
	"sync/atomic"
	"testing"
	"testing/synctest"
	"time"

	"github.com/MontFerret/api"
	"github.com/MontFerret/ferret/v2"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/vm"
)

func TestNativeOperationsAfterClose(t *testing.T) {
	engine, plan := newNativeLifecyclePlan(t)
	if err := plan.Close(); err != nil {
		t.Fatal(err)
	}

	option := func(api.SessionOptions) error {
		t.Error("closed plan applied a session option")

		return nil
	}

	for _, debug := range []bool{false, true} {
		session, err := createNativeLifecycleSession(plan, t.Context(), debug, option)
		if session != nil || !errors.Is(err, runtime.ErrInvalidOperation) {
			t.Fatalf("closed plan: session=%v err=%v", session, err)
		}
	}

	if data, err := plan.Marshal(); data != nil || !errors.Is(err, runtime.ErrInvalidOperation) {
		t.Fatalf("closed marshal: data=%v err=%v", data, err)
	}

	if err := engine.Close(); err != nil {
		t.Fatal(err)
	}

	for _, compile := range []func(context.Context, ferret.Source, ...ferret.PlanOption) (*ferret.Plan, error){engine.Compile, engine.CompileDebug} {
		created, err := compile(t.Context(), ferret.NewAnonymousSource("RETURN 1"), func(api.PlanOptions) error {
			t.Error("closed engine applied a plan option")

			return nil
		})
		if created != nil || !errors.Is(err, runtime.ErrInvalidOperation) {
			t.Fatalf("closed compile: plan=%v err=%v", created, err)
		}
	}

	if loaded, err := engine.Load(nil); loaded != nil || !errors.Is(err, runtime.ErrInvalidOperation) {
		t.Fatalf("closed load: plan=%v err=%v", loaded, err)
	}
}

func TestNativeEngineCloseDoesNotWaitForCompile(t *testing.T) {
	for _, debug := range []bool{false, true} {
		t.Run(map[bool]string{false: "ordinary", true: "debug"}[debug], func(t *testing.T) {
			synctest.Test(t, func(t *testing.T) {
				entered, release := make(chan struct{}), make(chan struct{})
				releaseHook := sync.OnceFunc(func() { close(release) })
				defer releaseHook()
				var engineCloses, planCloses atomic.Int32
				engine, err := ferret.New(
					ferret.WithBeforeCompileHook(func(context.Context) error {
						close(entered)
						<-release

						return nil
					}),
					ferret.WithEngineCloseHook(func() error {
						engineCloses.Add(1)

						return nil
					}),
					ferret.WithPlanCloseHook(func() error {
						planCloses.Add(1)

						return nil
					}),
				)
				if err != nil {
					t.Fatal(err)
				}

				defer func() { _ = engine.Close() }()
				compile := engine.Compile
				if debug {
					compile = engine.CompileDebug
				}

				var plan *ferret.Plan
				var compileErr error
				go func() { plan, compileErr = compile(t.Context(), ferret.NewAnonymousSource("RETURN @value")) }()
				<-entered

				if err := engine.Close(); err != nil {
					t.Fatal(err)
				}

				if engineCloses.Load() != 1 || planCloses.Load() != 0 {
					t.Fatal("engine close did not limit cleanup to its own resources")
				}

				releaseHook()
				synctest.Wait()
				if compileErr != nil || plan == nil {
					t.Fatalf("admitted compile: plan=%v err=%v", plan, compileErr)
				}

				if params := plan.Params(); len(params) != 1 || params[0] != "value" {
					t.Fatalf("invalid returned plan: params=%v", params)
				}

				if err := plan.Close(); err != nil || planCloses.Load() != 1 {
					t.Fatalf("caller plan cleanup: err=%v calls=%d", err, planCloses.Load())
				}
			})
		})
	}
}

func TestNativeCapacityWaitObservesPlanCloseAndContext(t *testing.T) {
	for _, debug := range []bool{false, true} {
		for _, parent := range []string{"plan", "engine"} {
			t.Run(map[bool]string{false: "ordinary/", true: "debug/"}[debug]+parent, func(t *testing.T) {
				synctest.Test(t, func(t *testing.T) {
					closing, release := make(chan struct{}), make(chan struct{})
					releaseHook := sync.OnceFunc(func() { close(release) })
					closeErr := errors.New("plan cleanup failed")
					engine, plan := newNativeLifecyclePlan(t, ferret.WithMaxActiveSessions(1), ferret.WithPlanCloseHook(func() error {
						close(closing)
						<-release

						return closeErr
					}))
					t.Cleanup(releaseHook)
					first, err := createNativeLifecycleSession(plan, t.Context(), debug)
					if err != nil {
						t.Fatal(err)
					}

					defer func() { _ = first.Close() }()
					ctx, cancel := context.WithCancel(t.Context())
					defer cancel()
					var waiting io.Closer
					var waitErr error
					finished := false
					go func() {
						waiting, waitErr = createNativeLifecycleSession(plan, ctx, debug)
						finished = true
					}()
					synctest.Wait()
					if finished {
						t.Fatal("session creation bypassed occupied capacity")
					}

					if parent == "plan" {
						closed := make(chan error, 1)
						go func() { closed <- plan.Close() }()
						<-closing
						synctest.Wait()
						if !finished || waiting != nil || !errors.Is(waitErr, runtime.ErrInvalidOperation) {
							t.Fatalf("wait during plan cleanup: finished=%v session=%v err=%v", finished, waiting, waitErr)
						}

						releaseHook()
						if err := <-closed; !errors.Is(err, closeErr) {
							t.Fatalf("lost plan cleanup error: %v", err)
						}
					} else {
						if err := engine.Close(); err != nil {
							t.Fatal(err)
						}

						synctest.Wait()
						if finished {
							t.Fatalf("engine close interfered with plan admission: %v", waitErr)
						}

						cancel()
						synctest.Wait()
						if !finished || waiting != nil || !errors.Is(waitErr, context.Canceled) {
							t.Fatalf("canceled wait: session=%v err=%v", waiting, waitErr)
						}
					}

					runNativeLifecycleSession(t, first)
				})
			})
		}
	}
}

func TestNativePlanCloseDuringSessionConstruction(t *testing.T) {
	for _, debug := range []bool{false, true} {
		for _, capacity := range []int{0, 1} {
			for _, admitted := range []bool{false, true} {
				name := map[bool]string{false: "ordinary/", true: "debug/"}[debug] + map[int]string{0: "unlimited/", 1: "limited/"}[capacity] + map[bool]string{false: "options", true: "environment"}[admitted]
				t.Run(name, func(t *testing.T) {
					synctest.Test(t, func(t *testing.T) {
						var closes atomic.Int32
						engine, plan := newNativeLifecyclePlan(t, ferret.WithMaxActiveSessions(capacity), ferret.WithSessionCloseHook(func() error {
							closes.Add(1)

							return nil
						}))
						sibling, err := engine.CompileDebug(t.Context(), ferret.NewAnonymousSource("RETURN 1"))
						if err != nil {
							t.Fatal(err)
						}

						t.Cleanup(func() { _ = sibling.Close() })
						entered, release := make(chan struct{}), make(chan struct{})
						releaseOption := sync.OnceFunc(func() { close(release) })
						defer releaseOption()
						block := func() {
							close(entered)
							<-release
						}
						var option ferret.SessionOption = func(api.SessionOptions) error {
							block()

							return nil
						}
						if admitted {
							option = ferret.WithEnvironmentOptions(vm.WithFunctionsRegistrar(func(runtime.FunctionDefs) { block() }))
						}

						var session io.Closer
						var createErr error
						go func() { session, createErr = createNativeLifecycleSession(plan, t.Context(), debug, option) }()
						<-entered
						if err := plan.Close(); err != nil {
							t.Fatal(err)
						}

						releaseOption()
						synctest.Wait()
						if admitted && debug {
							if createErr != nil || session == nil || closes.Load() != 0 {
								t.Fatalf("admitted debug construction: session=%v err=%v closes=%d", session, createErr, closes.Load())
							}

							runNativeLifecycleSession(t, session)
							if err := session.Close(); err != nil || closes.Load() != 1 {
								t.Fatalf("caller session cleanup: err=%v calls=%d", err, closes.Load())
							}
						} else if session != nil || !errors.Is(createErr, runtime.ErrInvalidOperation) {
							t.Fatalf("closed admission/pool: session=%v err=%v", session, createErr)
						}

						retry, err := createNativeLifecycleSession(sibling, t.Context(), debug)
						if err != nil {
							t.Fatalf("sibling lost capacity: %v", err)
						}

						if err := retry.Close(); err != nil {
							t.Fatal(err)
						}
					})
				})
			}
		}
	}
}

func TestNativeCanceledEnvironmentConstructionReleasesCapacity(t *testing.T) {
	for _, debug := range []bool{false, true} {
		t.Run(map[bool]string{false: "ordinary", true: "debug"}[debug], func(t *testing.T) {
			var closes atomic.Int32
			closeErr := errors.New("session cleanup failed")
			_, plan := newNativeLifecyclePlan(t, ferret.WithMaxActiveSessions(1), ferret.WithMaxVMsPerPlan(1), ferret.WithSessionCloseHook(func() error {
				closes.Add(1)

				return closeErr
			}))
			ctx, cancel := context.WithCancel(t.Context())
			defer cancel()
			session, err := createNativeLifecycleSession(plan, ctx, debug,
				ferret.WithSessionFSRoot(t.TempDir()),
				ferret.WithEnvironmentOptions(vm.WithFunctionsRegistrar(func(runtime.FunctionDefs) { cancel() })),
			)
			if session != nil || !errors.Is(err, context.Canceled) || !errors.Is(err, closeErr) || closes.Load() != 1 {
				t.Fatalf("canceled construction: session=%v err=%v closes=%d", session, err, closes.Load())
			}

			// A deadline bounds a leaked-permit failure without controlling ordering.
			retryCtx, deadline := context.WithTimeout(t.Context(), 5*time.Second)
			defer deadline()
			retry, err := createNativeLifecycleSession(plan, retryCtx, debug)
			if err != nil {
				t.Fatalf("capacity not released: %v", err)
			}

			if err := retry.Close(); !errors.Is(err, closeErr) || closes.Load() != 2 {
				t.Fatalf("retry close: err=%v calls=%d", err, closes.Load())
			}
		})
	}
}

func TestNativeSessionCreationRacesPlanClose(t *testing.T) {
	for _, debug := range []bool{false, true} {
		t.Run(map[bool]string{false: "ordinary", true: "debug"}[debug], func(t *testing.T) {
			engine, _ := newNativeLifecyclePlan(t, ferret.WithMaxActiveSessions(1))
			ctx, cancel := context.WithTimeout(t.Context(), 5*time.Second)
			defer cancel()

			for range 50 {
				plan, err := engine.CompileDebug(t.Context(), ferret.NewAnonymousSource("RETURN 1"))
				if err != nil {
					t.Fatal(err)
				}

				start := make(chan struct{})
				var session io.Closer
				var createErr, closeErr error
				var wait sync.WaitGroup
				wait.Go(func() {
					<-start
					session, createErr = createNativeLifecycleSession(plan, ctx, debug)
				})
				wait.Go(func() {
					<-start
					closeErr = plan.Close()
				})
				close(start)
				wait.Wait()
				if closeErr != nil {
					t.Fatal(closeErr)
				}

				if createErr != nil {
					if session != nil || !errors.Is(createErr, runtime.ErrInvalidOperation) {
						t.Fatalf("creation race: session=%v err=%v", session, createErr)
					}
				} else {
					runNativeLifecycleSession(t, session)
					if err := session.Close(); err != nil {
						t.Fatal(err)
					}
				}
			}
		})
	}
}

func newNativeLifecyclePlan(t *testing.T, opts ...ferret.Option) (*ferret.Engine, *ferret.Plan) {
	t.Helper()
	engine, err := ferret.New(opts...)
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() { _ = engine.Close() })
	plan, err := engine.CompileDebug(t.Context(), ferret.NewAnonymousSource("RETURN 1"))
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() { _ = plan.Close() })

	return engine, plan
}

func createNativeLifecycleSession(plan *ferret.Plan, ctx context.Context, debug bool, opts ...ferret.SessionOption) (io.Closer, error) {
	if debug {
		session, err := plan.NewDebugSession(ctx, opts...)
		if session == nil {
			return nil, err
		}

		return session, err
	}

	session, err := plan.NewSession(ctx, opts...)
	if session == nil {
		return nil, err
	}

	return session, err
}

func runNativeLifecycleSession(t *testing.T, session io.Closer) {
	t.Helper()
	var output *ferret.Output
	var err error
	switch session := session.(type) {
	case *ferret.Session:
		output, err = session.Run(t.Context())
	case *ferret.DebugSession:
		if _, err := session.Start(t.Context()); err != nil {
			t.Fatal(err)
		}

		event, runErr := session.Continue(t.Context())
		err = runErr
		if event != nil {
			output = event.Output
			if event.Error != nil {
				err = errors.Join(err, event.Error)
			}
		}
	default:
		t.Fatalf("unexpected session: %T", session)
	}

	if err != nil || output == nil || string(output.Content) != "1" {
		t.Fatalf("returned session is unusable: output=%v err=%v", output, err)
	}
}
