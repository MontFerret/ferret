package diagnostics

import (
	"github.com/MontFerret/ferret/v2/pkg/diagnostics"
	"github.com/MontFerret/ferret/v2/pkg/source"
)

type SyntaxErrorMatcher func(src *source.Source, err *diagnostics.Diagnostic, offending *TokenNode) bool

func AnalyzeSyntaxError(src *source.Source, err *diagnostics.Diagnostic, offending *TokenNode) bool {
	matchers := []SyntaxErrorMatcher{
		matchArrayOperatorErrors,
		matchQueryErrors,
		matchWhileLoopErrors,
		matchDispatchErrors,
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
