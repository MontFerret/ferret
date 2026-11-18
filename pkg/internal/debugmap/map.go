// Package debugmap bounds the work of inspecting an unordered in-memory map.
package debugmap

// MaxKeys bounds both the key sort and the backing table traversed to build it.
const MaxKeys = 128

// Comparisons of long common prefixes can dominate even a small key sort.
const maxKeyBytes = 64 << 10

// Bounds records allocation history that len(map) alone cannot recover.
// It is reset only when the backing table is replaced.
type Bounds struct {
	removed   uint8
	oversized bool
}

func (b *Bounds) Reserve(size int) { b.oversized = size > MaxKeys }

func (b *Bounds) AddKey(key string) {
	if len(key) > maxKeyBytes {
		b.oversized = true
	}
}

func (b *Bounds) Removed() {
	if b.removed <= MaxKeys {
		b.removed++
	}
}

func (b Bounds) Pageable(live int) bool {
	return !b.oversized && live <= MaxKeys-int(b.removed)
}
