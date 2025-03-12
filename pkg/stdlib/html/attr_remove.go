package html

import (
	"context"
	"github.com/MontFerret/ferret/pkg/drivers"
	"github.com/MontFerret/ferret/pkg/runtime/core"
)

// ATTR_REMOVE removes single or more attribute(s) of a given element.
// @param {HTMLPage | HTMLDocument | HTMLElement} node - Target node.
// @param {String, repeated} attrNames - Attribute name(s).
func AttributeRemove(ctx context.Context, args ...core.Value) (core.Value, error) {
	if err := core.ValidateArgs(args, 2, core.MaxArgs); err != nil {
		return core.None, err
	}

	el, err := drivers.ToElement(args[0])

	if err != nil {
		return core.None, err
	}

	attrs := args[1:]
	attrsStr := make([]core.String, 0, len(attrs))

	for _, attr := range attrs {
		str, err := core.CastString(attr)

		if err != nil {
			return core.None, err
		}

		attrsStr = append(attrsStr, str)
	}

	return core.None, el.RemoveAttribute(ctx, attrsStr...)
}
