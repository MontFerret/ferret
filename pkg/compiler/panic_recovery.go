package compiler

import (
	goruntime "runtime"

	"github.com/MontFerret/ferret/v2/pkg/diagnostics"
	parserd "github.com/MontFerret/ferret/v2/pkg/parser/diagnostics"
	"github.com/MontFerret/ferret/v2/pkg/source"
)

func addRecoveredAnalysisDiagnostic(src source.Source, errors *parserd.ErrorHandler, recovered any) {
	var diagnostic *diagnostics.Diagnostic

	buf := make([]byte, 1024)
	n := goruntime.Stack(buf, false)
	stackTrace := string(buf[:n])

	switch value := recovered.(type) {
	case string:
		diagnostic = diagnostics.NewUnexpectedError(src, value+"\n"+stackTrace)
	case error:
		diagnostic = diagnostics.NewUnexpectedErrorWith(src, "unhandled panic\n"+stackTrace, value)
	default:
		diagnostic = diagnostics.NewUnexpectedError(src, "unhandled panic\n"+stackTrace)
	}

	errors.Add(diagnostic)
}
