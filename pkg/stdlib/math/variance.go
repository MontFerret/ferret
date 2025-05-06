package math

import (
	"math"

	"github.com/MontFerret/ferret/pkg/runtime/internal"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values/types"
)

func variance(input *internal.Array, sample core.Int) core.Float {
	if input.Length() == 0 {
		return core.NewFloat(math.NaN())
	}

	m, _ := mean(input)

	var err error
	var variance core.Float

	input.ForEach(func(value core.Value, idx int) bool {
		err = core.ValidateType(value, types.Int, types.Float)

		if err != nil {
			return false
		}

		n := core.Float(toFloat(value))

		variance += (n - m) * (n - m)

		return true
	})

	// When getting the mean of the squared differences
	// "sample" will allow us to know if it's a sample
	// or population and whether to subtract by one or not
	l := core.Float(input.Length() - (1 * int(sample)))

	return variance / l
}
