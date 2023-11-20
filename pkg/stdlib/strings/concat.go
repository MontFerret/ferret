package strings

import (
	"context"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
)

// CONCAT concatenates one or more instances of String, or an Array.
// @param {String, repeated | String[]} src - The source string / array.
// @return {String} - A string value.
func Concat(_ context.Context, args ...core.Value) (core.Value, error) {
	if err := core.ValidateArgs(args, 1, core.MaxArgs); err != nil {
		return values.EmptyString, err
	}

	argsCount := len(args)

	res := values.EmptyString

	if argsCount == 1 {
		argv, ok := args[0].(*values.Array)

		if ok {
			argv.ForEach(func(value core.Value, _ int) bool {
				res = res.Concat(value)

				return true
			})

			return res, nil
		}
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

	separator, ok := args[0].(values.String)

	if !ok {
		separator = values.NewString(separator.String())
	}

	res := values.EmptyString

	for idx, arg := range args[1:] {
		switch argv := arg.(type) {
		case *values.Array:
			argv.ForEach(func(value core.Value, idx int) bool {
				if value != values.None {
					if idx > 0 {
						res = res.Concat(separator)
					}

					res = res.Concat(value)
				}

				return true
			})
		default:
			if argv != values.None {
				if idx > 0 {
					res = res.Concat(separator)
				}

				res = res.Concat(argv)
			}

		}
	}

	return res, nil
}
