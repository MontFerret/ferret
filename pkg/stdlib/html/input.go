package html

import (
	"context"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/MontFerret/ferret/pkg/runtime/values/types"
)

// Input types a value to an underlying input element.
// @param source (Document | Element) - Event target.
// @param valueOrSelector (String) - Selector or a value.
// @param value (String) - Target value.
// @param delay (Int, optional) - Waits delay milliseconds between keystrokes
// @returns (Boolean) - Returns true if an element was found.
func Input(_ context.Context, args ...core.Value) (core.Value, error) {
	err := core.ValidateArgs(args, 2, 4)

	if err != nil {
		return values.None, err
	}

	arg1 := args[0]
	err = core.ValidateType(arg1, types.HTMLDocument, types.HTMLElement)

	if err != nil {
		return values.False, err
	}

	switch args[0].(type) {
	case values.DHTMLDocument:
		doc, ok := arg1.(values.DHTMLDocument)

		if !ok {
			return values.False, core.Errors(core.ErrInvalidType, ErrNotDynamic)
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

		return doc.InputBySelector(arg2.(values.String), args[2], delay)
	case values.DHTMLNode:
		el, ok := arg1.(values.DHTMLNode)

		if !ok {
			return values.False, core.Errors(core.ErrInvalidType, ErrNotDynamic)
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

		err = el.Input(args[1], delay)

		if err != nil {
			return values.False, err
		}

		return values.True, nil
	default:
		return values.False, core.Errors(core.ErrInvalidArgument)
	}
}
