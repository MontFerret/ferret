package html

import (
	"context"

	"github.com/MontFerret/ferret/pkg/drivers"
	"github.com/MontFerret/ferret/pkg/runtime/core"
)

// INPUT_CLEAR clears a value from an underlying input element.
// @param {HTMLPage | HTMLDocument | HTMLElement} node - Target html node.
// @param {String} [selector] - CSS selector.
func InputClear(ctx context.Context, args ...core.Value) (core.Value, error) {
	err := core.ValidateArgs(args, 1, 2)

	if err != nil {
		return core.None, err
	}

	el, err := drivers.ToElement(args[0])

	if err != nil {
		return core.None, err
	}

	// CLEAR(el)
	if len(args) == 1 {
		return core.None, el.Clear(ctx)
	}

	selector, err := drivers.ToQuerySelector(args[1])

	if err != nil {
		return core.None, err
	}

	return core.True, el.ClearBySelector(ctx, selector)
}
