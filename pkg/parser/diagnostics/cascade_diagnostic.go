package diagnostics

import "github.com/MontFerret/ferret/v2/pkg/diagnostics"

func isCascadeStoppingSyntaxDiagnostic(err *diagnostics.Diagnostic) bool {
	return isFunctionBodySyntaxDiagnostic(err) || isArrayLiteralSeparatorDiagnostic(err)
}
