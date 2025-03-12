package strings

import (
	"context"
	"github.com/MontFerret/ferret/pkg/runtime/internal"
	"strings"

	"github.com/MontFerret/ferret/pkg/runtime/core"
)

// SPLIT splits the given string value into a list of strings, using the separator.
// @param {String} str - The string to split.
// @param {String} separator - The separator.
// @param {Int} limit - Limit the number of split values in the result. If no limit is given, the number of splits returned is not bounded.
// @return {String[]} - Array of strings.
func Split(_ context.Context, args ...core.Value) (core.Value, error) {
	err := core.ValidateArgs(args, 2, 3)

	if err != nil {
		return core.None, err
	}

	text := args[0].String()
	separator := args[1].String()
	limit := -1

	if len(args) > 2 {
		args2, ok := args[2].(core.Int)

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

	arr := internal.NewArray(len(strs))

	for _, str := range strs {
		arr.Push(core.NewString(str))
	}

	return arr, nil
}
