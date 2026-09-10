package universal

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/MontFerret/api"
	"github.com/MontFerret/ferret/v2/pkg/engine"
)

func TestCancellationAndHookFailureRetainBothCauses(t *testing.T) {
	for _, stage := range []string{"compile", "run", "debug"} {
		t.Run(stage, func(t *testing.T) {
			ctx, cancel := context.WithCancel(t.Context())
			t.Cleanup(cancel)
			failure := errors.New("hook failure")
			option := engine.WithBeforeRunHook(func(context.Context) (context.Context, error) {
				cancel()

				return nil, failure
			})

			if stage == "compile" {
				option = engine.WithBeforeCompileHook(func(context.Context) error {
					cancel()

					return failure
				})
			}

			engine, err := engine.New(option)
			if err != nil {
				t.Fatal(err)
			}

			t.Cleanup(func() { _ = engine.Close() })
			runtime := Wrap(engine)
			t.Cleanup(func() { _ = runtime.Close() })
			src := api.NewAnonymousSource("RETURN 1")

			switch stage {
			case "compile":
				_, err = runtime.Compile(ctx, src)
			case "run":
				_, err = runtime.Run(ctx, src)
			case "debug":
				plan, compileErr := runtime.CompileDebug(ctx, src)
				if compileErr != nil {
					t.Fatal(compileErr)
				}

				t.Cleanup(func() { _ = plan.Close() })

				session, createErr := plan.NewDebugSession(ctx)
				if createErr != nil {
					t.Fatal(createErr)
				}

				t.Cleanup(func() { _ = session.Close() })
				_, err = session.Start(ctx)
			}

			if !errors.Is(err, context.Canceled) || !errors.Is(err, failure) {
				t.Fatalf("lost cancellation or hook failure: %v", err)
			}
		})
	}
}

func TestCancellationAndOptionFailureRetainBothCauses(t *testing.T) {
	runtime := newTestRuntime(t)
	ctx, cancel := context.WithCancel(t.Context())
	t.Cleanup(cancel)
	failure := errors.New("option failure")

	plan, err := runtime.Compile(ctx, api.NewAnonymousSource("RETURN 1"))
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() { _ = plan.Close() })

	session, err := plan.NewSession(ctx, func(api.SessionOptions) error {
		cancel()

		return failure
	})
	if session != nil || !errors.Is(err, context.Canceled) || !errors.Is(err, failure) {
		t.Fatalf("session=%v lost cancellation or option failure: %v", session, err)
	}
}

func TestCanceledAdmissionDoesNotApplyOptionsOrPublish(t *testing.T) {
	runtime := newTestRuntime(t)

	plan, err := runtime.CompileDebug(t.Context(), api.NewAnonymousSource("RETURN @value"))
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() { _ = plan.Close() })
	ctx, cancel := context.WithCancel(t.Context())
	cancel()
	calls := 0
	planOption := func(api.PlanOptions) error {
		calls++

		return nil
	}

	sessionOption := func(api.SessionOptions) error {
		calls++

		return nil
	}

	for _, compile := range []func(context.Context, api.Source, ...api.PlanOption) (api.Plan, error){runtime.Compile, runtime.CompileDebug} {
		created, err := compile(ctx, api.NewAnonymousSource("RETURN 1"), planOption)
		if created != nil || !errors.Is(err, context.Canceled) {
			t.Fatalf("canceled compilation: plan=%v err=%v", created, err)
		}
	}

	created, err := plan.NewSession(ctx, sessionOption)
	if created != nil || !errors.Is(err, context.Canceled) {
		t.Fatalf("canceled session: session=%v err=%v", created, err)
	}

	debug, err := plan.NewDebugSession(ctx, sessionOption)
	if debug != nil || !errors.Is(err, context.Canceled) || calls != 0 {
		t.Fatalf("canceled debug session: session=%v err=%v calls=%d", debug, err, calls)
	}
}

