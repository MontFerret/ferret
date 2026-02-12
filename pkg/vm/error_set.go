package vm

import "github.com/MontFerret/ferret/v2/pkg/diagnostics"

type RuntimeErrorSet = diagnostics.Diagnostics[*RuntimeError]

func NewRuntimeErrorSet(size int) *RuntimeErrorSet {
	return diagnostics.NewDiagnostics[*RuntimeError](size)
}
