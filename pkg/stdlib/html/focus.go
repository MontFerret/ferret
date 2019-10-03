package html

import (
	"context"

	"github.com/MontFerret/ferret/pkg/drivers"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
)

// FOCUS Sets focus on the element.
// @param target (HTMLPage | HTMLDocument | HTMLElement) - Target node.
// @param selector (String, optional) - Optional CSS selector.
func Focus(ctx context.Context, args ...core.Value) (core.Value, error) {
	err := core.ValidateArgs(args, 1, 2)

	if err != nil {
		return values.None, err
	}

	el, err := drivers.ToElement(args[0])

	if err != nil {
		return values.None, err
	}

	if len(args) == 1 {
		return values.None, el.Focus(ctx)
	}

	return values.None, el.FocusBySelector(ctx, values.ToString(args[1]))
}
