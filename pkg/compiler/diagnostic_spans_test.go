package compiler

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/diagnostics"
	parserd "github.com/MontFerret/ferret/v2/pkg/parser/diagnostics"
	"github.com/MontFerret/ferret/v2/pkg/source"
)

func TestReturnDiagnosticInsertionSpans(t *testing.T) {
	for _, test := range []struct {
		query    string
		keyword  string
		position source.Position
	}{
		{"RETURN", "RETURN", source.Position{Line: 1, Column: 7}},
		{"let arr = []\n\nreturn", "return", source.Position{Line: 3, Column: 7}},
		{"RETURN DISTINCT", "RETURN DISTINCT", source.Position{Line: 1, Column: 16}},
		{"RETURN \t\n", "RETURN", source.Position{Line: 1, Column: 7}},
		{"RETURN DISTINCT \t\n", "RETURN DISTINCT", source.Position{Line: 1, Column: 16}},
		{"LET text = \"é😀\" RETURN", "RETURN", source.Position{Line: 1, Column: 27}},
	} {
		t.Run(test.query, func(t *testing.T) {
			src := source.New("return.fql", test.query)
			d := compileAndAnalyzeDiagnostics(t, src)[0]
			message := "Expected expression after '" + test.keyword + "'"
			hint := "Did you forget to provide a value to return?"
			if test.keyword == "RETURN DISTINCT" {
				hint = "RETURN DISTINCT treats DISTINCT as a modifier. To return an identifier named DISTINCT, wrap it in parentheses, e.g. RETURN (DISTINCT)."
			}

			if d.Kind != parserd.SyntaxError || d.Message != message || d.Hint != hint {
				t.Fatalf("diagnostic = %+v, want %q and hint %q", d, message, hint)
			}

			offset := strings.LastIndex(test.query, test.keyword) + len(test.keyword)
			want := diagnostics.NewMainErrorSpan(source.Span{Start: offset, End: offset}, "missing return value")
			if len(d.Spans) != 1 || d.Spans[0] != want || d.Source.ID() != src.ID() {
				t.Fatalf("diagnostic source/spans = %v / %+v, want %v / %+v", d.Source, d.Spans, src, want)
			}

			if got := d.Source.PositionAt(want.Span); got != test.position {
				t.Fatalf("position = %+v, want %+v", got, test.position)
			}
		})
	}
}

func TestDiagnosticSpanBounds(t *testing.T) {
	queries := []string{
		"RETURN", "RETURN DISTINCT", "let arr = []\n\nreturn",
		"LET", "LET x =", "VAR x = 0\nx +=", "DELETE",
		"FOR", "FOR x", "FOR x IN", "FOR x IN 0.. RETURN x", "RETURN 1..",
		"FOR x IN [] FILTER", "FOR x IN [] LIMIT", "FOR x IN [] LIMIT ,",
		"FOR x IN [] LIMIT 1,", "FOR x IN [] LIMIT 1,2,3",
		"FOR x IN [] COLLECT", "FOR x IN [] COLLECT key = x INTO",
		"FOR x IN [] COLLECT AGGREGATE",
		"FUNC f()", "FUNC f() =>", "RETURN (1", "LET x = F(",
		"RETURN true AND", "RETURN true ?", "RETURN true ? 1 :",
		"RETURN [", "RETURN {", "RETURN {x:", "RETURN `abc", "RETURN \"abc",
		"RETURN doc[~", "RETURN [1][* RETURN", "RETURN [1][? NONE]",
		"RETURN [1 2]", "RETURN MATCH 1 { 1 => 1 2 => 2 }",
		"RETURN )", "DELETE obj", "RETURN [missing, other]",
		"LET [x, x] = [1, 2] RETURN x",
	}
	for _, query := range queries {
		for _, prefix := range []string{"", "LET prefix = \"é😀\"\n"} {
			t.Run(prefix+query, func(t *testing.T) {
				src := source.New("bounds.fql", prefix+query)
				for _, d := range compileAndAnalyzeDiagnostics(t, src) {
					if d.Source.ID() != src.ID() || d.Source.Content() != src.Content() || d.Message == "" || d.Kind == "" || len(d.Spans) == 0 {
						t.Fatalf("missing diagnostic metadata: %+v", d)
					}

					for _, annotation := range d.Spans {
						span := annotation.Span
						if span.Start < 0 || span.Start > span.End || span.End > len(d.Source.Content()) {
							t.Errorf("%q: annotation %q has invalid span %+v for %d bytes", d.Message, annotation.Label, span, d.Source.Length())
						}
					}
				}
			})
		}
	}
}

