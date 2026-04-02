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
		if is(errorNode, "TIMEOUT") {
			err.Message = "Expected FAIL or RETURN after 'ON TIMEOUT'"
			err.Hint = "Use ON TIMEOUT FAIL to propagate timeout expiration or ON TIMEOUT RETURN <expr> to supply a fallback value."
		} else {
			err.Message = "Expected FAIL or RETURN after 'ON ERROR'"
			err.Hint = "Use ON ERROR FAIL to propagate failures or ON ERROR RETURN <expr> to supply a fallback value."
		}
		err.Spans = []diagnostics.ErrorSpan{
			diagnostics.NewMainErrorSpan(span, "missing recovery action"),
		}
		return true
	}

	if on := errorPolicyMissingErrorNode(offending); on != nil {
		span := spanFromTokenSafe(on.Token(), src)
		err.Message = "Expected ERROR or TIMEOUT after 'ON' in recovery tail"
		err.Hint = "Complete the tail as ON ERROR FAIL, ON ERROR RETURN <expr>, ON TIMEOUT FAIL, or ON TIMEOUT RETURN <expr>."
		err.Spans = []diagnostics.ErrorSpan{
			diagnostics.NewMainErrorSpan(span, "missing recovery condition"),
		}
		return true
	}

	return false
}

func isErrorPolicyPredicateFailure(msg string) bool {
	return has(msg, "errorkeyword failed predicate") ||
		has(msg, "timeoutkeyword failed predicate") ||
		has(msg, "failkeyword failed predicate") ||
		has(msg, "returnkeyword failed predicate")
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

	if is(offending, "<EOF>") && (is(offending.Prev(), "ERROR") || is(offending.Prev(), "TIMEOUT")) && hasPrevToken(offending, "ON", 3) {
		return offending.Prev()
	}

	if prev := offending.Prev(); prev != nil && (is(prev, "ERROR") || is(prev, "TIMEOUT")) && hasPrevToken(prev, "ON", 3) {
		return prev
	}

	return nil
}
