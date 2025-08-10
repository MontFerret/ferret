package diagnostics

import (
	"fmt"

	"github.com/MontFerret/ferret/pkg/file"
)

func matchMissingReturnValue(src *file.Source, err *CompilationError, offending *TokenNode) bool {
	extraneous := isExtraneous(err.Message)

	if !is(offending, "RETURN") && !extraneous {
		return false
	}

	if extraneous {
		span := spanFromTokenSafe(offending.Token(), src)
		err.Spans = []ErrorSpan{
			NewMainErrorSpan(span, "query must end with a value"),
		}

		err.Message = "Expected a RETURN or FOR clause at end of query"
		err.Hint = "All queries must return a value. Add a RETURN statement to complete the query."

		return true
	}

	span := spanFromTokenSafe(offending.Token(), src)
	span.Start = span.End
	span.End = span.Start + 1
	err.Message = fmt.Sprintf("Expected expression after '%s'", offending)
	err.Hint = "Did you forget to provide a value to return?"
	err.Spans = []ErrorSpan{
		NewMainErrorSpan(span, "missing return value"),
	}

	return true
}
