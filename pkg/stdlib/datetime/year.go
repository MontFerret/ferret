package datetime

import (
	"context"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
)

// DateYear returns the year extracted from the given date.
// @params date (DateTime) - source DateTime.
// @return (Int) - a year number.
func DateYear(_ context.Context, args ...core.Value) (core.Value, error) {
	err := core.ValidateArgs(args, 1, 1)
	if err != nil {
		return values.None, err
	}

	err = core.ValidateType(args[0], core.DateTimeType)
	if err != nil {
		return values.None, err
	}

	year := args[0].(values.DateTime).Year()

	return values.NewInt(int64(year)), nil
}
