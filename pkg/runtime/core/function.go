package core

import (
	"context"
	"fmt"
)

const MaxArgs = 65536

type (
	Function = func(ctx context.Context, args ...Value) (Value, error)

	Functions map[string]Function

	Namespace interface {
		Namespace(name string) Namespace
		RegisterFunction(name string, fun Function) error
		RegisterFunctions(funs Functions) error
		RegisteredFunctions() []string
		RemoveFunction(name string)
	}
)

func ValidateArgs(args []Value, minimum, maximum int) error {
	count := len(args)

	if count < minimum || count > maximum {
		return Error(
			ErrInvalidArgumentNumber,
			fmt.Sprintf(
				"expected number of arguments %d-%d, but got %d",
				minimum,
				maximum,
				len(args)))
	}

	return nil
}
