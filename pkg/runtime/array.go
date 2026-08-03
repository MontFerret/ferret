package runtime

import (
	"context"
	"encoding/binary"
	"hash/fnv"
	"strings"
)

type Array struct {
	data []Value
}

func EmptyArray() *Array {
	return &Array{data: make([]Value, 0)}
}

func NewArray(cap int) *Array {
	return &Array{data: make([]Value, 0, cap)}
}

func NewArray64(cap Int) *Array {
	return &Array{data: make([]Value, 0, cap)}
}

func NewSizedArray(size int) *Array {
	return &Array{data: make([]Value, size)}
}

func NewArrayWith(values ...Value) *Array {
	return &Array{data: values}
}

func NewArrayOf(values []Value) *Array {
	return &Array{data: values}
}

func (t *Array) Type() Type {
	return TypeArray
}

func (t *Array) String() string {
	var b strings.Builder

	writeArray(&b, t)

	return b.String()
}

func (t *Array) Equal(ctx context.Context, other Value) (bool, error) {
	if otherArray, ok := other.(*Array); ok {
		return t.equalArray(ctx, otherArray)
	}

	otherList, ok := other.(List)
	if !ok {
		return false, nil
	}

	otherSize, err := otherList.Length(ctx)
	if err != nil {
		return false, err
	}

	if Int(len(t.data)) != otherSize {
		return false, nil
	}

	for idx, value := range t.data {
		otherValue, err := otherList.At(ctx, Int(idx))
		if err != nil {
			return false, err
		}

		equal, err := EqualValues(ctx, value, otherValue)
		if err != nil {
			return false, err
		}

		if !equal {
			return false, nil
		}
	}

	return true, nil
}

func (t *Array) equalArray(ctx context.Context, other *Array) (bool, error) {
	if len(t.data) != len(other.data) {
		return false, nil
	}

	for idx, value := range t.data {
		equal, err := EqualValues(ctx, value, other.data[idx])
		if err != nil {
			return false, err
		}
		if !equal {
			return false, nil
		}
	}

	return true, nil
}

func (t *Array) Compare(ctx context.Context, other Value) (Ordering, error) {
	if otherArray, ok := other.(*Array); ok {
		return t.compareArray(ctx, otherArray)
	}

	otherList, ok := other.(List)
	if !ok {
		return Equal, incompatibleComparisonError(t, other)
	}

	otherSize, err := otherList.Length(ctx)
	if err != nil {
		return Equal, err
	}

	size := Int(len(t.data))

	if size < otherSize {
		return Less, nil
	}

	if size > otherSize {
		return Greater, nil
	}

	for idx, value := range t.data {
		otherValue, err := otherList.At(ctx, Int(idx))
		if err != nil {
			return Equal, err
		}

		comparison, err := CompareValues(ctx, value, otherValue)
		if err != nil {
			return Equal, err
		}

		if comparison != Equal {
			return comparison, nil
		}
	}

	return Equal, nil
}

func (t *Array) compareArray(ctx context.Context, other *Array) (Ordering, error) {
	if len(t.data) < len(other.data) {
		return Less, nil
	}

	if len(t.data) > len(other.data) {
		return Greater, nil
	}

	for idx, value := range t.data {
		comparison, err := CompareValues(ctx, value, other.data[idx])
		if err != nil {
			return Equal, err
		}
		if comparison != Equal {
			return comparison, nil
		}
	}

	return Equal, nil
}

