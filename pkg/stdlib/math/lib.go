package math

import (
	"math"

	"github.com/MontFerret/ferret/pkg/runtime"
)

const (
	RadToDeg  = 180 / math.Pi
	DegToRad  = math.Pi / 180
	RadToGrad = 200 / math.Pi
	GradToDeg = math.Pi / 200
)

func RegisterLib(ns runtime.Namespace) error {
	return ns.RegisterFunctions(
		runtime.NewFunctionsFromMap(map[string]runtime.Function{
			"ABS":                 Abs,
			"ACOS":                Acos,
			"ASIN":                Asin,
			"ATAN":                Atan,
			"ATAN2":               Atan2,
			"AVERAGE":             Average,
			"CEIL":                Ceil,
			"COS":                 Cos,
			"DEGREES":             Degrees,
			"EXP":                 Exp,
			"EXP2":                Exp2,
			"FLOOR":               Floor,
			"LOG":                 Log,
			"LOG2":                Log2,
			"LOG10":               Log10,
			"MAX":                 Max,
			"MEDIAN":              Median,
			"MIN":                 Min,
			"PERCENTILE":          Percentile,
			"PI":                  Pi,
			"POW":                 Pow,
			"RADIANS":             Radians,
			"RAND":                Rand,
			"RANGE":               Range,
			"ROUND":               Round,
			"SIN":                 Sin,
			"SQRT":                Sqrt,
			"STDDEV_POPULATION":   StandardDeviationPopulation,
			"STDDEV_SAMPLE":       StandardDeviationSample,
			"SUM":                 Sum,
			"TAN":                 Tan,
			"VARIANCE_POPULATION": PopulationVariance,
			"VARIANCE_SAMPLE":     SampleVariance,
		}))
}

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
