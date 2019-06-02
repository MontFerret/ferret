package html

import (
	"context"

	"github.com/MontFerret/ferret/pkg/drivers"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/MontFerret/ferret/pkg/runtime/values/types"
)

// InnerHTML Returns inner HTML string of a given or matched by CSS selector element
// @param doc (Open|Element) - Parent document or element.
// @param selector (String, optional) - String of CSS selector.
// @returns (String) - Inner HTML string if an element found, otherwise empty string.
func InnerHTML(ctx context.Context, args ...core.Value) (core.Value, error) {
	err := core.ValidateArgs(args, 1, 2)

	if err != nil {
		return values.EmptyString, err
	}

	err = core.ValidateType(args[0], drivers.HTMLDocumentType, drivers.HTMLElementType)

	if err != nil {
		return values.None, err
	}

	el, err := resolveElement(args[0])

	if err != nil {
		return values.None, err
	}

	if len(args) == 1 {
		return el.InnerHTML(ctx), nil
	}

	err = core.ValidateType(args[1], types.String)

	if err != nil {
		return values.None, err
	}

	selector := args[1].(values.String)

	return el.InnerHTMLBySelector(ctx, selector), nil
}
