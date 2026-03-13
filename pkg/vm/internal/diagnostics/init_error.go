package diagnostics

import (
	"fmt"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/diagnostics"
)

type (
	InitializationError struct {
		Cause   error
		Message string
		PC      int
		Dst     bytecode.Operand
	}

	InitializationErrorSet struct {
		diagnostics.Diagnostics[*InitializationError]
	}
)

func NewInitializationErrorSet(size int) *InitializationErrorSet {
	return &InitializationErrorSet{
		Diagnostics: *diagnostics.NewDiagnostics[*InitializationError](size),
	}
}

func (e *InitializationError) Error() string {
	if e == nil {
		return ""
	}

	if e.Cause != nil {
		return e.Message + ": " + e.Cause.Error()
	}

	return e.Message
}

func (e *InitializationError) Unwrap() error {
	return e.Cause
}

func (e *InitializationError) Format() string {
	return fmt.Sprintf("initialization error: %s", e.Message)
}

func (s *InitializationErrorSet) Add(err error, pc int, dst bytecode.Operand) {
	s.Diagnostics.Add(&InitializationError{
		Cause:   err,
		Message: fmt.Sprintf("initialization error at PC %d, dst %s", pc, dst),
		PC:      pc,
		Dst:     dst,
	})
}
