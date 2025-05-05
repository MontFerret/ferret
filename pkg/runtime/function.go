package runtime

import (
	"context"
	"fmt"
)

const MaxArgs = 65536

type (
	// Functions is a container for functions.
	Functions map[string]Function

	// Function is a common interface for all functions of FQL.
	Function = func(ctx context.Context, args ...Value) (Value, error)
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

// NewFunctions returns new empty Functions.
func NewFunctions() Functions {
	return make(map[string]Function)
}

// NewFunctionsFromMap creates new Functions from map, where
// key is the name of the function and value is the function.
func NewFunctionsFromMap(funcs map[string]Function) Functions {
	fns := NewFunctions()

	for name, fn := range funcs {
		fns.Set(name, fn)
	}

	return fns
}

// Has returns true if the function with the given name exists.
func (fns Functions) Has(name string) bool {
	_, exists := fns[name]
	return exists
}

// Get returns the function with the given name. If the function
// does not exist it returns nil, false.
func (fns Functions) Get(name string) (Function, bool) {
	fn, exists := fns[name]
	return fn, exists
}

// MustGet returns the function with the given name. If the function
// does not exist it panics.
func (fns Functions) MustGet(name string) Function {
	return fns[name]
}

// Set sets the function with the given name. If the function
// with the such name already exists it will be overwritten.
func (fns Functions) Set(name string, fn Function) {
	fns[name] = fn
}

// Unset delete the function with the given name.
func (fns Functions) Unset(name string) {
	delete(fns, name)
}

// Names returns the names of the internal functions.
func (fns Functions) Names() []string {
	names := make([]string, 0, len(fns))

	for name := range fns {
		names = append(names, name)
	}

	return names
}

// Unwrap returns the internal map of functions.
func (fns Functions) Unwrap() map[string]Function {
	return fns
}
