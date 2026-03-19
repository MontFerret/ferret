package diagnostics

import (
	"fmt"

	"github.com/MontFerret/ferret/v2/pkg/diagnostics"
	"github.com/MontFerret/ferret/v2/pkg/file"
)

func matchMissingAssignmentValue(src *file.Source, err *diagnostics.Diagnostic, offending *TokenNode) bool {
	if matchInvalidVarDiscard(src, err, offending) {
		return true
	}

	if matchAssignmentExpression(src, err, offending) {
		return true
	}

	if isExtraneous(err.Message) {
		return false
	}

	prev := offending.Prev()

	if node := anyIs(offending, prev, "="); node != nil {
		if keyword := node.Prev(); is(keyword, "COLLECT") || is(keyword, "AGGREGATE") {
			return false
		}

		span := spanFromTokenSafe(node.Token(), src)
		span.Start++
		span.End++

		prevText := ""
		if node.Prev() != nil {
			prevText = node.Prev().GetText()
		}

		err.Message = fmt.Sprintf("Expected expression after '=' for variable '%s'", prevText)
		err.Hint = "Did you forget to provide a value?"
		err.Spans = []diagnostics.ErrorSpan{
			diagnostics.NewMainErrorSpan(span, "missing value"),
		}

		return true
	}

	if is(offending, "LET") || is(offending, "VAR") {
		span := spanFromTokenSafe(offending.Token(), src)
		span.Start = span.End
		span.End = span.Start + 1
		err.Message = "Expected variable name"
		err.Hint = "Did you forget to provide a variable name?"
		err.Spans = []diagnostics.ErrorSpan{
			diagnostics.NewMainErrorSpan(span, "missing name"),
		}

		return true
	}

	return false
}

func matchInvalidVarDiscard(src *file.Source, err *diagnostics.Diagnostic, offending *TokenNode) bool {
	if offending == nil {
		return false
	}

	target := offending
	if is(offending, "VAR") && is(offending.Next(), "_") {
		target = offending.Next()
	} else if !is(offending, "_") || !is(offending.Prev(), "VAR") {
		return false
	}

	span := spanFromTokenSafe(target.Token(), src)

	err.Message = "VAR cannot use '_' as a variable name"
	err.Hint = "Use a real variable name for VAR, or use LET _ = ... to explicitly discard a value."
	err.Spans = []diagnostics.ErrorSpan{
		diagnostics.NewMainErrorSpan(span, "invalid variable name"),
	}

	return true
}

func matchAssignmentExpression(src *file.Source, err *diagnostics.Diagnostic, offending *TokenNode) bool {
	if offending == nil || (!isNoAlternative(err.Message) && !isMismatched(err.Message) && !isExtraneous(err.Message) && !isMissing(err.Message)) {
		return false
	}

	eq := findPrevToken(offending, "=", 8)
	if eq == nil {
		return false
	}

	prev := eq.Prev()
	if prev == nil || !isIdentifier(prev) {
		return false
	}

	if marker := eq.PrevAt(2); is(marker, "LET") || is(marker, "VAR") || is(marker, "STEP") || is(marker, "COLLECT") || is(marker, "AGGREGATE") || is(marker, "FOR") {
		return false
	}

	span := spanFromTokenSafe(eq.Token(), src)
	span.Start++
	span.End++

	err.Message = "Assignment is only allowed as a standalone statement"
	err.Hint = fmt.Sprintf("Move '%s = ...' to its own statement. Assignment cannot be used inside expressions.", prev.GetText())
	err.Spans = []diagnostics.ErrorSpan{
		diagnostics.NewMainErrorSpan(span, "assignment is not an expression"),
	}

	return true
}
