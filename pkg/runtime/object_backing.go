package runtime

import "github.com/MontFerret/ferret/v2/pkg/internal/debugmap"

const smallObjectFields = 8

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
		slots     [smallObjectFields]objectProperty
		smallBits uint64
		count     uint8
	}

	objectProperty struct {
		value Value
		key   string
	}
)

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
	return b.group.slots[:b.group.count]
}

func (t Object) fields() map[string]Value {
	if t.reserved != nil || t.backing == nil || t.backing.group == nil {
		return t.reserved
	}

	return t.backing.group.fields
}

func (t *Object) length() int {
	if t.reserved != nil || t.backing == nil {
		return len(t.reserved)
	}

	if g := t.backing.group; g != nil {
		return int(g.count) + len(g.fields)
	}

	return int(t.backing.firstLen)
}

func (t *Object) lookupValue(key Value) (Value, bool) {
	name, ok := key.(String)
	if !ok {
		return t.lookup(key.String())
	}
	if fields := t.reserved; fields != nil {
		value, found := fields[string(name)]
		return value, found
	}

	return t.lookup(string(name))
}

func (t *Object) lookup(key string) (Value, bool) {
	if fields := t.fields(); fields != nil || t.backing == nil {
		value, found := fields[key]
		return value, found
	}

	b := t.backing
	if g := b.group; g != nil {
		if i := g.find(key); i >= 0 {
			return g.slots[i].value, true
		}
	} else if b.firstLen != 0 && b.first[0].key == key {
		return b.first[0].value, true
	}

	return nil, false
}

func (g *objectProperties) find(key string) int {
	if g.smallBits&smallObjectBit(key) != 0 {
		for i := range int(g.count) {
			if g.slots[i].key == key {
				return i
			}
		}
	}

	return -1
}

func (t *Object) setKey(key string, value Value) {
	b := t.backing
	if b != nil {
		b.bounds.AddKey(key)
	}

	if fields := t.reserved; fields != nil || b == nil {
		fields[key] = value
		return
	}

	b.setKey(key, value)
}

func (b *objectBacking) setKey(key string, value Value) {
	g := b.group
	if g == nil {
		if b.firstLen == 0 || b.first[0].key == key {
			b.first[0] = objectProperty{key: key, value: value}
			b.firstLen = 1

			return
		}

		p := &objectProperties{count: 2}
		p.slots[0] = b.first[0]
		p.slots[1] = objectProperty{key: key, value: value}
		p.smallBits = smallObjectBit(p.slots[0].key) | smallObjectBit(key)
		b.first[0] = objectProperty{}
		b.firstLen = 0
		b.group = p

		return
	}

	if g.fields != nil {
		g.fields[key] = value

		return
	}

	idx := g.find(key)
	if idx >= 0 {
		g.slots[idx] = objectProperty{key: key, value: value}

		return
	}

	if g.count == smallObjectFields {
		fields := make(map[string]Value, 2*smallObjectFields)
		for _, property := range g.slots {
			fields[property.key] = property.value
		}

		fields[key] = value
		clear(g.slots[:])
		g.count = 0
		g.fields = fields

		return
	}

	g.slots[g.count] = objectProperty{key: key, value: value}
	g.count++
	g.smallBits |= smallObjectBit(key)
}

func (t *Object) deleteKey(key string) {
	b := t.backing
	if fields := t.fields(); fields != nil || b == nil {
		before := len(fields)
		delete(fields, key)
		if b != nil && len(fields) < before {
			b.bounds.Removed()
		}

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

	idx := g.find(key)
	if idx < 0 {
		return
	}

	last := g.count - 1
	g.slots[idx] = g.slots[last]
	g.slots[last] = objectProperty{}
	g.count = last

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
	var names [smallObjectFields]string
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
