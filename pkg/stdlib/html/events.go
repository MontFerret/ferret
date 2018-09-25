package html

import (
	"context"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/MontFerret/ferret/pkg/stdlib/html/driver/browser"
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

	doc, ok := arg.(*browser.HtmlDocument)

	if !ok {
		return values.False, core.Error(core.ErrInvalidType, "expected dynamic document")
	}

	return values.None, doc.WaitForSelector(values.NewString(selector), timeout)
}
