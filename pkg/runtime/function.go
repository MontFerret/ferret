package runtime

import (
	"context"
)

// MaxArgs defines the maximum number of arguments that a function can accept.
const MaxArgs = 65536

type (
	Functions interface {
		Has(name string) bool
		F() FunctionCollection[Function]
		F0() FunctionCollection[Function0]
		F1() FunctionCollection[Function1]
		F2() FunctionCollection[Function2]
		F3() FunctionCollection[Function3]
		F4() FunctionCollection[Function4]
		SetAll(otherFns Functions) Functions
		Unset(name string) Functions
		UnsetAll() Functions
		Names() []string
	}

	FunctionsBuilder interface {
		Set(name string, fn Function) FunctionsBuilder
		Set0(name string, fn Function0) FunctionsBuilder
		Set1(name string, fn Function1) FunctionsBuilder
		Set2(name string, fn Function2) FunctionsBuilder
		Set3(name string, fn Function3) FunctionsBuilder
		Set4(name string, fn Function4) FunctionsBuilder
		Build() Functions
	}

	FunctionConstraint interface {
		Function | Function0 | Function1 | Function2 | Function3 | Function4
	}

	FunctionCollection[T FunctionConstraint] interface {
		Has(name string) bool
		Set(name string, fn T) FunctionCollection[T]
		SetAll(otherFns FunctionCollection[T]) FunctionCollection[T]
		Get(name string) (T, bool)
		GetAll() map[string]T
		Unset(name string) FunctionCollection[T]
		UnsetAll() FunctionCollection[T]
		Names() []string
		Size() int
	}

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
