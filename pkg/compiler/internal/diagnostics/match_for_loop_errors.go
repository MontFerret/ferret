package diagnostics

import (
	"github.com/MontFerret/ferret/pkg/file"
)

func matchForLoopErrors(src *file.Source, err *CompilationError, offending *TokenNode) bool {
	prev := offending.Prev()

	if is(prev, "IN") {
		span := spanFromTokenSafe(prev.Token(), src)
		span.Start = span.End + 1
		span.End = span.Start + 1
		err.Message = "Expected expression after 'IN'"
		err.Hint = "Each FOR loop must iterate over a collection or range."
		err.Spans = []ErrorSpan{
			NewMainErrorSpan(span, "missing value"),
		}

		return true
	}

	if is(prev, "FOR") {
		span := spanFromTokenSafe(offending.Token(), src)
		span.Start = span.End
		span.End = span.Start + 1
		err.Message = "Expected 'IN' after loop variable"
		err.Hint = "Use 'FOR x IN [iterable]' syntax."
		err.Spans = []ErrorSpan{
			NewMainErrorSpan(span, "missing keyword"),
		}

		return true
	}

	if is(offending, "FOR") {
		span := spanFromTokenSafe(offending.Token(), src)
		span.Start = span.End
		span.End = span.Start + 1
		err.Message = "Expected loop variable before 'IN'"
		err.Hint = "FOR must declare a variable."
		err.Spans = []ErrorSpan{
			NewMainErrorSpan(span, "missing variable"),
		}

		return true
	}

	return false
}
