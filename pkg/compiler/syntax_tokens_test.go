package compiler

import (
	"strings"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/source"
)

func TestAnalyzeSyntaxTokensCoverPublicKindsAndTriviaContract(t *testing.T) {
	query := `// café
USE WEB::HTML AS html
LET value = "text"
RETURN [value + 42, 2s, @input]`

	analysis, err := New().Analyze(source.NewAnonymous(query))
	if err != nil {
		t.Fatal(err)
	}

	tokens := analysis.SyntaxTokens()
	wantKinds := []SyntaxTokenKind{
		SyntaxTokenKindIdentifier,
		SyntaxTokenKindNamespace,
		SyntaxTokenKindKeyword,
		SyntaxTokenKindString,
		SyntaxTokenKindNumber,
		SyntaxTokenKindDuration,
		SyntaxTokenKindComment,
		SyntaxTokenKindOperator,
		SyntaxTokenKindPunctuation,
	}

	for _, want := range wantKinds {
		if !hasSyntaxTokenKind(tokens, want) {
			t.Errorf("tokens do not contain kind %v: %+v", want, tokens)
		}
	}

	for i, token := range tokens {
		if token.Span.Start < 0 || token.Span.End <= token.Span.Start || token.Span.End > len(query) {
			t.Fatalf("token %d has invalid span %+v", i, token)
		}

		if i > 0 && tokens[i-1].Span.Start > token.Span.Start {
			t.Fatalf("tokens are not source ordered at %d: %+v", i, tokens)
		}

		text := query[token.Span.Start:token.Span.End]
		if strings.TrimSpace(text) == "" {
			t.Fatalf("whitespace trivia escaped into syntax tokens: %+v", token)
		}
	}

	valueOffset := strings.Index(query, "value")
	if valueOffset < 0 {
		t.Fatal("value marker not found")
	}

	if token, ok := syntaxTokenWithText(tokens, query, "value"); !ok || token.Span.Start != valueOffset {
		t.Fatalf("value token = %+v, %t, want UTF-8 byte offset %d", token, ok, valueOffset)
	}
}

func TestAnalyzeSyntaxTokensIncludeMalformedAndUnknownInput(t *testing.T) {
	query := "RETURN [1, §"

	analysis, err := New().Analyze(source.NewAnonymous(query))
	if err == nil {
		t.Fatal("Analyze succeeded for malformed input")
	}

	if analysis == nil {
		t.Fatal("Analyze returned no partial snapshot")
	}

	tokens := analysis.SyntaxTokens()
	unknown, ok := syntaxTokenWithText(tokens, query, "§")
	if !ok || unknown.Kind != SyntaxTokenKindUnknown {
		t.Fatalf("unknown token = %+v, %t", unknown, ok)
	}

	if got, want := unknown.Span.Start, strings.Index(query, "§"); got != want {
		t.Fatalf("unknown byte start = %d, want %d", got, want)
	}
}

func TestAnalyzeSyntaxTokensClassifyTemplateTextAndDelimitersAsStrings(t *testing.T) {
	query := "LET name = 'Ferret'\nRETURN `hello ${name}`"

	analysis, err := New().Analyze(source.NewAnonymous(query))
	if err != nil {
		t.Fatal(err)
	}

	tokens := analysis.SyntaxTokens()
	for _, text := range []string{"`", "hello "} {
		token, ok := syntaxTokenWithText(tokens, query, text)
		if !ok || token.Kind != SyntaxTokenKindString {
			t.Errorf("template token %q = %+v, %t", text, token, ok)
		}
	}

	open, ok := syntaxTokenWithText(tokens, query, "${")
	if !ok || open.Kind != SyntaxTokenKindPunctuation {
		t.Errorf("template expression open = %+v, %t", open, ok)
	}

	close, ok := syntaxTokenWithText(tokens, query, "}")
	if !ok || close.Kind != SyntaxTokenKindPunctuation {
		t.Errorf("template expression close = %+v, %t", close, ok)
	}
}

func TestAnalysisSyntaxTokensReturnsDefensiveCopy(t *testing.T) {
	analysis, err := New().Analyze(source.NewAnonymous("RETURN 1"))
	if err != nil {
		t.Fatal(err)
	}

	first := analysis.SyntaxTokens()
	if len(first) == 0 {
		t.Fatal("Analyze returned no syntax tokens")
	}

	first[0] = SyntaxToken{}
	second := analysis.SyntaxTokens()
	if second[0] == (SyntaxToken{}) {
		t.Fatal("SyntaxTokens exposed mutable analysis storage")
	}
}

func hasSyntaxTokenKind(tokens []SyntaxToken, kind SyntaxTokenKind) bool {
	for _, token := range tokens {
		if token.Kind == kind {
			return true
		}
	}

	return false
}

func syntaxTokenWithText(tokens []SyntaxToken, query, text string) (SyntaxToken, bool) {
	for _, token := range tokens {
		if query[token.Span.Start:token.Span.End] == text {
			return token, true
		}
	}

	return SyntaxToken{}, false
}
