package diagnostics

import (
	"fmt"

	"github.com/MontFerret/ferret/v2/pkg/diagnostics"
	"github.com/MontFerret/ferret/v2/pkg/source"
)

func matchDispatchErrors(src *source.Source, err *diagnostics.Diagnostic, offending *TokenNode) bool {
	if err == nil || offending == nil {
		return false
	}

	if !isMismatched(err.Message) && !isMissing(err.Message) && !isNoAlternative(err.Message) && !isExtraneous(err.Message) {
		return false
	}

	if spanNode := shorthandMissingEventNode(offending); spanNode != nil {
		span := spanFromTokenSafe(spanNode.Token(), src)
		err.Message = "Expected dispatch event before '->'"
		err.Hint = `Provide an event expression, e.g. "click" -> btn.`
		err.Spans = []diagnostics.ErrorSpan{
			diagnostics.NewMainErrorSpan(span, "missing dispatch event"),
		}
		return true
	}

	if spanNode := longFormMissingEventNode(offending); spanNode != nil {
		span := spanFromTokenSafe(spanNode.Token(), src)
		err.Message = "Expected dispatch event after 'DISPATCH'"
		err.Hint = `Provide an event expression, e.g. DISPATCH "click" IN btn.`
		err.Spans = []diagnostics.ErrorSpan{
			diagnostics.NewMainErrorSpan(span, "missing dispatch event"),
		}
		return true
	}

	if spanNode := shorthandMissingTargetNode(offending); spanNode != nil {
		span := spanFromTokenSafe(spanNode.Token(), src)
		err.Message = "Expected dispatch target after '->'"
		err.Hint = `Provide a dispatchable target, e.g. "click" -> btn.`
		err.Spans = []diagnostics.ErrorSpan{
			diagnostics.NewMainErrorSpan(span, "missing dispatch target"),
		}
		return true
	}

	if spanNode := longFormMissingTargetNode(offending); spanNode != nil {
		span := spanFromTokenSafe(spanNode.Token(), src)
		err.Message = "Expected dispatch target after 'IN'"
		err.Hint = `Provide a dispatchable target, e.g. DISPATCH "click" IN btn.`
		err.Spans = []diagnostics.ErrorSpan{
			diagnostics.NewMainErrorSpan(span, "missing dispatch target"),
		}
		return true
	}

	if clause, spanNode := shorthandUnsupportedClause(offending); clause != "" {
		span := spanFromTokenSafe(spanNode.Token(), src)
		err.Message = fmt.Sprintf("Dispatch shorthand does not support %s", clause)
		err.Hint = shorthandClauseHint(clause)
		err.Spans = []diagnostics.ErrorSpan{
			diagnostics.NewMainErrorSpan(span, "unsupported shorthand clause"),
		}
		return true
	}

	if clause, spanNode := longFormMissingClauseValue(offending); clause != "" {
		span := spanFromTokenSafe(spanNode.Token(), src)
		err.Message = fmt.Sprintf("Expected dispatch %s after %s", dispatchClauseValueName(clause), clause)
		err.Hint = longFormClauseHint(clause)
		err.Spans = []diagnostics.ErrorSpan{
			diagnostics.NewMainErrorSpan(span, "missing dispatch clause value"),
		}
		return true
	}

	return false
}

func shorthandMissingEventNode(offending *TokenNode) *TokenNode {
	if offending == nil {
		return nil
	}

	arrow := offending
	if !is(arrow, "->") {
		if next := findNextToken(offending, "->", 4); next != nil {
			arrow = next
		} else {
			_, arrow = prevTokenDistance(offending, "->", 12)
			if arrow == nil {
				return nil
			}
		}
	}

	prev := arrow.Prev()
	if prev == nil || is(prev, "=") || is(prev, "RETURN") || is(prev, "(") || is(prev, ",") || is(prev, "=>") || is(prev, "?") || is(prev, ":") {
		return arrow
	}

	return nil
}

func longFormMissingEventNode(offending *TokenNode) *TokenNode {
	if offending == nil {
		return nil
	}

	if is(offending, "IN") && is(offending.Prev(), "DISPATCH") {
		return offending.Prev()
	}

	return nil
}

