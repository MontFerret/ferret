package diagnostics

import (
	"fmt"
	"strings"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

type SpreadError struct {
	Actual   runtime.Type
	Target   runtime.Type
	Expected runtime.Type
}

func SpreadErrorOf(actual runtime.Value, target, expected runtime.Type) error {
	return &SpreadError{
		Actual:   runtime.TypeOf(actual),
		Target:   target,
		Expected: expected,
	}
}

func (e *SpreadError) Error() string {
	if e == nil {
		return ""
	}

	return fmt.Sprintf("cannot spread %s into %s", e.Actual, e.Target)
}

func (e *SpreadError) Unwrap() error {
	return runtime.ErrInvalidType
}

func (e *SpreadError) Label() string {
	if e == nil {
		return ""
	}

	return fmt.Sprintf("expected %s or none in %s literal", e.Expected, strings.ToLower(e.Target.String()))
}

func (e *SpreadError) Note() string {
	return e.Error()
}

func (e *SpreadError) Hint() string {
	if e == nil {
		return ""
	}

	article := "a"
	if e.Expected == runtime.TypeObject {
		article = "an"
	}

	return fmt.Sprintf(
		"Spread %s %s value or none inside an %s literal",
		article,
		e.Expected,
		strings.ToLower(e.Target.String()),
	)
}
