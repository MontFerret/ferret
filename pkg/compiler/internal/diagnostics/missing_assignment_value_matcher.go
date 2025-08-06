package diagnostics

import (
	"fmt"
	"github.com/MontFerret/ferret/pkg/file"
)

func missingAssignmentValueMatcher(src *file.Source, err *CompilationError, offending *TokenNode) bool {
	prev := offending.Prev()

	// CASE: LET x = [missing value]
	if is(offending, "LET") || is(prev, "=") {
		span := spanFromTokenSafe(prev.Token(), src)
		span.Start++
		span.End++
		err.Message = fmt.Sprintf("Expected expression after '=' for variable '%s'", prev.Prev())
		err.Hint = "Did you forget to provide a value?"
		err.Spans = []ErrorSpan{
			NewMainErrorSpan(span, "missing value"),
		}

		return true
	}

	return false
}
