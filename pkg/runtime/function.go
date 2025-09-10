package runtime

import (
	"context"
)

// MaxArgs defines the maximum number of arguments that a function can accept.
const MaxArgs = 65536

type (
	// Function is a common interface for functions with variable number of arguments.
	// All functions receive a context and a slice of values, returning a value and an error.
	Function = func(ctx context.Context, args ...Value) (Value, error)

	// Function0 is a common interface for functions with no arguments.
	Function0 = func(ctx context.Context) (Value, error)

	// Function1 is a common interface for functions with a single argument.
	Function1 = func(ctx context.Context, arg Value) (Value, error)

	// Function2 is a common interface for functions with two arguments.
	Function2 = func(ctx context.Context, arg1, arg2 Value) (Value, error)

	// Function3 is a common interface for functions with three arguments.
	Function3 = func(ctx context.Context, arg1, arg2, arg3 Value) (Value, error)

	// Function4 is a common interface for functions with four arguments.
	Function4 = func(ctx context.Context, arg1, arg2, arg3, arg4 Value) (Value, error)
)
