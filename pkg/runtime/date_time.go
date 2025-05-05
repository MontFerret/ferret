package runtime

import (
	"hash/fnv"
	"time"

	"github.com/wI2L/jettison"
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

func MustParseDateTime(input interface{}) DateTime {
	dt, err := ParseDateTime(input)

	if err != nil {
		panic(err)
	}

	return dt
}

func (dt DateTime) MarshalJSON() ([]byte, error) {
	return jettison.MarshalOpts(dt.Time, jettison.NoHTMLEscaping())
}

func (dt DateTime) Type() string {
	return TypeDateTime
}

func (dt DateTime) String() string {
	return dt.Time.String()
}

func (dt DateTime) Unwrap() interface{} {
	return dt.Time
}

func (dt DateTime) Hash() uint64 {
	h := fnv.New64a()

	h.Write([]byte(TypeDateTime))
	h.Write([]byte(":"))

	bytes, err := dt.Time.GobEncode()

	if err != nil {
		return 0
	}

	h.Write(bytes)

	return h.Sum64()
}

func (dt DateTime) Copy() Value {
	return NewDateTime(dt.Time)
}

func (dt DateTime) Compare(other Value) int64 {
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
