package html

import (
	"context"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
)

// ELEMENT_EXISTS returns a boolean value indicating whether there is an element matched by selector.
// @param docOrEl (HTMLDocument|HTMLNode) - Parent document or element.
// @param selector (String) - CSS selector.
// @returns (Boolean) - A boolean value indicating whether there is an element matched by selector.
func ElementExists(ctx context.Context, args ...core.Value) (core.Value, error) {
	el, selector, err := queryArgs(args)

	if err != nil {
		return values.None, err
	}

	return el.ExistsBySelector(ctx, selector)
}
