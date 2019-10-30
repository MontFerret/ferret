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
		RegisterFunctions(funs *Functions) error
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
	// Functions is a container for functions.
	Functions struct {
		functions map[string]Function
	}

	// Function is a common interface for all functions of FQL.
	Function = func(ctx context.Context, args ...Value) (Value, error)
)

// NewFunctions returns new empty Functions.
func NewFunctions() *Functions {
	return &Functions{
		functions: make(map[string]Function),
	}
}

// NewFunctionsFromMap creates new Functions from map, where
// key is the name of the function and value is the function.
func NewFunctionsFromMap(funcs map[string]Function) *Functions {
	fns := NewFunctions()

	for name, fn := range funcs {
		fns.Set(name, fn)
	}

	return fns
}

// Get returns the function with the given name. If the function
// does not exist it returns nil, false.
func (fns *Functions) Get(name string) (Function, bool) {
	fn, exists := fns.functions[strings.ToUpper(name)]
	return fn, exists
}

// Set sets the function with the given name. If the function
// with the such name already exists it will be overwritten.
func (fns *Functions) Set(name string, fn Function) {
	// the preferred way to create Functions is NewFunctions.
	// But just in case, if someone creates differently
	if fns.functions == nil {
		fns.functions = make(map[string]Function, 1)
	}

	fns.functions[strings.ToUpper(name)] = fn
}

// Unset delete the function with the given name.
func (fns *Functions) Unset(name string) {
	delete(fns.functions, strings.ToUpper(name))
}

// Names returns the names of the internal functions.
func (fns *Functions) Names() []string {
	names := make([]string, 0, len(fns.functions))

	for name := range fns.functions {
		names = append(names, name)
	}

	return names
}
