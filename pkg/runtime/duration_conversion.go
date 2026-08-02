package runtime

import (
	"context"
	"errors"
	"math"
	"math/big"
	"time"
)

const durationNumericUnit = int64(time.Millisecond)

// ToDuration converts a runtime value to a canonical Duration.
// Numeric values are interpreted as milliseconds, and fractional nanoseconds
// are truncated toward zero. Duration strings, NONE, Booleans, and supported
// list shapes follow the language's broader Duration conversion contract.
func ToDuration(ctx context.Context, input Value) (Duration, error) {
	if input == nil || input == None {
		return ZeroDuration, nil
	}

	switch value := input.(type) {
	case Duration:
		return value, nil
	case Int:
		result, ok := multiplyDurationInt64(int64(value), durationNumericUnit)
		if !ok {
			return ZeroDuration, durationConversionRangeError(input)
		}

		return Duration(result), nil
	case Float:
		inputFloat := float64(value)
		if math.IsNaN(inputFloat) || math.IsInf(inputFloat, 0) {
			return ZeroDuration, durationConversionArgumentError(input)
		}

		result := new(big.Rat).SetFloat64(inputFloat)
		result.Mul(result, new(big.Rat).SetInt64(durationNumericUnit))
		nanos, err := truncateDurationRat(result)
		if err != nil {
			return ZeroDuration, durationConversionRangeError(input)
		}

		return Duration(nanos), nil
	case String:
		duration, err := ParseDuration(value.String())
		if err == nil {
			return duration, nil
		}

		if errors.Is(err, ErrRange) {
			return ZeroDuration, durationConversionRangeError(input)
		}

		return ZeroDuration, durationConversionArgumentError(input)
	case Boolean:
		if value {
			return Duration(durationNumericUnit), nil
		}

		return ZeroDuration, nil
	case List:
		length, err := value.Length(ctx)
		if err != nil {
			return ZeroDuration, err
		}

		switch length {
		case 0:
			return ZeroDuration, nil
		case 1:
			item, err := value.At(ctx, ZeroInt)
			if err != nil {
				return ZeroDuration, err
			}

			return ToDuration(ctx, item)
		default:
			return ZeroDuration, durationConversionArgumentError(input)
		}
	default:
		return ZeroDuration, durationConversionTypeError(input)
	}
}

func durationConversionArgumentError(input Value) error {
	return newConversionError(TypeDuration, Errorf(
		ErrInvalidArgument,
		"cannot convert %s %q to Duration",
		TypeName(TypeOf(input)),
		input.String(),
	))
}

func durationConversionTypeError(input Value) error {
	return newConversionError(TypeDuration, Errorf(
		ErrInvalidType,
		"cannot convert %s %q to Duration",
		TypeName(TypeOf(input)),
		input.String(),
	))
}

func durationConversionRangeError(input Value) error {
	return newConversionError(TypeDuration, Errorf(
		ErrRange,
		"cannot convert %s %q to Duration: value exceeds the supported range",
		TypeName(TypeOf(input)),
		input.String(),
	))
}
