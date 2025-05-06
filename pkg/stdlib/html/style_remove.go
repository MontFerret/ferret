package html

import (
	"context"

	"github.com/MontFerret/ferret/pkg/drivers"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values/types"
)

// STYLE_REMOVE removes single or more style attribute value(s) of a given element.
// @param {HTMLElement} element - Target html element.
// @param {String, repeated} names - Style name(s).
func StyleRemove(ctx context.Context, args ...core.Value) (core.Value, error) {
	err := core.ValidateArgs(args, 2, core.MaxArgs)

	if err != nil {
		return core.None, err
	}

	el, err := drivers.ToElement(args[0])

	if err != nil {
		return core.None, err
	}

	attrs := args[1:]
	attrsStr := make([]core.String, 0, len(attrs))

	for _, attr := range attrs {
		str, ok := attr.(core.String)

		if !ok {
			return core.None, core.TypeError(attr.Type(), types.String)
		}

		attrsStr = append(attrsStr, str)
	}

	return core.None, el.RemoveStyle(ctx, attrsStr...)
}
