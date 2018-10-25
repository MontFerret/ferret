package html

import (
	"context"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
)

// ElementsCount returns a number of found HTML elements by a given CSS selector.
// Returns an empty array if element not found.
<<<<<<< HEAD
// @param docOrEl (HTMLDocument|HTMLNode) - Parent document or element.
// @param selector (String) - CSS selector.
// @returns (Int) - A number of found HTML elements by a given CSS selector.
func ElementsCount(ctx context.Context, args ...core.Value) (core.Value, error) {
=======
// @param docOrEl (HTMLDocument|HTMLElement) - Parent document or element.
// @param selector (String) - CSS selector.
// @returns (Int) - A number of found HTML elements by a given CSS selector.
func ElementsCount(_ context.Context, args ...core.Value) (core.Value, error) {
>>>>>>> 9f24172... rewrite comments
	el, selector, err := queryArgs(args)

	if err != nil {
		return values.None, err
	}

	return el.CountBySelector(ctx, selector), nil
}
