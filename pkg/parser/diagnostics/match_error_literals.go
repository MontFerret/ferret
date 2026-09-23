package diagnostics

import (
	"fmt"

	"github.com/MontFerret/ferret/v2/pkg/diagnostics"
	"github.com/MontFerret/ferret/v2/pkg/parser/fql"
	"github.com/MontFerret/ferret/v2/pkg/source"
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

func matchLiteralErrors(src source.Source, err *diagnostics.Diagnostic, offending *TokenNode) bool {
	if isUnclosedTemplateLiteral(offending) {
		span := insertionSpanAfterToken(offending.Token(), src)

		err.Message = "Unclosed string literal"
		err.Hint = "Add a matching '`' to close the string."
		err.Spans = []diagnostics.ErrorSpan{
			diagnostics.NewMainErrorSpan(span, "missing closing '`'"),
		}

		return true
	}

	if isMissing(err.Message) && has(err.Message, "`") {
		span := insertionSpanAfterToken(offending.Token(), src)

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
			var span source.Span
			var typeOfQuote string
			var quote string

			if isMissingClosingQuote {
				quote = token
				typeOfQuote = "closing"
				span = quotedLiteralInsertionSpan(src, err, offending, false)
			} else {
				quote = token[len(token)-1:]
				typeOfQuote = "opening"
				span = quotedLiteralInsertionSpan(src, err, offending, true)
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
		if is(offending.Prev(), "[") {
			var span source.Span

			if isKeyword(offending) {
				span = insertionSpanAfterToken(offending.Prev().Token(), src)
			} else {
				span = insertionSpanAfterToken(offending.Token(), src)
			}

			if next := offending.Next(); is(next, "[") {
				span = insertionSpanAfterToken(next.Token(), src)
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

		if is(offending, "[") || extractNoAlternativeInput(err.Message) == "[" {
			span := insertionSpanAfterToken(offending.Token(), src)
			if next := offending.Next(); is(next, "[") {
				span = insertionSpanAfterToken(next.Token(), src)
			}

			if isComputedPropertyPrefix(offending.Prev()) {
				err.Message = "Unclosed computed property expression"
				err.Hint = "Add a closing ']' to complete the computed property expression."
				err.Spans = []diagnostics.ErrorSpan{
					diagnostics.NewMainErrorSpan(span, "missing ']'"),
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
			var span source.Span

			if isKeyword(offending) {
				span = insertionSpanAfterToken(offending.Prev().Token(), src)
			} else {
				span = insertionSpanAfterToken(offending.Token(), src)
			}

			err.Message = "Unclosed object literal"
			err.Hint = "Add a closing '}' to complete the object."
			err.Spans = []diagnostics.ErrorSpan{
				diagnostics.NewMainErrorSpan(span, "missing '}'"),
			}

			return true
		}

		if is(offending, "{") && isNoAlternative(err.Message) {
			span := insertionSpanAfterToken(offending.Token(), src)
			if next := offending.Next(); is(next, ":") {
				span = spanFromTokenSafe(next.Token(), src)
				span.End = span.Start
			}

			err.Message = "Expected property name before ':'"
			err.Hint = "Object properties must have a name before the colon, e.g. { property: 123 }."
			err.Spans = []diagnostics.ErrorSpan{
				diagnostics.NewMainErrorSpan(span, "missing property name"),
			}

			return true
		}

		if is(offending, ":") && isIdentifier(offending.Prev()) {
			span := insertionSpanAfterToken(offending.Token(), src)
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
			var span source.Span
			var typeOfQuote string

			if isKeyword(offending) {
				span = quotedLiteralInsertionSpan(src, err, offending, false)
				typeOfQuote = "closing"
			} else {
				span = quotedLiteralInsertionSpan(src, err, offending, true)
				typeOfQuote = "opening"
			}

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
			span := insertionSpanAfterToken(offending.Token(), src)

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

func quotedLiteralInsertionSpan(src source.Source, err *diagnostics.Diagnostic, offending *TokenNode, opening bool) source.Span {
	text := []rune(src.Content())
	valid := func(span source.Span) bool {
		return span.Start >= 0 && span.Start <= span.End && span.End <= len(text)
	}

	// Recovery metadata is optional and remains in ANTLR character coordinates.
	// Reject unusable anchors before scanning instead of repairing their offsets.
	var span source.Span
	if err != nil && len(err.Spans) > 0 && valid(err.Spans[0].Span) {
		span = err.Spans[0].Span
	} else if offending != nil && offending.Token() != nil {
		span = SpanFromToken(offending.Token())
		if !valid(span) {
			return source.Span{}
		}
	} else {
		return source.Span{}
	}

	if opening {
		tokens := lexDefaultTokens(src.Content())
		idx := findDiagnosticSpanTokenIndex(tokens, err)
		if idx >= 0 && idx < len(tokens) {
			for i := idx - 1; i >= 0 && tokens[i].GetTokenType() == fql.FqlLexerIdentifier; i-- {
				if tokens[i].GetLine() != tokens[idx].GetLine() {
					break
				}

				span.Start = tokens[i].GetStart()
			}
		}

		span.End = span.Start

		return span
	}

	// Ordinary quoted strings end at the line boundary. Count ANTLR characters,
	// not bytes or words from the human-readable parser error.
	end := span.End
	for end < len(text) && text[end] != '\n' && text[end] != '\r' {
		end++
	}

	return source.Span{Start: end, End: end}
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
