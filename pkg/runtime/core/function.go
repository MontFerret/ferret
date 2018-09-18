package core

import (
	"context"
	"fmt"
)

type Function = func(ctx context.Context, args ...Value) (Value, error)

func ValidateArgs(inputs []Value, required int) error {
	if len(inputs) != required {
		return Error(ErrMissedArgument, fmt.Sprintf("expected %d, but got %d arguments", required, len(inputs)))
	}

	return nil
}
