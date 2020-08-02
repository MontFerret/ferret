package strings

import (
	"context"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/MontFerret/ferret/pkg/runtime/values/types"
)

// CONCAT concatenates one or more instances of String, or an Array.
// @param {String, repeated | String[]} src - The source string / array.
// @return {String} - A string value.
func Concat(_ context.Context, args ...core.Value) (core.Value, error) {
	err := core.ValidateArgs(args, 1, core.MaxArgs)

	if err != nil {
		return values.EmptyString, err
	}

	argsCount := len(args)

	res := values.EmptyString

	if argsCount == 1 && args[0].Type() == types.Array {
		arr := args[0].(*values.Array)

		arr.ForEach(func(value core.Value, _ int) bool {
			res = res.Concat(value)

			return true
		})

		return res, nil
	}

	for _, str := range args {
		res = res.Concat(str)
	}

	return res, nil
}

// CONCAT_SEPARATOR concatenates one or more instances of String, or an Array with a given separator.
// @param {String} separator - The separator string.
// @param {String, repeated | String[]} src - The source string / array.
// @return {String} - Concatenated string.
func ConcatWithSeparator(_ context.Context, args ...core.Value) (core.Value, error) {
	err := core.ValidateArgs(args, 2, core.MaxArgs)

	if err != nil {
		return values.EmptyString, err
	}

	separator := args[0]

	if separator.Type() != types.String {
		separator = values.NewString(separator.String())
	}

	res := values.EmptyString

	for idx, arg := range args[1:] {
		if arg.Type() != types.Array {
			if arg.Type() != types.None {
				if idx > 0 {
					res = res.Concat(separator)
				}

				res = res.Concat(arg)
			}
		} else {
			arr := arg.(*values.Array)

			arr.ForEach(func(value core.Value, idx int) bool {
				if value.Type() != types.None {
					if idx > 0 {
						res = res.Concat(separator)
					}

					res = res.Concat(value)
				}

				return true
			})
		}
	}

	return res, nil
}
