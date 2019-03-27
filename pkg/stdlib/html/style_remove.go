package html

import (
	"context"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/MontFerret/ferret/pkg/runtime/values/types"
)

// StyleRemove removes single or more style attribute value(s) of a given element.
// @param el (HTMLElement) - Target element.
// @param names (...String) - Style name(s).
func StyleRemove(ctx context.Context, args ...core.Value) (core.Value, error) {
	err := core.ValidateArgs(args, 2, core.MaxArgs)

	if err != nil {
		return values.None, err
	}

	el, err := resolveElement(args[0])

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

	return values.None, el.RemoveStyle(ctx, attrsStr...)
}
