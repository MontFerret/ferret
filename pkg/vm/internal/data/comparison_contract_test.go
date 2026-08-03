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

func TestRegexpComparisonUsesPatternAndStableType(t *testing.T) {
	ctx := context.Background()
	left, err := data.NewRegexp(runtime.NewString("a+"))
	if err != nil {
		t.Fatalf("compile left regexp: %v", err)
	}
	right, err := data.NewRegexp(runtime.NewString("b+"))
	if err != nil {
		t.Fatalf("compile right regexp: %v", err)
	}
	same, err := data.NewRegexp(runtime.NewString("a+"))
	if err != nil {
		t.Fatalf("compile same regexp: %v", err)
	}

	equal, err := left.Equal(ctx, same)
	if err != nil || !equal {
		t.Fatalf("expected equal regexp patterns, equal=%v err=%v", equal, err)
	}

	ordering, err := left.Compare(ctx, right)
	if err != nil {
		t.Fatalf("compare regexp patterns: %v", err)
	}
	if ordering != runtime.Less {
		t.Fatalf("expected left regexp to sort before right, got %v", ordering)
	}

	if left.Type() != data.TypeRegexp {
		t.Fatalf("expected stable regexp type %v, got %v", data.TypeRegexp, left.Type())
	}
	if left.Hash() != same.Hash() {
		t.Fatal("expected equal regexp patterns to have equal hashes")
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

func TestGroupingPropagatesEqualityErrorsAndCancellation(t *testing.T) {
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

	cancelled, cancel := context.WithCancel(context.Background())
	cancel()
	if err := collector.Set(cancelled, second, runtime.NewInt(2)); !errors.Is(err, context.Canceled) {
		t.Fatalf("expected cancellation, got %v", err)
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
