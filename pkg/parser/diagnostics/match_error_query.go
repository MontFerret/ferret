package diagnostics

import (
	"strings"

	"github.com/MontFerret/ferret/v2/pkg/diagnostics"
	"github.com/MontFerret/ferret/v2/pkg/parser/fql"
	"github.com/MontFerret/ferret/v2/pkg/source"
)

func matchQueryErrors(src *source.Source, err *diagnostics.Diagnostic, offending *TokenNode) bool {
	if err == nil || offending == nil {
		return false
	}

	if !isMismatched(err.Message) && !isMissing(err.Message) && !isNoAlternative(err.Message) && !isExtraneous(err.Message) && !has(err.Message, "queryexpression failed predicate") {
		return false
	}

	if !hasPrevToken(offending, "QUERY", 10) && !is(offending, "QUERY") && !hasPrevToken(offending, "USING", 10) {
		return false
	}

	if prev := offending.Prev(); prev != nil && is(prev, "WITH") {
		span := spanFromTokenSafe(prev.Token(), src)
		err.Message = "Expected options expression after WITH"
		err.Hint = "Provide an options expression, e.g. WITH { limit: 10 }."
		err.Spans = []diagnostics.ErrorSpan{
			diagnostics.NewMainErrorSpan(span, "missing expression"),
		}
		return true
	}

	if isMissingQueryLiteral(err.Message, offending) {
		spanNode := offending
		if prev := offending.Prev(); prev != nil && is(prev, "QUERY") {
			spanNode = prev
		}

		span := spanFromTokenSafe(spanNode.Token(), src)
		err.Message = "QUERY requires a query literal"
		err.Hint = "Provide a query literal, e.g. QUERY `.items` IN doc USING css or QUERY @q IN doc USING css."
		err.Spans = []diagnostics.ErrorSpan{
			diagnostics.NewMainErrorSpan(span, "missing query literal"),
		}
		return true
	}

	if expectsKeyword(err.Message, "in") && hasPrevToken(offending, "QUERY", 10) && hasQueryLiteralBetween(offending, 10) && !hasTokenBefore(offending, "QUERY", "IN", 10) {
		span := spanFromTokenSafe(offending.Token(), src)
		err.Message = "Expected IN after query literal"
		err.Hint = "Add IN <expr>, e.g. QUERY `.items` IN doc USING css."
		err.Spans = []diagnostics.ErrorSpan{
			diagnostics.NewMainErrorSpan(span, "missing 'IN'"),
		}
		return true
	}

	if (is(offending, "USING") || is(offending, "WITH") || isEOF(offending)) && hasPrevToken(offending, "IN", 2) {
		span := spanFromTokenSafe(offending.Token(), src)
		err.Message = "Expected expression after IN"
		err.Hint = "Provide a source expression, e.g. QUERY `.items` IN doc USING css."
		err.Spans = []diagnostics.ErrorSpan{
			diagnostics.NewMainErrorSpan(span, "missing expression"),
		}
		return true
	}

	if expectsKeyword(err.Message, "using") && hasTokenBefore(offending, "QUERY", "IN", 10) && !hasTokenBefore(offending, "QUERY", "USING", 10) {
		span := spanFromTokenSafe(offending.Token(), src)
		err.Message = "Expected USING <dialect> after IN expression"
		err.Hint = "Add USING <dialect>, e.g. QUERY `.items` IN doc USING css."
		err.Spans = []diagnostics.ErrorSpan{
			diagnostics.NewMainErrorSpan(span, "missing 'USING'"),
		}
		return true
	}

	if hasPrevToken(offending, "USING", 2) {
		if is(offending, "USING") {
			if next := offending.Next(); next != nil {
				if isEOF(next) {
					span := spanFromTokenSafe(offending.Token(), src)
					err.Message = "Expected dialect identifier after USING"
					err.Hint = "Provide a dialect identifier, e.g. USING css."
					err.Spans = []diagnostics.ErrorSpan{
						diagnostics.NewMainErrorSpan(span, "missing dialect"),
					}
					return true
				}

				if !isIdentifier(next) {
					span := spanFromTokenSafe(next.Token(), src)
					err.Message = "Dialect after USING must be an identifier"
					err.Hint = "Provide a dialect identifier such as css or xpath."
					err.Spans = []diagnostics.ErrorSpan{
						diagnostics.NewMainErrorSpan(span, "invalid dialect"),
					}
					return true
				}
			} else if dialectTokenInMessage(err.Message) {
				span := spanFromTokenSafe(offending.Token(), src)
				err.Message = "Dialect after USING must be an identifier"
				err.Hint = "Provide a dialect identifier such as css or xpath."
				err.Spans = []diagnostics.ErrorSpan{
					diagnostics.NewMainErrorSpan(span, "invalid dialect"),
				}
				return true
			} else if next := offending.Next(); next == nil || isEOF(next) || is(next, "WITH") {
				span := spanFromTokenSafe(offending.Token(), src)
				err.Message = "Expected dialect identifier after USING"
				err.Hint = "Provide a dialect identifier, e.g. USING css."
				err.Spans = []diagnostics.ErrorSpan{
					diagnostics.NewMainErrorSpan(span, "missing dialect"),
				}
				return true
			}
		}

		if is(offending, "WITH") || isEOF(offending) {
			span := spanFromTokenSafe(offending.Token(), src)
			err.Message = "Expected dialect identifier after USING"
			err.Hint = "Provide a dialect identifier, e.g. USING css."
			err.Spans = []diagnostics.ErrorSpan{
				diagnostics.NewMainErrorSpan(span, "missing dialect"),
			}
			return true
		}

		if !isIdentifier(offending) {
			span := spanFromTokenSafe(offending.Token(), src)
			err.Message = "Dialect after USING must be an identifier"
			err.Hint = "Provide a dialect identifier such as css or xpath."
			err.Spans = []diagnostics.ErrorSpan{
				diagnostics.NewMainErrorSpan(span, "invalid dialect"),
			}
			return true
		}
	}

	if (is(offending, ".") || isStringLiteral(offending)) && hasPrevToken(offending, "USING", 4) {
		span := spanFromTokenSafe(offending.Token(), src)
		err.Message = "Dialect after USING must be an identifier"
		err.Hint = "Provide a dialect identifier such as css or xpath."
		err.Spans = []diagnostics.ErrorSpan{
			diagnostics.NewMainErrorSpan(span, "invalid dialect"),
		}
		return true
	}

	return false
}

