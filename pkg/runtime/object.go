package runtime

import (
	"context"
	"encoding/binary"
	"hash/fnv"
	"sort"
	"sync/atomic"

	"github.com/wI2L/jettison"
)

type (
	objectShape struct {
		id          uint64
		fields      map[string]int
		names       []string
		transitions map[string]*objectShape
	}

	ObjectProperty struct {
		key   string
		value Value
	}

	Object struct {
		shape *objectShape
		slots []Value
		size  int
	}
)

var (
	shapeIDCounter   uint64
	emptyObjectShape = newObjectShape(nil, nil)
)

func nextShapeID() uint64 {
	return atomic.AddUint64(&shapeIDCounter, 1)
}

func newObjectShape(fields map[string]int, names []string) *objectShape {
	if fields == nil {
		fields = make(map[string]int)
	}

	if names == nil {
		names = make([]string, 0)
	}

	return &objectShape{
		id:          nextShapeID(),
		fields:      fields,
		names:       names,
		transitions: make(map[string]*objectShape),
	}
}

func (s *objectShape) transition(key string) *objectShape {
	if next, ok := s.transitions[key]; ok {
		return next
	}

	fields := make(map[string]int, len(s.fields)+1)
	for k, v := range s.fields {
		fields[k] = v
	}

	slot := len(s.names)
	fields[key] = slot

	names := make([]string, slot+1)
	copy(names, s.names)
	names[slot] = key

	next := newObjectShape(fields, names)
	s.transitions[key] = next

	return next
}

func NewObjectProperty(name string, value Value) *ObjectProperty {
	return &ObjectProperty{name, value}
}

func NewObject() *Object {
	return &Object{shape: emptyObjectShape, slots: make([]Value, 0)}
}

func NewObjectOf(size int) *Object {
	return &Object{shape: emptyObjectShape, slots: make([]Value, 0, size)}
}

func NewObjectWith(props ...*ObjectProperty) *Object {
	obj := &Object{shape: emptyObjectShape, slots: make([]Value, 0, len(props))}

	for _, prop := range props {
		obj.setString(prop.key, prop.value)
	}

	return obj
}

func (t *Object) ShapeID() uint64 {
	if t == nil || t.shape == nil {
		return 0
	}

	return t.shape.id
}

func (t *Object) LookupSlot(key string) (int, bool) {
	if t == nil || t.shape == nil {
		return 0, false
	}

	idx, ok := t.shape.fields[key]
	return idx, ok
}

func (t *Object) SlotValue(slot int) (Value, bool) {
	if t == nil || slot < 0 || slot >= len(t.slots) {
		return nil, false
	}

	val := t.slots[slot]
	if val == nil {
		return nil, false
	}

	return val, true
}

func (t *Object) setString(key string, value Value) {
	if value == nil {
		value = None
	}

	if idx, ok := t.shape.fields[key]; ok {
		if t.slots[idx] == nil {
			t.size++
		}

		t.slots[idx] = value
		return
	}

	t.shape = t.shape.transition(key)
	t.slots = append(t.slots, value)
	t.size++
}

func (t *Object) Type() string {
	return "object"
}

