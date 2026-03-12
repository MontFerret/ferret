package data

import (
	"context"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

type FastObject struct {
	cache         *ShapeCache
	shape         *fastShape
	dict          map[string]runtime.Value
	slots         []runtime.Value
	size          int
	dictThreshold int
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
		if _, exists := t.dict[key]; !exists {
			t.size++
		}

		t.dict[key] = value

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
		t.dict[key] = value
		t.size = len(t.dict)

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
		if _, exists := t.dict[key]; exists {
			delete(t.dict, key)

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

	dict := make(map[string]runtime.Value, t.size)
	for idx, key := range t.shape.names {
		if idx >= len(t.slots) {
			break
		}

		val := t.slots[idx]

		if val == nil {
			continue
		}

		dict[key] = val
	}

	t.dict = dict
	t.shape = nil
	t.slots = nil
	t.size = len(dict)
}

func (t *FastObject) len() int {
	if t.dict != nil {
		return len(t.dict)
	}

	return t.size
}

func (t *FastObject) keys() []string {
	if t.dict != nil {
		keys := make([]string, 0, len(t.dict))

		for k := range t.dict {
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
		return t.dict[key]
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
		val, ok := t.dict[key]

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
		for key, val := range t.dict {
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
		return runtime.NewObjectWith(t.dict)
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
