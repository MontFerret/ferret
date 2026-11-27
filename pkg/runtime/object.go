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

// Object is a handle. Value copies share its properties until Clear detaches
// one receiver. A pre-sized map has a shared history only when its original
// reservation permits bounded debugger enumeration.
type Object struct {
	backing  *objectBacking
	reserved map[string]Value
}

func NewObject() *Object {
	return NewObjectOf(0)
}

func NewObjectOf(size int) *Object {
	if size <= smallObjectFields {
		return &Object{backing: new(objectBacking)}
	}
	obj := &Object{reserved: make(map[string]Value, size)}
	if size <= debugmap.MaxKeys {
		obj.backing = new(objectBacking)
	}
	return obj
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

	if Int(t.length()) != otherSize {
		return false, nil
	}

	for key, value := range t.callbackEntries {
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

	size := Int(t.length())
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
			rightValue, _ = otherObject.lookup(rightKey)
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

		leftValue, _ := t.lookup(leftKey)
		comparison, err := CompareValues(ctx, leftValue, rightValue)
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
func (t *Object) DebugFieldsPageable() bool {
	return t.backing != nil && t.backing.bounds.Pageable(t.length())
}

// ReadDebugFields visits the selected interval in stable key order.
func (t *Object) ReadDebugFields(start, end int, visit func(string, Value) error) error {
	if fields := t.fields(); t.backing != nil && fields != nil {
		return debugmap.Read(fields, t.backing.bounds, start, end, visit)
	}

	size := t.length()
	if !t.DebugFieldsPageable() || start < 0 || end < start || end > size {
		return errors.ErrUnsupported
	}

	if start == end {
		return nil
	}

	view := *t
	keys := view.sortedKeys()
	for _, key := range keys[start:end] {
		value, _ := view.lookup(key)
		if err := visit(key, value); err != nil {
			return err
		}
	}

	return nil
}

func (t *Object) Hash() uint64 {
	h := fnv.New64a()

	h.Write([]byte("object:"))
	h.Write([]byte("{"))

	keys := t.sortedKeys()
	endIndex := len(keys) - 1

	for idx, key := range keys {
		h.Write([]byte(key))
		h.Write([]byte(":"))

		el, _ := t.lookup(key)
		if el == nil {
			el = None
		}

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
	var copied Object
	t.copyInto(&copied)

	return &copied
}

func (t *Object) Clone(ctx context.Context) (Cloneable, error) {
	return t.cloneInto(ctx, new(Object))
}

func (t *Object) Length(_ context.Context) (Int, error) {
	return Int(t.length()), nil
}

func (t *Object) IsEmpty(_ context.Context) (Boolean, error) {
	return t.length() == 0, nil
}

func (t *Object) Keys(_ context.Context) (List, error) {
	keys := make([]Value, 0, t.length())

	for k := range t.entries {
		keys = append(keys, NewString(k))
	}

	return NewArrayOf(keys), nil
}

func (t *Object) Values(_ context.Context) (List, error) {
	keys := make([]Value, 0, t.length())

	for _, v := range t.entries {
		keys = append(keys, v)
	}

	return NewArrayOf(keys), nil
}

func (t *Object) ForEach(ctx context.Context, predicate KeyReadablePredicate) error {
	// Capture the backing across Clear. Keeping the map loop with the supplied
	// predicate lets Go inline a concrete callback without boxing each key.
	return func(view Object) (err error) {
		visit := func(key string, value Value) bool {
			more, visitErr := predicate(ctx, value, String(key))
			if visitErr != nil {
				err = visitErr
				return false
			}
			return bool(more)
		}
		if fields := view.fields(); fields != nil {
			for key, value := range fields {
				if !visit(key, value) {
					break
				}
			}
		} else {
			view.callbackEntries(visit)
		}
		return err
	}(*t)
}

func (t *Object) Filter(ctx context.Context, predicate KeyReadablePredicate) (List, error) {
	res := NewArray(t.length())

	for key, val := range t.callbackEntries {
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
	for key, val := range t.callbackEntries {
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
	_, exists := t.lookupValue(key)

	return Boolean(exists), nil
}

func (t *Object) Contains(ctx context.Context, target Value) (Boolean, error) {
	for _, val := range t.callbackEntries {
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
	val, found := t.lookupValue(key)

	if found {
		return val, nil
	}

	return None, nil
}

func (t *Object) Lookup(_ context.Context, key Value) (Value, bool, error) {
	val, found := t.lookupValue(key)

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
	t.deleteKey(key.String())

	return nil
}

func (t *Object) Remove(ctx context.Context, value Value) error {
	for key, val := range t.callbackEntries {
		equal, err := EqualValues(ctx, value, val)
		if err != nil {
			return err
		}

		if equal {
			t.deleteKey(key)
			break
		}
	}

	return nil
}

func (t *Object) Merge(ctx context.Context, other Map) error {
	switch m := other.(type) {
	case *Object:
		for key, val := range m.entries {
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

// An embedded dictionary can take the resulting handle by value. Keeping the
// allocation at the caller lets it retain only the detached property storage.
func (t *Object) copyInto(dst *Object) {
	*dst = *NewObjectOf(t.length())
	for key, value := range t.entries {
		dst.setKey(key, value)
	}
}

func (t *Object) cloneInto(ctx context.Context, cloned *Object) (Cloneable, error) {
	*cloned = *NewObjectOf(t.length())
	for key, value := range t.callbackEntries {
		if value != nil {
			clone, err := CloneOrCopy(ctx, value)
			if err != nil {
				return nil, err
			}
			value = clone
		}

		cloned.setKey(key, value)
	}

	return cloned, nil
}

func (t *Object) sortedKeys() []string {
	keys := t.propertyNames()
	sort.Strings(keys)
	return keys
}