func (t *Object) MarshalJSON() ([]byte, error) {
	obj := make(map[string]Value, t.size)

	for idx, key := range t.shape.names {
		val := t.slots[idx]
		if val == nil {
			continue
		}

		obj[key] = val
	}

	return jettison.MarshalOpts(obj, jettison.NoHTMLEscaping())
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
func (t *Object) Compare(other Value) int64 {
	otherObject, ok := other.(*Object)

	if !ok {
		return CompareTypes(t, other)
	}

	size := t.size
	otherSize := otherObject.size

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

	for idx, k := range t.shape.names {
		if t.slots[idx] != nil {
			tKeys = append(tKeys, k)
		}
	}

	sortedT := sort.StringSlice(tKeys)
	sortedT.Sort()

	otherKeys := make([]string, 0, otherSize)

	for idx, k := range otherObject.shape.names {
		if otherObject.slots[idx] != nil {
			otherKeys = append(otherKeys, k)
		}
	}

	sortedOther := sort.StringSlice(otherKeys)
	sortedOther.Sort()

	var tVal, otherVal Value
	var tKey, otherKey string

	for i := 0; i < len(tKeys) && res == 0; i++ {
		tKey, otherKey = sortedT[i], sortedOther[i]

		if tKey == otherKey {
			tVal = t.slots[t.shape.fields[tKey]]
			otherVal = otherObject.slots[otherObject.shape.fields[tKey]]
			res = CompareValues(tVal, otherVal)

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
	obj := make(map[string]interface{}, t.size)

	for idx, key := range t.shape.names {
		val := t.slots[idx]
		if val == nil {
			continue
		}

		obj[key] = val.Unwrap()
	}

	return obj
}

func (t *Object) Hash() uint64 {
	h := fnv.New64a()

	h.Write([]byte("object:"))
	h.Write([]byte("{"))

	keys := make([]string, 0, t.size)

	for idx, key := range t.shape.names {
		if t.slots[idx] != nil {
			keys = append(keys, key)
		}
	}

	// order does not really matter
	// but it will give us a consistent hash sum
	sort.Strings(keys)
	endIndex := len(keys) - 1

	for idx, key := range keys {
		h.Write([]byte(key))
		h.Write([]byte(":"))

		el := t.slots[t.shape.fields[key]]

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

func (t *Object) Copy() Value {
	slots := make([]Value, len(t.slots))
	copy(slots, t.slots)

	return &Object{
		shape: t.shape,
		slots: slots,
		size:  t.size,
	}
}

func (t *Object) Clone(ctx context.Context) (Cloneable, error) {
	slots := make([]Value, len(t.slots))

	for idx, value := range t.slots {
		if value == nil {
			continue
		}

		cloneable, ok := value.(Cloneable)
		if ok {
			clone, err := cloneable.Clone(ctx)
			if err != nil {
				return nil, err
			}
			slots[idx] = clone
			continue
		}

		slots[idx] = value.Copy()
	}

	return &Object{
		shape: t.shape,
		slots: slots,
		size:  t.size,
	}, nil
}

func (t *Object) Length(_ context.Context) (Int, error) {
	return Int(t.size), nil
}

func (t *Object) IsEmpty(_ context.Context) (Boolean, error) {
	return t.size == 0, nil
}

func (t *Object) Keys(_ context.Context) (List, error) {
	keys := make([]Value, 0, t.size)

	for idx, k := range t.shape.names {
		if t.slots[idx] == nil {
			continue
		}
		keys = append(keys, NewString(k))
	}

	return NewArrayOf(keys), nil
}

func (t *Object) Values(_ context.Context) (List, error) {
	values := make([]Value, 0, t.size)

	for _, v := range t.slots {
		if v == nil {
			continue
		}
		values = append(values, v)
	}

	return NewArrayOf(values), nil
}

func (t *Object) ForEach(ctx context.Context, predicate KeyedPredicate) error {
	for idx, key := range t.shape.names {
		val := t.slots[idx]
		if val == nil {
			continue
		}

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

func (t *Object) Find(ctx context.Context, predicate KeyedPredicate) (List, error) {
	res := NewArray(t.size)

	for idx, key := range t.shape.names {
		val := t.slots[idx]
		if val == nil {
			continue
		}

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

func (t *Object) FindOne(ctx context.Context, predicate KeyedPredicate) (Value, Boolean, error) {
	for idx, key := range t.shape.names {
		val := t.slots[idx]
		if val == nil {
			continue
		}

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

func (t *Object) ContainsKey(_ context.Context, key Value) (Boolean, error) {
	idx, ok := t.shape.fields[key.String()]
	if !ok {
		return false, nil
	}

	return Boolean(t.slots[idx] != nil), nil
}

func (t *Object) ContainsValue(_ context.Context, target Value) (Boolean, error) {
	for _, val := range t.slots {
		if val == nil {
			continue
		}

		res := CompareValues(target, val)

		if res == 0 {
			return true, nil
		}
	}

	return false, nil
}

func (t *Object) Get(_ context.Context, key Value) (Value, error) {
	idx, ok := t.shape.fields[key.String()]
	if !ok {
		return None, ErrNotFound
	}

	val := t.slots[idx]
	if val == nil {
		return None, ErrNotFound
	}

	return val, nil
}

func (t *Object) Set(_ context.Context, key Value, value Value) error {
	t.setString(key.String(), value)
	return nil
}

func (t *Object) Remove(_ context.Context, key Value) error {
	idx, ok := t.shape.fields[key.String()]
	if !ok {
		return nil
	}

	if t.slots[idx] != nil {
		t.slots[idx] = nil
		t.size--
	}

	return nil
}

func (t *Object) Clear(_ context.Context) error {
	t.shape = emptyObjectShape
	t.slots = nil
	t.size = 0

	return nil
}

func (t *Object) Iterate(_ context.Context) (Iterator, error) {
	// TODO: implement channel based iterator
	return NewObjectIterator(t), nil
}
