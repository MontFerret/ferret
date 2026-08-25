package ferret

import (
	"context"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/source"
)

func BenchmarkDebugSessionRepeatedSourcePoints(b *testing.B) {
	engine, err := New()
	if err != nil {
		b.Fatal(err)
	}
	defer engine.Close()

	compiled, err := engine.CompileDebug(context.Background(), newAnonymousAPIFile(`
RETURN FOR i IN 1..100
  RETURN i
`))
	if err != nil {
		b.Fatal(err)
	}
	plan, ok := compiled.(*Plan)
	if !ok {
		b.Fatalf("unexpected plan type %T", compiled)
	}
	defer plan.Close()

	b.ReportAllocs()
	b.ResetTimer()

	for b.Loop() {
		session, sessionErr := plan.newDebugSession(context.Background())
		if sessionErr != nil {
			b.Fatal(sessionErr)
		}

		event, startErr := session.Start(context.Background())
		if startErr != nil {
			b.Fatal(startErr)
		}

		for event.Reason != DebugReasonCompleted {
			event, startErr = session.Step(context.Background())
			if startErr != nil {
				b.Fatal(startErr)
			}
		}

		if closeErr := session.Close(); closeErr != nil {
			b.Fatal(closeErr)
		}
	}
}

func BenchmarkDebugSessionPausedCallerFrameInspection(b *testing.B) {
	engine, err := New()
	if err != nil {
		b.Fatal(err)
	}
	defer engine.Close()

	compiled, err := engine.CompileDebug(context.Background(), apiFile(source.New("caller.fql", `LET base = 1
FUNC outer(p) {
  LET carried = base + p
  FUNC inner(q) {
    RETURN carried + q
  }
  LET result = inner(3)
  RETURN result
}
RETURN outer(2)`)))
	if err != nil {
		b.Fatal(err)
	}
	plan, ok := compiled.(*Plan)
	if !ok {
		b.Fatalf("unexpected plan type %T", compiled)
	}
	defer plan.Close()

	session, err := plan.newDebugSession(context.Background())
	if err != nil {
		b.Fatal(err)
	}
	defer session.Close()

	if _, err := session.SetBreakpoint("caller.fql", 5); err != nil {
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
		locals, localsErr := session.FrameLocals(2)
		if localsErr != nil {
			b.Fatal(localsErr)
		}
		if len(locals) == 0 {
			b.Fatal("caller frame has no locals")
		}
	}
}
