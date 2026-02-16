package diagnostics

import (
	"fmt"

	"github.com/MontFerret/ferret/v2/pkg/diagnostics"
	"github.com/MontFerret/ferret/v2/pkg/file"
	"github.com/MontFerret/ferret/v2/pkg/parser/fql"
)

func matchArrayOperatorErrors(src *file.Source, err *diagnostics.Diagnostic, offending *TokenNode) bool {
	if offending == nil {
		return false
	}

	if matchQueryOperatorErrors(src, err, offending) {
		return true
	}

	if matchArrayInlineReturnErrors(src, err, offending) {
		return true
	}

	if matchArrayQuestionQuantifierErrors(src, err, offending) {
		return true
	}

	if matchArrayOperatorUnclosed(src, err, offending) {
		return true
	}

	return false
}

func matchQueryOperatorErrors(src *file.Source, err *diagnostics.Diagnostic, offending *TokenNode) bool {
	if !(isMismatched(err.Message) || isMissing(err.Message) || isNoAlternative(err.Message)) {
		return false
	}

	if !isArrayApplyContext(offending) {
		return false
	}

	if isStringLiteral(offending) && hasPrevToken(offending, "~", 4) && !isMissingStringLiteral(err.Message) {
		span := spanFromTokenSafe(offending.Token(), src)
		err.Message = "Expected query type before query literal"
		err.Hint = "Provide a type name before the query string, e.g. doc[~ css`...`]."
		err.Spans = []diagnostics.ErrorSpan{
			diagnostics.NewMainErrorSpan(span, "missing query type"),
		}

		return true
	}

	if is(offending, "~") {
		span := spanFromTokenSafe(offending.Token(), src)
		span.Start = span.End
		span.End = span.Start + 1

		err.Message = "Expected query literal after '~'"
		err.Hint = "Provide a query literal, e.g. doc[~ css`...`]."
		err.Spans = []diagnostics.ErrorSpan{
			diagnostics.NewMainErrorSpan(span, "missing query literal"),
		}

		return true
	}

	if prev := offending.Prev(); prev != nil && is(prev, "~") {
		if isIdentifier(offending) {
			queryType := offending.GetText()
			span := spanFromTokenSafe(offending.Token(), src)
			span.Start = span.End
			span.End = span.Start + 1

			err.Message = fmt.Sprintf("Expected query string after '%s'", queryType)
			err.Hint = "Provide a query string, e.g. doc[~ css`...`]."
			err.Spans = []diagnostics.ErrorSpan{
				diagnostics.NewMainErrorSpan(span, "missing query string"),
			}

			return true
		}

		span := spanFromTokenSafe(prev.Token(), src)
		span.Start = span.End
		span.End = span.Start + 1

		err.Message = "Expected query literal after '~'"
		err.Hint = "Provide a query literal, e.g. doc[~ css`...`]."
		err.Spans = []diagnostics.ErrorSpan{
			diagnostics.NewMainErrorSpan(span, "missing query literal"),
		}

		return true
	}

	if prev := offending.Prev(); prev != nil && isIdentifier(prev) && hasPrevToken(prev, "~", 4) {
		queryType := prev.GetText()
		span := spanFromTokenSafe(prev.Token(), src)
		span.Start = span.End
		span.End = span.Start + 1

		err.Message = fmt.Sprintf("Expected query string after '%s'", queryType)
		err.Hint = "Provide a query string, e.g. doc[~ css`...`]."
		err.Spans = []diagnostics.ErrorSpan{
			diagnostics.NewMainErrorSpan(span, "missing query string"),
		}

		return true
	}

	if hasMissingClosingBracket(err.Message) {
		span := spanFromTokenSafe(offending.Token(), src)
		span.Start = span.End
		span.End = span.Start + 1

		err.Message = "Unclosed query operator"
		err.Hint = "Add a closing ']' to complete the query operator."
		err.Spans = []diagnostics.ErrorSpan{
			diagnostics.NewMainErrorSpan(span, "missing ']'"),
		}

		return true
	}

	return false
}

