package html

import (
	"context"

	"github.com/MontFerret/ferret/pkg/drivers"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
)

// SELECT selects a value from an underlying select element.
// @param source (Open | GetElement) - Event target.
// @param valueOrSelector (String | Array<String>) - Selector or a an array of strings as a value.
// @param value (Array<String) - Target value. Optional.
// @returns (Array<String>) - Returns an array of selected values.
func Select(ctx context.Context, args ...core.Value) (core.Value, error) {
	err := core.ValidateArgs(args, 2, 4)

	if err != nil {
		return values.None, err
	}

	el, err := drivers.ToElement(args[0])

	if err != nil {
		return values.None, err
	}

	if len(args) == 2 {
		arr := values.ToArray(ctx, args[1])

		return el.Select(ctx, arr)
	}

	selector := values.ToString(args[1])
	arr := values.ToArray(ctx, args[2])

	return el.SelectBySelector(ctx, selector, arr)
}
