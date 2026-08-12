package diagnostics

import (
	"fmt"
	"strings"

	"github.com/antlr4-go/antlr/v4"

	"github.com/MontFerret/ferret/v2/pkg/diagnostics"
	"github.com/MontFerret/ferret/v2/pkg/parser/fql"
	"github.com/MontFerret/ferret/v2/pkg/source"
)

func matchSpreadEntryErrors(src *source.Source, err *diagnostics.Diagnostic, offending *TokenNode) bool {
	if src == nil || err == nil || offending == nil {
		return false
	}

	if !isNoAlternative(err.Message) && !isMissing(err.Message) && !isMismatched(err.Message) && !isExtraneous(err.Message) {
		return false
	}

	ellipsis, target := findMissingSpread(src.Content(), err, offending)
	if target == "" {
		return false
	}

	lowerTarget := strings.ToLower(target)
	err.Message = fmt.Sprintf("Expected expression after '...' in %s spread", lowerTarget)
	err.Hint = fmt.Sprintf("Provide an %s expression or none after '...'.", target)
	err.Spans = []diagnostics.ErrorSpan{
		diagnostics.NewMainErrorSpan(spanFromTokenSafe(ellipsis, src), "missing expression"),
	}

	return true
}

func findMissingSpread(content string, err *diagnostics.Diagnostic, offending *TokenNode) (antlr.Token, string) {
	tokens := lexDefaultTokens(content)
	offendingIdx := findLexedTokenIndex(tokens, offending.Token())
	if offendingIdx < 0 {
		offendingIdx = findDiagnosticSpanTokenIndex(tokens, err)
	}

	for i, token := range tokens {
		if token.GetTokenType() != fql.FqlLexerEllipsis || i+1 >= len(tokens) {
			continue
		}

		if offendingIdx >= 0 && i != offendingIdx && i+1 != offendingIdx {
			continue
		}

		switch tokenText(tokens[i+1]) {
		case ",", "]", "}", "<EOF>":
			target := enclosingSpreadLiteral(tokens, i)
			if target != "" {
				return token, target
			}
		}
	}

	return nil, ""
}

func enclosingSpreadLiteral(tokens []antlr.Token, ellipsisIdx int) string {
	brackets := 0
	braces := 0
	parentheses := 0

	for i := ellipsisIdx - 1; i >= 0; i-- {
		switch tokenText(tokens[i]) {
		case "]":
			brackets++
		case "[":
			if brackets > 0 {
				brackets--

				continue
			}

			if braces == 0 && parentheses == 0 {
				return "Array"
			}
		case "}":
			braces++
		case "{":
			if braces > 0 {
				braces--

				continue
			}

			if brackets == 0 && parentheses == 0 {
				return "Object"
			}
		case ")":
			parentheses++
		case "(":
			if parentheses > 0 {
				parentheses--
			}
		}
	}

	return ""
}
