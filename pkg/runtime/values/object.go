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
	ObjectPredicate = func(value core.Value, key string) bool

	ObjectProperty struct {
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
	return jettison.MarshalOpts(t.value, jettison.NoHTMLEscaping())
}

func (t *Object) Type() core.Type {
	return types.Object
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
	if other.Type() == t.Type() {
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

		var res int64

		tKeys := make([]string, 0, len(t.value))

		for k := range t.value {
			tKeys = append(tKeys, k)
		}

		sortedT := sort.StringSlice(tKeys)
		sortedT.Sort()

		otherKeys := make([]string, 0, other.Length())

		other.ForEach(func(value core.Value, k string) bool {
			otherKeys = append(otherKeys, k)
			return true
		})

		sortedOther := sort.StringSlice(otherKeys)
		sortedOther.Sort()

		var tVal, otherVal core.Value
		var tKey, otherKey string

		for i := 0; i < len(t.value) && res == 0; i++ {
			tKey, otherKey = sortedT[i], sortedOther[i]

			if tKey == otherKey {
				tVal, _ = t.Get(NewString(tKey))
				otherVal, _ = other.Get(NewString(tKey))
				res = tVal.Compare(otherVal)

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

	return types.Compare(types.Object, other.Type())
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

func (t *Object) Keys() []String {
	keys := make([]String, 0, len(t.value))

	for k := range t.value {
		keys = append(keys, NewString(k))
	}

	return keys
}

func (t *Object) Values() []core.Value {
	keys := make([]core.Value, 0, len(t.value))

	for _, v := range t.value {
		keys = append(keys, v)
	}

	return keys
}

func (t *Object) ForEach(predicate ObjectPredicate) {
	for key, val := range t.value {
		if !predicate(val, key) {
			break
		}
	}
}

func (t *Object) Find(predicate ObjectPredicate) (core.Value, Boolean) {
	for idx, val := range t.value {
		if predicate(val, idx) {
			return val, True
		}
	}

	return None, False
}

func (t *Object) Has(key String) Boolean {
	_, exists := t.value[string(key)]

	return NewBoolean(exists)
}

func (t *Object) MustGet(key String) core.Value {
	val, _ := t.Get(key)

	return val
}

func (t *Object) MustGetOr(key String, defaultValue core.Value) core.Value {
	val, found := t.value[string(key)]

	if found {
		return val
	}

	return defaultValue
}

func (t *Object) Get(key String) (core.Value, Boolean) {
	val, found := t.value[string(key)]

	if found {
		return val, NewBoolean(found)
	}

	return None, NewBoolean(found)
}

func (t *Object) Set(key String, value core.Value) {
	if value != nil {
		t.value[string(key)] = value
	} else {
		t.value[string(key)] = None
	}
}

func (t *Object) Remove(key String) {
	delete(t.value, string(key))
}

func (t *Object) Clone() core.Cloneable {
	cloned := NewObject()

	var value core.Value
	var keyString String

	for key := range t.value {
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

func (t *Object) GetIn(ctx context.Context, path []core.Value) (core.Value, core.PathError) {
	if len(path) == 0 {
		return None, nil
	}

	segmentIdx := 0
	first, _ := t.Get(ToString(path[segmentIdx]))

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
