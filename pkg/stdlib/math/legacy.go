package math

import (
	"context"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

// Separate declarations attach deprecation metadata only to the global surface.
func registerLegacy(ns runtime.Namespace) {
	ns.Function().A0().
		Add("pi", legacyPi).
		Add("rand", rand0)

	ns.Function().A1().
		Add("abs", legacyAbs).
		Add("acos", legacyAcos).
		Add("asin", legacyAsin).
		Add("atan", legacyAtan).
		Add("average", Average).
		Add("ceil", legacyCeil).
		Add("cos", legacyCos).
		Add("degrees", legacyDegrees).
		Add("exp", legacyExp).
		Add("exp2", legacyExp2).
		Add("floor", legacyFloor).
		Add("log", legacyLog).
		Add("log2", legacyLog2).
		Add("log10", legacyLog10).
		Add("max", Max).
		Add("median", Median).
		Add("min", Min).
		Add("radians", legacyRadians).
		Add("rand", rand1).
		Add("round", legacyRound).
		Add("sin", legacySin).
		Add("sqrt", legacySqrt).
		Add("stddev_population", StandardDeviationPopulation).
		Add("stddev_sample", StandardDeviationSample).
		Add("sum", Sum).
		Add("tan", legacyTan).
		Add("variance_population", PopulationVariance).
		Add("variance_sample", SampleVariance)

	ns.Function().A2().
		Add("atan2", legacyAtan2).
		Add("percentile", percentile2).
		Add("pow", legacyPow).
		Add("rand", rand2).
		Add("range", range2)

	ns.Function().A3().
		Add("percentile", percentile3).
		Add("range", range3)
}

// pi returns Pi value.
// @return {Float} Pi value.
// @deprecated Use math::pi instead.
func legacyPi(ctx context.Context) (runtime.Value, error) {
	return Pi(ctx)
}

// abs returns the absolute value of a given number.
// @param number {Int | Float} Input number.
// @return {Float} The absolute value of a given number.
// @deprecated Use math::abs instead.
func legacyAbs(ctx context.Context, arg1 runtime.Value) (runtime.Value, error) {
	return Abs(ctx, arg1)
}

// acos returns the arccosine, in radians, of a given number.
// @param number {Int | Float} Input number.
// @return {Float} The arccosine, in radians, of a given number.
// @deprecated Use math::acos instead.
func legacyAcos(ctx context.Context, arg1 runtime.Value) (runtime.Value, error) {
	return Acos(ctx, arg1)
}

// asin returns the arcsine, in radians, of a given number.
// @param number {Int | Float} Input number.
// @return {Float} The arcsine, in radians, of a given number.
// @deprecated Use math::asin instead.
func legacyAsin(ctx context.Context, arg1 runtime.Value) (runtime.Value, error) {
	return Asin(ctx, arg1)
}

// atan returns the arctangent, in radians, of a given number.
// @param number {Int | Float} Input number.
// @return {Float} The arctangent, in radians, of a given number.
// @deprecated Use math::atan instead.
func legacyAtan(ctx context.Context, arg1 runtime.Value) (runtime.Value, error) {
	return Atan(ctx, arg1)
}

// ceil returns the least integer value greater than or equal to a given value.
// @param number {Int | Float} Input number.
// @return {Int} The least integer value greater than or equal to a given value.
// @deprecated Use math::ceil instead.
func legacyCeil(ctx context.Context, arg1 runtime.Value) (runtime.Value, error) {
	return Ceil(ctx, arg1)
}

// cos returns the cosine of a given number.
// @param number {Int | Float} Input number.
// @return {Float} The cosine of a given number.
// @deprecated Use math::cos instead.
func legacyCos(ctx context.Context, arg1 runtime.Value) (runtime.Value, error) {
	return Cos(ctx, arg1)
}

// degrees returns the angle converted from radians to degrees.
// @param number {Int | Float} The input number.
// @return {Float} The angle in degrees
// @deprecated Use math::degrees instead.
func legacyDegrees(ctx context.Context, arg1 runtime.Value) (runtime.Value, error) {
	return Degrees(ctx, arg1)
}

// exp returns Euler's constant (2.71828...) raised to the power of value.
// @param number {Int | Float} Input number.
// @return {Float} Euler's constant raised to the power of value.
// @deprecated Use math::exp instead.
func legacyExp(ctx context.Context, arg1 runtime.Value) (runtime.Value, error) {
	return Exp(ctx, arg1)
}

// exp2 returns 2 raised to the power of value.
// @param number {Int | Float} Input number.
// @return {Float} 2 raised to the power of value.
// @deprecated Use math::exp2 instead.
func legacyExp2(ctx context.Context, arg1 runtime.Value) (runtime.Value, error) {
	return Exp2(ctx, arg1)
}

// floor returns the greatest integer value less than or equal to a given value.
// @param number {Int | Float} Input number.
// @return {Int} The greatest integer value less than or equal to a given value.
// @deprecated Use math::floor instead.
func legacyFloor(ctx context.Context, arg1 runtime.Value) (runtime.Value, error) {
	return Floor(ctx, arg1)
}

// log returns the natural logarithm of a given value.
// @param number {Int | Float} Input number.
// @return {Float} The natural logarithm of a given value.
// @deprecated Use math::log instead.
func legacyLog(ctx context.Context, arg1 runtime.Value) (runtime.Value, error) {
	return Log(ctx, arg1)
}

// log2 returns the binary logarithm of a given value.
// @param number {Int | Float} Input number.
// @return {Float} The binary logarithm of a given value.
// @deprecated Use math::log2 instead.
func legacyLog2(ctx context.Context, arg1 runtime.Value) (runtime.Value, error) {
	return Log2(ctx, arg1)
}

// log10 returns the decimal logarithm of a given value.
// @param number {Int | Float} Input number.
// @return {Float} The decimal logarithm of a given value.
// @deprecated Use math::log10 instead.
func legacyLog10(ctx context.Context, arg1 runtime.Value) (runtime.Value, error) {
	return Log10(ctx, arg1)
}

// radians returns the angle converted from degrees to radians.
// @param number {Int | Float} The input number.
// @return {Float} The angle in radians.
// @deprecated Use math::radians instead.
func legacyRadians(ctx context.Context, arg1 runtime.Value) (runtime.Value, error) {
	return Radians(ctx, arg1)
}

// round returns the nearest integer, rounding half away from zero.
// @param number {Int | Float} Input number.
// @return {Int} The nearest integer, rounding half away from zero.
// @deprecated Use math::round instead.
func legacyRound(ctx context.Context, arg1 runtime.Value) (runtime.Value, error) {
	return Round(ctx, arg1)
}

// sin returns the sine of the radian argument.
// @param number {Int | Float} Input number.
// @return {Float} The sin, in radians, of a given number.
// @deprecated Use math::sin instead.
func legacySin(ctx context.Context, arg1 runtime.Value) (runtime.Value, error) {
	return Sin(ctx, arg1)
}

// sqrt returns the square root of a given number.
// @param value {Int | Float} A number.
// @return {Float} The square root.
// @deprecated Use math::sqrt instead.
func legacySqrt(ctx context.Context, arg1 runtime.Value) (runtime.Value, error) {
	return Sqrt(ctx, arg1)
}

// tan returns the tangent of a given number.
// @param number {Int | Float} A number.
// @return {Float} The tangent.
// @deprecated Use math::tan instead.
func legacyTan(ctx context.Context, arg1 runtime.Value) (runtime.Value, error) {
	return Tan(ctx, arg1)
}

// atan2 returns the arc tangent of y/x, using the signs of the two to determine the quadrant of the return value.
// @param number1 {Int | Float} Input number.
// @param number2 {Int | Float} Input number.
// @return {Float} The arc tangent of y/x, using the signs of the two to determine the quadrant of the return value.
// @deprecated Use math::atan2 instead.
func legacyAtan2(ctx context.Context, arg1, arg2 runtime.Value) (runtime.Value, error) {
	return Atan2(ctx, arg1, arg2)
}

// pow returns the base to the exponent value.
// @param base {Int | Float} The base value.
// @param exp {Int | Float} The exponent value.
// @return {Float} The exponentiated value.
// @deprecated Use math::pow instead.
func legacyPow(ctx context.Context, arg1, arg2 runtime.Value) (runtime.Value, error) {
	return Pow(ctx, arg1, arg2)
}
