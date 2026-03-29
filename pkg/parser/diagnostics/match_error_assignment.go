package diagnostics

import (
	"fmt"

	"github.com/MontFerret/ferret/v2/pkg/diagnostics"
)

func matchMissingAssignmentValue(src *source.Source, err *diagnostics.Diagnostic, offending *TokenNode) bool {
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

	if node := anyAssignmentOperator(offending, prev); node != nil {
		if keyword := node.Prev(); is(keyword, "COLLECT") || is(keyword, "AGGREGATE") {
			return false
		}

		span := spanFromTokenSafe(node.Token(), src)
		span.Start = span.End
		span.End = span.Start + 1

		prevText := ""
		if node.Prev() != nil {
			prevText = node.Prev().GetText()
		}

		err.Message = fmt.Sprintf("Expected expression after '%s' for variable '%s'", node.GetText(), prevText)
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

func matchInvalidVarDiscard(src *source.Source, err *diagnostics.Diagnostic, offending *TokenNode) bool {
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

func matchAssignmentExpression(src *source.Source, err *diagnostics.Diagnostic, offending *TokenNode) bool {
	if offending == nil || (!isNoAlternative(err.Message) && !isMismatched(err.Message) && !isExtraneous(err.Message) && !isMissing(err.Message)) {
		return false
	}

	operator := findPrevAssignmentOperator(offending, 16)
	if operator == nil {
		return false
	}

	prev := operator.Prev()
	if prev == nil || !isIdentifier(prev) {
		return false
	}

	if operator != offending {
		if next := operator.Next(); next == nil || next == offending || isEOF(next) {
			return false
		}
	}

	if marker := operator.PrevAt(2); is(marker, "LET") || is(marker, "VAR") || is(marker, "COLLECT") || is(marker, "AGGREGATE") || is(marker, "FOR") {
		return false
	}

	span := spanFromTokenSafe(operator.Token(), src)

	err.Message = "Assignment is only allowed as a standalone statement"
	err.Hint = fmt.Sprintf("Move '%s %s ...' to its own statement. Assignment cannot be used inside expressions.", prev.GetText(), operator.GetText())
	err.Spans = []diagnostics.ErrorSpan{
		diagnostics.NewMainErrorSpan(span, "assignment is not an expression"),
	}

	return true
}

func anyAssignmentOperator(nodes ...*TokenNode) *TokenNode {
	for _, node := range nodes {
		if isAssignmentOperator(node) {
			return node
		}
	}

	return nil
}

func findPrevAssignmentOperator(node *TokenNode, steps int) *TokenNode {
	current := node
	for i := 0; i < steps && current != nil; i++ {
		if isAssignmentOperator(current) {
			return current
		}

		current = current.Prev()
	}

	return nil
}

func isAssignmentOperator(node *TokenNode) bool {
	return is(node, "=") || is(node, "+=") || is(node, "-=") || is(node, "*=") || is(node, "/=")
}
