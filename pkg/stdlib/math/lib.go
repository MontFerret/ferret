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
		runtime.
			NewFunctionsBuilder().
			Set1("ABS", Abs).
			Set1("ACOS", Acos).
			Set1("ASIN", Asin).
			Set1("ATAN", Atan).
			Set2("ATAN2", Atan2).
			Set1("AVERAGE", Average).
			Set1("CEIL", Ceil).
			Set1("COS", Cos).
			Set1("DEGREES", Degrees).
			Set1("EXP", Exp).
			Set1("EXP2", Exp2).
			Set1("FLOOR", Floor).
			Set1("LOG", Log).
			Set1("LOG2", Log2).
			Set1("LOG10", Log10).
			Set1("MAX", Max).
			Set1("MEDIAN", Median).
			Set1("MIN", Min).
			Set("PERCENTILE", Percentile).
			Set0("PI", Pi).
			Set2("POW", Pow).
			Set1("RADIANS", Radians).
			Set("RAND", Rand).
			Set("RANGE", Range).
			Set1("ROUND", Round).
			Set1("SIN", Sin).
			Set1("SQRT", Sqrt).
			Set1("STDDEV_POPULATION", StandardDeviationPopulation).
			Set1("STDDEV_SAMPLE", StandardDeviationSample).
			Set1("SUM", Sum).
			Set1("TAN", Tan).
			Set1("VARIANCE_POPULATION", PopulationVariance).
			Set1("VARIANCE_SAMPLE", SampleVariance).
			Build(),
	)
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
