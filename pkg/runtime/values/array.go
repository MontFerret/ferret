package values

import (
	"context"
	"encoding/binary"
	"hash/fnv"
	"sort"

	"github.com/wI2L/jettison"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values/types"
)

type (
	ArrayPredicate = func(value core.Value, idx int) bool

	ArraySorter = func(first, second core.Value) bool

	Array struct {
		items []core.Value
	}
)

func EmptyArray() *Array {
	return &Array{items: make([]core.Value, 0, 0)}
}

func NewArray(size int) *Array {
	return &Array{items: make([]core.Value, 0, size)}
}

func NewArrayWith(values ...core.Value) *Array {
	return &Array{items: values}
}

func NewArrayOf(values []core.Value) *Array {
	return &Array{items: values}
}

func (t *Array) MarshalJSON() ([]byte, error) {
	return jettison.MarshalOpts(t.items, jettison.NoHTMLEscaping())
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
	if other.Type() == types.Array {
		other := other.(*Array)

		if t.Length() == 0 && other.Length() == 0 {
			return 0
		}

		if t.Length() < other.Length() {
			return -1
		}

		if t.Length() > other.Length() {
			return 1
		}

		var res int64
		var val core.Value

		other.ForEach(func(otherVal core.Value, idx int) bool {
			val = t.Get(NewInt(idx))
			res = val.Compare(otherVal)

			return res == 0
		})

		return res
	}

	return types.Compare(types.Array, other.Type())
}

func (t *Array) Unwrap() interface{} {
	arr := make([]interface{}, t.Length())

	for idx, val := range t.items {
		arr[idx] = val.Unwrap()
	}

	return arr
}

func (t *Array) Hash() uint64 {
	h := fnv.New64a()

	h.Write([]byte(t.Type().String()))
	h.Write([]byte(":"))
	h.Write([]byte("["))

	endIndex := len(t.items) - 1

	for i, el := range t.items {
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
	c := NewArray(len(t.items))

	for _, el := range t.items {
		c.Push(el)
	}

	return c
}

func (t *Array) Length() Int {
	return Int(len(t.items))
}

func (t *Array) ForEach(predicate ArrayPredicate) {
	for idx, val := range t.items {
		if !predicate(val, idx) {
			break
		}
	}
}

func (t *Array) First() core.Value {
	if len(t.items) > 0 {
		return t.items[0]
	}

	return None
}

func (t *Array) Last() core.Value {
	size := len(t.items)

	if size > 1 {
		return t.items[size-1]
	} else if size == 1 {
		return t.items[0]
	}

	return None
}

func (t *Array) Find(predicate ArrayPredicate) (*Array, Boolean) {
	result := NewArray(len(t.items))

	for idx, val := range t.items {
		if predicate(val, idx) {
			result.Push(val)
		}
	}

	return result, result.Length() > 0
}

func (t *Array) FindOne(predicate ArrayPredicate) (core.Value, Boolean) {
	for idx, val := range t.items {
		if predicate(val, idx) {
			return val, True
		}
	}

	return None, False
}

func (t *Array) Get(idx Int) core.Value {
	l := len(t.items) - 1

	if l < 0 {
		return None
	}

	if int(idx) > l {
		return None
	}

	return t.items[idx]
}

func (t *Array) Set(idx Int, value core.Value) error {
	last := len(t.items) - 1

	if last >= int(idx) {
		t.items[idx] = value

		return nil
	}

	return core.Error(core.ErrInvalidOperation, "out of bounds")
}

func (t *Array) Push(item core.Value) {
	t.items = append(t.items, item)
}

func (t *Array) Slice(from, to Int) *Array {
	length := t.Length()

	if from >= length {
		return NewArray(0)
	}

	if to > length {
		to = length
	}

	result := new(Array)
	result.items = t.items[from:to]

	return result
}

func (t *Array) IndexOf(item core.Value) Int {
	res := Int(-1)

	for idx, el := range t.items {
		if el.Compare(item) == 0 {
			res = Int(idx)
			break
		}
	}

	return res
}

func (t *Array) Insert(idx Int, value core.Value) {
	t.items = append(t.items[:idx], append([]core.Value{value}, t.items[idx:]...)...)
}

func (t *Array) RemoveAt(idx Int) {
	i := int(idx)
	max := len(t.items) - 1

	if i > max {
		return
	}

	t.items = append(t.items[:i], t.items[i+1:]...)
}

func (t *Array) Clone() core.Cloneable {
	cloned := NewArray(0)

	var value core.Value
	for idx := NewInt(0); idx < t.Length(); idx++ {
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
		return first.Compare(second) == -1
	})
}

func (t *Array) SortWith(sorter ArraySorter) *Array {
	c := make([]core.Value, len(t.items))
	copy(c, t.items)

	sort.SliceStable(c, func(i, j int) bool {
		return sorter(c[i], c[j])
	})

	res := new(Array)
	res.items = c

	return res
}

func (t *Array) GetIn(ctx context.Context, path []core.Value) (core.Value, core.PathError) {
	if len(path) == 0 {
		return None, nil
	}

	segmentIdx := 0

	if typ := path[segmentIdx].Type(); typ != types.Int {
		return None, core.NewPathError(core.TypeError(typ, types.Int), segmentIdx)
	}

	first := t.Get(path[segmentIdx].(Int))

	if len(path) == 1 {
		return first, nil
	}

	segmentIdx++

	if first == None || first == nil {
		return None, core.NewPathError(core.ErrInvalidPath, segmentIdx)
	}

	getter, ok := first.(core.Getter)

	if !ok {
		return GetIn(ctx, first, path[segmentIdx:])
	}

	return getter.GetIn(ctx, path[segmentIdx:])
}
