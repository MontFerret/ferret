package html

import (
	"context"

	"github.com/MontFerret/ferret/pkg/drivers"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/MontFerret/ferret/pkg/runtime/values/types"
)

// CLICK dispatches click event on a given element
// @param {HTMLPage | HTMLDocument | HTMLElement} node - Target html node.
// @param {String | Int} [cssSelectorOrClicks] - CSS selector or count of clicks.
// @param {Int} [clicks=1] - Count of clicks.
func Click(ctx context.Context, args ...core.Value) (core.Value, error) {
	err := core.ValidateArgs(args, 1, 3)

	if err != nil {
		return values.False, err
	}

	el, err := drivers.ToElement(args[0])

	if err != nil {
		return values.False, err
	}

	// CLICK(elOrDoc)
	if len(args) == 1 {
		return values.True, el.Click(ctx, 1)
	}

	if len(args) == 2 {
		err := core.ValidateType(args[1], types.String, types.Int, drivers.QuerySelectorType)

		if err != nil {
			return values.False, err
		}

		if args[1].Type() == types.String || args[1].Type() == drivers.QuerySelectorType {
			selector, err := drivers.ToQuerySelector(args[1])

			if err != nil {
				return values.None, err
			}

			exists, err := el.ExistsBySelector(ctx, selector)

			if err != nil {
				return values.False, err
			}

			if !exists {
				return exists, nil
			}

			return exists, el.ClickBySelector(ctx, selector, 1)
		}

		return values.True, el.Click(ctx, values.ToInt(args[1]))
	}

	err = core.ValidateType(args[2], types.Int)

	if err != nil {
		return values.False, err
	}

	// CLICK(doc, selector)
	selector, err := drivers.ToQuerySelector(args[1])

	if err != nil {
		return values.None, err
	}

	exists, err := el.ExistsBySelector(ctx, selector)

	if err != nil {
		return values.False, err
	}

	if !exists {
		return exists, nil
	}

	count := values.ToInt(args[2])

	return exists, el.ClickBySelector(ctx, selector, count)
}
