package runtime

import (
	"bytes"
	"hash/maphash"

	"github.com/MontFerret/ferret/v2/pkg/internal/debugmap"
)

const (
	smallObjectFields   = 8
	multipleObjectSlots = 255 // Above every one-based slot position.
)

type (
	// A backing holds the first property inline. Copies share the backing,
	// which never points back to a receiver that Clear could replace.
	objectBacking struct {
		group    *objectProperties
		first    [1]objectProperty
		bounds   debugmap.Bounds
		firstLen uint8
	}

	objectProperties struct {
		fields    map[string]Value
		index     *objectNameIndex
		slots     []objectProperty
		smallBits uint64
	}

	// A tag is absent, names one slot, or requires a bounded scan for
	// collisions. Every candidate is compared by its exact string key.
	objectNameIndex struct {
		candidate [256]uint8
		tags      [debugmap.MaxKeys]byte
	}

	objectProperty struct {
		value Value
		key   string
	}

	// Co-allocate the detached owner, its index and one bounded size class.
	indexedObject struct {
		objectBacking
		properties objectProperties
		names      objectNameIndex
	}

	indexedObjectBlock[Slots any] struct {
		slots Slots
		indexedObject
	}
)

var objectKeySeed = maphash.MakeSeed()

func (b *indexedObject) withSlots(slots []objectProperty) *objectBacking {
	b.properties.slots = slots
	b.properties.index = &b.names
	b.group = &b.properties
	return &b.objectBacking
}

func newObjectBacking(size int) *objectBacking {
	switch {
	case size <= smallObjectFields:
		return new(objectBacking)
	case size <= 16:
		p := new(indexedObjectBlock[[16]objectProperty])
		return p.withSlots(p.slots[:0])
	case size <= 32:
		p := new(indexedObjectBlock[[32]objectProperty])
		return p.withSlots(p.slots[:0])
	case size <= 64:
		p := new(indexedObjectBlock[[64]objectProperty])
		return p.withSlots(p.slots[:0])
	default:
		p := new(indexedObjectBlock[[debugmap.MaxKeys]objectProperty])
		return p.withSlots(p.slots[:0])
	}
}

// A clear bit rules out the key. Shared bits and old bits left by deletions
// only cause an exact scan of the small group.
func smallObjectBit(key string) uint64 {
	tag := len(key)
	if tag != 0 {
		tag ^= int(key[len(key)-1])
	}
	return uint64(1) << (tag & 63)
}

func (b *objectBacking) slots() []objectProperty {
	if b.group == nil {
		return b.first[:b.firstLen]
	}
	return b.group.slots
}

func (t Object) fields() map[string]Value {
	if b := t.backing; b != nil && b.group != nil {
		return b.group.fields
	}

	return t.large
}

func (t *Object) length() int {
	if b := t.backing; b != nil {
		if b.group == nil {
			return int(b.firstLen)
		}
		return len(b.group.slots) + len(b.group.fields)
	}

	return len(t.large)
}

func (t *Object) lookup(key string) (Value, bool) {
	fields := t.large
	if b := t.backing; b != nil {
		g := b.group
		if g == nil {
			p := &b.first[0]
			if b.firstLen != 0 && p.key == key {
				return p.value, true
			}
			return nil, false
		}
		if g.fields == nil {
			i, _ := g.find(key)
			if i >= 0 {
				return g.slots[i].value, true
			}
			return nil, false
		}
		fields = g.fields
	}

	v, ok := fields[key]

	return v, ok
}

func objectNameTag(key string) byte {
	return byte(maphash.Comparable(objectKeySeed, key))
}

