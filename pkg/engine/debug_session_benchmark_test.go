package engine

import (
	"context"
	"strconv"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/debugger"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/source"
)

// Time the complete public operation, including metadata and stop-scoped publication.
func BenchmarkDebugSessionVariablePage(b *testing.B) {
	engine, err := New()
	if err != nil {
		b.Fatal(err)
	}

	defer engine.Close()

	for _, kind := range []string{"array", "dataset"} {
		for _, size := range []int{100, 10_000, 1_000_000} {
			b.Run(kind+"/"+strconv.Itoa(size), func(b *testing.B) {
				query := "LET results = @items\nRETURN results"
				var option SessionOption
				if kind == "array" {
					values := make([]runtime.Value, size)

					for i := range values {
						values[i] = runtime.Int(42)
					}

					option = WithSessionRuntimeParam("items", runtime.NewArrayOf(values))
				} else {
					query = "LET results = (FOR n IN 0..@last RETURN 42)\nRETURN results"
					option = WithSessionRuntimeParam("last", runtime.Int(size-1))
				}

				plan, err := engine.CompileDebug(b.Context(), source.NewAnonymous(query))
				if err != nil {
					b.Fatal(err)
				}

				defer plan.Close()
				session, err := plan.NewDebugSession(b.Context(), option)
				if err != nil {
					b.Fatal(err)
				}

				defer session.Close()
				if _, err := session.SetBreakpoint(b.Context(), source.Location{Position: source.Position{Line: 2}}); err != nil {
					b.Fatal(err)
				}

				if _, err := session.Start(b.Context()); err != nil {
					b.Fatal(err)
				}

				if stop, err := session.Continue(b.Context()); err != nil || stop.Reason != debugger.ReasonBreakpoint {
					b.Fatalf("stop: %+v, %v", stop, err)
				}

				root, err := session.Evaluate(b.Context(), "results")
				if err != nil || root.IndexedCount == nil || *root.IndexedCount != int64(size) {
					b.Fatalf("root: %+v, %v", root, err)
				}

				page := debugger.VariablePage{Start: int64(size - 20), Count: 20, Filter: debugger.VariableFilterIndexed}
				b.ReportAllocs()

				for b.Loop() {
					got, err := session.Variables(b.Context(), root.Reference, page)
					if err != nil || len(got) != 20 {
						b.Fatalf("page length %d: %v", len(got), err)
					}
				}
			})
		}
	}
}

func BenchmarkDebugSessionRepeatedSourcePoints(b *testing.B) {
	engine, err := New()
	if err != nil {
		b.Fatal(err)
	}
	defer engine.Close()

	plan, err := engine.CompileDebug(context.Background(), source.NewAnonymous(`
RETURN FOR i IN 1..100
  RETURN i
`))
	if err != nil {
		b.Fatal(err)
	}
	defer plan.Close()

	b.ReportAllocs()
	b.ResetTimer()

	for b.Loop() {
		session, sessionErr := plan.NewDebugSession(context.Background())
		if sessionErr != nil {
			b.Fatal(sessionErr)
		}

		event, startErr := session.Start(context.Background())
		if startErr != nil {
			b.Fatal(startErr)
		}

		for event.Reason != debugger.ReasonCompleted {
			event, startErr = session.StepIn(context.Background())
			if startErr != nil {
				b.Fatal(startErr)
			}
		}

		if closeErr := session.Close(); closeErr != nil {
			b.Fatal(closeErr)
		}
	}
}

func BenchmarkDebugSessionBreakpointChecks(b *testing.B) {
	engine, err := New()
	if err != nil {
		b.Fatal(err)
	}
	defer engine.Close()

	plan, err := engine.CompileDebug(context.Background(), source.New("checks.fql", `FUNC unused() {
  RETURN 0
}
RETURN FOR i IN 1..100
  RETURN i`))
	if err != nil {
		b.Fatal(err)
	}
	defer plan.Close()

	for _, tc := range []struct {
		name string
		line int
	}{
		{name: "Empty"},
		{name: "Populated", line: 2},
		{name: "FrequentHits", line: 5},
	} {
		b.Run(tc.name, func(b *testing.B) {
			b.ReportAllocs()

			for b.Loop() {
				session, err := plan.NewDebugSession(context.Background())
				if err != nil {
					b.Fatal(err)
				}

				if tc.line != 0 {
					breakpoint, err := session.SetBreakpoint(context.Background(), source.Location{Position: source.Position{Line: tc.line}})
					if err != nil || !breakpoint.Bound {
						b.Fatalf("set breakpoint: %#v, %v", breakpoint, err)
					}
				}

				event, err := session.Start(context.Background())
				if err != nil {
					b.Fatal(err)
				}

				for event.Reason != debugger.ReasonCompleted {
					event, err = session.Continue(context.Background())
					if err != nil {
						b.Fatal(err)
					}
				}

				if err := session.Close(); err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}

func BenchmarkDebugSessionPausedCallerFrameInspection(b *testing.B) {
	engine, err := New()
	if err != nil {
		b.Fatal(err)
	}
	defer engine.Close()

	plan, err := engine.CompileDebug(context.Background(), source.New("caller.fql", `LET base = 1
FUNC outer(p) {
  LET carried = base + p
  FUNC inner(q) {
    RETURN carried + q
  }
  LET result = inner(3)
  RETURN result
}
RETURN outer(2)`))
	if err != nil {
		b.Fatal(err)
	}
	defer plan.Close()

	session, err := plan.NewDebugSession(context.Background())
	if err != nil {
		b.Fatal(err)
	}
	defer session.Close()

	if _, err := session.SetBreakpoint(context.Background(), source.Location{SourceName: "caller.fql", Position: source.Position{Line: 5}}); err != nil {
		b.Fatal(err)
	}
	if _, err := session.Start(context.Background()); err != nil {
		b.Fatal(err)
	}
	if _, err := session.Continue(context.Background()); err != nil {
		b.Fatal(err)
	}

	b.ReportAllocs()
	b.ResetTimer()

	for b.Loop() {
		locals, localsErr := session.FrameLocals(context.Background(), 2)
		if localsErr != nil {
			b.Fatal(localsErr)
		}
		if len(locals) == 0 {
			b.Fatal("caller frame has no locals")
		}
	}
}
