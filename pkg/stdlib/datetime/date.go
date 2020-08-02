package datetime

import (
	"context"
	"time"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/MontFerret/ferret/pkg/runtime/values/types"
)

// DATE converts RFC3339 date time string to DateTime object.
// @param {String} time - String in RFC3339 format.
// @return {DateTime} - New DateTime object derived from timeString.
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
