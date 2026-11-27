package runtime

import (
	"bytes"
	"encoding/binary"
	"hash/maphash"

	"github.com/MontFerret/ferret/v2/pkg/internal/debugmap"
)

const (
	smallObjectFields = 8
	// The name index stores one-based slot bounds in bytes.
	maxIndexedObjectFields uint8 = debugmap.MaxKeys
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

	// end[tag] bounds the prefix that can contain this tag. Each append
	// records its one-based slot. Swap-removal only moves entries left, so
	// stale bounds remain valid until replaced by a later append.
	objectNameIndex struct {
		end  [256]uint8
		tags [maxIndexedObjectFields]byte
	}

	objectProperty struct {
		value Value
		key   string
	}
)

var objectKeySeed = maphash.MakeSeed()

func newObjectBacking(size int) *objectBacking {
	if size <= smallObjectFields {
		return new(objectBacking)
	}

	// Keep aliases attached to the same owner when the slots grow into a
	// map. The separately allocated slots can then be collected.
	b := new(struct {
		objectBacking
		properties objectProperties
		names      objectNameIndex
	})
	b.properties.slots = make([]objectProperty, 0, size)
	b.properties.index = &b.names
	b.group = &b.properties

	return &b.objectBacking
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

func (t *Object) lookupValue(key Value) (Value, bool) {
	name, ok := key.(String)
	if !ok {
		return t.lookup(key.String())
	}
	if fields := t.large; fields != nil {
		value, found := fields[string(name)]
		return value, found
	}

	return t.lookup(string(name))
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
	n := len(key)
	if n > 8 {
		return byte(maphash.String(objectKeySeed, key))
	}

	if n >= 4 {
		// The first and last words cover the complete short name. Mix their
		// sum before taking the high byte to distinguish nearby names.
		first := binary.LittleEndian.Uint32([]byte(key[:4]))
		last := binary.LittleEndian.Uint32([]byte(key[n-4:]))

		return byte(((first + last) * 0x9e3779b1) >> 24)
	}

	var tag byte
	for i := 0; i < n; i++ {
		tag = tag*31 + key[i]
	}
	return tag
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
	end := min(int(b.index.end[tag]), len(b.slots))
	if end == 0 {
		return -1, tag
	}

	if b.slots[end-1].key == key {
		return end - 1, tag
	}

	tags := b.index.tags[:end-1]
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
		g.slots[idx] = objectProperty{key: key, value: value}

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
		g.index.end[tag] = uint8(len(g.slots))
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

	idx, _ := g.find(key)
	if idx < 0 {
		return
	}

	last := len(g.slots) - 1
	g.slots[idx] = g.slots[last]
	g.slots[last] = objectProperty{}
	g.slots = g.slots[:last]

	if ix := g.index; ix != nil {
		ix.tags[idx] = ix.tags[last]
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

	slots := t.backing.slots()
	if len(slots) <= 1 {
		if len(slots) == 1 {
			yield(slots[0].key, slots[0].value)
		}

		return
	}

	t.callbackSlotEntries(slots, yield)
}

func (t Object) callbackSlotEntries(slots []objectProperty, yield func(string, Value) bool) {
	// A callback may remove a field, move its slot or promote the group.
	// Preserve the original names and read their current values on delivery.
	var names [debugmap.MaxKeys]string
	for i := range slots {
		names[i] = slots[i].key
	}

	for i, key := range names[:len(slots)] {
		current := t.backing.slots()

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
