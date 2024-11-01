package values

import (
	"context"
	"encoding/binary"
	"hash/fnv"
	"sort"

	"github.com/MontFerret/ferret/pkg/runtime/values/types"

	"github.com/wI2L/jettison"

	"github.com/MontFerret/ferret/pkg/runtime/core"
)

type (
	ArrayPredicate = func(value core.Value, idx int) bool

	ArraySorter = func(first, second core.Value) bool

	Array struct {
		data []core.Value
	}
)

func EmptyArray() *Array {
	return &Array{data: make([]core.Value, 0, 0)}
}

func NewArray(size int) *Array {
	return &Array{data: make([]core.Value, 0, size)}
}

func NewSizedArray(size int) *Array {
	return &Array{data: make([]core.Value, size)}
}

func NewArrayWith(values ...core.Value) *Array {
	return &Array{data: values}
}

func (t *Array) Iterate(ctx context.Context) (core.Iterator, error) {
	return NewArrayIterator(t), nil
}

func (t *Array) MarshalJSON() ([]byte, error) {
	return jettison.MarshalOpts(t.data, jettison.NoHTMLEscaping())
}

func (t *Array) Type() core.Type {
	return types.Array
}

func (t *Array) String() string {
	marshaled, err := t.MarshalJSON()

	if err != nil {
		return "[]"
	}

	return string(marshaled)
}

func (t *Array) Compare(other core.Value) int64 {
	otherArr, ok := other.(*Array)

	if !ok {
		return types.Compare(types.Array, core.Reflect(other))
	}

	if t.Length() == 0 && otherArr.Length() == 0 {
		return 0
	}

	if t.Length() < otherArr.Length() {
		return -1
	}

	if t.Length() > otherArr.Length() {
		return 1
	}

	var res int64
	var val core.Value

	otherArr.ForEach(func(otherVal core.Value, idx int) bool {
		val = t.Get(idx)
		res = Compare(val, otherVal)

		return res == 0
	})

	return res
}

func (t *Array) Unwrap() interface{} {
	arr := make([]interface{}, t.Length())

	for idx, val := range t.data {
		arr[idx] = val.Unwrap()
	}

	return arr
}

func (t *Array) Hash() uint64 {
	h := fnv.New64a()

	h.Write([]byte(types.Array.String()))
	h.Write([]byte(":"))
	h.Write([]byte("["))

	endIndex := len(t.data) - 1

	for i, el := range t.data {
		bytes := make([]byte, 8)
		binary.LittleEndian.PutUint64(bytes, el.Hash())

		h.Write(bytes)

		if i != endIndex {
			h.Write([]byte(","))
		}
	}

	h.Write([]byte("]"))

	return h.Sum64()
}

func (t *Array) Copy() core.Value {
	c := NewArray(len(t.data))

	for _, el := range t.data {
		c.Push(el)
	}

	return c
}

func (t *Array) Length() int {
	return len(t.data)
}

func (t *Array) ForEach(predicate ArrayPredicate) {
	for idx, val := range t.data {
		if !predicate(val, idx) {
			break
		}
	}
}

func (t *Array) First() core.Value {
	if len(t.data) > 0 {
		return t.data[0]
	}

	return None
}

func (t *Array) Last() core.Value {
	size := len(t.data)

	if size > 1 {
		return t.data[size-1]
	} else if size == 1 {
		return t.data[0]
	}

	return None
}

func (t *Array) Find(predicate ArrayPredicate) (*Array, Boolean) {
	result := NewArray(len(t.data))

	for idx, val := range t.data {
		if predicate(val, idx) {
			result.Push(val)
		}
	}

	return result, result.Length() > 0
}

func (t *Array) FindOne(predicate ArrayPredicate) (core.Value, Boolean) {
	for idx, val := range t.data {
		if predicate(val, idx) {
			return val, True
		}
	}

	return None, False
}

func (t *Array) Get(idx int) core.Value {
	l := len(t.data) - 1

	if l < 0 {
		return None
	}

	if int(idx) > l {
		return None
	}

	return t.data[idx]
}

func (t *Array) Set(idx Int, value core.Value) error {
	last := len(t.data) - 1

	if last >= int(idx) {
		t.data[idx] = value

		return nil
	}

	return core.Error(core.ErrInvalidOperation, "out of bounds")
}

func (t *Array) MustSet(idx Int, value core.Value) {
	t.data[idx] = value
}

func (t *Array) Push(item core.Value) {
	t.data = append(t.data, item)
}

func (t *Array) Slice(from, to int) *Array {
	length := t.Length()

	if from >= length {
		return NewArray(0)
	}

	if to > length {
		to = length
	}

	result := new(Array)
	result.data = t.data[from:to]

	return result
}

func (t *Array) Contains(item core.Value) Boolean {
	return t.IndexOf(item) >= 0
}

func (t *Array) IndexOf(item core.Value) Int {
	res := Int(-1)

	for idx, el := range t.data {
		if Compare(item, el) == 0 {
			res = Int(idx)
			break
		}
	}

	return res
}

func (t *Array) Insert(idx Int, value core.Value) {
	t.data = append(t.data[:idx], append([]core.Value{value}, t.data[idx:]...)...)
}

func (t *Array) RemoveAt(idx Int) {
	i := int(idx)
	max := len(t.data) - 1

	if i > max {
		return
	}

	t.data = append(t.data[:i], t.data[i+1:]...)
}

func (t *Array) Clone() core.Cloneable {
	cloned := NewArray(0)

	var value core.Value
	for idx := 0; idx < t.Length(); idx++ {
		value = t.Get(idx)

		cloneable, ok := value.(core.Cloneable)

		if ok {
			value = cloneable.Clone()
		}

		cloned.Push(value)
	}

	return cloned
}

func (t *Array) Sort() *Array {
	return t.SortWith(func(first, second core.Value) bool {
		return Compare(first, second) == -1
	})
}

func (t *Array) SortWith(sorter ArraySorter) *Array {
	c := make([]core.Value, len(t.data))
	copy(c, t.data)

	sort.SliceStable(c, func(i, j int) bool {
		return sorter(c[i], c[j])
	})

	res := new(Array)
	res.data = c

	return res
}