func (t *Array) Hash() uint64 {
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

func (t *Array) Copy() Value {
	return &Array{data: t.copyInternal(0)}
}

func (t *Array) CopyWithGrowth(cap Int) *Array {
	return &Array{data: t.copyInternal(cap)}
}

func (t *Array) copyInternal(cap Int) []Value {
	c := make([]Value, 0, len(t.data)+int(cap))
	c = append(c, t.data...)

	return c
}

func (t *Array) Clone(ctx context.Context) (Cloneable, error) {
	size := len(t.data)
	res := &Array{data: make([]Value, size)}

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

func (t *Array) Iterate(_ context.Context) (Iterator, error) {
	return NewArrayIterator(t), nil
}

func (t *Array) Length(_ context.Context) (Int, error) {
	return Int(len(t.data)), nil
}

func (t *Array) Contains(ctx context.Context, value Value) (Boolean, error) {
	idx, err := t.IndexOf(ctx, value)

	if err != nil {
		return false, err
	}

	return idx >= 0, nil
}

func (t *Array) IndexOf(ctx context.Context, item Value) (Int, error) {
	for idx, el := range t.data {
		equal, err := EqualValues(ctx, item, el)
		if err != nil {
			return -1, err
		}

		if equal {
			return Int(idx), nil
		}
	}

	return -1, nil
}

func (t *Array) At(_ context.Context, idx Int) (Value, error) {
	l := Int(len(t.data) - 1)

	if l < 0 {
		return None, nil
	}

	if idx > l {
		return None, nil
	}

	return t.data[idx], nil
}

func (t *Array) LookupAt(_ context.Context, idx Int) (Value, bool, error) {
	l := Int(len(t.data) - 1)

	if l < 0 {
		return None, false, nil
	}

	if idx > l {
		return None, false, nil
	}

	return t.data[idx], true, nil
}

func (t *Array) First(_ context.Context) (Value, error) {
	if len(t.data) > 0 {
		return t.data[0], nil
	}

	return None, nil
}

func (t *Array) Last(_ context.Context) (Value, error) {
	size := len(t.data)

	if size > 1 {
		return t.data[size-1], nil
	} else if size == 1 {
		return t.data[0], nil
	}

	return None, nil
}

func (t *Array) Filter(ctx context.Context, predicate IndexReadablePredicate) (List, error) {
	result := NewArray(len(t.data))
	size := Int(len(t.data))

	for idx := Int(0); idx < size; idx++ {
		val := t.data[idx]
		res, err := predicate(ctx, val, idx)

		if err != nil {
			return nil, err
		}

		if res {
			_ = result.Append(ctx, val)
		}
	}

	return result, nil
}

func (t *Array) Find(ctx context.Context, predicate IndexReadablePredicate) (Value, Boolean, error) {
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

func (t *Array) Slice(_ context.Context, start, end Int) (List, error) {
	length := Int(len(t.data))

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

func (t *Array) SortAsc(ctx context.Context) error {
	return t.sort(ctx, true)
}

func (t *Array) SortDesc(ctx context.Context) error {
	return t.sort(ctx, false)
}

func (t *Array) sort(ctx context.Context, ascending Boolean) error {
	return SortSlice(ctx, t.data, ascending)
}

func (t *Array) SortWith(ctx context.Context, comparator Comparator) error {
	return SortSliceWith(ctx, t.data, comparator)
}

func (t *Array) ForEach(ctx context.Context, predicate IndexReadablePredicate) error {
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

func (t *Array) Append(_ context.Context, value Value) error {
	t.data = append(t.data, value)

	return nil
}

func (t *Array) SetAt(_ context.Context, idx Int, value Value) error {
	last := Int(len(t.data) - 1)

	if last >= idx {
		t.data[idx] = value

		return nil
	}

	return Error(ErrInvalidOperation, "out of bounds")
}

func (t *Array) Insert(_ context.Context, idx Int, value Value) error {
	t.data = append(t.data[:idx], append([]Value{value}, t.data[idx:]...)...)

	return nil
}

func (t *Array) Clear(_ context.Context) error {
	t.data = make([]Value, 0)

	return nil
}

func (t *Array) Empty(_ context.Context) (List, error) {
	return NewArray(0), nil
}

func (t *Array) Remove(ctx context.Context, value Value) error {
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

func (t *Array) RemoveAt(_ context.Context, idx Int) (Value, error) {
	edge := Int(len(t.data) - 1)

	if idx > edge {
		return None, nil
	}

	item := t.data[idx]

	t.data = append(t.data[:idx], t.data[idx+1:]...)

	return item, nil
}

func (t *Array) Swap(_ context.Context, i, j Int) error {
	t.data[i], t.data[j] = t.data[j], t.data[i]

	return nil
}

func (t *Array) Concat(ctx context.Context, other List) error {
	switch list := other.(type) {
	case *Array:
		t.data = append(t.data, list.data...)

		return nil
	default:
		return ForEach(ctx, other, func(ctx context.Context, value, _ Value) (Boolean, error) {
			t.data = append(t.data, value)

			return true, nil
		})
	}
}
