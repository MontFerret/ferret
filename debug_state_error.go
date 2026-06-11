package ferret

import (
	"fmt"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

// DebugStateError reports a debugger command that is invalid for the current
// session state.
type DebugStateError struct {
	Operation string
	State     string
}

func (e *DebugStateError) Error() string {
	return fmt.Sprintf("cannot %s while debug session is %s", e.Operation, e.State)
}

func (e *DebugStateError) Unwrap() error {
	return runtime.ErrInvalidOperation
}

func (e *DebugStateError) Is(target error) bool {
	_, ok := target.(*DebugStateError)
	return ok
}
