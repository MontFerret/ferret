package vm

import "github.com/MontFerret/ferret/v2/pkg/runtime"

func negate(input runtime.Value) runtime.Value {
	switch val := input.(type) {
	case runtime.Int:
		return -val
	case runtime.Float:
		return -val
	case runtime.Boolean:
		return !val
	default:
		return runtime.None
	}
}

func negative(input runtime.Value) runtime.Value {
	switch value := input.(type) {
	case runtime.Int:
		return -value
	case runtime.Float:
		return -value
	default:
		return runtime.None
	}
}

func positive(input runtime.Value) runtime.Value {
	switch value := input.(type) {
	case runtime.Int:
		return +value
	case runtime.Float:
		return +value
	default:
		return runtime.None
	}
}
