package universal

import (
	"context"
	"errors"
	"testing"

	"github.com/MontFerret/api"
	apidebugger "github.com/MontFerret/api/debugger"

	encodingjson "github.com/MontFerret/ferret/v2/pkg/encoding/json"
	"github.com/MontFerret/ferret/v2/pkg/engine"
)

func TestConfiguredOutputAndCleanupFailuresRemainAvailable(t *testing.T) {
	for _, debug := range []bool{false, true} {
		t.Run(map[bool]string{false: "ordinary", true: "debug"}[debug], func(t *testing.T) {
			failure := errors.New("after run")
			r := newTestRuntime(t,
				engine.WithEncodingCodec("application/custom", encodingjson.Default),
				engine.WithAfterRunHook(func(context.Context, error) error { return failure }),
			)
			p, err := r.CompileDebug(t.Context(), api.NewAnonymousSource("RETURN @value"))
			if err != nil {
				t.Fatal(err)
			}

			t.Cleanup(func() { _ = p.Close() })
			var outputs []api.Output
			for range 2 {
				opts := []api.SessionOption{api.WithParam("value", 42), api.WithOutputContentType(" application/custom ")}
				var out api.Output
				if debug {
					s, err := p.NewDebugSession(t.Context(), opts...)
					if err != nil {
						t.Fatal(err)
					}

					if _, err := s.Start(t.Context()); err != nil {
						t.Fatal(err)
					}

					event, runErr := s.Continue(t.Context())
					if !errors.Is(runErr, failure) || event == nil || event.Reason != apidebugger.ReasonCompleted || event.Output == nil {
						t.Fatalf("debug event=%+v err=%v", event, runErr)
					}

					out = *event.Output
					_ = s.Close()
				} else {
					s, err := p.NewSession(t.Context(), opts...)
					if err != nil {
						t.Fatal(err)
					}

					var runErr error
					out, runErr = s.Run(t.Context())
					if !errors.Is(runErr, failure) {
						t.Fatalf("run error=%v", runErr)
					}

					_ = s.Close()
				}

				if out.ContentType != "application/custom" || string(out.Content) != "42" {
					t.Fatalf("output=%+v", out)
				}

				outputs = append(outputs, out)
			}

			_ = p.Close()
			_ = r.Close()
			outputs[0].Content[0] = '9'
			if string(outputs[1].Content) != "42" {
				t.Fatal("output aliases another execution")
			}
		})
	}
}

func TestMissingOutputCodecFailsWithoutOutput(t *testing.T) {
	r := newTestRuntime(t)
	out, err := r.Run(t.Context(), api.NewAnonymousSource("RETURN 42"), api.WithOutputContentType("unknown/codec"))
	if err == nil || out.Content != nil || out.ContentType != "" {
		t.Fatalf("output=%+v err=%v", out, err)
	}
}

func TestOutputValueDereferencesWithoutCopyingBytes(t *testing.T) {
	if out := outputValue(nil); out.Content != nil || out.ContentType != "" {
		t.Fatalf("nil Native output=%+v", out)
	}

	native := &api.Output{Content: []byte("42"), ContentType: "application/json"}
	out := outputValue(native)
	if out.ContentType != native.ContentType || string(out.Content) != "42" || &out.Content[0] != &native.Content[0] {
		t.Fatalf("pointer-to-value bridge changed output: %+v", out)
	}
}
