package parser_test

import (
	"testing"

	"github.com/antlr4-go/antlr/v4"

	"github.com/MontFerret/ferret/v2/pkg/parser/fql"
)

func BenchmarkIdentifierLexing(b *testing.B) {
	input := "RETURN ENCODING::JSON_PARSE(CONCAT(foo_123, sha256, value1x, namespace2::value_))"
	b.ReportAllocs()

	for b.Loop() {
		lexer := fql.NewFqlLexer(antlr.NewInputStream(input))
		lexer.GetAllTokens()
	}
}

func TestIdentifierDigitUnderscore(t *testing.T) {
	for _, input := range []string{"base64_encode", "base64_decode", "a1_b", "value1_", "foo_123", "x1_y2_z3", "sha256"} {
		t.Run(input, func(t *testing.T) {
			tokens := fql.NewFqlLexer(antlr.NewInputStream(input)).GetAllTokens()
			if len(tokens) != 1 || tokens[0].GetTokenType() != fql.FqlLexerIdentifier || tokens[0].GetText() != input {
				t.Fatalf("identifier %q was split: %v", input, tokens)
			}

			tokens = fql.NewFqlLexer(antlr.NewInputStream(input + "::")).GetAllTokens()
			if len(tokens) != 1 || tokens[0].GetTokenType() != fql.FqlLexerNamespaceSegment {
				t.Fatalf("namespace %q was split: %v", input, tokens)
			}
		})
	}
}
