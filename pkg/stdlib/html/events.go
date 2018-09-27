package html

import (
	"context"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/MontFerret/ferret/pkg/stdlib/html/driver/dynamic"
)

func WaitElement(_ context.Context, args ...core.Value) (core.Value, error) {
	err := core.ValidateArgs(args, 2, 3)

	if err != nil {
		return values.None, err
	}

	arg := args[0]
	selector := args[1].String()
	timeout := values.NewInt(5000)

	if len(args) > 2 {
		if args[2].Type() == core.IntType {
			timeout = args[2].(values.Int)
		}
	}

	err = core.ValidateType(arg, core.HtmlDocumentType)

	if err != nil {
		return values.None, err
	}

	doc, ok := arg.(*dynamic.HtmlDocument)

	if !ok {
		return values.False, core.Error(core.ErrInvalidType, "expected dynamic document")
	}

	return values.None, doc.WaitForSelector(values.NewString(selector), timeout)
}

func WaitNavigation(_ context.Context, args ...core.Value) (core.Value, error) {
	err := core.ValidateArgs(args, 1, 2)

	if err != nil {
		return values.None, err
	}

	err = core.ValidateType(args[0], core.HtmlDocumentType)

	if err != nil {
		return values.None, err
	}

	doc, ok := args[0].(*dynamic.HtmlDocument)

	if !ok {
		return values.None, core.Error(core.ErrInvalidType, "expected dynamic document")
	}

	timeout := values.NewInt(5000)

	if len(args) > 1 {
		err = core.ValidateType(args[1], core.IntType)

		if err != nil {
			return values.None, err
		}

		timeout = args[1].(values.Int)
	}

	return values.None, doc.WaitForNavigation(timeout)
}
