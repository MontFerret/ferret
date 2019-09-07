package html

import (
	"context"

	"github.com/MontFerret/ferret/pkg/drivers"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/MontFerret/ferret/pkg/runtime/values/types"
)

// STYLE_SET sets or updates a single or more style attribute value of a given element.
// @param el (HTMLElement) - Target element.
// @param nameOrObj (String | Object) - Style name or an object representing a key-value pair of attributes.
// @param value (String) - If a second parameter is a string value, this parameter represent a style value.
func StyleSet(ctx context.Context, args ...core.Value) (core.Value, error) {
	err := core.ValidateArgs(args, 2, core.MaxArgs)

	if err != nil {
		return values.None, err
	}

	el, err := drivers.ToElement(args[0])

	if err != nil {
		return values.None, err
	}

	switch arg1 := args[1].(type) {
	case values.String:
		// STYLE_SET(el, name, value)
		err = core.ValidateArgs(args, 3, 3)

		if err != nil {
			return values.None, nil
		}

		return values.None, el.SetStyle(ctx, arg1, args[2])
	case *values.Object:
		// STYLE_SET(el, values)
		return values.None, el.SetStyles(ctx, arg1)
	default:
		return values.None, core.TypeError(arg1.Type(), types.String, types.Object)
	}
}
