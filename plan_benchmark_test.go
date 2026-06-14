package ferret

import (
	"context"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/source"
)

func BenchmarkPlanNewSession(b *testing.B) {
	engine, err := New(WithMaxIdleVMsPerPlan(1), WithMaxVMsPerPlan(1))
	if err != nil {
		b.Fatal(err)
	}
	defer func() {
		_ = engine.Close()
	}()

	plan, err := engine.Compile(context.Background(), source.NewAnonymous("RETURN 1"))
	if err != nil {
		b.Fatal(err)
	}
	defer func() {
		_ = plan.Close()
	}()

	b.ReportAllocs()
	b.ResetTimer()

	for b.Loop() {
		session, sessionErr := plan.NewSession(context.Background())
		if sessionErr != nil {
			b.Fatal(sessionErr)
		}
		if closeErr := session.Close(); closeErr != nil {
			b.Fatal(closeErr)
		}
	}
}

func BenchmarkPlanNewDebugSession(b *testing.B) {
	engine, err := New()
	if err != nil {
		b.Fatal(err)
	}
	defer func() {
		_ = engine.Close()
	}()

	plan, err := engine.CompileDebug(context.Background(), source.NewAnonymous("RETURN 1"))
	if err != nil {
		b.Fatal(err)
	}
	defer func() {
		_ = plan.Close()
	}()

	b.ReportAllocs()
	b.ResetTimer()

	for b.Loop() {
		session, sessionErr := plan.NewDebugSession(context.Background())
		if sessionErr != nil {
			b.Fatal(sessionErr)
		}
		if closeErr := session.Close(); closeErr != nil {
			b.Fatal(closeErr)
		}
	}
}
