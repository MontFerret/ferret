package html

import (
	"context"

	"github.com/MontFerret/ferret/pkg/drivers"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/MontFerret/ferret/pkg/runtime/values/types"
)

// ATTR_REMOVE removes single or more attribute(s) of a given element.
// @param {HTMLPage | HTMLDocument | HTMLElement} node - Target node.
// @param {String, repeated} attrNames - Attribute name(s).
func AttributeRemove(ctx context.Context, args ...core.Value) (core.Value, error) {
	err := core.ValidateArgs(args, 2, core.MaxArgs)

	if err != nil {
		return values.None, err
	}

	el, err := drivers.ToElement(args[0])

	if err != nil {
		return values.None, err
	}

	attrs := args[1:]
	attrsStr := make([]values.String, 0, len(attrs))

	for _, attr := range attrs {
		str, ok := attr.(values.String)

		if !ok {
			return values.None, core.TypeError(attr.Type(), types.String)
		}

		attrsStr = append(attrsStr, str)
	}

	return values.None, el.RemoveAttribute(ctx, attrsStr...)
}
