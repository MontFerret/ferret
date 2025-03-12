package internal

import (
	"context"
	"encoding/binary"
	"hash/fnv"
	"sort"

	"github.com/wI2L/jettison"

	"github.com/MontFerret/ferret/pkg/runtime/core"
)

type (
	ObjectProperty struct {
		key   string
		value core.Value
	}

	Object struct {
		data map[string]core.Value
	}
)

func NewObjectProperty(name string, value core.Value) *ObjectProperty {
	return &ObjectProperty{name, value}
}

func NewObject() core.Map {
	return &Object{make(map[string]core.Value)}
}

func NewObjectWith(props ...*ObjectProperty) core.Map {
	obj := &Object{make(map[string]core.Value)}

	for _, prop := range props {
		obj.data[prop.key] = prop.value
	}

	return obj
}

func (t *Object) Type() string {
	return "object"
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
func (t *Object) Compare(ctx context.Context, other core.Value) (int64, error) {
	otherObject, ok := other.(*Object)

	if !ok {
		// TODO: Maybe we should return an error here
		return core.CompareTypes(t, other), nil
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

	otherObject.ForEach(ctx, func(_ context.Context, value, k core.Value) (bool, error) {
		otherKeys = append(otherKeys, k.String())
		return true, nil
	})

	sortedOther := sort.StringSlice(otherKeys)
	sortedOther.Sort()

	var tVal, otherVal core.Value
	var tKey, otherKey string

	for i := 0; i < len(t.data) && res == 0; i++ {
		tKey, otherKey = sortedT[i], sortedOther[i]

		if tKey == otherKey {
			tVal, _ = t.Get(ctx, core.NewString(tKey))
			otherVal, _ = otherObject.Get(ctx, core.NewString(tKey))
			comp, err := core.CompareValues(ctx, tVal, otherVal)

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

func (t *Object) Copy() core.Value {
	c := &Object{make(map[string]core.Value)}

	for k, v := range t.data {
		c.data[k] = v
	}

	return c
}

func (t *Object) Clone(ctx context.Context) (core.Cloneable, error) {
	cloned := &Object{make(map[string]core.Value)}

	var value core.Value

	for key := range t.data {
		value, _ = t.data[key]

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

		cloned.data[key] = value
	}

	return cloned, nil
}

func (t *Object) Length(_ context.Context) (core.Int, error) {
	return core.Int(len(t.data)), nil
}

func (t *Object) IsEmpty(_ context.Context) (bool, error) {
	return len(t.data) == 0, nil
}

func (t *Object) Keys(_ context.Context) ([]core.Value, error) {
	keys := make([]core.Value, 0, len(t.data))

	for k := range t.data {
		keys = append(keys, core.NewString(k))
	}

	return keys, nil
}

func (t *Object) Values(_ context.Context) ([]core.Value, error) {
	keys := make([]core.Value, 0, len(t.data))

	for _, v := range t.data {
		keys = append(keys, v)
	}

	return keys, nil
}

func (t *Object) ForEach(ctx context.Context, predicate core.KeyedPredicate) error {
	for key, val := range t.data {
		doContinue, err := predicate(ctx, val, core.String(key))

		if err != nil {
			return err
		}

		if !doContinue {
			break
		}
	}

	return nil
}

func (t *Object) Find(ctx context.Context, predicate core.KeyedPredicate) (core.List, error) {
	res := NewArray(len(t.data))

	for key, val := range t.data {
		match, err := predicate(ctx, val, core.String(key))

		if err != nil {
			return nil, err
		}

		if match {
			res.Add(ctx, val)
		}
	}

	return res, nil
}

func (t *Object) FindOne(ctx context.Context, predicate core.KeyedPredicate) (core.Value, core.Boolean, error) {
	for key, val := range t.data {
		res, err := predicate(ctx, val, core.String(key))

		if err != nil {
			return core.None, false, err
		}

		if res {
			return val, true, nil
		}
	}

	return core.None, false, nil
}

func (t *Object) ContainsKey(_ context.Context, key core.Value) (core.Boolean, error) {
	_, exists := t.data[key.String()]

	return core.Boolean(exists), nil
}

func (t *Object) ContainsValue(ctx context.Context, target core.Value) (core.Boolean, error) {
	for _, val := range t.data {
		res, err := core.CompareValues(ctx, target, val)

		if err != nil {
			return false, err
		}

		if res == 0 {
			return true, nil
		}
	}

	return false, nil
}

func (t *Object) Get(_ context.Context, key core.Value) (core.Value, error) {
	val, found := t.data[key.String()]

	if found {
		return val, nil
	}

	return core.None, nil
}

func (t *Object) Set(_ context.Context, key core.Value, value core.Value) error {
	if value != nil {
		value = core.None
	}

	t.data[key.String()] = value

	return nil
}

func (t *Object) Remove(_ context.Context, key core.Value) error {
	delete(t.data, key.String())

	return nil
}

func (t *Object) Clear(_ context.Context) error {
	t.data = make(map[string]core.Value)

	return nil
}

func (t *Object) Iterate(_ context.Context) (core.Iterator, error) {
	// TODO: implement channel based iterator
	return NewObjectIterator(t), nil
}