func shorthandMissingTargetNode(offending *TokenNode) *TokenNode {
	if offending == nil {
		return nil
	}

	arrow := offending
	if !is(arrow, "->") {
		if next := findNextToken(offending, "->", 4); next != nil {
			arrow = next
		} else {
			_, arrow = prevTokenDistance(offending, "->", 12)
			if arrow == nil {
				return nil
			}
		}
	}

	next := arrow.Next()
	// next == offending: the parser reported the token immediately after '->' as
	// the offending symbol, meaning the arrow's successor is itself the problem token.
	if next == nil || next == offending || isEOF(next) ||
		is(next, "RETURN") || is(next, "WITH") || is(next, "OPTIONS") ||
		is(next, ",") || is(next, ")") || is(next, "]") {
		return arrow
	}

	if is(offending, "<EOF>") && is(offending.Prev(), "->") {
		return offending.Prev()
	}

	prev := offending.Prev()
	if prev == nil || !is(prev, "->") {
		return nil
	}

	if is(offending, "RETURN") || is(offending, "WITH") || is(offending, "OPTIONS") || is(offending, ",") || is(offending, ")") || is(offending, "]") {
		return prev
	}

	return nil
}

func longFormMissingTargetNode(offending *TokenNode) *TokenNode {
	if offending == nil {
		return nil
	}

	if is(offending, "<EOF>") && is(offending.Prev(), "IN") && hasPrevToken(offending.Prev(), "DISPATCH", 12) {
		return offending.Prev()
	}

	prev := offending.Prev()
	if prev == nil || !is(prev, "IN") || !hasPrevToken(prev, "DISPATCH", 12) {
		return nil
	}

	if is(offending, "RETURN") || is(offending, "WITH") || is(offending, "OPTIONS") || is(offending, ",") || is(offending, ")") || is(offending, "]") {
		return prev
	}

	return nil
}

func shorthandUnsupportedClause(offending *TokenNode) (string, *TokenNode) {
	if offending == nil {
		return "", nil
	}

	for _, clause := range []string{"WITH", "OPTIONS"} {
		if distance, node := prevTokenDistance(offending, clause, 12); node != nil {
			// distance == 0: offending token IS the clause keyword itself.
			// The second condition guards against a long-form clause whose preceding
			// token is something other than '->' (which would mean it belongs to a
			// DISPATCH...IN form, not a shorthand). We only flag it as an unsupported
			// shorthand clause when '->' appears somewhere before it without DISPATCH.
			clauseIsOffendingOrNotAfterArrow := distance == 0 || (node.Prev() != nil && !is(node.Prev(), "->"))
			if clauseIsOffendingOrNotAfterArrow && hasPrevToken(node, "->", 8) && !hasPrevToken(node, "DISPATCH", 12) {
				return node.GetText(), node
			}
		}
	}

	return "", nil
}

func longFormMissingClauseValue(offending *TokenNode) (string, *TokenNode) {
	if offending == nil {
		return "", nil
	}

	prev := offending.Prev()
	if prev != nil && (is(prev, "WITH") || is(prev, "OPTIONS")) && hasPrevToken(prev, "DISPATCH", 12) {
		return prev.GetText(), prev
	}

	if (is(offending, "WITH") || is(offending, "OPTIONS")) && hasPrevToken(offending, "DISPATCH", 12) {
		next := offending.Next()
		if next == nil || is(next, "<EOF>") {
			return offending.GetText(), offending
		}
	}

	return "", nil
}

func shorthandClauseHint(clause string) string {
	switch clause {
	case "WITH":
		return `Use the long form instead, e.g. DISPATCH "input" IN field WITH { value: "x" }.`
	case "OPTIONS":
		return `Use the long form instead, e.g. DISPATCH "click" IN btn OPTIONS { bubbles: true }.`
	default:
		return "Use the long form DISPATCH syntax for configured dispatch."
	}
}

func longFormClauseHint(clause string) string {
	switch clause {
	case "WITH":
		return `Provide a payload expression, e.g. DISPATCH "input" IN field WITH { value: "x" }.`
	case "OPTIONS":
		return `Provide an options expression, e.g. DISPATCH "click" IN btn OPTIONS { bubbles: true }.`
	default:
		return "Provide a dispatch clause value."
	}
}

func dispatchClauseValueName(clause string) string {
	switch clause {
	case "WITH":
		return "payload"
	case "OPTIONS":
		return "options"
	default:
		return "value"
	}
}
