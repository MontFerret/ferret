package diagnostics

import (
	"github.com/MontFerret/ferret/v2/pkg/diagnostics"
	"github.com/MontFerret/ferret/v2/pkg/source"
)

func matchErrorPolicyErrors(src *source.Source, err *diagnostics.Diagnostic, offending *TokenNode) bool {
	if err == nil || offending == nil {
		return false
	}

	if !isMismatched(err.Message) && !isMissing(err.Message) && !isNoAlternative(err.Message) && !isExtraneous(err.Message) && !isErrorPolicyPredicateFailure(err.Message) {
		return false
	}

	if errorNode := errorPolicyMissingActionNode(offending); errorNode != nil {
		span := spanFromTokenSafe(errorNode.Token(), src)
		err.Message = "Expected SUPPRESS or THROW after 'ON ERROR'"
		err.Hint = "Use ON ERROR SUPPRESS to swallow failures or ON ERROR THROW to propagate them."
		err.Spans = []diagnostics.ErrorSpan{
			diagnostics.NewMainErrorSpan(span, "missing error policy"),
		}
		return true
	}

	if on := errorPolicyMissingErrorNode(offending); on != nil {
		span := spanFromTokenSafe(on.Token(), src)
		err.Message = "Expected ERROR after 'ON' in error policy tail"
		err.Hint = "Complete the tail as ON ERROR SUPPRESS or ON ERROR THROW."
		err.Spans = []diagnostics.ErrorSpan{
			diagnostics.NewMainErrorSpan(span, "missing ERROR"),
		}
		return true
	}

	return false
}

func isErrorPolicyPredicateFailure(msg string) bool {
	return has(msg, "errorkeyword failed predicate") || has(msg, "suppresskeyword failed predicate")
}

func errorPolicyMissingErrorNode(offending *TokenNode) *TokenNode {
	if offending == nil {
		return nil
	}

	if is(offending, "<EOF>") && is(offending.Prev(), "ON") {
		return offending.Prev()
	}

	if is(offending, "ON") {
		next := offending.Next()
		if next == nil || is(next, "<EOF>") {
			return offending
		}
	}

	if prev := offending.Prev(); prev != nil && is(prev, "ON") {
		return prev
	}

	return nil
}

func errorPolicyMissingActionNode(offending *TokenNode) *TokenNode {
	if offending == nil {
		return nil
	}

	if is(offending, "<EOF>") && is(offending.Prev(), "ERROR") && hasPrevToken(offending, "ON", 3) {
		return offending.Prev()
	}

	if prev := offending.Prev(); prev != nil && is(prev, "ERROR") && hasPrevToken(prev, "ON", 3) {
		return prev
	}

	return nil
}
