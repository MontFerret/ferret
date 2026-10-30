package formatter_test

import (
	"bytes"
	"errors"
	"reflect"
	"strings"
	"testing"
	"unicode/utf8"

	"github.com/MontFerret/ferret/v2/pkg/compiler"
	"github.com/MontFerret/ferret/v2/pkg/diagnostics"
	"github.com/MontFerret/ferret/v2/pkg/formatter"
	"github.com/MontFerret/ferret/v2/pkg/source"
)

func TestFormatterUnicodeSource(t *testing.T) {
	for _, test := range []struct{ name, input, want string }{
		{"ASCII", "// ascii\nRETURN 1", "// ascii\nreturn 1"},
		{"BMP comment", "// é\nRETURN 1", "// é\nreturn 1"},
		{"supplementary comment", "// 😀\nRETURN 1", "// 😀\nreturn 1"},
		{"comment after literal", `RETURN ["😀", /* keep */ 1]`, "return [\n    \"😀\", /* keep */\n    1\n]"},
		{"trailing comment", `RETURN "é😀" /* 尾😀 */`, `return "é😀" /* 尾😀 */`},
		{"literal at EOF", `RETURN "é😀"`, `return "é😀"`},
		{"function without host", `RETURN UNAVAILABLE::HOST("é😀")`, `return UNAVAILABLE::HOST("é😀")`},
		{"CRLF", "// é😀\r\nRETURN \"😀\"", "// é😀\r\nreturn \"😀\""},
		{"multiline trivia", "/* é\n * 😀\n */\nRETURN 1\n// 尾😀\n", "/* é\n * 😀\n */\nreturn 1\n// 尾😀\n"},
		{"block gaps", "FUNC f() {\n// 😀\nLET x=\"é\"\n/* 中 */\nRETURN x\n// 尾😀\n}\nRETURN f()", "func f() {\n    // 😀\n    let x = \"é\"\n    /* 中 */\n    return x\n    // 尾😀\n}\nreturn f()"},
		{"parenthesized comments", `RETURN (/* é😀 */ "中" /* 尾😀 */)`, `return (/* é😀 */ "中" /* 尾😀 */)`},
	} {
		t.Run(test.name, func(t *testing.T) {
			f, err := formatter.New()
			if err != nil {
				t.Fatal(err)
			}

			var out bytes.Buffer
			if err := f.Format(&out, source.NewAnonymous(test.input)); err != nil {
				t.Fatal(err)
			}

			if got := out.String(); !utf8.ValidString(got) || got != test.want {
				t.Fatalf("output = %q (valid UTF-8: %v), want %q", got, utf8.ValidString(got), test.want)
			}

			var again bytes.Buffer
			if err := f.Format(&again, source.NewAnonymous(out.String())); err != nil || again.String() != out.String() {
				t.Fatalf("second format = %q, %v; want %q", again.String(), err, out.String())
			}
		})
	}
}

func TestFormatterDiagnosticByteSpans(t *testing.T) {
	f, err := formatter.New()
	if err != nil {
		t.Fatal(err)
	}

	c, err := compiler.New()
	if err != nil {
		t.Fatal(err)
	}

	for _, input := range []string{
		"", "RETURN 1 +", `RETURN "é" +`, `RETURN "😀" +`,
		"// é😀\nRETURN \"中\" +", "// é😀\r\nRETURN \"中\" +",
		`RETURN ["😀"`, "RETURN [\"é😀\"\r\n", "RETURN",
		"LET x =\nLET y =\nRETURN \"😀\" +",
	} {
		t.Run(input, func(t *testing.T) {
			src := source.New("unicode.fql", input)
			var out bytes.Buffer
			formatErr := f.Format(&out, src)
			_, compileErr := c.Compile(t.Context(), src)
			formatted := collectDiagnostics(formatErr)
			compiled := collectDiagnostics(compileErr)
			if len(formatted) == 0 || len(formatted) != len(compiled) {
				t.Fatalf("diagnostics: format=%v, compile=%v", formatErr, compileErr)
			}

			for i, got := range formatted {
				want := compiled[i]
				if got.Kind != want.Kind || got.Message != want.Message || got.Hint != want.Hint || got.Note != want.Note ||
					got.Source.ID() != src.ID() || !reflect.DeepEqual(got.Spans, want.Spans) {
					t.Fatalf("diagnostic %d: got kind=%s message=%q hint=%q note=%q source=%v spans=%+v; want %+v spans=%+v", i, got.Kind, got.Message, got.Hint, got.Note, got.Source.ID(), got.Spans, want, want.Spans)
				}

				for _, annotation := range got.Spans {
					span := annotation.Span
					if span.Start < 0 || span.End < span.Start || span.End > len(input) {
						t.Fatalf("invalid native span %+v for %q", span, input)
					}
				}
			}

			if input == `RETURN "😀" +` {
				if got := formatted[0].Spans[0].Span; got != (source.Span{Start: 14, End: 15}) {
					t.Fatalf("operator byte span = %+v, want [14,15)", got)
				}
			}

			if strings.HasPrefix(input, `RETURN [`) {
				span := formatted[0].Spans[0].Span
				end := len(strings.TrimRight(input, "\r\n"))

				if span != (source.Span{Start: end, End: end}) {
					t.Fatalf("insertion span = %+v, want [%d,%d)", span, end, end)
				}
			}

			if input == "" && len(formatted[0].Spans) != 0 {
				t.Fatal("empty source acquired a fabricated location")
			}

			if out.Len() != 0 {
				t.Fatalf("syntax failure wrote %q", out.String())
			}

			if err := f.Format(&out, source.NewAnonymous("RETURN 1")); err != nil || out.String() != "return 1" {
				t.Fatalf("reuse after syntax failure: output=%q error=%v", out.String(), err)
			}
		})
	}
}

func collectDiagnostics(err error) []*diagnostics.Diagnostic {
	if diagnostic, ok := err.(*diagnostics.Diagnostic); ok {
		return []*diagnostics.Diagnostic{diagnostic}
	}

	var set *diagnostics.DiagnosticSet
	if !errors.As(err, &set) {
		return nil
	}

	var result []*diagnostics.Diagnostic
	for _, diagnostic := range set.Errors() {
		result = append(result, diagnostic)
	}

	return result
}
