package diagnostics

import (
	"fmt"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

type DestructureError struct {
	Actual   runtime.Type
	Expected string
}

func DestructureErrorOf(actual runtime.Value, expected string) error {
	return &DestructureError{
		Actual:   runtime.TypeOf(actual),
		Expected: expected,
	}
}

func (e *DestructureError) Error() string {
	if e == nil {
		return ""
	}

	return fmt.Sprintf("cannot destructure %s as %s", e.Actual, e.Expected)
}

func (e *DestructureError) Unwrap() error {
	return runtime.ErrInvalidType
}

func (e *DestructureError) Label() string {
	if e == nil {
		return ""
	}

	return fmt.Sprintf("value cannot be destructured as %s", e.Expected)
}

func (e *DestructureError) Note() string {
	return e.Error()
}

func (e *DestructureError) Hint() string {
	if e == nil {
		return ""
	}

	return fmt.Sprintf("Ensure the value supports %s destructuring", e.Expected)
}
