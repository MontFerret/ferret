package strings

import (
	"context"
	"strings"

	"github.com/MontFerret/ferret/pkg/runtime/core"
)

// TRIM returns the string value with whitespace stripped from the start and/or end.
// @param {String} str - The string.
// @param {String} chars - Overrides the characters that should be removed from the string. It defaults to \r\n \t.
// @return {String} - The string without chars on both sides.
func Trim(_ context.Context, args ...core.Value) (core.Value, error) {
	err := core.ValidateArgs(args, 1, 2)

	if err != nil {
		return core.EmptyString, err
	}

	text := args[0].String()

	if len(args) > 1 {
		return core.NewString(strings.Trim(text, args[1].String())), nil
	}

	return core.NewString(strings.TrimSpace(text)), nil
}

// LTRIM returns the string value with whitespace stripped from the start only.
// @param {String} str - The string.
// @param {String} chars - Overrides the characters that should be removed from the string. It defaults to \r\n \t.
// @return {String} - The string without chars at the left-hand side.
func LTrim(_ context.Context, args ...core.Value) (core.Value, error) {
	err := core.ValidateArgs(args, 1, 2)

	if err != nil {
		return core.EmptyString, err
	}

	text := args[0].String()
	chars := " "

	if len(args) > 1 {
		chars = args[1].String()
	}

	return core.NewString(strings.TrimLeft(text, chars)), nil
}

// RTRIM returns the string value with whitespace stripped from the end only.
// @param {String} str - The string.
// @param {String} chars - Overrides the characters that should be removed from the string. It defaults to \r\n \t.
// @return {String} - The string without chars at the right-hand side.
func RTrim(_ context.Context, args ...core.Value) (core.Value, error) {
	err := core.ValidateArgs(args, 1, 2)

	if err != nil {
		return core.EmptyString, err
	}

	text := args[0].String()
	chars := " "

	if len(args) > 1 {
		chars = args[1].String()
	}

	return core.NewString(strings.TrimRight(text, chars)), nil
}
