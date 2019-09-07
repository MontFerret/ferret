package html

import (
	"context"

	"github.com/MontFerret/ferret/pkg/drivers"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/MontFerret/ferret/pkg/runtime/values/types"
)

// CLICK_ALL dispatches click event on all matched element
// @param source (Open) - Open.
// @param selector (String) - Selector.
// @param count (Int, optional) - Optional count of clicks.
// @returns (Boolean) - Returns true if matched at least one element.
func ClickAll(ctx context.Context, args ...core.Value) (core.Value, error) {
	err := core.ValidateArgs(args, 2, 3)

	if err != nil {
		return values.False, err
	}

	el, err := drivers.ToElement(args[0])

	if err != nil {
		return values.False, err
	}

	selector := values.ToString(args[1])

	exists, err := el.ExistsBySelector(ctx, selector)

	if err != nil {
		return values.False, err
	}

	if !exists {
		return values.False, nil
	}

	count := values.NewInt(1)

	if len(args) == 3 {
		err := core.ValidateType(args[2], types.Int)

		if err != nil {
			return values.False, err
		}

		count = values.ToInt(args[2])
	}

	return values.True, el.ClickBySelectorAll(ctx, selector, count)
}
