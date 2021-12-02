package html

import (
	"context"

	"github.com/MontFerret/ferret/pkg/drivers"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
)

// ATTR_QUERY finds a single or more attribute(s) by an query selector.
// @param {HTMLPage | HTMLDocument | HTMLElement} node - Target node.
// @param {String} selector - Query selector.
// @param {String, repeated} attrName - Attr name(s).
// @return {Object} - Key-value pairs of attribute values.
func AttributeQuery(ctx context.Context, args ...core.Value) (core.Value, error) {
	err := core.ValidateArgs(args, 2, core.MaxArgs)

	if err != nil {
		return values.None, err
	}

	parent, err := drivers.ToElement(args[0])

	if err != nil {
		return values.None, err
	}

	selector, err := drivers.ToQuerySelector(args[1])

	if err != nil {
		return values.None, err
	}

	found, err := parent.QuerySelector(ctx, selector)

	if err != nil {
		return values.None, err
	}

	el, err := drivers.ToElement(found)

	if err != nil {
		return values.None, err
	}

	names := args[2:]
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
