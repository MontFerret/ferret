package diagnostics

import (
	"fmt"
	"strings"

	"github.com/MontFerret/ferret/pkg/file"
)

func matchLiteralErrors(src *file.Source, err *CompilationError, offending *TokenNode) bool {
	if isNoAlternative(err.Message) {
		input := extractNoAlternativeInputs(err.Message)
		token := input[len(input)-1]

		isMissingClosingQuote := isQuote(token)
		isMissingOpeningQuote := isKeyword(offending.Prev()) && isQuote(token[len(token)-1:]) && !isValidString(token)

		if isMissingClosingQuote || isMissingOpeningQuote {
			var span file.Span
			var typeOfQuote string
			var quote string

			if isMissingClosingQuote {
				quote = token
				typeOfQuote = "closing"
				span = spanFromTokenSafe(offending.Token(), src)
				inputRaw := extractNoAlternativeInput(err.Message)
				spaces := strings.Count(inputRaw, " ") + 1
				span.Start += spaces
				span.End += spaces
			} else {
				quote = token[len(token)-1:]
				typeOfQuote = "opening"
				span = spanFromTokenSafe(offending.Prev().Token(), src)
				span.Start += 2
				span.End += 2
			}

			err.Message = "Unclosed string literal"

			if quote == "'" {
				err.Hint = fmt.Sprintf("Add a matching \"%s\" to close the string.", quote)
				err.Spans = []ErrorSpan{
					NewMainErrorSpan(span, fmt.Sprintf("missing %s \"%s\"", typeOfQuote, quote)),
				}
			} else {
				err.Hint = fmt.Sprintf("Add a matching '%s' to close the string.", quote)
				err.Spans = []ErrorSpan{
					NewMainErrorSpan(span, fmt.Sprintf("missing %s '%s'", typeOfQuote, quote)),
				}
			}

			return true
		}
	}

	if isNoAlternative(err.Message) || isMissing(err.Message) {
		if is(offending.Prev(), "[") {
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

			if !isKeyword(offending.PrevAt(2)) {
				err.Message = "Unclosed computed property expression"
				err.Hint = "Add a closing ']' to complete the computed property expression."
				err.Spans = []ErrorSpan{
					NewMainErrorSpan(span, "missing ']'"),
				}
			} else {
				err.Message = "Unclosed array literal"
				err.Hint = "Add a closing ']' to complete the array."
				err.Spans = []ErrorSpan{
					NewMainErrorSpan(span, "missing ']'"),
				}
			}

			return true
		}

		if is(offending, "[") {
			span := spanFromTokenSafe(offending.Token(), src)
			span.Start++
			span.End++

			err.Message = "Expected a valid list of values"
			err.Hint = "Did you forget to provide a value?"
			err.Spans = []ErrorSpan{
				NewMainErrorSpan(span, "missing value"),
			}

			return true
		}

		if is(offending.Prev(), "{") {
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

			err.Message = "Unclosed object literal"
			err.Hint = "Add a closing '}' to complete the object."
			err.Spans = []ErrorSpan{
				NewMainErrorSpan(span, "missing '}'"),
			}

			return true
		}

		if is(offending, "{") && isNoAlternative(err.Message) {
			span := spanFromTokenSafe(offending.Token(), src)
			span.Start++
			span.End++

			err.Message = "Expected property name before ':'"
			err.Hint = "Object properties must have a name before the colon, e.g. { property: 123 }."
			err.Spans = []ErrorSpan{
				NewMainErrorSpan(span, "missing property name"),
			}

			return true
		}

		if is(offending, ":") && isIdentifier(offending.Prev()) {
			span := spanFromTokenSafe(offending.Token(), src)
			span.Start++
			span.End++
			property := offending.Prev().GetText()

			err.Message = "Expected value after object property name"
			err.Hint = fmt.Sprintf("Provide a value for the property, e.g. { %s: 123 }.", property)
			err.Spans = []ErrorSpan{
				NewMainErrorSpan(span, "missing value"),
			}

			return true
		}
	}

	if isExtraneous(err.Message) || isMismatched(err.Message) {
		var token string

		if isExtraneous(err.Message) {
			token = extractExtraneousInput(err.Message)
		} else {
			token = extractMismatchedInput(err.Message)
		}

		if isQuote(token) {
			var span file.Span
			var typeOfQuote string

			if isKeyword(offending) {
				span = spanFromTokenSafe(offending.Token(), src)
				typeOfQuote = "closing"
			} else {
				span = spanFromTokenSafe(offending.Prev().Token(), src)
				typeOfQuote = "opening"
			}

			span.Start += 2
			span.End += 2
			err.Message = "Unclosed string literal"

			if token == "'" {
				err.Hint = fmt.Sprintf("Add a matching \"%s\" to close the string.", token)
				err.Spans = []ErrorSpan{
					NewMainErrorSpan(span, fmt.Sprintf("missing %s \"%s\"", typeOfQuote, token)),
				}
			} else {
				err.Hint = fmt.Sprintf("Add a matching '%s' to close the string.", token)
				err.Spans = []ErrorSpan{
					NewMainErrorSpan(span, fmt.Sprintf("missing %s '%s'", typeOfQuote, token)),
				}
			}

			return true
		}

		if is(offending, "[") && token == "]" {
			span := spanFromTokenSafe(offending.Token(), src)
			span.End++

			val := offending.Prev().String()
			err.Message = "Expected expression inside computed property brackets"
			err.Hint = fmt.Sprintf("Provide a property key or index inside '[ ]', e.g. %s[0] or %s[\"key\"].", val, val)
			err.Spans = []ErrorSpan{
				NewMainErrorSpan(span, "missing expression"),
			}

			return true
		}
	}

	return false
}
