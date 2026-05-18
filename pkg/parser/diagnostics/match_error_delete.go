package diagnostics

import (
	"github.com/MontFerret/ferret/v2/pkg/diagnostics"
	"github.com/MontFerret/ferret/v2/pkg/source"
)

const deleteTargetMessage = "DELETE requires a property or computed-key target"

func matchDeleteStatementErrors(src *source.Source, err *diagnostics.Diagnostic, offending *TokenNode) bool {
	if offending == nil || isEOF(offending) {
		return false
	}

	if !isNoAlternative(err.Message) && !isMismatched(err.Message) && !isExtraneous(err.Message) && !isMissing(err.Message) {
		return false
	}

	deleteToken := findPrevDeleteToken(offending, 16)
	if deleteToken == nil {
		return false
	}

	span := spanFromTokenSafe(offending.Token(), src)
	if isMissing(err.Message) {
		span = spanFromTokenSafe(deleteToken.Token(), src)
		span.Start = span.End
		span.End = span.Start + 1
	}

	err.Message = deleteTargetMessage
	err.Hint = `Use DELETE obj.foo or DELETE obj["foo"] to remove a property.`
	err.Spans = []diagnostics.ErrorSpan{
		diagnostics.NewMainErrorSpan(span, "invalid delete target"),
	}

	return true
}

func findPrevDeleteToken(node *TokenNode, steps int) *TokenNode {
	current := node
	for i := 0; i < steps && current != nil; i++ {
		if is(current, "DELETE") {
			return current
		}

		current = current.Prev()
	}

	return nil
}
