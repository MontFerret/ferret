package objects_test

import (
	"context"
	"errors"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/stdlib/objects"
)

func TestImmutableMergeUsesFirstMapFactory(t *testing.T) {
	for name, merge := range map[string]func(context.Context, ...runtime.Value) (runtime.Value, error){"shallow": objects.Merge, "deep": objects.MergeDeep} {
		t.Run(name, func(t *testing.T) {
			ctx := t.Context()
			nested := runtime.NewArrayWith(runtime.Int(1))
			first := &configuredMap{Object: runtime.NewObjectWith(map[string]runtime.Value{"nested": nested, "a": runtime.Int(1)}), backend: "remote/utf8/limit=100"}
			second := &configuredMap{Object: runtime.NewObjectWith(map[string]runtime.Value{"a": runtime.Int(2)}), backend: "different"}
			empty, err := first.New(ctx)
			if err != nil {
				t.Fatal(err)
			}

			if size, err := empty.Length(ctx); err != nil || size != 0 {
				t.Fatal("factory copied elements")
			}

			if empty.(*configuredMap).backend != first.backend {
				t.Fatal("factory lost configuration")
			}

			if err := empty.Set(ctx, runtime.String("a"), runtime.Int(9)); err != nil {
				t.Fatal(err)
			}

			if value, err := first.Get(ctx, runtime.String("a")); err != nil || value != runtime.Int(1) {
				t.Fatal("factory aliases source")
			}

			result, err := merge(ctx, first, second)
			if err != nil {
				t.Fatal(err)
			}

			destination, ok := result.(*configuredMap)
			if !ok || destination == first || destination.backend != first.backend || first.creations != 2 || second.creations != 0 {
				t.Fatalf("wrong destination: %T", result)
			}

			if value, err := destination.Get(ctx, runtime.String("a")); err != nil || value != runtime.Int(2) {
				t.Fatal("merge precedence changed")
			}

			value, err := destination.Get(ctx, runtime.String("nested"))
			if err != nil || value == nested {
				t.Fatal("merge did not clone incoming list")
			}

			if err := value.(runtime.List).Append(ctx, runtime.Int(3)); err != nil {
				t.Fatal(err)
			}

			if size, err := nested.Length(ctx); err != nil || size != 1 {
				t.Fatal("merge changed source value")
			}

			if destination.closes != 0 || first.closes != 0 || second.closes != 0 {
				t.Fatal("successful merge closed a result or borrowed source")
			}

			if err := destination.Close(); err != nil {
				t.Fatal(err)
			}

			if err := empty.(*configuredMap).Close(); err != nil {
				t.Fatal(err)
			}

			primary, cleanup := errors.New("factory failed"), errors.New("factory rollback failed")
			first.newErr = errors.Join(primary, cleanup)
			result, err = merge(ctx, first)
			if result != runtime.None || !errors.Is(err, primary) || !errors.Is(err, cleanup) {
				t.Fatalf("factory errors lost: %v, %v", result, err)
			}

			canceled, cancel := context.WithCancel(ctx)
			cancel()
			before := first.creations
			if result, err := first.New(canceled); result != nil || !errors.Is(err, context.Canceled) {
				t.Fatalf("factory cancellation: %v, %v", result, err)
			}

			if result, err := merge(canceled, first); result != runtime.None || !errors.Is(err, context.Canceled) {
				t.Fatalf("merge cancellation: %v, %v", result, err)
			}

			if before != first.creations {
				t.Fatal("pre-canceled construction performed host work")
			}
		})
	}
}
