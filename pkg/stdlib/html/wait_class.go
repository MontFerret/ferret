package html

import (
	"context"

	"github.com/MontFerret/ferret/pkg/drivers"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/MontFerret/ferret/pkg/runtime/values/types"
)

// WAIT_CLASS waits for a class to appear on a given element.
// Stops the execution until the navigation ends or operation times out.
// @param {HTMLPage | HTMLDocument | HTMLElement} node - Target html node.
// @param {String} selectorOrClass - If document is passed, this param must represent an element selector. Otherwise target class.
// @param {String | Int} [classOrTimeout] - If document is passed, this param must represent target class name. Otherwise timeout.
// @param {Int} [timeout] - If document is passed, this param must represent timeout. Otherwise not passed.
func WaitClass(ctx context.Context, args ...core.Value) (core.Value, error) {
	return waitClassWhen(ctx, args, drivers.WaitEventPresence)
}

// WAIT_NO_CLASS waits for a class to disappear on a given element.
// Stops the execution until the navigation ends or operation times out.
// @param {HTMLPage | HTMLDocument | HTMLElement} node - Target html node.
// @param {String} selectorOrClass - If document is passed, this param must represent an element selector. Otherwise target class.
// @param {String | Int} [classOrTimeout] - If document is passed, this param must represent target class name. Otherwise timeout.
// @param {Int} [timeout] - If document is passed, this param must represent timeout. Otherwise not passed.
func WaitNoClass(ctx context.Context, args ...core.Value) (core.Value, error) {
	return waitClassWhen(ctx, args, drivers.WaitEventAbsence)
}

func waitClassWhen(ctx context.Context, args []core.Value, when drivers.WaitEvent) (core.Value, error) {
	err := core.ValidateArgs(args, 2, 4)

	if err != nil {
		return values.None, err
	}

	// document or element
	arg1 := args[0]
	err = core.ValidateType(arg1, drivers.HTMLPageType, drivers.HTMLDocumentType, drivers.HTMLElementType)

	if err != nil {
		return values.None, err
	}

	timeout := values.NewInt(drivers.DefaultWaitTimeout)

	// if a document is passed
	if arg1.Type() == drivers.HTMLPageType || arg1.Type() == drivers.HTMLDocumentType {
		// revalidate args with more accurate amount
		err := core.ValidateArgs(args, 3, 4)

		if err != nil {
			return values.None, err
		}

		selector, err := drivers.ToQuerySelector(args[1])

		if err != nil {
			return values.None, err
		}

		// class
		err = core.ValidateType(args[2], types.String)

		if err != nil {
			return values.None, err
		}

		el, err := drivers.ToElement(arg1)

		if err != nil {
			return values.None, err
		}

		class := args[2].(values.String)

		if len(args) == 4 {
			err = core.ValidateType(args[3], types.Int)

			if err != nil {
				return values.None, err
			}

			timeout = args[3].(values.Int)
		}

		ctx, fn := waitTimeout(ctx, timeout)
		defer fn()

		return values.True, el.WaitForClassBySelector(ctx, selector, class, when)
	}

	el := arg1.(drivers.HTMLElement)
	class := args[1].(values.String)

	if len(args) == 3 {
		err = core.ValidateType(args[2], types.Int)

		if err != nil {
			return values.None, err
		}

		timeout = args[2].(values.Int)
	}

	ctx, fn := waitTimeout(ctx, timeout)
	defer fn()

	return values.True, el.WaitForClass(ctx, class, when)
}
