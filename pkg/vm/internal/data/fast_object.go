package data

import (
	"context"
	"errors"

	"github.com/MontFerret/ferret/v2/pkg/internal/debugmap"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

// FastObject uses shape slots until growth selects its detached native dictionary.
type FastObject struct {
	cache *ShapeCache
	shape *fastShape
	dict  runtime.Object
	slots []runtime.Value
	size  int
}

// NewFastObject starts at the root shape of this cache.
func (c *ShapeCache) NewFastObject() *FastObject {
	return &FastObject{
		cache: c,
		shape: c.Root(),
	}
}

func (*FastObject) VMUntracked() {}

func (t *FastObject) Shape() *Shape {
	if t == nil {
		return nil
	}

	return t.shape
}

// ShapeIDFor returns a cache-local ID only to the cache that allocated it.
func (t *FastObject) ShapeIDFor(owner *ShapeCache) uint64 {
	if t == nil || t.shape == nil || t.cache != owner {
		return 0
	}

	return t.shape.id
}

func (t *FastObject) LookupSlot(key string) (int, bool) {
	if t == nil || t.shape == nil {
		return 0, false
	}

	idx, ok := t.shape.fields[key]

	return idx, ok
}

func (t *FastObject) SlotValue(slot int) (runtime.Value, bool) {
	if t == nil || t.shape == nil || slot < 0 || slot >= len(t.slots) {
		return nil, false
	}

	val := t.slots[slot]
	if val == nil {
		return nil, false
	}

	return val, true
}

func (t *FastObject) SetSlotWithShape(next *Shape, slot int, value runtime.Value) bool {
	if t == nil || t.shape == nil || next == nil || slot < 0 {
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

	shape := t.shape
	if shape == nil {
		if t.cache != nil {
			t.dict.SetString(key, value)
		}

		return nil, nil, 0, false
	}

	if idx, ok := shape.fields[key]; ok {
		slot := &t.slots[idx]
		if *slot == nil {
			t.size++
		}

		*slot = value

		return shape, shape, idx, true
	}

	if t.cache.dictThreshold > 0 && t.size+1 > t.cache.dictThreshold {
		t.toDict()
		t.dict.SetString(key, value)

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

func (t *FastObject) toDict() {
	dict := *runtime.NewObjectOf(t.size + 1)
	for idx, key := range t.shape.names {
		if idx >= len(t.slots) {
			break
		}

		val := t.slots[idx]

		if val == nil {
			continue
		}

		dict.SetString(key, val)
	}

	t.dict = dict
	t.shape = nil
	t.slots = nil
	t.size = 0
}

func (t *FastObject) len() int {
	if t.shape == nil {
		size, _ := t.dict.Length(context.Background())

		return int(size)
	}

	return t.size
}

// DebugFieldsPageable checks the actual shape width, including tombstones.
func (t *FastObject) DebugFieldsPageable() bool {
	if t.shape == nil {
		return t.dict.DebugFieldsPageable()
	}

	if len(t.shape.names) != len(t.slots) {
		return false
	}

	return t.size == len(t.slots) || len(t.slots) <= debugmap.MaxKeys
}

// ReadDebugFields pages dense shapes directly in slot order. A small sparse
// shape needs a bounded occupancy scan to locate the requested logical offset.
func (t *FastObject) ReadDebugFields(start, end int, visit func(string, runtime.Value) error) error {
	if t.shape == nil {
		return t.dict.ReadDebugFields(start, end, visit)
	}

	if !t.DebugFieldsPageable() || start < 0 || end < start || end > t.size {
		return errors.ErrUnsupported
	}

	names, slots := t.shape.names, t.slots
	if t.size == len(slots) {
		for i := start; i < end; i++ {
			if err := visit(names[i], slots[i]); err != nil {
				return err
			}
		}

		return nil
	}

	pos := 0

	for slot, value := range slots {
		if pos == end {
			break
		}

		if value == nil {
			continue
		}

		if pos >= start {
			if err := visit(names[slot], value); err != nil {
				return err
			}
		}

		pos++
	}

	return nil
}

func (t *FastObject) keys() []runtime.String {
	keys := make([]runtime.String, 0, t.len())
	t.forEachKV(func(key string, _ runtime.Value) {
		keys = append(keys, runtime.String(key))
	})

	return keys
}

// Snapshot keys are already addressable. Passing a pointer to the native
// dictionary avoids allocating an interface box for every compared property.
func (t *FastObject) getByKey(key *runtime.String) runtime.Value {
	if t.shape == nil {
		value, found, _ := t.dict.Lookup(context.Background(), key)
		if !found {
			return nil
		}

		return value
	}

	if idx, ok := t.shape.fields[string(*key)]; ok {
		if idx < len(t.slots) {
			return t.slots[idx]
		}
	}

	return runtime.None
}

func (t *FastObject) getShapeSlot(key string) (runtime.Value, bool) {
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
	if t.shape == nil {
		_ = t.dict.ForEach(context.Background(), func(_ context.Context, value, key runtime.Value) (runtime.Boolean, error) {
			fn(string(key.(runtime.String)), value)

			return true, nil
		})

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
	if t.shape == nil {
		return t.dict.Copy().(*runtime.Object)
	}

	obj := runtime.NewObject()

	for idx, key := range t.shape.names {
		if t.slots[idx] == nil {
			continue
		}

		obj.SetString(key, t.slots[idx])
	}

	return obj
}
