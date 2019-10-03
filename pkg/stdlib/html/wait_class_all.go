package html

import (
	"context"

	"github.com/MontFerret/ferret/pkg/drivers"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/MontFerret/ferret/pkg/runtime/values/types"
)

// WAIT_CLASS_ALL waits for a class to appear on all matched elements.
// Stops the execution until the navigation ends or operation times out.
// @param doc (HTMLDocument) - Parent document.
// @param selector (String) - String of CSS selector.
// @param class (String) - String of target CSS class.
// @param timeout (Int, optional) - Optional timeout.
func WaitClassAll(ctx context.Context, args ...core.Value) (core.Value, error) {
	return waitClassAllWhen(ctx, args, drivers.WaitEventPresence)
}

// WAIT_NO_CLASS_ALL waits for a class to disappear on all matched elements.
// Stops the execution until the navigation ends or operation times out.
// @param doc (HTMLDocument) - Parent document.
// @param selector (String) - String of CSS selector.
// @param class (String) - String of target CSS class.
// @param timeout (Int, optional) - Optional timeout.
func WaitNoClassAll(ctx context.Context, args ...core.Value) (core.Value, error) {
	return waitClassAllWhen(ctx, args, drivers.WaitEventAbsence)
}

func waitClassAllWhen(ctx context.Context, args []core.Value, when drivers.WaitEvent) (core.Value, error) {
	err := core.ValidateArgs(args, 3, 4)

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

	// class
	err = core.ValidateType(args[2], types.String)

	if err != nil {
		return values.None, err
	}

	selector := args[1].(values.String)
	class := args[2].(values.String)
	timeout := values.NewInt(drivers.DefaultWaitTimeout)

	if len(args) == 4 {
		err = core.ValidateType(args[3], types.Int)

		if err != nil {
			return values.None, err
		}

		timeout = args[3].(values.Int)
	}

	ctx, fn := waitTimeout(ctx, timeout)
	defer fn()

	return values.None, doc.WaitForClassBySelectorAll(ctx, selector, class, when)
}
