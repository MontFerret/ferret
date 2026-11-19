package data

import (
	"context"
	"errors"

	"github.com/MontFerret/ferret/v2/pkg/internal/debugmap"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

type FastObject struct {
	cache         *ShapeCache
	shape         *fastShape
	dict          *fastDictionary
	slots         []runtime.Value
	size          int
	dictThreshold int
}

// Keep unordered-table history with the dictionary so shape-backed VM objects
// pay no additional allocation or storage cost.
type fastDictionary struct {
	values map[string]runtime.Value
	bounds debugmap.Bounds
}

func NewFastObject(cache *ShapeCache, dictThreshold int) *FastObject {
	if cache == nil {
		cache = NewShapeCache(0)
	}

	if dictThreshold < 0 {
		dictThreshold = 0
	}

	return &FastObject{
		cache:         cache,
		shape:         cache.Root(),
		slots:         make([]runtime.Value, 0),
		dictThreshold: dictThreshold,
	}
}

func NewFastObjectOf(cache *ShapeCache, dictThreshold int, size int) *FastObject {
	obj := NewFastObject(cache, dictThreshold)

	if size > 0 {
		obj.slots = make([]runtime.Value, 0, size)
	}

	return obj
}

func (*FastObject) VMUntracked() {}

func (t *FastObject) Shape() *Shape {
	if t == nil || t.dict != nil {
		return nil
	}

	return t.shape
}

func (t *FastObject) ShapeID() uint64 {
	if t == nil || t.shape == nil || t.dict != nil {
		return 0
	}

	return t.shape.id
}

func (t *FastObject) LookupSlot(key string) (int, bool) {
	if t == nil || t.shape == nil || t.dict != nil {
		return 0, false
	}

	idx, ok := t.shape.fields[key]

	return idx, ok
}

func (t *FastObject) SlotValue(slot int) (runtime.Value, bool) {
	if t == nil || t.dict != nil || slot < 0 || slot >= len(t.slots) {
		return nil, false
	}

	val := t.slots[slot]
	if val == nil {
		return nil, false
	}

	return val, true
}

func (t *FastObject) SetSlotWithShape(next *Shape, slot int, value runtime.Value) bool {
	if t == nil || t.dict != nil || next == nil || slot < 0 {
		return false
	}

	if value == nil {
		value = runtime.None
	}

	if len(t.slots) < len(next.names) {
		slots := make([]runtime.Value, len(next.names))
		copy(slots, t.slots)
		t.slots = slots
	}

	if slot >= len(t.slots) {
		return false
	}

	t.shape = next

	if t.slots[slot] == nil {
		t.size++
	}

	t.slots[slot] = value

	return true
}

func (t *FastObject) SetStringCached(key string, value runtime.Value) (*Shape, *Shape, int, bool) {
	if t == nil {
		return nil, nil, 0, false
	}

	if value == nil {
		value = runtime.None
	}

	if t.dict != nil {
		t.dict.bounds.AddKey(key)
		t.dict.values[key] = value
		t.size = len(t.dict.values)

		return nil, nil, 0, false
	}

	shape := t.shape
	if shape == nil {
		return nil, nil, 0, false
	}

	if idx, ok := shape.fields[key]; ok {
		if t.slots[idx] == nil {
			t.size++
		}

		t.slots[idx] = value

		return shape, shape, idx, true
	}

	if t.dictThreshold > 0 && t.size+1 > t.dictThreshold {
		t.toDict()
		t.dict.bounds.AddKey(key)
		t.dict.values[key] = value
		t.size = len(t.dict.values)

		return nil, nil, 0, false
	}

	next := t.cache.Transition(shape, key)
	if next == nil {
		return nil, nil, 0, false
	}

	if len(t.slots) < len(next.names) {
		slots := make([]runtime.Value, len(next.names))
		copy(slots, t.slots)
		t.slots = slots
	}

	t.shape = next
	slot := next.fields[key]
	t.slots[slot] = value
	t.size++

	return shape, next, slot, true
}

func (t *FastObject) setString(key string, value runtime.Value) {
	t.SetStringCached(key, value)
}

func (t *FastObject) removeString(key string) {
	if t.dict != nil {
		if _, exists := t.dict.values[key]; exists {
			delete(t.dict.values, key)
			t.dict.bounds.Removed()

			t.size--
		}

		return
	}

	idx, ok := t.shape.fields[key]
	if !ok {
		return
	}

	if t.slots[idx] != nil {
		t.slots[idx] = nil
		t.size--
	}
}

func (t *FastObject) toDict() {
	if t.dict != nil {
		return
	}

	dict := &fastDictionary{values: make(map[string]runtime.Value, t.size)}
	for idx, key := range t.shape.names {
		if idx >= len(t.slots) {
			break
		}

		val := t.slots[idx]

		if val == nil {
			continue
		}

		dict.bounds.AddKey(key)
		dict.values[key] = val
	}

	t.dict = dict
	t.shape = nil
	t.slots = nil
	t.size = len(dict.values)
}

func (t *FastObject) len() int {
	if t.dict != nil {
		return len(t.dict.values)
	}

	return t.size
}

// DebugFieldsPageable checks the actual shape width, including tombstones.
func (t *FastObject) DebugFieldsPageable() bool {
	if t.dict != nil {
		return t.dict.bounds.Pageable(len(t.dict.values))
	}

	if t.shape == nil || len(t.shape.names) != len(t.slots) {
		return false
	}

	return t.size == len(t.slots) || len(t.slots) <= debugmap.MaxKeys
}

// ReadDebugFields pages dense shapes directly in slot order. A small sparse
// shape needs a bounded occupancy scan to locate the requested logical offset.
func (t *FastObject) ReadDebugFields(start, end int, visit func(string, runtime.Value) error) error {
	if t.dict != nil {
		return debugmap.Read(t.dict.values, t.dict.bounds, start, end, visit)
	}

	if !t.DebugFieldsPageable() || start < 0 || end < start || end > t.size {
		return errors.ErrUnsupported
	}

	if t.size == len(t.slots) {
		for i := start; i < end; i++ {
			if err := visit(t.shape.names[i], t.slots[i]); err != nil {
				return err
			}
		}

		return nil
	}

	pos := 0

	for slot, value := range t.slots {
		if pos == end {
			break
		}

		if value == nil {
			continue
		}

		if pos >= start {
			if err := visit(t.shape.names[slot], value); err != nil {
				return err
			}
		}

		pos++
	}

	return nil
}

func (t *FastObject) keys() []string {
	if t.dict != nil {
		keys := make([]string, 0, len(t.dict.values))

		for k := range t.dict.values {
			keys = append(keys, k)
		}

		return keys
	}

	keys := make([]string, 0, t.size)
	for idx, k := range t.shape.names {
		if t.slots[idx] == nil {
			continue
		}

		keys = append(keys, k)
	}

	return keys
}

func (t *FastObject) getByKey(key string) runtime.Value {
	if t.dict != nil {
		return t.dict.values[key]
	}

	if idx, ok := t.shape.fields[key]; ok {
		if idx < len(t.slots) {
			return t.slots[idx]
		}
	}

	return runtime.None
}

func (t *FastObject) getSlot(key string) (runtime.Value, bool) {
	if t.dict != nil {
		val, ok := t.dict.values[key]

		return val, ok
	}

	if idx, ok := t.shape.fields[key]; ok {
		val := t.slots[idx]
		if val == nil {
			return runtime.None, false
		}

		return val, true
	}

	return runtime.None, false
}

func (t *FastObject) forEachKV(fn func(key string, val runtime.Value)) {
	if t.dict != nil {
		for key, val := range t.dict.values {
			fn(key, val)
		}

		return
	}

	for idx, key := range t.shape.names {
		val := t.slots[idx]

		if val == nil {
			continue
		}

		fn(key, val)
	}
}

func (t *FastObject) toMap() runtime.Map {
	if t.dict != nil {
		return runtime.NewObjectWith(t.dict.values)
	}

	ctx := context.Background()
	obj := runtime.NewObject()

	for idx, key := range t.shape.names {
		if t.slots[idx] == nil {
			continue
		}

		_ = obj.Set(ctx, runtime.String(key), t.slots[idx])
	}

	return obj
}
