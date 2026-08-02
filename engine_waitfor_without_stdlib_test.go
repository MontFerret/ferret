package ferret

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/source"
)

func TestTimedWaitForRunsWithoutStdlib(t *testing.T) {
	engine := mustNewEngine(t, WithoutStdlib())
	t.Cleanup(func() { _ = engine.Close() })

	tests := []struct {
		name  string
		query string
		want  string
	}{
		{
			name:  "duration timeout",
			query: `LET ready = false RETURN WAITFOR ready TIMEOUT 1ms EVERY 0ms`,
			want:  "false",
		},
		{
			name:  "computed timeout",
			query: `LET ready = false LET timeout = 0.5ms RETURN WAITFOR ready TIMEOUT timeout * 2 EVERY 0ms`,
			want:  "false",
		},
		{
			name:  "timeout fallback",
			query: `RETURN WAITFOR VALUE "ready" WHEN false TIMEOUT 1ms EVERY 0ms ON TIMEOUT RETURN "timeout"`,
			want:  `"timeout"`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			output, err := engine.Run(context.Background(), source.NewAnonymous(tc.query))
			if err != nil {
				t.Fatalf("run failed: %v", err)
			}
			if got := string(output.Content); got != tc.want {
				t.Fatalf("unexpected output: got %q, want %q", got, tc.want)
			}
		})
	}
}

func TestTimedWaitForWithoutStdlibHonorsCancellation(t *testing.T) {
	engine := mustNewEngine(t, WithoutStdlib())
	t.Cleanup(func() { _ = engine.Close() })

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := engine.Run(
		ctx,
		source.NewAnonymous(`LET ready = false RETURN WAITFOR ready TIMEOUT 10s EVERY 10s`),
	)
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("expected cancellation, got %v", err)
	}
}

func TestTimedWaitForWithoutStdlibSupportsRepeatedSessionRuns(t *testing.T) {
	engine := mustNewEngine(t, WithoutStdlib())
	t.Cleanup(func() { _ = engine.Close() })

	plan := mustCompilePlan(t, engine, `LET ready = false RETURN WAITFOR ready TIMEOUT 0.5ms EVERY 0ms`)
	t.Cleanup(func() { _ = plan.Close() })

	session := mustNewSession(t, plan)
	t.Cleanup(func() { _ = session.Close() })

	for run := 0; run < 2; run++ {
		output, err := session.Run(context.Background())
		if err != nil {
			t.Fatalf("run %d failed: %v", run+1, err)
		}
		if got := string(output.Content); got != "false" {
			t.Fatalf("run %d returned %q, want false", run+1, got)
		}
	}
}

func TestExplicitNowStillRequiresStdlib(t *testing.T) {
	engine := mustNewEngine(t, WithoutStdlib())
	t.Cleanup(func() { _ = engine.Close() })

	_, err := engine.Run(context.Background(), source.NewAnonymous(`RETURN NOW()`))
	if err == nil || !strings.Contains(err.Error(), "unresolved function") {
		t.Fatalf("expected unresolved NOW function, got %v", err)
	}
}
