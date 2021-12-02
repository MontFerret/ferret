package html

import (
	"context"

	"github.com/MontFerret/ferret/pkg/drivers"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
)

// INNER_HTML returns inner HTML string of a given or matched by CSS selector element
// @param {HTMLPage | HTMLDocument | HTMLElement} node - Target html node.
// @param {String} [selector] - String of CSS selector.
// @return {String} - Inner HTML string if a matched element, otherwise empty string.
func GetInnerHTML(ctx context.Context, args ...core.Value) (core.Value, error) {
	err := core.ValidateArgs(args, 1, 2)

	if err != nil {
		return values.EmptyString, err
	}

	el, err := drivers.ToElement(args[0])

	if err != nil {
		return values.None, err
	}

	if len(args) == 1 {
		return el.GetInnerHTML(ctx)
	}

	selector, err := drivers.ToQuerySelector(args[1])

	if err != nil {
		return values.None, err
	}

	return el.GetInnerHTMLBySelector(ctx, selector)
}
