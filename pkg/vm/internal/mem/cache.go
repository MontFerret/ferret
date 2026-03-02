package mem

import (
	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/vm/internal/data"
)

type (
	Cache struct {
		FuncHash        uint64
		Functions       []*CachedFunction
		RegexpsWarmed   bool
		Regexps         []*CachedRegexp
		LoadKeyICs      []*LoadKeyCache
		LoadKeyConstICs []*LoadKeyConstCache
		ShapeCache      *data.ShapeCache
	}

	CachedFunction struct {
		Fn0 runtime.Function0
		Fn1 runtime.Function1
		Fn2 runtime.Function2
		Fn3 runtime.Function3
		Fn4 runtime.Function4
		FnV runtime.Function
	}

	CachedRegexp struct {
		Pattern string
		Regexp  *data.Regexp
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
		shapeID     uint64
		slot        int
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
	return &LoadKeyConstCache{slot: -1}
}

func (c *LoadKeyConstCache) Lookup(shapeID uint64) (int, bool) {
	if c == nil || c.megamorphic {
		return 0, false
	}

	if c.shapeID == shapeID {
		return c.slot, true
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

	if c.shapeID == 0 {
		c.shapeID = shapeID
		c.slot = slot
		return
	}

	if c.shapeID == shapeID {
		c.slot = slot
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

func NewCache(bytecodeLen, shapeCacheLimit int) *Cache {
	return &Cache{
		Functions:       make([]*CachedFunction, bytecodeLen),
		Regexps:         make([]*CachedRegexp, bytecodeLen),
		LoadKeyICs:      make([]*LoadKeyCache, bytecodeLen),
		LoadKeyConstICs: make([]*LoadKeyConstCache, bytecodeLen),
		ShapeCache:      data.NewShapeCache(shapeCacheLimit),
	}
}
