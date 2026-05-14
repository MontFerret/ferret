package diagnostics

import (
	"github.com/antlr4-go/antlr/v4"

	"github.com/MontFerret/ferret/v2/pkg/diagnostics"
	"github.com/MontFerret/ferret/v2/pkg/parser/fql"
	"github.com/MontFerret/ferret/v2/pkg/source"
)

const (
	missingArrayItemCommaMessage = "Expected ',' between array items"
	missingArrayItemCommaHint    = "Separate array items with commas, e.g. [1, 2, 3]."
)

func matchArrayLiteralSeparatorErrors(src *source.Source, err *diagnostics.Diagnostic, offending *TokenNode) bool {
	if src == nil || err == nil || offending == nil || offending.Token() == nil {
		return false
	}

	if !isNoAlternative(err.Message) && !isMissing(err.Message) && !isMismatched(err.Message) {
		return false
	}

	tokens := lexDefaultTokens(src.Content())
	offendingIdx := findLexedTokenIndex(tokens, offending.Token())
	if offendingIdx < 0 {
		offendingIdx = findDiagnosticSpanTokenIndex(tokens, err)
	}
	if offendingIdx < 0 {
		return false
	}

	openIdx := findEnclosingArrayLiteralOpen(tokens, offendingIdx)
	if openIdx < 0 {
		return false
	}

	closeIdx := findMatchingCloseBracket(tokens, openIdx)
	if closeIdx < 0 || offendingIdx >= closeIdx {
		return false
	}

	if isComputedPropertyNameOpen(tokens, openIdx, closeIdx) {
		return false
	}

	span, ok := arrayItemSeparatorInsertionSpan(tokens, openIdx, closeIdx, src)
	if !ok {
		return false
	}

	err.Message = missingArrayItemCommaMessage
	err.Hint = missingArrayItemCommaHint
	err.Spans = []diagnostics.ErrorSpan{
		diagnostics.NewMainErrorSpan(span, "missing comma"),
	}

	return true
}

func isArrayLiteralSeparatorDiagnostic(err *diagnostics.Diagnostic) bool {
	return err != nil && err.Kind == SyntaxError && err.Message == missingArrayItemCommaMessage
}

func findEnclosingArrayLiteralOpen(tokens []antlr.Token, offendingIdx int) int {
	if offendingIdx <= 0 || offendingIdx > len(tokens) {
		return -1
	}

	stack := make([]int, 0, 8)
	for i := 0; i < offendingIdx; i++ {
		text := tokenText(tokens[i])

		switch text {
		case "(", "[", "{":
			stack = append(stack, i)
		case ")", "]", "}":
			if len(stack) == 0 {
				return -1
			}

			openIdx := stack[len(stack)-1]
			if !delimitersMatch(tokenText(tokens[openIdx]), text) {
				return -1
			}

			stack = stack[:len(stack)-1]
		}
	}

	if len(stack) == 0 {
		return -1
	}

	for i := len(stack) - 1; i >= 0; i-- {
		openIdx := stack[i]
		if isTokenText(tokens[openIdx], "[") && isArrayLiteralOpen(tokens, openIdx) {
			return openIdx
		}
	}

	return -1
}

func isArrayLiteralOpen(tokens []antlr.Token, openIdx int) bool {
	if openIdx < 0 || openIdx >= len(tokens) || !isTokenText(tokens[openIdx], "[") {
		return false
	}

	if openIdx == 0 {
		return true
	}

	prev := tokens[openIdx-1]
	if isTokenText(prev, ".") {
		return false
	}

	return !tokenCanEndExpression(prev)
}

func findMatchingCloseBracket(tokens []antlr.Token, openIdx int) int {
	if openIdx < 0 || openIdx >= len(tokens) || !isTokenText(tokens[openIdx], "[") {
		return -1
	}

	depth := 0
	for i := openIdx + 1; i < len(tokens); i++ {
		switch tokenText(tokens[i]) {
		case "[":
			depth++
		case "]":
			if depth == 0 {
				return i
			}

			depth--
		}
	}

	return -1
}

