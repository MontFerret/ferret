package html

import (
	"context"

	"github.com/MontFerret/ferret/pkg/drivers"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/MontFerret/ferret/pkg/runtime/values/types"
)

// WAIT_STYLE_ALL waits until a target style value appears on all matched elements with a given value.
// @param {HTMLPage | HTMLDocument | HTMLElement} node - Target html node.
// @param {String} styleNameOrSelector - Style name or CSS selector.
// @param {String | Any} valueOrStyleName - Style value or name.
// @param {Any | Int} [valueOrTimeout] - Style value or wait timeout.
// @param {Int} [timeout=5000] - Timeout.
func WaitStyleAll(ctx context.Context, args ...core.Value) (core.Value, error) {
	return waitStyleAllWhen(ctx, args, drivers.WaitEventPresence)
}

// WAIT_NO_STYLE_ALL waits until a target style value disappears on all matched elements with a given value.
// @param {HTMLPage | HTMLDocument | HTMLElement} node - Target html node.
// @param {String} styleNameOrSelector - Style name or CSS selector.
// @param {String | Any} valueOrStyleName - Style value or name.
// @param {Any | Int} [valueOrTimeout] - Style value or wait timeout.
// @param {Int} [timeout=5000] - Timeout.
func WaitNoStyleAll(ctx context.Context, args ...core.Value) (core.Value, error) {
	return waitStyleAllWhen(ctx, args, drivers.WaitEventAbsence)
}

func waitStyleAllWhen(ctx context.Context, args []core.Value, when drivers.WaitEvent) (core.Value, error) {
	err := core.ValidateArgs(args, 4, 5)

	if err != nil {
		return values.None, err
	}

	el, err := drivers.ToElement(args[0])

	if err != nil {
		return values.None, err
	}

	// selector
	selector, err := drivers.ToQuerySelector(args[1])

	if err != nil {
		return values.None, err
	}

	// attr name
	err = core.ValidateType(args[2], types.String)

	if err != nil {
		return values.None, err
	}

	name := args[2].(values.String)
	value := args[3]
	timeout := values.NewInt(drivers.DefaultWaitTimeout)

	if len(args) == 5 {
		err = core.ValidateType(args[4], types.Int)

		if err != nil {
			return values.None, err
		}

		timeout = args[4].(values.Int)
	}

	ctx, fn := waitTimeout(ctx, timeout)
	defer fn()

	return values.True, el.WaitForStyleBySelectorAll(ctx, selector, name, value, when)
}
