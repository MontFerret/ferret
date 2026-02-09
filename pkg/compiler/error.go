package compiler

import (
	"github.com/MontFerret/ferret/pkg/compiler/internal/diagnostics"
)

type CompilationError = diagnostics.CompilationError
type MultiCompilationError = diagnostics.CompilationErrorSet

var (
	SyntaxError = diagnostics.SyntaxError
	NameError   = diagnostics.NameError
)
