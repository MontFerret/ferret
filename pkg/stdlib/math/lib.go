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
		Add("pi", Pi).
		Add("rand", rand0)

	ns.Function().A1().
		Add("abs", Abs).
		Add("acos", Acos).
		Add("asin", Asin).
		Add("atan", Atan).
		Add("average", Average).
		Add("ceil", Ceil).
		Add("cos", Cos).
		Add("degrees", Degrees).
		Add("exp", Exp).
		Add("exp2", Exp2).
		Add("floor", Floor).
		Add("log", Log).
		Add("log2", Log2).
		Add("log10", Log10).
		Add("max", Max).
		Add("median", Median).
		Add("min", Min).
		Add("radians", Radians).
		Add("rand", rand1).
		Add("round", Round).
		Add("sin", Sin).
		Add("sqrt", Sqrt).
		Add("stddev_population", StandardDeviationPopulation).
		Add("stddev_sample", StandardDeviationSample).
		Add("sum", Sum).
		Add("tan", Tan).
		Add("variance_population", PopulationVariance).
		Add("variance_sample", SampleVariance)

	ns.Function().A2().
		Add("atan2", Atan2).
		Add("percentile", percentile2).
		Add("pow", Pow).
		Add("rand", rand2).
		Add("range", range2)

	ns.Function().A3().
		Add("percentile", percentile3).
		Add("range", range3)
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
