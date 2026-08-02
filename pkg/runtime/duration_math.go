package runtime

import (
	"context"
	"errors"
	"math"
	"math/big"
)

// AddChecked applies checked numeric, Duration, and DateTime addition.
func AddChecked(ctx context.Context, left, right Value) (Value, error) {
	var leftDuration, rightDuration Duration
	leftIsDuration := false
	rightIsDuration := false

	switch value := left.(type) {
	case DateTime:
		return addDateTimeChecked(ctx, left, right)
	case Duration:
		leftDuration = value
		leftIsDuration = true
	}

	switch value := right.(type) {
	case DateTime:
		return addDateTimeChecked(ctx, left, right)
	case Duration:
		rightDuration = value
		rightIsDuration = true
	}

	if !leftIsDuration && !rightIsDuration {
		return Add(ctx, left, right), nil
	}

	if !leftIsDuration {
		var err error
		leftDuration, err = ToDuration(ctx, left)

		if err != nil {
			return None, err
		}
	}

	if !rightIsDuration {
		var err error
		rightDuration, err = ToDuration(ctx, right)

		if err != nil {
			return None, err
		}
	}

	result, ok := addDurationInt64(int64(leftDuration), int64(rightDuration))
	if !ok {
		return None, durationRangeError("addition")
	}

	return Duration(result), nil
}

// SubtractChecked applies checked numeric, Duration, and DateTime subtraction.
func SubtractChecked(ctx context.Context, left, right Value) (Value, error) {
	var leftDuration, rightDuration Duration
	leftIsDuration := false
	rightIsDuration := false

	switch value := left.(type) {
	case DateTime:
		return subtractDateTimeChecked(ctx, left, right)
	case Duration:
		leftDuration = value
		leftIsDuration = true
	}

	switch value := right.(type) {
	case DateTime:
		return subtractDateTimeChecked(ctx, left, right)
	case Duration:
		rightDuration = value
		rightIsDuration = true
	}

	if !leftIsDuration && !rightIsDuration {
		return Subtract(ctx, left, right), nil
	}

	if !leftIsDuration {
		var err error
		leftDuration, err = ToDuration(ctx, left)

		if err != nil {
			return None, err
		}
	}

	if !rightIsDuration {
		var err error
		rightDuration, err = ToDuration(ctx, right)

		if err != nil {
			return None, err
		}
	}

	result, ok := subtractDurationInt64(int64(leftDuration), int64(rightDuration))
	if !ok {
		return None, durationRangeError("subtraction")
	}

	return Duration(result), nil
}

// MultiplyChecked applies checked numeric and Duration multiplication.
func MultiplyChecked(ctx context.Context, left, right Value) (Value, error) {
	var leftDuration, rightDuration Duration
	leftIsDuration := false
	rightIsDuration := false

	switch value := left.(type) {
	case DateTime:
		return None, temporalBinaryTypeError("*", left, right)
	case Duration:
		leftDuration = value
		leftIsDuration = true
	}

	switch value := right.(type) {
	case DateTime:
		return None, temporalBinaryTypeError("*", left, right)
	case Duration:
		rightDuration = value
		rightIsDuration = true
	}

	if !leftIsDuration && !rightIsDuration {
		return Multiply(ctx, left, right), nil
	}

	if leftIsDuration && rightIsDuration {
		return None, durationBinaryTypeError("*", left, right)
	}

	if leftIsDuration {
		return scaleDuration(ctx, leftDuration, right, false, false)
	}

	return scaleDuration(ctx, rightDuration, left, false, false)
}

// DivideChecked applies checked numeric and Duration division.
func DivideChecked(ctx context.Context, left, right Value) (Value, error) {
	var leftDuration, rightDuration Duration
	leftIsDuration := false
	rightIsDuration := false

	switch value := left.(type) {
	case DateTime:
		return None, temporalBinaryTypeError("/", left, right)
	case Duration:
		leftDuration = value
		leftIsDuration = true
	}

	switch value := right.(type) {
	case DateTime:
		return None, temporalBinaryTypeError("/", left, right)
	case Duration:
		rightDuration = value
		rightIsDuration = true
	}

	if !leftIsDuration && !rightIsDuration {
		return Divide(ctx, left, right), nil
	}

	if !leftIsDuration {
		return None, durationBinaryTypeError("/", left, right)
	}

	if rightIsDuration {
		return divideDurationRatio(leftDuration, rightDuration)
	}

	if rightString, ok := right.(String); ok {
		duration, err := ToDuration(ctx, rightString)
		if err == nil {
			return divideDurationRatio(leftDuration, duration)
		}

		if errors.Is(err, ErrRange) {
			return None, err
		}
	}

	return scaleDuration(ctx, leftDuration, right, true, true)
}

// ModulusChecked rejects temporal modulus and preserves numeric behavior.
func ModulusChecked(ctx context.Context, left, right Value) (Value, error) {
	switch left.(type) {
	case DateTime:
		return None, temporalBinaryTypeError("%", left, right)
	case Duration:
		return None, durationBinaryTypeError("%", left, right)
	}

	switch right.(type) {
	case DateTime:
		return None, temporalBinaryTypeError("%", left, right)
	case Duration:
		return None, durationBinaryTypeError("%", left, right)
	}

	return Modulus(ctx, left, right), nil
}

// IncrementChecked rejects incrementing a duration.
func IncrementChecked(ctx context.Context, value Value) (Value, error) {
	if _, ok := value.(DateTime); ok {
		return None, Error(ErrInvalidOperation, "increment is not supported for DateTime")
	}

	if _, ok := value.(Duration); ok {
		return None, Error(ErrInvalidOperation, "increment is not supported for Duration")
	}

	return Increment(ctx, value), nil
}

