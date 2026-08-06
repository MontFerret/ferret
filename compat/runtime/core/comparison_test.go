package core_test

import (
	"context"
	"encoding/json"
	"errors"
	"strconv"
	"testing"

	"github.com/MontFerret/ferret/v2/compat/runtime/core"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

type legacyComparisonValue struct {
	typeName string
	value    int64
}

func (v *legacyComparisonValue) Type() core.Type {
	name := v.typeName
	if name == "" {
		name = "legacy-comparison"
	}

	return core.NewType(name)
}

func (v *legacyComparisonValue) String() string {
	return strconv.FormatInt(v.value, 10)
}

func (v *legacyComparisonValue) Compare(other core.Value) int64 {
	otherValue := other.(*legacyComparisonValue)
	switch {
	case v.value < otherValue.value:
		return -99
	case v.value > otherValue.value:
		return 99
	default:
		return 0
	}
}

func (v *legacyComparisonValue) Unwrap() interface{} {
	return v.value
}

func (v *legacyComparisonValue) Hash() uint64 {
	return uint64(v.value)
}

func (v *legacyComparisonValue) Copy() core.Value {
	return &legacyComparisonValue{value: v.value, typeName: v.typeName}
}

func (v *legacyComparisonValue) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func TestCustomValueBridgeImplementsFallibleCapabilities(t *testing.T) {
	left := core.UnwrapValue(&legacyComparisonValue{value: 1})
	right := core.UnwrapValue(&legacyComparisonValue{value: 2})
	same := core.UnwrapValue(&legacyComparisonValue{value: 1})

	if _, ok := left.(runtime.Equatable); !ok {
		t.Fatalf("custom v1 value %T does not implement Equatable", left)
	}
	if _, ok := left.(runtime.Comparable); !ok {
		t.Fatalf("custom v1 value %T does not implement Comparable", left)
	}

	equal, err := runtime.EqualValues(t.Context(), left, same)
	if err != nil || !equal {
		t.Fatalf("equal custom values: equal=%v err=%v", equal, err)
	}
	ordering, err := runtime.CompareValues(t.Context(), left, right)
	if err != nil || ordering != runtime.Less {
		t.Fatalf("compare custom values: ordering=%v err=%v", ordering, err)
	}
	ordering, err = runtime.CompareValues(t.Context(), right, left)
	if err != nil || ordering != runtime.Greater {
		t.Fatalf("reverse custom comparison: ordering=%v err=%v", ordering, err)
	}
	ordering, err = runtime.CompareValues(t.Context(), left, same)
	if err != nil || ordering != runtime.Equal {
		t.Fatalf("equal custom comparison: ordering=%v err=%v", ordering, err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	_, err = left.(runtime.Comparable).Compare(ctx, right)
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("canceled comparison error = %v", err)
	}
	_, err = left.(runtime.Equatable).Equal(ctx, same)
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("canceled equality error = %v", err)
	}
}

func TestV1CompareCollapsesStrictErrorsToLegacyTypeRank(t *testing.T) {
	left := core.WrapValue(runtime.NewBox(1))
	right := core.WrapValue(runtime.NewBox(2))

	if comparison := left.Compare(right); comparison != 0 {
		t.Fatalf("unsupported same-rank comparison = %d, want legacy zero", comparison)
	}

	if comparison := core.WrapValue(runtime.NewInt(1)).Compare(core.WrapValue(runtime.NewInt(2))); comparison != -1 {
		t.Fatalf("native comparison = %d, want -1", comparison)
	}
}

func TestCustomValueBridgePreservesLegacyTypeDomains(t *testing.T) {
	left := core.UnwrapValue(&legacyComparisonValue{value: 1, typeName: "left"})
	right := core.UnwrapValue(&legacyComparisonValue{value: 1, typeName: "right"})

	if runtime.TypeName(runtime.TypeOf(left)) == runtime.TypeName(runtime.TypeOf(right)) {
		t.Fatalf("distinct legacy types share runtime domain %q", runtime.TypeName(runtime.TypeOf(left)))
	}

	equal, err := left.(runtime.Equatable).Equal(t.Context(), right)
	if err != nil || equal {
		t.Fatalf("mismatched legacy equality: equal=%v err=%v", equal, err)
	}
	_, err = left.(runtime.Comparable).Compare(t.Context(), right)
	if !errors.Is(err, runtime.ErrInvalidOperation) {
		t.Fatalf("mismatched legacy ordering error = %v, want invalid operation", err)
	}
}
