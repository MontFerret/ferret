package html

import (
	"context"

	"github.com/MontFerret/ferret/pkg/runtime/core"
)

// ELEMENTS_COUNT returns a number of found HTML elements by a given CSS selector.
// Returns an empty array if element not found.
// @param {HTMLPage | HTMLDocument | HTMLElement} node - Target html node.
// @param {String} selector - CSS selector.
// @return {Int} - A number of matched HTML elements by a given CSS selector.
func ElementsCount(ctx context.Context, args ...core.Value) (core.Value, error) {
	el, selector, err := queryArgs(args)

	if err != nil {
		return core.None, err
	}

	return el.CountBySelector(ctx, selector)
}
