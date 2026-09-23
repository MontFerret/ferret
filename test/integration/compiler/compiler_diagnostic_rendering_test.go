package compiler_test

import (
	"strings"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/diagnostics"
	parserd "github.com/MontFerret/ferret/v2/pkg/parser/diagnostics"
	"github.com/MontFerret/ferret/v2/pkg/source"
)

func TestSyntaxDiagnosticRendering(t *testing.T) {
	for _, test := range []struct {
		name     string
		query    string
		message  string
		hint     string
		label    string
		location string
		snippet  string
		span     source.Span
	}{
		{
			name:     "missing return value",
			query:    "RETURN",
			message:  "Expected expression after 'RETURN'",
			hint:     "Did you forget to provide a value to return?",
			label:    "missing return value",
			location: " --> rendering.fql:1:7\n",
			snippet:  "1 | RETURN\n  |       ^ missing return value\n",
			span:     source.Span{Start: 6, End: 6},
		},
		{
			name:     "missing distinct return value",
			query:    "RETURN DISTINCT",
			message:  "Expected expression after 'RETURN DISTINCT'",
			hint:     "RETURN DISTINCT treats DISTINCT as a modifier. To return an identifier named DISTINCT, wrap it in parentheses, e.g. RETURN (DISTINCT).",
			label:    "missing return value",
			location: " --> rendering.fql:1:16\n",
			snippet:  "1 | RETURN DISTINCT\n  |                ^ missing return value\n",
			span:     source.Span{Start: 15, End: 15},
		},
		{
			name:     "multiline missing return value",
			query:    "let arr = []\n\nreturn",
			message:  "Expected expression after 'return'",
			hint:     "Did you forget to provide a value to return?",
			label:    "missing return value",
			location: " --> rendering.fql:3:7\n",
			snippet:  "3 | return\n  |       ^ missing return value\n",
			span:     source.Span{Start: 20, End: 20},
		},
		{
			name:     "missing comma",
			query:    "RETURN [1 2]",
			message:  "Expected ',' between array items",
			hint:     "Separate array items with commas, e.g. [1, 2, 3].",
			label:    "missing comma",
			location: " --> rendering.fql:1:10\n",
			snippet:  "1 | RETURN [1 2]\n  |          ^ missing comma\n",
			span:     source.Span{Start: 9, End: 9},
		},
		{
			name:     "missing closing quote",
			query:    "RETURN \"abc",
			message:  "Unclosed string literal",
			hint:     "Add a matching '\"' to close the string.",
			label:    "missing closing '\"'",
			location: " --> rendering.fql:1:12\n",
			snippet:  "1 | RETURN \"abc\n  |            ^ missing closing '\"'\n",
			span:     source.Span{Start: 11, End: 11},
		},
		{
			name:     "offending token",
			query:    "RETURN )",
			message:  "Expected expression after 'RETURN'",
			hint:     "Did you forget to provide a value to return?",
			label:    "missing return value",
			location: " --> rendering.fql:1:8\n",
			snippet:  "1 | RETURN )\n  |        ^ missing return value\n",
			span:     source.Span{Start: 7, End: 8},
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			src := source.New("rendering.fql", test.query)
			_, err := mustNewCompiler(t).Compile(t.Context(), src)
			if err == nil {
				t.Fatal("expected compilation error")
			}

			diag := firstCompilationError(err)
			if diag == nil {
				t.Fatalf("expected native diagnostic, got %T: %v", err, err)
			}

			if diag.Kind != parserd.SyntaxError || diag.Message != test.message || diag.Hint != test.hint {
				t.Fatalf("diagnostic = %+v, want SyntaxError with message %q and hint %q", diag, test.message, test.hint)
			}

			wantSpan := diagnostics.NewMainErrorSpan(test.span, test.label)
			if len(diag.Spans) != 1 || diag.Spans[0] != wantSpan || diag.Source.ID() != src.ID() {
				t.Fatalf("source/spans = %v / %+v, want %v / %+v", diag.Source, diag.Spans, src, wantSpan)
			}

			formatted := diagnostics.Format(err)
			for _, fragment := range []string{
				"SyntaxError: " + test.message + "\n",
				test.location,
				test.snippet,
				"Hint: " + test.hint + "\n",
			} {
				if !strings.Contains(formatted, fragment) {
					t.Errorf("formatted diagnostic lacks %q:\n%s", fragment, formatted)
				}
			}
		})
	}
}
