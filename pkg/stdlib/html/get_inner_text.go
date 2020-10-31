package html

import (
	"context"

	"github.com/MontFerret/ferret/pkg/drivers"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/MontFerret/ferret/pkg/runtime/values/types"
)

// INNER_TEXT returns inner text string of a given or matched by CSS selector element
// @param {HTMLPage | HTMLDocument | HTMLElement} node - Target html node.
// @param {String} [selector] - String of CSS selector.
// @return {String} - Inner text if a matched element, otherwise empty string.
func GetInnerText(ctx context.Context, args ...core.Value) (core.Value, error) {
	err := core.ValidateArgs(args, 1, 2)

	if err != nil {
		return values.EmptyString, err
	}

	el, err := drivers.ToElement(args[0])

	if err != nil {
		return values.None, err
	}

	if len(args) == 1 {
		return el.GetInnerText(ctx)
	}

	err = core.ValidateType(args[1], types.String)

	if err != nil {
		return values.None, err
	}

	selector := args[1].(values.String)

	return el.GetInnerTextBySelector(ctx, selector)
}
