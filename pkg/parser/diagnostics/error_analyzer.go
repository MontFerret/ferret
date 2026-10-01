package diagnostics

import (
	"github.com/MontFerret/ferret/v2/pkg/diagnostics"
	"github.com/MontFerret/ferret/v2/pkg/source"
)

type SyntaxErrorMatcher func(src source.Source, err *diagnostics.Diagnostic, offending *TokenNode) bool

func AnalyzeSyntaxError(src source.Source, err *diagnostics.Diagnostic, offending *TokenNode) bool {
	return analyzeSyntaxError(src, err, offending, matchLiteralErrors)
}

func analyzeSyntaxError(src source.Source, err *diagnostics.Diagnostic, offending *TokenNode, literals SyntaxErrorMatcher) bool {
	matchers := []SyntaxErrorMatcher{
		matchCoalesceErrors,
		matchArrayOperatorErrors,
		matchErrorPolicyErrors,
		matchQueryErrors,
		matchWaitForErrors,
		matchWhileLoopErrors,
		matchDispatchErrors,
		matchSpreadEntryErrors,
		matchArrayLiteralSeparatorErrors,
		matchMissingFunctionParamsClose,
		matchMixedFunctionBodySyntax,
		literals,
		matchMissingAssignmentValue,
		matchDeleteStatementErrors,
		matchForLoopErrors,
		matchMatchArmSeparatorErrors,
		matchMissingReturnDistinctValue,
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
