package html

import (
	"context"

<<<<<<< HEAD
	"github.com/MontFerret/ferret/pkg/drivers"
=======
>>>>>>> 9f24172... rewrite comments
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/MontFerret/ferret/pkg/runtime/values/types"
)

<<<<<<< HEAD
// InnerHTML Returns inner HTML string of a given or matched by CSS selector element
// @param doc (Document|Element) - Parent document or element.
// @param selector (String, optional) - String of CSS selector.
// @returns (String) - Inner HTML string if an element found, otherwise empty string.
func InnerHTML(ctx context.Context, args ...core.Value) (core.Value, error) {
	err := core.ValidateArgs(args, 1, 2)
=======
// InnerHTML Returns inner HTML string of a matched element
// @param doc (Document|Element) - Parent document or element.
// @param selector (String) - String of CSS selector.
// @returns (String) - Inner HTML string if an element found, otherwise empty string.
func InnerHTML(_ context.Context, args ...core.Value) (core.Value, error) {
	err := core.ValidateArgs(args, 2, 2)
>>>>>>> 9f24172... rewrite comments

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
