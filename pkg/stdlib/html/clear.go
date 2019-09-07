package html

import (
	"context"

	"github.com/MontFerret/ferret/pkg/drivers"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
)

// INPUT_CLEAR clears a value from an underlying input element.
// @param source (HTMLPage | HTMLDocument | HTMLElement) - Event target.
// @param selector (String, options) - Selector.
func InputClear(ctx context.Context, args ...core.Value) (core.Value, error) {
	err := core.ValidateArgs(args, 1, 2)

	if err != nil {
		return values.None, err
	}

	el, err := drivers.ToElement(args[0])

	if err != nil {
		return values.None, err
	}

	// CLEAR(el)
	if len(args) == 1 {
		return values.None, el.Clear(ctx)
	}

	return values.None, el.ClearBySelector(ctx, values.ToString(args[1]))
}
