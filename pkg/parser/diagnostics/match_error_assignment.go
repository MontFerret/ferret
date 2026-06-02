package diagnostics

import (
	"fmt"

	"github.com/antlr4-go/antlr/v4"

	"github.com/MontFerret/ferret/v2/pkg/diagnostics"
	"github.com/MontFerret/ferret/v2/pkg/source"
)

const (
	filterAssignmentMessage        = "Assignment is not valid in a FILTER predicate"
	filterComparisonAssignmentHint = "Use '==' to compare values, e.g. FILTER user.active == true."
	filterStatementAssignmentHint  = "Move the assignment to a standalone statement before FILTER. FILTER predicates must be expressions."
)

func matchMissingAssignmentValue(src *source.Source, err *diagnostics.Diagnostic, offending *TokenNode) bool {
	if matchInvalidVarDiscard(src, err, offending) {
		return true
	}

	if matchAssignmentExpression(src, err, offending) {
		return true
	}

	if isExtraneous(err.Message) {
		return false
	}

	prev := offending.Prev()

	if node := anyAssignmentOperator(offending, prev); node != nil {
		if keyword := node.Prev(); is(keyword, "COLLECT") || is(keyword, "AGGREGATE") {
			return false
		}

		span := spanFromTokenSafe(node.Token(), src)
		span.Start = span.End
		span.End = span.Start + 1

		prevText := ""
		if node.Prev() != nil {
			prevText = node.Prev().GetText()
		}

		err.Message = fmt.Sprintf("Expected expression after '%s' for variable '%s'", node.GetText(), prevText)
		err.Hint = "Did you forget to provide a value?"
		err.Spans = []diagnostics.ErrorSpan{
			diagnostics.NewMainErrorSpan(span, "missing value"),
		}

		return true
	}

	if is(offending, "LET") || is(offending, "VAR") {
		span := spanFromTokenSafe(offending.Token(), src)
		span.Start = span.End
		span.End = span.Start + 1
		err.Message = "Expected variable name"
		err.Hint = "Did you forget to provide a variable name?"
		err.Spans = []diagnostics.ErrorSpan{
			diagnostics.NewMainErrorSpan(span, "missing name"),
		}

		return true
	}

	return false
}

func matchInvalidVarDiscard(src *source.Source, err *diagnostics.Diagnostic, offending *TokenNode) bool {
	if offending == nil {
		return false
	}

	target := offending
	if is(offending, "VAR") && is(offending.Next(), "_") {
		target = offending.Next()
	} else if !is(offending, "_") || !is(offending.Prev(), "VAR") {
		return false
	}

	span := spanFromTokenSafe(target.Token(), src)

	err.Message = "VAR cannot use '_' as a variable name"
	err.Hint = "Use a real variable name for VAR, or use LET _ = ... to explicitly discard a value."
	err.Spans = []diagnostics.ErrorSpan{
		diagnostics.NewMainErrorSpan(span, "invalid variable name"),
	}

	return true
}

func matchAssignmentExpression(src *source.Source, err *diagnostics.Diagnostic, offending *TokenNode) bool {
	if offending == nil || (!isNoAlternative(err.Message) && !isMismatched(err.Message) && !isExtraneous(err.Message) && !isMissing(err.Message)) {
		return false
	}

	if matchFilterAssignmentExpression(src, err, offending) {
		return true
	}

	operator := findPrevAssignmentOperator(offending, 16)
	if operator == nil {
		return false
	}

	prev := operator.Prev()
	if prev == nil || !isIdentifier(prev) {
		return false
	}

	if operator != offending {
		if next := operator.Next(); next == nil || next == offending || isEOF(next) {
			return false
		}
	}

	if marker := operator.PrevAt(2); is(marker, "LET") || is(marker, "VAR") || is(marker, "COLLECT") || is(marker, "AGGREGATE") || is(marker, "FOR") {
		return false
	}

	span := spanFromTokenSafe(operator.Token(), src)

	err.Message = "Assignment is only allowed as a standalone statement"
	err.Hint = fmt.Sprintf("Move '%s %s ...' to its own statement. Assignment cannot be used inside expressions.", prev.GetText(), operator.GetText())
	err.Spans = []diagnostics.ErrorSpan{
		diagnostics.NewMainErrorSpan(span, "assignment is not an expression"),
	}

	return true
}

