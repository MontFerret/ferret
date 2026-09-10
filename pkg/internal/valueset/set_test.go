package valueset_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/MontFerret/ferret/v2/pkg/internal/valueset"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

func TestSetPropagatesEqualityErrorsWithoutMutation(t *testing.T) {
	ctx := context.Background()
	sentinel := errors.New("equality failed")
	set := valueset.New(2)

	added, err := set.Add(ctx, collisionValue{label: "first", err: sentinel})
	if err != nil || !added {
		t.Fatalf("add first: added=%t err=%v", added, err)
	}

	added, err = set.Add(ctx, collisionValue{label: "second"})
	if added {
		t.Fatal("failed add must not mutate the set")
	}
	if !errors.Is(err, sentinel) {
		t.Fatalf("expected equality error, got %v", err)
	}
	if set.Len() != 1 {
		t.Fatalf("expected one entry, got %d", set.Len())
	}
}

func TestSetTracksDistinctValues(t *testing.T) {
	ctx := context.Background()
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
		got, err := set.Add(ctx, tc.value)
		if err != nil {
			t.Fatalf("Add(%v): %v", tc.value, err)
		}
		if got != tc.added {
			t.Fatalf("Add(%v): expected %t, got %t", tc.value, tc.added, got)
		}
	}

	if got := set.Len(); got != 3 {
		t.Fatalf("expected length 3, got %d", got)
	}
}

func TestSetUsesFerretEquality(t *testing.T) {
	ctx := context.Background()
	set := valueset.New(2)
	first := runtime.NewObjectWith(map[string]runtime.Value{
		"name": runtime.NewString("Ada"),
		"role": runtime.NewString("admin"),
	})
	equal := runtime.NewObjectWith(map[string]runtime.Value{
		"role": runtime.NewString("admin"),
		"name": runtime.NewString("Ada"),
	})

	added, err := set.Add(ctx, first)
	if err != nil {
		t.Fatalf("add first: %v", err)
	}
	if !added {
		t.Fatal("expected first object to be added")
	}

	added, err = set.Add(ctx, equal)
	if err != nil {
		t.Fatalf("add equal: %v", err)
	}
	if added {
		t.Fatal("expected equal reordered object to be rejected")
	}

	if got := set.Len(); got != 1 {
		t.Fatalf("expected length 1, got %d", got)
	}
}

func TestSetUsesNumericEqualityAcrossIntAndFloat(t *testing.T) {
	set := valueset.New(2)

	added, err := set.Add(t.Context(), runtime.NewInt(1))
	if err != nil || !added {
		t.Fatalf("add Int: added=%t err=%v", added, err)
	}

	added, err = set.Add(t.Context(), runtime.NewFloat(1))
	if err != nil {
		t.Fatalf("add Float: %v", err)
	}
	if added {
		t.Fatal("equal Float must not be added")
	}
	if set.Len() != 1 {
		t.Fatalf("expected one numeric value, got %d", set.Len())
	}
}

func TestSetUsesStrictDurationEquality(t *testing.T) {
	set := valueset.New(4)

	for _, test := range []struct {
		value runtime.Value
		added bool
	}{
		{value: runtime.NewString("1s"), added: true},
		{value: runtime.NewInt(1000), added: true},
		{value: runtime.NewDuration(time.Second), added: true},
		{value: runtime.NewDuration(1000 * time.Millisecond), added: false},
	} {
		added, err := set.Add(t.Context(), test.value)
		if err != nil || added != test.added {
			t.Fatalf("Add(%T) = %t, %v, want %t, nil", test.value, added, err, test.added)
		}
	}

	if set.Len() != 3 {
		t.Fatalf("expected String, Int, and one Duration value, got %d", set.Len())
	}
}

func TestSetSeparatesHashCollisions(t *testing.T) {
	ctx := context.Background()
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
		got, err := set.Add(ctx, tc.value)
		if err != nil {
			t.Fatalf("Add(%v): %v", tc.value, err)
		}
		if got != tc.added {
			t.Fatalf("Add(%v): expected %t, got %t", tc.value, tc.added, got)
		}
	}

	if got := set.Len(); got != 2 {
		t.Fatalf("expected length 2, got %d", got)
	}
}

func TestSetMembershipAndRemoval(t *testing.T) {
	values := []runtime.Value{collisionValue{label: "first"}, collisionValue{label: "second"}, collisionValue{label: "third"}}
	for _, order := range [][]int{{0, 1, 2}, {1, 2, 0}, {2, 0, 1}} {
		set := valueset.New(0)
		for _, value := range values {
			if _, err := set.Add(t.Context(), value); err != nil {
				t.Fatal(err)
			}
		}

		for step, index := range order {
			found, err := set.Contains(t.Context(), values[index])
			if err != nil || !found {
				t.Fatalf("missing value %d before removal: %v", index, err)
			}

			removed, err := set.Remove(t.Context(), values[index])
			if err != nil || !removed || set.Len() != len(values)-step-1 {
				t.Fatalf("remove %d: %t, %v; size=%d", index, removed, err, set.Len())
			}

			found, err = set.Contains(t.Context(), values[index])
			if err != nil || found {
				t.Fatalf("removed value still present: %t, %v", found, err)
			}

			removed, err = set.Remove(t.Context(), values[index])
			if err != nil || removed {
				t.Fatalf("duplicate removal: %t, %v", removed, err)
			}
		}

		if added, err := set.Add(t.Context(), values[0]); err != nil || !added || set.Len() != 1 {
			t.Fatalf("reuse emptied set: %t, %v; size=%d", added, err, set.Len())
		}
	}
}

func TestSetLookupErrorsPreserveMembership(t *testing.T) {
	sentinel := errors.New("comparison failed")
	set := valueset.New(0)
	if _, err := set.Add(t.Context(), collisionValue{label: "first", err: sentinel}); err != nil {
		t.Fatal(err)
	}

	for _, lookup := range []func(context.Context, runtime.Value) (bool, error){set.Contains, set.Remove} {
		found, err := lookup(t.Context(), collisionValue{label: "second"})
		if found || !errors.Is(err, sentinel) || set.Len() != 1 {
			t.Fatalf("failed lookup: %t, %v; size=%d", found, err, set.Len())
		}
	}

	canceled, cancel := context.WithCancel(t.Context())
	cancel()
	if found, err := set.Remove(canceled, collisionValue{label: "first"}); found || !errors.Is(err, context.Canceled) || set.Len() != 1 {
		t.Fatalf("canceled removal: %t, %v; size=%d", found, err, set.Len())
	}
}

func TestSetMembershipUsesNumericAndNestedEquality(t *testing.T) {
	for _, pair := range [][2]runtime.Value{
		{runtime.Int(1), runtime.Float(1)},
		{runtime.NewArrayWith(runtime.Int(1)), runtime.NewArrayWith(runtime.Float(1))},
		{runtime.NewObjectWith(map[string]runtime.Value{"a": runtime.Int(1)}), runtime.NewObjectWith(map[string]runtime.Value{"a": runtime.Float(1)})},
	} {
		set := valueset.New(0)
		if _, err := set.Add(t.Context(), pair[0]); err != nil {
			t.Fatal(err)
		}

		if found, err := set.Contains(t.Context(), pair[1]); err != nil || !found {
			t.Fatalf("equal value absent: %t, %v", found, err)
		}

		if removed, err := set.Remove(t.Context(), pair[1]); err != nil || !removed || set.Len() != 0 {
			t.Fatalf("equal value not removed: %t, %v", removed, err)
		}
	}
}
