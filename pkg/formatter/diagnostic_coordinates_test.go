package formatter_test

import (
	"io"
	"testing"

	"github.com/antlr4-go/antlr/v4"

	"github.com/MontFerret/ferret/v2/pkg/compiler"
	"github.com/MontFerret/ferret/v2/pkg/diagnostics"
	"github.com/MontFerret/ferret/v2/pkg/formatter"
	"github.com/MontFerret/ferret/v2/pkg/parser"
	parserd "github.com/MontFerret/ferret/v2/pkg/parser/diagnostics"
	"github.com/MontFerret/ferret/v2/pkg/source"
)

func TestParserDiagnosticCoordinates(t *testing.T) {
	f, err := formatter.New()
	if err != nil {
		t.Fatal(err)
	}

	c, err := compiler.New()
	if err != nil {
		t.Fatal(err)
	}

	for _, test := range []struct {
		input, selected, label string
		runeStart, runeEnd     int
		byteStart, byteEnd     int
	}{
		{"// 😀\nFOR 1 WHILE TRUE RETURN 1", "1", "invalid binding", 9, 10, 12, 13},
		{"// é\nFOR 123 WHILE TRUE RETURN 1", "123", "invalid binding", 9, 12, 10, 13},
		{"// 😀\nRETURN QUERY `.x` IN @doc USING \"css\" /* ok */", `"css"`, "invalid dialect", 37, 42, 40, 45},
		{`RETURN WAITFOR EVENT "😀" IN @obs TIMEOUT 5s EVERY 250ms`, "EVERY", "unsupported clause", 44, 49, 47, 52},
		{"// é\nRETURN WAITFOR EVENT \"message\" @source", "@", "missing IN", 36, 37, 37, 38},
		{"FOR 1 WHILE TRUE RETURN 1", "1", "invalid binding", 4, 5, 4, 5},
		{"FOR 123 DO WHILE TRUE RETURN 1", "123", "invalid binding", 4, 7, 4, 7},
		{"FOR é WHILE TRUE RETURN 1", "é", "invalid binding", 4, 5, 4, 6},
		{"return query `.x` in @doc using 'css'", "'css'", "invalid dialect", 32, 37, 32, 37},
		{`RETURN WAITFOR EVENT "x" IN @obs TIMEOUT 5s EVERY 250ms`, "EVERY", "unsupported clause", 44, 49, 44, 49},
		{`RETURN WAITFOR EVENT "message" @source`, "@", "missing IN", 31, 32, 31, 32},
		{"// ſ\r\nreturn query `.x` in @doc using \"é css\"/* USING */", `"é css"`, "invalid dialect", 38, 45, 39, 47},
		{"// 😀\r\nRETURN WAITFOR EVENT \"é\" @source", "@", "missing IN", 31, 32, 35, 36},
		{"// é\nRETURN QUERY `.x` IN @doc USING `css` /* ok */", "`css`", "invalid dialect", 37, 42, 38, 43},
		{"// é\nRETURN QUERY `.x` IN @doc USING 'café", "'café", "invalid dialect", 37, 42, 38, 44},
		{"RETURN QUERY `.x` IN @doc USING `a ${f(`中`)} b` /* ok */", "`a ${f(`中`)} b`", "invalid dialect", 32, 47, 32, 49},
		{"// é😀\r\nRETURN", "", "missing return value", 13, 13, 17, 17},
		{"", "", "", 0, 0, 0, 0},
	} {
		t.Run(test.input, func(t *testing.T) {
			src := source.New("coordinates.fql", test.input)
			if !src.Empty() {
				t.Run("parser", func(t *testing.T) {
					history := parserd.NewTokenHistory(64)
					handler := parserd.NewErrorHandler(src, 10)
					p := parser.New(test.input, func(stream antlr.TokenStream) antlr.TokenStream {
						return parserd.NewTrackingTokenStream(stream, history)
					})
					p.RemoveErrorListeners()
					p.AddErrorListener(parserd.NewErrorListener(src, handler, history))
					p.Program()

					if handler.Errors().Size() != 1 {
						t.Fatalf("got %d parser diagnostics, want one", handler.Errors().Size())
					}
					spans := handler.Errors().First().Spans
					want := source.Span{Start: test.runeStart, End: test.runeEnd}
					if len(spans) != 1 || spans[0] != diagnostics.NewMainErrorSpan(want, test.label) {
						t.Fatalf("parser annotations = %+v, want %q at rune span %+v", spans, test.label, want)
					}
				})
			}

			formatErr := f.Format(io.Discard, src)
			_, compileErr := c.Compile(t.Context(), src)
			analysis, analyzeErr := c.Analyze(src)
			for _, result := range []struct {
				name string
				ds   []*diagnostics.Diagnostic
			}{
				{"Format", collectDiagnostics(formatErr)},
				{"Compile", collectDiagnostics(compileErr)},
				{"Analyze", collectDiagnostics(analyzeErr)},
				{"Analysis.Diagnostics", analysis.Diagnostics()},
			} {
				t.Run(result.name, func(t *testing.T) {
					if len(result.ds) != 1 {
						t.Fatalf("got %d diagnostics, want one", len(result.ds))
					}
					d := result.ds[0]
					if d.Kind != parserd.SyntaxError || d.Source.ID() != src.ID() {
						t.Fatalf("unexpected diagnostic kind or source: %+v", d)
					}
					if src.Empty() {
						if len(d.Spans) != 0 {
							t.Fatalf("empty source acquired annotations: %+v", d.Spans)
						}
						return
					}
					if len(d.Spans) != 1 || !d.Spans[0].Main || d.Spans[0].Label != test.label {
						t.Fatalf("%q: unexpected annotations %+v", d.Message, d.Spans)
					}
					span := d.Spans[0].Span
					if span.Start < 0 || span.Start > span.End || span.End > len(test.input) {
						t.Fatalf("invalid byte span %+v", span)
					}
					got := test.input[span.Start:span.End]
					if span.Start != test.byteStart || span.End != test.byteEnd || got != test.selected {
						t.Fatalf("got [%d,%d) %q, want [%d,%d) %q", span.Start, span.End, got,
							test.byteStart, test.byteEnd, test.selected)
					}
				})
			}
		})
	}
}
