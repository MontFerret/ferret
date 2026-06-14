package debugger

import (
	"fmt"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

// StateError reports a debugger command that is invalid for the current
// session state.
type StateError struct {
	Operation string
	State     string
}

func (e *StateError) Error() string {
	return fmt.Sprintf("cannot %s while debug session is %s", e.Operation, e.State)
}

func (e *StateError) Unwrap() error {
	return runtime.ErrInvalidOperation
}

func (e *StateError) Is(target error) bool {
	_, ok := target.(*StateError)
	return ok
}
