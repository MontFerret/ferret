package html

import (
	"context"
	"github.com/MontFerret/ferret/pkg/drivers"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/internal"
	"github.com/MontFerret/ferret/pkg/runtime/values/types"
)

// CLICK_ALL dispatches click event on all matched element
// @param {HTMLPage | HTMLDocument | HTMLElement} node - Target html node.
// @param {String} selector - CSS selector.
// @param {Int} [clicks=1] - Optional count of clicks.
// @return {Boolean} - True if matched at least one element.
func ClickAll(ctx context.Context, args ...core.Value) (core.Value, error) {
	err := core.ValidateArgs(args, 2, 3)

	if err != nil {
		return core.False, err
	}

	el, err := drivers.ToElement(args[0])

	if err != nil {
		return core.False, err
	}

	selector, err := drivers.ToQuerySelector(args[1])

	if err != nil {
		return core.None, err
	}

	exists, err := el.ExistsBySelector(ctx, selector)

	if err != nil {
		return core.False, err
	}

	if !exists {
		return core.False, nil
	}

	count := core.NewInt(1)

	if len(args) == 3 {
		err := core.ValidateType(args[2], types.Int)

		if err != nil {
			return core.False, err
		}

		count = internal.ToInt(args[2])
	}

	return core.True, el.ClickBySelectorAll(ctx, selector, count)
}
