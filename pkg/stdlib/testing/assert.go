package testing

import (
	"context"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/MontFerret/ferret/pkg/runtime/values/types"
)

// Assert checks whether value is a boolean value.
// @param value (Boolean|Array) - Expression(s) to test for truthiness.
// @param (String) - Message to display on error.
func Assert(_ context.Context, args ...core.Value) (core.Value, error) {
	err := core.ValidateArgs(args, 1, 2)

	if err != nil {
		return values.None, err
	}

	if err := core.ValidateType(args[0], types.Boolean, types.Array); err != nil {
		return values.None, err
	}

	result := values.True

	switch v := args[0].(type) {
	case values.Boolean:
		result = v
		break
	case *values.Array:
		v.ForEach(func(value core.Value, idx int) bool {
			if value.Compare(values.False) == 0 {
				result = values.False
				return false
			}

			return true
		})
	}

	if result {
		return values.None, nil
	}

	var message = "Expected to be true"

	if len(args) > 1 {
		message = args[1].String()
	}

	return values.None, core.Error(ErrAssertion, message)
}
