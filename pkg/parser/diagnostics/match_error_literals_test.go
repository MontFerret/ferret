package diagnostics

import (
	"testing"

	"github.com/antlr4-go/antlr/v4"

	"github.com/MontFerret/ferret/v2/pkg/diagnostics"
	"github.com/MontFerret/ferret/v2/pkg/source"
)

func TestQuotedLiteralInsertionSpanRecoveryMetadata(t *testing.T) {
	diagnostic := func(start, end int) *diagnostics.Diagnostic {
		return &diagnostics.Diagnostic{Spans: []diagnostics.ErrorSpan{
			diagnostics.NewMainErrorSpan(source.Span{Start: start, End: end}, "quote"),
		}}
	}
	token := func(start, end int) *TokenNode {
		return &TokenNode{token: antlr.NewCommonToken(&antlr.TokenSourceCharStreamPair{}, 0, antlr.TokenDefaultChannel, start, end-1)}
	}

	for _, test := range []struct {
		err         *diagnostics.Diagnostic
		node        *TokenNode
		name        string
		query       string
		wantOpening int
		wantClosing int
	}{
		{name: "no metadata", query: "abc"},
		{name: "no spans", query: "abc", err: &diagnostics.Diagnostic{}},
		{name: "nil token", query: "abc", node: &TokenNode{}},
		{name: "nil diagnostic with token", query: "abc\ndef", node: token(1, 2), wantOpening: 1, wantClosing: 3},
		{name: "no spans with token", query: "abc\ndef", err: &diagnostics.Diagnostic{}, node: token(1, 2), wantOpening: 1, wantClosing: 3},
		{name: "nil node with diagnostic", query: "abc", err: diagnostic(1, 2), wantOpening: 1, wantClosing: 3},
		{name: "nil token with diagnostic", query: "abc", err: diagnostic(1, 2), node: &TokenNode{}, wantOpening: 1, wantClosing: 3},
		{name: "diagnostic takes precedence", query: "abc", err: diagnostic(1, 2), node: token(2, 3), wantOpening: 1, wantClosing: 3},
		{name: "negative diagnostic", query: "abc", err: diagnostic(-2, -1)},
		{name: "negative start fallback", query: "abc", err: diagnostic(-1, 1), node: token(1, 2), wantOpening: 1, wantClosing: 3},
		{name: "negative end fallback", query: "abc", err: diagnostic(-2, -1), node: token(1, 2), wantOpening: 1, wantClosing: 3},
		{name: "reversed fallback", query: "abc", err: diagnostic(2, 1), node: token(1, 2), wantOpening: 1, wantClosing: 3},
		{name: "oversized fallback", query: "abc", err: diagnostic(1, 4), node: token(1, 2), wantOpening: 1, wantClosing: 3},
		{name: "negative token", query: "abc", node: token(-1, 0)},
		{name: "reversed token", query: "abc", node: token(2, 1)},
		{name: "oversized token", query: "abc", node: token(1, 4)},
		{name: "unmatched diagnostic span", query: "abc def", err: diagnostic(1, 2), wantOpening: 1, wantClosing: 7},
		{name: "empty source", err: diagnostic(0, 0)},
		{name: "empty source invalid metadata", err: diagnostic(0, 1), node: token(0, 1)},
		{name: "Unicode rune coordinates", query: "é🙂x", err: diagnostic(1, 2), wantOpening: 1, wantClosing: 3},
		{name: "byte offset rejected", query: "é🙂x", err: diagnostic(6, 7), node: token(1, 2), wantOpening: 1, wantClosing: 3},
		{name: "byte token offset rejected", query: "é🙂x", node: token(6, 7)},
		{name: "zero width EOF", query: "é🙂x", err: diagnostic(3, 3), wantOpening: 3, wantClosing: 3},
		{name: "zero width EOF token", query: "é🙂x", node: token(3, 3), wantOpening: 3, wantClosing: 3},
		{name: "LF boundary", query: "abc\nnext", err: diagnostic(1, 2), wantOpening: 1, wantClosing: 3},
		{name: "CR boundary", query: "abc\rnext", err: diagnostic(1, 2), wantOpening: 1, wantClosing: 3},
		{name: "CRLF boundary", query: "abc\r\nnext", err: diagnostic(1, 2), wantOpening: 1, wantClosing: 3},
	} {
		t.Run(test.name, func(t *testing.T) {
			src := source.New("literal.fql", test.query)
			for _, opening := range []bool{true, false} {
				offset := test.wantClosing
				if opening {
					offset = test.wantOpening
				}

				want := source.Span{Start: offset, End: offset}
				if got := quotedLiteralInsertionSpan(src, test.err, test.node, opening); got != want {
					t.Errorf("opening=%t: insertion = %+v, want %+v", opening, got, want)
				}
			}
		})
	}
}

func TestQuotedLiteralInsertionSpanOpeningAnchor(t *testing.T) {
	for _, test := range []struct {
		query string
		want  int
	}{
		{query: "RETURN foo bar\"", want: 7},
		{query: "foo\nbar\"", want: 4},
		{query: "foo\r\nbar\"", want: 5},
		{query: "LET x = 'é🙂'\nRETURN foo bar\"", want: 20},
		{query: "\"", want: 0},
	} {
		t.Run(test.query, func(t *testing.T) {
			src := source.New("literal.fql", test.query)
			tokens := lexDefaultTokens(test.query)
			quote := tokens[len(tokens)-1]
			err := &diagnostics.Diagnostic{Spans: []diagnostics.ErrorSpan{
				diagnostics.NewMainErrorSpan(SpanFromToken(quote), "missing quote"),
			}}
			want := source.Span{Start: test.want, End: test.want}
			if got := quotedLiteralInsertionSpan(src, err, &TokenNode{token: quote}, true); got != want {
				t.Fatalf("insertion = %+v, want %+v", got, want)
			}
		})
	}
}
