package diagnostics

import (
	"github.com/MontFerret/ferret/pkg/diagnostics"
)

type CompilationErrorSet struct {
	*diagnostics.Diagnostics[*CompilationError]
}

func NewCompilationErrorSet(errors []*CompilationError) error {
	return &CompilationErrorSet{
		Diagnostics: diagnostics.NewDiagnostics[*CompilationError](errors),
	}
}
