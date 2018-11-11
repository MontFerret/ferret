package html

import (
	"context"
	"github.com/MontFerret/ferret/pkg/html/dynamic"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
)

// Select selects a value from an underlying select element.
// @param source (Document | Element) - Event target.
// @param valueOrSelector (String | Array<String>) - Selector or a an array of strings as a value.
// @param value (Array<String) - Target value. Optional.
// @returns (Array<String>) - Returns an array of selected values.
func Select(_ context.Context, args ...core.Value) (core.Value, error) {
	err := core.ValidateArgs(args, 2, 4)

	if err != nil {
		return values.None, err
	}

	arg1 := args[0]
	err = core.ValidateType(arg1, core.HTMLDocumentType, core.HTMLElementType)

	if err != nil {
		return values.False, err
	}

	switch args[0].(type) {
	case *dynamic.HTMLDocument:
		doc, ok := arg1.(*dynamic.HTMLDocument)

		if !ok {
			return values.False, core.Errors(core.ErrInvalidType, ErrNotDynamic)
		}

		// selector
		arg2 := args[1]
		err = core.ValidateType(arg2, core.StringType)

		if err != nil {
			return values.False, err
		}

		arg3 := args[2]
		err = core.ValidateType(arg3, core.ArrayType)

		if err != nil {
			return values.False, err
		}

		return doc.SelectBySelector(arg2.(values.String), arg3.(*values.Array))
	case *dynamic.HTMLElement:
		el, ok := arg1.(*dynamic.HTMLElement)

		if !ok {
			return values.False, core.Errors(core.ErrInvalidType, ErrNotDynamic)
		}

		arg2 := args[1]

		err = core.ValidateType(arg2, core.ArrayType)

		if err != nil {
			return values.False, err
		}

		return el.Select(arg2.(*values.Array))
	default:
		return values.False, core.Errors(core.ErrInvalidArgument)
	}
}
