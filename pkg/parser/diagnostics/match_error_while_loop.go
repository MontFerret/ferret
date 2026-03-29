package diagnostics

import (
	"regexp"
	"sort"
	"strings"

	"github.com/MontFerret/ferret/v2/pkg/diagnostics"
	"github.com/MontFerret/ferret/v2/pkg/source"
)

var invalidWhileLoopBindingPatterns = []*regexp.Regexp{
	regexp.MustCompile(`(?is)\bFOR\b\s+([^\s]+)\s+WHILE\b`),
	regexp.MustCompile(`(?is)\bFOR\b\s+([^\s]+)\s+DO\s+WHILE\b`),
}

type whileLoopBindingMatch struct {
	binding     string
	bindingSpan source.Span
	headerStart int
	skipDo      bool
}

func matchWhileLoopErrors(src *source.Source, err *diagnostics.Diagnostic, offending *TokenNode) bool {
	if matchInvalidWhileLoopBinding(src, err, offending) {
		return true
	}

	if matchMissingWhileLoopCondition(src, err, offending) {
		return true
	}

	if matchStandaloneWhileLoop(src, err, offending) {
		return true
	}

	return false
}

func matchInvalidWhileLoopBinding(src *source.Source, err *diagnostics.Diagnostic, offending *TokenNode) bool {
	span, ok := findInvalidWhileLoopBindingSpan(src)
	if !ok {
		return false
	}

	err.Message = "Expected identifier before 'WHILE'"
	err.Hint = "Use 'FOR WHILE [condition]' or 'FOR x WHILE [condition]' syntax."
	err.Spans = []diagnostics.ErrorSpan{
		diagnostics.NewMainErrorSpan(span, "invalid binding"),
	}

	return true
}

func findInvalidWhileLoopBindingSpan(src *source.Source) (source.Span, bool) {
	if src == nil {
		return source.Span{}, false
	}

	content := src.Content()
	matches := make([]whileLoopBindingMatch, 0, len(invalidWhileLoopBindingPatterns))

	for i, pattern := range invalidWhileLoopBindingPatterns {
		for _, indexes := range pattern.FindAllStringSubmatchIndex(content, -1) {
			if len(indexes) < 4 {
				continue
			}

			matches = append(matches, whileLoopBindingMatch{
				headerStart: indexes[0],
				bindingSpan: source.Span{
					Start: indexes[2],
					End:   indexes[3],
				},
				binding: content[indexes[2]:indexes[3]],
				skipDo:  i == 0,
			})
		}
	}

	sort.Slice(matches, func(i, j int) bool {
		if matches[i].headerStart != matches[j].headerStart {
			return matches[i].headerStart < matches[j].headerStart
		}

		return matches[i].bindingSpan.Start < matches[j].bindingSpan.Start
	})

	for _, match := range matches {
		if match.skipDo && strings.EqualFold(match.binding, "DO") {
			continue
		}

		if isValidWhileLoopBindingText(match.binding) {
			continue
		}

		return match.bindingSpan, true
	}

	return source.Span{}, false
}

func isValidWhileLoopBindingText(text string) bool {
	if text == "_" {
		return true
	}

	if text == "" {
		return false
	}

	for i, ch := range text {
		if i == 0 {
			if ch != '_' && (ch < 'A' || ch > 'Z') && (ch < 'a' || ch > 'z') {
				return false
			}
			continue
		}

		if ch != '_' && (ch < '0' || ch > '9') && (ch < 'A' || ch > 'Z') && (ch < 'a' || ch > 'z') {
			return false
		}
	}

	switch strings.ToUpper(text) {
	case "RETURN", "DISPATCH", "QUERY", "USING", "NONE", "NULL", "LET", "VAR", "USE", "WAITFOR",
		"WHILE", "DO", "IN", "LIKE", "NOT", "FOR", "TRUE", "FALSE", "THROW", "MATCH", "WHEN", "FUNC":
		return false
	default:
		return true
	}
}

