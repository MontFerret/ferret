package html

import (
	"context"
	"github.com/MontFerret/ferret/pkg/runtime/internal"

	"github.com/MontFerret/ferret/pkg/drivers"
	"github.com/MontFerret/ferret/pkg/runtime/core"
)

// STYLE_GET gets single or more style attribute value(s) of a given element.
// @param {HTMLElement} element - Target html element.
// @param {String, repeated} names - Style name(s).
// @return {hashMap} - Collection of key-value pairs of style values.
func StyleGet(ctx context.Context, args ...core.Value) (core.Value, error) {
	err := core.ValidateArgs(args, 2, core.MaxArgs)

	if err != nil {
		return core.None, err
	}

	el, err := drivers.ToElement(args[0])

	if err != nil {
		return core.None, err
	}

	names := args[1:]
	result := internal.NewObject()

	for _, n := range names {
		name := core.NewString(n.String())
		val, err := el.GetStyle(ctx, name)

		if err != nil {
			return core.None, err
		}

		if val != core.None {
			result.Set(name, val)
		}
	}

	return result, nil
}
