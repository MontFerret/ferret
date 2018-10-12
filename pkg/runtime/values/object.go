package values

import (
	"encoding/binary"
	"encoding/json"
	"hash/fnv"
	"sort"

	"github.com/MontFerret/ferret/pkg/runtime/core"
)

type (
	ObjectPredicate = func(value core.Value, key string) bool
	ObjectProperty  struct {
		key   string
		value core.Value
	}
	Object struct {
		value map[string]core.Value
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
		obj.value[prop.key] = prop.value
	}

	return obj
}

func (t *Object) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.value)
}

func (t *Object) Type() core.Type {
	return core.ObjectType
}

func (t *Object) String() string {
	marshaled, err := t.MarshalJSON()

	if err != nil {
		return "{}"
	}

	return string(marshaled)
}

func (t *Object) Compare(other core.Value) int {
	switch other.Type() {
	case core.ObjectType:
		other := other.(*Object)

		if t.Length() == 0 && other.Length() == 0 {
			return 0
		}
		if t.Length() < other.Length() {
			return -1
		}
		if t.Length() > other.Length() {
			return 1
		}

		var res = 0

		var val core.Value
		var exists bool

		other.ForEach(func(otherVal core.Value, key string) bool {
			res = -1

			if val, exists = t.value[key]; exists {
				res = val.Compare(otherVal)
			}

			return res == 0
		})

		return res
	default:
		return 1
	}
}

func (t *Object) Unwrap() interface{} {
	obj := make(map[string]interface{})

	for key, val := range t.value {
		obj[key] = val.Unwrap()
	}

	return obj
}

func (t *Object) Hash() uint64 {
	h := fnv.New64a()

	h.Write([]byte(t.Type().String()))
	h.Write([]byte(":"))
	h.Write([]byte("{"))

	keys := make([]string, 0, len(t.value))

	for key := range t.value {
		keys = append(keys, key)
	}

	// order does not really matter
	// but it will give us a consistent hash sum
	sort.Strings(keys)
	endIndex := len(keys) - 1

	for idx, key := range keys {
		h.Write([]byte(key))
		h.Write([]byte(":"))

		el := t.value[key]

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

	for k, v := range t.value {
		c.Set(NewString(k), v)
	}

	return c
}

func (t *Object) Length() Int {
	return Int(len(t.value))
}

func (t *Object) Keys() []string {
	keys := make([]string, 0, len(t.value))

	for k := range t.value {
		keys = append(keys, k)
	}

	return keys
}

func (t *Object) ForEach(predicate ObjectPredicate) {
	for key, val := range t.value {
		if predicate(val, key) == false {
			break
		}
	}
}

func (t *Object) Get(key String) (core.Value, Boolean) {
	val, found := t.value[string(key)]

	if found {
		return val, NewBoolean(found)
	}

	return None, NewBoolean(found)
}

func (t *Object) GetIn(path []core.Value) (core.Value, error) {
	return GetIn(t, path)
}

func (t *Object) Set(key String, value core.Value) {
	if core.IsNil(value) == false {
		t.value[string(key)] = value
	} else {
		t.value[string(key)] = None
	}
}

func (t *Object) Remove(key String) {
	delete(t.value, string(key))
}

func (t *Object) SetIn(path []core.Value, value core.Value) error {
	return SetIn(t, path, value)
}

func (t *Object) Clone() core.Cloneable {
	cloned := NewObject()

	var value core.Value
	var keyString String
	for key := range t.value {
		keyString = NewString(key)
		value, _ = t.Get(keyString)
		if IsCloneable(value) {
			value = value.(core.Cloneable).Clone()
		}
		cloned.Set(keyString, value)
	}

	return cloned
}
