package testing

import (
	"context"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
)

// True asserts that value is true.
// @param (Mixed) - Value to test.
// @param (String) - Message to display on error.
func True(_ context.Context, args ...core.Value) (core.Value, error) {
	err := core.ValidateArgs(args, 1, 2)

	if err != nil {
		return values.None, err
	}

	if args[0] == values.True {
		return values.None, nil
	}

	if len(args) > 1 {
		return values.None, core.Error(ErrAssertion, args[0].String())
	}

	return values.None, core.Errorf(ErrAssertion, "expected %s to be true", args[0])
}
