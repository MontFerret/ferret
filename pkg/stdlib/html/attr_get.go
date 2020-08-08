package html

import (
	"context"

	"github.com/MontFerret/ferret/pkg/drivers"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
)

// ATTR_GET gets single or more attribute(s) of a given element.
// @param {HTMLPage | HTMLDocument | HTMLElement} node - Target node.
// @param {String, repeated} attrNames - Attribute name(s).
// @return {Object} - Key-value pairs of attribute values.
func AttributeGet(ctx context.Context, args ...core.Value) (core.Value, error) {
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
	attrs, err := el.GetAttributes(ctx)

	if err != nil {
		return values.None, err
	}

	for _, n := range names {
		name := values.NewString(n.String())
		val, exists := attrs.Get(name)

		if exists {
			result.Set(name, val)
		}
	}

	return result, nil
}
