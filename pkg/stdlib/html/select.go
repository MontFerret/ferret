package html

import (
	"context"

	"github.com/MontFerret/ferret/pkg/drivers"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/MontFerret/ferret/pkg/runtime/values/types"
)

// Select selects a value from an underlying select element.
// @param source (Open | GetElement) - Event target.
// @param valueOrSelector (String | Array<String>) - Selector or a an array of strings as a value.
// @param value (Array<String) - Target value. Optional.
// @returns (Array<String>) - Returns an array of selected values.
func Select(ctx context.Context, args ...core.Value) (core.Value, error) {
	err := core.ValidateArgs(args, 2, 4)

	if err != nil {
		return values.None, err
	}

	arg1 := args[0]
	err = core.ValidateType(arg1, drivers.HTMLDocumentType, drivers.HTMLElementType)

	if err != nil {
		return values.False, err
	}

	if arg1.Type() == drivers.HTMLDocumentType {
		doc := arg1.(drivers.HTMLDocument)

		// selector
		arg2 := args[1]
		err = core.ValidateType(arg2, types.String)

		if err != nil {
			return values.False, err
		}

		arg3 := args[2]
		err = core.ValidateType(arg3, types.Array)

		if err != nil {
			return values.False, err
		}

		return doc.SelectBySelector(ctx, arg2.(values.String), arg3.(*values.Array))
	}

	el := arg1.(drivers.HTMLElement)
	arg2 := args[1]

	err = core.ValidateType(arg2, types.Array)

	if err != nil {
		return values.False, err
	}

	return el.Select(ctx, arg2.(*values.Array))
}
