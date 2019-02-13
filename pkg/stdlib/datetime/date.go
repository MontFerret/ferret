package datetime

import (
	"context"
	"time"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/MontFerret/ferret/pkg/runtime/values/types"
)

// Date convert RFC3339 date time string to DateTime object.
// @params timeString (String) - string in RFC3339 format.
// @return (DateTime) - new DateTime object derived from timeString.
func Date(_ context.Context, args ...core.Value) (core.Value, error) {
	err := core.ValidateArgs(args, 1, 1)
	if err != nil {
		return values.None, err
	}

	err = core.ValidateType(args[0], types.String)
	if err != nil {
		return values.None, err
	}

	timeStrings := args[0].(values.String)

	t, err := time.Parse(values.DefaultTimeLayout, timeStrings.String())
	if err != nil {
		return values.None, err
	}

	return values.NewDateTime(t), nil
}
