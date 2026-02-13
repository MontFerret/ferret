package diagnostics

import (
	"fmt"
	"strings"

	"github.com/MontFerret/ferret/v2/pkg/diagnostics"
	"github.com/MontFerret/ferret/v2/pkg/file"
)

func matchWaitForErrors(src *file.Source, err *CompilationError, offending *TokenNode) bool {
	if err == nil || offending == nil {
		return false
	}

	if has(err.Message, "waitforpredicate failed predicate") {
		if keyword, spanNode := waitForPredicateKeyword(offending); keyword != "" {
			span := spanFromTokenSafe(spanNode.Token(), src)
			err.Message = fmt.Sprintf("Expected expression after '%s' in WAITFOR predicate", keyword)
			err.Hint = fmt.Sprintf("Provide an expression after %s, e.g. WAITFOR %s x.", keyword, keyword)
			err.Spans = []diagnostics.ErrorSpan{
				diagnostics.NewMainErrorSpan(span, "missing expression"),
			}
			return true
		}
	}

	if keyword, spanNode := waitForPredicateKeyword(offending); keyword != "" {
		if is(offending, "RETURN") || isMissing(err.Message) || isNoAlternative(err.Message) {
			span := spanFromTokenSafe(spanNode.Token(), src)
			err.Message = fmt.Sprintf("Expected expression after '%s' in WAITFOR predicate", keyword)
			err.Hint = fmt.Sprintf("Provide an expression after %s, e.g. WAITFOR %s x.", keyword, keyword)
			err.Spans = []diagnostics.ErrorSpan{
				diagnostics.NewMainErrorSpan(span, "missing expression"),
			}
			return true
		}
	}

	if clause, spanNode := waitForMissingClauseValue(offending); clause != "" {
		span := spanFromTokenSafe(spanNode.Token(), src)
		err.Message = fmt.Sprintf("Expected value after '%s' in WAITFOR clause", clause)
		if clause == "BACKOFF" {
			err.Hint = "Provide a backoff strategy, e.g. BACKOFF LINEAR."
		} else if clause == "JITTER" {
			err.Hint = "Provide a jitter value between 0 and 1, e.g. JITTER 0.2."
		} else {
			err.Hint = fmt.Sprintf("Provide a duration, e.g. %s 100ms.", clause)
		}
		err.Spans = []diagnostics.ErrorSpan{
			diagnostics.NewMainErrorSpan(span, "missing value"),
		}
		return true
	}

	return false
}

func waitForPredicateKeyword(offending *TokenNode) (string, *TokenNode) {
	if offending == nil {
		return "", nil
	}

	if is(offending, "EXISTS") {
		if is(offending.Prev(), "NOT") {
			return "NOT EXISTS", offending
		}
		return "EXISTS", offending
	}

	if is(offending, "VALUE") {
		return "VALUE", offending
	}

	if is(offending.Prev(), "EXISTS") {
		if is(offending.Prev().Prev(), "NOT") {
			return "NOT EXISTS", offending.Prev()
		}
		return "EXISTS", offending.Prev()
	}

	if is(offending.Prev(), "VALUE") {
		return "VALUE", offending.Prev()
	}

	return "", nil
}

func waitForMissingClauseValue(offending *TokenNode) (string, *TokenNode) {
	if offending == nil {
		return "", nil
	}

	if is(offending, "TIMEOUT") || is(offending, "EVERY") || is(offending, "BACKOFF") || is(offending, "JITTER") {
		if hasWaitforBefore(offending) {
			return strings.ToUpper(offending.GetText()), offending
		}
	}

	prev := offending.Prev()
	if prev != nil {
		if is(prev, "TIMEOUT") || is(prev, "EVERY") || is(prev, "BACKOFF") || is(prev, "JITTER") {
			if hasWaitforBefore(prev) {
				return strings.ToUpper(prev.GetText()), prev
			}
		}
	}

	return "", nil
}

func hasWaitforBefore(node *TokenNode) bool {
	for curr := node.Prev(); curr != nil; curr = curr.Prev() {
		if is(curr, "WAITFOR") {
			return true
		}
	}

	return false
}
