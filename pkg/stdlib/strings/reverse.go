package strings

import (
	"context"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
)

// Reverse returns the reverse of the string value.
// @param text (String) - The string to revers
// @returns (String) - Returns a reversed version of the string.
func Reverse(_ context.Context, args ...core.Value) (core.Value, error) {
	err := core.ValidateArgs(args, 1, 1)

	if err != nil {
		return values.EmptyString, err
	}

	text := args[0].String()
	runes := []rune(text)
	size := len(runes)

	// Reverse
	for i := 0; i < size/2; i++ {
		runes[i], runes[size-1-i] = runes[size-1-i], runes[i]
	}

	return values.NewString(string(runes)), nil
}
