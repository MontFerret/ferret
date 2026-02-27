package math

import (
	"math"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

const (
	RadToDeg  = 180 / math.Pi
	DegToRad  = math.Pi / 180
	RadToGrad = 200 / math.Pi
	GradToDeg = math.Pi / 200
)

func RegisterLib(ns runtime.Namespace) {
	ns.Function().A0().
		Add("PI", Pi)

	ns.Function().A1().
		Add("ABS", Abs).
		Add("ACOS", Acos).
		Add("ASIN", Asin).
		Add("ATAN", Atan).
		Add("AVERAGE", Average).
		Add("CEIL", Ceil).
		Add("COS", Cos).
		Add("DEGREES", Degrees).
		Add("EXP", Exp).
		Add("EXP2", Exp2).
		Add("FLOOR", Floor).
		Add("LOG", Log).
		Add("LOG2", Log2).
		Add("LOG10", Log10).
		Add("MAX", Max).
		Add("MEDIAN", Median).
		Add("MIN", Min).
		Add("RADIANS", Radians).
		Add("ROUND", Round).
		Add("SIN", Sin).
		Add("SQRT", Sqrt).
		Add("STDDEV_POPULATION", StandardDeviationPopulation).
		Add("STDDEV_SAMPLE", StandardDeviationSample).
		Add("SUM", Sum).
		Add("TAN", Tan).
		Add("VARIANCE_POPULATION", PopulationVariance).
		Add("VARIANCE_SAMPLE", SampleVariance)

	ns.Function().A2().
		Add("ATAN2", Atan2).
		Add("POW", Pow)

	ns.Function().Var().
		Add("PERCENTILE", Percentile).
		Add("RAND", Rand).
		Add("RANGE", Range)
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
