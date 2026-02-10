package mem

import (
	"regexp"

	"github.com/MontFerret/ferret/pkg/runtime"
)

type (
	Cache struct {
		FuncHash        uint64
		Functions       map[int]*CachedFunction
		Regexps         map[int]*regexp.Regexp
		LoadKeyICs      []*LoadKeyCache
		LoadKeyConstICs []*LoadKeyConstCache
	}

	CachedFunction struct {
		Fn0 runtime.Function0
		Fn1 runtime.Function1
		Fn2 runtime.Function2
		Fn3 runtime.Function3
		Fn4 runtime.Function4
		FnV runtime.Function
	}

	LoadKeyCache struct {
		entries     [loadKeyICEntries]LoadKeyCacheEntry
		size        uint8
		megamorphic bool
	}

	LoadKeyCacheEntry struct {
		ShapeID uint64
		Key     string
		Slot    int
	}

	LoadKeyConstCache struct {
		entries     [loadKeyICEntries]LoadKeyConstCacheEntry
		size        uint8
		megamorphic bool
	}

	LoadKeyConstCacheEntry struct {
		ShapeID uint64
		Slot    int
	}
)

const loadKeyICEntries = 4

func NewLoadKeyCache() *LoadKeyCache {
	return &LoadKeyCache{}
}

func (c *LoadKeyCache) Lookup(shapeID uint64, key string) (int, bool) {
	if c == nil || c.megamorphic {
		return 0, false
	}

	for i := 0; i < int(c.size); i++ {
		entry := c.entries[i]
		if entry.ShapeID == shapeID && entry.Key == key {
			return entry.Slot, true
		}
	}

	return 0, false
}

func (c *LoadKeyCache) Add(shapeID uint64, key string, slot int) {
	if c == nil || c.megamorphic {
		return
	}

	for i := 0; i < int(c.size); i++ {
		if c.entries[i].ShapeID == shapeID && c.entries[i].Key == key {
			c.entries[i].Slot = slot
			return
		}
	}

	if c.size < loadKeyICEntries {
		c.entries[c.size] = LoadKeyCacheEntry{
			ShapeID: shapeID,
			Key:     key,
			Slot:    slot,
		}
		c.size++
		return
	}

	c.megamorphic = true
}

func NewLoadKeyConstCache() *LoadKeyConstCache {
	return &LoadKeyConstCache{}
}

func (c *LoadKeyConstCache) Lookup(shapeID uint64) (int, bool) {
	if c == nil || c.megamorphic {
		return 0, false
	}

	for i := 0; i < int(c.size); i++ {
		entry := c.entries[i]
		if entry.ShapeID == shapeID {
			return entry.Slot, true
		}
	}

	return 0, false
}

func (c *LoadKeyConstCache) Add(shapeID uint64, slot int) {
	if c == nil || c.megamorphic {
		return
	}

	for i := 0; i < int(c.size); i++ {
		if c.entries[i].ShapeID == shapeID {
			c.entries[i].Slot = slot
			return
		}
	}

	if c.size < loadKeyICEntries {
		c.entries[c.size] = LoadKeyConstCacheEntry{
			ShapeID: shapeID,
			Slot:    slot,
		}
		c.size++
		return
	}

	c.megamorphic = true
}

func NewCache(bytecodeLen int) *Cache {
	return &Cache{
		Functions:       make(map[int]*CachedFunction),
		Regexps:         make(map[int]*regexp.Regexp),
		LoadKeyICs:      make([]*LoadKeyCache, bytecodeLen),
		LoadKeyConstICs: make([]*LoadKeyConstCache, bytecodeLen),
	}
}
