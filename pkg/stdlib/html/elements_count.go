package html

import (
	"context"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
)

/*
 * Returns a number of found HTML elements by a given CSS selector.
 * Returns an empty array if element not found.
 * @param docOrEl (HTMLDocument|HTMLElement) - Parent document or element.
 * @param selector (String) - CSS selector.
 * @returns (Array) - A number of found HTML elements by a given CSS selector.
 */
func ElementsCount(_ context.Context, args ...core.Value) (core.Value, error) {
	el, selector, err := queryArgs(args)

	if err != nil {
		return values.None, err
	}

	return el.CountBySelector(selector), nil
}
