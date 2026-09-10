package valueset_test

import (
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/internal/valueset"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

func BenchmarkSetInsert(b *testing.B) {
	for _, test := range []struct {
		name   string
		values []runtime.Value
	}{
		{name: "distinct", values: []runtime.Value{runtime.Int(1), runtime.Int(2), runtime.Int(3), runtime.Int(4)}},
		{name: "repeated", values: []runtime.Value{runtime.Int(1), runtime.Int(2), runtime.Int(1), runtime.Int(2)}},
		{name: "mixed_numeric", values: []runtime.Value{runtime.Int(1), runtime.Int(2), runtime.Float(1), runtime.Float(2)}},
		{name: "collisions", values: []runtime.Value{collisionValue{label: "a"}, collisionValue{label: "b"}, collisionValue{label: "c"}, collisionValue{label: "d"}}},
	} {
		b.Run(test.name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				set := valueset.New(0)
				for _, value := range test.values {
					if _, err := set.Add(b.Context(), value); err != nil {
						b.Fatal(err)
					}
				}
			}
		})
	}
}

func BenchmarkSetLookup(b *testing.B) {
	for _, test := range []struct {
		value  runtime.Value
		name   string
		values []runtime.Value
		found  bool
	}{
		{name: "integer", values: []runtime.Value{runtime.Int(1)}, value: runtime.Int(1), found: true},
		{name: "mixed_numeric", values: []runtime.Value{runtime.Int(1)}, value: runtime.Float(1), found: true},
		{name: "missing_hash", values: []runtime.Value{runtime.Int(1)}, value: runtime.Int(2)},
		{name: "collision_first", values: []runtime.Value{collisionValue{label: "a"}, collisionValue{label: "b"}}, value: collisionValue{label: "a"}, found: true},
		{name: "collision_later", values: []runtime.Value{collisionValue{label: "a"}, collisionValue{label: "b"}}, value: collisionValue{label: "b"}, found: true},
		{name: "collision_missing", values: []runtime.Value{collisionValue{label: "a"}, collisionValue{label: "b"}}, value: collisionValue{label: "c"}},
	} {
		b.Run(test.name, func(b *testing.B) {
			set := valueset.New(0)
			for _, value := range test.values {
				if _, err := set.Add(b.Context(), value); err != nil {
					b.Fatal(err)
				}
			}

			b.Run("contains", func(b *testing.B) {
				b.ReportAllocs()
				for b.Loop() {
					found, err := set.Contains(b.Context(), test.value)
					if err != nil || found != test.found {
						b.Fatalf("contains: %t, %v", found, err)
					}
				}
			})

			if !test.found {
				return
			}

			b.Run("add", func(b *testing.B) {
				b.ReportAllocs()
				for b.Loop() {
					added, err := set.Add(b.Context(), test.value)
					if err != nil || added {
						b.Fatalf("duplicate add: %t, %v", added, err)
					}
				}
			})
		})
	}
}
