package values

import (
	"encoding/binary"
	"hash/fnv"
	"strconv"

	"github.com/wI2L/jettison"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values/types"
)

type Int int64

const ZeroInt = Int(0)

func NewInt(input int) Int {
	return Int(int64(input))
}

func ParseInt(input interface{}) (Int, error) {
	if core.IsNil(input) {
		return ZeroInt, nil
	}

	switch val := input.(type) {
	case int:
		return Int(val), nil
	case int64:
		return Int(val), nil
	case int32:
		return Int(val), nil
	case int16:
		return Int(val), nil
	case int8:
		return Int(val), nil
	case string:
		i, err := strconv.Atoi(val)

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

func MustParseInt(input interface{}) Int {
	res, err := ParseInt(input)

	if err != nil {
		panic(err)
	}

	return res
}

func (t Int) MarshalJSON() ([]byte, error) {
	return jettison.MarshalOpts(int64(t), jettison.NoHTMLEscaping())
}

func (t Int) Type() core.Type {
	return types.Int
}

func (t Int) String() string {
	return strconv.Itoa(int(t))
}

func (t Int) Compare(other core.Value) int64 {
	otherType := other.Type()

	if otherType == types.Int {
		i := other.(Int)

		if t == i {
			return 0
		}

		if t < i {
			return -1
		}

		return +1
	}

	if otherType == types.Float {
		f := other.(Float)
		f2 := Float(t)

		if f2 == f {
			return 0
		}

		if f2 < f {
			return -1
		}

		return +1
	}

	return types.Compare(types.Int, otherType)
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
