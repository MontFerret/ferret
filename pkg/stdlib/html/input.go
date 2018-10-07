package html

import (
	"context"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/MontFerret/ferret/pkg/stdlib/html/driver/dynamic"
)

/*
 * Sends a value to an underlying input element.
 * @param source (Document | Element) - Event target.
 * @param valueOrSelector (String) - Selector or a value.
 * @param value (String) - Target value.
 * @returns (Boolean) - Returns true if an element was found.
 */
func Input(_ context.Context, args ...core.Value) (core.Value, error) {
	err := core.ValidateArgs(args, 2, 3)

	if err != nil {
		return values.None, err
	}

	// TYPE(el, "foobar")
	if len(args) == 2 {
		arg1 := args[0]

		err := core.ValidateType(arg1, core.HTMLElementType)

		if err != nil {
			return values.False, err
		}

		el, ok := arg1.(*dynamic.HTMLElement)

		if !ok {
			return values.False, core.Errors(core.ErrInvalidType, ErrNotDynamic)
		}

		err = el.Input(args[1])

		if err != nil {
			return values.False, err
		}

		return values.True, nil
	}

	arg1 := args[0]

	err = core.ValidateType(arg1, core.HTMLDocumentType)

	if err != nil {
		return values.False, err
	}

	arg2 := args[1]

	err = core.ValidateType(arg2, core.StringType)

	if err != nil {
		return values.False, err
	}

	doc, ok := arg1.(*dynamic.HTMLDocument)

	if !ok {
		return values.False, core.Errors(core.ErrInvalidType, ErrNotDynamic)
	}

	return doc.InputBySelector(arg2.(values.String), args[2])
}
