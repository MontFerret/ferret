package runtime

import (
	"context"
	"encoding/binary"
	"errors"
	"hash/fnv"
	"sort"
	"strings"

	"github.com/MontFerret/ferret/v2/pkg/internal/debugmap"
)

type Object struct {
	data map[string]Value
	// Go value copies share the backing map and its scan history.
	bounds *debugmap.Bounds
}

func NewObject() *Object {
	return NewObjectOf(0)
}

func NewObjectOf(size int) *Object {
	o := &Object{data: make(map[string]Value, size), bounds: new(debugmap.Bounds)}
	o.bounds.Reserve(size)

	return o
}

// NewObjectWith creates a new object with the provided properties.
// The properties are copied to ensure that the original map can be safely modified without affecting the created object.
func NewObjectWith(props map[string]Value) *Object {
	o := NewObjectOf(len(props))

	for k, v := range props {
		o.setKey(k, v)
	}

	return o
}

func (t *Object) Type() Type {
	return TypeObject
}

func (t *Object) ObjectLike() {}

func (t *Object) String() string {
	var b strings.Builder
	writeObject(&b, t)

	return b.String()
}

func (t *Object) Equal(ctx context.Context, other Value) (bool, error) {
	otherMap, ok := other.(Map)
	if !ok {
		return false, nil
	}

	otherSize, err := otherMap.Length(ctx)
	if err != nil {
		return false, err
	}

	if Int(len(t.data)) != otherSize {
		return false, nil
	}

	for key, value := range t.data {
		otherValue, found, err := otherMap.Lookup(ctx, String(key))
		if err != nil {
			return false, err
		}

		if !found {
			return false, nil
		}

		equal, err := EqualValues(ctx, value, otherValue)
		if err != nil {
			return false, err
		}

		if !equal {
			return false, nil
		}
	}

	return true, nil
}

// Compare preserves Ferret's structural object ordering, including its
// historical descending comparison for the first distinct sorted key.
func (t *Object) Compare(ctx context.Context, other Value) (Ordering, error) {
	otherMap, ok := other.(Map)
	if !ok {
		return Equal, incompatibleComparisonError(t, other)
	}

	otherSize, err := otherMap.Length(ctx)
	if err != nil {
		return Equal, err
	}

	size := Int(len(t.data))
	if size < otherSize {
		return Less, nil
	}

	if size > otherSize {
		return Greater, nil
	}

	leftKeys := t.sortedKeys()
	var rightKeys []string
	otherObject, nativeObject := otherMap.(*Object)

	if nativeObject {
		rightKeys = otherObject.sortedKeys()
	} else {
		rightKeys, err = sortedMapKeys(ctx, otherMap)
		if err != nil {
			return Equal, err
		}
	}

	for idx, leftKey := range leftKeys {
		rightKey := rightKeys[idx]
		if leftKey != rightKey {
			if leftKey < rightKey {
				return Greater, nil
			}

			return Less, nil
		}

		var rightValue Value
		if nativeObject {
			rightValue = otherObject.data[rightKey]
		} else {
			var found bool
			rightValue, found, err = otherMap.Lookup(ctx, String(rightKey))

			if err != nil {
				return Equal, err
			}

			if !found {
				return Equal, Errorf(ErrInvalidOperation, "object key %q disappeared during comparison", rightKey)
			}
		}

		comparison, err := CompareValues(ctx, t.data[leftKey], rightValue)
		if err != nil {
			return Equal, err
		}

		if comparison != Equal {
			return comparison, nil
		}
	}

	return Equal, nil
}

// DebugFieldsPageable reports whether a lexicographic key sort is bounded.
func (t *Object) DebugFieldsPageable() bool { return t.bounds != nil && t.bounds.Pageable(len(t.data)) }

// ReadDebugFields visits the selected interval in stable key order.
func (t *Object) ReadDebugFields(start, end int, visit func(string, Value) error) error {
	if t.bounds == nil {
		return errors.ErrUnsupported
	}

	return debugmap.Read(t.data, *t.bounds, start, end, visit)
}

func (t *Object) setKey(key string, value Value) {
	t.bounds.AddKey(key)
	t.data[key] = value
}

func (t *Object) sortedKeys() []string {
	keys := make([]string, 0, len(t.data))
	for key := range t.data {
		keys = append(keys, key)
	}

	sort.Strings(keys)
	return keys
}