func matchMissingWhileLoopCondition(src *source.Source, err *diagnostics.Diagnostic, offending *TokenNode) bool {
	whileToken := findWhileLoopHeaderToken(offending)
	if whileToken == nil {
		return false
	}

	span := spanFromTokenSafe(whileToken.Token(), src)
	err.Message = "Expected condition after 'WHILE'"
	err.Hint = "Use 'FOR WHILE [condition]' or 'FOR x WHILE [condition]' syntax."
	err.Spans = []diagnostics.ErrorSpan{
		diagnostics.NewMainErrorSpan(span, "missing condition"),
	}

	return true
}

func matchStandaloneWhileLoop(src *source.Source, err *diagnostics.Diagnostic, offending *TokenNode) bool {
	loopKind, span, ok := findStandaloneWhileLoopSpan(src, err, offending)
	if !ok {
		return false
	}

	switch loopKind {
	case "DO WHILE":
		err.Message = "Standalone DO WHILE loops are not supported"
		err.Hint = "Use 'FOR DO WHILE [condition]' or 'FOR x DO WHILE [condition]' syntax."
	default:
		err.Message = "Standalone WHILE loops are not supported"
		err.Hint = "Use 'FOR WHILE [condition]' or 'FOR x WHILE [condition]' syntax."
	}

	err.Spans = []diagnostics.ErrorSpan{
		diagnostics.NewMainErrorSpan(span, "unsupported loop"),
	}

	return true
}

func findStandaloneWhileLoopSpan(src *source.Source, err *diagnostics.Diagnostic, offending *TokenNode) (string, source.Span, bool) {
	whileToken := findStandaloneWhileLoopToken(offending)
	if whileToken == nil {
		if is(offending, "DO") && err != nil && isNoAlternative(err.Message) && has(err.Message, "do while") && !hasPrevToken(offending, "FOR", 4) {
			return "DO WHILE", spanFromTokenSafe(offending.Token(), src), true
		}

		return "", source.Span{}, false
	}

	if is(whileToken.Prev(), "DO") {
		doSpan := spanFromTokenSafe(whileToken.Prev().Token(), src)
		whileSpan := spanFromTokenSafe(whileToken.Token(), src)

		return "DO WHILE", source.Span{
			Start: doSpan.Start,
			End:   whileSpan.End,
		}, true
	}

	return "WHILE", spanFromTokenSafe(whileToken.Token(), src), true
}

func findStandaloneWhileLoopToken(offending *TokenNode) *TokenNode {
	for _, node := range whileLoopTokenCandidates(offending) {
		if !is(node, "WHILE") {
			continue
		}

		if isWhileLoopHeaderToken(node) || hasPrevToken(node, "FOR", 4) {
			continue
		}

		return node
	}

	return nil
}

func findWhileLoopHeaderToken(offending *TokenNode) *TokenNode {
	for _, node := range whileLoopTokenCandidates(offending) {
		if isWhileLoopHeaderToken(node) {
			return node
		}
	}

	return nil
}

func isWhileLoopHeaderToken(node *TokenNode) bool {
	if !is(node, "WHILE") {
		return false
	}

	prev := node.Prev()
	switch {
	case is(prev, "FOR"):
		return true
	case is(prev, "DO") && is(prev.Prev(), "FOR"):
		return true
	case isLoopVariableToken(prev) && is(prev.Prev(), "FOR"):
		return true
	case is(prev, "DO") && isLoopVariableToken(prev.Prev()) && is(prev.PrevAt(2), "FOR"):
		return true
	default:
		return false
	}
}

func whileLoopTokenCandidates(offending *TokenNode) []*TokenNode {
	if offending == nil {
		return nil
	}

	candidates := []*TokenNode{offending}

	for i := 1; i <= 3; i++ {
		candidates = append(candidates, offending.PrevAt(i), offending.NextAt(i))
	}

	return candidates
}
