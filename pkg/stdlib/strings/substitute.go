package strings

import (
	"context"
	"strings"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

// SUBSTITUTE replaces search values in the string value.
// @param {String} str - The string to modify
// @param {String} search - The string representing a search pattern
// @param {String} replace - The string representing a replace value
// @param {Int} limit - The cap the number of replacements to this value.
// @return {String} - Returns a string with replace substring.
func Substitute(ctx context.Context, args ...runtime.Value) (runtime.Value, error) {
	err := runtime.ValidateArgs(args, 2, 4)

	if err != nil {
		return runtime.EmptyString, err
	}

	switch len(args) {
	case 2:
		return substitute2(ctx, args[0], args[1])
	case 3:
		return substitute3(ctx, args[0], args[1], args[2])
	default:
		return substitute4(ctx, args[0], args[1], args[2], args[3])
	}
}

func substitute2(ctx context.Context, arg1, arg2 runtime.Value) (runtime.Value, error) {
	return substitute(ctx, arg1, arg2, runtime.EmptyString, -1)
}

func substitute3(ctx context.Context, arg1, arg2, arg3 runtime.Value) (runtime.Value, error) {
	return substitute(ctx, arg1, arg2, arg3, -1)
}

func substitute4(ctx context.Context, arg1, arg2, arg3, arg4 runtime.Value) (runtime.Value, error) {
	limit := runtime.CastOr[runtime.Int](arg4, runtime.Int(-1))
	return substitute(ctx, arg1, arg2, arg3, int(limit))
}

func substitute(_ context.Context, arg1, arg2, arg3 runtime.Value, limit int) (runtime.Value, error) {
	text := arg1.String()
	search := arg2.String()
	replace := arg3.String()

	out := strings.Replace(text, search, replace, limit)

	return runtime.NewString(out), nil
}
