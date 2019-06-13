package html

import (
	"context"

	"github.com/MontFerret/ferret/pkg/drivers"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/MontFerret/ferret/pkg/runtime/values/types"
)

// Input types a value to an underlying input element.
// @param source (Open | GetElement) - Event target.
// @param valueOrSelector (String) - Selector or a value.
// @param value (String) - Target value.
// @param delay (Int, optional) - Waits delay milliseconds between keystrokes
// @returns (Boolean) - Returns true if an element was found.
func Input(ctx context.Context, args ...core.Value) (core.Value, error) {
	err := core.ValidateArgs(args, 2, 4)

	if err != nil {
		return values.None, err
	}

	arg1 := args[0]
	err = core.ValidateType(arg1, drivers.HTMLPageType, drivers.HTMLDocumentType, drivers.HTMLElementType)

	if err != nil {
		return values.False, err
	}

	if arg1.Type() == drivers.HTMLPageType || arg1.Type() == drivers.HTMLDocumentType {
		doc, err := drivers.ToDocument(arg1)

		if err != nil {
			return values.False, err
		}

		// selector
		arg2 := args[1]
		err = core.ValidateType(arg2, types.String)

		if err != nil {
			return values.False, err
		}

		delay := values.Int(0)

		if len(args) == 4 {
			arg4 := args[3]

			err = core.ValidateType(arg4, types.Int)

			if err != nil {
				return values.False, err
			}

			delay = arg4.(values.Int)
		}

		return doc.InputBySelector(ctx, arg2.(values.String), args[2], delay)
	}

	el, err := drivers.ToElement(arg1)

	if err != nil {
		return values.None, err
	}

	delay := values.Int(0)

	if len(args) == 3 {
		arg3 := args[2]

		err = core.ValidateType(arg3, types.Int)

		if err != nil {
			return values.False, err
		}

		delay = arg3.(values.Int)
	}

	err = el.Input(ctx, args[1], delay)

	if err != nil {
		return values.False, err
	}

	return values.True, nil
}
