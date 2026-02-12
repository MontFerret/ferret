package runtime

import (
	"encoding/binary"
	"hash/fnv"

	"github.com/wI2L/jettison"
)

type Array struct {
	data []Value
}

func newArray(cap int) *Array {
	return &Array{data: make([]Value, 0, cap)}
}

func (t *Array) MarshalJSON() ([]byte, error) {
	return jettison.MarshalOpts(t.data, jettison.NoHTMLEscaping())
}

func (t *Array) String() string {
	marshaled, err := t.MarshalJSON()

	if err != nil {
		return "[]"
	}

	return string(marshaled)
}

func (t *Array) Compare(_ Context, other Value) int64 {
	otherArr, ok := other.(*Array)

	if !ok {
		return CompareTypes(t, other)
	}

	size := len(t.data)
	otherArrSize := len(otherArr.data)

	if size == 0 && otherArrSize == 0 {
		return 0
	}

	if size < otherArrSize {
		return -1
	}

	if size > otherArrSize {
		return 1
	}

	var res int64

	for i := 0; i < size; i++ {
		thisVal := t.data[i]
		otherVal := otherArr.data[i]

		comp := CompareValues(nil, thisVal, otherVal)

		if comp != 0 {
			return comp
		}

		res = comp
	}

	return res
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

func (t *Array) Copy(Context) (Value, error) {
	return &Array{data: t.copyInternal(0)}, nil
}

func (t *Array) CopyWithGrowth(cap Int) *Array {
	return &Array{data: t.copyInternal(cap)}
}

func (t *Array) copyInternal(cap Int) []Value {
	c := make([]Value, 0, len(t.data)+int(cap))
	c = append(c, t.data...)

	return c
}

func (t *Array) Clone(ctx Context) (Cloneable, error) {
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

func (t *Array) Iterate(_ Context) (Iterator, error) {
	return NewArrayIterator(t), nil
}

func (t *Array) Length(_ Context) (Int, error) {
	return Int(len(t.data)), nil
}

func (t *Array) IsEmpty(_ Context) (Boolean, error) {
	return len(t.data) == 0, nil
}

func (t *Array) Contains(ctx Context, value Value) (Boolean, error) {
	idx, err := t.IndexOf(ctx, value)

	if err != nil {
		return false, err
	}

	return idx >= 0, nil
}

func (t *Array) IndexOf(ctx Context, value Value) (Int, error) {
	for idx, el := range t.data {
		comp := CompareValues(ctx, value, el)

		if comp == 0 {
			return Int(idx), nil
		}
	}

	return -1, nil
}

func (t *Array) Get(_ Context, idx Int) (Value, error) {
	l := Int(len(t.data) - 1)

	if l < 0 {
		return None, nil
	}

	if idx > l {
		return None, nil
	}

	return t.data[idx], nil
}

func (t *Array) First(_ Context) (Value, error) {
	if len(t.data) > 0 {
		return t.data[0], nil
	}

	return None, nil
}

func (t *Array) Last(_ Context) (Value, error) {
	size := len(t.data)

	if size > 1 {
		return t.data[size-1], nil
	} else if size == 1 {
		return t.data[0], nil
	}

	return None, nil
}

func (t *Array) Find(ctx Context, predicate IndexReadablePredicate) (List, error) {
	result := ctx.Alloc().Array(len(t.data))
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

func (t *Array) FindOne(ctx Context, predicate IndexReadablePredicate) (Value, Boolean, error) {
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

func (t *Array) Slice(ctx Context, start, end Int) (List, error) {
	length := Int(len(t.data))

	if start >= length {
		return ctx.Alloc().Array(0), nil
	}

	if end > length {
		end = length
	}

	result := new(Array)
	result.data = t.data[start:end]

	return result, nil
}

func (t *Array) SortAsc(ctx Context) error {
	return t.sort(ctx, true)
}

func (t *Array) SortDesc(ctx Context) error {
	return t.sort(ctx, false)
}

func (t *Array) sort(ctx Context, ascending Boolean) error {
	SortSlice(ctx, t.data, ascending)

	return nil
}

func (t *Array) SortWith(ctx Context, comparator Comparator) error {
	c := make([]Value, len(t.data))
	copy(c, t.data)

	SortSliceWith(ctx, t.data, comparator)

	res := new(Array)
	res.data = c

	return nil
}

func (t *Array) ForEach(ctx Context, predicate IndexReadablePredicate) error {
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

func (t *Array) Append(_ Context, val Value) error {
	t.data = append(t.data, val)

	return nil
}

func (t *Array) Set(_ Context, idx Int, value Value) error {
	last := Int(len(t.data) - 1)

	if last >= idx {
		t.data[idx] = value

		return nil
	}

	return Error(ErrInvalidOperation, "out of bounds")
}

func (t *Array) Insert(_ Context, idx Int, value Value) error {
	t.data = append(t.data[:idx], append([]Value{value}, t.data[idx:]...)...)

	return nil
}

func (t *Array) Clear(_ Context) error {
	t.data = make([]Value, 0, cap(t.data))

	return nil
}

func (t *Array) RemoveValue(ctx Context, value Value) error {
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

func (t *Array) RemoveAt(_ Context, idx Int) (Value, error) {
	edge := Int(len(t.data) - 1)

	if idx > edge {
		return None, nil
	}

	item := t.data[idx]

	t.data = append(t.data[:idx], t.data[idx+1:]...)

	return item, nil
}

func (t *Array) Swap(_ Context, i, j Int) error {
	t.data[i], t.data[j] = t.data[j], t.data[i]

	return nil
}