func isMissingQueryLiteral(msg string, offending *TokenNode) bool {
	if offending == nil {
		return false
	}

	if is(offending, "IN") && is(offending.Prev(), "QUERY") {
		return true
	}

	if is(offending, "QUERY") && is(offending.Next(), "IN") {
		return true
	}

	if hasPrevToken(offending, "QUERY", 10) && !hasQueryLiteralBetween(offending, 10) {
		return true
	}

	if has(msg, "stringliteral") || has(msg, "backtickopen") {
		return true
	}

	if isMismatched(msg) || isNoAlternative(msg) {
		if prev := offending.Prev(); prev != nil && is(prev, "QUERY") {
			return true
		}
	}

	if has(msg, "queryexpression failed predicate") {
		return true
	}

	return false
}

func expectsKeyword(msg string, keyword string) bool {
	return strings.Contains(strings.ToLower(msg), strings.ToLower(keyword))
}

func isEOF(node *TokenNode) bool {
	return node != nil && is(node, "<EOF>")
}

func hasTokenBefore(node *TokenNode, stopToken string, tokenText string, max int) bool {
	current := node
	for i := 0; i < max && current != nil; i++ {
		if is(current, stopToken) {
			return false
		}
		if is(current, tokenText) {
			return true
		}
		current = current.Prev()
	}
	return false
}

func hasQueryLiteralBetween(node *TokenNode, max int) bool {
	current := node
	for i := 0; i < max && current != nil; i++ {
		if is(current, "QUERY") {
			return false
		}
		if isStringLiteral(current) || isBacktickToken(current) || isParamToken(current) || isIdentifier(current) || (isSafeReservedWordToken(current) && !isQueryModifierToken(current)) {
			return true
		}
		current = current.Prev()
	}
	return false
}

func dialectTokenInMessage(msg string) bool {
	if has(msg, "stringliteral") || has(msg, "backtickopen") {
		return true
	}

	mismatched := extractMismatchedInput(msg)
	if mismatched == "" {
		mismatched = extractExtraneousInput(msg)
	}
	if isValidString(mismatched) || mismatched == "." {
		return true
	}

	return false
}

func isBacktickToken(node *TokenNode) bool {
	if node == nil || node.Token() == nil {
		return false
	}

	tt := node.Token().GetTokenType()
	return tt == fql.FqlLexerBacktickOpen || tt == fql.FqlLexerBacktickClose
}

func isParamToken(node *TokenNode) bool {
	if node == nil || node.Token() == nil {
		return false
	}

	return node.Token().GetTokenType() == fql.FqlLexerParam
}

func isQueryModifierToken(node *TokenNode) bool {
	if node == nil || node.Token() == nil {
		return false
	}

	switch node.Token().GetTokenType() {
	case fql.FqlLexerExists,
		fql.FqlLexerAny,
		fql.FqlLexerValue,
		fql.FqlLexerCount,
		fql.FqlLexerOne:
		return true
	default:
		return false
	}
}

func isSafeReservedWordToken(node *TokenNode) bool {
	if node == nil || node.Token() == nil {
		return false
	}

	switch node.Token().GetTokenType() {
	case fql.FqlLexerAnd,
		fql.FqlLexerOr,
		fql.FqlLexerAs,
		fql.FqlLexerDistinct,
		fql.FqlLexerFilter,
		fql.FqlLexerSort,
		fql.FqlLexerLimit,
		fql.FqlLexerCollect,
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
		fql.FqlLexerCount,
		fql.FqlLexerOne:
		return true
	default:
		return false
	}
}