func isComputedPropertyNameOpen(tokens []antlr.Token, openIdx, closeIdx int) bool {
	if openIdx <= 0 || closeIdx <= openIdx || closeIdx+1 >= len(tokens) {
		return false
	}

	prev := tokenText(tokens[openIdx-1])
	return (prev == "{" || prev == ",") && isTokenText(tokens[closeIdx+1], ":")
}

func arrayItemSeparatorInsertionSpan(tokens []antlr.Token, openIdx, closeIdx int, src *source.Source) (source.Span, bool) {
	previousIdx := missingArrayItemSeparatorPreviousToken(tokens, openIdx, closeIdx)
	if previousIdx < 0 || isTokenText(tokens[previousIdx], ",") {
		return source.Span{}, false
	}

	previous := tokens[previousIdx]
	offset := previous.GetStop() + 1
	if offset < 0 || offset >= len(src.Content()) {
		return spanFromTokenSafe(previous, src), true
	}

	return source.Span{Start: offset, End: offset + 1}, true
}

func missingArrayItemSeparatorPreviousToken(tokens []antlr.Token, openIdx, closeIdx int) int {
	depth := 0
	previousEndIdx := -1

	for i := openIdx + 1; i < closeIdx; i++ {
		text := tokenText(tokens[i])

		if depth == 0 {
			if text == "," {
				previousEndIdx = -1
				continue
			}

			if isTopLevelArrayExpressionContinuation(tokens, i) {
				previousEndIdx = -1
				if text == "(" || text == "[" || text == "{" {
					depth++
				}
				continue
			}

			if previousEndIdx >= 0 && isArrayItemStartToken(tokens[i]) {
				return previousEndIdx
			}
		}

		switch text {
		case "(", "[", "{":
			depth++
		case ")", "]", "}":
			if depth == 0 {
				return -1
			}

			depth--
			if depth == 0 {
				previousEndIdx = i
			}
		default:
			if depth == 0 {
				if tokenCanEndExpression(tokens[i]) {
					previousEndIdx = i
				} else {
					previousEndIdx = -1
				}
			}
		}
	}

	return -1
}

func tokenCanEndExpression(token antlr.Token) bool {
	if token == nil {
		return false
	}

	switch token.GetTokenType() {
	case fql.FqlLexerIdentifier,
		fql.FqlLexerIgnoreIdentifier,
		fql.FqlLexerStringLiteral,
		fql.FqlLexerBacktickOpen,
		fql.FqlLexerDurationLiteral,
		fql.FqlLexerIntegerLiteral,
		fql.FqlLexerFloatLiteral,
		fql.FqlLexerNone,
		fql.FqlLexerNull,
		fql.FqlLexerBooleanLiteral,
		fql.FqlLexerCloseBracket,
		fql.FqlLexerCloseParen,
		fql.FqlLexerCloseBrace:
		return true
	default:
		return false
	}
}

func isArrayItemStartToken(token antlr.Token) bool {
	if token == nil {
		return false
	}

	switch token.GetTokenType() {
	case fql.FqlLexerOpenBrace,
		fql.FqlLexerOpenBracket,
		fql.FqlLexerOpenParen,
		fql.FqlLexerIdentifier,
		fql.FqlLexerIgnoreIdentifier,
		fql.FqlLexerStringLiteral,
		fql.FqlLexerBacktickOpen,
		fql.FqlLexerDurationLiteral,
		fql.FqlLexerIntegerLiteral,
		fql.FqlLexerFloatLiteral,
		fql.FqlLexerParam,
		fql.FqlLexerNone,
		fql.FqlLexerNull,
		fql.FqlLexerBooleanLiteral,
		fql.FqlLexerMatch,
		fql.FqlLexerQuery,
		fql.FqlLexerWaitfor,
		fql.FqlLexerDispatch,
		fql.FqlLexerFor,
		fql.FqlLexerAs,
		fql.FqlLexerDistinct,
		fql.FqlLexerFilter,
		fql.FqlLexerSort,
		fql.FqlLexerLimit,
		fql.FqlLexerCollect,
		fql.FqlLexerCount,
		fql.FqlLexerSortDirection,
		fql.FqlLexerInto,
		fql.FqlLexerKeep,
		fql.FqlLexerWith,
		fql.FqlLexerAt,
		fql.FqlLexerLeast,
		fql.FqlLexerAggregate,
		fql.FqlLexerEvent,
		fql.FqlLexerTimeout,
		fql.FqlLexerOptions,
		fql.FqlLexerEvery,
		fql.FqlLexerBackoff,
		fql.FqlLexerJitter,
		fql.FqlLexerExists,
		fql.FqlLexerValue,
		fql.FqlLexerOne:
		return true
	default:
		return false
	}
}

