package diagnostics

import (
	"fmt"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/diagnostics"
)

type (
	WarmupError struct {
		Cause error
		PC    int
		Dst   bytecode.Operand
	}

	WarmupErrorSet struct {
		diagnostics.Diagnostics[*WarmupError]
	}
)

func NewWarmupErrorSet(size int) *WarmupErrorSet {
	return &WarmupErrorSet{
		Diagnostics: *diagnostics.NewDiagnostics[*WarmupError](size),
	}
}

func (e *WarmupError) Error() string {
	return e.Cause.Error()
}

func (e *WarmupError) Unwrap() error {
	return e.Cause
}

func (e *WarmupError) Format() string {
	return fmt.Sprintf("warmup error at PC %d, dst %s: %v", e.PC, e.Dst, e.Cause)
}

func (s *WarmupErrorSet) Add(err error, pc int, dst bytecode.Operand) {
	s.Diagnostics.Add(&WarmupError{
		Cause: err,
		PC:    pc,
		Dst:   dst,
	})
}
