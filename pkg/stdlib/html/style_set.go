package html

import (
	"context"

	"github.com/MontFerret/ferret/pkg/drivers"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/internal"
	"github.com/MontFerret/ferret/pkg/runtime/values/types"
)

// STYLE_SET sets or updates a single or more style attribute value of a given element.
// @param {HTMLElement} element - Target html element.
// @param {String | hashMap} nameOrObj - Style name or an object representing a key-value pair of attributes.
// @param {String} value - If a second parameter is a string value, this parameter represent a style value.
func StyleSet(ctx context.Context, args ...core.Value) (core.Value, error) {
	err := core.ValidateArgs(args, 2, 3)

	if err != nil {
		return core.None, err
	}

	el, err := drivers.ToElement(args[0])

	if err != nil {
		return core.None, err
	}

	switch arg1 := args[1].(type) {
	case core.String:
		// STYLE_SET(el, name, value)
		err = core.ValidateArgs(args, 3, 3)

		if err != nil {
			return core.None, nil
		}

		return core.None, el.SetStyle(ctx, arg1, core.NewString(args[2].String()))
	case *internal.Object:
		// STYLE_SET(el, values)
		return core.None, el.SetStyles(ctx, arg1)
	default:
		return core.None, core.TypeError(arg1.Type(), types.String, types.Object)
	}
}
