package math

import (
	"github.com/MontFerret/ferret/pkg/runtime/internal"
	"math"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values/types"
)

func mean(input *internal.Array) (core.Float, error) {
	if input.Length() == 0 {
		return core.NewFloat(math.NaN()), nil
	}

	var err error
	var sum float64

	input.ForEach(func(value core.Value, idx int) bool {
		err = core.ValidateType(value, types.Int, types.Float)

		if err != nil {
			return false
		}

		sum += toFloat(value)

		return true
	})

	if err != nil {
		return 0, err
	}

	return core.NewFloat(sum / float64(input.Length())), nil
}
