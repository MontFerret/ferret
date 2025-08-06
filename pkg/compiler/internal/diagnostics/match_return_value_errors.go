package diagnostics

import (
	"fmt"

	"github.com/MontFerret/ferret/pkg/file"
)

func matchMissingReturnValue(src *file.Source, err *CompilationError, offending *TokenNode) bool {
	if !is(offending, "RETURN") {
		return false
	}

	span := spanFromTokenSafe(offending.Token(), src)
	err.Message = fmt.Sprintf("Expected expression after '%s'", offending)
	err.Hint = "Did you forget to provide a value to return?"
	err.Spans = []ErrorSpan{
		NewMainErrorSpan(span, "missing return value"),
	}

	return true
}
