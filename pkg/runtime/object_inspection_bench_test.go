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
			obj := NewObjectOf(tc.reserve)
			obj.SetString(survivor, True)

			for i := 1; i < tc.grow; i++ {
				obj.SetString(strconv.Itoa(i), False)
			}
			for i := 1; i < tc.grow; i++ {
				if err := obj.RemoveKey(b.Context(), String(strconv.Itoa(i))); err != nil {
					b.Fatal(err)
				}
			}
			if len(obj.data) != 1 {
				b.Fatalf("got %d live properties", len(obj.data))
			}

			b.Run("first_key", func(b *testing.B) {
				for b.Loop() {
					var key string
					for key = range obj.data {
						break
					}
					if key != survivor {
						b.Fatalf("discovered %q", key)
					}
				}
			})
			b.Run("known_key", func(b *testing.B) {
				for b.Loop() {
					if obj.data[survivor] != True {
						b.Fatal("lost known property")
					}
				}
			})
		})
	}
}
