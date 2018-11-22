package values

import (
	"encoding/binary"
	"encoding/json"
	"hash/fnv"
	"strconv"

	"github.com/MontFerret/ferret/pkg/runtime/core"
)

type Int int64

var ZeroInt = Int(0)

func NewInt(input int) Int {
	return Int(int64(input))
}

func ParseInt(input interface{}) (Int, error) {
	if core.IsNil(input) {
		return ZeroInt, nil
	}

	switch input.(type) {
	case int:
		return Int(input.(int)), nil
	case int64:
		return Int(input.(int64)), nil
	case int32:
		return Int(input.(int32)), nil
	case int16:
		return Int(input.(int16)), nil
	case int8:
		return Int(input.(int8)), nil
	case string:
		i, err := strconv.Atoi(input.(string))

		if err == nil {
			if i == 0 {
				return ZeroInt, nil
			}

			return Int(i), nil
		}

		return ZeroInt, err
	default:
		return ZeroInt, core.Error(core.ErrInvalidType, "expected 'int'")
	}
}

func ParseIntP(input interface{}) Int {
	res, err := ParseInt(input)

	if err != nil {
		panic(err)
	}

	return res
}

func (t Int) MarshalJSON() ([]byte, error) {
	return json.Marshal(int64(t))
}

func (t Int) Type() core.Type {
	return core.IntType
}

func (t Int) String() string {
	return strconv.Itoa(int(t))
}

func (t Int) Compare(other core.Value) int {
	switch other.Type() {
	case core.IntType:
		i := other.(Int)

		if t == i {
			return 0
		}

		if t < i {
			return -1
		}

		return +1
	case core.FloatType:
		f := other.(Float)
		f2 := Float(t)

		if f2 == f {
			return 0
		}

		if f2 < f {
			return -1
		}

		return +1
	case core.BooleanType, core.NoneType:
		return 1
	default:
		return -1
	}
}

func (t Int) Unwrap() interface{} {
	return int(t)
}

func (t Int) Hash() uint64 {
	h := fnv.New64a()

	h.Write([]byte(t.Type().String()))
	h.Write([]byte(":"))

	bytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(bytes, uint64(t))
	h.Write(bytes)

	return h.Sum64()
}

func (t Int) Copy() core.Value {
	return t
}