func TestRuntimeRunPreservesOutputAndAllCleanupErrors(t *testing.T) {
	sessionErr, planErr := errors.New("session cleanup"), errors.New("plan cleanup")
	var sessions, plans atomic.Int32

	engine, err := engine.New(
		engine.WithSessionCloseHook(func() error {
			sessions.Add(1)

			return sessionErr
		}),
		engine.WithPlanCloseHook(func() error {
			plans.Add(1)

			return planErr
		}),
	)
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() { _ = engine.Close() })
	runtime := Wrap(engine)
	t.Cleanup(func() { _ = runtime.Close() })

	output, err := runtime.Run(t.Context(), api.NewAnonymousSource("RETURN 42"))
	if string(output.Content) != "42" || !errors.Is(err, sessionErr) || !errors.Is(err, planErr) || sessions.Load() != 1 || plans.Load() != 1 {
		t.Fatalf("output=%+v err=%v cleanup=%d/%d", output, err, sessions.Load(), plans.Load())
	}
}

func TestPlanSupportsConcurrentIndependentSessions(t *testing.T) {
	runtime := newTestRuntime(t)

	plan, err := runtime.Compile(t.Context(), api.NewAnonymousSource("RETURN @value"))
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() { _ = plan.Close() })
	var workers sync.WaitGroup
	for i := range 8 {
		workers.Go(func() {
			session, err := plan.NewSession(t.Context(), api.WithParam("value", i))
			if err != nil {
				t.Error(err)

				return
			}

			output, runErr := session.Run(t.Context())
			if err := errors.Join(runErr, session.Close()); err != nil || string(output.Content) != string(rune('0'+i)) {
				t.Errorf("session %d output=%q err=%v", i, output.Content, err)
			}
		})
	}

	workers.Wait()
}

func TestSessionRootsAndSiblingCleanupRemainIndependent(t *testing.T) {
	runtime := newTestRuntime(t)

	plan, err := runtime.CompileDebug(t.Context(), api.NewSource("buffer://roots", `RETURN TO_STRING(IO::FS::READ("value.txt"))`))
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() { _ = plan.Close() })

	firstRoot, secondRoot := t.TempDir(), t.TempDir()
	for root, value := range map[string]string{firstRoot: "first", secondRoot: "second"} {
		if err := os.WriteFile(filepath.Join(root, "value.txt"), []byte(value), 0o600); err != nil {
			t.Fatal(err)
		}
	}

	first, err := plan.NewSession(t.Context(), api.WithFSRoot(firstRoot))
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() { _ = first.Close() })

	second, err := plan.NewDebugSession(t.Context(), api.WithFSRoot(secondRoot))
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() { _ = second.Close() })

	if _, err := second.Start(t.Context()); err != nil {
		t.Fatal(err)
	}

	output, err := first.Run(t.Context())
	if err != nil || string(output.Content) != `"first"` {
		t.Fatalf("first output=%q error=%v", output.Content, err)
	}

	if err := first.Close(); err != nil {
		t.Fatal(err)
	}

	event, err := second.Continue(t.Context())
	if err != nil || event == nil || event.Output == nil || string(event.Output.Content) != `"second"` {
		t.Fatalf("second event=%+v error=%v", event, err)
	}
}

func TestCancellationDuringRunHookSurvivesContextReplacement(t *testing.T) {
	ctx, cancel := context.WithCancel(t.Context())
	t.Cleanup(cancel)

	engine, err := engine.New(engine.WithBeforeRunHook(func(context.Context) (context.Context, error) {
		cancel()

		return context.Background(), nil
	}))
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() { _ = engine.Close() })
	runtime := Wrap(engine)
	t.Cleanup(func() { _ = runtime.Close() })

	output, err := runtime.Run(ctx, api.NewAnonymousSource("RETURN 42"))
	if !errors.Is(err, context.Canceled) || len(output.Content) != 0 {
		t.Fatalf("hook cancellation output=%+v error=%v", output, err)
	}
}
