package diagnostics

import (
	"regexp"
	"strings"

	"github.com/antlr4-go/antlr/v4"

	"github.com/MontFerret/ferret/pkg/parser/fql"

	"github.com/MontFerret/ferret/pkg/file"
)

func SpanFromRuleContext(ctx antlr.ParserRuleContext) file.Span {
	start := ctx.GetStart()
	stop := ctx.GetStop()

	if start == nil || stop == nil {
		return file.Span{Start: 0, End: 0}
	}

	return file.Span{Start: start.GetStart(), End: stop.GetStop() + 1}
}

func SpanFromToken(tok antlr.Token) file.Span {
	if tok == nil {
		return file.Span{Start: 0, End: 0}
	}

	return file.Span{Start: tok.GetStart(), End: tok.GetStop() + 1}
}

func spanFromTokenSafe(tok antlr.Token, src *file.Source) file.Span {
	if tok == nil {
		return file.Span{Start: 0, End: 1}
	}

	start := tok.GetStart()
	end := tok.GetStop() + 1 // exclusive end

	if start < 0 {
		start = 0
	}

	if end <= start {
		end = start + 1
	}

	// clamp to source length
	maxLen := len(src.Content())

	if end > maxLen {
		end = maxLen
	}
	if start > maxLen {
		start = maxLen - 1
	}

	return file.Span{Start: start, End: end}
}

func isIdentifier(node *TokenNode) bool {
	if node == nil {
		return false
	}

	token := node.Token()

	if token == nil {
		return false
	}

	tt := token.GetTokenType()

	return tt == fql.FqlLexerIdentifier || tt == fql.FqlLexerIgnoreIdentifier
}

func isKeyword(node *TokenNode) bool {
	if node == nil {
		return false
	}

	token := node.Token()

	if token == nil {
		return false
	}

	ttype := token.GetTokenType()

	// 0 is usually invalid; <EOF> is -1
	if ttype <= 0 || ttype >= len(fql.FqlLexerLexerStaticData.LiteralNames) {
		return false
	}

	lit := fql.FqlLexerLexerStaticData.LiteralNames[ttype]

	return strings.HasPrefix(lit, "'") && strings.HasSuffix(lit, "'")
}

func isQuote(input string) bool {
	if input == "\"" || input == "'" || input == "`" {
		return true
	}

	return false
}

func isValidString(input string) bool {
	if input == "" {
		return false
	}

	if isQuote(input) {
		return true
	}

	if len(input) >= 2 {
		first := input[0:1]
		last := input[len(input)-1:]
		return isQuote(first) && isQuote(last) && first == last
	}

	return false
}

func is(node *TokenNode, expected string) bool {
	if node == nil {
		return false
	}

	if node.GetText() == "" {
		return false
	}

	return strings.ToUpper(node.GetText()) == expected
}

func anyIs(first, second *TokenNode, expected string) *TokenNode {
	if is(first, expected) {
		return first
	}

	if is(second, expected) {
		return second
	}

	return nil
}

func has(msg string, substr string) bool {
	return strings.Contains(strings.ToLower(msg), strings.ToLower(substr))
}

func isMismatched(msg string) bool {
	return has(msg, "mismatched input")
}

func isNoAlternative(msg string) bool {
	return has(msg, "no viable alternative at input")
}

func isMissing(msg string) bool {
	return has(msg, "missing")
}

func isMissingToken(msg string, token string) bool {
	return has(msg, "missing") && has(msg, token)
}

func extractNoAlternativeInput(msg string) string {
	re := regexp.MustCompile(`no viable alternative at input\s+(?P<input>.+)`)
	match := re.FindStringSubmatch(msg)

	if match == nil || len(match) <= re.SubexpIndex("input") {
		return ""
	}

	return strings.Trim(match[re.SubexpIndex("input")], "'")
}

func extractNoAlternativeInputs(msg string) []string {
	re := regexp.MustCompile(`no viable alternative at input\s+(?P<input>.+)`)
	match := re.FindStringSubmatch(msg)

	if match == nil || len(match) <= re.SubexpIndex("input") {
		return []string{}
	}

	input := match[re.SubexpIndex("input")]
	input = strings.TrimPrefix(input, "'")
	input = strings.TrimSuffix(input, "'")

	return strings.Fields(input)
}

func isExtraneous(msg string) bool {
	return has(msg, "extraneous input")
}

func extractExtraneousInput(msg string) string {
	re := regexp.MustCompile(`extraneous input\s+(?P<input>.+?)\s+expecting`)
	match := re.FindStringSubmatch(msg)

	if match == nil || len(match) <= re.SubexpIndex("input") {
		return ""
	}

	input := match[re.SubexpIndex("input")]
	input = strings.TrimPrefix(input, "'")
	input = strings.TrimSuffix(input, "'")

	return input
}

func extractMismatchedInput(msg string) string {
	re := regexp.MustCompile(`mismatched input\s+(?P<input>.+?)\s+expecting`)
	match := re.FindStringSubmatch(msg)

	if match == nil || len(match) <= re.SubexpIndex("input") {
		return ""
	}

	input := match[re.SubexpIndex("input")]
	input = strings.TrimPrefix(input, "'")
	input = strings.TrimSuffix(input, "'")

	return input
}
