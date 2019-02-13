package html

import (
	"context"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/MontFerret/ferret/pkg/runtime/values/types"
)

// Element finds an element by a given CSS selector.
// Returns NONE if element not found.
// @param docOrEl (HTMLDocument|HTMLElement) - Parent document or element.
// @param selector (String) - CSS selector.
// @returns (HTMLElement | None) - Returns an HTMLElement if found, otherwise NONE.
func Element(_ context.Context, args ...core.Value) (core.Value, error) {
	el, selector, err := queryArgs(args)

	if err != nil {
		return values.None, err
	}

	return el.QuerySelector(selector), nil
}

func queryArgs(args []core.Value) (values.HTMLNode, values.String, error) {
	err := core.ValidateArgs(args, 2, 2)

	if err != nil {
		return nil, values.EmptyString, err
	}

	err = core.ValidateType(args[0], types.HTMLDocument, types.HTMLElement)

	if err != nil {
		return nil, values.EmptyString, err
	}

	err = core.ValidateType(args[1], types.String)

	if err != nil {
		return nil, values.EmptyString, err
	}

	return args[0].(values.HTMLNode), args[1].(values.String), nil
}
