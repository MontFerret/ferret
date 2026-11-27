package runtime

import (
	"strconv"
	"testing"
)

// Measure how much work the first named page would need without backing
// history. Every fixture has the same live property when inspection starts.
func BenchmarkObjectInspectionKeyDiscovery(b *testing.B) {
	for _, tc := range []struct {
		name    string
		reserve int
		grow    int
	}{
		{name: "fresh"},
		{name: "reserved_65536", reserve: 65_536},
		{name: "reserved_1000000", reserve: 1_000_000},
		{name: "deleted_1000000", grow: 1_000_000},
	} {
		b.Run(tc.name, func(b *testing.B) {
			const survivor = "keep"
			values := make(map[string]Value, tc.reserve)
			values[survivor] = True

			for i := 1; i < tc.grow; i++ {
				values[strconv.Itoa(i)] = False
			}
			for i := 1; i < tc.grow; i++ {
				delete(values, strconv.Itoa(i))
			}
			if len(values) != 1 {
				b.Fatalf("got %d live properties", len(values))
			}

			b.Run("first_key", func(b *testing.B) {
				for b.Loop() {
					var key string
					for key = range values {
						break
					}
					if key != survivor {
						b.Fatalf("discovered %q", key)
					}
				}
			})
			b.Run("known_key", func(b *testing.B) {
				for b.Loop() {
					if values[survivor] != True {
						b.Fatal("lost known property")
					}
				}
			})
		})
	}
}
