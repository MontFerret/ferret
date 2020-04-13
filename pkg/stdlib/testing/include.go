package testing

import (
	"context"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
)

// Include asserts that haystack includes needle.
// Can be used to assert the inclusion of a value in an array, a substring in a string, or a subset of properties in an object.
// @param (Mixed) - Haystack value.
// @param (Mixed) - Needle value.
// @param (String) - Message to display on error.
func Include(_ context.Context, args ...core.Value) (core.Value, error) {
	err := core.ValidateArgs(args, 2, 3)

	if err != nil {
		return values.None, err
	}

	haystack := args[0]
	needle := args[1]

}
