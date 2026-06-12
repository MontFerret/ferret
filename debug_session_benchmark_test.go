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

	plan, err := engine.CompileDebug(context.Background(), source.NewAnonymous(`
FOR i IN 1..100
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
