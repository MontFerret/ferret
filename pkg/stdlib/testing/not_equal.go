package testing

import (
	"context"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
)

// NotEqual asserts inequality of actual and expected values.
// @param (Mixed) - Actual value.
// @param (Mixed) - Expected value.
// @param (String) - Message to display on error.
func NotEqual(_ context.Context, args ...core.Value) (core.Value, error) {
	err := core.ValidateArgs(args, 2, 3)

	if err != nil {
		return values.None, err
	}

	actual := args[0]
	expected := args[1]

	if actual.Compare(expected) != 0 {
		return values.None, nil
	}

	if len(args) > 2 {
		return values.None, core.Error(ErrAssertion, args[2].String())
	}

	return values.None, core.Errorf(ErrAssertion, "expected %s to not be %s", actual, expected)
}
