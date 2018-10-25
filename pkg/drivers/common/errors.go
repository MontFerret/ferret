package common

import "github.com/MontFerret/ferret/pkg/runtime/core"

var (
	ErrReadOnly    = core.Error(core.ErrInvalidOperation, "read only")
	ErrInvalidPath = core.Error(core.ErrInvalidOperation, "invalid path")
)
