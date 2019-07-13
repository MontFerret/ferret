package html

import (
	"context"

	"github.com/MontFerret/ferret/pkg/drivers"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/MontFerret/ferret/pkg/runtime/values/types"
)

// SetInnerText sets inner text string to a given or matched by CSS selector element
// @param doc (Open|GetElement) - Parent document or element.
// @param selector (String, optional) - String of CSS selector.
// @param innerText (String) - String of inner text.
func SetInnerText(ctx context.Context, args ...core.Value) (core.Value, error) {
	err := core.ValidateArgs(args, 2, 3)

	if err != nil {
		return values.None, err
	}

	el, err := drivers.ToElement(args[0])

	if err != nil {
		return values.None, err
	}

	if len(args) == 2 {
		err := core.ValidateType(args[1], types.String)

		if err != nil {
			return values.None, err
		}

		return values.None, el.SetInnerText(ctx, values.ToString(args[1]))
	}

	err = core.ValidateType(args[1], types.String)

	if err != nil {
		return values.None, err
	}

	err = core.ValidateType(args[2], types.String)

	if err != nil {
		return values.None, err
	}

	selector := values.ToString(args[1])
	innerHTML := values.ToString(args[2])

	return values.None, el.SetInnerTextBySelector(ctx, selector, innerHTML)
}
