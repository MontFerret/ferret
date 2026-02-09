package diagnostics

import (
	"github.com/MontFerret/ferret/pkg/diagnostics"
)

type CompilationErrorSet struct {
	*diagnostics.Diagnostics[*CompilationError]
}

func NewCompilationErrorSet(size int) *CompilationErrorSet {
	return &CompilationErrorSet{
		Diagnostics: diagnostics.NewDiagnostics[*CompilationError](size),
	}
}
