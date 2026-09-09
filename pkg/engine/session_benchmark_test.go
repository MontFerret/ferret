package engine

import (
	"context"
	"errors"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/debugger"
	"github.com/MontFerret/ferret/v2/pkg/logging"
	"github.com/MontFerret/ferret/v2/pkg/source"
)

func BenchmarkNewSessionOptions(b *testing.B) {
	for _, configured := range []bool{false, true} {
		name := "default"
		var setters []SessionOption

		if configured {
			name = "configured"
			setters = []SessionOption{
				WithSessionParam("value", 42),
				WithOutputContentType("application/json"),
				WithSessionLogLevel(logging.InfoLevel),
				WithDebugFormat(debugger.FormatOptions{MaxDepth: 4, MaxItems: 12, MaxBytes: 2048}),
			}
		}

		b.Run(name, func(b *testing.B) {
			b.ReportAllocs()

			for b.Loop() {
				if _, err := newSessionConfig(setters); err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}

func BenchmarkSessionRun(b *testing.B) {
	for _, mode := range []string{"no_hook", "after_success", "after_failure"} {
		b.Run(mode, func(b *testing.B) {
			var hookErr error
			var setters []Option

			if mode == "after_failure" {
				hookErr = errors.New("after run failed")
			}

			if mode != "no_hook" {
				setters = append(setters, WithAfterRunHook(func(context.Context, error) error { return hookErr }))
			}

			engine, err := New(setters...)
			if err != nil {
				b.Fatal(err)
			}

			defer func() { _ = engine.Close() }()
			plan, err := engine.Compile(b.Context(), source.NewAnonymous("RETURN @value + 1"))
			if err != nil {
				b.Fatal(err)
			}

			defer func() { _ = plan.Close() }()
			session, err := plan.NewSession(b.Context(), WithSessionParam("value", 41))
			if err != nil {
				b.Fatal(err)
			}

			defer func() { _ = session.Close() }()
			b.ReportAllocs()
			b.ResetTimer()

			for b.Loop() {
				output, err := session.Run(b.Context())
				if !errors.Is(err, hookErr) {
					b.Fatal(err)
				}

				if hookErr == nil && (output == nil || string(output.Content) != "42") {
					b.Fatalf("unexpected output: %+v", output)
				}
			}
		})
	}
}
