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

func TestAggregateCollectorAllowsEmptyStringGroupKey(t *testing.T) {
	ctx := context.Background()
	collector := data.NewAggregateCollector(bytecode.NewAggregatePlan(
		[]runtime.String{runtime.NewString("cnt")},
		[]bytecode.AggregateKind{bytecode.AggregateCount},
		false,
	))

	if err := collector.Set(ctx, runtime.NewString("cnt"), runtime.NewInt(1)); err != nil {
		t.Fatalf("set aggregate value: %v", err)
	}

	if err := collector.Set(ctx, runtime.NewString(""), runtime.NewString("a")); err != nil {
		t.Fatalf("set empty-string group key (first): %v", err)
	}

	if err := collector.Set(ctx, runtime.NewString(""), runtime.NewString("b")); err != nil {
		t.Fatalf("set empty-string group key (second): %v", err)
	}

	group, err := collector.Get(ctx, runtime.NewString(""))
	if err != nil {
		t.Fatalf("get empty-string group key: %v", err)
	}

	if got := listLength(t, ctx, group); got != 2 {
		t.Fatalf("expected empty-string group size 2, got %d", got)
	}

	if got := measurableLength(t, ctx, collector); got != 2 {
		t.Fatalf("expected collector length 2 (1 aggregate + 1 group), got %d", got)
	}
}

func TestKeyGroupCollectorAllowsEmptyStringKey(t *testing.T) {
	ctx := context.Background()
	collector := data.NewKeyGroupCollector()

	if err := collector.Set(ctx, runtime.NewString(""), runtime.NewInt(1)); err != nil {
		t.Fatalf("set empty-string key (first): %v", err)
	}

	if err := collector.Set(ctx, runtime.NewString(""), runtime.NewInt(2)); err != nil {
		t.Fatalf("set empty-string key (second): %v", err)
	}

	group, err := collector.Get(ctx, runtime.NewString(""))
	if err != nil {
		t.Fatalf("get empty-string key group: %v", err)
	}

	if got := listLength(t, ctx, group); got != 2 {
		t.Fatalf("expected empty-string group size 2, got %d", got)
	}

	if got := measurableLength(t, ctx, collector); got != 1 {
		t.Fatalf("expected one grouped entry, got %d", got)
	}
}

func TestGroupedAggregateCollectorAllowsEmptyStringKey(t *testing.T) {
	ctx := context.Background()
	collector := data.NewGroupedAggregateCollector(bytecode.NewAggregatePlan(
		[]runtime.String{runtime.NewString("cnt")},
		[]bytecode.AggregateKind{bytecode.AggregateCount},
		true,
	))

	if err := collector.Set(ctx, runtime.NewString(""), runtime.NewString("row1")); err != nil {
		t.Fatalf("set grouped value: %v", err)
	}

	aggKey := data.NewAggregateKey(runtime.NewString(""), 0)

	if err := collector.Set(ctx, aggKey, runtime.NewInt(10)); err != nil {
		t.Fatalf("set aggregate update for empty-string group key: %v", err)
	}

	if got := measurableLength(t, ctx, collector); got != 1 {
		t.Fatalf("expected one grouped entry after aggregate update, got %d", got)
	}

	countValue, err := collector.Get(ctx, aggKey)
	if err != nil {
		t.Fatalf("get aggregate value for empty-string group key: %v", err)
	}

	count, ok := countValue.(runtime.Int)
	if !ok {
		t.Fatalf("expected runtime.Int count value, got %T", countValue)
	}

	if count != 1 {
		t.Fatalf("expected count 1, got %d", count)
	}
}

func TestKeyCounterCollectorAllowsEmptyStringKey(t *testing.T) {
	ctx := context.Background()
	collector := data.NewKeyCounterCollector()

	if err := collector.Set(ctx, runtime.NewString(""), runtime.NewInt(1)); err != nil {
		t.Fatalf("set empty-string key (first): %v", err)
	}

	if err := collector.Set(ctx, runtime.NewString(""), runtime.NewInt(1)); err != nil {
		t.Fatalf("set empty-string key (second): %v", err)
	}

	if got := measurableLength(t, ctx, collector); got != 1 {
		t.Fatalf("expected one grouped entry, got %d", got)
	}

	iter, err := collector.Iterate(ctx)
	if err != nil {
		t.Fatalf("iterate: %v", err)
	}

	value, key, err := iter.Next(ctx)
	if errors.Is(err, io.EOF) {
		t.Fatal("expected one iterator item")
	}
	if err != nil {
		t.Fatalf("next: %v", err)
	}

	keyString, ok := key.(runtime.String)
	if !ok {
		t.Fatalf("expected runtime.String key, got %T", key)
	}

	if keyString.String() != "" {
		t.Fatalf("expected empty-string key, got %q", keyString.String())
	}

	count, ok := value.(runtime.Int)
	if !ok {
		t.Fatalf("expected runtime.Int count value, got %T", value)
	}

	if count != 2 {
		t.Fatalf("expected count 2, got %d", count)
	}
}

func listLength(t *testing.T, ctx context.Context, value runtime.Value) runtime.Int {
	t.Helper()

	list, ok := value.(runtime.List)
	if !ok {
		t.Fatalf("expected runtime.List, got %T", value)
	}

	length, err := list.Length(ctx)
	if err != nil {
		t.Fatalf("list length: %v", err)
	}

	return length
}

func measurableLength(t *testing.T, ctx context.Context, value runtime.Measurable) runtime.Int {
	t.Helper()

	length, err := value.Length(ctx)
	if err != nil {
		t.Fatalf("length: %v", err)
	}

	return length
}
