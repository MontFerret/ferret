package strings

import (
	"context"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/MontFerret/ferret/pkg/runtime/values/types"
)

// Concat concatenates one or more instances of Read, or an Array.
// @params src (String...|Array) - The source string / array.
// @returns String
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

// ConcatWithSeparator concatenates one or more instances of Read, or an Array with a given separator.
// @params separator (string) - The separator string.
// @params src (string...|array) - The source string / array.
// @returns string
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
