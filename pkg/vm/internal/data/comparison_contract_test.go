package data_test

import (
	"context"
	"errors"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/vm/internal/data"
)

type opaqueComparisonValue struct {
	equalErr   error
	compareErr error
	id         int
	hash       uint64
}

func (*opaqueComparisonValue) String() string {
	panic("opaque comparison value must not be stringified")
}

func (v *opaqueComparisonValue) Hash() uint64 {
	return v.hash
}

func (v *opaqueComparisonValue) Copy() runtime.Value {
	return v
}

func (v *opaqueComparisonValue) Equal(ctx context.Context, other runtime.Value) (bool, error) {
	if err := ctx.Err(); err != nil {
		return false, err
	}

	if v.equalErr != nil {
		return false, v.equalErr
	}

	otherValue, ok := other.(*opaqueComparisonValue)
	if !ok {
		return false, nil
	}

	return v.id == otherValue.id, nil
}

func (v *opaqueComparisonValue) Compare(ctx context.Context, other runtime.Value) (runtime.Ordering, error) {
	if err := ctx.Err(); err != nil {
		return runtime.Equal, err
	}

	if v.compareErr != nil {
		return runtime.Equal, v.compareErr
	}

	otherValue, ok := other.(*opaqueComparisonValue)
	if !ok {
		return runtime.Equal, runtime.Error(runtime.ErrInvalidOperation, "incompatible opaque values")
	}

	switch {
	case v.id < otherValue.id:
		return runtime.Less, nil
	case v.id > otherValue.id:
		return runtime.Greater, nil
	default:
		return runtime.Equal, nil
	}
}

func TestDataSetUsesListSemanticsAndBackingHash(t *testing.T) {
	ctx := context.Background()
	set := data.NewDataSet(false)
	array := runtime.NewArrayOf([]runtime.Value{runtime.NewInt(1), runtime.NewString("two")})

	for _, value := range []runtime.Value{runtime.NewInt(1), runtime.NewString("two")} {
		if err := set.Append(ctx, value); err != nil {
			t.Fatalf("append dataset value: %v", err)
		}
	}

	for _, operands := range [][2]runtime.Value{{set, array}, {array, set}} {
		equal, err := runtime.EqualValues(ctx, operands[0], operands[1])
		if err != nil {
			t.Fatalf("equal dataset and array: %v", err)
		}
		if !equal {
			t.Fatal("expected dataset and array to be structurally equal")
		}

		ordering, err := runtime.CompareValues(ctx, operands[0], operands[1])
		if err != nil {
			t.Fatalf("compare dataset and array: %v", err)
		}
		if ordering != runtime.Equal {
			t.Fatalf("expected equal ordering, got %v", ordering)
		}
	}

	if set.Hash() != array.Hash() {
		t.Fatalf("expected dataset hash %d to match backing-list hash %d", set.Hash(), array.Hash())
	}
}

func TestFastObjectUsesObjectLikeSemantics(t *testing.T) {
	ctx := context.Background()
	fast := data.NewFastObject(nil, 0)
	object := runtime.NewObject()

	for key, value := range map[string]runtime.Value{
		"a": runtime.NewInt(1),
		"b": runtime.NewString("two"),
	} {
		if err := fast.Set(ctx, runtime.NewString(key), value); err != nil {
			t.Fatalf("set fast object: %v", err)
		}
		if err := object.Set(ctx, runtime.NewString(key), value); err != nil {
			t.Fatalf("set runtime object: %v", err)
		}
	}

	for _, operands := range [][2]runtime.Value{{fast, object}, {object, fast}} {
		equal, err := runtime.EqualValues(ctx, operands[0], operands[1])
		if err != nil {
			t.Fatalf("equal fast object and object: %v", err)
		}
		if !equal {
			t.Fatal("expected fast object and runtime object to be structurally equal")
		}

		ordering, err := runtime.CompareValues(ctx, operands[0], operands[1])
		if err != nil {
			t.Fatalf("compare fast object and object: %v", err)
		}
		if ordering != runtime.Equal {
			t.Fatalf("expected equal ordering, got %v", ordering)
		}
	}

	if fast.Hash() != object.Hash() {
		t.Fatalf("expected fast-object hash %d to match object hash %d", fast.Hash(), object.Hash())
	}
}

func TestKeyCollectorsKeepIntAndStringKeysDistinct(t *testing.T) {
	ctx := context.Background()
	keyCollector := data.NewKeyCollector()
	keyCounter := data.NewKeyCounterCollector()

	for _, key := range []runtime.Value{runtime.NewInt(1), runtime.NewString("1")} {
		if err := keyCollector.Set(ctx, key, runtime.None); err != nil {
			t.Fatalf("collect key: %v", err)
		}
		if err := keyCounter.Set(ctx, key, runtime.None); err != nil {
			t.Fatalf("count key: %v", err)
		}
	}

	for name, collector := range map[string]runtime.Measurable{
		"key":     keyCollector,
		"counter": keyCounter,
	} {
		length, err := collector.Length(ctx)
		if err != nil {
			t.Fatalf("%s collector length: %v", name, err)
		}
		if length != 2 {
			t.Fatalf("expected %s collector to keep Int(1) and String(1) distinct, got %d entries", name, length)
		}
	}
}

