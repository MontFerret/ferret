package strings

import (
	"context"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	"strings"
)

func arg1(inputs []core.Value) (core.Value, error) {
	if len(inputs) == 0 {
		return values.None, core.Error(core.ErrMissedArgument, "value")
	}

	return inputs[0], nil
}

func Contains(_ context.Context, inputs ...core.Value) (core.Value, error) {
	if len(inputs) < 2 {
		return values.None, core.ErrMissedArgument
	}

	var text values.String
	var search values.String
	returnIndex := values.False

	err := core.ValidateType(inputs[0], core.StringType)

	if err != nil {
		return values.None, err
	}

	err = core.ValidateType(inputs[1], core.StringType)

	if err != nil {
		return values.None, err
	}

	text = inputs[0].(values.String)
	search = inputs[1].(values.String)

	if len(inputs) > 2 {
		err = core.ValidateType(inputs[2], core.BooleanType)

		if err != nil {
			return values.None, err
		}

		returnIndex = inputs[2].(values.Boolean)
	}

	if returnIndex == values.True {
		return text.IndexOf(search), nil
	}

	return text.Contains(search), nil
}

func Concat(_ context.Context, inputs ...core.Value) (core.Value, error) {
	if len(inputs) == 0 {
		return values.EmptyString, nil
	}

	res := values.EmptyString

	for _, str := range inputs {
		res = res.Concat(str)
	}

	return res, nil
}

func ConcatWithSeparator(_ context.Context, inputs ...core.Value) (core.Value, error) {
	if len(inputs) == 0 || len(inputs) == 1 {
		return values.EmptyString, nil
	}

	separator := inputs[0]

	if separator.Type() != core.StringType {
		separator = values.NewString(separator.String())
	}

	res := values.EmptyString

	for _, str := range inputs[1:] {
		res = res.Concat(separator).Concat(str)
	}

	return res, nil
}

func FindFirst(_ context.Context, _ ...core.Value) (core.Value, error) {
	return values.None, core.Error(core.ErrNotImplemented, "FIND_FIRST")
}

func FindLast(_ context.Context, _ ...core.Value) (core.Value, error) {
	return values.None, core.Error(core.ErrNotImplemented, "FIND_LAST")
}

func JsonParse(_ context.Context, _ ...core.Value) (core.Value, error) {
	return values.None, core.Error(core.ErrNotImplemented, "JSON_PARSE")
}

func JsonStringify(_ context.Context, _ ...core.Value) (core.Value, error) {
	return values.None, core.Error(core.ErrNotImplemented, "JSON_STRINGIFY")
}

func Left(_ context.Context, _ ...core.Value) (core.Value, error) {
	return values.None, core.Error(core.ErrNotImplemented, "LEFT")
}

func Like(_ context.Context, _ ...core.Value) (core.Value, error) {
	return values.None, core.Error(core.ErrNotImplemented, "LIKE")
}

func Lower(_ context.Context, _ ...core.Value) (core.Value, error) {
	return values.None, core.Error(core.ErrNotImplemented, "LOWER")
}

func LTrim(_ context.Context, _ ...core.Value) (core.Value, error) {
	return values.None, core.Error(core.ErrNotImplemented, "LTRIM")
}

func Right(_ context.Context, _ ...core.Value) (core.Value, error) {
	return values.None, core.Error(core.ErrNotImplemented, "RIGHT")
}

func RTrim(_ context.Context, _ ...core.Value) (core.Value, error) {
	return values.None, core.Error(core.ErrNotImplemented, "RTRIM")
}

func Split(_ context.Context, _ ...core.Value) (core.Value, error) {
	return values.None, core.Error(core.ErrNotImplemented, "SPLIT")
}

func Substitute(_ context.Context, _ ...core.Value) (core.Value, error) {
	return values.None, core.Error(core.ErrNotImplemented, "SUBSTITUTE")
}

func Substring(_ context.Context, _ ...core.Value) (core.Value, error) {
	return values.None, core.Error(core.ErrNotImplemented, "SUBSTRING")
}

func Trim(_ context.Context, inputs ...core.Value) (core.Value, error) {
	val, err := arg1(inputs)

	if err != nil {
		return values.None, err
	}

	return values.NewString(strings.TrimSpace(val.String())), nil
}

func Upper(_ context.Context, _ ...core.Value) (core.Value, error) {
	return values.None, core.Error(core.ErrNotImplemented, "UPPER")
}
