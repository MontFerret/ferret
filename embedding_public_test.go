package ferret_test

import (
	"context"
	"errors"
	"io/fs"
	"reflect"
	"strings"
	"testing"

	"github.com/MontFerret/ferret/v2"
	"github.com/MontFerret/ferret/v2/pkg/diagnostics"
)

func TestPublicEngineRunPreservesOptionsOutputAndCleanup(t *testing.T) {
	hookErr := errors.New("after run failed")
	var closed []string

	eng, err := ferret.New(
		nil,
		ferret.WithParam("value", 1),
		ferret.WithParams(map[string]any{"value": 2, "base": 10}),
		ferret.WithAfterRunHook(func(context.Context, error) error { return hookErr }),
		ferret.WithSessionCloseHook(func() error {
			closed = append(closed, "session")

			return nil
		}),
		ferret.WithPlanCloseHook(func() error {
			closed = append(closed, "plan")

			return nil
		}),
	)
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		if err := eng.Close(); err != nil {
			t.Errorf("close engine: %v", err)
		}
	})

	output, err := eng.Run(t.Context(), ferret.NewAnonymousSource("RETURN @base + @value"),
		ferret.WithSessionParam("value", 3),
		nil,
		ferret.WithSessionParams(map[string]any{"value": 4}),
	)
	if !errors.Is(err, hookErr) {
		t.Fatalf("lost hook error: %v", err)
	}

	if output == nil || output.ContentType != "application/json" || string(output.Content) != "14" {
		t.Fatalf("unexpected output: %+v", output)
	}

	if !reflect.DeepEqual(closed, []string{"session", "plan"}) {
		t.Fatalf("cleanup before caller teardown = %v, want [session plan]", closed)
	}
}

func TestPublicConstructorPreservesErrorIdentity(t *testing.T) {
	hookErr := errors.New("initialization failed")
	eng, err := ferret.New(ferret.WithEngineInitHook(func() error { return hookErr }))
	if eng != nil {
		_ = eng.Close()
		t.Fatal("failed constructor returned an engine")
	}

	if !errors.Is(err, hookErr) {
		t.Fatalf("lost initialization error: %v", err)
	}
}

func TestPublicEngineCloseErrorsAppearOnce(t *testing.T) {
	for _, rollback := range []bool{false, true} {
		t.Run(map[bool]string{false: "shutdown", true: "rollback"}[rollback], func(t *testing.T) {
			hookErr := &fs.PathError{Op: "flush", Path: "hook", Err: errors.New("failed")}
			initErr := errors.New("initialization failed")
			closes := 0
			eng, err := ferret.New(
				ferret.WithEngineInitHook(func() error {
					if rollback {
						return initErr
					}

					return nil
				}),
				ferret.WithEngineCloseHook(func() error {
					closes++

					return hookErr
				}),
			)

			wantError := "close hooks: flush hook: failed"

			if rollback {
				if eng != nil {
					_ = eng.Close()

					t.Fatal("failed constructor returned an engine")
				}

				if !errors.Is(err, initErr) {
					t.Fatalf("lost initialization error: %v", err)
				}

				wantError = "init hooks: initialization failed\nclose engine: " + wantError
			} else {
				if err != nil {
					t.Fatal(err)
				}

				err = eng.Close()
				if repeated := eng.Close(); repeated != err {
					t.Fatalf("repeated close changed the error: %v", repeated)
				}
			}

			var typed *fs.PathError
			if !errors.Is(err, hookErr) || !errors.As(err, &typed) || typed != hookErr {
				t.Fatalf("lost hook error identity: %v", err)
			}

			if err.Error() != wantError {
				t.Errorf("cleanup error = %q, want %q", err.Error(), wantError)
			}

			if closes != 1 {
				t.Fatalf("close hook calls = %d, want one", closes)
			}
		})
	}
}

