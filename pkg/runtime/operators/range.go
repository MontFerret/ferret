package operators

import (
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
)

func Range(left, right core.Value) (core.Value, error) {
	start := values.ToInt(left)
	end := values.ToInt(right)

	if start > end {
		return values.None, nil
	}

	return values.NewRange(uint64(start), uint64(end)), nil
}
