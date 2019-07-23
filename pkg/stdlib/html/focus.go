package html

import (
	"context"
	"github.com/MontFerret/ferret/pkg/drivers"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
)

// FOCUS Calls focus on the element.
// @param target (HTMLPage | HTMLDocument | HTMLElement) - Target node.
// @param selector (String, optional) - Optional CSS selector. Required when target is HTMLPage or HTMLDocument.
func Focus(ctx context.Context, args ...core.Value) (core.Value, error) {
	err := core.ValidateArgs(args, 1, 2)

	if err != nil {
		return values.None, err
	}

	// Document with selector
	if len(args) == 2 {
		doc, err := drivers.ToDocument(args[0])

		if err != nil {
			return values.None, err
		}

		return values.None, doc.FocusBySelector(ctx, values.ToString(args[1]))
	}

	el, err := drivers.ToElement(args[0])

	if err != nil {
		return values.None, err
	}

	return values.None, el.Focus(ctx)
}
