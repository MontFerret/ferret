package operators

import (
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
)

func Range(left, right core.Value) (core.Value, error) {
	start := values.ToInt(left)
	end := values.ToInt(right)

	return values.NewRange(int64(start), int64(end)), nil
}