func matchFilterAssignmentExpression(src *source.Source, err *diagnostics.Diagnostic, offending *TokenNode) bool {
	if src == nil || err == nil || offending == nil || offending.Token() == nil {
		return false
	}

	tokens := lexDefaultTokens(src.Content())
	offendingIdx := findFilterAssignmentOffendingIndex(tokens, offending)
	if offendingIdx < 0 {
		return false
	}

	filterIdx := findFilterPredicateStart(tokens, offendingIdx)
	if filterIdx < 0 {
		return false
	}

	predicateStartIdx := filterIdx + 1
	predicateEndIdx := findFilterPredicateEnd(tokens, predicateStartIdx)
	if offendingIdx < predicateStartIdx || offendingIdx > predicateEndIdx {
		return false
	}

	operatorIdx := findFilterAssignmentOperator(tokens, predicateStartIdx, predicateEndIdx, offendingIdx)
	if operatorIdx < 0 {
		return false
	}

	operator := tokens[operatorIdx]
	span := spanFromTokenSafe(operator, src)
	label := "assignment is not an expression"
	hint := filterStatementAssignmentHint

	if isTokenText(operator, "=") {
		label = "use '==' for comparison"
		hint = filterComparisonAssignmentHint
	}

	err.Message = filterAssignmentMessage
	err.Hint = hint
	err.Spans = []diagnostics.ErrorSpan{
		diagnostics.NewMainErrorSpan(span, label),
	}

	return true
}

func findFilterAssignmentOffendingIndex(tokens []antlr.Token, offending *TokenNode) int {
	for node := offending; node != nil; node = node.Prev() {
		idx := findLexedTokenIndex(tokens, node.Token())
		if idx >= 0 {
			return idx
		}
	}

	return -1
}

func findFilterPredicateEnd(tokens []antlr.Token, startIdx int) int {
	depth := 0

	for i := startIdx; i < len(tokens); i++ {
		text := tokenText(tokens[i])

		switch text {
		case "(", "[", "{":
			depth++
		case ")", "]", "}":
			if depth == 0 {
				return i
			}

			depth--
		default:
			if depth == 0 && isFilterPredicateBoundary(text) {
				return i
			}
		}
	}

	return len(tokens)
}

func findFilterPredicateStart(tokens []antlr.Token, offendingIdx int) int {
	depth := 0

	for i := offendingIdx; i >= 0; i-- {
		text := tokenText(tokens[i])

		switch text {
		case ")", "]", "}":
			depth++
		case "(", "[", "{":
			if depth > 0 {
				depth--
			}
		default:
			if depth != 0 {
				continue
			}

			if text == "FILTER" {
				return i
			}

			if i != offendingIdx && isFilterPredicateBoundary(text) {
				return -1
			}
		}
	}

	return -1
}

func findFilterAssignmentOperator(tokens []antlr.Token, startIdx, endIdx, offendingIdx int) int {
	if startIdx < 0 || endIdx > len(tokens) || startIdx >= endIdx || offendingIdx < startIdx {
		return -1
	}

	if offendingIdx >= endIdx {
		offendingIdx = endIdx - 1
	}

	if operatorIdx := findForwardFilterAssignmentOperator(tokens, startIdx, endIdx, offendingIdx); operatorIdx >= 0 {
		return operatorIdx
	}

	return findBackwardFilterAssignmentOperator(tokens, startIdx, endIdx, offendingIdx)
}

func findForwardFilterAssignmentOperator(tokens []antlr.Token, startIdx, endIdx, offendingIdx int) int {
	depth := 0

	for i := offendingIdx; i < endIdx; i++ {
		text := tokenText(tokens[i])

		switch text {
		case "(", "[", "{":
			depth++
		case ")", "]", "}":
			if depth == 0 {
				return -1
			}

			depth--
		default:
			if depth == 0 && i != offendingIdx && isFilterAssignmentSearchBoundary(tokens, i) {
				return -1
			}
		}

		if !isAssignmentOperatorText(text) {
			continue
		}

		if isFilterAssignmentOperatorCandidate(tokens, startIdx, endIdx, i) {
			return i
		}

		return -1
	}

	return -1
}

