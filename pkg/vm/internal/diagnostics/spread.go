package diagnostics

import (
	"fmt"
	"strings"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

type SpreadError struct {
	Actual runtime.Type
	Target string
}

func SpreadErrorOf(actual runtime.Value, target string) error {
	return &SpreadError{
		Actual: runtime.TypeOf(actual),
		Target: target,
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

	return fmt.Sprintf("expected %s or none in %s literal", e.Target, strings.ToLower(e.Target))
}

func (e *SpreadError) Note() string {
	return e.Error()
}

func (e *SpreadError) Hint() string {
	if e == nil {
		return ""
	}

	return fmt.Sprintf("Spread an %s value or none inside an %s literal", e.Target, strings.ToLower(e.Target))
}
