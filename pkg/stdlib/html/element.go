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

// Element finds an element by a given CSS selector.
// Returns NONE if element not found.
// @param docOrEl (HTMLDocument|HTMLElement) - Parent document or element.
// @param selector (String) - CSS selector.
// @returns (HTMLElement | None) - Returns an HTMLElement if found, otherwise NONE.
<<<<<<< HEAD
func Element(ctx context.Context, args ...core.Value) (core.Value, error) {
=======
func Element(_ context.Context, args ...core.Value) (core.Value, error) {
>>>>>>> 9f24172... rewrite comments
	el, selector, err := queryArgs(args)

	if err != nil {
		return values.None, err
	}

	return el.QuerySelector(ctx, selector), nil
}

func queryArgs(args []core.Value) (drivers.HTMLNode, values.String, error) {
	err := core.ValidateArgs(args, 2, 2)

	if err != nil {
		return nil, values.EmptyString, err
	}

	err = core.ValidateType(args[0], drivers.HTMLDocumentType, drivers.HTMLElementType)

	if err != nil {
		return nil, values.EmptyString, err
	}

	err = core.ValidateType(args[1], types.String)

	if err != nil {
		return nil, values.EmptyString, err
	}

	return args[0].(drivers.HTMLNode), args[1].(values.String), nil
}
