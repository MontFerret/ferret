package html

import (
	"context"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
)

// ELEMENTS finds HTML elements by a given CSS selector.
// Returns an empty array if element not found.
// @param docOrEl (HTMLDocument|HTMLNode) - Parent document or element.
// @param selector (String) - CSS selector.
// @returns (Array) - Returns an array of found HTML element.
func Elements(ctx context.Context, args ...core.Value) (core.Value, error) {
	el, selector, err := queryArgs(args)

	if err != nil {
		return values.None, err
	}

	return el.QuerySelectorAll(ctx, selector)
}
