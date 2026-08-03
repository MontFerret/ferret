package data_test

import (
	"context"
	"errors"
	"io"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/vm/internal/data"
)

func TestKeyCounterKeepsSemanticEntriesAfterSorting(t *testing.T) {
	ctx := context.Background()
	collector := data.NewKeyCounterCollector()

	for _, key := range []runtime.Value{runtime.NewString("b"), runtime.NewString("a")} {
		if err := collector.Set(ctx, key, runtime.None); err != nil {
			t.Fatalf("set initial key: %v", err)
		}
	}

	if _, err := collector.Iterate(ctx); err != nil {
		t.Fatalf("sort initial keys: %v", err)
	}

	if err := collector.Set(ctx, runtime.NewString("b"), runtime.None); err != nil {
		t.Fatalf("increment key after sorting: %v", err)
	}

	index, err := collector.Get(ctx, runtime.NewString("b"))
	if err != nil {
		t.Fatalf("lookup key after sorting: %v", err)
	}
	if index != runtime.NewInt(1) {
		t.Fatalf("expected sorted index 1 for key b, got %v", index)
	}

	iter, err := collector.Iterate(ctx)
	if err != nil {
		t.Fatalf("iterate updated counts: %v", err)
	}

	want := map[runtime.String]runtime.Int{
		runtime.NewString("a"): 1,
		runtime.NewString("b"): 2,
	}
	for {
		value, key, err := iter.Next(ctx)
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			t.Fatalf("iterate count: %v", err)
		}

		stringKey := key.(runtime.String)
		if value != want[stringKey] {
			t.Fatalf("count for %s: got %v, want %v", stringKey, value, want[stringKey])
		}
	}
}

func TestSorterResortsAfterAppending(t *testing.T) {
	ctx := context.Background()
	sorter := data.NewSorter(runtime.SortDirectionAsc)

	if err := sorter.Set(ctx, runtime.NewInt(2), runtime.NewString("two")); err != nil {
		t.Fatalf("set initial value: %v", err)
	}
	if _, err := sorter.Iterate(ctx); err != nil {
		t.Fatalf("sort initial value: %v", err)
	}
	if err := sorter.Set(ctx, runtime.NewInt(1), runtime.NewString("one")); err != nil {
		t.Fatalf("append value after sorting: %v", err)
	}

	iter, err := sorter.Iterate(ctx)
	if err != nil {
		t.Fatalf("resort values: %v", err)
	}
	_, key, err := iter.Next(ctx)
	if err != nil {
		t.Fatalf("read first sorted value: %v", err)
	}
	if key != runtime.NewInt(1) {
		t.Fatalf("expected key 1 first after resort, got %v", key)
	}
}

func TestMultiSorterResortsAfterAppending(t *testing.T) {
	ctx := context.Background()
	sorter := data.NewMultiSorter([]runtime.SortDirection{runtime.SortDirectionAsc})
	key := func(value runtime.Int) runtime.Value {
		return runtime.NewArrayWith(value)
	}

	if err := sorter.Set(ctx, key(2), runtime.NewString("two")); err != nil {
		t.Fatalf("set initial value: %v", err)
	}
	if _, err := sorter.Iterate(ctx); err != nil {
		t.Fatalf("sort initial value: %v", err)
	}
	if err := sorter.Set(ctx, key(1), runtime.NewString("one")); err != nil {
		t.Fatalf("append value after sorting: %v", err)
	}

	iter, err := sorter.Iterate(ctx)
	if err != nil {
		t.Fatalf("resort values: %v", err)
	}
	_, firstKey, err := iter.Next(ctx)
	if err != nil {
		t.Fatalf("read first sorted value: %v", err)
	}
	firstPart, err := firstKey.(runtime.List).At(ctx, runtime.ZeroInt)
	if err != nil {
		t.Fatalf("read first multi-sort key: %v", err)
	}
	if firstPart != runtime.NewInt(1) {
		t.Fatalf("expected key [1] first after resort, got %v", firstPart)
	}
}
