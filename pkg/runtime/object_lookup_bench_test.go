package runtime

import (
	"fmt"
	"strconv"
	"strings"
	"testing"
)

func BenchmarkObjectLookupKeys(b *testing.B) {
	for _, size := range []int{16, 128} {
		for _, pattern := range []string{"short", "four_bytes", "eight_bytes", "common_suffix"} {
			obj := NewObjectOf(size)
			hits := make([]Value, size)
			misses := make([]Value, size)
			for i := range size {
				name := strconv.Itoa(i)
				missing := strconv.Itoa(i + size)
				switch pattern {
				case "four_bytes":
					name = fmt.Sprintf("k%03d", i)
					missing = fmt.Sprintf("k%03d", i+size)
				case "eight_bytes":
					name = fmt.Sprintf("field%03d", i)
					missing = fmt.Sprintf("field%03d", i+size)
				case "common_suffix":
					name = fmt.Sprintf("property_%03d_common_suffix", i)
					missing = fmt.Sprintf("property_%03d_common_suffix", i+size)
				}
				_ = obj.Set(b.Context(), String(name), True)
				hits[i] = String(strings.Clone(name))
				misses[i] = String(missing)
			}
			for _, tc := range []struct {
				name string
				keys []Value
			}{{"hit", hits}, {"miss", misses}} {
				b.Run(strconv.Itoa(size)+"/"+pattern+"/"+tc.name, func(b *testing.B) {
					ctx := b.Context()
					i := 0
					for b.Loop() {
						benchmarkObjectValue, _, _ = obj.Lookup(ctx, tc.keys[i])
						i = (i + 1) & (size - 1)
					}
				})
			}
		}
	}
}
