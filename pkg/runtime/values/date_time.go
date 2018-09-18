package values

import (
	"crypto/sha512"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"time"
)

const defaultTimeLayout = time.RFC3339

type DateTime struct {
	time.Time
}

func NewDateTime(time time.Time) DateTime {
	return DateTime{time}
}

func ParseDateTime(input interface{}) (DateTime, error) {
	return ParseDateTimeWith(input, defaultTimeLayout)
}

func ParseDateTimeWith(input interface{}, layout string) (DateTime, error) {
	switch input.(type) {
	case string:
		t, err := time.Parse(layout, input.(string))

		if err != nil {
			return DateTime{time.Now()}, err
		}

		return DateTime{t}, nil
	default:
		return DateTime{time.Now()}, core.ErrInvalidType
	}
}

func ParseDateTimeP(input interface{}) DateTime {
	dt, err := ParseDateTime(input)

	if err != nil {
		panic(err)
	}

	return dt
}

func (t DateTime) MarshalJSON() ([]byte, error) {
	return t.Time.MarshalJSON()
}

func (t DateTime) Type() core.Type {
	return core.DateTimeType
}

func (t DateTime) String() string {
	return t.Time.String()
}

func (t DateTime) Compare(other core.Value) int {
	switch other.Type() {
	case core.DateTimeType:
		other := other.(DateTime)

		if t.After(other.Time) {
			return 1
		}

		if t.Before(other.Time) {
			return -1
		}

		return 0
	default:
		if other.Type() > core.DateTimeType {
			return -1
		}

		return 1
	}
}

func (t DateTime) Unwrap() interface{} {
	return t.Time
}

func (t DateTime) Hash() int {
	h := sha512.New()

	t.Time.MarshalJSON()

	out, err := h.Write([]byte(t.Time.String()))

	if err != nil {
		return 0
	}

	return out
}
