package data

import (
	"context"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

func (t *FastObject) Type() string {
	return "object"
}

func (t *FastObject) ObjectLike() {}

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

	return runtime.None, nil
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

func (t *FastObject) Empty(_ context.Context) (runtime.Map, error) {
	return NewFastObject(t.cache, t.dictThreshold), nil
}

func (t *FastObject) Merge(ctx context.Context, other runtime.Map) error {
	return other.ForEach(ctx, func(c context.Context, value, key runtime.Value) (runtime.Boolean, error) {
		t.setString(key.String(), value)

		return true, nil
	})
}