func TestPublicDiagnosticFormatting(t *testing.T) {
	var observed error
	eng, err := ferret.New(ferret.WithAfterCompileHook(func(_ context.Context, err error) error {
		observed = err

		return nil
	}))
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		if err := eng.Close(); err != nil {
			t.Errorf("close engine: %v", err)
		}
	})

	_, err = eng.Compile(t.Context(), ferret.NewSource("public.fql", "RETURN missing"))
	if err == nil || err != observed {
		t.Fatalf("diagnostic identity changed: returned=%v observed=%v", err, observed)
	}

	var diagnostic *diagnostics.Diagnostic
	if !errors.As(err, &diagnostic) {
		t.Fatalf("lost compilation diagnostic: %T", err)
	}

	// Formatting remains usable as a function value after removing reassignment.
	format := ferret.FormatError
	formatted := format(err)
	if formatted != diagnostics.Format(err) || !strings.Contains(formatted, "public.fql:1") || !strings.Contains(formatted, "RETURN missing") {
		t.Fatalf("unexpected diagnostic formatting:\n%s", formatted)
	}

	plainErr := errors.New("plain error")
	if got := format(plainErr); got != plainErr.Error() {
		t.Fatalf("plain error formatting = %q", got)
	}

	ctx, cancel := context.WithCancel(t.Context())
	cancel()

	_, err = eng.Compile(ctx, ferret.NewAnonymousSource("RETURN 1"))
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("lost cancellation identity: %v", err)
	}
}

func TestPublicDebugSession(t *testing.T) {
	eng, err := ferret.New(ferret.WithOptimizationLevel(ferret.OptimizationBasic))
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		if err := eng.Close(); err != nil {
			t.Errorf("close engine: %v", err)
		}
	})

	plan, err := eng.CompileDebug(t.Context(), ferret.NewSource("public-debug.fql", "LET value = @input\nRETURN value + 1"))
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		if err := plan.Close(); err != nil {
			t.Errorf("close plan: %v", err)
		}
	})

	session, err := plan.NewDebugSession(t.Context(),
		ferret.WithSessionParam("input", 41),
		ferret.WithDebugFormat(ferret.DebugFormatOptions{MaxDepth: 3, MaxItems: 10, MaxBytes: 256}),
	)
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		if err := session.Close(); err != nil {
			t.Errorf("close debug session: %v", err)
		}
	})

	breakpoint, err := session.SetBreakpoint(ferret.DebugSourceLocation{
		SourceName: "public-debug.fql",
		Position:   ferret.Position{Line: 2},
	})
	if err != nil {
		t.Fatal(err)
	}

	if !breakpoint.Bound || breakpoint.Location.Line != 2 {
		t.Fatalf("unexpected breakpoint: %+v", breakpoint)
	}

	event, err := session.Start(t.Context())
	if err != nil {
		t.Fatal(err)
	}

	if event == nil || event.Reason != ferret.DebugReasonEntry {
		t.Fatalf("unexpected entry: %+v", event)
	}

	event, err = session.Continue(t.Context())
	if err != nil {
		t.Fatal(err)
	}

	if event == nil || event.Reason != ferret.DebugReasonBreakpoint || event.Location.Line != 2 ||
		len(event.HitBreakpointIDs) != 1 || event.HitBreakpointIDs[0] != breakpoint.ID {
		t.Fatalf("unexpected breakpoint stop: %+v", event)
	}

	value, err := session.Evaluate(t.Context(), "value")
	if err != nil {
		t.Fatal(err)
	}

	if value.Display != "41" {
		t.Fatalf("unexpected paused value: %+v", value)
	}

	event, err = session.Continue(t.Context())
	if err != nil {
		t.Fatal(err)
	}

	if event == nil || event.Reason != ferret.DebugReasonCompleted || event.Error != nil ||
		event.Output == nil || string(event.Output.Content) != "42" {
		t.Fatalf("unexpected completion: %+v", event)
	}
}
