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

// RegisterLib registers canonical math functions and deprecated global compatibility functions.
// @namespace math
func RegisterLib(ns runtime.Namespace) {
	registerCanonical(ns.Namespace("math"))
	registerLegacy(ns)
}

func registerCanonical(ns runtime.Namespace) {
	ns.Function().A0().
		Add("e", E).
		Add("pi", Pi)

	ns.Function().A1().
		Add("abs", Abs).
		Add("acos", Acos).
		Add("asin", Asin).
		Add("atan", Atan).
		Add("mean", canonicalMean).
		Add("cbrt", Cbrt).
		Add("ceil", Ceil).
		Add("cos", Cos).
		Add("degrees", Degrees).
		Add("exp", Exp).
		Add("exp2", Exp2).
		Add("expm1", Expm1).
		Add("floor", Floor).
		Add("log", Log).
		Add("log1p", Log1p).
		Add("log2", Log2).
		Add("log10", Log10).
		Add("max", canonicalMax).
		Add("median", canonicalMedian).
		Add("min", canonicalMin).
		Add("radians", Radians).
		Add("round", Round).
		Add("sign", Sign).
		Add("sin", Sin).
		Add("sqrt", Sqrt).
		Add("stddev", canonicalStddev).
		Add("stddev_sample", canonicalStddevSample).
		Add("sum", canonicalSum).
		Add("tan", Tan).
		Add("trunc", Trunc).
		Add("variance", canonicalVariance).
		Add("variance_sample", canonicalVarianceSample)

	ns.Function().A2().
		Add("atan2", Atan2).
		Add("hypot", Hypot).
		Add("percentile", canonicalPercentile).
		Add("pow", Pow)

	ns.Function().A3().
		Add("clamp", Clamp)
}