func TestMissingSyntaxInsertionAnchors(t *testing.T) {
	for _, test := range []struct {
		query  string
		before string
		label  string
	}{
		{"LET", "", "missing name"},
		{"LET = 1", "= 1", "missing name"},
		{"LET x =", "", "missing value"},
		{"VAR x = 0\nx +=\nRETURN x", "\nRETURN x", "missing value"},
		{"DELETE", "", "invalid delete target"},
		{"FOR", "", "missing variable"},
		{"FOR x", "", "missing keyword"},
		{"FOR IN [] RETURN 1", "IN [] RETURN 1", "missing variable"},
		{"FOR x IN\nRETURN x", "\nRETURN x", "missing value"},
		{"FOR x IN [] FILTER\nRETURN x", "\nRETURN x", "missing expression"},
		{"FOR x IN [] LIMIT\nRETURN x", "\nRETURN x", "missing expression"},
		{"FOR x IN []\nCOLLECT\nRETURN x", "\nRETURN x", "missing grouping or aggregation"},
		{"FOR x IN []\nCOLLECT AGGREGATE\nRETURN x", "\nRETURN x", "missing variable assignment"},
		{"RETURN 1..", "", "missing value"},
		{"FOR x IN 0.. RETURN x", " RETURN x", "missing value"},
		{"FUNC f()", "", "missing function body"},
		{"FUNC f() =>", "", "missing expression"},
		{"LET x = (1 OR\nRETURN x", "\nRETURN x", "missing expression"},
		{"LET x = 1 ?\nRETURN x", "\nRETURN x", "missing expression"},
		{"LET x = 1 ? 2 :\nRETURN x", "\nRETURN x", "missing expression"},
		{"RETURN (1", "", "missing ')'"},
		{"LET x = F(", "", "missing ')'"},
		{"LET x = F(,1)", ",1)", "missing value"},
		{"LET x = F(1", "", "missing ')'"},
		{"LET x = F(1\nRETURN x", "\nRETURN x", "missing ')'"},
		{"LET x = F([1, 2]", "", "missing ')'"},
		{"RETURN [", "", "missing value"},
		{"RETURN {", "", "missing property name"},
		{"RETURN { : 1 }", ": 1 }", "missing property name"},
		{"RETURN {x:", "", "missing value"},
		{"RETURN \"abc", "", "missing closing '\"'"},
		{"LET x = \"abc\r\nRETURN x", "\r\nRETURN x", "missing closing '\"'"},
		{"RETURN abc\"", "abc\"", "missing opening '\"'"},
		{"LET x = foo bar\"\nRETURN x", "foo bar\"\nRETURN x", "missing opening '\"'"},
		{"RETURN `abc", "", "missing closing '`'"},
		{"LET x = obj[", "", "missing ']'"},
		{"RETURN doc[~", "", "missing query literal"},
		{"RETURN doc[~ css()]", "()]", "missing query string"},
		{"RETURN [1][* RETURN", "", "missing expression"},
		{"RETURN [1][? NONE]", "]", "missing 'FILTER'"},
		{"RETURN [1 2]", " 2]", "missing comma"},
		{"RETURN MATCH 1 { 1 => 1 2 => 2 }", " 2 => 2 }", "missing comma"},
		{"FOR x IN [] LIMIT ,", "", "missing value"},
	} {
		for _, prefix := range []string{"", "LET prefix = \"é😀\"\n"} {
			t.Run(prefix+test.query, func(t *testing.T) {
				query := prefix + test.query
				d := compileAndAnalyzeDiagnostics(t, source.New("insertion.fql", query))[0]
				offset := len(query) - len(test.before)
				want := diagnostics.NewMainErrorSpan(source.Span{Start: offset, End: offset}, test.label)
				if len(d.Spans) != 1 || d.Spans[0] != want {
					t.Fatalf("%q: spans = %+v, want %+v", d.Message, d.Spans, want)
				}
			})
		}
	}
}

