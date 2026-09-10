package uapi

import (
	"context"
	"errors"
	"io"
	"sync/atomic"
	"testing"

	"github.com/MontFerret/api"
	apidebugger "github.com/MontFerret/api/debugger"

	"github.com/MontFerret/ferret/v2/pkg/encoding"
	encodingjson "github.com/MontFerret/ferret/v2/pkg/encoding/json"
	"github.com/MontFerret/ferret/v2/pkg/engine"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
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

func TestMissingOutputCodecFailsAtEncodingAndPreservesCleanup(t *testing.T) {
	for _, mode := range []string{"run", "session", "debug session"} {
		t.Run(mode, func(t *testing.T) {
			var callbacks, runs, sessions, plans atomic.Int32
			result := &outputResource{}
			r := newTestRuntime(t,
				engine.WithFunctionsRegistrar(func(ns runtime.Namespace) {
					ns.Function().A0().Add("MAKE_OUTPUT", func(context.Context) (runtime.Value, error) {
						runs.Add(1)

						return result, nil
					})
				}),
				engine.WithSessionCloseHook(func() error {
					sessions.Add(1)

					return nil
				}),
				engine.WithPlanCloseHook(func() error {
					plans.Add(1)

					return nil
				}),
			)
			option := func(opts api.SessionOptions) error {
				callbacks.Add(1)
				err := api.WithOutputContentType("unknown/codec")(opts)
				if err != nil {
					t.Errorf("codec checked during portable callback: %v", err)
				}

				return err
			}

			src := api.NewAnonymousSource("RETURN MAKE_OUTPUT()")
			var out api.Output
			var operationErr error
			var p api.Plan
			var child io.Closer
			if mode == "run" {
				out, operationErr = r.Run(t.Context(), src, option)
			} else {
				var err error
				p, err = r.CompileDebug(t.Context(), src)
				if err != nil {
					t.Fatal(err)
				}

				t.Cleanup(func() { _ = p.Close() })
				if mode == "session" {
					s, err := p.NewSession(t.Context(), option)
					if err != nil {
						t.Fatalf("session creation checked codec availability: %v", err)
					}

					child = s
					t.Cleanup(func() { _ = s.Close() })
					out, operationErr = s.Run(t.Context())
				} else {
					s, err := p.NewDebugSession(t.Context(), option)
					if err != nil {
						t.Fatalf("debug session creation checked codec availability: %v", err)
					}

					child = s
					t.Cleanup(func() { _ = s.Close() })
					if _, err := s.Start(t.Context()); err != nil {
						t.Fatal(err)
					}

					var event *apidebugger.Event
					event, operationErr = s.Continue(t.Context())
					if event == nil || event.Reason != apidebugger.ReasonCompleted || event.Output != nil {
						t.Fatalf("debug completion event=%+v err=%v", event, operationErr)
					}
				}
			}

			if !errors.Is(operationErr, encoding.ErrCodecNotFound) || out.Content != nil || out.ContentType != "" {
				t.Fatalf("output=%+v err=%v", out, operationErr)
			}

			if callbacks.Load() != 1 || runs.Load() != 1 || result.closes.Load() != 1 {
				t.Fatalf("callbacks=%d query calls=%d result closes=%d, want 1 each", callbacks.Load(), runs.Load(), result.closes.Load())
			}

			if child != nil {
				if sessions.Load() != 0 || plans.Load() != 0 {
					t.Fatalf("encoding failure closed caller-owned session/plan: %d/%d", sessions.Load(), plans.Load())
				}

				for range 2 {
					if err := child.Close(); err != nil {
						t.Fatal(err)
					}

					if err := p.Close(); err != nil {
						t.Fatal(err)
					}
				}
			}

			if sessions.Load() != 1 || plans.Load() != 1 || result.closes.Load() != 1 {
				t.Fatalf("session/plan/result closes=%d/%d/%d, want 1 each", sessions.Load(), plans.Load(), result.closes.Load())
			}
		})
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
