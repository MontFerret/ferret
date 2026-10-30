package diagnostics

import (
	"github.com/antlr4-go/antlr/v4"

	"github.com/MontFerret/ferret/v2/pkg/diagnostics"
	"github.com/MontFerret/ferret/v2/pkg/source"
)

type SyntaxErrorMatcher func(src source.Source, err *diagnostics.Diagnostic, offending *TokenNode) bool

func AnalyzeSyntaxError(src source.Source, err *diagnostics.Diagnostic, offending *TokenNode) bool {
	return analyzeSyntaxError(src, err, offending, nil, matchLiteralErrors)
}

func analyzeSyntaxError(src source.Source, err *diagnostics.Diagnostic, offending *TokenNode, stream antlr.TokenStream, literals SyntaxErrorMatcher) bool {
	matchers := []SyntaxErrorMatcher{
		matchCoalesceErrors,
		matchArrayOperatorErrors,
		matchErrorPolicyErrors,
		func(src source.Source, err *diagnostics.Diagnostic, offending *TokenNode) bool {
			return matchQueryErrors(src, err, offending, stream)
		},
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
