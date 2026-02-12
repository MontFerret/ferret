package types

import (
	"github.com/MontFerret/ferret/pkg/runtime"
)

// TO_DATETIME takes an input value of any type and converts it into the appropriate date time value.
// @param {Any} value - Input value of arbitrary type.
// @return {DateTime} - Parsed date time.
func ToDateTime(_ runtime.Context, arg runtime.Value) (runtime.Value, error) {
	return runtime.ParseDateTime(arg)
}
