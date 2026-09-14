package math

import "github.com/MontFerret/ferret/v2/pkg/runtime"

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