func isTopLevelArrayExpressionContinuation(tokens []antlr.Token, idx int) bool {
	if idx < 0 || idx >= len(tokens) {
		return false
	}

	token := tokens[idx]
	prev := func() antlr.Token {
		if idx == 0 {
			return nil
		}

		return tokens[idx-1]
	}()

	switch token.GetTokenType() {
	case fql.FqlLexerOpenParen:
		return tokenCanBeFunctionName(prev)
	case fql.FqlLexerOpenBracket:
		return tokenCanEndExpression(prev)
	case fql.FqlLexerDot,
		fql.FqlLexerQuestionMark,
		fql.FqlLexerPlus,
		fql.FqlLexerMinus,
		fql.FqlLexerMulti,
		fql.FqlLexerDiv,
		fql.FqlLexerMod,
		fql.FqlLexerGt,
		fql.FqlLexerLt,
		fql.FqlLexerEq,
		fql.FqlLexerGte,
		fql.FqlLexerLte,
		fql.FqlLexerNeq,
		fql.FqlLexerRegexMatch,
		fql.FqlLexerRegexNotMatch,
		fql.FqlLexerAnd,
		fql.FqlLexerOr,
		fql.FqlLexerIn,
		fql.FqlLexerLike,
		fql.FqlLexerNot,
		fql.FqlLexerAll,
		fql.FqlLexerAny:
		return true
	default:
		return false
	}
}

func tokenCanBeFunctionName(token antlr.Token) bool {
	if token == nil {
		return false
	}

	switch token.GetTokenType() {
	case fql.FqlLexerIdentifier,
		fql.FqlLexerIgnoreIdentifier,
		fql.FqlLexerAnd,
		fql.FqlLexerOr,
		fql.FqlLexerAs,
		fql.FqlLexerDistinct,
		fql.FqlLexerFilter,
		fql.FqlLexerSort,
		fql.FqlLexerLimit,
		fql.FqlLexerCollect,
		fql.FqlLexerCount,
		fql.FqlLexerSortDirection,
		fql.FqlLexerInto,
		fql.FqlLexerKeep,
		fql.FqlLexerWith,
		fql.FqlLexerAll,
		fql.FqlLexerAny,
		fql.FqlLexerAt,
		fql.FqlLexerLeast,
		fql.FqlLexerAggregate,
		fql.FqlLexerEvent,
		fql.FqlLexerTimeout,
		fql.FqlLexerOptions,
		fql.FqlLexerEvery,
		fql.FqlLexerBackoff,
		fql.FqlLexerJitter,
		fql.FqlLexerExists,
		fql.FqlLexerValue,
		fql.FqlLexerOne,
		fql.FqlLexerReturn,
		fql.FqlLexerDispatch,
		fql.FqlLexerQuery,
		fql.FqlLexerUsing,
		fql.FqlLexerNone,
		fql.FqlLexerNull,
		fql.FqlLexerLet,
		fql.FqlLexerVar,
		fql.FqlLexerUse,
		fql.FqlLexerWaitfor,
		fql.FqlLexerWhile,
		fql.FqlLexerDo,
		fql.FqlLexerIn,
		fql.FqlLexerLike,
		fql.FqlLexerNot,
		fql.FqlLexerFor,
		fql.FqlLexerBooleanLiteral,
		fql.FqlLexerMatch,
		fql.FqlLexerWhen:
		return true
	default:
		return false
	}
}

func delimitersMatch(open, close string) bool {
	switch open {
	case "(":
		return close == ")"
	case "[":
		return close == "]"
	case "{":
		return close == "}"
	default:
		return false
	}
}
