package internal

import (
	"context"
	"encoding/binary"
	"go/types"
	"hash/fnv"
	"sort"

	"github.com/wI2L/jettison"

	"github.com/MontFerret/ferret/pkg/runtime/core"
)

type Array struct {
	data []core.Value
}

func EmptyArray() core.List {
	return &Array{data: make([]core.Value, 0)}
}

func NewArray(cap int) core.List {
	return &Array{data: make([]core.Value, 0, cap)}
}

func NewSizedArray(size int) core.List {
	return &Array{data: make([]core.Value, size)}
}

func NewArrayWith(values ...core.Value) core.List {
	return &Array{data: values}
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

func (t *Array) Compare(ctx context.Context, other core.Value) (int64, error) {
	otherArr, ok := other.(*Array)

	if !ok {
		return types.Compare(types.Array, core.Reflect(other)), nil
	}

	size := len(t.data)
	otherArrSize := len(otherArr.data)

	if size == 0 && otherArrSize == 0 {
		return 0, nil
	}

	if size < otherArrSize {
		return -1, nil
	}

	if size > otherArrSize {
		return 1, nil
	}

	var res int64

	otherArr.ForEach(ctx, func(c context.Context, otherVal core.Value, idx int) (bool, error) {
		val, _ := t.Get(ctx, idx)
		comp, err := core.CompareValues(ctx, val, otherVal)

		if err != nil {
			return false, err
		}

		res = comp

		return res == 0, nil
	})

	return res, nil
}

func (t *Array) Unwrap() interface{} {
	arr := make([]interface{}, len(t.data))

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
	ctx := context.Background()
	c := NewArray(len(t.data))

	for _, el := range t.data {
		c.Add(ctx, el)
	}

	return c
}

func (t *Array) Clone(ctx context.Context) (core.Cloneable, error) {
	cloned := NewArray(0)

	var value core.Value
	for idx := 0; idx < len(t.data); idx++ {
		value, _ = t.Get(ctx, idx)

		cloneable, ok := value.(core.Cloneable)

		if ok {
			clone, err := cloneable.Clone(ctx)

			if err != nil {
				return nil, err
			}

			value = clone
		} else {
			value = value.Copy()
		}

		cloned.Add(ctx, value)
	}

	return cloned, nil
}

func (t *Array) Iterate(_ context.Context) (core.Iterator, error) {
	return NewArrayIterator(t), nil
}

func (t *Array) Length(_ context.Context) (int, error) {
	return len(t.data), nil
}

func (t *Array) IsEmpty(_ context.Context) (bool, error) {
	return len(t.data) == 0, nil
}

func (t *Array) Contains(ctx context.Context, value core.Value) (bool, error) {
	idx, err := t.IndexOf(ctx, value)

	if err != nil {
		return false, err
	}

	return idx >= 0, nil
}

func (t *Array) IndexOf(ctx context.Context, item core.Value) (int, error) {
	for idx, el := range t.data {
		comp, err := core.CompareValues(ctx, item, el)

		if err != nil {
			return -1, err
		}

		if comp == 0 {
			return idx, nil
		}
	}

	return -1, nil
}

func (t *Array) Get(_ context.Context, idx int) (core.Value, error) {
	l := len(t.data) - 1

	if l < 0 {
		return core.None, nil
	}

	if idx > l {
		return core.None, nil
	}

	return t.data[idx], nil
}

func (t *Array) First(_ context.Context) (core.Value, error) {
	if len(t.data) > 0 {
		return t.data[0], nil
	}

	return core.None, nil
}

func (t *Array) Last(_ context.Context) (core.Value, error) {
	size := len(t.data)

	if size > 1 {
		return t.data[size-1], nil
	} else if size == 1 {
		return t.data[0], nil
	}

	return core.None, nil
}

func (t *Array) Find(ctx context.Context, predicate core.IndexedPredicate) (core.List, error) {
	result := NewArray(len(t.data))

	for idx, val := range t.data {
		res, err := predicate(ctx, val, idx)

		if err != nil {
			return nil, err
		}

		if res {
			result.Add(ctx, val)
		}
	}

	return result, nil
}

func (t *Array) FindOne(ctx context.Context, predicate core.IndexedPredicate) (core.Value, bool, error) {
	for idx, val := range t.data {
		res, err := predicate(ctx, val, idx)

		if err != nil {
			return core.None, false, err
		}

		if res {
			return val, true, nil
		}
	}

	return core.None, false, nil
}

func (t *Array) Slice(_ context.Context, start, end int) (core.List, error) {
	length := len(t.data)

	if start >= length {
		return NewArray(0), nil
	}

	if end > length {
		end = length
	}

	result := new(Array)
	result.data = t.data[start:end]

	return result, nil
}

func (t *Array) Sort(ctx context.Context, ascending bool) (core.List, error) {
	var pivot int64 = -1

	if ascending {
		pivot = 1
	}

	return t.SortWith(ctx, func(c context.Context, first, second core.Value) (int64, error) {
		comp, err := core.CompareValues(c, first, second)

		if err != nil {
			return 0, err
		}

		return pivot * comp, nil
	})
}

func (t *Array) SortWith(ctx context.Context, comparator core.Comparator) (core.List, error) {
	c := make([]core.Value, len(t.data))
	copy(c, t.data)

	var err error

	sort.SliceStable(c, func(i, j int) bool {
		comp, e := comparator(ctx, c[i], c[j])

		if e != nil {
			err = e

			return true
		}

		return comp == 0
	})

	if err != nil {
		return nil, err
	}

	res := new(Array)
	res.data = c

	return res, nil
}

func (t *Array) ForEach(ctx context.Context, predicate core.IndexedPredicate) error {
	for idx, val := range t.data {
		res, err := predicate(ctx, val, idx)

		if err != nil {
			return err
		}

		if !res {
			break
		}
	}

	return nil
}

func (t *Array) Add(_ context.Context, value core.Value) error {
	t.data = append(t.data, value)

	return nil
}

func (t *Array) Set(_ context.Context, idx int, value core.Value) error {
	last := len(t.data) - 1

	if last >= idx {
		t.data[idx] = value

		return nil
	}

	return core.Error(core.ErrInvalidOperation, "out of bounds")
}

func (t *Array) Insert(_ context.Context, idx int, value core.Value) error {
	t.data = append(t.data[:idx], append([]core.Value{value}, t.data[idx:]...)...)

	return nil
}

func (t *Array) Clear(_ context.Context) error {
	t.data = make([]core.Value, 0)

	return nil
}

func (t *Array) Remove(ctx context.Context, value core.Value) error {
	idx, err := t.IndexOf(ctx, value)

	if err != nil {
		return err
	}

	if idx < 0 {
		return nil
	}

	_, err = t.RemoveAt(ctx, idx)

	return err
}

func (t *Array) RemoveAt(_ context.Context, idx int) (core.Value, error) {
	edge := len(t.data) - 1

	if idx > edge {
		return core.None, nil
	}

	item := t.data[idx]

	t.data = append(t.data[:idx], t.data[idx+1:]...)

	return item, nil
}

func (t *Array) Swap(_ context.Context, i, j int) error {
	t.data[i], t.data[j] = t.data[j], t.data[i]

	return nil
}
