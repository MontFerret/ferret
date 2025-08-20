package operators

import "github.com/MontFerret/ferret/pkg/runtime"

func Negate(input runtime.Value) runtime.Value {
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

func Negative(input runtime.Value) runtime.Value {
	switch value := input.(type) {
	case runtime.Int:
		return -value
	case runtime.Float:
		return -value
	default:
		return runtime.None
	}
}

func Positive(input runtime.Value) runtime.Value {
	switch value := input.(type) {
	case runtime.Int:
		return +value
	case runtime.Float:
		return +value
	default:
		return runtime.None
	}
}

func ToNumberOrString(input runtime.Value) runtime.Value {
	switch value := input.(type) {
	case runtime.Int, runtime.Float, runtime.String:
		return value
	default:
		return runtime.ToString(value)
	}
}