// DecrementChecked rejects decrementing a duration.
func DecrementChecked(ctx context.Context, value Value) (Value, error) {
	if _, ok := value.(DateTime); ok {
		return None, Error(ErrInvalidOperation, "decrement is not supported for DateTime")
	}

	if _, ok := value.(Duration); ok {
		return None, Error(ErrInvalidOperation, "decrement is not supported for Duration")
	}

	return Decrement(ctx, value), nil
}

// NegateChecked applies unary negation while detecting duration overflow.
func NegateChecked(value Value) (Value, error) {
	if _, ok := value.(DateTime); ok {
		return None, Error(ErrInvalidOperation, "logical negation is not supported for DateTime")
	}

	duration, ok := value.(Duration)
	if !ok {
		switch value := value.(type) {
		case Int:
			return -value, nil
		case Float:
			return -value, nil
		case Boolean:
			return !value, nil
		default:
			return None, nil
		}
	}

	if duration == Duration(math.MinInt64) {
		return None, durationRangeError("negation")
	}

	return -duration, nil
}

// NegativeChecked applies unary minus while detecting duration overflow.
func NegativeChecked(value Value) (Value, error) {
	if _, ok := value.(DateTime); ok {
		return None, Error(ErrInvalidOperation, "unary minus is not supported for DateTime")
	}

	duration, ok := value.(Duration)
	if ok {
		if duration == Duration(math.MinInt64) {
			return None, durationRangeError("negation")
		}

		return -duration, nil
	}

	switch value := value.(type) {
	case Int:
		return -value, nil
	case Float:
		return -value, nil
	default:
		return None, nil
	}
}

// PositiveChecked applies unary positive to numeric and duration values.
func PositiveChecked(value Value) (Value, error) {
	switch value := value.(type) {
	case DateTime:
		return None, Error(ErrInvalidOperation, "unary plus is not supported for DateTime")
	case Duration:
		return value, nil
	case Int:
		return +value, nil
	case Float:
		return +value, nil
	default:
		return None, nil
	}
}

func scaleDuration(ctx context.Context, duration Duration, scalar Value, divide, allowNumericString bool) (Value, error) {
	number := scalar
	if _, ok := number.(String); ok && allowNumericString {
		converted, err := ToNumber(ctx, scalar)
		if err != nil {
			return None, Errorf(
				ErrInvalidArgument,
				"cannot use %s %q as a numeric Duration scale",
				TypeName(TypeOf(scalar)),
				scalar.String(),
			)
		}

		number = converted
	}

	switch value := number.(type) {
	case Int:
		if divide {
			if value == 0 {
				return None, Error(ErrInvalidOperation, "division by zero")
			}

			if duration == Duration(math.MinInt64) && value == -1 {
				return None, durationRangeError("division")
			}

			return Duration(int64(duration) / int64(value)), nil
		}

		result, ok := multiplyDurationInt64(int64(duration), int64(value))
		if !ok {
			return None, durationRangeError("multiplication")
		}

		return Duration(result), nil
	case Float:
		scalarFloat := float64(value)

		if math.IsNaN(scalarFloat) || math.IsInf(scalarFloat, 0) {
			return None, Error(ErrInvalidOperation, "duration scale must be finite")
		}

		if divide && scalarFloat == 0 {
			return None, Error(ErrInvalidOperation, "division by zero")
		}

		result := new(big.Rat).SetInt64(int64(duration))
		scalarRat := new(big.Rat).SetFloat64(scalarFloat)

		if divide {
			result.Quo(result, scalarRat)
		} else {
			result.Mul(result, scalarRat)
		}

		nanos, err := truncateDurationRat(result)
		if err != nil {
			return None, err
		}

		return Duration(nanos), nil
	default:
		operator := "*"
		if divide {
			operator = "/"
		}

		return None, durationBinaryTypeError(operator, duration, scalar)
	}
}

func divideDurationRatio(left, right Duration) (Value, error) {
	if right == 0 {
		return None, Error(ErrInvalidOperation, "division by zero")
	}

	leftNanos := int64(left)
	rightNanos := int64(right)

	if leftNanos == math.MinInt64 && rightNanos == -1 {
		return Float(float64(leftNanos) / float64(rightNanos)), nil
	}

	if leftNanos%rightNanos == 0 {
		return Int(leftNanos / rightNanos), nil
	}

	return Float(float64(leftNanos) / float64(rightNanos)), nil
}

func addDurationInt64(left, right int64) (int64, bool) {
	result := left + right
	if (right > 0 && result < left) || (right < 0 && result > left) {
		return 0, false
	}

	return result, true
}

func subtractDurationInt64(left, right int64) (int64, bool) {
	result := left - right
	if (right > 0 && result > left) || (right < 0 && result < left) {
		return 0, false
	}

	return result, true
}

func multiplyDurationInt64(left, right int64) (int64, bool) {
	if left == 0 || right == 0 {
		return 0, true
	}
	if (left == math.MinInt64 && right == -1) || (right == math.MinInt64 && left == -1) {
		return 0, false
	}

	result := left * right
	if result/right != left {
		return 0, false
	}

	return result, true
}

func durationBinaryTypeError(operator string, left, right Value) error {
	return Errorf(
		ErrInvalidOperation,
		"operator %s is not supported for %s and %s",
		operator,
		TypeName(TypeOf(left)),
		TypeName(TypeOf(right)),
	)
}

func durationRangeError(operation string) error {
	return Errorf(ErrRange, "duration %s exceeds the supported range", operation)
}
