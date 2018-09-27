package html

import (
	"context"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/MontFerret/ferret/pkg/stdlib/html/driver"
	"github.com/MontFerret/ferret/pkg/stdlib/html/driver/static"
)

func Document(ctx context.Context, inputs ...core.Value) (core.Value, error) {
	url, dynamic, err := documentArgs(inputs)

	if err != nil {
		return values.None, err
	}

	var drv driver.Driver

	if !dynamic {
		drv, err = driver.FromContext(ctx, driver.Static)
	} else {
		drv, err = driver.FromContext(ctx, driver.Dynamic)
	}

	if err != nil {
		return values.None, err
	}

	return drv.GetDocument(ctx, url.String())
}

func DocumentParse(ctx context.Context, inputs ...core.Value) (core.Value, error) {
	arg1 := values.EmptyString

	if len(inputs) == 0 {
		return arg1, core.Error(core.ErrMissedArgument, "document string")
	}

	a1 := inputs[0]

	if a1.Type() != core.StringType {
		return arg1, core.Error(core.TypeError(a1.Type(), core.StringType), "arg 1")
	}

	drv, err := driver.FromContext(ctx, driver.Static)

	if err != nil {
		return values.None, err
	}

	return drv.(*static.Driver).ParseDocument(ctx, arg1.String())
}

func documentArgs(inputs []core.Value) (values.String, values.Boolean, error) {
	arg1 := values.EmptyString
	arg2 := values.False

	if len(inputs) == 0 {
		return arg1, arg2, core.Error(core.ErrMissedArgument, "element and useJs")
	}

	a1 := inputs[0]

	if a1.Type() != core.StringType {
		return arg1, arg2, core.Error(core.TypeError(a1.Type(), core.StringType), "arg 1")
	}

	arg1 = a1.(values.String)

	if len(inputs) == 2 {
		a2 := inputs[1]

		if a2.Type() != core.BooleanType {
			return arg1, arg2, core.Error(core.TypeError(a2.Type(), core.BooleanType), "arg 2")
		}

		arg2 = a2.(values.Boolean)
	}

	return arg1, arg2, nil
}
