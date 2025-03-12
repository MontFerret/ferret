package html

import (
	"context"

	"github.com/MontFerret/ferret/pkg/drivers"
	"github.com/MontFerret/ferret/pkg/runtime/core"
)

// BLUR Calls blur on the element.
// @param {HTMLPage | HTMLDocument | HTMLElement} node - Target node.
// @param {String} [selector] - CSS selector.
func Blur(ctx context.Context, args ...core.Value) (core.Value, error) {
	err := core.ValidateArgs(args, 1, 2)

	if err != nil {
		return core.None, err
	}

	el, err := drivers.ToElement(args[0])

	if err != nil {
		return core.None, err
	}

	if len(args) == 1 {
		return core.None, el.Blur(ctx)
	}

	selector, err := drivers.ToQuerySelector(args[1])

	if err != nil {
		return core.None, err
	}

	return core.None, el.BlurBySelector(ctx, selector)
}
