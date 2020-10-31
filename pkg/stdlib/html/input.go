package html

import (
	"context"
	"github.com/MontFerret/ferret/pkg/drivers"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/MontFerret/ferret/pkg/runtime/values/types"
)

// INPUT types a value to an underlying input element.
// @param {HTMLPage | HTMLDocument | HTMLElement} node - Target html node.
// @param {String} valueOrSelector - CSS selector or a value.
// @param {String} value - Target value.
// @param {Int} [delay] - Target value.
// @return {Boolean} - Returns true if an element was found.
func Input(ctx context.Context, args ...core.Value) (core.Value, error) {
	err := core.ValidateArgs(args, 2, 4)

	if err != nil {
		return values.False, err
	}

	el, err := drivers.ToElement(args[0])

	if err != nil {
		return values.False, err
	}

	delay := values.NewInt(drivers.DefaultKeyboardDelay)

	// INPUT(el, value)
	if len(args) == 2 {
		return values.True, el.Input(ctx, args[1], delay)
	}

	var selector values.String
	var value core.Value

	// INPUT(el, valueOrSelector, valueOrOpts)
	if len(args) == 3 {
		switch v := args[2].(type) {
		// INPUT(el, value, delay)
		case values.Int, values.Float:
			value = args[1]
			delay = values.ToInt(v)

			return values.True, el.Input(ctx, value, delay)
		default:
			// INPUT(el, selector, value)
			if err := core.ValidateType(args[1], types.String); err != nil {
				return values.False, err
			}

			selector = values.ToString(args[1])
			value = args[2]
		}
	} else {
		// INPUT(el, selector, value, delay)
		if err := core.ValidateType(args[1], types.String); err != nil {
			return values.False, err
		}

		if err := core.ValidateType(args[3], types.Int); err != nil {
			return values.False, err
		}

		selector = values.ToString(args[1])
		value = args[2]
		delay = values.ToInt(args[3])
	}

	exists, err := el.ExistsBySelector(ctx, selector)

	if err != nil {
		return values.False, err
	}

	if !exists {
		return values.False, nil
	}

	return values.True, el.InputBySelector(ctx, selector, value, delay)
}
