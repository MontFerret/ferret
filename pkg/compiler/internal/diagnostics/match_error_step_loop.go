package diagnostics

import (
	"regexp"
	"strings"

	"github.com/MontFerret/ferret/pkg/file"
)

func matchStepLoopErrors(src *file.Source, err *CompilationError, offending *TokenNode) bool {
	prev := offending.Prev()

	// Handle "WHILE STEP" case - missing condition after WHILE
	if is(offending, "STEP") && is(prev, "WHILE") {
		span := spanFromTokenSafe(prev.Token(), src)
		span.Start = span.End + 1
		span.End = span.Start + 1
		err.Message = "Expected expression after 'WHILE'"
		err.Hint = "STEP loops require a condition after WHILE, e.g., 'FOR i = 0 WHILE i < 10 STEP i = i + 1'."
		err.Spans = []ErrorSpan{
			NewMainErrorSpan(span, "missing condition"),
		}

		return true
	}

	// Handle incomplete STEP clause - missing variable assignment after STEP
	if is(prev, "STEP") && !isIdentifier(offending) {
		span := spanFromTokenSafe(prev.Token(), src)
		span.Start = span.End + 1
		span.End = span.Start + 1
		err.Message = "Expected variable assignment after 'STEP'"
		err.Hint = "STEP requires a variable assignment, e.g., 'STEP i = i + 1'."
		err.Spans = []ErrorSpan{
			NewMainErrorSpan(span, "missing assignment"),
		}

		return true
	}

	// Handle "Expected a RETURN or FOR clause at end of query" when preceded by STEP
	if has(err.Message, "Expected a RETURN or FOR clause at end of query") && hasPrevToken(offending, "STEP", 5) {
		stepNode := findPrevToken(offending, "STEP", 5)
		if stepNode != nil {
			span := spanFromTokenSafe(stepNode.Token(), src)
			span.Start = span.End + 1
			span.End = span.Start + 1
			err.Message = "Expected variable assignment after 'STEP'"
			err.Hint = "STEP requires a variable assignment, e.g., 'STEP i = i + 1'."
			err.Spans = []ErrorSpan{
				NewMainErrorSpan(span, "missing assignment"),
			}

			return true
		}
	}

	// Handle "mismatched input '<EOF>' expecting '='" in STEP context
	if isMismatched(err.Message) && has(err.Message, "expecting '='") && hasPrevToken(offending, "STEP", 3) {
		stepNode := findPrevToken(offending, "STEP", 3)
		if stepNode != nil {
			span := spanFromTokenSafe(stepNode.Token(), src)
			span.Start = span.End + 1
			span.End = span.Start + 1
			err.Message = "Incomplete STEP clause"
			err.Hint = "STEP requires a complete variable assignment, e.g., 'STEP i = i + 1'."
			err.Spans = []ErrorSpan{
				NewMainErrorSpan(span, "incomplete assignment"),
			}

			return true
		}
	}

	// Handle extraneous RETURN after STEP (and similar cases)
	if isExtraneous(err.Message) && is(offending, "RETURN") && hasPrevToken(offending, "STEP", 2) {
		stepNode := findPrevToken(offending, "STEP", 2)
		if stepNode != nil {
			span := spanFromTokenSafe(stepNode.Token(), src)
			span.Start = span.End + 1
			span.End = span.Start + 1
			err.Message = "Expected variable assignment after 'STEP'"
			err.Hint = "STEP requires a variable assignment, e.g., 'STEP i = i + 1'."
			err.Spans = []ErrorSpan{
				NewMainErrorSpan(span, "missing assignment"),
			}

			return true
		}
	}

	// Handle missing '=' in STEP assignment using multiple error patterns
	if isMissing(err.Message) && isMissingToken(err.Message, "=") && hasPrevToken(offending, "STEP", 3) {
		span := spanFromTokenSafe(offending.Token(), src)
		span.Start = span.End
		span.End = span.Start + 1
		err.Message = "Expected '=' after variable in STEP clause"
		err.Hint = "STEP assignments require '=', e.g., 'STEP i = i + 1'."
		err.Spans = []ErrorSpan{
			NewMainErrorSpan(span, "missing '='"),
		}

		return true
	}

	// Handle "no viable alternative" patterns for STEP loops
	if isNoAlternative(err.Message) {
		tokens := extractNoAlternativeInputs(err.Message)
		if len(tokens) == 0 {
			return false
		}

		upper := toUpperTokens(tokens)

		// Missing STEP keyword: "... WHILE <cond> <ident> = ..."
		if !containsToken(upper, "STEP") && containsToken(upper, "WHILE") {
			// Guard against missing initial assignment value: "FOR i = WHILE ..."
			if hasForEqualsWhilePattern(upper) {
				return false
			}

			last := normalizeToken(tokens[len(tokens)-1])
			if isIdentifierToken(last) {
				span := spanFromTokenSafe(offending.Token(), src)
				span.Start = span.End
				span.End = span.Start + 1
				err.Message = "Syntax error: missing 'STEP' at '" + last + "'"
				err.Spans = []ErrorSpan{
					NewMainErrorSpan(span, "missing 'STEP'"),
				}

				return true
			}
		}

		stepIdx := lastIndexOfToken(upper, "STEP")
		if stepIdx != -1 {
			after := upper[stepIdx+1:]
			rawAfter := tokens[stepIdx+1:]
			stepNode := findPrevToken(offending, "STEP", 6)

			span := spanFromTokenSafe(offending.Token(), src)
			if stepNode != nil {
				span = spanFromTokenSafe(stepNode.Token(), src)
				span.Start = span.End + 1
				span.End = span.Start + 1
			}

			if len(after) == 0 {
				err.Message = "Incomplete STEP clause"
				err.Hint = "STEP requires a complete variable assignment, e.g., 'STEP i = i + 1'."
				err.Spans = []ErrorSpan{
					NewMainErrorSpan(span, "incomplete assignment"),
				}

				return true
			}

			switch after[0] {
			case "RETURN", "FILTER", "SORT", "LIMIT":
				err.Message = "Expected a RETURN or FOR clause at end of query"
				err.Hint = "All queries must return a value. Add a RETURN statement to complete the query."
				err.Spans = []ErrorSpan{
					NewMainErrorSpan(span, "unexpected clause"),
				}

				return true
			}

			if isIdentifierToken(rawAfter[0]) {
				if len(after) == 1 {
					err.Message = "Incomplete STEP clause"
					err.Hint = "STEP requires a complete variable assignment, e.g., 'STEP i = i + 1'."
					err.Spans = []ErrorSpan{
						NewMainErrorSpan(span, "incomplete assignment"),
					}

					return true
				}

				if after[1] != "=" {
					err.Message = "Expected '=' after variable in STEP clause"
					err.Hint = "STEP assignments require '=', e.g., 'STEP i = i + 1'."
					err.Spans = []ErrorSpan{
						NewMainErrorSpan(span, "missing '='"),
					}

					return true
				}
			}
		}
	}

	return false
}

