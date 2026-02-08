package compiler

import (
	"github.com/MontFerret/ferret/pkg/compiler/internal/diagnostics"
)

type ErrorKind = diagnostics.ErrorKind
type CompilationError = diagnostics.CompilationError
type MultiCompilationError = diagnostics.MultiCompilationError

var (
	UnknownError     = diagnostics.UnknownError
	SyntaxError      = diagnostics.SyntaxError
	NameError        = diagnostics.NameError
	TypeError        = diagnostics.TypeError
	SemanticError    = diagnostics.SemanticError
	UnsupportedError = diagnostics.UnsupportedError
	InternalError    = diagnostics.InternalError
)
