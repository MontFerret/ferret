package html

import (
	"context"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
)

// ELEMENT_EXISTS returns a boolean value indicating whether there is an element matched by selector.
// @param {HTMLPage | HTMLDocument | HTMLElement} node - Target html node.
// @param {String} selector - CSS selector.
// @return {Boolean} - A boolean value indicating whether there is an element matched by selector.
func ElementExists(ctx context.Context, args ...core.Value) (core.Value, error) {
	el, selector, err := queryArgs(args)

	if err != nil {
		return values.None, err
	}

	return el.ExistsBySelector(ctx, selector)
}
