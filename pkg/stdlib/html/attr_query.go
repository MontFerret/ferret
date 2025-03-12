package html

import (
	"context"
	"github.com/MontFerret/ferret/pkg/runtime/internal"

	"github.com/MontFerret/ferret/pkg/drivers"
	"github.com/MontFerret/ferret/pkg/runtime/core"
)

// ATTR_QUERY finds a single or more attribute(s) by an query selector.
// @param {HTMLPage | HTMLDocument | HTMLElement} node - Target node.
// @param {String} selector - Query selector.
// @param {String, repeated} attrName - Attr name(s).
// @return {Object} - First-value pairs of attribute values.
func AttributeQuery(ctx context.Context, args ...core.Value) (core.Value, error) {
	err := core.ValidateArgs(args, 2, core.MaxArgs)

	if err != nil {
		return core.None, err
	}

	parent, err := drivers.ToElement(args[0])

	if err != nil {
		return core.None, err
	}

	selector, err := drivers.ToQuerySelector(args[1])

	if err != nil {
		return core.None, err
	}

	found, err := parent.QuerySelector(ctx, selector)

	if err != nil {
		return core.None, err
	}

	el, err := drivers.ToElement(found)

	if err != nil {
		return core.None, err
	}

	names := args[2:]
	result := internal.NewObject()
	attrs, err := el.GetAttributes(ctx)

	if err != nil {
		return core.None, err
	}

	for _, n := range names {
		name := core.NewString(n.String())
		val, exists := attrs.Get(name)

		if exists {
			result.Set(name, val)
		}
	}

	return result, nil
}
