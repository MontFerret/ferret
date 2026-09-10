package objects_test

import (
	"context"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/stdlib/objects"
)

func BenchmarkObjectTransform(b *testing.B) {
	ctx := context.Background()
	left := runtime.NewObjectWith(map[string]runtime.Value{
		"name":   runtime.String("Ferret"),
		"nested": runtime.NewObjectWith(map[string]runtime.Value{"left": runtime.Int(1)}),
		"items":  runtime.NewArrayWith(runtime.Int(1), runtime.Int(2)),
	})
	right := runtime.NewObjectWith(map[string]runtime.Value{
		"nested":  runtime.NewObjectWith(map[string]runtime.Value{"right": runtime.Int(2)}),
		"version": runtime.Int(2),
	})
	keys := runtime.NewArrayWith(runtime.String("left"), runtime.String("right"))
	values := runtime.NewArrayWith(left, right)
	entries := runtime.NewArrayWith(runtime.NewArrayWith(runtime.String("left"), left), runtime.NewArrayWith(runtime.String("right"), right))
	cases := map[string]func() (runtime.Value, error){
		"merge":        func() (runtime.Value, error) { return objects.Merge(ctx, left, right) },
		"merge_deep":   func() (runtime.Value, error) { return objects.MergeDeep(ctx, left, right) },
		"keep_keys":    func() (runtime.Value, error) { return objects.KeepKeys(ctx, left, runtime.String("nested")) },
		"values":       func() (runtime.Value, error) { return objects.Values(ctx, left) },
		"entries":      func() (runtime.Value, error) { return objects.Entries(ctx, left) },
		"from_entries": func() (runtime.Value, error) { return objects.FromEntries(ctx, entries) },
		"omit_keys":    func() (runtime.Value, error) { return objects.OmitKeys(ctx, left, runtime.String("name")) },
		"zip":          func() (runtime.Value, error) { return objects.Zip(ctx, keys, values) },
	}
	for name, run := range cases {
		b.Run(name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				if _, err := run(); err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}

func BenchmarkMutableObjectTransform(b *testing.B) {
	ctx := context.Background()
	patch := runtime.NewObjectWith(map[string]runtime.Value{
		"nested":  runtime.NewObjectWith(map[string]runtime.Value{"right": runtime.Int(2)}),
		"version": runtime.Int(2),
	})
	for _, test := range []struct {
		name  string
		apply runtime.Function
		args  []runtime.Value
	}{
		{name: "merge", apply: objects.MergeMutable, args: []runtime.Value{patch}},
		{name: "merge_deep", apply: objects.MergeDeepMutable, args: []runtime.Value{patch}},
		{name: "keep_keys", apply: objects.KeepKeysMutable, args: []runtime.Value{runtime.String("nested")}},
		{name: "omit_keys", apply: objects.OmitKeysMutable, args: []runtime.Value{runtime.String("name")}},
	} {
		b.Run(test.name, func(b *testing.B) {
			args := make([]runtime.Value, len(test.args)+1)
			copy(args[1:], test.args)
			b.ReportAllocs()
			for b.Loop() {
				b.StopTimer()
				args[0] = runtime.NewObjectWith(map[string]runtime.Value{
					"name":   runtime.String("Ferret"),
					"nested": runtime.NewObjectWith(map[string]runtime.Value{"left": runtime.Int(1)}),
					"items":  runtime.NewArrayWith(runtime.Int(1), runtime.Int(2)),
				})
				b.StartTimer()

				if _, err := test.apply(ctx, args...); err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}
