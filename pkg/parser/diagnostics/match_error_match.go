package diagnostics

import (
	"strings"

	"github.com/antlr4-go/antlr/v4"

	"github.com/MontFerret/ferret/v2/pkg/diagnostics"
	"github.com/MontFerret/ferret/v2/pkg/parser/fql"
	"github.com/MontFerret/ferret/v2/pkg/source"
)

const (
	missingMatchArmCommaMessage = "Expected ',' between MATCH arms"
	missingMatchArmCommaHint    = "Separate MATCH arms with commas, e.g. 0 => 0, 1 => 1, _ => 0."
)

func matchMatchArmSeparatorErrors(src *source.Source, err *diagnostics.Diagnostic, offending *TokenNode) bool {
	if !(isNoAlternative(err.Message) || isMissing(err.Message) || isMismatched(err.Message)) {
		return false
	}

	if src == nil || offending == nil || offending.Token() == nil {
		return false
	}

	tokens := lexDefaultTokens(src.Content())
	offendingIdx := findLexedTokenIndex(tokens, offending.Token())
	if offendingIdx < 0 {
		return false
	}

	for arrowIdx := 0; arrowIdx < len(tokens); arrowIdx++ {
		if !isTokenText(tokens[arrowIdx], "=>") || !isInsideMatchArms(tokens, arrowIdx) {
			continue
		}

		nextArrowIdx, ok := findNextTopLevelArrowBeforeSeparator(tokens, arrowIdx+1)
		if !ok || nextArrowIdx <= arrowIdx+1 {
			continue
		}

		if offendingIdx < arrowIdx || offendingIdx > nextArrowIdx {
			continue
		}

		span, ok := matchArmSeparatorInsertionSpan(tokens, arrowIdx, nextArrowIdx, src)
		if !ok {
			continue
		}

		err.Message = missingMatchArmCommaMessage
		err.Hint = missingMatchArmCommaHint
		err.Spans = []diagnostics.ErrorSpan{
			diagnostics.NewMainErrorSpan(span, "missing comma"),
		}

		return true
	}

	return false
}

func lexDefaultTokens(input string) []antlr.Token {
	stream := antlr.NewInputStream(input)
	lexer := fql.NewFqlLexer(stream)
	lexer.RemoveErrorListeners()
	tokens := make([]antlr.Token, 0)

	for {
		token := lexer.NextToken()
		if token == nil || token.GetTokenType() == antlr.TokenEOF {
			break
		}

		if token.GetChannel() == antlr.TokenDefaultChannel {
			tokens = append(tokens, token)
		}
	}

	return tokens
}

func findLexedTokenIndex(tokens []antlr.Token, token antlr.Token) int {
	if token == nil {
		return -1
	}

	for idx, candidate := range tokens {
		if candidate.GetStart() == token.GetStart() && candidate.GetStop() == token.GetStop() {
			return idx
		}
	}

	return -1
}

func isInsideMatchArms(tokens []antlr.Token, idx int) bool {
	openIdx := findEnclosingOpenParen(tokens, idx)
	if openIdx < 0 {
		return false
	}

	depth := 0
	for i := openIdx - 1; i >= 0; i-- {
		text := tokenText(tokens[i])

		switch text {
		case ")", "]", "}":
			depth++
		case "(", "[", "{":
			if depth > 0 {
				depth--
			}
		case "MATCH":
			if depth == 0 {
				return true
			}
		case "FUNC", "RETURN", "LET", "VAR", "FOR", "USE":
			if depth == 0 {
				return false
			}
		}
	}

	return false
}

func findEnclosingOpenParen(tokens []antlr.Token, idx int) int {
	depth := 0

	for i := idx - 1; i >= 0; i-- {
		text := tokenText(tokens[i])

		switch text {
		case ")":
			depth++
		case "(":
			if depth == 0 {
				return i
			}

			depth--
		}
	}

	return -1
}

func findNextTopLevelArrowBeforeSeparator(tokens []antlr.Token, startIdx int) (int, bool) {
	depth := 0

	for i := startIdx; i < len(tokens); i++ {
		text := tokenText(tokens[i])

		switch text {
		case "(", "[", "{":
			depth++
		case ")", "]", "}":
			if depth == 0 {
				return -1, false
			}

			depth--
		case ",":
			if depth == 0 {
				return -1, false
			}
		case "=>":
			if depth == 0 {
				return i, true
			}
		}
	}

	return -1, false
}

func matchArmSeparatorInsertionSpan(tokens []antlr.Token, currentArrowIdx, nextArrowIdx int, src *source.Source) (source.Span, bool) {
	nextArmStartIdx := matchArmStartTokenIndex(tokens, currentArrowIdx, nextArrowIdx)
	if nextArmStartIdx <= currentArrowIdx+1 {
		return source.Span{}, false
	}

	previousValue := tokens[nextArmStartIdx-1]
	if previousValue == nil {
		return source.Span{}, false
	}

	offset := previousValue.GetStop() + 1
	if offset < 0 || offset >= len(src.Content()) {
		return spanFromTokenSafe(previousValue, src), true
	}

	return source.Span{Start: offset, End: offset + 1}, true
}

func matchArmStartTokenIndex(tokens []antlr.Token, currentArrowIdx, nextArrowIdx int) int {
	depth := 0

	for i := nextArrowIdx - 1; i > currentArrowIdx; i-- {
		text := tokenText(tokens[i])

		switch text {
		case ")", "]", "}":
			depth++
		case "(", "[", "{":
			if depth > 0 {
				depth--
			}
		case "WHEN":
			if depth == 0 {
				return i
			}
		}
	}

	return nextArrowIdx - 1
}

func isTokenText(token antlr.Token, expected string) bool {
	return tokenText(token) == expected
}

func tokenText(token antlr.Token) string {
	if token == nil {
		return ""
	}

	return strings.ToUpper(token.GetText())
}
