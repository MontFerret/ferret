package testing

import (
	"context"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
)

// Fail returns an error.
// @param (String) - Message to display on error.
func Fail(_ context.Context, args ...core.Value) (core.Value, error) {
	err := core.ValidateArgs(args, 0, 1)

	if err != nil {
		return values.None, err
	}

	if len(args) > 0 {
		return values.None, core.Error(ErrAssertion, args[0].String())
	}

	return values.None, ErrAssertion
}
