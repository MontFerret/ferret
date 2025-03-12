package strings

import (
	"context"
	"github.com/MontFerret/ferret/pkg/runtime/internal"

	"github.com/MontFerret/ferret/pkg/runtime/core"
)

// CONCAT concatenates one or more instances of String, or an Array.
// @param {String, repeated | String[]} src - The source string / array.
// @return {String} - A string value.
func Concat(_ context.Context, args ...core.Value) (core.Value, error) {
	if err := core.ValidateArgs(args, 1, core.MaxArgs); err != nil {
		return core.EmptyString, err
	}

	argsCount := len(args)

	res := core.EmptyString

	if argsCount == 1 {
		argv, ok := args[0].(*internal.Array)

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
		return core.EmptyString, err
	}

	separator := args[0]

	separator, ok := args[0].(core.String)

	if !ok {
		separator = core.NewString(separator.String())
	}

	res := core.EmptyString

	for idx, arg := range args[1:] {
		switch argv := arg.(type) {
		case *internal.Array:
			argv.ForEach(func(value core.Value, idx int) bool {
				if value != core.None {
					if idx > 0 {
						res = res.Concat(separator)
					}

					res = res.Concat(value)
				}

				return true
			})
		default:
			if argv != core.None {
				if idx > 0 {
					res = res.Concat(separator)
				}

				res = res.Concat(argv)
			}

		}
	}

	return res, nil
}
