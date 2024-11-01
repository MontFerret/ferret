package values

import (
	"encoding/binary"
	"hash/fnv"
	"strconv"

	"github.com/MontFerret/ferret/pkg/runtime/values/types"

	"github.com/wI2L/jettison"

	"github.com/MontFerret/ferret/pkg/runtime/core"
)

type Int int64

const ZeroInt = Int(0)

func NewInt(input int) Int {
	return Int(int64(input))
}

func NewInt64(input int64) Int {
	return Int(input)
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
	switch otherVal := other.(type) {
	case Int:
		if t == otherVal {
			return 0
		}

		if t < otherVal {
			return -1
		}

		return +1

	case Float:
		f := Float(t)

		if f == otherVal {
			return 0
		}

		if f < otherVal {
			return -1
		}

		return +1

	default:
		return types.Compare(types.Int, core.Reflect(other))
	}
}

func (t Int) Unwrap() interface{} {
	return int(t)
}

func (t Int) Hash() uint64 {
	h := fnv.New64a()

	h.Write([]byte(types.Int.String()))
	h.Write([]byte(":"))

	bytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(bytes, uint64(t))
	h.Write(bytes)

	return h.Sum64()
}

func (t Int) Copy() core.Value {
	return t
}
