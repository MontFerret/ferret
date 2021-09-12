package html

import (
	"context"

	"github.com/MontFerret/ferret/pkg/drivers"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
)

// SELECT selects a value from an underlying select element.
// @param {HTMLElement} element - Target html element.
// @param {String | String[]} valueOrSelector - Selector or a an array of strings as a value.
// @param {String[]} value - Target value. Optional.
// @return {String[]} - Array of selected values.
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

	selector, err := drivers.ToQuerySelector(args[1])

	if err != nil {
		return values.None, err
	}

	arr := values.ToArray(ctx, args[2])

	return el.SelectBySelector(ctx, selector, arr)
}
