package html

import (
	"context"

	"github.com/MontFerret/ferret/pkg/drivers"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/MontFerret/ferret/pkg/runtime/values/types"
)

func WaitAttribute(ctx context.Context, args ...core.Value) (core.Value, error) {
	return waitAttributeWhen(ctx, args, drivers.WaitEventPresence)
}

func WaitNoAttribute(ctx context.Context, args ...core.Value) (core.Value, error) {
	return waitAttributeWhen(ctx, args, drivers.WaitEventAbsence)
}

func waitAttributeWhen(ctx context.Context, args []core.Value, when drivers.WaitEvent) (core.Value, error) {
	err := core.ValidateArgs(args, 3, 5)

	if err != nil {
		return values.None, err
	}

	// document or element
	arg1 := args[0]
	err = core.ValidateType(arg1, drivers.HTMLDocumentType, drivers.HTMLElementType)

	if err != nil {
		return values.None, err
	}

	// selector or attr name
	err = core.ValidateType(args[1], types.String)

	if err != nil {
		return values.None, err
	}

	timeout := values.NewInt(defaultTimeout)

	// if a document is passed
	// WAIT_ATTR(doc, selector, attrName, attrValue, timeout)
	if arg1.Type() == drivers.HTMLDocumentType {
		// revalidate args with more accurate amount
		err := core.ValidateArgs(args, 4, 5)

		if err != nil {
			return values.None, err
		}

		// attr name
		err = core.ValidateType(args[2], types.String)

		if err != nil {
			return values.None, err
		}

		doc := arg1.(drivers.HTMLDocument)
		selector := args[1].(values.String)
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

		return values.None, doc.WaitForAttributeBySelector(ctx, selector, name, value, when)
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

	return values.None, el.WaitForAttribute(ctx, name, value, when)
}
