package diagnostics

import (
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/diagnostics"
	"github.com/MontFerret/ferret/v2/pkg/source"
)

// Exercise matcher entry points that parser recovery can reach with different
// token histories. Compiler tests separately cover the currently selected routes.
func TestInsertionSpansAcrossRecoveryContexts(t *testing.T) {
	for _, test := range []struct {
		matcher SyntaxErrorMatcher
		query   string
		message string
		label   string
		before  string
		index   int
	}{
		{query: "COLLECT =", message: "no viable alternative at input 'COLLECT ='", label: "missing variable", before: "=", matcher: matchForLoopErrors, index: 0},
		{query: "COLLECT =", message: "no viable alternative at input 'COLLECT ='", label: "missing variable", before: "=", matcher: matchForLoopErrors, index: 1},
		{query: "INTO", message: "no viable alternative at input 'INTO'", label: "missing variable name", matcher: matchForLoopErrors},
		{query: "AGGREGATE", message: "no viable alternative at input 'AGGREGATE'", label: "missing variable assignment", matcher: matchForLoopErrors},
		{query: "DELETE", message: "missing Identifier at '<EOF>'", label: "invalid delete target", matcher: matchDeleteStatementErrors},
		{query: "doc[~ css`x`", message: "missing ']' at '<EOF>'", label: "missing ']'", matcher: matchQueryOperatorErrors, index: 6},
		{query: "items[* FILTER true", message: "missing ']' at '<EOF>'", label: "missing ']'", matcher: matchArrayOperatorUnclosed, index: 4},
	} {
		t.Run(test.query+"/"+test.before, func(t *testing.T) {
			src := source.New("recovery.fql", test.query)
			tokens := lexDefaultTokens(test.query)
			history := NewTokenHistory(len(tokens))
			for _, token := range tokens {
				history.Add(token)
			}

			node := history.Last().PrevAt(len(tokens) - 1 - test.index)
			err := &diagnostics.Diagnostic{Kind: SyntaxError, Source: src, Message: test.message}
			if !test.matcher(src, err, node) {
				t.Fatal("expected matcher to report missing syntax")
			}

			offset := len(test.query) - len(test.before)
			want := diagnostics.NewMainErrorSpan(source.Span{Start: offset, End: offset}, test.label)
			if len(err.Spans) != 1 || err.Spans[0] != want {
				t.Fatalf("spans = %+v, want %+v", err.Spans, want)
			}

			if err.Message == test.message || err.Hint == "" || err.Kind != SyntaxError {
				t.Fatalf("missing diagnostic metadata: %+v", err)
			}
		})
	}
}
