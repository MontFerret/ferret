package runtime

import "context"

// MaxArgs defines the maximum number of arguments that a function can accept.
const MaxArgs = 65536

type (
	// LegacyFunction represents a function type used by previous versions of the engine.
	LegacyFunction = func(ctx context.Context, args ...Value) (Value, error)

	// Function is a common interface for functions with variable number of arguments.
	// All functions receive a context and a slice of values, returning a value and an error.
	Function = func(ctx Context, args ...Value) (Value, error)

	// Function0 is a common interface for functions with no arguments.
	Function0 = func(ctx Context) (Value, error)

	// Function1 is a common interface for functions with a single argument.
	Function1 = func(ctx Context, arg Value) (Value, error)

	// Function2 is a common interface for functions with two arguments.
	Function2 = func(ctx Context, arg1, arg2 Value) (Value, error)

	// Function3 is a common interface for functions with three arguments.
	Function3 = func(ctx Context, arg1, arg2, arg3 Value) (Value, error)

	// Function4 is a common interface for functions with four arguments.
	Function4 = func(ctx Context, arg1, arg2, arg3, arg4 Value) (Value, error)
)

// FromLegacyFunction converts a LegacyFunction to the new Function type.
func FromLegacyFunction(fn LegacyFunction) Function {
	return func(ctx Context, args ...Value) (Value, error) {
		return fn(ctx, args...)
	}
}
