package html

import (
	"context"

	"github.com/MontFerret/ferret/pkg/drivers"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/MontFerret/ferret/pkg/runtime/values/types"
)

// INNER_HTML returns inner HTML string of a given or matched by CSS selector element
// @param doc (Open|GetElement) - Parent document or element.
// @param selector (String, optional) - String of CSS selector.
// @returns (String) - Inner HTML string if an element found, otherwise empty string.
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

	err = core.ValidateType(args[1], types.String)

	if err != nil {
		return values.None, err
	}

	selector := args[1].(values.String)

	return el.GetInnerHTMLBySelector(ctx, selector)
}
