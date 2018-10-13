package math

import (
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	"math"
)

func mean(input *values.Array) (values.Float, error) {
	if input.Length() == 0 {
		return values.NewFloat(math.NaN()), nil
	}

	var err error
	var sum float64

	input.ForEach(func(value core.Value, idx int) bool {
		err = core.ValidateType(value, core.FloatType, core.IntType)

		if err != nil {
			return false
		}

		sum += toFloat(value)

		return true
	})

	if err != nil {
		return 0, err
	}

	return values.NewFloat(sum / float64(input.Length())), nil
}
