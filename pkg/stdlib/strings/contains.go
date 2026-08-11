package strings

import (
	"context"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

// CONTAINS returns a value indicating whether a specified substring occurs within a string.
// @param str {String} The source string.
// @param search {String} The string to seek.
// @param returnIndex {Boolean} Values which indicates whether to return the character position of the match is returned instead of a boolean.
// @return {Boolean | Int} A value indicating whether a specified substring occurs within a string.
func Contains(ctx context.Context, args ...runtime.Value) (runtime.Value, error) {
	if err := runtime.ValidateArgs(args, 2, 3); err != nil {
		return runtime.False, err
	}

	if len(args) == 2 {
		return contains2(ctx, args[0], args[1])
	}

	return contains3(ctx, args[0], args[1], args[2])
}

// CONTAINS returns a value indicating whether a specified substring occurs within a string.
// @param str {String} The source string.
// @param search {String} The string to seek.
// @return {Boolean | Int} A value indicating whether a specified substring occurs within a string.
func contains2(ctx context.Context, arg1, arg2 runtime.Value) (runtime.Value, error) {
	return contains3(ctx, arg1, arg2, runtime.False)
}

// CONTAINS returns a value indicating whether a specified substring occurs within a string.
// @param str {String} The source string.
// @param search {String} The string to seek.
// @param returnIndex {Boolean} Values which indicates whether to return the character position of the match is returned instead of a boolean.
// @return {Boolean | Int} A value indicating whether a specified substring occurs within a string.
func contains3(_ context.Context, arg1, arg2, arg3 runtime.Value) (runtime.Value, error) {
	text := runtime.CastOr[runtime.String](arg1, runtime.EmptyString)
	search := runtime.CastOr[runtime.String](arg2, runtime.EmptyString)
	returnIndex := runtime.CastOr[runtime.Boolean](arg3, runtime.False)

	if returnIndex == runtime.True {
		return text.IndexOf(search), nil
	}

	return text.Contains(search), nil
}
