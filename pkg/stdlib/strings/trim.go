package strings

import (
	"context"
	"strings"

	"github.com/MontFerret/ferret/pkg/runtime"
)

// TRIM returns the string value with whitespace stripped from the start and/or end.
// @param {String} str - The string.
// @param {String} chars - Overrides the characters that should be removed from the string. It defaults to \r\n \t.
// @return {String} - The string without chars on both sides.
func Trim(_ context.Context, args ...runtime.Value) (runtime.Value, error) {
	err := runtime.ValidateArgs(args, 1, 2)

	if err != nil {
		return runtime.EmptyString, err
	}

	text := args[0].String()

	if len(args) > 1 {
		return runtime.NewString(strings.Trim(text, args[1].String())), nil
	}

	return runtime.NewString(strings.TrimSpace(text)), nil
}

// LTRIM returns the string value with whitespace stripped from the start only.
// @param {String} str - The string.
// @param {String} chars - Overrides the characters that should be removed from the string. It defaults to \r\n \t.
// @return {String} - The string without chars at the left-hand side.
func LTrim(_ context.Context, args ...runtime.Value) (runtime.Value, error) {
	err := runtime.ValidateArgs(args, 1, 2)

	if err != nil {
		return runtime.EmptyString, err
	}

	text := args[0].String()
	chars := " "

	if len(args) > 1 {
		chars = args[1].String()
	}

	return runtime.NewString(strings.TrimLeft(text, chars)), nil
}

// RTRIM returns the string value with whitespace stripped from the end only.
// @param {String} str - The string.
// @param {String} chars - Overrides the characters that should be removed from the string. It defaults to \r\n \t.
// @return {String} - The string without chars at the right-hand side.
func RTrim(_ context.Context, args ...runtime.Value) (runtime.Value, error) {
	err := runtime.ValidateArgs(args, 1, 2)

	if err != nil {
		return runtime.EmptyString, err
	}

	text := args[0].String()
	chars := " "

	if len(args) > 1 {
		chars = args[1].String()
	}

	return runtime.NewString(strings.TrimRight(text, chars)), nil
}
