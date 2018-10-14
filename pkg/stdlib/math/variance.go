package math

import (
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	"math"
)

func variance(input *values.Array, sample values.Int) (values.Float, error) {
	if input.Length() == 0 {
		return values.NewFloat(math.NaN()), nil
	}

	m, _ := mean(input)

	var err error
	var variance values.Float

	input.ForEach(func(value core.Value, idx int) bool {
		err = core.ValidateType(value, core.IntType, core.FloatType)

		if err != nil {
			return false
		}

		n := values.Float(toFloat(value))

		variance += (n - m) * (n - m)

		return true
	})

	// When getting the mean of the squared differences
	// "sample" will allow us to know if it's a sample
	// or population and wether to subtract by one or not
	l := values.Float(input.Length() - (1 * sample))

	return variance / l, nil
}
