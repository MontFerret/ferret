package html

import (
	"context"

	"github.com/MontFerret/ferret/pkg/drivers"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/MontFerret/ferret/pkg/runtime/values/types"
)

// ATTR_SET sets or updates a single or more attribute(s) of a given element.
// @param el (HTMLElement) - Target element.
// @param nameOrObj (String | Object) - Attribute name or an object representing a key-value pair of attributes.
// @param value (String) - If a second parameter is a string value, this parameter represent an attribute value.
func AttributeSet(ctx context.Context, args ...core.Value) (core.Value, error) {
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
		// ATTR_SET(el, name, value)
		err = core.ValidateArgs(args, 3, 3)

		if err != nil {
			return values.None, nil
		}

		arg2, ok := args[2].(values.String)

		if !ok {
			return values.None, core.TypeError(arg1.Type(), types.String, types.Object)
		}

		return values.None, el.SetAttribute(ctx, arg1, arg2)
	case *values.Object:
		// ATTR_SET(el, values)
		return values.None, el.SetAttributes(ctx, arg1)
	default:
		return values.None, core.TypeError(arg1.Type(), types.String, types.Object)
	}
}
