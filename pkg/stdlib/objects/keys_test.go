package objects_test

import (
	"context"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/stdlib/objects"
)

func TestKeys(t *testing.T) {
	ctx := context.Background()
	for _, source := range []*runtime.Object{
		runtime.NewObject(),
		runtime.NewObjectWith(map[string]runtime.Value{"a": runtime.Int(1), "é 😀": runtime.None, "a.b": runtime.True}),
	} {
		result, err := objects.Keys(ctx, source)
		if err != nil {
			t.Fatal(err)
		}

		list := result.(runtime.List)
		got, err := list.Length(ctx)
		if err != nil {
			t.Fatal(err)
		}

		want, err := source.Length(ctx)
		if err != nil || got != want {
			t.Fatalf("length = %d, want %d: %v", got, want, err)
		}

		err = list.ForEach(ctx, func(ctx context.Context, key runtime.Value, _ runtime.Int) (runtime.Boolean, error) {
			present, err := source.ContainsKey(ctx, key)
			if err != nil || !present {
				t.Fatalf("unexpected key %v: %v", key, err)
			}

			return true, nil
		})
		if err != nil {
			t.Fatal(err)
		}

		if err := list.Append(ctx, runtime.String("extra")); err != nil {
			t.Fatal(err)
		}

		present, err := source.ContainsKey(ctx, runtime.String("extra"))
		if err != nil || present {
			t.Fatalf("keys alias source: %v", err)
		}
	}

	result, err := objects.Keys(ctx, runtime.Int(1))
	if err == nil || result != runtime.None {
		t.Fatalf("invalid map: %v, %v", result, err)
	}
}
