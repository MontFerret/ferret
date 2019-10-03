package html

import (
	"context"

	"github.com/MontFerret/ferret/pkg/drivers"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/MontFerret/ferret/pkg/runtime/values/types"
)

// WAIT_ATTR_ALL waits for an attribute to appear on all matched elements with a given value.
// Stops the execution until the navigation ends or operation times out.
// @param doc (HTMLDocument) - Parent document.
// @param selector (String) - String of CSS selector.
// @param class (String) - String of target CSS class.
// @param timeout (Int, optional) - Optional timeout.
func WaitAttributeAll(ctx context.Context, args ...core.Value) (core.Value, error) {
	return waitAttributeAllWhen(ctx, args, drivers.WaitEventPresence)
}

// WAIT_NO_ATTR_ALL waits for an attribute to disappear on all matched elements by a given value.
// Stops the execution until the navigation ends or operation times out.
// @param doc (HTMLDocument) - Parent document.
// @param selector (String) - String of CSS selector.
// @param class (String) - String of target CSS class.
// @param timeout (Int, optional) - Optional timeout.
func WaitNoAttributeAll(ctx context.Context, args ...core.Value) (core.Value, error) {
	return waitAttributeAllWhen(ctx, args, drivers.WaitEventAbsence)
}

func waitAttributeAllWhen(ctx context.Context, args []core.Value, when drivers.WaitEvent) (core.Value, error) {
	err := core.ValidateArgs(args, 4, 5)

	if err != nil {
		return values.None, err
	}

	doc, err := drivers.ToDocument(args[0])

	if err != nil {
		return values.None, err
	}

	// selector
	err = core.ValidateType(args[1], types.String)

	if err != nil {
		return values.None, err
	}

	// attr name
	err = core.ValidateType(args[2], types.String)

	if err != nil {
		return values.None, err
	}

	selector := args[1].(values.String)
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

	return values.None, doc.WaitForAttributeBySelectorAll(ctx, selector, name, value, when)
}
