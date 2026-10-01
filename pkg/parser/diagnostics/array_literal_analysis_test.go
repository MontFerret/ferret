package diagnostics

import (
	"testing"
	"unicode/utf8"

	"github.com/antlr4-go/antlr/v4"

	"github.com/MontFerret/ferret/v2/pkg/diagnostics"
	"github.com/MontFerret/ferret/v2/pkg/parser/fql"
	"github.com/MontFerret/ferret/v2/pkg/source"
)

func TestArrayLiteralAnalysisRecoveryAnchors(t *testing.T) {
	for _, query := range []string{"[1, 2, 3", "[1, 2,", "[\n  \"é🙂\",", "[[1, [2,"} {
		t.Run(query, func(t *testing.T) {
			lexer := fql.NewFqlLexer(antlr.NewInputStream(query))
			stream := antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel)
			stream.Fill()
			stream.Seek(1)
			originalIndex := stream.Index()
			ctx := fql.NewExpressionAtomContext(nil, nil, -1)
			ctx.SetStart(stream.Get(0))
			analysis := newArrayLiteralAnalysis(stream)
			src := source.New("array.fql", query)
			for _, anchor := range []int{0, 1, stream.Size() - 1} {
				diagnostic := &diagnostics.Diagnostic{
					Kind: SyntaxError, Message: "no viable alternative at input 'unrelated recovery text'", Source: src,
					Spans: []diagnostics.ErrorSpan{diagnostics.NewMainErrorSpan(SpanFromToken(stream.Get(anchor)), "")},
				}
				if !analysis.match(src, diagnostic, ctx) {
					t.Fatal("expected unclosed array independent of recovery anchor")
				}

				offset := utf8.RuneCountInString(query)
				want := diagnostics.NewMainErrorSpan(source.Span{Start: offset, End: offset}, "expected ']'")
				if len(diagnostic.Spans) != 1 || diagnostic.Spans[0] != want {
					t.Fatalf("parser spans = %+v, want ANTLR character span %+v", diagnostic.Spans, want)
				}

				if stream.Index() != originalIndex {
					t.Fatal("diagnostic validation moved the original parser cursor")
				}
			}
		})
	}
}

func TestArrayEntryValidationRejectsIncompleteContents(t *testing.T) {
	for _, query := range []string{
		"[1,,2", "[1 2", "[1 2, [3", "[1, 2 +", "[1, ...", "[1, (2 + 3", "[1, { value: 2",
		"[1, 2)", "[value[1", "[1, \"unterminated", "[1, `unterminated", "[1, 2 ? 3", "[1, 2 + + +",
	} {
		t.Run(query, func(t *testing.T) {
			lexer := fql.NewFqlLexer(antlr.NewInputStream(query))
			lexer.RemoveErrorListeners()
			stream := antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel)
			stream.Fill()
			analysis := newArrayLiteralAnalysis(stream)
			if got := analysis.analyze(0); got.message == "Unclosed array literal" {
				t.Fatalf("incomplete contents misclassified: %+v", got)
			}
		})
	}
}
