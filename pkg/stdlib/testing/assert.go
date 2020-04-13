package testing

import (
	"context"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/MontFerret/ferret/pkg/runtime/values/types"
)

// Assert checks whether value is a boolean value.
// @param value (Boolean) - Expression to test for truthiness.
// @param (String) - Message to display on error.
func Assert(_ context.Context, args ...core.Value) (core.Value, error) {
	err := core.ValidateArgs(args, 1, 2)

	if err != nil {
		return values.None, err
	}

	if err := core.ValidateType(args[0], types.Boolean); err != nil {
		return values.None, err
	}

	exp := args[0].(values.Boolean)

	if exp {
		return values.None, nil
	}

	var message = "Expected to be true"

	if len(args) > 1 {
		message = args[1].String()
	}

	return values.None, core.Error(ErrAssertion, message)
}
