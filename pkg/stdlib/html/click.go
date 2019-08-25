package html

import (
	"context"

	"github.com/MontFerret/ferret/pkg/drivers"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
)

// Click dispatches click event on a given element
// @param source (Open | GetElement) - Event source.
// @param selector (String, optional) - Optional selector.
func Click(ctx context.Context, args ...core.Value) (core.Value, error) {
	err := core.ValidateArgs(args, 1, 2)

	if err != nil {
		return values.False, err
	}

	el, err := drivers.ToElement(args[0])

	if err != nil {
		return values.False, err
	}

	// CLICK(elOrDoc)
	if len(args) == 1 {
		return values.True, el.Click(ctx)
	}

	// CLICK(doc, selector)
	selector := values.ToString(args[1])
	exists, err := el.ExistsBySelector(ctx, selector)

	if err != nil {
		return values.False, err
	}

	if !exists {
		return exists, nil
	}

	return exists, el.ClickBySelector(ctx, selector)
}
