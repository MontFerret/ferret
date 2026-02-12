package runtime

import (
	"context"
	"encoding/binary"
	"hash/fnv"
	"sort"

	"github.com/wI2L/jettison"
)

type (
	ObjectProperty struct {
		key   string
		value Value
	}

	// Object represents a JSON object.
	// It is a collection of key-value pairs where keys are strings and values are of any type.
	Object struct {
		data map[string]Value
	}
)

func NewObjectProperty(name string, value Value) *ObjectProperty {
	return &ObjectProperty{name, value}
}

func newObjectOf(size int) *Object {
	return &Object{make(map[string]Value, size)}
}

func newObjectWith(props ...*ObjectProperty) *Object {
	obj := &Object{make(map[string]Value)}

	for _, prop := range props {
		obj.data[prop.key] = prop.value
	}

	return obj
}

func (t *Object) MarshalJSON() ([]byte, error) {
	return jettison.MarshalOpts(t.data, jettison.NoHTMLEscaping())
}

func (t *Object) String() string {
	marshaled, err := t.MarshalJSON()

	if err != nil {
		return "{}"
	}

	return string(marshaled)
}

// Compare compares the source object with other core.Value
// The behavior of the Compare is similar
// to the comparison of objects in ArangoDB
func (t *Object) Compare(ctx Context, other Value) int64 {
	otherObject, ok := other.(*Object)

	if !ok {
		return CompareTypes(t, other)
	}

	size := len(t.data)
	otherSize := len(otherObject.data)

	if size == 0 && otherSize == 0 {
		return 0
	}

	if size < otherSize {
		return -1
	}

	if size > otherSize {
		return 1
	}

	var res int64

	tKeys := make([]string, 0, size)

	for k := range t.data {
		tKeys = append(tKeys, k)
	}

	sortedT := sort.StringSlice(tKeys)
	sortedT.Sort()

	otherKeys := make([]string, 0, otherSize)

	for k := range otherObject.data {
		otherKeys = append(otherKeys, k)
	}

	sortedOther := sort.StringSlice(otherKeys)
	sortedOther.Sort()

	var tVal, otherVal Value
	var tKey, otherKey string

	for i := 0; i < len(t.data) && res == 0; i++ {
		tKey, otherKey = sortedT[i], sortedOther[i]

		if tKey == otherKey {
			tVal = t.data[tKey]
			otherVal = otherObject.data[tKey]
			res = CompareValues(ctx, tVal, otherVal)

			continue
		}

		if tKey < otherKey {
			res = 1
		} else {
			res = -1
		}

		break
	}

	return res
}

func (t *Object) Unwrap() interface{} {
	obj := make(map[string]interface{})

	for key, val := range t.data {
		obj[key] = val.Unwrap()
	}

	return obj
}

func (t *Object) Hash() uint64 {
	h := fnv.New64a()

	h.Write([]byte("object:"))
	h.Write([]byte("{"))

	keys := make([]string, 0, len(t.data))

	for key := range t.data {
		keys = append(keys, key)
	}

	// order does not really matter
	// but it will give us a consistent hash sum
	sort.Strings(keys)
	endIndex := len(keys) - 1

	for idx, key := range keys {
		h.Write([]byte(key))
		h.Write([]byte(":"))

		el := t.data[key]

		bytes := make([]byte, 8)
		binary.LittleEndian.PutUint64(bytes, el.Hash())

		h.Write(bytes)

		if idx != endIndex {
			h.Write([]byte(","))
		}
	}

	h.Write([]byte("}"))

	return h.Sum64()
}

func (t *Object) Copy(ctx Context) (Value, error) {
	c := ctx.Alloc().Object(len(t.data))

	for k, v := range t.data {
		c.data[k] = v
	}

	return c, nil
}

func (t *Object) Clone(ctx Context) (Cloneable, error) {
	cloned := ctx.Alloc().Object(len(t.data))

	var value Value

	for key := range t.data {
		value = t.data[key]
		cloneable, ok := value.(Cloneable)

		if ok {
			clone, err := cloneable.Clone(ctx)

			if err != nil {
				return nil, err
			}

			value = clone
		} else {
			value, _ = value.Copy(nil)
		}

		cloned.data[key] = value
	}

	return cloned, nil
}

func (t *Object) Length(_ Context) (Int, error) {
	return Int(len(t.data)), nil
}

func (t *Object) IsEmpty(_ context.Context) (Boolean, error) {
	return len(t.data) == 0, nil
}

func (t *Object) Keys(ctx Context) (List, error) {
	keys := ctx.Alloc().Array(len(t.data))

	for k := range t.data {
		_ = keys.Append(ctx, NewString(k))
	}

	return keys, nil
}

func (t *Object) Values(ctx Context) (List, error) {
	keys := ctx.Alloc().Array(len(t.data))

	for _, v := range t.data {
		_ = keys.Append(ctx, v)
	}

	return keys, nil
}

func (t *Object) ForEach(ctx Context, predicate KeyReadablePredicate) error {
	for key, val := range t.data {
		doContinue, err := predicate(ctx, val, String(key))

		if err != nil {
			return err
		}

		if !doContinue {
			break
		}
	}

	return nil
}

func (t *Object) Find(ctx Context, predicate KeyReadablePredicate) (List, error) {
	res := ctx.Alloc().Array(len(t.data))

	for key, val := range t.data {
		match, err := predicate(ctx, val, String(key))

		if err != nil {
			return nil, err
		}

		if match {
			_ = res.Append(ctx, val)
		}
	}

	return res, nil
}

func (t *Object) FindOne(ctx Context, predicate KeyReadablePredicate) (Value, Boolean, error) {
	for key, val := range t.data {
		res, err := predicate(ctx, val, String(key))

		if err != nil {
			return None, false, err
		}

		if res {
			return val, true, nil
		}
	}

	return None, false, nil
}

func (t *Object) ContainsKey(_ Context, key Value) (Boolean, error) {
	_, exists := t.data[key.String()]

	return Boolean(exists), nil
}

func (t *Object) ContainsValue(_ Context, value Value) (Boolean, error) {
	for _, val := range t.data {
		res := CompareValues(nil, value, val)

		if res == 0 {
			return true, nil
		}
	}

	return false, nil
}

func (t *Object) Get(_ Context, key Value) (Value, error) {
	val, found := t.data[key.String()]

	if found {
		return val, nil
	}

	return None, ErrNotFound
}

func (t *Object) Set(_ Context, key, value Value) error {
	if value == nil {
		value = None
	}

	t.data[key.String()] = value

	return nil
}

func (t *Object) RemoveValue(ctx Context, value Value) error {
	for key, val := range t.data {
		if CompareValues(ctx, value, val) == 0 {
			delete(t.data, key)

			break
		}
	}

	return nil
}

func (t *Object) RemoveKey(_ Context, key Value) error {
	delete(t.data, key.String())

	return nil
}

func (t *Object) Clear(_ Context) error {
	t.data = make(map[string]Value)

	return nil
}

func (t *Object) Iterate(ctx Context) (Iterator, error) {
	// TODO: implement channel based iterator
	return NewObjectIterator(t), nil
}
