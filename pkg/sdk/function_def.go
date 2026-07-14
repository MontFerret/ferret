package sdk

import "github.com/MontFerret/ferret/v2/pkg/runtime"

// FunctionDef describes a named Ferret host function before registration.
type FunctionDef struct {
	function any
	name     string
}

// Func creates a function definition for RegisterFunctions.
func Func[T runtime.FunctionConstraint](name string, fn T) FunctionDef {
	return FunctionDef{
		name:     name,
		function: fn,
	}
}
