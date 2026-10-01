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
	for _, query := range []string{
		"[1, 2, 3", "[1, 2,", "[\n  \"é🙂\",", "[[1, [2,",
		"[FAIL() ON ERROR RETURN NONE", "[FAIL() ON ERROR RETRY 3 OR RETURN NONE",
		"[FAIL() ON /* recovery */ ERROR RETURN 'é🙂'", "[[FAIL() ON ERROR RETURN NONE",
		"[FAIL() ON ERROR FAIL", "[FAIL() ON ERROR RETRY 3",
		"[FAIL() ON ERROR RETRY 3 DELAY 1MS BACKOFF CONSTANT OR FAIL",
		"[FAIL() ON ERROR RETURN FAIL() ON ERROR RETURN NONE",
		"[{ ON: 1 }.ON",
	} {
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

func TestArrayEntryValidationRecoveryReuse(t *testing.T) {
	entries := []struct {
		query string
		valid bool
	}{
		{query: "FAIL() ON ERROR RETURN NONE", valid: true},
		{query: "1", valid: true},
		{query: "FAIL() ON ERROR RETURN"},
		{query: "2", valid: true},
		{query: "FAIL() ON ERROR RETRY 3 OR RETURN NONE", valid: true},
		{query: "FAIL() ON ERROR RETRY 3 OR"},
		{query: "FAIL() ON ERROR FAIL", valid: true},
	}
	var query string
	for _, entry := range entries {
		query += entry.query + ","
	}

	lexer := fql.NewFqlLexer(antlr.NewInputStream(query))
	stream := antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel)
	stream.Fill()
	stream.Seek(1)
	originalIndex := stream.Index()
	validation := newArrayEntryValidation(stream)
	start := 0
	entryIndex := 0
	for end := 0; end < stream.Size(); end++ {
		if stream.Get(end).GetTokenType() != fql.FqlLexerComma {
			continue
		}

		cursor, valid := validation.parse(start, end)
		if valid != entries[entryIndex].valid || cursor.LA(1) != antlr.TokenEOF {
			t.Fatalf("entry %q: valid = %v, next token = %q", entries[entryIndex].query, valid, cursor.LT(1).GetText())
		}

		if stream.Index() != originalIndex {
			t.Fatal("validation moved the original parser cursor")
		}

		if validation.parser.BuildParseTrees || len(validation.parser.GetParseListeners()) != 0 {
			t.Fatal("validation retained a parse tree or recovery listener")
		}

		start = end + 1
		entryIndex++
	}

	if entryIndex != len(entries) {
		t.Fatalf("validated %d entries, want %d", entryIndex, len(entries))
	}
}

func TestArrayEntryValidationRejectsIncompleteContents(t *testing.T) {
	for _, query := range []string{
		"[1,,2", "[1 2", "[1 2, [3", "[1, 2 +", "[1, ...", "[1, (2 + 3", "[1, { value: 2",
		"[1, 2)", "[value[1", "[1, \"unterminated", "[1, `unterminated", "[1, 2 ? 3", "[1, 2 + + +",
		"[FAIL() ON", "[FAIL() ON ERROR", "[FAIL() ON ERROR RETURN", "[FAIL() ON ERROR RETRY",
		"[FAIL() ON ERROR RETRY 3 DELAY", "[FAIL() ON ERROR RETRY 3 DELAY 1MS BACKOFF",
		"[FAIL() ON ERROR RETRY 3 OR", "[FAIL() ON ERROR RETRY 3 OR RETURN",
		"[FAIL() ON ERROR FAIL OR", "[FAIL() ON, 2", "[FAIL() ON ERROR, 2",
		"[FAIL() ON ERROR RETURN, 2", "[FAIL() ON ERROR RETRY, 2",
		"[FAIL() ON ERROR RETRY 3 DELAY, 2", "[FAIL() ON ERROR RETRY 3 DELAY 1MS BACKOFF, 2",
		"[FAIL() ON ERROR RETRY 3 OR, 2", "[FAIL() ON ERROR RETRY 3 OR RETURN, 2",
		"[1 OR RETURN 2", "[1 OR RETURN 2, 3",
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