// Helper function to check if a previous token exists within n steps
func hasPrevToken(node *TokenNode, tokenText string, steps int) bool {
	current := node
	for i := 0; i < steps && current != nil; i++ {
		if is(current, tokenText) {
			return true
		}
		current = current.Prev()
	}
	return false
}

// Helper function to find a previous token within n steps
func findPrevToken(node *TokenNode, tokenText string, steps int) *TokenNode {
	current := node
	for i := 0; i < steps && current != nil; i++ {
		if is(current, tokenText) {
			return current
		}
		current = current.Prev()
	}
	return nil
}

func toUpperTokens(tokens []string) []string {
	upper := make([]string, len(tokens))
	for i, t := range tokens {
		upper[i] = strings.ToUpper(normalizeToken(t))
	}

	return upper
}

func lastIndexOfToken(tokens []string, token string) int {
	needle := strings.ToUpper(token)
	for i := len(tokens) - 1; i >= 0; i-- {
		if tokens[i] == needle {
			return i
		}
	}

	return -1
}

func containsToken(tokens []string, token string) bool {
	return lastIndexOfToken(tokens, token) >= 0
}

func hasForEqualsWhilePattern(tokens []string) bool {
	for i := 0; i+3 < len(tokens); i++ {
		if tokens[i] == "FOR" && tokens[i+2] == "=" && tokens[i+3] == "WHILE" {
			return true
		}
	}

	return false
}

var identifierToken = regexp.MustCompile(`^[A-Za-z_][A-Za-z0-9_]*$`)

func isIdentifierToken(token string) bool {
	return identifierToken.MatchString(normalizeToken(token))
}

func normalizeToken(token string) string {
	if token == "" {
		return token
	}

	// Remove escaped whitespace markers that appear in error messages.
	token = strings.ReplaceAll(token, `\n`, "")
	token = strings.ReplaceAll(token, `\t`, "")
	token = strings.ReplaceAll(token, `\r`, "")

	return token
}
