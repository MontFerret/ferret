package operators

import (
	"context"
	"strings"

	"github.com/MontFerret/ferret/pkg/runtime"
)

func Contains(ctx context.Context, input runtime.Value, value runtime.Value) runtime.Boolean {
	switch val := input.(type) {
	case runtime.List:
		idx, err := val.IndexOf(ctx, value)

		if err != nil {
			return runtime.False
		}

		return idx > -1
	case runtime.Map:
		containsValue, err := val.ContainsValue(ctx, value)

		if err != nil {
			return runtime.False
		}

		return containsValue
	case runtime.String:
		return runtime.Boolean(strings.Contains(val.String(), value.String()))
	default:
		return false
	}
}
