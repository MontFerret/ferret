package strings

import (
	"context"
	"strings"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

// TRIM returns the string value with whitespace stripped from the start and/or end.
// @param {String} str - The string.
// @param {String} chars - Overrides the characters that should be removed from the string. It defaults to \r\n \t.
// @return {String} - The string without chars on both sides.
func Trim(ctx context.Context, args ...runtime.Value) (runtime.Value, error) {
	err := runtime.ValidateArgs(args, 1, 2)

	if err != nil {
		return runtime.EmptyString, err
	}

	if len(args) == 1 {
		return trim1(ctx, args[0])
	}

	return trim2(ctx, args[0], args[1])
}

func trim1(_ context.Context, arg1 runtime.Value) (runtime.Value, error) {
	return runtime.NewString(strings.TrimSpace(arg1.String())), nil
}

func trim2(_ context.Context, arg1, arg2 runtime.Value) (runtime.Value, error) {
	return runtime.NewString(strings.Trim(arg1.String(), arg2.String())), nil
}

// LTRIM returns the string value with whitespace stripped from the start only.
// @param {String} str - The string.
// @param {String} chars - Overrides the characters that should be removed from the string. It defaults to \r\n \t.
// @return {String} - The string without chars at the left-hand side.
func LTrim(ctx context.Context, args ...runtime.Value) (runtime.Value, error) {
	err := runtime.ValidateArgs(args, 1, 2)

	if err != nil {
		return runtime.EmptyString, err
	}

	if len(args) == 1 {
		return lTrim1(ctx, args[0])
	}

	return lTrim2(ctx, args[0], args[1])
}

func lTrim1(_ context.Context, arg1 runtime.Value) (runtime.Value, error) {
	return runtime.NewString(strings.TrimLeft(arg1.String(), " ")), nil
}

func lTrim2(_ context.Context, arg1, arg2 runtime.Value) (runtime.Value, error) {
	return runtime.NewString(strings.TrimLeft(arg1.String(), arg2.String())), nil
}

// RTRIM returns the string value with whitespace stripped from the end only.
// @param {String} str - The string.
// @param {String} chars - Overrides the characters that should be removed from the string. It defaults to \r\n \t.
// @return {String} - The string without chars at the right-hand side.
func RTrim(ctx context.Context, args ...runtime.Value) (runtime.Value, error) {
	err := runtime.ValidateArgs(args, 1, 2)

	if err != nil {
		return runtime.EmptyString, err
	}

	if len(args) == 1 {
		return rTrim1(ctx, args[0])
	}

	return rTrim2(ctx, args[0], args[1])
}

func rTrim1(_ context.Context, arg1 runtime.Value) (runtime.Value, error) {
	return runtime.NewString(strings.TrimRight(arg1.String(), " ")), nil
}

func rTrim2(_ context.Context, arg1, arg2 runtime.Value) (runtime.Value, error) {
	return runtime.NewString(strings.TrimRight(arg1.String(), arg2.String())), nil
}
