package math

import (
	"context"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/stdlib/arrays"
)

// range returns an array of numbers in the specified range, optionally with increments other than 1.
// @param start {Int | Float} The value to start the range at (inclusive).
// @param end {Int | Float} The value to end the range with (inclusive).
// @param step {Int | Float} How much to change the value in every step. Positive steps ascend, negative steps descend, and zero is invalid.
// @return {Float[]} Numbers in the inclusive range, using a default step of one.
// @deprecated Use arrays::range instead.
func Range(ctx context.Context, args ...runtime.Value) (runtime.Value, error) {
	return arrays.Range(ctx, args...)
}

// range returns an array of numbers in the specified range, optionally with increments other than 1.
// @param start {Int | Float} The value to start the range at (inclusive).
// @param end {Int | Float} The value to end the range with (inclusive).
// @return {Float[]} Numbers in the inclusive range, using a default step of one.
// @deprecated Use arrays::range instead.
func range2(ctx context.Context, arg1, arg2 runtime.Value) (runtime.Value, error) {
	return arrays.Range(ctx, arg1, arg2)
}

// range returns an array of numbers in the specified range, optionally with increments other than 1.
// @param start {Int | Float} The value to start the range at (inclusive).
// @param end {Int | Float} The value to end the range with (inclusive).
// @param step {Int | Float} How much to change the value in every step. Positive steps ascend, negative steps descend, and zero is invalid.
// @return {Float[]} Numbers in the inclusive range, using a default step of one.
// @deprecated Use arrays::range instead.
func range3(ctx context.Context, arg1, arg2, arg3 runtime.Value) (runtime.Value, error) {
	return arrays.Range(ctx, arg1, arg2, arg3)
}
