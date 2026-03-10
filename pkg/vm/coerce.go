package vm

import (
	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

func coerceBool(input runtime.Value) runtime.Boolean {
	if input == nil || input == runtime.None {
		return runtime.False
	}

	switch val := input.(type) {
	case runtime.Boolean:
		return val
	case runtime.String:
		return val != ""
	case runtime.Int:
		return val != 0
	case runtime.Float:
		return val != 0
	case runtime.DateTime:
		return val.IsZero() != true
	default:
		return runtime.True
	}
}
