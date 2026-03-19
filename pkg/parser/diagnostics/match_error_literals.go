package diagnostics

import (
	"fmt"
	"strings"

	"github.com/MontFerret/ferret/v2/pkg/diagnostics"
	"github.com/MontFerret/ferret/v2/pkg/file"
	"github.com/MontFerret/ferret/v2/pkg/parser/fql"
)

func isComputedPropertyPrefix(node *TokenNode) bool {
	if node == nil {
		return false
	}

	// Common cases: identifiers or closing delimiters of an expression.
	if isIdentifier(node) || is(node, "]") || is(node, ")") || is(node, "}") {
		return true
	}

	// Literals (numbers/strings) are not keywords.
	if !isKeyword(node) {
		return true
	}

	return false
}

func matchLiteralErrors(src *file.Source, err *diagnostics.Diagnostic, offending *TokenNode) bool {
	if isUnclosedTemplateLiteral(offending) {
		span := spanFromTokenSafe(offending.Token(), src)
		span.Start = span.End
		span.End = span.Start + 1

		err.Message = "Unclosed string literal"
		err.Hint = "Add a matching '`' to close the string."
		err.Spans = []diagnostics.ErrorSpan{
			diagnostics.NewMainErrorSpan(span, "missing closing '`'"),
		}

		return true
	}

	if isMissing(err.Message) && has(err.Message, "`") {
		span := spanFromTokenSafe(offending.Token(), src)
		span.Start = span.End
		span.End = span.Start + 1

		err.Message = "Unclosed string literal"
		err.Hint = "Add a matching '`' to close the string."
		err.Spans = []diagnostics.ErrorSpan{
			diagnostics.NewMainErrorSpan(span, "missing closing '`'"),
		}

		return true
	}

	if isNoAlternative(err.Message) {
		input := extractNoAlternativeInputs(err.Message)
		if len(input) == 0 {
			return false
		}

		token := input[len(input)-1]

		isMissingClosingQuote := isQuote(token)
		isMissingOpeningQuote := len(token) > 0 && isQuote(token[len(token)-1:]) && !isValidString(token)

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
				span = spanFromTokenSafe(offending.Token(), src)
				span.Start = span.End
				span.End = span.Start + 1
			}

			err.Message = "Unclosed string literal"

			if quote == "'" {
				err.Hint = fmt.Sprintf("Add a matching \"%s\" to close the string.", quote)
				err.Spans = []diagnostics.ErrorSpan{
					diagnostics.NewMainErrorSpan(span, fmt.Sprintf("missing %s \"%s\"", typeOfQuote, quote)),
				}
			} else {
				err.Hint = fmt.Sprintf("Add a matching '%s' to close the string.", quote)
				err.Spans = []diagnostics.ErrorSpan{
					diagnostics.NewMainErrorSpan(span, fmt.Sprintf("missing %s '%s'", typeOfQuote, quote)),
				}
			}

			return true
		}
	}

	if isNoAlternative(err.Message) || isMissing(err.Message) || isMismatched(err.Message) {
		if is(offending, "RETURN") {
			if bracket := findPrevToken(offending, "[", 6); bracket != nil && isComputedPropertyPrefix(bracket.Prev()) {
				span := spanFromTokenSafe(offending.Token(), src)
				span.Start = span.End
				span.End = span.Start + 1

				err.Message = "Expected a RETURN or FOR clause at end of query"
				err.Hint = "All queries must return a value. Add a RETURN statement to complete the query."
				err.Spans = []diagnostics.ErrorSpan{
					diagnostics.NewMainErrorSpan(span, "incomplete expression"),
				}

				return true
			}
		}

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
				err.Spans = []diagnostics.ErrorSpan{
					diagnostics.NewMainErrorSpan(span, "missing ']'"),
				}
			} else {
				err.Message = "Unclosed array literal"
				err.Hint = "Add a closing ']' to complete the array."
				err.Spans = []diagnostics.ErrorSpan{
					diagnostics.NewMainErrorSpan(span, "missing ']'"),
				}
			}

			return true
		}

		if is(offending, "[") {
			span := spanFromTokenSafe(offending.Token(), src)
			span.Start++
			span.End++

			if isComputedPropertyPrefix(offending.Prev()) {
				val := offending.Prev().String()
				err.Message = "Expected expression inside computed property brackets"
				err.Hint = fmt.Sprintf("Provide a property key or index inside '[ ]', e.g. %s[0] or %s[\"key\"].", val, val)
				err.Spans = []diagnostics.ErrorSpan{
					diagnostics.NewMainErrorSpan(span, "missing expression"),
				}

				return true
			}

			err.Message = "Expected a valid list of values"
			err.Hint = "Did you forget to provide a value?"
			err.Spans = []diagnostics.ErrorSpan{
				diagnostics.NewMainErrorSpan(span, "missing value"),
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
			err.Spans = []diagnostics.ErrorSpan{
				diagnostics.NewMainErrorSpan(span, "missing '}'"),
			}

			return true
		}

		if is(offending, "{") && isNoAlternative(err.Message) {
			span := spanFromTokenSafe(offending.Token(), src)
			span.Start++
			span.End++

			err.Message = "Expected property name before ':'"
			err.Hint = "Object properties must have a name before the colon, e.g. { property: 123 }."
			err.Spans = []diagnostics.ErrorSpan{
				diagnostics.NewMainErrorSpan(span, "missing property name"),
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
			err.Spans = []diagnostics.ErrorSpan{
				diagnostics.NewMainErrorSpan(span, "missing value"),
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
				err.Spans = []diagnostics.ErrorSpan{
					diagnostics.NewMainErrorSpan(span, fmt.Sprintf("missing %s \"%s\"", typeOfQuote, token)),
				}
			} else {
				err.Hint = fmt.Sprintf("Add a matching '%s' to close the string.", token)
				err.Spans = []diagnostics.ErrorSpan{
					diagnostics.NewMainErrorSpan(span, fmt.Sprintf("missing %s '%s'", typeOfQuote, token)),
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
			err.Spans = []diagnostics.ErrorSpan{
				diagnostics.NewMainErrorSpan(span, "missing expression"),
			}

			return true
		}
	}

	return false
}

func isUnclosedTemplateLiteral(node *TokenNode) bool {
	if node == nil || node.Token() == nil {
		return false
	}

	// If we recently saw a closing backtick, don't treat it as unclosed.
	if hasPrevTokenType(node, fql.FqlLexerBacktickClose, 12) {
		return false
	}

	ttype := node.Token().GetTokenType()
	if ttype == fql.FqlLexerBacktickOpen || ttype == fql.FqlLexerTemplateChars || ttype == fql.FqlLexerTemplateExprStart || ttype == fql.FqlLexerTemplateExprEnd {
		return true
	}

	return hasPrevTokenType(node, fql.FqlLexerBacktickOpen, 12)
}

func hasPrevTokenType(node *TokenNode, tokenType int, steps int) bool {
	current := node
	for i := 0; i < steps && current != nil; i++ {
		if current.Token() != nil && current.Token().GetTokenType() == tokenType {
			return true
		}
		current = current.Prev()
	}
	return false
}