func matchArrayInlineReturnErrors(src *file.Source, err *diagnostics.Diagnostic, offending *TokenNode) bool {
	if !isArrayOperatorContext(offending) {
		return false
	}

	if is(offending, "RETURN") {
		span := spanFromTokenSafe(offending.Token(), src)
		span.Start = span.End
		span.End = span.Start + 1

		err.Message = "Expected expression after 'RETURN' in array operator"
		err.Hint = "Provide a projection expression, e.g. [* RETURN CURRENT]."
		err.Spans = []diagnostics.ErrorSpan{
			diagnostics.NewMainErrorSpan(span, "missing expression"),
		}

		return true
	}

	prev := offending.Prev()
	if prev == nil || !is(prev, "RETURN") {
		return false
	}

	span := spanFromTokenSafe(prev.Token(), src)
	span.Start = span.End
	span.End = span.Start + 1

	err.Message = "Expected expression after 'RETURN' in array operator"
	err.Hint = "Provide a projection expression, e.g. [* RETURN CURRENT]."
	err.Spans = []diagnostics.ErrorSpan{
		diagnostics.NewMainErrorSpan(span, "missing expression"),
	}

	return true
}

func matchArrayQuestionQuantifierErrors(src *file.Source, err *diagnostics.Diagnostic, offending *TokenNode) bool {
	if !(isMismatched(err.Message) || isMissing(err.Message) || isNoAlternative(err.Message)) {
		return false
	}

	if !isArrayQuestionHead(offending) {
		return false
	}

	if hasPrevToken(offending, "FILTER", 4) {
		return false
	}

	prev := offending.Prev()
	span := spanFromTokenSafe(offending.Token(), src)

	if prev != nil {
		span = spanFromTokenSafe(prev.Token(), src)
	}

	span.Start = span.End
	span.End = span.Start + 1

	err.Message = "Expected FILTER after quantifier in array filter"
	err.Hint = "Add a FILTER expression, e.g. [? NONE FILTER <expr>]."
	err.Spans = []diagnostics.ErrorSpan{
		diagnostics.NewMainErrorSpan(span, "missing 'FILTER'"),
	}

	return true
}

func matchArrayOperatorUnclosed(src *file.Source, err *diagnostics.Diagnostic, offending *TokenNode) bool {
	if !(isMismatched(err.Message) || isMissing(err.Message) || isNoAlternative(err.Message)) {
		return false
	}

	if !hasMissingClosingBracket(err.Message) {
		return false
	}

	if !isArrayOperatorContext(offending) && !isArrayQuestionHead(offending) {
		return false
	}

	span := spanFromTokenSafe(offending.Token(), src)
	span.Start = span.End
	span.End = span.Start + 1

	err.Message = "Unclosed array operator"
	err.Hint = "Add a closing ']' to complete the array operator."
	err.Spans = []diagnostics.ErrorSpan{
		diagnostics.NewMainErrorSpan(span, "missing ']'"),
	}

	return true
}

func isStringLiteral(node *TokenNode) bool {
	if node == nil || node.Token() == nil {
		return false
	}

	return node.Token().GetTokenType() == fql.FqlLexerStringLiteral
}

func isMissingStringLiteral(msg string) bool {
	return isMissing(msg) && has(msg, "stringliteral")
}

func hasMissingClosingBracket(msg string) bool {
	if isMissingToken(msg, "]") {
		return true
	}

	return isMismatched(msg) && has(msg, "']'")
}

func isArrayOperatorContext(node *TokenNode) bool {
	_, star := prevTokenDistance(node, "*", 8)
	if star == nil {
		return false
	}

	return isArrayStarHead(star)
}

func isArrayQuestionHead(node *TokenNode) bool {
	_, question := prevTokenDistance(node, "?", 8)
	if question == nil {
		return false
	}

	return question.Prev() != nil && is(question.Prev(), "[")
}

func isArrayApplyContext(node *TokenNode) bool {
	_, tilde := prevTokenDistance(node, "~", 8)
	if tilde == nil {
		return false
	}

	return tilde.Prev() != nil && is(tilde.Prev(), "[")
}

func isArrayStarHead(star *TokenNode) bool {
	current := star

	for i := 0; i < 4 && current != nil; i++ {
		prev := current.Prev()
		if prev == nil {
			return false
		}

		if is(prev, "[") {
			return true
		}

		if !is(prev, "*") {
			return false
		}

		current = prev
	}

	return false
}

func prevTokenDistance(node *TokenNode, token string, max int) (int, *TokenNode) {
	current := node

	for i := 0; i < max && current != nil; i++ {
		if is(current, token) {
			return i, current
		}
		current = current.Prev()
	}

	return -1, nil
}
