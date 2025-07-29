package compiler

import "github.com/MontFerret/ferret/pkg/compiler/internal/core"

type ErrorKind = core.ErrorKind
type CompilationError = core.CompilationError
type MultiCompilationError = core.MultiCompilationError

var (
	UnknownError     = core.UnknownError
	SyntaxError      = core.SyntaxError
	NameError        = core.NameError
	TypeError        = core.TypeError
	SemanticError    = core.SemanticError
	UnsupportedError = core.UnsupportedError
	InternalError    = core.InternalError
)
