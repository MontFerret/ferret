package html

import (
	"context"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
)

func Element(_ context.Context, inputs ...core.Value) (core.Value, error) {
	el, selector, err := elementArgs(inputs)

	if err != nil {
		return values.None, err
	}

	return el.QuerySelector(selector), nil
}

func Elements(_ context.Context, inputs ...core.Value) (core.Value, error) {
	el, selector, err := elementArgs(inputs)

	if err != nil {
		return values.None, err
	}

	return el.QuerySelectorAll(selector), nil
}

func elementArgs(inputs []core.Value) (values.HtmlNode, values.String, error) {
	if len(inputs) == 0 {
		return nil, values.EmptyString, core.Error(core.ErrMissedArgument, "element and arg2")
	}

	if len(inputs) == 1 {
		return nil, values.EmptyString, core.Error(core.ErrMissedArgument, "arg2")
	}

	arg1 := inputs[0]
	arg2 := inputs[1]

	if arg1.Type() != core.HtmlDocumentType &&
		arg1.Type() != core.HtmlElementType {
		return nil, values.EmptyString, core.TypeError(arg1.Type(), core.HtmlDocumentType, core.HtmlElementType)
	}

	if arg2.Type() != core.StringType {
		return nil, values.EmptyString, core.TypeError(arg2.Type(), core.StringType)
	}

	return arg1.(values.HtmlNode), arg2.(values.String), nil
}
