package strings

import (
	"context"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

// concat concatenates values using their string representations. With one List argument, concatenates its elements. None contributes no text; nested lists are not flattened.
// @param src {Any, repeated} Values to concatenate, or a single List of values.
// @return {String} A string value.
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
