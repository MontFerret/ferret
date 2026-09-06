package strings

import (
	"context"
	"strings"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

// replace replaces literal, non-overlapping occurrences. An empty search matches before each rune and at the end; replacement text is literal.
// @param str {String} The string to modify
// @param search {String} The string representing a search pattern
// @param replace {String} The string representing a replace value
// @param limit {Int} Non-negative maximum replacements. Zero leaves text unchanged; omitted means unlimited.
// @return {String} Returns a string with replace substring.
func Replace(ctx context.Context, args ...runtime.Value) (runtime.Value, error) {
	err := runtime.ValidateArgs(args, 3, 4)
	if err != nil {
		return runtime.EmptyString, err
	}

	switch len(args) {
	case 3:
		return replace3(ctx, args[0], args[1], args[2])
	default:
		return replace4(ctx, args[0], args[1], args[2], args[3])
	}
}

// replace replaces literal, non-overlapping occurrences. An empty search matches before each rune and at the end; replacement text is literal.
// @param str {String} The string to modify
// @param search {String} The string representing a search pattern
// @param replace {String} The string representing a replace value
// @return {String} Returns a string with replace substring.
func replace3(ctx context.Context, arg1, arg2, arg3 runtime.Value) (runtime.Value, error) {
	return replace(ctx, arg1, arg2, arg3, -1)
}

// replace replaces literal, non-overlapping occurrences. An empty search matches before each rune and at the end; replacement text is literal.
// @param str {String} The string to modify
// @param search {String} The string representing a search pattern
// @param replace {String} The string representing a replace value
// @param limit {Int} Non-negative maximum replacements. Zero leaves text unchanged; omitted means unlimited.
// @return {String} Returns a string with replace substring.
func replace4(ctx context.Context, arg1, arg2, arg3, arg4 runtime.Value) (runtime.Value, error) {
	limit, err := nonNegativeLimit(arg4, 3)
	if err != nil {
		return runtime.None, err
	}

	return replace(ctx, arg1, arg2, arg3, limit)
}

func replace(_ context.Context, arg1, arg2, arg3 runtime.Value, limit int) (runtime.Value, error) {
	text, search, replacement, err := runtime.CastArgs3[runtime.String, runtime.String, runtime.String](arg1, arg2, arg3)
	if err != nil {
		return runtime.None, err
	}

	out := strings.Replace(string(text), string(search), string(replacement), limit)

	return runtime.NewString(out), nil
}
