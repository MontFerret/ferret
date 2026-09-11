package runtime_test

import (
	"context"
	"errors"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

var (
	_ runtime.Factory[runtime.List] = (*runtime.Array)(nil)
	_ runtime.Factory[runtime.Map]  = (*runtime.Object)(nil)
	_ runtime.List                  = (*runtime.Array)(nil)
	_ runtime.Map                   = (*runtime.Object)(nil)
)

func TestNativeCollectionFactories(t *testing.T) {
	ctx := t.Context()
	array := runtime.NewArrayWith(runtime.Int(7))
	list, err := array.New(ctx)
	if err != nil {
		t.Fatal(err)
	}

	if _, ok := list.(*runtime.Array); !ok || list == array {
		t.Fatalf("list family: %T", list)
	}

	if length, err := list.Length(ctx); err != nil || length != 0 {
		t.Fatalf("new list is not empty: %d, %v", length, err)
	}

	if err := list.Append(ctx, runtime.Int(9)); err != nil {
		t.Fatal(err)
	}

	if value, err := array.At(ctx, 0); err != nil || value != runtime.Int(7) {
		t.Fatal("list source changed")
	}

	if err := array.SetAt(ctx, 0, runtime.Int(8)); err != nil {
		t.Fatal(err)
	}

	if value, err := list.At(ctx, 0); err != nil || value != runtime.Int(9) {
		t.Fatal("list destination aliases source")
	}

	object := runtime.NewObjectWith(map[string]runtime.Value{"a": runtime.Int(7)})
	mapping, err := object.New(ctx)
	if err != nil {
		t.Fatal(err)
	}

	if _, ok := mapping.(*runtime.Object); !ok || mapping == object {
		t.Fatalf("map family: %T", mapping)
	}

	if length, err := mapping.Length(ctx); err != nil || length != 0 {
		t.Fatalf("new map is not empty: %d, %v", length, err)
	}

	if err := mapping.Set(ctx, runtime.String("a"), runtime.Int(9)); err != nil {
		t.Fatal(err)
	}

	if value, err := object.Get(ctx, runtime.String("a")); err != nil || value != runtime.Int(7) {
		t.Fatal("map source changed")
	}

	if err := object.Set(ctx, runtime.String("a"), runtime.Int(8)); err != nil {
		t.Fatal(err)
	}

	if value, err := mapping.Get(ctx, runtime.String("a")); err != nil || value != runtime.Int(9) {
		t.Fatal("map destination aliases source")
	}

	canceled, cancel := context.WithCancel(ctx)
	cancel()
	if value, err := array.New(canceled); value != nil || !errors.Is(err, context.Canceled) {
		t.Fatalf("array cancellation: %v, %v", value, err)
	}

	if value, err := object.New(canceled); value != nil || !errors.Is(err, context.Canceled) {
		t.Fatalf("object cancellation: %v, %v", value, err)
	}
}
