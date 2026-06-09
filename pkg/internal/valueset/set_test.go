package valueset_test

import (
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/internal/valueset"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

type collisionValue struct {
	label string
}

func (v collisionValue) String() string {
	return v.label
}

func (v collisionValue) Hash() uint64 {
	return 7
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

func TestSetTracksDistinctValues(t *testing.T) {
	set := valueset.New(4)

	for _, tc := range []struct {
		value runtime.Value
		added bool
	}{
		{runtime.NewInt(1), true},
		{runtime.NewString("1"), true},
		{runtime.NewBoolean(true), true},
		{runtime.NewInt(1), false},
		{runtime.NewString("1"), false},
	} {
		if got := set.Add(tc.value); got != tc.added {
			t.Fatalf("Add(%v): expected %t, got %t", tc.value, tc.added, got)
		}
	}

	if got := set.Len(); got != 3 {
		t.Fatalf("expected length 3, got %d", got)
	}
}

func TestSetUsesFerretEquality(t *testing.T) {
	set := valueset.New(2)
	first := runtime.NewObjectWith(map[string]runtime.Value{
		"name": runtime.NewString("Ada"),
		"role": runtime.NewString("admin"),
	})
	equal := runtime.NewObjectWith(map[string]runtime.Value{
		"role": runtime.NewString("admin"),
		"name": runtime.NewString("Ada"),
	})

	if !set.Add(first) {
		t.Fatal("expected first object to be added")
	}

	if set.Add(equal) {
		t.Fatal("expected equal reordered object to be rejected")
	}

	if got := set.Len(); got != 1 {
		t.Fatalf("expected length 1, got %d", got)
	}
}

func TestSetSeparatesHashCollisions(t *testing.T) {
	set := valueset.New(3)
	first := collisionValue{label: "first"}
	second := collisionValue{label: "second"}

	for _, tc := range []struct {
		value collisionValue
		added bool
	}{
		{first, true},
		{second, true},
		{first, false},
		{second, false},
	} {
		if got := set.Add(tc.value); got != tc.added {
			t.Fatalf("Add(%v): expected %t, got %t", tc.value, tc.added, got)
		}
	}

	if got := set.Len(); got != 2 {
		t.Fatalf("expected length 2, got %d", got)
	}
}
