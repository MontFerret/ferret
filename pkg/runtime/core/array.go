package core

import (
	"context"
	"encoding/binary"
	"hash/fnv"
	"sort"

	"github.com/wI2L/jettison"
)

type arrayList struct {
	data []Value
}

func EmptyArray() List {
	return &arrayList{data: make([]Value, 0)}
}

func NewArray(cap int) List {
	return &arrayList{data: make([]Value, 0, cap)}
}

func NewSizedArray(size int) List {
	return &arrayList{data: make([]Value, size)}
}

func NewArrayWith(values ...Value) List {
	return &arrayList{data: values}
}

func (t *arrayList) MarshalJSON() ([]byte, error) {
	return jettison.MarshalOpts(t.data, jettison.NoHTMLEscaping())
}

func (t *arrayList) String() string {
	marshaled, err := t.MarshalJSON()

	if err != nil {
		return "[]"
	}

	return string(marshaled)
}

func (t *arrayList) Compare(ctx context.Context, other Value) (int64, error) {
	otherArr, ok := other.(*arrayList)

	if !ok {
		return CompareTypes(t, other), nil
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

	for i := 0; i < size; i++ {
		thisVal := t.data[i]
		otherVal := otherArr.data[i]

		comp, err := CompareValues(ctx, thisVal, otherVal)

		if err != nil {
			return 0, err
		}

		if comp != 0 {
			return comp, nil
		}

		res = comp
	}

	return res, nil
}

func (t *arrayList) Unwrap() interface{} {
	arr := make([]interface{}, len(t.data))

	for idx, val := range t.data {
		arr[idx] = val.Unwrap()
	}

	return arr
}

func (t *arrayList) Hash() uint64 {
	h := fnv.New64a()

	h.Write([]byte("array:"))
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

func (t *arrayList) Copy() Value {
	ctx := context.Background()
	c := NewArray(len(t.data))

	for _, el := range t.data {
		_ = c.Add(ctx, el)
	}

	return c
}

func (t *arrayList) Clone(ctx context.Context) (Cloneable, error) {
	size := len(t.data)
	res := &arrayList{data: make([]Value, size)}

	var value Value

	for idx := 0; idx < len(t.data); idx++ {
		value = t.data[idx]
		cloned, err := CloneOrCopy(ctx, value)

		if err != nil {
			return nil, err
		}

		res.data[idx] = cloned
	}

	return res, nil
}

func (t *arrayList) Iterate(_ context.Context) (Iterator, error) {
	return NewArrayIterator(t), nil
}

func (t *arrayList) Length(_ context.Context) (Int, error) {
	return Int(len(t.data)), nil
}

func (t *arrayList) IsEmpty(_ context.Context) (Boolean, error) {
	return len(t.data) == 0, nil
}

func (t *arrayList) Contains(ctx context.Context, value Value) (Boolean, error) {
	idx, err := t.IndexOf(ctx, value)

	if err != nil {
		return false, err
	}

	return idx >= 0, nil
}

func (t *arrayList) IndexOf(ctx context.Context, item Value) (Int, error) {
	for idx, el := range t.data {
		comp, err := CompareValues(ctx, item, el)

		if err != nil {
			return -1, err
		}

		if comp == 0 {
			return Int(idx), nil
		}
	}

	return -1, nil
}

func (t *arrayList) Get(_ context.Context, idx Int) (Value, error) {
	l := Int(len(t.data) - 1)

	if l < 0 {
		return None, nil
	}

	if idx > l {
		return None, nil
	}

	return t.data[idx], nil
}

func (t *arrayList) First(_ context.Context) (Value, error) {
	if len(t.data) > 0 {
		return t.data[0], nil
	}

	return None, nil
}

func (t *arrayList) Last(_ context.Context) (Value, error) {
	size := len(t.data)

	if size > 1 {
		return t.data[size-1], nil
	} else if size == 1 {
		return t.data[0], nil
	}

	return None, nil
}

func (t *arrayList) Find(ctx context.Context, predicate IndexedPredicate) (List, error) {
	result := NewArray(len(t.data))
	size := Int(len(t.data))

	for idx := Int(0); idx < size; idx++ {
		val := t.data[idx]
		res, err := predicate(ctx, val, idx)

		if err != nil {
			return nil, err
		}

		if res {
			_ = result.Add(ctx, val)
		}
	}

	return result, nil
}

func (t *arrayList) FindOne(ctx context.Context, predicate IndexedPredicate) (Value, Boolean, error) {
	size := Int(len(t.data))

	for idx := Int(0); idx < size; idx++ {
		val := t.data[idx]
		res, err := predicate(ctx, val, idx)

		if err != nil {
			return None, false, err
		}

		if res {
			return val, true, nil
		}
	}

	return None, false, nil
}

func (t *arrayList) Slice(_ context.Context, start, end Int) (List, error) {
	length := Int(len(t.data))

	if start >= length {
		return NewArray(0), nil
	}

	if end > length {
		end = length
	}

	result := new(arrayList)
	result.data = t.data[start:end]

	return result, nil
}

func (t *arrayList) Sort(ctx context.Context, ascending Boolean) (List, error) {
	var pivot int64 = -1

	if ascending {
		pivot = 1
	}

	return t.SortWith(ctx, func(c context.Context, first, second Value) (int64, error) {
		comp, err := CompareValues(c, first, second)

		if err != nil {
			return 0, err
		}

		return pivot * comp, nil
	})
}

func (t *arrayList) SortWith(ctx context.Context, comparator Comparator) (List, error) {
	c := make([]Value, len(t.data))
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

	res := new(arrayList)
	res.data = c

	return res, nil
}

func (t *arrayList) ForEach(ctx context.Context, predicate IndexedPredicate) error {
	size := Int(len(t.data))

	for idx := Int(0); idx < size; idx++ {
		val := t.data[idx]
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

func (t *arrayList) Add(_ context.Context, value Value) error {
	t.data = append(t.data, value)

	return nil
}

func (t *arrayList) Set(_ context.Context, idx Int, value Value) error {
	last := Int(len(t.data) - 1)

	if last >= idx {
		t.data[idx] = value

		return nil
	}

	return Error(ErrInvalidOperation, "out of bounds")
}

func (t *arrayList) Insert(_ context.Context, idx Int, value Value) error {
	t.data = append(t.data[:idx], append([]Value{value}, t.data[idx:]...)...)

	return nil
}

func (t *arrayList) Clear(_ context.Context) error {
	t.data = make([]Value, 0)

	return nil
}

func (t *arrayList) Remove(ctx context.Context, value Value) error {
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

func (t *arrayList) RemoveAt(_ context.Context, idx Int) (Value, error) {
	edge := Int(len(t.data) - 1)

	if idx > edge {
		return None, nil
	}

	item := t.data[idx]

	t.data = append(t.data[:idx], t.data[idx+1:]...)

	return item, nil
}

func (t *arrayList) Swap(_ context.Context, i, j Int) error {
	t.data[i], t.data[j] = t.data[j], t.data[i]

	return nil
}
