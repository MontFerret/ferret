package data_test

import (
	"context"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/vm/internal/data"
)

type collisionValue struct {
	label string
	hash  uint64
}

func (v collisionValue) String() string {
	return v.label
}

func (v collisionValue) Hash() uint64 {
	return v.hash
}

func (v collisionValue) Copy() runtime.Value {
	return v
}

func (v collisionValue) Compare(other runtime.Value) int {
	o, ok := other.(collisionValue)
	if !ok {
		return runtime.CompareTypes(v, other)
	}

	switch {
	case v.label < o.label:
		return -1
	case v.label > o.label:
		return 1
	default:
		return 0
	}
}

func TestKeyGroupCollectorSeparatesHashCollisions(t *testing.T) {
	ctx := context.Background()
	collector := data.NewKeyGroupCollector()
	first := collisionValue{hash: 7, label: "first"}
	second := collisionValue{hash: 7, label: "second"}

	if err := collector.Set(ctx, first, runtime.NewInt(1)); err != nil {
		t.Fatalf("set first: %v", err)
	}

	if err := collector.Set(ctx, second, runtime.NewInt(2)); err != nil {
		t.Fatalf("set second: %v", err)
	}

	assertGroupedIntValue(t, ctx, collector, first, 1)
	assertGroupedIntValue(t, ctx, collector, second, 2)
}

func TestGroupedAggregateCollectorSeparatesHashCollisions(t *testing.T) {
	ctx := context.Background()
	collector := data.NewGroupedAggregateCollector(bytecode.NewAggregatePlan(
		[]runtime.String{runtime.NewString("cnt")},
		[]bytecode.AggregateKind{bytecode.AggregateCount},
		true,
	))

	first := collisionValue{hash: 11, label: "first"}
	second := collisionValue{hash: 11, label: "second"}

	if err := collector.Set(ctx, first, runtime.NewString("row-a")); err != nil {
		t.Fatalf("set first group row: %v", err)
	}

	if err := collector.Set(ctx, second, runtime.NewString("row-b")); err != nil {
		t.Fatalf("set second group row: %v", err)
	}

	if err := collector.Set(ctx, data.NewAggregateKey(first, 0), runtime.NewInt(1)); err != nil {
		t.Fatalf("set first aggregate: %v", err)
	}

	if err := collector.Set(ctx, data.NewAggregateKey(second, 0), runtime.NewInt(1)); err != nil {
		t.Fatalf("set second aggregate: %v", err)
	}

	assertGroupedStringValue(t, ctx, collector, first, "row-a")
	assertGroupedStringValue(t, ctx, collector, second, "row-b")
	assertAggregateIntValue(t, ctx, collector, first, 1)
	assertAggregateIntValue(t, ctx, collector, second, 1)
}

func TestGroupedAggregateCollectorTrackGroupValuesFalseSkipsStoredGroups(t *testing.T) {
	ctx := context.Background()
	key := runtime.NewString("only")
	collector := data.NewGroupedAggregateCollector(bytecode.NewAggregatePlan(
		[]runtime.String{runtime.NewString("cnt")},
		[]bytecode.AggregateKind{bytecode.AggregateCount},
		false,
	))

	if err := collector.Set(ctx, key, runtime.NewString("row")); err != nil {
		t.Fatalf("set grouped row: %v", err)
	}

	if err := collector.Set(ctx, data.NewAggregateKey(key, 0), runtime.NewInt(1)); err != nil {
		t.Fatalf("set aggregate value: %v", err)
	}

	group, err := collector.Get(ctx, key)
	if err != nil {
		t.Fatalf("get group: %v", err)
	}

	if group != runtime.None {
		t.Fatalf("expected trackGroupValues=false group lookup to return NONE, got %v", group)
	}

	iter, err := collector.Iterate(ctx)
	if err != nil {
		t.Fatalf("iterate: %v", err)
	}

	value, iterKey, err := iter.Next(ctx)
	if err != nil {
		t.Fatalf("next: %v", err)
	}

	if value != runtime.None {
		t.Fatalf("expected iterator value NONE when group values are disabled, got %v", value)
	}

	if runtime.CompareValues(iterKey, key) != 0 {
		t.Fatalf("expected iterator key %v, got %v", key, iterKey)
	}

	assertAggregateIntValue(t, ctx, collector, key, 1)
}

func TestGroupedAggregateCollectorTrackGroupValuesTrueKeepsGroups(t *testing.T) {
	ctx := context.Background()
	key := runtime.NewString("only")
	collector := data.NewGroupedAggregateCollector(bytecode.NewAggregatePlan(
		[]runtime.String{runtime.NewString("cnt")},
		[]bytecode.AggregateKind{bytecode.AggregateCount},
		true,
	))

	if err := collector.Set(ctx, key, runtime.NewString("row")); err != nil {
		t.Fatalf("set grouped row: %v", err)
	}

	group, err := collector.Get(ctx, key)
	if err != nil {
		t.Fatalf("get group: %v", err)
	}

	if got := listLength(t, ctx, group); got != 1 {
		t.Fatalf("expected tracked group size 1, got %d", got)
	}
}

func assertGroupedIntValue(t *testing.T, ctx context.Context, collector runtime.KeyReadable, key runtime.Value, want int) {
	t.Helper()

	group, err := collector.Get(ctx, key)
	if err != nil {
		t.Fatalf("get group %v: %v", key, err)
	}

	item, err := group.(runtime.List).At(ctx, 0)
	if err != nil {
		t.Fatalf("read group item: %v", err)
	}

	got, ok := item.(runtime.Int)
	if !ok {
		t.Fatalf("expected runtime.Int group item, got %T", item)
	}

	if got != runtime.NewInt(want) {
		t.Fatalf("expected %d, got %d", want, got)
	}
}

func assertGroupedStringValue(t *testing.T, ctx context.Context, collector runtime.KeyReadable, key runtime.Value, want string) {
	t.Helper()

	group, err := collector.Get(ctx, key)
	if err != nil {
		t.Fatalf("get group %v: %v", key, err)
	}

	item, err := group.(runtime.List).At(ctx, 0)
	if err != nil {
		t.Fatalf("read group item: %v", err)
	}

	got, ok := item.(runtime.String)
	if !ok {
		t.Fatalf("expected runtime.String group item, got %T", item)
	}

	if got.String() != want {
		t.Fatalf("expected %q, got %q", want, got.String())
	}
}

func assertAggregateIntValue(t *testing.T, ctx context.Context, collector runtime.KeyReadable, key runtime.Value, want int) {
	t.Helper()

	value, err := collector.Get(ctx, data.NewAggregateKey(key, 0))
	if err != nil {
		t.Fatalf("get aggregate %v: %v", key, err)
	}

	got, ok := value.(runtime.Int)
	if !ok {
		t.Fatalf("expected runtime.Int aggregate value, got %T", value)
	}

	if got != runtime.NewInt(want) {
		t.Fatalf("expected %d, got %d", want, got)
	}
}
