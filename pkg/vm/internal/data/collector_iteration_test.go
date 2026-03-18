package data_test

import (
	"context"
	"errors"
	"io"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/vm/internal/data"
)

func TestDecodeAggregateKeySupportsDirectValues(t *testing.T) {
	groupKey, idx, ok := data.DecodeAggregateKey(data.NewAggregateKey(runtime.NewString(""), 0))
	if !ok {
		t.Fatal("expected aggregate key to decode successfully")
	}

	groupKeyString, ok := groupKey.(runtime.String)
	if !ok {
		t.Fatalf("expected runtime.String group key, got %T", groupKey)
	}

	if groupKeyString.String() != "" {
		t.Fatalf("expected empty-string group key, got %q", groupKeyString.String())
	}

	if idx != 0 {
		t.Fatalf("expected selector index 0, got %d", idx)
	}
}

func TestCounterCollectorIterateReturnsSingleCountEntry(t *testing.T) {
	ctx := context.Background()
	collector := data.NewCounterCollector()

	if err := collector.Set(ctx, runtime.None, runtime.None); err != nil {
		t.Fatalf("first set: %v", err)
	}

	if err := collector.Set(ctx, runtime.None, runtime.None); err != nil {
		t.Fatalf("second set: %v", err)
	}

	iter, err := collector.Iterate(ctx)
	if err != nil {
		t.Fatalf("iterate: %v", err)
	}

	value, key, err := iter.Next(ctx)
	if err != nil {
		t.Fatalf("next: %v", err)
	}

	count, ok := value.(runtime.Int)
	if !ok {
		t.Fatalf("expected runtime.Int count value, got %T", value)
	}

	if count != 2 {
		t.Fatalf("expected count 2, got %d", count)
	}

	if key != runtime.ZeroInt {
		t.Fatalf("expected zero key, got %v", key)
	}

	_, _, err = iter.Next(ctx)
	if !errors.Is(err, io.EOF) {
		t.Fatalf("expected iterator EOF after single count entry, got %v", err)
	}
}

func TestAggregateCollectorIterateReturnsAggregatesThenSortedGroups(t *testing.T) {
	ctx := context.Background()
	collector := data.NewAggregateCollector(bytecode.NewAggregatePlan(
		[]runtime.String{runtime.NewString("cnt"), runtime.NewString("sum")},
		[]bytecode.AggregateKind{bytecode.AggregateCount, bytecode.AggregateSum},
	))

	if err := collector.Set(ctx, runtime.NewString("cnt"), runtime.NewString("row")); err != nil {
		t.Fatalf("set count aggregate: %v", err)
	}

	if err := collector.Set(ctx, runtime.NewString("sum"), runtime.NewInt(7)); err != nil {
		t.Fatalf("set sum aggregate: %v", err)
	}

	if err := collector.Set(ctx, runtime.NewString("b"), runtime.NewString("row-b")); err != nil {
		t.Fatalf("set group b: %v", err)
	}

	if err := collector.Set(ctx, runtime.NewString("a"), runtime.NewString("row-a")); err != nil {
		t.Fatalf("set group a: %v", err)
	}

	iter, err := collector.Iterate(ctx)
	if err != nil {
		t.Fatalf("iterate: %v", err)
	}

	assertAggregateIterEntry(t, ctx, iter, runtime.NewString("cnt"), runtime.NewInt(1))
	assertAggregateIterEntry(t, ctx, iter, runtime.NewString("sum"), runtime.NewInt(7))
	assertAggregateIterGroup(t, ctx, iter, runtime.NewString("a"), runtime.NewString("row-a"))
	assertAggregateIterGroup(t, ctx, iter, runtime.NewString("b"), runtime.NewString("row-b"))

	_, _, err = iter.Next(ctx)
	if !errors.Is(err, io.EOF) {
		t.Fatalf("expected iterator EOF after aggregate and group entries, got %v", err)
	}
}

func assertAggregateIterEntry(t *testing.T, ctx context.Context, iter runtime.Iterator, wantKey runtime.String, wantValue runtime.Value) {
	t.Helper()

	value, key, err := iter.Next(ctx)
	if err != nil {
		t.Fatalf("next aggregate entry: %v", err)
	}

	if key != wantKey {
		t.Fatalf("expected key %q, got %v", wantKey, key)
	}

	if value.String() != wantValue.String() {
		t.Fatalf("expected value %v, got %v", wantValue, value)
	}
}

func assertAggregateIterGroup(t *testing.T, ctx context.Context, iter runtime.Iterator, wantKey runtime.String, wantItem runtime.Value) {
	t.Helper()

	value, key, err := iter.Next(ctx)
	if err != nil {
		t.Fatalf("next group entry: %v", err)
	}

	if key != wantKey {
		t.Fatalf("expected group key %q, got %v", wantKey, key)
	}

	group, ok := value.(runtime.List)
	if !ok {
		t.Fatalf("expected runtime.List group value, got %T", value)
	}

	if got := listLength(t, ctx, group); got != 1 {
		t.Fatalf("expected group size 1, got %d", got)
	}

	item, err := group.At(ctx, 0)
	if err != nil {
		t.Fatalf("group item: %v", err)
	}

	if item != wantItem {
		t.Fatalf("expected group item %v, got %v", wantItem, item)
	}
}
