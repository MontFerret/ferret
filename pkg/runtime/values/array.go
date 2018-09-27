package values

import (
	"crypto/sha512"
	"encoding/json"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/pkg/errors"
)

type (
	ArrayPredicate = func(value core.Value, idx int) bool
	Array          struct {
		value []core.Value
	}
)

func NewArray(size int) *Array {
	return &Array{value: make([]core.Value, 0, size)}
}

func NewArrayWith(values ...core.Value) *Array {
	return &Array{value: values}
}

func (t *Array) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.value)
}

func (t *Array) Type() core.Type {
	return core.ArrayType
}

func (t *Array) String() string {
	marshaled, err := t.MarshalJSON()

	if err != nil {
		return "[]"
	}

	return string(marshaled)
}

func (t *Array) Compare(other core.Value) int {
	switch other.Type() {
	case core.ArrayType:
		arr := other.(*Array)

		if t.Length() == 0 && arr.Length() == 0 {
			return 0
		}

		var res = 1

		for _, val := range t.value {
			arr.ForEach(func(otherVal core.Value, idx int) bool {
				res = val.Compare(otherVal)

				return res != -1
			})
		}

		return res
	case core.ObjectType:
		return -1
	default:
		return 1
	}
}

func (t *Array) Unwrap() interface{} {
	arr := make([]interface{}, t.Length())

	for idx, val := range t.value {
		arr[idx] = val.Unwrap()
	}

	return arr
}

func (t *Array) Hash() int {
	bytes, err := t.MarshalJSON()

	if err != nil {
		return 0
	}

	h := sha512.New()

	out, err := h.Write(bytes)

	if err != nil {
		return 0
	}

	return out
}

func (t *Array) Length() Int {
	return Int(len(t.value))
}

func (t *Array) ForEach(predicate ArrayPredicate) {
	for idx, val := range t.value {
		if predicate(val, idx) == false {
			break
		}
	}
}

func (t *Array) Get(idx Int) core.Value {
	l := len(t.value) - 1

	if int(idx) > l {
		return None
	}

	return t.value[idx]
}

func (t *Array) Set(idx Int, value core.Value) error {
	last := len(t.value) - 1

	if last >= int(idx) {
		t.value[idx] = value

		return nil
	}

	return errors.Wrap(core.ErrInvalidOperation, "out of bounds")
}

func (t *Array) Push(item core.Value) {
	t.value = append(t.value, item)
}

func (t *Array) Slice(from, to Int) []core.Value {
	return t.value[from:to]
}

func (t *Array) IndexOf(item core.Value) Int {
	res := Int(-1)

	for idx, el := range t.value {
		if el.Compare(item) == 0 {
			res = Int(idx)
			break
		}
	}

	return res
}

func (t *Array) Insert(idx Int, value core.Value) {
	t.value = append(t.value[:idx], append([]core.Value{value}, t.value[idx:]...)...)
}

func (t *Array) RemoveAt(idx Int) {
	i := int(idx)
	max := len(t.value) - 1

	if i > max {
		return
	}

	t.value = append(t.value[:i], t.value[i+1:]...)
}
