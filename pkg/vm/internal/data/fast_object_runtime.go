package data

import (
	"context"
	"encoding/binary"
	"hash/fnv"
	"slices"

	ferretencoding "github.com/MontFerret/ferret/v2/pkg/encoding"
	"github.com/MontFerret/ferret/v2/pkg/encoding/json"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

func (t *FastObject) MarshalJSON() ([]byte, error) {
	return ferretencoding.EncodeBytes(context.Background(), json.Default, t.toMap())
}

func (t *FastObject) String() string {
	marshaled, err := t.MarshalJSON()

	if err != nil {
		return "{}"
	}

	return string(marshaled)
}

func (t *FastObject) Equal(ctx context.Context, other runtime.Value) (bool, error) {
	otherObject, ok := other.(runtime.ObjectLike)
	if !ok {
		return false, nil
	}

	return equalObjectLike(ctx, t, otherObject)
}

func (t *FastObject) Compare(ctx context.Context, other runtime.Value) (runtime.Ordering, error) {
	otherObject, ok := other.(runtime.ObjectLike)
	if !ok {
		return runtime.Equal, runtime.Errorf(
			runtime.ErrInvalidOperation,
			"comparison cannot be applied to %s and %s",
			runtime.TypeName(t.Type()),
			runtime.TypeName(runtime.TypeOf(other)),
		)
	}

	return compareObjectLike(ctx, t, otherObject)
}

func (t *FastObject) Hash() uint64 {
	if t.shape == nil {
		return t.dict.Hash()
	}

	h := fnv.New64a()

	h.Write([]byte("object:"))
	h.Write([]byte("{"))

	keys := t.keys()
	slices.Sort(keys)
	endIndex := len(keys) - 1

	for idx, key := range keys {
		h.Write([]byte(key))
		h.Write([]byte(":"))

		el := t.getByKey(&keys[idx])

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

func (t *FastObject) Type() runtime.Type {
	return runtime.TypeObject
}

func (t *FastObject) ObjectLike() {}

func (t *FastObject) Copy() runtime.Value {
	if t.shape == nil {
		return &FastObject{cache: t.cache, dict: *t.dict.Copy().(*runtime.Object)}
	}

	slots := make([]runtime.Value, len(t.slots))
	copy(slots, t.slots)

	return &FastObject{
		cache: t.cache,
		shape: t.shape,
		slots: slots,
		size:  t.size,
	}
}

func (t *FastObject) Clone(ctx context.Context) (runtime.Cloneable, error) {
	if t.shape == nil {
		cloned, err := t.dict.Clone(ctx)
		if err != nil {
			return nil, err
		}
		return &FastObject{cache: t.cache, dict: *cloned.(*runtime.Object)}, nil
	}

	// Child cloning may clear or promote the source while these slots are read.
	shape := t.shape
	source := t.slots
	slots := make([]runtime.Value, len(source))
	size := 0

	for idx, value := range source {
		if value == nil {
			continue
		}

		clone, err := runtime.CloneOrCopy(ctx, value)
		if err != nil {
			return nil, err
		}

		if clone == nil {
			clone = runtime.None
		}
		slots[idx] = clone
		size++
	}

	return &FastObject{
		cache: t.cache,
		shape: shape,
		slots: slots,
		size:  size,
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
	if t.shape == nil {
		return t.dict.ForEach(ctx, predicate)
	}

	// Visit the original names. A callback may grow or rebuild the shape:
	// a matching slot supplies its current value, a detached slot its old value.
	shape, slots := t.shape, t.slots
	for idx, key := range shape.names {
		val := slots[idx]
		current := t.shape
		if current == shape || (current != nil && idx < len(current.names) && current.names[idx] == key) {
			val = t.slots[idx]
		}

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

func (t *FastObject) ContainsKey(ctx context.Context, key runtime.Value) (runtime.Boolean, error) {
	_, found, err := t.Lookup(ctx, key)
	return runtime.Boolean(found), err
}

func (t *FastObject) ContainsValue(ctx context.Context, target runtime.Value) (runtime.Boolean, error) {
	keys := t.keys()
	for i := range keys {
		equal, err := runtime.EqualValues(ctx, target, t.getByKey(&keys[i]))
		if err != nil {
			return runtime.False, err
		}

		if equal {
			return runtime.True, nil
		}
	}

	return runtime.False, nil
}

func (t *FastObject) Contains(ctx context.Context, target runtime.Value) (runtime.Boolean, error) {
	return t.ContainsValue(ctx, target)
}

func (t *FastObject) Get(ctx context.Context, key runtime.Value) (runtime.Value, error) {
	value, _, err := t.Lookup(ctx, key)
	return value, err
}

// Reuse native string boxes when forwarding an already resolved name.
func dictionaryKey(key runtime.Value, name string) runtime.Value {
	switch key.(type) {
	case runtime.String, *runtime.String:
		return key
	default:
		return runtime.String(name)
	}
}

func (t *FastObject) Lookup(ctx context.Context, key runtime.Value) (runtime.Value, bool, error) {
	// A host key can promote the object in String.
	name := key.String()
	if t.shape == nil {
		return t.dict.Lookup(ctx, dictionaryKey(key, name))
	}

	value, found := t.getShapeSlot(name)
	return value, found, nil
}

func (t *FastObject) Set(_ context.Context, key runtime.Value, value runtime.Value) error {
	t.setString(key.String(), value)

	return nil
}

func (t *FastObject) RemoveKey(ctx context.Context, key runtime.Value) error {
	name := key.String()
	if t.shape == nil {
		return t.dict.RemoveKey(ctx, dictionaryKey(key, name))
	}

	idx, ok := t.shape.fields[name]
	if ok && t.slots[idx] != nil {
		t.slots[idx] = nil
		t.size--
	}

	return nil
}

func (t *FastObject) Remove(ctx context.Context, value runtime.Value) error {
	keys := t.keys()
	for i := range keys {
		equal, err := runtime.EqualValues(ctx, value, t.getByKey(&keys[i]))
		if err != nil {
			return err
		}

		if equal {
			return t.RemoveKey(ctx, &keys[i])
		}
	}

	return nil
}

func (t *FastObject) Clear(ctx context.Context) error {
	if t.shape == nil {
		// Other Go value copies retain the old dictionary and its inspection history.
		return t.dict.Clear(ctx)
	}

	t.shape = t.cache.Root()
	t.slots = nil
	t.size = 0

	return nil
}

func (t *FastObject) Iterate(ctx context.Context) (runtime.Iterator, error) {
	if t.shape == nil {
		return t.dict.Iterate(ctx)
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

// New creates an independent empty collection with the receiver's configuration.
func (t *FastObject) New(ctx context.Context) (runtime.Map, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	cache := t.cache
	if cache == nil {
		cache = NewShapeCache(0, 0)
	}

	return cache.NewFastObject(), nil
}

func (t *FastObject) Merge(ctx context.Context, other runtime.Map) error {
	return other.ForEach(ctx, func(c context.Context, value, key runtime.Value) (runtime.Boolean, error) {
		t.setString(key.String(), value)

		return true, nil
	})
}
