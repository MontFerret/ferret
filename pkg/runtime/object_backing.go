package runtime

import (
	"bytes"
	"hash/maphash"

	"github.com/MontFerret/ferret/v2/pkg/internal/debugmap"
)

const smallObjectFields = 8

// Backings are shared by Go value copies and never refer to the receiver.
// A bounded vector can filter candidate names by a one-byte fingerprint.
// A fingerprint match always requires an exact Go string comparison.
type objectBacking struct {
	fields map[string]Value
	index  *objectNameIndex
	slots  []objectProperty
	bounds debugmap.Bounds
}

// Each tag remembers one candidate slot. On a collision the byte vector
// finds other candidates without comparing every property name. A removal may
// leave a stale hint. Only zero proves the tag has never been inserted.
type objectNameIndex struct {
	candidate [256]uint8
	tags      [debugmap.MaxKeys]byte
}

type objectProperty struct {
	value Value
	key   string
}

var objectKeySeed = maphash.MakeSeed()

func newObjectBacking(size int) *objectBacking {
	b := new(objectBacking)
	if size > smallObjectFields {
		b.slots = make([]objectProperty, 0, size)
		b.index = new(objectNameIndex)
	}

	return b
}

func (t *Object) length() int {
	if b := t.backing; b != nil {
		return len(b.slots) + len(b.fields)
	}

	return len(t.large)
}

func (t *Object) lookup(key string) (Value, bool) {
	fields := t.large
	if b := t.backing; b != nil {
		if b.fields == nil {
			i, _ := b.find(key)
			if i >= 0 {
				return b.slots[i].value, true
			}

			return nil, false
		}

		fields = b.fields
	}

	v, ok := fields[key]

	return v, ok
}

func (b *objectBacking) find(key string) (int, byte) {
	if b.index == nil {
		for i := range b.slots {
			if b.slots[i].key == key {
				return i, 0
			}
		}

		return -1, 0
	}

	tag := byte(maphash.String(objectKeySeed, key))
	hint := b.index.candidate[tag]
	if hint == 0 {
		return -1, tag
	}

	slot := int(hint) - 1
	if slot < len(b.slots) && b.slots[slot].key == key {
		return slot, tag
	}

	tags := b.index.tags[:len(b.slots)]
	for offset := 0; offset < len(tags); {
		pos := bytes.IndexByte(tags[offset:], tag)
		if pos < 0 {
			break
		}

		i := offset + pos
		if b.slots[i].key == key {
			return i, tag
		}

		offset = i + 1
	}

	return -1, tag
}

func (t *Object) setKey(key string, value Value) {
	b := t.backing
	if b == nil {
		t.large[key] = value

		return
	}

	b.bounds.AddKey(key)
	if b.fields != nil {
		b.fields[key] = value

		return
	}

	idx, tag := b.find(key)
	if idx >= 0 {
		b.slots[idx].value = value

		return
	}

	if b.slots == nil {
		b.slots = make([]objectProperty, 0, smallObjectFields)
	}

	if len(b.slots) == cap(b.slots) {
		// Skip the first small Go map growth when spilling an unhinted
		// vector. Keep the reservation within the proven inspection limit.
		size := min(debugmap.MaxKeys, 2*len(b.slots))
		fields := make(map[string]Value, size)
		for _, property := range b.slots {
			fields[property.key] = property.value
		}

		fields[key] = value
		b.slots = nil
		b.index = nil
		b.fields = fields

		return
	}

	b.slots = append(b.slots, objectProperty{key: key, value: value})
	if b.index != nil {
		b.index.tags[len(b.slots)-1] = tag
		b.index.candidate[tag] = uint8(len(b.slots))
	}
}

func (t *Object) deleteKey(key string) {
	b := t.backing
	if b == nil {
		delete(t.large, key)

		return
	}

	if b.fields != nil {
		before := len(b.fields)
		delete(b.fields, key)
		if len(b.fields) < before {
			b.bounds.Removed()
		}

		return
	}

	idx, _ := b.find(key)
	if idx < 0 {
		return
	}

	last := len(b.slots) - 1
	b.slots[idx] = b.slots[last]
	b.slots[last] = objectProperty{}
	b.slots = b.slots[:last]
	if b.index != nil && idx != last {
		tag := b.index.tags[last]
		b.index.tags[idx] = tag
		b.index.candidate[tag] = uint8(idx + 1)
	}

	b.bounds.Removed()
}

// entries is used when the consumer cannot mutate this backing.
func (t Object) entries(yield func(string, Value) bool) {
	fields := t.large
	if b := t.backing; b != nil {
		if b.fields == nil {
			for _, p := range b.slots {
				if !yield(p.key, p.value) {
					return
				}
			}

			return
		}

		fields = b.fields
	}

	for key, value := range fields {
		if !yield(key, value) {
			return
		}
	}
}

func (t Object) callbackEntries(yield func(string, Value) bool) {
	b := t.backing
	if b == nil || b.fields != nil || len(b.slots) <= 1 {
		t.entries(yield)

		return
	}

	// Callbacks may delete, promote, or clear. Snapshot the names and read
	// each surviving value from the original backing.
	var names [debugmap.MaxKeys]string
	for i := range b.slots {
		names[i] = b.slots[i].key
	}
	for _, key := range names[:len(b.slots)] {
		if v, ok := t.lookup(key); ok && !yield(key, v) {
			return
		}
	}
}

func (t *Object) propertyNames() []string {
	fields := t.large
	if b := t.backing; b != nil {
		if b.fields == nil {
			keys := make([]string, len(b.slots))
			for i := range b.slots {
				keys[i] = b.slots[i].key
			}

			return keys
		}

		fields = b.fields
	}

	keys := make([]string, 0, len(fields))
	for key := range fields {
		keys = append(keys, key)
	}

	return keys
}
