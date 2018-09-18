package types

import (
	"context"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
)

func isTypeof(value core.Value, ctype core.Type) core.Value {
	return values.NewBoolean(value.Type() == ctype)
}

func IsNone(_ context.Context, inputs ...core.Value) (core.Value, error) {
	if len(inputs) == 0 {
		return values.None, core.ErrMissedArgument
	}

	return isTypeof(inputs[0], core.NoneType), nil
}

func IsBool(_ context.Context, inputs ...core.Value) (core.Value, error) {
	if len(inputs) == 0 {
		return values.None, core.ErrMissedArgument
	}

	return isTypeof(inputs[0], core.BooleanType), nil
}

func IsInt(_ context.Context, inputs ...core.Value) (core.Value, error) {
	if len(inputs) == 0 {
		return values.None, core.ErrMissedArgument
	}

	return isTypeof(inputs[0], core.IntType), nil
}

func IsFloat(_ context.Context, inputs ...core.Value) (core.Value, error) {
	if len(inputs) == 0 {
		return values.None, core.ErrMissedArgument
	}

	return isTypeof(inputs[0], core.FloatType), nil
}

func IsString(_ context.Context, inputs ...core.Value) (core.Value, error) {
	if len(inputs) == 0 {
		return values.None, core.ErrMissedArgument
	}

	return isTypeof(inputs[0], core.StringType), nil
}

func IsDateTime(_ context.Context, inputs ...core.Value) (core.Value, error) {
	if len(inputs) == 0 {
		return values.None, core.ErrMissedArgument
	}

	return isTypeof(inputs[0], core.DateTimeType), nil
}

func IsArray(_ context.Context, inputs ...core.Value) (core.Value, error) {
	if len(inputs) == 0 {
		return values.None, core.ErrMissedArgument
	}

	return isTypeof(inputs[0], core.ArrayType), nil
}

func IsObject(_ context.Context, inputs ...core.Value) (core.Value, error) {
	if len(inputs) == 0 {
		return values.None, core.ErrMissedArgument
	}

	return isTypeof(inputs[0], core.ObjectType), nil
}

func IsHtmlElement(_ context.Context, inputs ...core.Value) (core.Value, error) {
	if len(inputs) == 0 {
		return values.None, core.ErrMissedArgument
	}

	return isTypeof(inputs[0], core.HtmlElementType), nil
}

func IsHtmlDocument(_ context.Context, inputs ...core.Value) (core.Value, error) {
	if len(inputs) == 0 {
		return values.None, core.ErrMissedArgument
	}

	return isTypeof(inputs[0], core.HtmlDocumentType), nil
}

func IsBinary(_ context.Context, inputs ...core.Value) (core.Value, error) {
	if len(inputs) == 0 {
		return values.None, core.ErrMissedArgument
	}

	return isTypeof(inputs[0], core.BinaryType), nil
}

func TypeName(_ context.Context, inputs ...core.Value) (core.Value, error) {
	if len(inputs) == 0 {
		return values.None, core.ErrMissedArgument
	}

	return values.NewString(inputs[0].Type().String()), nil
}
