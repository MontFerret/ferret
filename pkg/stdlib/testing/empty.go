package testing

import (
	"context"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/MontFerret/ferret/pkg/stdlib/collections"
)

// Empty asserts that the target does not contain any values.
// @param value (Measurable|Binary|Object|Array|String) - Value to test.
// @param (String) - Message to display on error.
func Empty(ctx context.Context, args ...core.Value) (core.Value, error) {
	err := core.ValidateArgs(args, 1, 2)

	if err != nil {
		return values.None, err
	}

	size, err := collections.Length(ctx, args[0])

	if err != nil {
		return values.None, err
	}

	if values.ToInt(size) == 0 {
		return values.None, nil
	}

	if len(args) > 1 {
		return values.None, core.Error(ErrAssertion, args[1].String())
	}

	return values.None, core.Error(ErrAssertion, "expected to be empty")
}
