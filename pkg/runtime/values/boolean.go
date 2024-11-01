package values

import (
	"hash/fnv"
	"strings"

	"github.com/MontFerret/ferret/pkg/runtime/values/types"

	"github.com/wI2L/jettison"

	"github.com/MontFerret/ferret/pkg/runtime/core"
)

type Boolean bool

const (
	False = Boolean(false)
	True  = Boolean(true)
)

func NewBoolean(input bool) Boolean {
	return Boolean(input)
}

func ParseBoolean(input interface{}) (Boolean, error) {
	b, ok := input.(bool)

	if ok {
		if b {
			return True, nil
		}

		return False, nil
	}

	s, ok := input.(string)

	if ok {
		return Boolean(strings.ToLower(s) == "true"), nil
	}

	return False, core.Error(core.ErrInvalidType, "expected 'bool'")
}

func MustParseBoolean(input interface{}) Boolean {
	res, err := ParseBoolean(input)

	if err != nil {
		panic(err)
	}

	return res
}

func (t Boolean) MarshalJSON() ([]byte, error) {
	return jettison.MarshalOpts(bool(t), jettison.NoHTMLEscaping())
}

func (t Boolean) Type() core.Type {
	return types.Boolean
}

func (t Boolean) String() string {
	if t {
		return "true"
	}

	return "false"
}

func (t Boolean) Compare(other core.Value) int64 {
	otherBool, ok := other.(Boolean)

	if !ok {
		return types.Compare(types.Boolean, core.Reflect(other))
	}

	if t == otherBool {
		return 0
	}

	if !t && otherBool {
		return -1
	}

	return +1
}

func (t Boolean) Unwrap() interface{} {
	return bool(t)
}

func (t Boolean) Hash() uint64 {
	h := fnv.New64a()

	h.Write([]byte(types.Boolean.String()))
	h.Write([]byte(":"))
	h.Write([]byte(t.String()))

	return h.Sum64()
}

func (t Boolean) Copy() core.Value {
	return t
}
