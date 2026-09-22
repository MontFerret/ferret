package math

import (
	"math"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

func toFloat(arg runtime.Value) float64 {
	switch v := arg.(type) {
	case runtime.Float:
		return float64(v)
	case runtime.Int:
		return float64(v)
	default:
		return 0
	}
}

// roundedInt uses an exclusive upper bound: float64(math.MaxInt64) is 2^63.
func roundedInt(value float64) (runtime.Value, error) {
	if math.IsNaN(value) || value < -0x1p63 || value >= 0x1p63 {
		return runtime.None, runtime.ArgError(runtime.Error(runtime.ErrRange, "rounded value exceeds the signed 64-bit integer range"), 0)
	}

	return runtime.NewInt64(int64(value)), nil
}
