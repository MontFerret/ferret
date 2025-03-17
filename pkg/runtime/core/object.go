package core

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

	hashMap struct {
		data map[string]Value
	}
)

func NewObjectProperty(name string, value Value) *ObjectProperty {
	return &ObjectProperty{name, value}
}

func NewObject() Map {
	return &hashMap{make(map[string]Value)}
}

func NewObjectWith(props ...*ObjectProperty) Map {
	obj := &hashMap{make(map[string]Value)}

	for _, prop := range props {
		obj.data[prop.key] = prop.value
	}

	return obj
}

func (t *hashMap) Type() string {
	return "object"
}

func (t *hashMap) MarshalJSON() ([]byte, error) {
	return jettison.MarshalOpts(t.data, jettison.NoHTMLEscaping())
}

func (t *hashMap) String() string {
	marshaled, err := t.MarshalJSON()

	if err != nil {
		return "{}"
	}

	return string(marshaled)
}

// Compare compares the source object with other core.Value
// The behavior of the Compare is similar
// to the comparison of objects in ArangoDB
func (t *hashMap) Compare(ctx context.Context, other Value) (int64, error) {
	otherObject, ok := other.(*hashMap)

	if !ok {
		// TODO: Maybe we should return an error here
		return CompareTypes(t, other), nil
	}

	size := len(t.data)
	otherSize := len(otherObject.data)

	if size == 0 && otherSize == 0 {
		return 0, nil
	}

	if size < otherSize {
		return -1, nil
	}

	if size > otherSize {
		return 1, nil
	}

	var res int64

	tKeys := make([]string, 0, size)

	for k := range t.data {
		tKeys = append(tKeys, k)
	}

	sortedT := sort.StringSlice(tKeys)
	sortedT.Sort()

	otherKeys := make([]string, 0, otherSize)

	_ = otherObject.ForEach(ctx, func(_ context.Context, value, k Value) (Boolean, error) {
		otherKeys = append(otherKeys, k.String())
		return true, nil
	})

	sortedOther := sort.StringSlice(otherKeys)
	sortedOther.Sort()

	var tVal, otherVal Value
	var tKey, otherKey string

	for i := 0; i < len(t.data) && res == 0; i++ {
		tKey, otherKey = sortedT[i], sortedOther[i]

		if tKey == otherKey {
			tVal, _ = t.Get(ctx, NewString(tKey))
			otherVal, _ = otherObject.Get(ctx, NewString(tKey))
			comp, err := CompareValues(ctx, tVal, otherVal)

			if err != nil {
				return 0, err
			}

			res = comp

			continue
		}

		if tKey < otherKey {
			res = 1
		} else {
			res = -1
		}

		break
	}

	return res, nil
}

func (t *hashMap) Unwrap() interface{} {
	obj := make(map[string]interface{})

	for key, val := range t.data {
		obj[key] = val.Unwrap()
	}

	return obj
}

func (t *hashMap) Hash() uint64 {
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

func (t *hashMap) Copy() Value {
	c := &hashMap{make(map[string]Value)}

	for k, v := range t.data {
		c.data[k] = v
	}

	return c
}

func (t *hashMap) Clone(ctx context.Context) (Cloneable, error) {
	cloned := &hashMap{make(map[string]Value)}

	var value Value

	for key := range t.data {
		value, _ = t.data[key]

		cloneable, ok := value.(Cloneable)

		if ok {
			clone, err := cloneable.Clone(ctx)

			if err != nil {
				return nil, err
			}

			value = clone
		} else {
			value = value.Copy()
		}

		cloned.data[key] = value
	}

	return cloned, nil
}

func (t *hashMap) Length(_ context.Context) (Int, error) {
	return Int(len(t.data)), nil
}

func (t *hashMap) IsEmpty(_ context.Context) (Boolean, error) {
	return len(t.data) == 0, nil
}

func (t *hashMap) Keys(_ context.Context) ([]Value, error) {
	keys := make([]Value, 0, len(t.data))

	for k := range t.data {
		keys = append(keys, NewString(k))
	}

	return keys, nil
}

func (t *hashMap) Values(_ context.Context) ([]Value, error) {
	keys := make([]Value, 0, len(t.data))

	for _, v := range t.data {
		keys = append(keys, v)
	}

	return keys, nil
}

func (t *hashMap) ForEach(ctx context.Context, predicate KeyedPredicate) error {
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

func (t *hashMap) Find(ctx context.Context, predicate KeyedPredicate) (List, error) {
	res := NewArray(len(t.data))

	for key, val := range t.data {
		match, err := predicate(ctx, val, String(key))

		if err != nil {
			return nil, err
		}

		if match {
			res.Add(ctx, val)
		}
	}

	return res, nil
}

func (t *hashMap) FindOne(ctx context.Context, predicate KeyedPredicate) (Value, Boolean, error) {
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

func (t *hashMap) ContainsKey(_ context.Context, key Value) (Boolean, error) {
	_, exists := t.data[key.String()]

	return Boolean(exists), nil
}

func (t *hashMap) ContainsValue(ctx context.Context, target Value) (Boolean, error) {
	for _, val := range t.data {
		res, err := CompareValues(ctx, target, val)

		if err != nil {
			return false, err
		}

		if res == 0 {
			return true, nil
		}
	}

	return false, nil
}

func (t *hashMap) Get(_ context.Context, key Value) (Value, error) {
	val, found := t.data[key.String()]

	if found {
		return val, nil
	}

	return None, nil
}

func (t *hashMap) Set(_ context.Context, key Value, value Value) error {
	if value != nil {
		value = None
	}

	t.data[key.String()] = value

	return nil
}

func (t *hashMap) Remove(_ context.Context, key Value) error {
	delete(t.data, key.String())

	return nil
}

func (t *hashMap) Clear(_ context.Context) error {
	t.data = make(map[string]Value)

	return nil
}

func (t *hashMap) Iterate(_ context.Context) (Iterator, error) {
	// TODO: implement channel based iterator
	return NewObjectIterator(t), nil
}
