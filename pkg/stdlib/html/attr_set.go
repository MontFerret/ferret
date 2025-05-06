package html

import (
	"context"

	"github.com/MontFerret/ferret/pkg/drivers/common"
	"github.com/MontFerret/ferret/pkg/runtime/internal"

	"github.com/MontFerret/ferret/pkg/drivers"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values/types"
)

// ATTR_SET sets or updates a single or more attribute(s) of a given element.
// @param {HTMLPage | HTMLDocument | HTMLElement} node - Target node.
// @param {String | hashMap} nameOrObj - Attribute name or an object representing a key-value pair of attributes.
// @param {String} value - If a second parameter is a string value, this parameter represent an attribute value.
func AttributeSet(ctx context.Context, args ...core.Value) (core.Value, error) {
	err := core.ValidateArgs(args, 2, core.MaxArgs)

	if err != nil {
		return core.None, err
	}

	el, err := drivers.ToElement(args[0])

	if err != nil {
		return core.None, err
	}

	switch arg1 := args[1].(type) {
	case core.String:
		// ATTR_SET(el, name, value)
		err = core.ValidateArgs(args, 3, 3)

		if err != nil {
			return core.None, nil
		}

		switch arg2 := args[2].(type) {
		case core.String:
			return core.None, el.SetAttribute(ctx, arg1, arg2)
		case *internal.Object:
			if arg1 == common.AttrNameStyle {
				return core.None, el.SetAttribute(ctx, arg1, common.SerializeStyles(ctx, arg2))
			}

			return core.None, el.SetAttribute(ctx, arg1, core.NewString(arg2.String()))
		default:
			return core.None, core.TypeError(arg1.Type(), types.String, types.Object)
		}
	case *internal.Object:
		// ATTR_SET(el, values)
		return core.None, el.SetAttributes(ctx, arg1)
	default:
		return core.None, core.TypeError(arg1.Type(), types.String, types.Object)
	}
}
