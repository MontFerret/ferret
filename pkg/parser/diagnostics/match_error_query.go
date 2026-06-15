package diagnostics

import (
	"fmt"
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

	if with := outOfOrderQueryWithNode(src, offending); with != nil {
		span := spanFromTokenSafe(with.Token(), src)
		err.Message = "WITH must appear before OPTIONS in QUERY"
		err.Hint = "Move the WITH clause before OPTIONS."
		err.Spans = []diagnostics.ErrorSpan{
			diagnostics.NewMainErrorSpan(span, "out-of-order WITH clause"),
		}
		return true
	}

	if !hasPrevToken(offending, "QUERY", 64) && !is(offending, "QUERY") && !hasPrevToken(offending, "USING", 64) && !sourceHasQueryBefore(src, offending) {
		return false
	}

	if is(offending, "WITH") && (hasTokenBefore(offending, "QUERY", "OPTIONS", 64) || sourceHasQueryOptionsBefore(src, offending)) {
		span := spanFromTokenSafe(offending.Token(), src)
		err.Message = "WITH must appear before OPTIONS in QUERY"
		err.Hint = "Move the WITH clause before OPTIONS."
		err.Spans = []diagnostics.ErrorSpan{
			diagnostics.NewMainErrorSpan(span, "out-of-order WITH clause"),
		}
		return true
	}

	if (is(offending, "WITH") || is(offending, "OPTIONS")) && !hasPrevToken(offending, "USING", 2) {
		clause := offending.GetText()
		span := spanFromTokenSafe(offending.Token(), src)
		err.Message = fmt.Sprintf("Expected query %s expression after %s", queryClauseValueName(clause), clause)
		err.Hint = queryClauseHint(clause)
		err.Spans = []diagnostics.ErrorSpan{
			diagnostics.NewMainErrorSpan(span, "missing expression"),
		}
		return true
	}

	if prev := offending.Prev(); prev != nil && (is(prev, "WITH") || is(prev, "OPTIONS")) {
		span := spanFromTokenSafe(prev.Token(), src)
		clause := prev.GetText()
		err.Message = fmt.Sprintf("Expected query %s expression after %s", queryClauseValueName(clause), clause)
		err.Hint = queryClauseHint(clause)
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

	if (is(offending, "USING") || is(offending, "WITH") || is(offending, "OPTIONS") || isEOF(offending)) && hasPrevToken(offending, "IN", 2) {
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
			} else if next := offending.Next(); next == nil || isEOF(next) || is(next, "WITH") || is(next, "OPTIONS") {
				span := spanFromTokenSafe(offending.Token(), src)
				err.Message = "Expected dialect identifier after USING"
				err.Hint = "Provide a dialect identifier, e.g. USING css."
				err.Spans = []diagnostics.ErrorSpan{
					diagnostics.NewMainErrorSpan(span, "missing dialect"),
				}
				return true
			}
		}

		if is(offending, "WITH") || is(offending, "OPTIONS") || isEOF(offending) {
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

func queryClauseHint(clause string) string {
	switch clause {
	case "WITH":
		return "Provide a params expression, e.g. WITH { params: [1] }."
	case "OPTIONS":
		return "Provide an options expression, e.g. OPTIONS { timeout: 5000 }."
	default:
		return "Provide a query clause expression."
	}
}

func queryClauseValueName(clause string) string {
	switch clause {
	case "WITH":
		return "params"
	case "OPTIONS":
		return "options"
	default:
		return "clause"
	}
}

func sourceHasQueryBefore(src *source.Source, offending *TokenNode) bool {
	prefix := sourcePrefixBefore(src, offending)
	query := strings.LastIndex(prefix, "QUERY")
	using := strings.LastIndex(prefix, "USING")
	dispatch := strings.LastIndex(prefix, "DISPATCH")

	return query >= 0 && using > query && dispatch < query
}

func sourceHasQueryOptionsBefore(src *source.Source, offending *TokenNode) bool {
	prefix := sourcePrefixBefore(src, offending)
	query := strings.LastIndex(prefix, "QUERY")
	options := strings.LastIndex(prefix, "OPTIONS")

	return query >= 0 && options > query
}

func sourceHasOutOfOrderQueryWith(src *source.Source) bool {
	if src == nil {
		return false
	}

	content := strings.ToUpper(src.Content())
	query := strings.LastIndex(content, "QUERY")
	options := strings.LastIndex(content, "OPTIONS")
	with := strings.LastIndex(content, "WITH")
	dispatch := strings.LastIndex(content, "DISPATCH")

	return query >= 0 && options > query && with > options && dispatch < query
}

func outOfOrderQueryWithNode(src *source.Source, offending *TokenNode) *TokenNode {
	if !sourceHasOutOfOrderQueryWith(src) {
		return nil
	}
	if is(offending, "WITH") {
		return offending
	}
	if next := offending.Next(); is(next, "WITH") {
		return next
	}

	return nil
}

func sourcePrefixBefore(src *source.Source, offending *TokenNode) string {
	if src == nil || offending == nil || offending.Token() == nil {
		return ""
	}

	content := []rune(src.Content())
	end := offending.Token().GetStart()
	if end < 0 || end > len(content) {
		end = len(content)
	}

	return strings.ToUpper(string(content[:end]))
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
