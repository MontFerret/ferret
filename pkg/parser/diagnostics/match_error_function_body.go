package diagnostics

import (
	"github.com/antlr4-go/antlr/v4"

	"github.com/MontFerret/ferret/v2/pkg/diagnostics"
	"github.com/MontFerret/ferret/v2/pkg/source"
)

const (
	missingFunctionParamsCloseMessage = "Expected function parameters before function body"
	missingFunctionParamsCloseLabel   = "missing parameter list before function body"
	missingFunctionParamsCloseHint    = "Add a parameter list before the block body, e.g. FUNC fib(n) ( ... RETURN expr ). Use FUNC fib() ( ... ) for no parameters."

	mixedFunctionBodyMessage = "Cannot combine arrow and block function body syntax"
	mixedFunctionBodyLabel   = "RETURN is only valid in a block function body"
	mixedFunctionBodyNote    = "Remove '=>' to use a block body, or remove RETURN and keep a single expression after '=>'."
	mixedFunctionBodyHint    = "Use either 'FUNC f(x) => expr' or 'FUNC f(x) ( ... RETURN expr )'."
)

func matchMissingFunctionParamsClose(src *source.Source, err *diagnostics.Diagnostic, offending *TokenNode) bool {
	if src == nil || err == nil || offending == nil {
		return false
	}

	if !isNoAlternative(err.Message) && !isMissing(err.Message) && !isMismatched(err.Message) {
		return false
	}

	tokens := lexDefaultTokens(src.Content())
	offendingIdx := findDiagnosticSpanTokenIndex(tokens, err)
	if offendingIdx < 0 || !isFunctionBlockBodyStartToken(tokens[offendingIdx]) {
		return false
	}

	paramsOpenIdx := findUnclosedFunctionParamsOpen(tokens, offendingIdx)
	if paramsOpenIdx < 0 {
		return false
	}

	err.Message = missingFunctionParamsCloseMessage
	err.Hint = missingFunctionParamsCloseHint
	err.Spans = []diagnostics.ErrorSpan{
		diagnostics.NewMainErrorSpan(spanFromTokenSafe(tokens[paramsOpenIdx], src), missingFunctionParamsCloseLabel),
	}

	return true
}

func matchMixedFunctionBodySyntax(src *source.Source, err *diagnostics.Diagnostic, offending *TokenNode) bool {
	if src == nil || err == nil || offending == nil {
		return false
	}

	if !isNoAlternative(err.Message) && !isMissing(err.Message) && !isMismatched(err.Message) {
		return false
	}

	if !is(offending, "RETURN") {
		return false
	}

	openParen := offending.Prev()
	if !is(openParen, "(") || !is(openParen.Prev(), "=>") {
		return false
	}

	if !isFunctionDeclarationArrow(openParen.Prev()) {
		return false
	}

	err.Message = mixedFunctionBodyMessage
	err.Hint = mixedFunctionBodyHint
	err.Note = mixedFunctionBodyNote
	err.Spans = []diagnostics.ErrorSpan{
		diagnostics.NewMainErrorSpan(spanFromTokenSafe(offending.Token(), src), mixedFunctionBodyLabel),
	}

	return true
}

func isFunctionBodySyntaxDiagnostic(err *diagnostics.Diagnostic) bool {
	return isMissingFunctionParamsCloseDiagnostic(err) || isMixedFunctionBodySyntaxDiagnostic(err)
}

func isMissingFunctionParamsCloseDiagnostic(err *diagnostics.Diagnostic) bool {
	return err != nil && err.Kind == SyntaxError && err.Message == missingFunctionParamsCloseMessage
}

func isMixedFunctionBodySyntaxDiagnostic(err *diagnostics.Diagnostic) bool {
	return err != nil && err.Kind == SyntaxError && err.Message == mixedFunctionBodyMessage
}

func isFunctionDeclarationArrow(arrow *TokenNode) bool {
	if !is(arrow, "=>") {
		return false
	}

	paramsClose := arrow.Prev()
	if !is(paramsClose, ")") {
		return false
	}

	depth := 0
	for current := paramsClose; current != nil; current = current.Prev() {
		switch {
		case is(current, ")"):
			depth++
		case is(current, "("):
			depth--
			if depth == 0 {
				name := current.Prev()
				return isFunctionNameToken(name) && is(name.Prev(), "FUNC")
			}
		}
	}

	return false
}

func findDiagnosticSpanTokenIndex(tokens []antlr.Token, err *diagnostics.Diagnostic) int {
	if err == nil || len(err.Spans) == 0 {
		return -1
	}

	span := err.Spans[0].Span
	for idx, token := range tokens {
		if token.GetStart() == span.Start && token.GetStop()+1 == span.End {
			return idx
		}
	}

	return -1
}

func isFunctionBlockBodyStartToken(token antlr.Token) bool {
	return isTokenText(token, "RETURN") ||
		isTokenText(token, "LET") ||
		isTokenText(token, "VAR") ||
		isTokenText(token, "FUNC") ||
		isTokenText(token, "WAITFOR") ||
		isTokenText(token, "DISPATCH")
}

func findUnclosedFunctionParamsOpen(tokens []antlr.Token, offendingIdx int) int {
	depth := 0

	for i := offendingIdx - 1; i >= 0; i-- {
		switch tokenText(tokens[i]) {
		case ")":
			depth++
		case "(":
			if depth > 0 {
				depth--
				continue
			}

			if i >= 2 && isTokenText(tokens[i-2], "FUNC") {
				return i
			}

			return -1
		}
	}

	return -1
}

func isFunctionNameToken(node *TokenNode) bool {
	if node == nil || node.Token() == nil {
		return false
	}

	return isIdentifier(node) || isKeyword(node)
}
