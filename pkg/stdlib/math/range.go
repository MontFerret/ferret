package math

import (
	"context"
	stdmath "math"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

// RANGE returns an array of numbers in the specified range, optionally with increments other than 1.
// @param start {Int | Float} The value to start the range at (inclusive).
// @param end {Int | Float} The value to end the range with (inclusive).
// @param step {Int | Float} How much to change the value in every step. Positive steps ascend, negative steps descend, and zero is invalid.
// @return {Int[] | Float[]} arrayList of numbers in the specified range, optionally with increments other than 1.
func Range(ctx context.Context, args ...runtime.Value) (runtime.Value, error) {
	if err := runtime.ValidateArgs(args, 2, 3); err != nil {
		return runtime.None, err
	}

	if len(args) == 2 {
		return range2(ctx, args[0], args[1])
	}

	return range3(ctx, args[0], args[1], args[2])
}

// RANGE returns an array of numbers in the specified range, optionally with increments other than 1.
// @param start {Int | Float} The value to start the range at (inclusive).
// @param end {Int | Float} The value to end the range with (inclusive).
// @return {Int[] | Float[]} arrayList of numbers in the specified range, optionally with increments other than 1.
func range2(ctx context.Context, arg1, arg2 runtime.Value) (runtime.Value, error) {
	return range3(ctx, arg1, arg2, runtime.Float(1))
}

// RANGE returns an array of numbers in the specified range, optionally with increments other than 1.
// @param start {Int | Float} The value to start the range at (inclusive).
// @param end {Int | Float} The value to end the range with (inclusive).
// @param step {Int | Float} How much to change the value in every step. Positive steps ascend, negative steps descend, and zero is invalid.
// @return {Int[] | Float[]} arrayList of numbers in the specified range, optionally with increments other than 1.
func range3(ctx context.Context, arg1, arg2, arg3 runtime.Value) (runtime.Value, error) {
	if err := runtime.ValidateArgValue(arg1, 0, runtime.AssertNumber); err != nil {
		return runtime.None, err
	}

	if err := runtime.ValidateArgValue(arg2, 1, runtime.AssertNumber); err != nil {
		return runtime.None, err
	}

	if err := runtime.ValidateArgValue(arg3, 2, runtime.AssertNumber); err != nil {
		return runtime.None, err
	}

	step := toFloat(arg3)
	start := toFloat(arg1)
	end := toFloat(arg2)

	if err := validateRangeNumber(start, 0, "start"); err != nil {
		return runtime.None, err
	}

	if err := validateRangeNumber(end, 1, "end"); err != nil {
		return runtime.None, err
	}

	if err := validateRangeNumber(step, 2, "step"); err != nil {
		return runtime.None, err
	}

	if step == 0 {
		return runtime.None, rangeStepError("step must not be zero")
	}

	ascending := step > 0
	if (ascending && start > end) || (!ascending && start < end) {
		return runtime.NewArray(0), nil
	}

	if start != end && start+step == start {
		return runtime.None, rangeStepError("step is too small to advance the range")
	}

	capacity, err := rangeCapacity(start, end, step)
	if err != nil {
		return runtime.None, err
	}

	arr := runtime.NewArray(capacity)

	if ascending {
		err = appendAscendingRange(ctx, arr, start, end, step)
	} else {
		err = appendDescendingRange(ctx, arr, start, end, step)
	}

	if err != nil {
		return runtime.None, err
	}

	return arr, nil
}

func appendAscendingRange(ctx context.Context, arr *runtime.Array, start, end, step float64) error {
	for value := start; value <= end; {
		_ = arr.Append(ctx, runtime.NewFloat(value))

		if value == end {
			return nil
		}

		next := value + step
		if next == value {
			return rangeStepError("step is too small to advance the range")
		}

		value = next
	}

	return nil
}

func appendDescendingRange(ctx context.Context, arr *runtime.Array, start, end, step float64) error {
	for value := start; value >= end; {
		_ = arr.Append(ctx, runtime.NewFloat(value))

		if value == end {
			return nil
		}

		next := value + step
		if next == value {
			return rangeStepError("step is too small to advance the range")
		}

		value = next
	}

	return nil
}

func validateRangeNumber(value float64, pos int, name string) error {
	if stdmath.IsNaN(value) || stdmath.IsInf(value, 0) {
		return runtime.ArgError(
			runtime.Error(runtime.ErrInvalidArgument, name+" must be finite"),
			pos,
		)
	}

	return nil
}

func rangeStepError(message string) error {
	return runtime.ArgError(runtime.Error(runtime.ErrInvalidArgument, message), 2)
}

// rangeCapacity divides the endpoints separately when their subtraction overflows,
// preserving finite cardinalities for wide ranges that also use large steps.
func rangeCapacity(start, end, step float64) (int, error) {
	distance := stdmath.Abs(end - start)
	stepSize := stdmath.Abs(step)

	if stdmath.IsInf(distance, 0) {
		distance = stdmath.Abs(end/stepSize - start/stepSize)
	} else {
		distance /= stepSize
	}

	capacity := stdmath.Floor(distance) + 1
	if stdmath.IsNaN(capacity) || stdmath.IsInf(capacity, 0) || capacity >= float64(stdmath.MaxInt) {
		return 0, runtime.Error(runtime.ErrRange, "range length exceeds array capacity")
	}

	return int(capacity), nil
}
