package diagnostics

import (
	"fmt"

	"github.com/MontFerret/ferret/v2/pkg/diagnostics"
	"github.com/MontFerret/ferret/v2/pkg/file"
)

func matchMissingAssignmentValue(src *file.Source, err *diagnostics.Diagnostic, offending *TokenNode) bool {
	if isExtraneous(err.Message) {
		return false
	}

	prev := offending.Prev()

	if node := anyIs(offending, prev, "="); node != nil {
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

	if is(offending, "LET") {
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
