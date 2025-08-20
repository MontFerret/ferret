package strings

import (
	"context"

	"github.com/MontFerret/ferret/pkg/runtime"
)

// CONCAT concatenates one or more instances of String, or an arrayList.
// @param {String, repeated | String[]} src - The source string / array.
// @return {String} - A string value.
func Concat(ctx context.Context, args ...runtime.Value) (runtime.Value, error) {
	if err := runtime.ValidateArgs(args, 1, runtime.MaxArgs); err != nil {
		return runtime.EmptyString, err
	}

	argsCount := len(args)

	res := runtime.EmptyString

	if argsCount == 1 {
		argv, ok := args[0].(runtime.List)

		if ok {
			err := argv.ForEach(ctx, func(c context.Context, value runtime.Value, _ runtime.Int) (runtime.Boolean, error) {
				res = res.Concat(value)

				return true, nil
			})

			if err != nil {
				return runtime.None, err
			}

			return res, nil
		}
	}

	for _, str := range args {
		res = res.Concat(str)
	}

	return res, nil
}

// CONCAT_SEPARATOR concatenates one or more instances of String, or an arrayList with a given separator.
// @param {String} separator - The separator string.
// @param {String, repeated | String[]} src - The source string / array.
// @return {String} - Concatenated string.
func ConcatWithSeparator(ctx context.Context, args ...runtime.Value) (runtime.Value, error) {
	err := runtime.ValidateArgs(args, 2, runtime.MaxArgs)

	if err != nil {
		return runtime.EmptyString, err
	}

	separator, ok := args[0].(runtime.String)

	if !ok {
		separator = runtime.NewString(separator.String())
	}

	res := runtime.EmptyString

	for idx, arg := range args[1:] {
		switch argv := arg.(type) {
		case runtime.List:
			err = argv.ForEach(ctx, func(c context.Context, value runtime.Value, idx runtime.Int) (runtime.Boolean, error) {
				if value != runtime.None {
					if idx > 0 {
						res = res.Concat(separator)
					}

					res = res.Concat(value)
				}

				return true, nil
			})

			if err != nil {
				return runtime.None, err
			}
		default:
			if argv != runtime.None {
				if idx > 0 {
					res = res.Concat(separator)
				}

				res = res.Concat(argv)
			}

		}
	}

	return res, nil
}
