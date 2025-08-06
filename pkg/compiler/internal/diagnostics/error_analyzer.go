package diagnostics

import (
	"github.com/MontFerret/ferret/pkg/file"
)

type SyntaxErrorMatcher func(src *file.Source, err *CompilationError, offending *TokenNode) bool

func AnalyzeSyntaxError(src *file.Source, err *CompilationError, offending *TokenNode) bool {
	matchers := []SyntaxErrorMatcher{
		matchMissingAssignmentValue,
		matchForLoopErrors,
		matchMissingReturnValue,
	}

	for _, matcher := range matchers {
		if matcher(src, err, offending) {
			return true
		}
	}

	return false
}
