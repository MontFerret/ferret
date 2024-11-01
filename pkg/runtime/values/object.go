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
	ObjectPredicate = func(value core.Value, key string) bool

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

func NewObject() *Object {
	return &Object{make(map[string]core.Value)}
}

func NewObjectWith(props ...*ObjectProperty) *Object {
	obj := NewObject()

	for _, prop := range props {
		obj.data[prop.key] = prop.value
	}

	return obj
}

func (t *Object) Type() core.Type {
	return types.Object
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
func (t *Object) Compare(other core.Value) int64 {
	otherObject, ok := other.(*Object)

	if !ok {
		return types.Compare(types.Object, core.Reflect(other))
	}

	if t.Length() == 0 && otherObject.Length() == 0 {
		return 0
	}

	if t.Length() < otherObject.Length() {
		return -1
	}

	if t.Length() > otherObject.Length() {
		return 1
	}

	var res int64

	tKeys := make([]string, 0, len(t.data))

	for k := range t.data {
		tKeys = append(tKeys, k)
	}

	sortedT := sort.StringSlice(tKeys)
	sortedT.Sort()

	otherKeys := make([]string, 0, otherObject.Length())

	otherObject.ForEach(func(value core.Value, k string) bool {
		otherKeys = append(otherKeys, k)
		return true
	})

	sortedOther := sort.StringSlice(otherKeys)
	sortedOther.Sort()

	var tVal, otherVal core.Value
	var tKey, otherKey string

	for i := 0; i < len(t.data) && res == 0; i++ {
		tKey, otherKey = sortedT[i], sortedOther[i]

		if tKey == otherKey {
			tVal, _ = t.Get(NewString(tKey))
			otherVal, _ = otherObject.Get(NewString(tKey))
			res = Compare(tVal, otherVal)

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

	h.Write([]byte(types.Object.String()))
	h.Write([]byte(":"))
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
	c := NewObject()

	for k, v := range t.data {
		c.Set(NewString(k), v)
	}

	return c
}

func (t *Object) Length() Int {
	return Int(len(t.data))
}

func (t *Object) Keys() []String {
	keys := make([]String, 0, len(t.data))

	for k := range t.data {
		keys = append(keys, NewString(k))
	}

	return keys
}

func (t *Object) Values() []core.Value {
	keys := make([]core.Value, 0, len(t.data))

	for _, v := range t.data {
		keys = append(keys, v)
	}

	return keys
}

func (t *Object) ForEach(predicate ObjectPredicate) {
	for key, val := range t.data {
		if !predicate(val, key) {
			break
		}
	}
}

func (t *Object) Find(predicate ObjectPredicate) (core.Value, Boolean) {
	for idx, val := range t.data {
		if predicate(val, idx) {
			return val, True
		}
	}

	return None, False
}

func (t *Object) Has(key String) Boolean {
	_, exists := t.data[string(key)]

	return NewBoolean(exists)
}

func (t *Object) MustGet(key String) core.Value {
	val, _ := t.Get(key)

	return val
}

func (t *Object) MustGetOr(key String, defaultValue core.Value) core.Value {
	val, found := t.data[string(key)]

	if found {
		return val
	}

	return defaultValue
}

func (t *Object) Get(key String) (core.Value, Boolean) {
	val, found := t.data[string(key)]

	if found {
		return val, NewBoolean(found)
	}

	return None, NewBoolean(found)
}

func (t *Object) Set(key String, value core.Value) {
	if value != nil {
		t.data[string(key)] = value
	} else {
		t.data[string(key)] = None
	}
}

func (t *Object) Remove(key String) {
	delete(t.data, string(key))
}

func (t *Object) Clone() core.Cloneable {
	cloned := NewObject()

	var value core.Value
	var keyString String

	for key := range t.data {
		keyString = NewString(key)
		value, _ = t.Get(keyString)

		cloneable, ok := value.(core.Cloneable)

		if ok {
			value = cloneable.Clone()
		}
		cloned.Set(keyString, value)
	}

	return cloned
}

func (t *Object) Iterate(_ context.Context) (core.Iterator, error) {
	// TODO: implement channel based iterator
	return NewObjectIterator(t), nil
}
