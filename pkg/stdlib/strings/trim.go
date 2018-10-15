package strings

import (
	"context"
	"strings"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
)

// Trim returns the string value with whitespace stripped from the start and/or end.
// @param value (String) - The string.
// @param chars (String) - Overrides the characters that should be removed from the string. It defaults to \r\n \t.
// @returns (String) - The string without chars on both sides.
func Trim(_ context.Context, args ...core.Value) (core.Value, error) {
	err := core.ValidateArgs(args, 1, 2)

	if err != nil {
		return values.EmptyString, err
	}

	text := args[0].String()

	if len(args) > 1 {
		return values.NewString(strings.Trim(text, args[1].String())), nil
	}

	return values.NewString(strings.TrimSpace(text)), nil
}

// LTrim returns the string value with whitespace stripped from the start only.
// @param value (String) - The string.
// @param chars (String) - Overrides the characters that should be removed from the string. It defaults to \r\n \t.
// @returns (String) - The string without chars at the left-hand side.
func LTrim(_ context.Context, args ...core.Value) (core.Value, error) {
	err := core.ValidateArgs(args, 1, 2)

	if err != nil {
		return values.EmptyString, err
	}

	text := args[0].String()
	chars := " "

	if len(args) > 1 {
		chars = args[1].String()
	}

	return values.NewString(strings.TrimLeft(text, chars)), nil
}

// RTrim returns the string value with whitespace stripped from the end only.
// @param value (String) - The string.
// @param chars (String) - Overrides the characters that should be removed from the string. It defaults to \r\n \t.
// @returns (String) - The string without chars at the right-hand side.
func RTrim(_ context.Context, args ...core.Value) (core.Value, error) {
	err := core.ValidateArgs(args, 1, 2)

	if err != nil {
		return values.EmptyString, err
	}

	text := args[0].String()
	chars := " "

	if len(args) > 1 {
		chars = args[1].String()
	}

	return values.NewString(strings.TrimRight(text, chars)), nil
}
