package html

import (
	"context"

	"github.com/MontFerret/ferret/pkg/drivers"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/MontFerret/ferret/pkg/runtime/values/types"
)

// WAIT_ELEMENT waits for element to appear in the DOM.
// Stops the execution until it finds an element or operation times out.
// @param {HTMLPage | HTMLDocument | HTMLElement} node - Target html node.
// @param {String} selector - Target element's selector.
// @param {Int} [timeout=5000] - Wait timeout.
func WaitElement(ctx context.Context, args ...core.Value) (core.Value, error) {
	return waitElementWhen(ctx, args, drivers.WaitEventPresence)
}

// WAIT_NO_ELEMENT waits for element to disappear in the DOM.
// Stops the execution until it does not find an element or operation times out.
// @param {HTMLPage | HTMLDocument | HTMLElement} node - Target html node.
// @param {String} selector - Target element's selector.
// @param {Int} [timeout=5000] - Wait timeout.
func WaitNoElement(ctx context.Context, args ...core.Value) (core.Value, error) {
	return waitElementWhen(ctx, args, drivers.WaitEventAbsence)
}

func waitElementWhen(ctx context.Context, args []core.Value, when drivers.WaitEvent) (core.Value, error) {
	err := core.ValidateArgs(args, 2, 3)

	if err != nil {
		return values.None, err
	}

	el, err := drivers.ToElement(args[0])

	if err != nil {
		return values.None, err
	}

	selector, err := drivers.ToQuerySelector(args[1])

	if err != nil {
		return values.None, err
	}

	timeout := values.NewInt(drivers.DefaultWaitTimeout)

	if len(args) > 2 {
		err = core.ValidateType(args[2], types.Int)

		if err != nil {
			return values.None, err
		}

		timeout = args[2].(values.Int)
	}

	ctx, fn := waitTimeout(ctx, timeout)
	defer fn()

	return values.True, el.WaitForElement(ctx, selector, when)
}