func findBackwardFilterAssignmentOperator(tokens []antlr.Token, startIdx, endIdx, offendingIdx int) int {
	depth := 0

	for i := offendingIdx; i >= startIdx; i-- {
		text := tokenText(tokens[i])

		switch text {
		case ")", "]", "}":
			depth++
		case "(", "[", "{":
			if depth == 0 {
				return -1
			}
			depth--
		default:
			if depth == 0 && i != offendingIdx && isFilterAssignmentSearchBoundary(tokens, i) {
				return -1
			}
		}

		if !isAssignmentOperatorText(text) {
			continue
		}

		if isFilterAssignmentOperatorCandidate(tokens, startIdx, endIdx, i) {
			return i
		}

		return -1
	}

	return -1
}

func isFilterAssignmentOperatorCandidate(tokens []antlr.Token, startIdx, endIdx, operatorIdx int) bool {
	if isFilterDeclarationAssignmentOperator(tokens, operatorIdx) {
		return false
	}

	if !hasFilterAssignmentValue(tokens, operatorIdx+1, endIdx) {
		return false
	}

	targetStartIdx := findFilterAssignmentTargetStart(tokens, startIdx, operatorIdx)
	return targetStartIdx >= 0 && isFilterAssignmentTarget(tokens, targetStartIdx, operatorIdx)
}

func isFilterAssignmentSearchBoundary(tokens []antlr.Token, idx int) bool {
	text := tokenText(tokens[idx])
	if isFilterPredicateBoundary(text) {
		return true
	}

	switch text {
	case ",", ":", "AND", "OR", "&&", "||":
		return true
	case "?":
		return idx+1 >= len(tokens) || !isTokenText(tokens[idx+1], ".")
	default:
		return false
	}
}

func isFilterDeclarationAssignmentOperator(tokens []antlr.Token, operatorIdx int) bool {
	if operatorIdx < 2 {
		return false
	}

	switch tokenText(tokens[operatorIdx-2]) {
	case "LET", "VAR", "COLLECT", "AGGREGATE", "FOR":
		return true
	default:
		return false
	}
}

func findFilterAssignmentTargetStart(tokens []antlr.Token, startIdx, operatorIdx int) int {
	if startIdx < 0 || startIdx >= operatorIdx {
		return -1
	}

	depth := 0

	for i := operatorIdx - 1; i >= startIdx; i-- {
		text := tokenText(tokens[i])

		switch text {
		case ")", "]", "}":
			depth++
		case "(", "[", "{":
			if depth == 0 {
				return i + 1
			}

			depth--
		default:
			if depth == 0 && isFilterAssignmentTargetBoundary(tokens, i) {
				return i + 1
			}
		}
	}

	return startIdx
}

func isFilterAssignmentTargetBoundary(tokens []antlr.Token, idx int) bool {
	text := tokenText(tokens[idx])
	if isFilterPredicateBoundary(text) {
		return true
	}

	switch text {
	case ",", ":", "AND", "OR", "&&", "||":
		return true
	case "?":
		return idx+1 >= len(tokens) || !isTokenText(tokens[idx+1], ".")
	default:
		return false
	}
}

func hasFilterAssignmentValue(tokens []antlr.Token, startIdx, endIdx int) bool {
	if startIdx >= endIdx {
		return false
	}

	depth := 0
	seenValue := false

	for i := startIdx; i < endIdx; i++ {
		text := tokenText(tokens[i])

		switch text {
		case "(", "[", "{":
			depth++
			seenValue = true
		case ")", "]", "}":
			if depth == 0 {
				return seenValue
			}

			depth--
		case ",":
			if depth == 0 {
				return seenValue
			}
		default:
			if depth == 0 && isFilterPredicateBoundary(text) {
				return seenValue
			}

			seenValue = true
			if depth == 0 {
				return true
			}
		}
	}

	return seenValue && depth == 0
}

func isFilterAssignmentTarget(tokens []antlr.Token, startIdx, operatorIdx int) bool {
	if startIdx >= operatorIdx || startIdx < 0 || operatorIdx > len(tokens) {
		return false
	}

	if isTokenText(tokens[startIdx], ".") || isFilterAssignmentSafeImplicitTargetStart(tokens, startIdx) {
		return isFilterImplicitAssignmentTarget(tokens, startIdx, operatorIdx)
	}

	if !isFilterAssignmentIdentifier(tokens[startIdx]) {
		return false
	}

	return isFilterAssignmentPathTail(tokens, startIdx+1, operatorIdx)
}

