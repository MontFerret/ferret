package sdk

import "github.com/MontFerret/ferret/v2/pkg/runtime"

// FunctionDef describes a named Ferret host function before registration.
type FunctionDef struct {
	function any
	name     string
}

// Func creates an arity-specific function definition for RegisterFunctions.
// RegisterFunctions stores the qualified name as lowercase and resolves it case-insensitively in FQL.
// Definitions may share a name when their fixed arities differ or one is variadic.
func Func[T runtime.FunctionConstraint](name string, fn T) FunctionDef {
	return FunctionDef{
		name:     name,
		function: fn,
	}
}
