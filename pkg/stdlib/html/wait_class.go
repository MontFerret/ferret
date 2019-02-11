package html

import (
	"context"

	"github.com/MontFerret/ferret/pkg/drivers"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/MontFerret/ferret/pkg/runtime/values/types"
)

// WaitClass waits for a class to appear on a given element.
// Stops the execution until the navigation ends or operation times out.
// @param docOrEl (HTMLDocument|HTMLElement) - Target document or element.
// @param selectorOrClass (String) - If document is passed, this param must represent an element selector.
// Otherwise target class.
// @param classOrTimeout (String|Int, optional) - If document is passed, this param must represent target class name.
// Otherwise timeout.
// @param timeout (Int, optional) - If document is passed, this param must represent timeout.
// Otherwise not passed.
func WaitClass(_ context.Context, args ...core.Value) (core.Value, error) {
	err := core.ValidateArgs(args, 2, 4)

	if err != nil {
		return values.None, err
	}

	// document or element
	arg1 := args[0]
	err = core.ValidateType(arg1, drivers.HTMLDocumentType, drivers.HTMLElementType)

	if err != nil {
		return values.None, err
	}

	// selector or class
	err = core.ValidateType(args[1], types.String)

	if err != nil {
		return values.None, err
	}

	timeout := values.NewInt(defaultTimeout)

	// if a document is passed
	if arg1.Type() == drivers.HTMLDocumentType {
		// revalidate args with more accurate amount
		err := core.ValidateArgs(args, 3, 4)

		if err != nil {
			return values.None, err
		}

		// class
		err = core.ValidateType(args[2], types.String)

		if err != nil {
			return values.None, err
		}

		doc := arg1.(drivers.HTMLDocument)
		selector := args[1].(values.String)
		class := args[2].(values.String)

		if len(args) == 4 {
			err = core.ValidateType(args[3], types.Int)

			if err != nil {
				return values.None, err
			}

			timeout = args[3].(values.Int)
		}

		return values.None, doc.WaitForClassBySelector(selector, class, timeout)
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

	return values.None, el.WaitForClass(class, timeout)
}
