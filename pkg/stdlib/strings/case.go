package strings

import (
	"context"
	"strings"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
)

// Lower converts strings to their lower-case counterparts. All other characters are returned unchanged.
// @param src (String) - The source string.
// @returns (String) - THis string in lower case.
func Lower(_ context.Context, args ...core.Value) (core.Value, error) {
	err := core.ValidateArgs(args, 1, 1)

	if err != nil {
		return values.EmptyString, err
	}

	text := strings.ToLower(args[0].String())

	return values.NewString(text), nil
}

// Upper converts strings to their upper-case counterparts. All other characters are returned unchanged.
// @param src (String) - The source string.
// @returns (String) - THis string in upper case.
func Upper(_ context.Context, args ...core.Value) (core.Value, error) {
	err := core.ValidateArgs(args, 1, 1)

	if err != nil {
		return values.EmptyString, err
	}

	text := strings.ToUpper(args[0].String())

	return values.NewString(text), nil
}
