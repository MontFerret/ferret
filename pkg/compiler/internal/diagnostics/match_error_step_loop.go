package diagnostics

import (
	"strings"

	"github.com/MontFerret/ferret/pkg/file"
)

func matchStepLoopErrors(src *file.Source, err *CompilationError, offending *TokenNode) bool {
	prev := offending.Prev()

	// Handle very specific case: "Incomplete STEP clause at end" - highest priority
	if isNoAlternative(err.Message) && strings.Contains(err.Message, "STEP\\n\\t\\t") {
		// Find the STEP token
		stepNode := findPrevToken(offending, "STEP", 5)
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

	// Handle "STEP followed by RETURN" case - expecting "Expected a RETURN or FOR clause at end of query"
	if isNoAlternative(err.Message) && has(err.Message, "STEP RETURN") {
		stepNode := findPrevToken(offending, "STEP", 5)
		if stepNode != nil {
			span := spanFromTokenSafe(stepNode.Token(), src)
			span.Start = span.End + 1
			span.End = span.Start + 1
			err.Message = "Expected a RETURN or FOR clause at end of query"
			err.Hint = "All queries must return a value. Add a RETURN statement to complete the query."
			err.Spans = []ErrorSpan{
				NewMainErrorSpan(span, "missing assignment"),
			}

			return true
		}
	}

	// Handle "Incomplete STEP clause" cases - first priority for exact STEP endings
	if isNoAlternative(err.Message) && has(err.Message, "STEP") {
		input := extractNoAlternativeInput(err.Message)
		// Handle cases like "FOR i = 0 WHILE i < 5 STEP\n\t\t" where STEP is at the end
		if strings.HasSuffix(strings.TrimSpace(input), "STEP") {
			stepNode := findPrevToken(offending, "STEP", 5)
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
		// Also handle cases that end with just the variable name after STEP
		tokens := strings.Fields(input)
		for i, token := range tokens {
			if token == "STEP" {
				// If STEP is last token or followed by only one more token before newline/end
				if i == len(tokens)-1 || (i+1 == len(tokens)-1 && !strings.Contains(tokens[i+1], "=")) {
					stepNode := findPrevToken(offending, "STEP", 5)
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
			}
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

	// Handle "Missing '=' in STEP assignment" - expecting "Expected '=' after variable in STEP clause"
	if isNoAlternative(err.Message) && has(err.Message, "STEP") && has(err.Message, "RETURN") {
		input := extractNoAlternativeInput(err.Message)
		tokens := strings.Fields(input)
		// Look for pattern like "STEP i RETURN" or "STEP variable RETURN"
		for i, token := range tokens {
			if token == "STEP" && i+2 < len(tokens) && tokens[i+2] == "RETURN" {
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
		}
	}

	// Handle "Missing STEP keyword" case - but not when it's actually missing initial value
	if isNoAlternative(err.Message) && has(err.Message, "WHILE") && has(err.Message, "i") && !has(err.Message, "STEP") {
		input := extractNoAlternativeInput(err.Message)
		// Pattern like "FOR i = 0 WHILE i < 5 i" (missing STEP)
		// But NOT pattern like "FOR i = WHILE i < 5 STEP i = i + 1 RETURN i" (missing initial value)
		if strings.Contains(input, "WHILE") && !strings.Contains(input, "STEP") && !strings.Contains(input, "= WHILE") {
			span := spanFromTokenSafe(offending.Token(), src)
			err.Message = "Syntax error: missing 'STEP' at 'i'"
			err.Hint = "STEP loops require the STEP keyword between WHILE and the increment expression."
			err.Spans = []ErrorSpan{
				NewMainErrorSpan(span, "missing 'STEP'"),
			}
			return true
		}
	}

	// Handle case where we have "STEP" but no variable after it (e.g., "STEP RETURN")
	if isNoAlternative(err.Message) && hasPrevToken(offending, "STEP", 2) {
		// Check if this looks like a STEP clause with missing parts
		input := extractNoAlternativeInput(err.Message)
		if input != "" {
			tokens := strings.Fields(input)
			if len(tokens) > 0 && (tokens[0] == "RETURN" || tokens[0] == "FILTER" || tokens[0] == "SORT" || tokens[0] == "LIMIT") {
				span := spanFromTokenSafe(findPrevToken(offending, "STEP", 2).Token(), src)
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