func TestOffendingSyntaxUsesActualRanges(t *testing.T) {
	for _, test := range []struct{ query, token string }{
		{"RETURN )", ")"},
		{"FOR x IN [] LIMIT 1,2,3", ","},
		{"FOR x IN [] LIMIT 100,200,", ","},
		{"DELETE obj", "DELETE obj"},
	} {
		t.Run(test.query, func(t *testing.T) {
			d := compileAndAnalyzeDiagnostics(t, source.New("offending.fql", test.query))[0]
			start := strings.LastIndex(test.query, test.token)
			want := source.Span{Start: start, End: start + len(test.token)}
			if len(d.Spans) != 1 || d.Spans[0].Span != want || !d.Spans[0].Main {
				t.Fatalf("%q: spans = %+v, want primary %+v", d.Message, d.Spans, want)
			}
		})
	}
}

func TestCollectNativeDiagnosticsTraversesErrorTree(t *testing.T) {
	first := &diagnostics.Diagnostic{Message: "first"}
	second := &diagnostics.Diagnostic{Message: "second"}
	third := &diagnostics.Diagnostic{Message: "third"}
	first.Cause = fmt.Errorf("cause: %w", third)
	err := errors.Join(fmt.Errorf("set: %w", diagnostics.NewDiagnosticsOf([]*diagnostics.Diagnostic{first, second})), third)
	if got := collectNativeDiagnostics(err); !reflect.DeepEqual(got, []*diagnostics.Diagnostic{first, third, second}) {
		t.Fatalf("diagnostics = %v, want first, third, second", got)
	}
}

func compileAndAnalyzeDiagnostics(t *testing.T, src source.Source) []*diagnostics.Diagnostic {
	t.Helper()
	c := mustNewCompiler(t)
	_, compileErr := c.Compile(t.Context(), src)
	_, analysisErr := c.Analyze(src)
	compiled := collectNativeDiagnostics(compileErr)
	analyzed := collectNativeDiagnostics(analysisErr)
	if compileErr == nil || analysisErr == nil || len(compiled) == 0 || len(analyzed) == 0 {
		t.Fatalf("expected native diagnostics: compile=%v analyze=%v", compileErr, analysisErr)
	}

	if !reflect.DeepEqual(compiled, analyzed) {
		t.Fatalf("Compile and Analyze disagree: compile=%+v analyze=%+v", compiled, analyzed)
	}

	return compiled
}

func collectNativeDiagnostics(err error) []*diagnostics.Diagnostic {
	var result []*diagnostics.Diagnostic
	seen := make(map[*diagnostics.Diagnostic]bool)
	var visit func(error)
	visit = func(current error) {
		if d, ok := current.(*diagnostics.Diagnostic); ok {
			if seen[d] {
				return
			}

			seen[d] = true
			result = append(result, d)
		}

		switch current := current.(type) {
		case interface{ Unwrap() []error }:
			for _, child := range current.Unwrap() {
				visit(child)
			}
		case interface{ Unwrap() error }:
			visit(current.Unwrap())
		}
	}
	visit(err)

	return result
}
