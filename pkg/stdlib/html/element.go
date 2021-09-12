package html

import (
	"context"

	"github.com/MontFerret/ferret/pkg/drivers"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/MontFerret/ferret/pkg/runtime/values/types"
)

// ELEMENT finds an element by a given CSS selector.
// Returns NONE if element not found.
// @param {HTMLPage | HTMLDocument | HTMLElement} node - Target html node.
// @param {String} selector - CSS selector.
// @return {HTMLElement} - A matched HTML element
func Element(ctx context.Context, args ...core.Value) (core.Value, error) {
	el, selector, err := queryArgs(args)

	if err != nil {
		return values.None, err
	}

	return el.QuerySelector(ctx, drivers.ToQuerySelector(selector))
}

func queryArgs(args []core.Value) (drivers.HTMLElement, drivers.QuerySelector, error) {
	err := core.ValidateArgs(args, 2, 2)

	if err != nil {
		return nil, drivers.QuerySelector{}, err
	}

	el, err := drivers.ToElement(args[0])

	if err != nil {
		return nil, drivers.QuerySelector{}, err
	}

	err = core.ValidateType(args[1], types.String, drivers.QuerySelectorType)

	if err != nil {
		return nil, drivers.QuerySelector{}, err
	}

	return el, drivers.ToQuerySelector(args[1]), nil
}
