package data

import (
	"context"
	"encoding/binary"
	"hash/fnv"
	"sort"

	"github.com/MontFerret/ferret/v2/pkg/encoding/json"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

type (
	FastObject struct {
		cache         *ShapeCache
		shape         *fastShape
		slots         []runtime.Value
		size          int
		dict          map[string]runtime.Value
		dictThreshold int
	}

	fastShape struct {
		id     uint64
		fields map[string]int
		names  []string
	}

	ShapeCache struct {
		limit       int
		nextID      uint64
		transitions map[shapeKey]*fastShape
		root        *fastShape
	}

	shapeKey struct {
		shapeID uint64
		key     string
	}
)

type Shape = fastShape

func NewShapeCache(limit int) *ShapeCache {
	if limit < 0 {
		limit = 0
	}

	cache := &ShapeCache{
		limit:       limit,
		transitions: make(map[shapeKey]*fastShape),
	}
	cache.root = cache.newShape(nil, nil)

	return cache
}

func (c *ShapeCache) Root() *Shape {
	if c == nil {
		return nil
	}

	return c.root
}

func (c *ShapeCache) Transition(shape *Shape, key string) *Shape {
	if c == nil || shape == nil {
		return nil
	}

	if c.limit > 0 {
		k := shapeKey{shapeID: shape.id, key: key}
		if next, ok := c.transitions[k]; ok {
			return next
		}

		if len(c.transitions) >= c.limit {
			// Hard cap: stop caching new transitions.
			return c.newShapeFrom(shape, key)
		}

		next := c.newShapeFrom(shape, key)
		c.transitions[k] = next
		return next
	}

	return c.newShapeFrom(shape, key)
}

func (c *ShapeCache) nextShapeID() uint64 {
	c.nextID++
	return c.nextID
}

func (c *ShapeCache) newShape(fields map[string]int, names []string) *fastShape {
	if fields == nil {
		fields = make(map[string]int)
	}

	if names == nil {
		names = make([]string, 0)
	}

	return &fastShape{
		id:     c.nextShapeID(),
		fields: fields,
		names:  names,
	}
}

func (c *ShapeCache) newShapeFrom(prev *fastShape, key string) *fastShape {
	fields := make(map[string]int, len(prev.fields)+1)
	for k, v := range prev.fields {
		fields[k] = v
	}

	slot := len(prev.names)
	fields[key] = slot

	names := make([]string, slot+1)
	copy(names, prev.names)
	names[slot] = key

	return c.newShape(fields, names)
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

func (t *FastObject) Type() string {
	return "object"
}

func (t *FastObject) ObjectLike() {}

func (t *FastObject) MarshalJSON() ([]byte, error) {
	return json.Default.Encode(t.toMap())
}

func (t *FastObject) String() string {
	marshaled, err := t.MarshalJSON()

	if err != nil {
		return "{}"
	}

	return string(marshaled)
}

func (t *FastObject) Compare(other runtime.Value) int {
	otherObject, ok := other.(*FastObject)

	if !ok {
		if otherLike, ok := other.(runtime.ObjectLike); ok {
			return runtime.CompareTypes(t, otherLike)
		}

		return runtime.CompareTypes(t, other)
	}

	size := t.len()
	otherSize := otherObject.len()

	if size == 0 && otherSize == 0 {
		return 0
	}

	if size < otherSize {
		return -1
	}

	if size > otherSize {
		return 1
	}

	tKeys := t.keys()
	sort.Strings(tKeys)

	otherKeys := otherObject.keys()
	sort.Strings(otherKeys)

	var res int

	for i := 0; i < len(tKeys) && res == 0; i++ {
		tKey, otherKey := tKeys[i], otherKeys[i]

		if tKey == otherKey {
			tVal := t.getByKey(tKey)
			otherVal := otherObject.getByKey(otherKey)
			res = runtime.CompareValues(tVal, otherVal)
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

func (t *FastObject) Hash() uint64 {
	h := fnv.New64a()

	h.Write([]byte("object:"))
	h.Write([]byte("{"))

	keys := t.keys()
	sort.Strings(keys)
	endIndex := len(keys) - 1

	for idx, key := range keys {
		h.Write([]byte(key))
		h.Write([]byte(":"))

		el := t.getByKey(key)

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

func (t *FastObject) Copy() runtime.Value {
	if t.dict != nil {
		dict := make(map[string]runtime.Value, len(t.dict))
		for k, v := range t.dict {
			dict[k] = v
		}

		return &FastObject{
			cache:         t.cache,
			shape:         t.shape,
			slots:         nil,
			size:          len(dict),
			dict:          dict,
			dictThreshold: t.dictThreshold,
		}
	}

	slots := make([]runtime.Value, len(t.slots))
	copy(slots, t.slots)

	return &FastObject{
		cache:         t.cache,
		shape:         t.shape,
		slots:         slots,
		size:          t.size,
		dictThreshold: t.dictThreshold,
	}
}

func (t *FastObject) Clone(ctx context.Context) (runtime.Cloneable, error) {
	if t.dict != nil {
		cloned := make(map[string]runtime.Value, len(t.dict))

		for key, value := range t.dict {
			if cloneable, ok := value.(runtime.Cloneable); ok {
				clone, err := cloneable.Clone(ctx)
				if err != nil {
					return nil, err
				}
				cloned[key] = clone
				continue
			}

			cloned[key] = value.Copy()
		}

		return &FastObject{
			cache:         t.cache,
			shape:         t.shape,
			dict:          cloned,
			size:          len(cloned),
			dictThreshold: t.dictThreshold,
		}, nil
	}

	slots := make([]runtime.Value, len(t.slots))

	for idx, value := range t.slots {
		if value == nil {
			continue
		}

		if cloneable, ok := value.(runtime.Cloneable); ok {
			clone, err := cloneable.Clone(ctx)
			if err != nil {
				return nil, err
			}
			slots[idx] = clone
			continue
		}

		slots[idx] = value.Copy()
	}

	return &FastObject{
		cache:         t.cache,
		shape:         t.shape,
		slots:         slots,
		size:          t.size,
		dictThreshold: t.dictThreshold,
	}, nil
}

func (t *FastObject) Length(_ context.Context) (runtime.Int, error) {
	return runtime.Int(t.len()), nil
}

func (t *FastObject) IsEmpty(_ context.Context) (runtime.Boolean, error) {
	return t.len() == 0, nil
}

func (t *FastObject) Keys(_ context.Context) (runtime.List, error) {
	keys := make([]runtime.Value, 0, t.len())
	t.forEachKV(func(key string, _ runtime.Value) {
		keys = append(keys, runtime.NewString(key))
	})

	return runtime.NewArrayOf(keys), nil
}

func (t *FastObject) Values(_ context.Context) (runtime.List, error) {
	values := make([]runtime.Value, 0, t.len())
	t.forEachKV(func(_ string, val runtime.Value) {
		values = append(values, val)
	})

	return runtime.NewArrayOf(values), nil
}

func (t *FastObject) ForEach(ctx context.Context, predicate runtime.KeyReadablePredicate) error {
	if t.dict != nil {
		for key, val := range t.dict {
			doContinue, err := predicate(ctx, val, runtime.String(key))

			if err != nil {
				return err
			}

			if !doContinue {
				break
			}
		}

		return nil
	}

	for idx, key := range t.shape.names {
		val := t.slots[idx]
		if val == nil {
			continue
		}

		doContinue, err := predicate(ctx, val, runtime.String(key))

		if err != nil {
			return err
		}

		if !doContinue {
			break
		}
	}

	return nil
}

func (t *FastObject) Filter(ctx context.Context, predicate runtime.KeyReadablePredicate) (runtime.List, error) {
	res := runtime.NewArray(t.len())

	err := t.ForEach(ctx, func(c context.Context, value, key runtime.Value) (runtime.Boolean, error) {
		match, err := predicate(c, value, key)

		if err != nil {
			return runtime.False, err
		}

		if match {
			if err := res.Append(c, value); err != nil {
				return runtime.False, err
			}
		}

		return true, nil
	})

	if err != nil {
		return nil, err
	}

	return res, nil
}

func (t *FastObject) Find(ctx context.Context, predicate runtime.KeyReadablePredicate) (runtime.Value, runtime.Boolean, error) {
	var out runtime.Value
	var ok bool

	err := t.ForEach(ctx, func(c context.Context, value, key runtime.Value) (runtime.Boolean, error) {
		match, err := predicate(c, value, key)

		if err != nil {
			return runtime.False, err
		}

		if match {
			out = value
			ok = true
			return false, nil
		}

		return true, nil
	})

	if err != nil {
		return runtime.None, false, err
	}

	if ok {
		return out, true, nil
	}

	return runtime.None, false, nil
}

func (t *FastObject) ContainsKey(_ context.Context, key runtime.Value) (runtime.Boolean, error) {
	_, exists := t.getSlot(key.String())

	return runtime.Boolean(exists), nil
}

func (t *FastObject) ContainsValue(_ context.Context, target runtime.Value) (runtime.Boolean, error) {
	found := false

	t.forEachKV(func(_ string, val runtime.Value) {
		if runtime.CompareValues(target, val) == 0 {
			found = true
		}
	})

	return runtime.Boolean(found), nil
}

func (t *FastObject) Contains(ctx context.Context, target runtime.Value) (runtime.Boolean, error) {
	return t.ContainsValue(ctx, target)
}

func (t *FastObject) Get(_ context.Context, key runtime.Value) (runtime.Value, error) {
	val, ok := t.getSlot(key.String())

	if ok {
		return val, nil
	}

	return runtime.None, runtime.ErrNotFound
}

func (t *FastObject) Set(_ context.Context, key runtime.Value, value runtime.Value) error {
	t.setString(key.String(), value)

	return nil
}

func (t *FastObject) RemoveKey(_ context.Context, key runtime.Value) error {
	t.removeString(key.String())

	return nil
}

func (t *FastObject) Remove(_ context.Context, value runtime.Value) error {
	if t.dict != nil {
		for key, val := range t.dict {
			if runtime.CompareValues(value, val) == 0 {
				t.removeString(key)
				break
			}
		}

		return nil
	}

	if t.shape == nil {
		return nil
	}

	for idx, key := range t.shape.names {
		if idx >= len(t.slots) {
			break
		}

		val := t.slots[idx]
		if val == nil {
			continue
		}

		if runtime.CompareValues(value, val) == 0 {
			t.removeString(key)
			break
		}
	}

	return nil
}

func (t *FastObject) Clear(_ context.Context) error {
	if t.dict != nil {
		t.dict = make(map[string]runtime.Value)
		t.size = 0

		return nil
	}

	t.shape = t.cache.Root()
	t.slots = nil
	t.size = 0

	return nil
}

func (t *FastObject) Iterate(_ context.Context) (runtime.Iterator, error) {
	if t.dict != nil {
		keys := make([]string, 0, len(t.dict))

		for key := range t.dict {
			keys = append(keys, key)
		}

		return &fastObjectDictIterator{
			keys: keys,
			dict: t.dict,
		}, nil
	}

	entries := make([]fastObjectEntry, 0, t.size)

	for idx, key := range t.shape.names {
		if t.slots[idx] == nil {
			continue
		}

		entries = append(entries, fastObjectEntry{key: key, slot: idx})
	}

	return &fastObjectIterator{entries: entries, slots: t.slots}, nil
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

func (t *FastObject) Empty(_ context.Context) (runtime.Map, error) {
	return NewFastObject(t.cache, t.dictThreshold), nil
}

func (t *FastObject) Merge(ctx context.Context, other runtime.Map) error {
	return other.ForEach(ctx, func(c context.Context, value, key runtime.Value) (runtime.Boolean, error) {
		t.setString(key.String(), value)

		return true, nil
	})
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