func (b *objectProperties) find(key string) (int, byte) {
	if b.index == nil {
		if b.smallBits&smallObjectBit(key) == 0 {
			return -1, 0
		}
		for i := range b.slots {
			if b.slots[i].key == key {
				return i, 0
			}
		}

		return -1, 0
	}

	tag := objectNameTag(key)
	hint := b.index.candidate[tag]
	if hint == 0 {
		return -1, tag
	}

	if hint != multipleObjectSlots {
		slot := int(hint) - 1
		if b.slots[slot].key == key {
			return slot, tag
		}
		return -1, tag
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
	if t.backing == nil {
		t.large[key] = value
	} else {
		t.setTrackedKey(key, value)
	}
}

func (t *Object) setTrackedKey(key string, value Value) {
	b := t.backing

	b.bounds.AddKey(key)
	g := b.group
	if g == nil {
		if b.firstLen == 0 || b.first[0].key == key {
			b.first[0] = objectProperty{key: key, value: value}
			b.firstLen = 1
			return
		}
		p := new(struct {
			inline [smallObjectFields]objectProperty
			objectProperties
		})
		p.inline[0] = b.first[0]
		p.inline[1] = objectProperty{key: key, value: value}
		p.slots = p.inline[:2]
		p.smallBits = smallObjectBit(p.inline[0].key) | smallObjectBit(key)
		b.first[0] = objectProperty{}
		b.firstLen = 0
		b.group = &p.objectProperties
		return
	}

	if g.fields != nil {
		g.fields[key] = value

		return
	}

	idx, tag := g.find(key)
	if idx >= 0 {
		g.slots[idx].value = value

		return
	}

	if len(g.slots) == cap(g.slots) {
		// Skip the first small Go map growth when spilling an unhinted
		// vector. Keep the reservation within the proven inspection limit.
		size := min(debugmap.MaxKeys, 2*len(g.slots))
		fields := make(map[string]Value, size)
		for _, property := range g.slots {
			fields[property.key] = property.value
		}

		fields[key] = value
		clear(g.slots)
		g.slots = nil
		g.index = nil
		g.fields = fields

		return
	}

	g.slots = append(g.slots, objectProperty{key: key, value: value})
	if g.index == nil {
		g.smallBits |= smallObjectBit(key)
	} else {
		g.index.tags[len(g.slots)-1] = tag
		if g.index.candidate[tag] == 0 {
			g.index.candidate[tag] = uint8(len(g.slots))
		} else {
			g.index.candidate[tag] = multipleObjectSlots
		}
	}
}

func (t *Object) deleteKey(key string) {
	b := t.backing
	if b == nil {
		delete(t.large, key)

		return
	}

	g := b.group
	if g == nil {
		if b.firstLen != 0 && b.first[0].key == key {
			b.first[0] = objectProperty{}
			b.firstLen = 0
			b.bounds.Removed()
		}
		return
	}

	if g.fields != nil {
		before := len(g.fields)
		delete(g.fields, key)
		if len(g.fields) < before {
			b.bounds.Removed()
		}

		return
	}

	idx, removedTag := g.find(key)
	if idx < 0 {
		return
	}

	last := len(g.slots) - 1
	g.slots[idx] = g.slots[last]
	g.slots[last] = objectProperty{}
	g.slots = g.slots[:last]
	if ix := g.index; ix != nil {
		if ix.candidate[removedTag] != multipleObjectSlots {
			ix.candidate[removedTag] = 0
		}
		if idx != last {
			tag := ix.tags[last]
			ix.tags[idx] = tag
			if ix.candidate[tag] != multipleObjectSlots {
				ix.candidate[tag] = uint8(idx + 1)
			}
		}
	}

	b.bounds.Removed()
}

// entries is used when the consumer cannot mutate this backing.
func (t Object) entries(yield func(string, Value) bool) {
	fields := t.fields()
	if fields == nil && t.backing != nil {
		for _, p := range t.backing.slots() {
			if !yield(p.key, p.value) {
				return
			}
		}
		return
	}

	for key, value := range fields {
		if !yield(key, value) {
			return
		}
	}
}

func (t Object) callbackEntries(yield func(string, Value) bool) {
	if fields := t.fields(); fields != nil || t.backing == nil {
		for key, value := range fields {
			if !yield(key, value) {
				return
			}
		}
		return
	}

	b := t.backing
	slots := b.slots()
	if len(slots) <= 1 {
		if len(slots) == 1 {
			yield(slots[0].key, slots[0].value)
		}
		return
	}

	// A callback may remove a field, move its slot or promote the group.
	// Preserve the original names and read their current values on delivery.
	var names [debugmap.MaxKeys]string
	for i := range slots {
		names[i] = slots[i].key
	}
	for i, key := range names[:len(slots)] {
		current := b.slots()
		if i < len(current) && current[i].key == key {
			if !yield(key, current[i].value) {
				return
			}
		} else if v, ok := t.lookup(key); ok && !yield(key, v) {
			return
		}
	}
}

func (t *Object) propertyNames() []string {
	fields := t.fields()
	if fields == nil && t.backing != nil {
		slots := t.backing.slots()
		keys := make([]string, len(slots))
		for i := range slots {
			keys[i] = slots[i].key
		}

		return keys
	}

	keys := make([]string, 0, len(fields))
	for key := range fields {
		keys = append(keys, key)
	}

	return keys
}
