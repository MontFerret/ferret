package operators

import (
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/MontFerret/ferret/pkg/runtime/values/types"
)

func Range(left, right core.Value) (core.Value, error) {
	err := core.ValidateType(left, types.Int, types.Float)

	if err != nil {
		return values.None, err
	}

	err = core.ValidateType(right, types.Int, types.Float)

	if err != nil {
		return values.None, err
	}

	var start int
	var end int

	if left.Type() == types.Float {
		start = int(left.(values.Float))
	} else {
		start = int(left.(values.Int))
	}

	if right.Type() == types.Float {
		end = int(right.(values.Float))
	} else {
		end = int(right.(values.Int))
	}

	arr := values.NewArray(10)

	for i := start; i <= end; i++ {
		arr.Push(values.NewInt(i))
	}

	return arr, nil
}
