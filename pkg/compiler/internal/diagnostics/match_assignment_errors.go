package diagnostics

import (
	"fmt"

	"github.com/MontFerret/ferret/pkg/file"
)

func matchMissingAssignmentValue(src *file.Source, err *CompilationError, offending *TokenNode) bool {
	if isExtraneous(err.Message) {
		return false
	}

	prev := offending.Prev()

	if node := anyIs(offending, prev, "="); node != nil {
		span := spanFromTokenSafe(node.Token(), src)
		span.Start++
		span.End++
		err.Message = fmt.Sprintf("Expected expression after '=' for variable '%s'", node.Prev())
		err.Hint = "Did you forget to provide a value?"
		err.Spans = []ErrorSpan{
			NewMainErrorSpan(span, "missing value"),
		}

		return true
	}

	if is(offending, "LET") {
		span := spanFromTokenSafe(offending.Token(), src)
		span.Start = span.End
		span.End = span.Start + 1
		err.Message = "Expected variable name"
		err.Hint = "Did you forget to provide a variable name?"
		err.Spans = []ErrorSpan{
			NewMainErrorSpan(span, "missing name"),
		}

		return true
	}

	return false
}
