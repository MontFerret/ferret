package html

import (
	"context"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
)

// ELEMENTS_COUNT returns a number of found HTML elements by a given CSS selector.
// Returns an empty array if element not found.
// @param parent (HTMLPage | HTMLDocument | HTMLElement) - Parent document or element.
// @param selector (String) - CSS selector.
// @return (Int) - A number of matched HTML elements by a given CSS selector.
func ElementsCount(ctx context.Context, args ...core.Value) (core.Value, error) {
	el, selector, err := queryArgs(args)

	if err != nil {
		return values.None, err
	}

	return el.CountBySelector(ctx, selector)
}
