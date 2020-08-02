package html

import (
	"context"

	"github.com/MontFerret/ferret/pkg/drivers"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
)

// STYLE_GET gets single or more style attribute value(s) of a given element.
// @param {HTMLElement} element - Target html element.
// @param {String, repeated} names - Style name(s).
// @return {Object} - Collection of key-value pairs of style values.
func StyleGet(ctx context.Context, args ...core.Value) (core.Value, error) {
	err := core.ValidateArgs(args, 2, core.MaxArgs)

	if err != nil {
		return values.None, err
	}

	el, err := drivers.ToElement(args[0])

	if err != nil {
		return values.None, err
	}

	names := args[1:]
	result := values.NewObject()

	for _, n := range names {
		name := values.NewString(n.String())
		val, err := el.GetStyle(ctx, name)

		if err != nil {
			return values.None, err
		}

		if val != values.None {
			result.Set(name, val)
		}
	}

	return result, nil
}
