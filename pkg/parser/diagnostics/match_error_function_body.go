package diagnostics

import (
	"github.com/MontFerret/ferret/v2/pkg/diagnostics"
	"github.com/MontFerret/ferret/v2/pkg/source"
)

const (
	mixedFunctionBodyMessage = "Cannot combine arrow and block function body syntax"
	mixedFunctionBodyLabel   = "RETURN is only valid in a block function body"
	mixedFunctionBodyNote    = "Remove '=>' to use a block body, or remove RETURN and keep a single expression after '=>'."
	mixedFunctionBodyHint    = "Use either 'FUNC f(x) => expr' or 'FUNC f(x) ( ... RETURN expr )'."
)

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

func isFunctionNameToken(node *TokenNode) bool {
	if node == nil || node.Token() == nil {
		return false
	}

	return isIdentifier(node) || isKeyword(node)
}
