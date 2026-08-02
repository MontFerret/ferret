package runtime

import (
	"context"
	"math"
	"math/big"
	"strings"
	"time"
)

const (
	dateTimeEpochCanonicalUnits = `"s", "ms", "us", or "ns"`
	dateTimeUnixToInternal      = int64(62_135_596_800)
	dateTimeNanosPerSecond      = int64(time.Second)
)

// ToDateTimeEpoch converts a numeric Unix epoch offset with an explicit unit
// to a UTC DateTime. Supported units are seconds, milliseconds, microseconds,
// and nanoseconds. Aliases include sec, second, seconds, millisecond,
// milliseconds, µs, μs, microsecond, microseconds, nanosecond, and nanoseconds.
// Fractional nanoseconds are truncated toward zero.
func ToDateTimeEpoch(_ context.Context, input, unit Value) (DateTime, error) {
	if input == nil {
		input = None
	}

	if unit == nil {
		unit = None
	}

	unitString, ok := unit.(String)
	if !ok {
		return ZeroDateTime, dateTimeEpochUnitTypeError(unit)
	}

	var epoch *big.Rat
	switch value := input.(type) {
	case Int:
		epoch = new(big.Rat).SetInt64(int64(value))
	case Float:
		valueFloat := float64(value)
		if math.IsNaN(valueFloat) || math.IsInf(valueFloat, 0) {
			return ZeroDateTime, dateTimeEpochArgumentError(
				input,
				unitString,
				"numeric epoch value must be finite",
			)
		}

		epoch = new(big.Rat).SetFloat64(valueFloat)
	case DateTime, String:
		return ZeroDateTime, dateTimeEpochArgumentError(
			input,
			unitString,
			"epoch units are only valid for Int or Float inputs",
		)
	default:
		return ZeroDateTime, dateTimeEpochInputTypeError(input, unitString)
	}

	multiplier, ok := parseDateTimeEpochUnit(unitString)
	if !ok {
		return ZeroDateTime, dateTimeEpochArgumentError(
			input,
			unitString,
			"unsupported epoch unit; expected "+dateTimeEpochCanonicalUnits,
		)
	}

	totalNanos := new(big.Rat).Mul(epoch, new(big.Rat).SetInt64(multiplier))
	truncatedNanos := new(big.Int).Quo(totalNanos.Num(), totalNanos.Denom())
	seconds, nanos := new(big.Int), new(big.Int)
	seconds.QuoRem(truncatedNanos, big.NewInt(dateTimeNanosPerSecond), nanos)

	if nanos.Sign() < 0 {
		seconds.Sub(seconds, big.NewInt(1))
		nanos.Add(nanos, big.NewInt(dateTimeNanosPerSecond))
	}

	if !seconds.IsInt64() {
		return ZeroDateTime, dateTimeEpochRangeError(input, unitString)
	}

	secondsValue := seconds.Int64()
	nanosValue := nanos.Int64()
	if secondsValue > math.MaxInt64-dateTimeUnixToInternal {
		return ZeroDateTime, dateTimeEpochRangeError(input, unitString)
	}

	result := time.Unix(secondsValue, nanosValue).UTC()
	if result.Unix() != secondsValue || int64(result.Nanosecond()) != nanosValue {
		return ZeroDateTime, dateTimeEpochRangeError(input, unitString)
	}

	return NewDateTime(result), nil
}

func parseDateTimeEpochUnit(unit String) (int64, bool) {
	switch strings.ToLower(unit.String()) {
	case "s", "sec", "second", "seconds":
		return int64(time.Second), true
	case "ms", "millisecond", "milliseconds":
		return int64(time.Millisecond), true
	case "us", "µs", "μs", "microsecond", "microseconds":
		return int64(time.Microsecond), true
	case "ns", "nanosecond", "nanoseconds":
		return int64(time.Nanosecond), true
	default:
		return 0, false
	}
}

func dateTimeEpochUnitTypeError(unit Value) error {
	return newConversionError(TypeDateTime, Errorf(
		ErrInvalidType,
		"cannot use %s %q as a DateTime epoch unit: expected String containing %s",
		TypeName(TypeOf(unit)),
		unit.String(),
		dateTimeEpochCanonicalUnits,
	))
}

func dateTimeEpochInputTypeError(input Value, unit String) error {
	return newConversionError(TypeDateTime, Errorf(
		ErrInvalidType,
		"cannot convert %s %q to DateTime with epoch unit %q: expected Int or Float",
		TypeName(TypeOf(input)),
		input.String(),
		unit.String(),
	))
}

func dateTimeEpochArgumentError(input Value, unit String, detail string) error {
	return newConversionError(TypeDateTime, Errorf(
		ErrInvalidArgument,
		"cannot convert %s %q to DateTime with epoch unit %q: %s",
		TypeName(TypeOf(input)),
		input.String(),
		unit.String(),
		detail,
	))
}

func dateTimeEpochRangeError(input Value, unit String) error {
	return newConversionError(TypeDateTime, Errorf(
		ErrRange,
		"cannot convert %s %q to DateTime with epoch unit %q: value exceeds the supported DateTime range",
		TypeName(TypeOf(input)),
		input.String(),
		unit.String(),
	))
}
