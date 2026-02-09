package runtime

import (
	"context"
)

type ObjectIterator struct {
	entries []objectIterEntry
	slots   []Value
	pos     int
}

type objectIterEntry struct {
	key  string
	slot int
}

func NewObjectIterator(obj *Object) Iterator {
	entries := make([]objectIterEntry, 0, obj.size)

	for idx, key := range obj.shape.names {
		if obj.slots[idx] == nil {
			continue
		}
		entries = append(entries, objectIterEntry{key: key, slot: idx})
	}

	return &ObjectIterator{entries: entries, slots: obj.slots}
}

func (iter *ObjectIterator) HasNext(_ context.Context) (bool, error) {
	return len(iter.entries) > iter.pos, nil
}

func (iter *ObjectIterator) Next(_ context.Context) (Value, Value, error) {
	if iter.pos >= len(iter.entries) {
		return None, None, Error(ErrInvalidOperation, "no more elements")
	}

	entry := iter.entries[iter.pos]
	value := iter.slots[entry.slot]
	if value == nil {
		value = None
	}
	iter.pos++

	return value, String(entry.key), nil
}
