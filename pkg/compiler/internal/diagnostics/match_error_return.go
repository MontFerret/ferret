package diagnostics

import (
	"fmt"

	"github.com/MontFerret/ferret/pkg/file"
)

func matchMissingReturnValue(src *file.Source, err *CompilationError, offending *TokenNode) bool {
	// Prefer range-specific error when the parser trips on an incomplete range like "0.. RETURN".
	if is(offending, "..") || is(offending.Prev(), "..") || has(err.Message, "..") {
		span := spanFromTokenSafe(offending.Token(), src)
		span.Start += 2
		span.End += 2

		start := ""
		if is(offending, "..") && offending.Prev() != nil {
			start = offending.Prev().GetText()
		} else if is(offending.Prev(), "..") && offending.Prev().Prev() != nil {
			start = offending.Prev().Prev().GetText()
		} else {
			start = extractRangeStart(err.Message)
		}

		err.Message = "Expected end value after '..' in range expression"
		err.Hint = fmt.Sprintf("Provide an end value to complete the range, e.g. %s..10.", start)
		err.Spans = []ErrorSpan{
			NewMainErrorSpan(span, "missing value"),
		}

		return true
	}

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
