package core

import (
	"context"
	"fmt"
)

type Function = func(ctx context.Context, args ...Value) (Value, error)

func ValidateArgs(args []Value, required int) error {
	if len(args) != required {
		return Error(ErrMissedArgument, fmt.Sprintf("expected %d, but got %d arguments", required, len(args)))
	}

	return nil
}