func TestKeyCollectorsDoNotStringifyOpaqueKeys(t *testing.T) {
	ctx := context.Background()
	keyCollector := data.NewKeyCollector()
	keyCounter := data.NewKeyCounterCollector()
	first := &opaqueComparisonValue{id: 1, hash: 17}
	same := &opaqueComparisonValue{id: 1, hash: 17}

	for _, collector := range []data.Transformer{keyCollector, keyCounter} {
		if err := collector.Set(ctx, first, runtime.None); err != nil {
			t.Fatalf("set first opaque key: %v", err)
		}
		if err := collector.Set(ctx, same, runtime.None); err != nil {
			t.Fatalf("set duplicate opaque key: %v", err)
		}
		if _, err := collector.Get(ctx, same); err != nil {
			t.Fatalf("get opaque key: %v", err)
		}
	}
}

func TestGroupingPropagatesEqualityErrors(t *testing.T) {
	sentinel := errors.New("equality failed")
	ctx := context.Background()
	collector := data.NewKeyGroupCollector()
	first := &opaqueComparisonValue{id: 1, hash: 23}
	second := &opaqueComparisonValue{id: 2, hash: 23, equalErr: sentinel}

	if err := collector.Set(ctx, first, runtime.NewInt(1)); err != nil {
		t.Fatalf("set first key: %v", err)
	}
	if err := collector.Set(ctx, second, runtime.NewInt(2)); !errors.Is(err, sentinel) {
		t.Fatalf("expected equality error, got %v", err)
	}

	if err := collector.Set(ctx, second, runtime.NewInt(2)); !errors.Is(err, sentinel) {
		t.Fatalf("expected repeated equality error, got %v", err)
	}
}

func TestGroupingDoesNotPollCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	collector := data.NewKeyGroupCollector()
	for _, key := range []runtime.Value{runtime.NewInt(1), runtime.NewInt(2), runtime.NewInt(3)} {
		if err := collector.Set(ctx, key, key); err != nil {
			t.Fatalf("set key %v with canceled context: %v", key, err)
		}
	}

	value, err := collector.Get(ctx, runtime.NewInt(2))
	if err != nil {
		t.Fatalf("get key with canceled context: %v", err)
	}
	group, ok := value.(runtime.List)
	if !ok {
		t.Fatalf("group type = %T, want runtime.List", value)
	}
	length, err := group.Length(ctx)
	if err != nil || length != 1 {
		t.Fatalf("group length = %d, %v, want 1, nil", length, err)
	}
}

func TestCoreComparisonHelpersDoNotPollCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	left := data.NewFastObject(nil, 0)
	right := data.NewFastObject(nil, 0)
	for _, object := range []*data.FastObject{left, right} {
		if err := object.Set(ctx, runtime.NewString("value"), runtime.NewInt(1)); err != nil {
			t.Fatalf("set fast object value: %v", err)
		}
	}

	equal, err := runtime.EqualValues(ctx, left, right)
	if err != nil || !equal {
		t.Fatalf("fast object equality = %v, %v, want true, nil", equal, err)
	}
	ordering, err := runtime.CompareValues(ctx, left, right)
	if err != nil || ordering != runtime.Equal {
		t.Fatalf("fast object ordering = %v, %v, want Equal, nil", ordering, err)
	}

	leftRegexp, err := runtime.NewRegexp("a+")
	if err != nil {
		t.Fatal(err)
	}
	rightRegexp, err := runtime.NewRegexp("a+")
	if err != nil {
		t.Fatal(err)
	}
	equal, err = runtime.EqualValues(ctx, leftRegexp, rightRegexp)
	if err != nil || !equal {
		t.Fatalf("regexp equality = %v, %v, want true, nil", equal, err)
	}
}

func TestCollectorSortingPropagatesOrderingErrors(t *testing.T) {
	sentinel := errors.New("ordering failed")
	ctx := context.Background()
	collector := data.NewKeyGroupCollector()
	first := &opaqueComparisonValue{id: 1, hash: 29, compareErr: sentinel}
	second := &opaqueComparisonValue{id: 2, hash: 31, compareErr: sentinel}

	if err := collector.Set(ctx, first, runtime.NewInt(1)); err != nil {
		t.Fatalf("set first key: %v", err)
	}
	if err := collector.Set(ctx, second, runtime.NewInt(2)); err != nil {
		t.Fatalf("set second key: %v", err)
	}
	if _, err := collector.Iterate(ctx); !errors.Is(err, sentinel) {
		t.Fatalf("expected ordering error, got %v", err)
	}
}

func TestMultiSorterPropagatesSortKeyReadErrors(t *testing.T) {
	sentinel := errors.New("read sort key")
	ctx := context.Background()
	sorter := data.NewMultiSorter([]runtime.SortDirection{runtime.SortDirectionAsc})
	first := &failingAtList{List: runtime.NewArrayOf([]runtime.Value{runtime.NewInt(1)}), err: sentinel}
	second := &failingAtList{List: runtime.NewArrayOf([]runtime.Value{runtime.NewInt(1)}), err: sentinel}

	if err := sorter.Set(ctx, first, runtime.NewString("a")); err != nil {
		t.Fatalf("set first row: %v", err)
	}
	if err := sorter.Set(ctx, second, runtime.NewString("b")); err != nil {
		t.Fatalf("set second row: %v", err)
	}

	if _, err := sorter.Iterate(ctx); !errors.Is(err, sentinel) {
		t.Fatalf("expected sort-key read error, got %v", err)
	}
}
