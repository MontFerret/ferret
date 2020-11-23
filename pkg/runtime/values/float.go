package values

import (
	"encoding/binary"
	"fmt"
	"hash/fnv"
	"math"
	"strconv"

	"github.com/wI2L/jettison"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values/types"
)

type Float float64

var NaN = Float(math.NaN())

const ZeroFloat = Float(0.0)

func NewFloat(input float64) Float {
	return Float(input)
}

func ParseFloat(input interface{}) (Float, error) {
	if core.IsNil(input) {
		return ZeroFloat, nil
	}

	i, ok := input.(float64)

	if ok {
		if i == 0 {
			return ZeroFloat, nil
		}

		return Float(i), nil
	}

	// try to cast
	str, ok := input.(string)

	if ok {
		i, err := strconv.Atoi(str)

		if err == nil {
			if i == 0 {
				return ZeroFloat, nil
			}

			return Float(i), nil
		}
	}

	return ZeroFloat, core.Error(core.ErrInvalidType, "expected 'float'")
}

func MustParseFloat(input interface{}) Float {
	res, err := ParseFloat(input)

	if err != nil {
		panic(err)
	}

	return res
}

func IsNaN(input Float) Boolean {
	return NewBoolean(math.IsNaN(float64(input)))
}

func IsInf(input Float, sign Int) Boolean {
	return NewBoolean(math.IsInf(float64(input), int(sign)))
}

func (t Float) MarshalJSON() ([]byte, error) {
	return jettison.MarshalOpts(float64(t), jettison.NoHTMLEscaping())
}

func (t Float) Type() core.Type {
	return types.Float
}

func (t Float) String() string {
	return fmt.Sprintf("%v", float64(t))
}

func (t Float) Compare(other core.Value) int64 {
	otherType := other.Type()
	raw := float64(t)

	if otherType == types.Float {
		f := other.Unwrap().(float64)

		if raw == f {
			return 0
		}

		if raw < f {
			return -1
		}

		return +1
	}

	if otherType == types.Int {
		i := other.Unwrap().(int)
		f := float64(i)

		if raw == f {
			return 0
		}

		if raw < f {
			return -1
		}

		return +1
	}

	return types.Compare(types.Float, otherType)
}

func (t Float) Unwrap() interface{} {
	return float64(t)
}

func (t Float) Hash() uint64 {
	h := fnv.New64a()

	h.Write([]byte(t.Type().String()))
	h.Write([]byte(":"))

	bytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(bytes, math.Float64bits(float64(t)))
	h.Write(bytes)

	return h.Sum64()
}

func (t Float) Copy() core.Value {
	return t
}
