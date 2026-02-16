package diagnostics

import (
	"fmt"
	"regexp"

	"github.com/MontFerret/ferret/v2/pkg/diagnostics"
	"github.com/MontFerret/ferret/v2/pkg/file"
)

func matchCommonErrors(src *file.Source, err *diagnostics.Diagnostic, offending *TokenNode) bool {
	if isNoAlternative(err.Message) {
		prevLogical := isLogicalOperator(offending.Prev())
		if prevLogical && !isExpressionStart(offending) {
			span := spanFromTokenSafe(offending.Prev().Token(), src)
			span.Start += 2
			span.End += 2

			operator := offending.Prev().GetText()
			err.Message = fmt.Sprintf("Expected right-hand expression after '%s'", operator)
			err.Hint = fmt.Sprintf("Provide an expression after the logical operator, e.g. (a %s b).", operator)
			err.Spans = []diagnostics.ErrorSpan{
				diagnostics.NewMainErrorSpan(span, "missing expression"),
			}

			return true
		}

		if groupParen := findGroupingParen(offending, 8); groupParen != nil {
			span := spanFromTokenSafe(offending.Token(), src)
			span.Start++
			span.End++

			err.Message = "Unclosed parenthesized expression"
			err.Hint = "Add a closing ')' to complete the expression."
			err.Spans = []diagnostics.ErrorSpan{
				diagnostics.NewMainErrorSpan(span, "missing ')'"),
			}

			return true
		}

		if is(offending.Prev(), ",") {
			span := spanFromTokenSafe(offending.Prev().Token(), src)
			span.Start++
			span.End++

			err.Message = "Expected expression after ','"
			err.Hint = "Did you forget to provide a value?"
			err.Spans = []diagnostics.ErrorSpan{
				diagnostics.NewMainErrorSpan(span, "missing value"),
			}

			return true
		}

		// Ternary operator, incomplete expression
		if is(offending.Prev(), "?") {
			span := spanFromTokenSafe(offending.Prev().Token(), src)
			span.Start++
			span.End++

			err.Message = "Expected expression after '?' in ternary operator"
			err.Hint = "Provide an expression after the question mark to complete the ternary operation."
			err.Spans = []diagnostics.ErrorSpan{
				diagnostics.NewMainErrorSpan(span, "missing expression"),
			}

			return true
		}

		// Ternary operator, missing the right-hand expression
		if is(offending.Prev(), ":") {
			span := spanFromTokenSafe(offending.Prev().Token(), src)
			span.Start++
			span.End++

			err.Message = "Expected expression after ':' in ternary operator"
			err.Hint = "Provide an expression after the colon to complete the ternary operation."
			err.Spans = []diagnostics.ErrorSpan{
				diagnostics.NewMainErrorSpan(span, "missing expression"),
			}

			return true
		}
	}

	if isMissing(err.Message) {
		if is(offending.Prev(), "..") || is(offending, "..") || has(err.Message, "..") {
			span := spanFromTokenSafe(offending.Token(), src)
			span.Start += 2
			span.End += 2

			start := ""
			if is(offending, "..") && offending.Prev() != nil {
				start = offending.Prev().GetText()
			} else if is(offending.Prev(), "..") && offending.Prev().Prev() != nil {
				start = offending.Prev().Prev().GetText()
			} else {
				start = extractRangeStart(err.Message)
			}

			err.Message = "Expected end value after '..' in range expression"
			err.Hint = fmt.Sprintf("Provide an end value to complete the range, e.g. %s..10.", start)
			err.Spans = []diagnostics.ErrorSpan{
				diagnostics.NewMainErrorSpan(span, "missing value"),
			}

			return true
		}
	}

	if isNoAlternative(err.Message) || isMissing(err.Message) {
		if is(offending.Prev(), "(") {
			var span file.Span

			if isKeyword(offending) {
				span = spanFromTokenSafe(offending.Prev().Token(), src)
				span.Start++
				span.End++
			} else {
				span = spanFromTokenSafe(offending.Token(), src)
				span.Start++
				span.End++
			}

			err.Message = "Unclosed function call"
			err.Hint = "Add a closing ')' to complete the function call."
			err.Spans = []diagnostics.ErrorSpan{
				diagnostics.NewMainErrorSpan(span, "missing ')'"),
			}

			return true
		}
	}

	if isMissing(err.Message) {
		if isMissingToken(err.Message, ")") {
			span := spanFromTokenSafe(offending.Token(), src)
			span.Start++
			span.End++

			err.Message = "Unclosed parenthesized expression"
			err.Hint = "Add a closing ')' to complete the expression."
			err.Spans = []diagnostics.ErrorSpan{
				diagnostics.NewMainErrorSpan(span, "missing ')'"),
			}

			return true
		}

		if is(offending.Prev(), "..") {
			span := spanFromTokenSafe(offending.Prev().Token(), src)
			span.Start += 2
			span.End += 2

			start := ""
			if prevPrev := offending.Prev().Prev(); prevPrev != nil {
				start = prevPrev.GetText()
			}

			err.Message = "Expected end value after '..' in range expression"
			err.Hint = fmt.Sprintf("Provide an end value to complete the range, e.g. %s..10.", start)
			err.Spans = []diagnostics.ErrorSpan{
				diagnostics.NewMainErrorSpan(span, "missing value"),
			}

			return true
		}

		if is(offending, "..") || has(err.Message, "..") {
			span := spanFromTokenSafe(offending.Token(), src)
			span.Start += 2
			span.End += 2

			start := ""
			if is(offending, "..") && offending.Prev() != nil {
				start = offending.Prev().GetText()
			} else {
				start = extractRangeStart(err.Message)
			}

			err.Message = "Expected end value after '..' in range expression"
			err.Hint = fmt.Sprintf("Provide an end value to complete the range, e.g. %s..10.", start)
			err.Spans = []diagnostics.ErrorSpan{
				diagnostics.NewMainErrorSpan(span, "missing value"),
			}

			return true
		}
	}

	if isExtraneous(err.Message) {
		if is(offending, "(") {
			span := spanFromTokenSafe(offending.Token(), src)
			span.Start++
			span.End++

			err.Message = "Expected a valid list of arguments"
			err.Hint = "Did you forget to provide a value?"
			err.Spans = []diagnostics.ErrorSpan{
				diagnostics.NewMainErrorSpan(span, "missing value"),
			}

			return true
		}

		if is(offending, "..") {
			span := spanFromTokenSafe(offending.Token(), src)
			span.Start += 2
			span.End += 2

			start := offending.Prev().GetText()
			err.Message = "Expected end value after '..' in range expression"
			err.Hint = fmt.Sprintf("Provide an end value to complete the range, e.g. %s..10.", start)
			err.Spans = []diagnostics.ErrorSpan{
				diagnostics.NewMainErrorSpan(span, "missing value"),
			}

			return true
		}
	}

	return false
}

