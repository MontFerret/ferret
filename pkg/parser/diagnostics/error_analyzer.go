package diagnostics

import (
	"github.com/MontFerret/ferret/v2/pkg/diagnostics"
	"github.com/MontFerret/ferret/v2/pkg/file"
)

type SyntaxErrorMatcher func(src *file.Source, err *diagnostics.Diagnostic, offending *TokenNode) bool

func AnalyzeSyntaxError(src *file.Source, err *diagnostics.Diagnostic, offending *TokenNode) bool {
	matchers := []SyntaxErrorMatcher{
		matchArrayOperatorErrors,
		matchQueryErrors,
		matchWhileLoopErrors,
		matchLiteralErrors,
		matchWaitForErrors,
		matchMissingAssignmentValue,
		matchForLoopErrors,
		matchCommonErrors,
		matchMissingReturnValue,
	}

	for _, matcher := range matchers {
		if matcher(src, err, offending) {
			return true
		}
	}

	return false
}
