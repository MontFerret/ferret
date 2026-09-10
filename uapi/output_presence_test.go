package uapi_test

import (
	"context"
	"errors"
	"testing"

	"github.com/MontFerret/api"
	apidiagnostics "github.com/MontFerret/api/diagnostics"

	"github.com/MontFerret/ferret/v2"
	"github.com/MontFerret/ferret/v2/pkg/diagnostics"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/source"
	"github.com/MontFerret/ferret/v2/uapi"
)

func TestExecutionPreservesOutputPresence(t *testing.T) {
	for _, mode := range []string{"runtime", "session"} {
		for _, tc := range []struct {
			name        string
			query       string
			contentType string
			content     string
			wantOutput  bool
			wantError   bool
			afterRun    bool
		}{
			{name: "populated output", query: "RETURN 42", content: "42", wantOutput: true},
			{name: "empty value", query: `RETURN ""`, content: `""`, wantOutput: true},
			{name: "empty bytes", query: "RETURN 42", contentType: "application/empty", wantOutput: true},
			{name: "execution failure", query: "RETURN FAIL_OUTPUT()", wantError: true},
			{name: "output and error", query: "RETURN 42", content: "42", wantOutput: true, wantError: true, afterRun: true},
			{name: "empty bytes and error", query: "RETURN 42", contentType: "application/empty", wantOutput: true, wantError: true, afterRun: true},
		} {
			t.Run(mode+"/"+tc.name, func(t *testing.T) {
				src := api.NewSource("output.fql", tc.query)
				cause := errors.New("operation failed")
				hookDiagnostic := diagnostics.NewUnexpectedErrorWith(source.New(src.Name, src.Content), "after run failed", cause)
				opts := []ferret.Option{
					ferret.WithEncodingCodec("application/empty", emptyOutputCodec{}),
					ferret.WithFunctionsRegistrar(func(ns runtime.Namespace) {
						ns.Function().A0().Add("FAIL_OUTPUT", func(context.Context) (runtime.Value, error) {
							return runtime.None, cause
						})
					}),
				}
				if tc.afterRun {
					opts = append(opts, ferret.WithAfterRunHook(func(context.Context, error) error {
						return hookDiagnostic
					}))
				}

				portable, err := uapi.New(opts...)
				if err != nil {
					t.Fatal(err)
				}

				t.Cleanup(func() { _ = portable.Close() })
				contentType := tc.contentType
				if contentType == "" {
					contentType = "application/json"
				}

				output, err := runOutput(t, portable, mode, t.Context(), src, api.WithOutputContentType(contentType))
				if (output != nil) != tc.wantOutput {
					t.Fatalf("output=%+v err=%v, want output present=%t", output, err, tc.wantOutput)
				}

				if output != nil && (output.ContentType != contentType || string(output.Content) != tc.content) {
					t.Fatalf("output=%+v, want content type=%q content=%q", output, contentType, tc.content)
				}

				if !tc.wantError {
					if err != nil {
						t.Fatal(err)
					}

					return
				}

				var native *diagnostics.Diagnostic
				var projected apidiagnostics.Diagnostics
				if !errors.Is(err, cause) || !errors.As(err, &native) || !errors.As(err, &projected) || len(projected) != 1 {
					t.Fatalf("lost Native cause or Universal diagnostic: %v", err)
				}

				if projected[0].Source != src || projected[0].Message != native.Message {
					t.Fatalf("projected diagnostic=%+v Native diagnostic=%+v", projected[0], native)
				}

				if tc.afterRun && native != hookDiagnostic {
					t.Fatalf("lost Native hook diagnostic identity: %v", err)
				}
			})
		}
	}
}

func TestRuntimeCompilationFailureReturnsNoOutput(t *testing.T) {
	portable, err := uapi.New()
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() { _ = portable.Close() })
	var executor api.Runtime = portable
	output, err := executor.Run(t.Context(), api.NewSource("invalid.fql", "RETURN"))
	if output != nil || err == nil {
		t.Fatalf("output=%+v err=%v", output, err)
	}

	var native *diagnostics.Diagnostic
	var projected apidiagnostics.Diagnostics
	if !errors.As(err, &native) || !errors.As(err, &projected) || len(projected) == 0 {
		t.Fatalf("lost compilation diagnostics: %v", err)
	}
}

func TestCanceledExecutionReturnsNoOutput(t *testing.T) {
	for _, mode := range []string{"runtime", "session"} {
		t.Run(mode, func(t *testing.T) {
			portable, err := uapi.New()
			if err != nil {
				t.Fatal(err)
			}

			t.Cleanup(func() { _ = portable.Close() })
			ctx, cancel := context.WithCancel(t.Context())
			cancel()
			output, err := runOutput(t, portable, mode, ctx, api.NewAnonymousSource("RETURN 42"))
			if output != nil || !errors.Is(err, context.Canceled) {
				t.Fatalf("canceled execution output=%+v err=%v", output, err)
			}
		})
	}
}

func runOutput(t *testing.T, portable api.Runtime, mode string, ctx context.Context, src api.Source, opts ...api.SessionOption) (*api.Output, error) {
	t.Helper()

	if mode == "runtime" {
		return portable.Run(ctx, src, opts...)
	}

	plan, err := portable.Compile(t.Context(), src)
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() { _ = plan.Close() })
	session, err := plan.NewSession(t.Context(), opts...)
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() { _ = session.Close() })

	return session.Run(ctx)
}