func findGroupingParen(node *TokenNode, max int) *TokenNode {
	current := node

	for i := 0; i < max && current != nil; i++ {
		if is(current, "(") {
			prev := current.Prev()
			if isIdentifier(prev) {
				return nil
			}

			return current
		}

		current = current.Prev()
	}

	return nil
}

func isLogicalOperator(node *TokenNode) bool {
	if node == nil {
		return false
	}

	return is(node, "||") || is(node, "OR") || is(node, "&&") || is(node, "AND")
}

func isExpressionStart(node *TokenNode) bool {
	if node == nil {
		return false
	}

	if isIdentifier(node) {
		return true
	}

	if !isKeyword(node) {
		return true
	}

	return is(node, "(") || is(node, "[") || is(node, "{") || is(node, "@") || is(node, "FOR") || is(node, "WAITFOR") || is(node, "DISPATCH")
}

var rangeStartRe = regexp.MustCompile(`([A-Za-z0-9_]+)\\.\\.`)

func extractRangeStart(msg string) string {
	if msg == "" {
		return ""
	}

	matches := rangeStartRe.FindAllStringSubmatch(msg, -1)
	if len(matches) == 0 {
		return ""
	}

	last := matches[len(matches)-1]
	if len(last) < 2 {
		return ""
	}

	return last[1]
}
