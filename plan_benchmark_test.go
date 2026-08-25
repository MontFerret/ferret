package ferret

import (
	"context"
	"testing"
)

func BenchmarkPlanNewSession(b *testing.B) {
	engine, err := New(WithMaxIdleVMsPerPlan(1), WithMaxVMsPerPlan(1))
	if err != nil {
		b.Fatal(err)
	}
	defer func() {
		_ = engine.Close()
	}()

	plan, err := engine.Compile(context.Background(), newAnonymousAPIFile("RETURN 1"))
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

	plan, err := engine.CompileDebug(context.Background(), newAnonymousAPIFile("RETURN 1"))
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
