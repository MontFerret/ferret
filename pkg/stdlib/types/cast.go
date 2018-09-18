package types

import (
	"context"
	"github.com/MontFerret/ferret/pkg/runtime/collections"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
)

func ToBool(_ context.Context, inputs ...core.Value) (core.Value, error) {
	if len(inputs) == 0 {
		return values.None, core.ErrMissedArgument
	}

	return values.ParseBoolean(inputs[0].Unwrap())
}

func ToInt(_ context.Context, inputs ...core.Value) (core.Value, error) {
	if len(inputs) == 0 {
		return values.None, core.ErrMissedArgument
	}

	return values.ParseInt(inputs[0].Unwrap())
}

func ToFloat(_ context.Context, inputs ...core.Value) (core.Value, error) {
	if len(inputs) == 0 {
		return values.None, core.ErrMissedArgument
	}

	return values.ParseFloat(inputs[0].Unwrap())
}

func ToString(_ context.Context, inputs ...core.Value) (core.Value, error) {
	if len(inputs) == 0 {
		return values.None, core.ErrMissedArgument
	}

	return values.NewString(inputs[0].String()), nil
}

func ToDateTime(_ context.Context, inputs ...core.Value) (core.Value, error) {
	if len(inputs) == 0 {
		return values.None, core.ErrMissedArgument
	}

	return values.ParseDateTime(inputs[0].String())
}

func ToArray(_ context.Context, inputs ...core.Value) (core.Value, error) {
	if len(inputs) == 0 {
		return values.None, core.ErrMissedArgument
	}

	value := inputs[0]

	switch value.Type() {
	case core.BooleanType,
		core.IntType,
		core.FloatType,
		core.StringType,
		core.DateTimeType,
		core.HtmlElementType,
		core.HtmlDocumentType:
		return values.NewArrayWith(value), nil
	case core.ArrayType:
		return value, nil
	case core.ObjectType:
		return collections.ToArray(collections.NewObjectIterator(value.(*values.Object)))
	default:
		return values.NewArray(0), nil
	}
}
