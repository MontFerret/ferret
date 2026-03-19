package diagnostics

import (
	"regexp"

	"github.com/MontFerret/ferret/v2/pkg/diagnostics"
	"github.com/MontFerret/ferret/v2/pkg/file"
)

var legacyStepLoopPattern = regexp.MustCompile(`(?is)\bFOR\b[\s\S]*?\b[A-Za-z_][A-Za-z0-9_]*\s*=[\s\S]*?\bWHILE\b[\s\S]*?\b(STEP)\b`)

func matchStepLoopErrors(src *file.Source, err *diagnostics.Diagnostic, offending *TokenNode) bool {
	span, ok := findLegacyStepLoopSpan(src, offending)
	if !ok {
		return false
	}

	err.Message = "STEP is no longer supported in FOR loops"
	err.Hint = "Use VAR state with 'FOR _ WHILE ...' and update the counter inside the loop body."
	err.Spans = []diagnostics.ErrorSpan{
		diagnostics.NewMainErrorSpan(span, "removed syntax"),
	}

	return true
}

func findLegacyStepLoopSpan(src *file.Source, node *TokenNode) (file.Span, bool) {
	if step := findLegacyStepLoopToken(node); step != nil {
		return spanFromTokenSafe(step.Token(), src), true
	}

	if src == nil {
		return file.Span{}, false
	}

	indexes := legacyStepLoopPattern.FindStringSubmatchIndex(src.Content())
	if len(indexes) < 4 {
		return file.Span{}, false
	}

	return file.Span{
		Start: indexes[2],
		End:   indexes[3],
	}, true
}

func findLegacyStepLoopToken(node *TokenNode) *TokenNode {
	step := node
	if !is(step, "STEP") {
		step = findPrevToken(node, "STEP", 64)
		if step == nil {
			step = findNextToken(node, "STEP", 64)
		}
	}

	if step == nil {
		return nil
	}

	if !hasPrevToken(step, "WHILE", 64) {
		return nil
	}

	if !hasPrevToken(step, "FOR", 96) {
		return nil
	}

	return step
}

func hasPrevToken(node *TokenNode, tokenText string, steps int) bool {
	return findPrevToken(node, tokenText, steps) != nil
}

func findPrevToken(node *TokenNode, tokenText string, steps int) *TokenNode {
	current := node
	for i := 0; i < steps && current != nil; i++ {
		if is(current, tokenText) {
			return current
		}
		current = current.Prev()
	}

	return nil
}

func findNextToken(node *TokenNode, tokenText string, steps int) *TokenNode {
	current := node
	for i := 0; i < steps && current != nil; i++ {
		if is(current, tokenText) {
			return current
		}
		current = current.Next()
	}

	return nil
}
