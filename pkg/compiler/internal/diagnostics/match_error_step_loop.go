package diagnostics

import (
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
	// BUT: do NOT convert this message for "STEP RETURN" patterns, as they should keep the general message
	if has(err.Message, "Expected a RETURN or FOR clause at end of query") && hasPrevToken(offending, "STEP", 5) {
		// Extract the input to check if this is a "STEP RETURN" pattern
		if isNoAlternative(err.Message) {
			input := extractNoAlternativeInput(err.Message)
			if has(input, "STEP RETURN") {
				return false // Keep the original message for "STEP RETURN" patterns
			}
		}
		
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
