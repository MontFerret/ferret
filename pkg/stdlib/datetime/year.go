package datetime

import (
	"context"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
<<<<<<< HEAD
	"github.com/MontFerret/ferret/pkg/runtime/values/types"
=======
>>>>>>> db2ab4b... added DATE_YEAR function
)

// DateYear returns the year extracted from the given date.
// @params date (DateTime) - source DateTime.
// @return (Int) - a year number.
func DateYear(_ context.Context, args ...core.Value) (core.Value, error) {
	err := core.ValidateArgs(args, 1, 1)
	if err != nil {
		return values.None, err
	}

<<<<<<< HEAD
	err = core.ValidateType(args[0], types.DateTime)
=======
	err = core.ValidateType(args[0], core.DateTimeType)
>>>>>>> db2ab4b... added DATE_YEAR function
	if err != nil {
		return values.None, err
	}

	year := args[0].(values.DateTime).Year()

	return values.NewInt(year), nil
}
