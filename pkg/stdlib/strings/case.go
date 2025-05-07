package strings

import (
	"context"
	"strings"

	"github.com/MontFerret/ferret/pkg/runtime"
)

// LOWER converts strings to their lower-case counterparts. All other characters are returned unchanged.
// @param {String} str - The source string.
// @return {String} - THis string in lower case.
func Lower(_ context.Context, args ...runtime.Value) (runtime.Value, error) {
	err := runtime.ValidateArgs(args, 1, 1)

	if err != nil {
		return runtime.EmptyString, err
	}

	text := strings.ToLower(args[0].String())

	return runtime.NewString(text), nil
}

// UPPER converts strings to their upper-case counterparts. All other characters are returned unchanged.
// @param {String} str - The source string.
// @return {String} - THis string in upper case.
func Upper(_ context.Context, args ...runtime.Value) (runtime.Value, error) {
	err := runtime.ValidateArgs(args, 1, 1)

	if err != nil {
		return runtime.EmptyString, err
	}

	text := strings.ToUpper(args[0].String())

	return runtime.NewString(text), nil
}
