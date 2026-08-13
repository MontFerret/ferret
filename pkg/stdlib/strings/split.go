package strings

import (
	"context"
	"strings"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

// split splits the given string value into a list of strings, using the separator.
// @param str {String} The string to split.
// @param separator {String} The separator.
// @param limit {Int} Limit the number of split values in the result. If no limit is given, the number of splits returned is not bounded.
// @return {String[]} arrayList of strings.
func Split(ctx context.Context, args ...runtime.Value) (runtime.Value, error) {
	err := runtime.ValidateArgs(args, 2, 3)

	if err != nil {
		return runtime.None, err
	}

	if len(args) == 2 {
		return split2(ctx, args[0], args[1])
	}

	return split3(ctx, args[0], args[1], args[2])
}

// split splits the given string value into a list of strings, using the separator.
// @param str {String} The string to split.
// @param separator {String} The separator.
// @return {String[]} arrayList of strings.
func split2(ctx context.Context, arg1, arg2 runtime.Value) (runtime.Value, error) {
	return split(ctx, arg1, arg2, -1)
}

// split splits the given string value into a list of strings, using the separator.
// @param str {String} The string to split.
// @param separator {String} The separator.
// @param limit {Int} Limit the number of split values in the result. If no limit is given, the number of splits returned is not bounded.
// @return {String[]} arrayList of strings.
func split3(ctx context.Context, arg1, arg2, arg3 runtime.Value) (runtime.Value, error) {
	limit := runtime.CastOr[runtime.Int](arg3, runtime.Int(-1))
	return split(ctx, arg1, arg2, int(limit))
}

func split(ctx context.Context, arg1, arg2 runtime.Value, limit int) (runtime.Value, error) {
	text := arg1.String()
	separator := arg2.String()

	var strs []string

	if limit < 0 {
		strs = strings.Split(text, separator)
	} else {
		strs = strings.SplitN(text, separator, limit)
	}

	arr := runtime.NewArray(len(strs))

	for _, str := range strs {
		_ = arr.Append(ctx, runtime.NewString(str))
	}

	return arr, nil
}
