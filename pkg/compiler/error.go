package compiler

import (
	"github.com/MontFerret/ferret/v2/pkg/compiler/internal/diagnostics"
)

type CompilationError = diagnostics.CompilationError
type CompilationErrorSet = diagnostics.CompilationErrorSet

var (
	SyntaxError = diagnostics.SyntaxError
	NameError   = diagnostics.NameError
)
