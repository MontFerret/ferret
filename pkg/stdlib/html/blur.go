package html

import (
	"context"
	"github.com/MontFerret/ferret/pkg/drivers"
	"github.com/MontFerret/ferret/pkg/runtime/values"

	"github.com/MontFerret/ferret/pkg/runtime/core"
)

// BLUR Calls blur on the element.
// @param target (HTMLPage | HTMLDocument | HTMLElement) - Target node.
// @param selector (String, optional) - Optional CSS selector.
func Blur(ctx context.Context, args ...core.Value) (core.Value, error) {
	err := core.ValidateArgs(args, 1, 2)

	if err != nil {
		return values.None, err
	}

	el, err := drivers.ToElement(args[0])

	if err != nil {
		return values.None, err
	}

	if len(args) == 1 {
		return values.None, el.Blur(ctx)
	}

	return values.None, el.BlurBySelector(ctx, values.ToString(args[1]))
}
