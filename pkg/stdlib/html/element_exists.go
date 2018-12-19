package html

import (
	"context"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
)

// ElementExists returns a boolean value indicating whether there is an element matched by selector.
// @param docOrEl (HTMLDocument|HTMLElement) - Parent document or element.
// @param selector (String) - CSS selector.
// @returns (Boolean) - A boolean value indicating whether there is an element matched by selector.
func ElementExists(_ context.Context, args ...core.Value) (core.Value, error) {
	el, selector, err := queryArgs(args)

	if err != nil {
		return values.None, err
	}

	return el.ExistsBySelector(selector), nil
}
