package testing

import (
	"context"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/MontFerret/ferret/pkg/stdlib/collections"
)

// Len asserts that a measurable value has a length or size with the expected value.
// @param (Measurable) - Measurable value.
// @param (Mixed) - Target length.
// @param (String) - Message to display on error.
func Len(ctx context.Context, args ...core.Value) (core.Value, error) {
	err := core.ValidateArgs(args, 2, 3)

	if err != nil {
		return values.None, err
	}

	col := args[0]
	size := args[1]

	out, err := collections.Length(ctx, col)

	if err != nil {
		return values.None, err
	}

	if out.Compare(size) == 0 {
		return values.None, nil
	}

	if len(args) > 2 {
		return values.None, core.Error(ErrAssertion, args[2].String())
	}

	return values.None, core.Errorf(ErrAssertion, "expected %s to has size %d", col, size)
}
