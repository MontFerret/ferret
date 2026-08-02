package runtime

import (
	"context"
	"hash/fnv"
	"time"
)

const DefaultTimeLayout = time.RFC3339

type DateTime struct {
	time.Time
}

var ZeroDateTime = DateTime{
	time.Time{},
}

func NewCurrentDateTime() DateTime {
	return DateTime{time.Now()}
}

func NewDateTime(time time.Time) DateTime {
	return DateTime{time}
}

func ParseDateTime(input interface{}) (DateTime, error) {
	return ParseDateTimeWith(input, DefaultTimeLayout)
}

func ParseDateTimeWith(input interface{}, layout string) (DateTime, error) {
	switch value := input.(type) {
	case string:
		t, err := time.Parse(layout, value)

		if err != nil {
			return DateTime{time.Now()}, err
		}

		return DateTime{t}, nil
	default:
		return DateTime{time.Now()}, ErrInvalidType
	}
}

// ToDateTime converts a native DateTime or RFC3339 runtime String to DateTime.
func ToDateTime(_ context.Context, input Value) (DateTime, error) {
	if input == nil {
		input = None
	}

	switch value := input.(type) {
	case DateTime:
		return value, nil
	case String:
		parsed, err := time.Parse(DefaultTimeLayout, value.String())
		if err != nil {
			return ZeroDateTime, newConversionError(TypeDateTime, Errorf(
				ErrInvalidArgument,
				"cannot convert String %q to DateTime",
				value.String(),
			))
		}

		return NewDateTime(parsed), nil
	default:
		return ZeroDateTime, newConversionError(TypeDateTime, Errorf(
			ErrInvalidType,
			"cannot convert %s %q to DateTime",
			TypeName(TypeOf(input)),
			input.String(),
		))
	}
}

func MustParseDateTime(input interface{}) DateTime {
	dt, err := ParseDateTime(input)

	if err != nil {
		panic(err)
	}

	return dt
}

func (dt DateTime) Type() Type {
	return TypeDateTime
}

func (dt DateTime) String() string {
	return dt.Time.String()
}

func (dt DateTime) Unwrap() any {
	return dt.Time
}

func (dt DateTime) Hash() uint64 {
	h := fnv.New64a()

	h.Write([]byte(TypeDateTime.Name()))
	h.Write([]byte(":"))

	bytes, err := dt.GobEncode()

	if err != nil {
		return 0
	}

	h.Write(bytes)

	return h.Sum64()
}

func (dt DateTime) Copy() Value {
	return NewDateTime(dt.Time)
}

func (dt DateTime) Compare(other Value) int {
	otherDt, ok := other.(DateTime)

	if !ok {
		return CompareTypes(dt, other)
	}

	if dt.After(otherDt.Time) {
		return 1
	}

	if dt.Before(otherDt.Time) {
		return -1
	}

	return 0
}
