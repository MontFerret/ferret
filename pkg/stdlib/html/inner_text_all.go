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

// InnerTextAll returns an array of inner text of matched elements.
// @param doc (HTMLDocument|HTMLElement) - Parent document or element.
// @param selector (String) - String of CSS selector.
// @returns (String) - An array of inner text if any element found, otherwise empty array.
<<<<<<< HEAD
func InnerTextAll(ctx context.Context, args ...core.Value) (core.Value, error) {
=======
func InnerTextAll(_ context.Context, args ...core.Value) (core.Value, error) {
>>>>>>> 9f24172... rewrite comments
	err := core.ValidateArgs(args, 2, 2)

	if err != nil {
		return values.None, err
	}

	err = core.ValidateType(args[0], drivers.HTMLDocumentType, drivers.HTMLElementType)

	if err != nil {
		return values.None, err
	}

	err = core.ValidateType(args[1], types.String)

	if err != nil {
		return values.None, err
	}

	el, err := resolveElement(args[0])

	if err != nil {
		return values.None, err
	}

	selector := args[1].(values.String)

	return el.InnerTextBySelectorAll(ctx, selector), nil
}
