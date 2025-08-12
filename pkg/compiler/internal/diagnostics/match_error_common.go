package diagnostics

import (
	"fmt"

	"github.com/MontFerret/ferret/pkg/file"
)

func matchCommonErrors(src *file.Source, err *CompilationError, offending *TokenNode) bool {
	if isNoAlternative(err.Message) {
		if is(offending.Prev(), ",") {
			span := spanFromTokenSafe(offending.Prev().Token(), src)
			span.Start++
			span.End++

			err.Message = "Expected expression after ','"
			err.Hint = "Did you forget to provide a value?"
			err.Spans = []ErrorSpan{
				NewMainErrorSpan(span, "missing value"),
			}

			return true
		}

		if is(offending.Prev(), "||") || is(offending.Prev(), "OR") ||
			is(offending.Prev(), "&&") || is(offending.Prev(), "AND") {
			span := spanFromTokenSafe(offending.Prev().Token(), src)
			span.Start += 2
			span.End += 2

			operator := offending.Prev().GetText()
			err.Message = fmt.Sprintf("Expected right-hand expression after '%s'", operator)
			err.Hint = fmt.Sprintf("Provide an expression after the logical operator, e.g. (a %s b).", operator)
			err.Spans = []ErrorSpan{
				NewMainErrorSpan(span, "missing expression"),
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
			err.Spans = []ErrorSpan{
				NewMainErrorSpan(span, "missing expression"),
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
			err.Spans = []ErrorSpan{
				NewMainErrorSpan(span, "missing expression"),
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
			err.Spans = []ErrorSpan{
				NewMainErrorSpan(span, "missing ')'"),
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
			err.Spans = []ErrorSpan{
				NewMainErrorSpan(span, "missing ')'"),
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
			err.Spans = []ErrorSpan{
				NewMainErrorSpan(span, "missing value"),
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
			err.Spans = []ErrorSpan{
				NewMainErrorSpan(span, "missing value"),
			}

			return true
		}
	}

	return false
}
