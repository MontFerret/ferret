package vm

import "github.com/MontFerret/ferret/pkg/diagnostics"

type RuntimeErrorSet = diagnostics.Diagnostics[*RuntimeError]

func NewRuntimeErrorSet(size int) *RuntimeErrorSet {
	return diagnostics.NewDiagnostics[*RuntimeError](size)
}
