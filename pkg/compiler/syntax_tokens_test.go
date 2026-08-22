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

	analysis, err := mustNewCompiler(t).Analyze(source.NewAnonymous(query))
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

	analysis, err := mustNewCompiler(t).Analyze(source.NewAnonymous(query))
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

	analysis, err := mustNewCompiler(t).Analyze(source.NewAnonymous(query))
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
	analysis, err := mustNewCompiler(t).Analyze(source.NewAnonymous("RETURN 1"))
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

func TestSyntaxWordsReturnCanonicalCategorizedDefensiveMetadata(t *testing.T) {
	words := SyntaxWords()
	if got, want := len(words), int(SyntaxWordWith); got != want {
		t.Fatalf("SyntaxWords length = %d, want %d", got, want)
	}

	seen := make(map[SyntaxWord]SyntaxWordInfo, len(words))
	previous := ""

	for _, word := range words {
		if word.Word == SyntaxWordUnknown {
			t.Fatalf("SyntaxWords contains unknown word: %+v", word)
		}

		if _, ok := seen[word.Word]; ok {
			t.Fatalf("SyntaxWords contains duplicate identity %d", word.Word)
		}

		if word.Spelling != strings.ToUpper(word.Spelling) {
			t.Fatalf("SyntaxWords spelling %q is not canonical uppercase", word.Spelling)
		}

		if previous != "" && previous >= word.Spelling {
			t.Fatalf("SyntaxWords are not in spelling order: %q before %q", previous, word.Spelling)
		}

		seen[word.Word] = word
		previous = word.Spelling
	}

	for word := SyntaxWord(1); word <= SyntaxWordWith; word++ {
		if _, ok := seen[word]; !ok {
			t.Errorf("SyntaxWords omits identity %d", word)
		}
	}

	wantCategories := map[SyntaxWord]SyntaxWordCategory{
		SyntaxWordLet:   SyntaxWordCategoryKeyword,
		SyntaxWordAnd:   SyntaxWordCategoryOperator,
		SyntaxWordTrue:  SyntaxWordCategoryLiteral,
		SyntaxWordRetry: SyntaxWordCategoryContextual,
	}

	for word, want := range wantCategories {
		if got := seen[word].Category; got != want {
			t.Errorf("SyntaxWords category for %q = %d, want %d", seen[word].Spelling, got, want)
		}
	}

	words[0] = SyntaxWordInfo{}
	if next := SyntaxWords(); len(next) == 0 || next[0] == (SyntaxWordInfo{}) {
		t.Fatal("SyntaxWords exposed mutable metadata storage")
	}
}

func TestAnalyzeSyntaxTokensExposeCanonicalWordIdentity(t *testing.T) {
	query := "let value = true\nRETURN value AND false ASC DESC && ! ON ERROR FAIL RETRY DELAY"
	analysis, _ := mustNewCompiler(t).Analyze(source.NewAnonymous(query))
	if analysis == nil {
		t.Fatal("Analyze returned no partial snapshot")
	}

	wantWords := map[string]SyntaxWord{
		"let":    SyntaxWordLet,
		"true":   SyntaxWordTrue,
		"RETURN": SyntaxWordReturn,
		"AND":    SyntaxWordAnd,
		"false":  SyntaxWordFalse,
		"ASC":    SyntaxWordAsc,
		"DESC":   SyntaxWordDesc,
		"ON":     SyntaxWordOn,
		"ERROR":  SyntaxWordError,
		"FAIL":   SyntaxWordFail,
		"RETRY":  SyntaxWordRetry,
		"DELAY":  SyntaxWordDelay,
	}

	tokens := analysis.SyntaxTokens()
	for text, want := range wantWords {
		token, ok := syntaxTokenWithText(tokens, query, text)
		if !ok || token.Word != want {
			t.Errorf("syntax word %q = %+v, %t, want identity %d", text, token, ok, want)
		}
	}

	for _, text := range []string{"&&", "!"} {
		token, ok := syntaxTokenWithText(tokens, query, text)
		if !ok || token.Word != SyntaxWordUnknown {
			t.Errorf("symbolic operator %q word = %+v, %t", text, token, ok)
		}
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
