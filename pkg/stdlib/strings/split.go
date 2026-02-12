package strings

import (
	"strings"

	"github.com/MontFerret/ferret/pkg/runtime"
)

// SPLIT splits the given string value into a list of strings, using the separator.
// @param {String} str - The string to split.
// @param {String} separator - The separator.
// @param {Int} limit - Limit the number of split values in the result. If no limit is given, the number of splits returned is not bounded.
// @return {String[]} - arrayList of strings.
func Split(ctx runtime.Context, args ...runtime.Value) (runtime.Value, error) {
	err := runtime.ValidateArgs(args, 2, 3)

	if err != nil {
		return runtime.None, err
	}

	text := args[0].String()
	separator := args[1].String()
	limit := -1

	if len(args) > 2 {
		args2, ok := args[2].(runtime.Int)

		if ok {
			limit = int(args2)
		}
	}

	var strs []string

	if limit < 0 {
		strs = strings.Split(text, separator)
	} else {
		strs = strings.SplitN(text, separator, limit)
	}

	arr := ctx.Alloc().Array(len(strs))

	for _, str := range strs {
		_ = arr.Append(ctx, runtime.NewString(str))
	}

	return arr, nil
}