func isFilterImplicitAssignmentTarget(tokens []antlr.Token, startIdx, operatorIdx int) bool {
	i := startIdx

	if isTokenText(tokens[i], "?") {
		i++
	}

	if i >= operatorIdx || !isTokenText(tokens[i], ".") {
		return false
	}

	i++
	if i >= operatorIdx {
		return false
	}

	if isTokenText(tokens[i], "[") {
		closeIdx := findFilterAssignmentComputedEnd(tokens, i, operatorIdx)
		if closeIdx < 0 {
			return false
		}

		i = closeIdx + 1
	} else {
		if !isFilterAssignmentPathToken(tokens[i]) {
			return false
		}

		i++
	}

	return isFilterAssignmentPathTail(tokens, i, operatorIdx)
}

func isFilterAssignmentSafeImplicitTargetStart(tokens []antlr.Token, startIdx int) bool {
	return isTokenText(tokens[startIdx], "?") && startIdx+1 < len(tokens) && isTokenText(tokens[startIdx+1], ".")
}

func isFilterAssignmentPathTail(tokens []antlr.Token, startIdx, operatorIdx int) bool {
	for i := startIdx; i < operatorIdx; {
		text := tokenText(tokens[i])

		switch text {
		case ".":
			i++
			if i >= operatorIdx || !isFilterAssignmentPathToken(tokens[i]) {
				return false
			}

			i++
		case "?":
			if i+1 >= operatorIdx || !isTokenText(tokens[i+1], ".") {
				return false
			}

			i += 2
			if i >= operatorIdx {
				return false
			}

			if isTokenText(tokens[i], "[") {
				closeIdx := findFilterAssignmentComputedEnd(tokens, i, operatorIdx)
				if closeIdx < 0 {
					return false
				}

				i = closeIdx + 1
				continue
			}

			if !isFilterAssignmentPathToken(tokens[i]) {
				return false
			}

			i++
		case "[":
			closeIdx := findFilterAssignmentComputedEnd(tokens, i, operatorIdx)
			if closeIdx < 0 {
				return false
			}

			i = closeIdx + 1
		default:
			return false
		}
	}

	return true
}

func findFilterAssignmentComputedEnd(tokens []antlr.Token, openIdx, limitIdx int) int {
	depth := 0

	for i := openIdx; i < limitIdx; i++ {
		text := tokenText(tokens[i])

		switch text {
		case "(", "[", "{":
			depth++
		case ")", "]", "}":
			depth--
			if depth < 0 {
				return -1
			}

			if depth == 0 {
				if delimitersMatch(tokenText(tokens[openIdx]), text) {
					return i
				}

				return -1
			}
		}
	}

	return -1
}

func isFilterAssignmentIdentifier(token antlr.Token) bool {
	if token == nil || tokenText(token) == "_" {
		return false
	}

	return isLoopVariableToken(&TokenNode{token: token})
}

func isFilterAssignmentPathToken(token antlr.Token) bool {
	if token == nil {
		return false
	}

	switch tokenText(token) {
	case "", "(", ")", "[", "]", "{", "}", ".", "?", ",", ":", "+", "-", "*", "/", "%", "=", "==", "!=", ">", "<", ">=", "<=", "=~", "!~", "=>", "<-", "&&", "||", "AND", "OR":
		return false
	default:
		return true
	}
}

func isFilterPredicateBoundary(text string) bool {
	switch text {
	case "RETURN", "FOR", "FILTER", "SORT", "LIMIT", "COLLECT", "LET", "VAR", "DELETE", "WAITFOR", "DISPATCH", "FUNC", "USE":
		return true
	default:
		return false
	}
}

func anyAssignmentOperator(nodes ...*TokenNode) *TokenNode {
	for _, node := range nodes {
		if isAssignmentOperator(node) {
			return node
		}
	}

	return nil
}

func findPrevAssignmentOperator(node *TokenNode, steps int) *TokenNode {
	current := node
	for i := 0; i < steps && current != nil; i++ {
		if isAssignmentOperator(current) {
			return current
		}

		current = current.Prev()
	}

	return nil
}

func isAssignmentOperator(node *TokenNode) bool {
	if node == nil {
		return false
	}

	return isAssignmentOperatorText(node.GetText())
}

func isAssignmentOperatorText(text string) bool {
	switch text {
	case "=", "+=", "-=", "*=", "/=":
		return true
	default:
		return false
	}
}
