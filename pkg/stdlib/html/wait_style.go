package html

import (
	"context"

	"github.com/MontFerret/ferret/pkg/drivers"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/MontFerret/ferret/pkg/runtime/values/types"
)

// WAIT_STYLE waits until a target style value appears
// @param {HTMLPage | HTMLDocument | HTMLElement} node - Target html node.
// @param {String} styleNameOrSelector - Style name or CSS selector.
// @param {String | Any} valueOrStyleName - Style value or name.
// @param {Any | Int} [valueOrTimeout] - Style value or wait timeout.
// @param {Int} [timeout=5000] - Wait timeout.
func WaitStyle(ctx context.Context, args ...core.Value) (core.Value, error) {
	return waitStyleWhen(ctx, args, drivers.WaitEventPresence)
}

// WAIT_NO_STYLE waits until a target style value disappears
// @param {HTMLPage | HTMLDocument | HTMLElement} node - Target html node.
// @param {String} styleNameOrSelector - Style name or CSS selector.
// @param {String | Any} valueOrStyleName - Style value or name.
// @param {Any | Int} [valueOrTimeout] - Style value or wait timeout.
// @param {Int} [timeout=5000] - Wait timeout.
func WaitNoStyle(ctx context.Context, args ...core.Value) (core.Value, error) {
	return waitStyleWhen(ctx, args, drivers.WaitEventAbsence)
}

func waitStyleWhen(ctx context.Context, args []core.Value, when drivers.WaitEvent) (core.Value, error) {
	err := core.ValidateArgs(args, 3, 5)

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
	// WAIT_ATTR(doc, selector, attrName, attrValue, timeout)
	if arg1.Type() == drivers.HTMLPageType || arg1.Type() == drivers.HTMLDocumentType {
		// revalidate args with more accurate amount
		err := core.ValidateArgs(args, 4, 5)

		if err != nil {
			return values.None, err
		}

		selector, err := drivers.ToQuerySelector(args[1])

		if err != nil {
			return values.None, err
		}

		// attr name
		err = core.ValidateType(args[2], types.String)

		if err != nil {
			return values.None, err
		}

		el, err := drivers.ToElement(arg1)

		if err != nil {
			return values.None, err
		}

		name := args[2].(values.String)
		value := args[3]

		if len(args) == 5 {
			err = core.ValidateType(args[4], types.Int)

			if err != nil {
				return values.None, err
			}

			timeout = args[4].(values.Int)
		}

		ctx, fn := waitTimeout(ctx, timeout)
		defer fn()

		return values.True, el.WaitForStyleBySelector(ctx, selector, name, value, when)
	}

	el := arg1.(drivers.HTMLElement)
	name := args[1].(values.String)
	value := args[2]

	if len(args) == 4 {
		err = core.ValidateType(args[3], types.Int)

		if err != nil {
			return values.None, err
		}

		timeout = args[3].(values.Int)
	}

	ctx, fn := waitTimeout(ctx, timeout)
	defer fn()

	return values.True, el.WaitForStyle(ctx, name, value, when)
}
