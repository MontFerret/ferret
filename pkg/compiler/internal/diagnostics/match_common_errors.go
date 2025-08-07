package diagnostics

import "github.com/MontFerret/ferret/pkg/file"

func matchCommonErrors(src *file.Source, err *CompilationError, offending *TokenNode) bool {
	if isNoAlternative(err.Message) {
		if is(offending.Prev(), ",") {
			span := spanFromTokenSafe(offending.Prev().Token(), src)
			span.Start++
			span.End++

			err.Message = "Expected expression after ','"
			err.Hint = "Did you forget to provide a value?"
			err.Spans = []ErrorSpan{
				NewMainErrorSpan(span, "missing value"),
			}

			return true
		}
	}

	return false
}
