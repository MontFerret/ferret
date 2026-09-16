package uapi

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/MontFerret/api"
	apidebugger "github.com/MontFerret/api/debugger"
)

func TestUnicodeDebuggerLocationsAddressOriginalBytes(t *testing.T) {
	for _, prefix := range []string{"abc", "é", "中", "😀", "😀😀", "é中😀", "e\u0301", "\t", "\xff", "\xf0\x9f", "é\x80😀\xff"} {
		for _, separator := range []string{" ", "\n", "\r\n"} {
			t.Run(prefix+"/"+separator, func(t *testing.T) {
				runtime := newTestRuntime(t)
				query := "LET text = \"" + prefix + "\"" + separator + "RETURN text"
				src := api.NewSource("cell:unicode", query)

				plan, err := runtime.CompileDebug(t.Context(), src)
				if err != nil {
					t.Fatal(err)
				}

				t.Cleanup(func() { _ = plan.Close() })

				session, err := plan.NewDebugSession(t.Context())
				if err != nil {
					t.Fatal(err)
				}

				t.Cleanup(func() { _ = session.Close() })

				start := strings.Index(query, "RETURN")
				position := api.Position{
					Line:   strings.Count(query[:start], "\n") + 1,
					Column: start - strings.LastIndex(query[:start], "\n"),
				}
				want := api.Range{
					Location: api.Location{SourceName: src.Name, Position: position},
					Span:     api.Span{Start: start, End: len(query)},
				}

				breakpoint, err := session.SetBreakpointAt(t.Context(), want.Location, apidebugger.BreakpointOptions{BindingMode: apidebugger.BreakpointBindExact})
				if err != nil || !breakpoint.Bound || breakpoint.Location != want {
					t.Fatalf("breakpoint = %+v, error = %v, want %+v", breakpoint, err, want)
				}

				if _, err := session.Start(t.Context()); err != nil {
					t.Fatal(err)
				}

				event, err := session.Continue(t.Context())
				if err != nil || event.Reason != apidebugger.ReasonBreakpoint || event.Location != want {
					t.Fatalf("stop = %+v, error = %v, want %+v", event, err, want)
				}

				frames, err := session.Frames(t.Context())
				if err != nil || len(frames) != 1 || frames[0].Location != want.Location {
					t.Fatalf("frames = %+v, error = %v, want %+v", frames, err, want)
				}
			})
		}
	}
}

func TestMalformedUTF8DebuggerPreservesPerByteDecoding(t *testing.T) {
	for _, test := range []struct{ input, want string }{
		{"\xff", "\ufffd"},
		{"\xf0\x9f", "\ufffd\ufffd"},
		{"é\x80😀\xff", "é\ufffd😀\ufffd"},
	} {
		t.Run(test.input, func(t *testing.T) {
			runtime := newTestRuntime(t)

			plan, err := runtime.CompileDebug(t.Context(), api.NewSource("cell:malformed", "RETURN \""+test.input+"\""))
			if err != nil {
				t.Fatal(err)
			}

			t.Cleanup(func() { _ = plan.Close() })

			session, err := plan.NewDebugSession(t.Context())
			if err != nil {
				t.Fatal(err)
			}

			t.Cleanup(func() { _ = session.Close() })

			if _, err := session.Start(t.Context()); err != nil {
				t.Fatal(err)
			}

			event, err := session.Continue(t.Context())
			if err != nil || event.Reason != apidebugger.ReasonCompleted || event.Output == nil {
				t.Fatalf("completion = %+v, error = %v", event, err)
			}

			var got string
			if err := json.Unmarshal(event.Output.Content, &got); err != nil {
				t.Fatal(err)
			}

			if got != test.want {
				t.Fatalf("decoded = %q, want %q", got, test.want)
			}
		})
	}
}