func (t *Object) Hash() uint64 {
	h := fnv.New64a()

	h.Write([]byte("object:"))
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

func (t *Object) Copy() Value {
	c := NewObject()

	for k, v := range t.data {
		c.setKey(k, v)
	}

	return c
}

func (t *Object) Clone(ctx context.Context) (Cloneable, error) {
	cloned := NewObject()

	var value Value

	for key := range t.data {
		value = t.data[key]
		cloneable, ok := value.(Cloneable)

		if ok {
			clone, err := cloneable.Clone(ctx)

			if err != nil {
				return nil, err
			}

			value = clone
		} else {
			value = value.Copy()
		}

		cloned.setKey(key, value)
	}

	return cloned, nil
}

func (t *Object) Length(_ context.Context) (Int, error) {
	return Int(len(t.data)), nil
}

func (t *Object) IsEmpty(_ context.Context) (Boolean, error) {
	return len(t.data) == 0, nil
}

func (t *Object) Keys(_ context.Context) (List, error) {
	keys := make([]Value, 0, len(t.data))

	for k := range t.data {
		keys = append(keys, NewString(k))
	}

	return NewArrayOf(keys), nil
}

func (t *Object) Values(_ context.Context) (List, error) {
	keys := make([]Value, 0, len(t.data))

	for _, v := range t.data {
		keys = append(keys, v)
	}

	return NewArrayOf(keys), nil
}

func (t *Object) ForEach(ctx context.Context, predicate KeyReadablePredicate) error {
	for key, val := range t.data {
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

func (t *Object) Filter(ctx context.Context, predicate KeyReadablePredicate) (List, error) {
	res := NewArray(len(t.data))

	for key, val := range t.data {
		match, err := predicate(ctx, val, String(key))

		if err != nil {
			return nil, err
		}

		if match {
			_ = res.Append(ctx, val)
		}
	}

	return res, nil
}

func (t *Object) Find(ctx context.Context, predicate KeyReadablePredicate) (Value, Boolean, error) {
	for key, val := range t.data {
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
	_, exists := t.data[key.String()]

	return Boolean(exists), nil
}

func (t *Object) Contains(ctx context.Context, target Value) (Boolean, error) {
	for _, val := range t.data {
		equal, err := EqualValues(ctx, target, val)
		if err != nil {
			return false, err
		}

		if equal {
			return true, nil
		}
	}

	return false, nil
}

func (t *Object) Get(_ context.Context, key Value) (Value, error) {
	val, found := t.data[key.String()]

	if found {
		return val, nil
	}

	return None, nil
}

func (t *Object) Lookup(_ context.Context, key Value) (Value, bool, error) {
	val, found := t.data[key.String()]

	if found {
		return val, true, nil
	}

	return None, false, nil
}

func (t *Object) Set(_ context.Context, key, value Value) error {
	t.SetString(key.String(), value)

	return nil
}

// SetString stores a value under a native string key.
func (t *Object) SetString(key string, value Value) {
	if value == nil {
		value = None
	}

	t.setKey(key, value)
}

func (t *Object) RemoveKey(_ context.Context, key Value) error {
	keyName := key.String()
	before := len(t.data)
	delete(t.data, keyName)
	if len(t.data) != before {
		t.bounds.Removed()
	}

	return nil
}

func (t *Object) Remove(ctx context.Context, value Value) error {
	for key, val := range t.data {
		equal, err := EqualValues(ctx, value, val)
		if err != nil {
			return err
		}

		if equal {
			delete(t.data, key)
			t.bounds.Removed()
			break
		}
	}

	return nil
}

func (t *Object) Merge(ctx context.Context, other Map) error {
	switch m := other.(type) {
	case *Object:
		for key, val := range m.data {
			t.setKey(key, val)
		}

		return nil
	default:
		return m.ForEach(ctx, func(ctx context.Context, value, key Value) (Boolean, error) {
			t.setKey(key.String(), value)
			return true, nil
		})
	}
}

func (t *Object) Clear(_ context.Context) error {
	*t = *NewObject()

	return nil
}

// New creates an independent empty collection with the receiver's configuration.
func (t *Object) New(ctx context.Context) (Map, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	return NewObject(), nil
}

func (t *Object) Iterate(_ context.Context) (Iterator, error) {
	return NewObjectIterator(t), nil
}
