package diagnostics

import (
	"regexp"
	"strings"

	"github.com/MontFerret/ferret/v2/pkg/diagnostics"
	"github.com/MontFerret/ferret/v2/pkg/file"
)

var invalidWhileLoopBindingPatterns = []*regexp.Regexp{
	regexp.MustCompile(`(?is)\bFOR\b\s+([^\s]+)\s+WHILE\b`),
	regexp.MustCompile(`(?is)\bFOR\b\s+([^\s]+)\s+DO\s+WHILE\b`),
}

func matchWhileLoopErrors(src *file.Source, err *diagnostics.Diagnostic, offending *TokenNode) bool {
	if matchInvalidWhileLoopBinding(src, err, offending) {
		return true
	}

	if matchMissingWhileLoopCondition(src, err, offending) {
		return true
	}

	return false
}

func matchInvalidWhileLoopBinding(src *file.Source, err *diagnostics.Diagnostic, offending *TokenNode) bool {
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

func findInvalidWhileLoopBindingSpan(src *file.Source) (file.Span, bool) {
	if src == nil {
		return file.Span{}, false
	}

	content := src.Content()

	for i, pattern := range invalidWhileLoopBindingPatterns {
		indexes := pattern.FindStringSubmatchIndex(content)
		if len(indexes) < 4 {
			continue
		}

		binding := content[indexes[2]:indexes[3]]
		if i == 0 && strings.EqualFold(binding, "DO") {
			continue
		}

		if isValidWhileLoopBindingText(binding) {
			continue
		}

		return file.Span{
			Start: indexes[2],
			End:   indexes[3],
		}, true
	}

	return file.Span{}, false
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

func matchMissingWhileLoopCondition(src *file.Source, err *diagnostics.Diagnostic, offending *TokenNode) bool {
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

func findWhileLoopHeaderToken(offending *TokenNode) *TokenNode {
	candidates := []*TokenNode{offending}
	if offending != nil {
		candidates = append(candidates, offending.Prev())
	}

	for _, node := range candidates {
		if !is(node, "WHILE") {
			continue
		}

		prev := node.Prev()
		switch {
		case is(prev, "FOR"):
			return node
		case is(prev, "DO") && is(prev.Prev(), "FOR"):
			return node
		case isLoopVariableToken(prev) && is(prev.Prev(), "FOR"):
			return node
		case is(prev, "DO") && isLoopVariableToken(prev.Prev()) && is(prev.PrevAt(2), "FOR"):
			return node
		}
	}

	return nil
}
