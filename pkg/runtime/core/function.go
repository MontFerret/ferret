package core

import (
	"context"
	"fmt"
	"strings"
)

const MaxArgs = 65536

type (
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

type (
	// Functions is a map of functions and their names.
	Functions map[string]Function

	// Function is a common interface for all functions of FQL.
	Function = func(ctx context.Context, args ...Value) (Value, error)
)

// Get returns the function with the given name. If the function
// does not exist it returns nil, false.
func (fns Functions) Get(name string) (Function, bool) {
	fn, exists := fns[strings.ToUpper(name)]
	return fn, exists
}

// Set sets the function with the given name. If the function
// with the such name already exists it will be overwritten.
func (fns Functions) Set(name string, fn Function) {
	fns[strings.ToUpper(name)] = fn
}

// Unset delete the function with the given name.
func (fns Functions) Unset(name string) {
	delete(fns, strings.ToUpper(name))
}
