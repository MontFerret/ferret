package vm

import "github.com/MontFerret/ferret/v2/pkg/vm/internal/diagnostic"

type RuntimeErrorSet = diagnostic.RuntimeErrorSet

func NewRuntimeErrorSet(size int) *RuntimeErrorSet {
	return diagnostic.NewRuntimeErrorSet(size)
}
