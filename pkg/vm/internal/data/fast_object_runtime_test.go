package data_test

import (
	"context"
	"strings"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/vm/internal/data"
)

func TestFastObjectGetMissingReturnsNoneWithoutError(t *testing.T) {
	obj := data.NewFastObject(nil, 0)

	val, err := obj.Get(context.Background(), runtime.NewString("missing"))
	if err != nil {
		t.Fatalf("expected missing fast-object key to return no error, got %v", err)
	}

	if val != runtime.None {
		t.Fatalf("expected runtime.None for missing fast-object key, got %v", val)
	}
}

func TestFastObjectDictionaryBackingLifetime(t *testing.T) {
	ctx := t.Context()
	object := data.NewFastObject(nil, 1)
	object.SetStringCached("keep", runtime.Int(1))
	longKey := runtime.String(strings.Repeat("x", 1<<17))
	object.SetStringCached(string(longKey), runtime.True)
	if err := object.RemoveKey(ctx, longKey); err != nil {
		t.Fatal(err)
	}

	if object.DebugFieldsPageable() {
		t.Fatal("original dictionary lost its key history")
	}

	copied := object.Copy().(*data.FastObject)
	cloned, err := object.Clone(ctx)
	if err != nil {
		t.Fatal(err)
	}

	iter, err := object.Iterate(ctx)
	if err != nil {
		t.Fatal(err)
	}

	alias := *object
	object.SetStringCached("keep", runtime.Int(2))
	if err := object.Clear(ctx); err != nil {
		t.Fatal(err)
	}

	object.SetStringCached("keep", runtime.Int(3))
	value, key, err := iter.Next(ctx)
	if err != nil || key != runtime.String("keep") || value != runtime.Int(2) {
		t.Fatalf("iterator after clear: %v=%v, %v", key, value, err)
	}

	value, found, err := alias.Lookup(ctx, runtime.String("keep"))
	if err != nil || !found || value != runtime.Int(2) || alias.DebugFieldsPageable() {
		t.Fatalf("old backing after clear: %v, found %t, pageable %t, %v", value, found, alias.DebugFieldsPageable(), err)
	}

	for _, snapshot := range []*data.FastObject{copied, cloned.(*data.FastObject)} {
		value, found, err := snapshot.Lookup(ctx, runtime.String("keep"))
		if err != nil || !found || value != runtime.Int(1) || !snapshot.DebugFieldsPageable() {
			t.Fatalf("independent dictionary: %v, found %t, pageable %t, %v", value, found, snapshot.DebugFieldsPageable(), err)
		}
	}
}
