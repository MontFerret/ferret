package values

import (
	"hash/fnv"
	"time"

	"github.com/MontFerret/ferret/pkg/runtime/values/types"

	"github.com/wI2L/jettison"

	"github.com/MontFerret/ferret/pkg/runtime/core"
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
		return DateTime{time.Now()}, core.ErrInvalidType
	}
}

func MustParseDateTime(input interface{}) DateTime {
	dt, err := ParseDateTime(input)

	if err != nil {
		panic(err)
	}

	return dt
}

func (t DateTime) MarshalJSON() ([]byte, error) {
	return jettison.MarshalOpts(t.Time, jettison.NoHTMLEscaping())
}

func (t DateTime) Type() core.Type {
	return types.DateTime
}

func (t DateTime) String() string {
	return t.Time.String()
}

func (t DateTime) Compare(other core.Value) int64 {
	otherDt, ok := other.(DateTime)

	if !ok {
		return types.Compare(types.DateTime, core.Reflect(other))
	}

	if t.After(otherDt.Time) {
		return 1
	}

	if t.Before(otherDt.Time) {
		return -1
	}

	return 0
}

func (t DateTime) Unwrap() interface{} {
	return t.Time
}

func (t DateTime) Hash() uint64 {
	h := fnv.New64a()

	h.Write([]byte(types.DateTime.String()))
	h.Write([]byte(":"))

	bytes, err := t.Time.GobEncode()

	if err != nil {
		return 0
	}

	h.Write(bytes)

	return h.Sum64()
}

func (t DateTime) Copy() core.Value {
	return NewDateTime(t.Time)
}
