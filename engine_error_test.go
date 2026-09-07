package ferret

import (
	"context"
	"errors"
	"reflect"
	"sync"
	"testing"
	"time"

	"github.com/MontFerret/ferret/v2/pkg/compiler"
	"github.com/MontFerret/ferret/v2/pkg/diagnostics"
	"github.com/MontFerret/ferret/v2/pkg/source"
)

func TestEnginePreservesDiagnosticFormatting(t *testing.T) {
	for _, query := range []string{"RETURN missing", "LET a = missing\nLET b = other\nRETURN 1"} {
		for _, debug := range []bool{false, true} {
			t.Run(query+map[bool]string{false: "/normal", true: "/debug"}[debug], func(t *testing.T) {
				var observed error
				engine := mustNewEngine(t, WithAfterCompileHook(func(_ context.Context, err error) error {
					observed = err

					return nil
				}))
				t.Cleanup(func() { _ = engine.Close() })
				c, err := compiler.New()
				if err != nil {
					t.Fatal(err)
				}

				src := source.New("diagnostics.fql", query)
				_, directErr := c.Compile(t.Context(), src)
				if directErr == nil {
					t.Fatal("expected compiler diagnostic")
				}

				if query != "RETURN missing" {
					var set *diagnostics.DiagnosticSet
					if !errors.As(directErr, &set) || set.Size() != 2 {
						t.Fatalf("expected two diagnostics, got %v", directErr)
					}
				}

				compile := engine.Compile
				if debug {
					compile = engine.CompileDebug
				}

				_, err = compile(t.Context(), src)
				if err == nil {
					t.Fatal("expected compiler diagnostic")
				}

				if err != observed || FormatError(err) != FormatError(directErr) {
					t.Fatalf("diagnostic changed: hook=%T returned=%T\nwant:\n%s\ngot:\n%s", observed, err, FormatError(directErr), FormatError(err))
				}
			})
		}
	}
}

func TestEngineRunPreservesRuntimeDiagnosticFormatting(t *testing.T) {
	engine := mustNewEngine(t)
	t.Cleanup(func() { _ = engine.Close() })
	for _, query := range []string{"RETURN 1 / @zero", "RETURN [@left, @right]"} {
		t.Run(query, func(t *testing.T) {
			src := source.New("runtime.fql", query)
			plan, err := engine.Compile(t.Context(), src)
			if err != nil {
				t.Fatal(err)
			}

			t.Cleanup(func() { _ = plan.Close() })
			session, err := plan.NewSession(t.Context(), WithSessionParam("zero", 0))
			if err != nil {
				t.Fatal(err)
			}

			t.Cleanup(func() { _ = session.Close() })
			_, directErr := session.Run(t.Context())
			_, runErr := engine.Run(t.Context(), src, WithSessionParam("zero", 0))
			if directErr == nil || runErr == nil {
				t.Fatalf("expected runtime failures: direct=%v engine=%v", directErr, runErr)
			}

			if FormatError(runErr) != FormatError(directErr) {
				t.Fatalf("want:\n%s\ngot:\n%s", FormatError(directErr), FormatError(runErr))
			}
		})
	}
}

func TestEngineCompilationPreservesDiagnosticAndAdditionalFailures(t *testing.T) {
	ctx, cancel := context.WithTimeout(t.Context(), 10*time.Second)
	t.Cleanup(cancel)
	entered, release := make(chan struct{}), make(chan struct{})
	releaseHook := sync.OnceFunc(func() { close(release) })
	hookErr := errors.New("after compile failed")
	engine := mustNewEngine(t, WithAfterCompileHook(func(context.Context, error) error {
		close(entered)
		<-release

		return hookErr
	}))
	t.Cleanup(func() { releaseHook(); _ = engine.Close() })
	errorsCh := make(chan error, 1)
	go func() {
		_, err := engine.Compile(ctx, NewAnonymousSource("RETURN missing"))
		errorsCh <- err
	}()
	select {
	case <-entered:
	case <-ctx.Done():
		t.Fatal(ctx.Err())
	}

	cancel()
	releaseHook()
	err := <-errorsCh
	var diagnostic *diagnostics.Diagnostic
	if !errors.As(err, &diagnostic) || !errors.Is(err, hookErr) || !errors.Is(err, context.Canceled) {
		t.Fatalf("lost failure cause: %v", err)
	}

	if FormatError(err) != err.Error() {
		t.Fatal("rendering of joined failures changed")
	}
}

func TestEngineCompilationDoesNotRepeatCancellation(t *testing.T) {
	for _, debug := range []bool{false, true} {
		t.Run(map[bool]string{false: "normal", true: "debug"}[debug], func(t *testing.T) {
			ctx, cancel := context.WithCancel(t.Context())
			t.Cleanup(cancel)
			engine := mustNewEngine(t, WithBeforeCompileHook(func(context.Context) error {
				cancel()

				return nil
			}))
			t.Cleanup(func() { _ = engine.Close() })
			compile := engine.Compile
			if debug {
				compile = engine.CompileDebug
			}

			_, err := compile(ctx, NewAnonymousSource("RETURN 1"))
			if err != context.Canceled {
				t.Fatalf("cancellation was rewrapped or duplicated: %v", err)
			}
		})
	}
}

func TestEngineRunPreservesOutputAndAllCleanupFailures(t *testing.T) {
	for _, query := range []string{"RETURN 42", "RETURN 1 / @zero"} {
		t.Run(query, func(t *testing.T) {
			sessionErr, planErr := errors.New("session cleanup"), errors.New("plan cleanup")
			var closed []string
			engine := mustNewEngine(t,
				WithSessionCloseHook(func() error {
					closed = append(closed, "session")

					return sessionErr
				}),
				WithPlanCloseHook(func() error {
					closed = append(closed, "plan")

					return planErr
				}),
			)
			t.Cleanup(func() { _ = engine.Close() })
			output, err := engine.Run(t.Context(), NewAnonymousSource(query), WithSessionParam("zero", 0))
			if !errors.Is(err, sessionErr) || !errors.Is(err, planErr) {
				t.Fatalf("lost cleanup cause: %v", err)
			}

			if query == "RETURN 42" {
				if output == nil || string(output.Content) != "42" {
					t.Fatalf("lost successful output: %+v", output)
				}
			} else {
				var diagnostic *diagnostics.Diagnostic
				if output != nil || !errors.As(err, &diagnostic) {
					t.Fatalf("lost runtime diagnostic: output=%+v err=%v", output, err)
				}
			}

			if !reflect.DeepEqual(closed, []string{"session", "plan"}) {
				t.Fatalf("cleanup order before teardown: %v", closed)
			}
		})
	}
}
