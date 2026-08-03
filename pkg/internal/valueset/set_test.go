package valueset_test

import (
	"context"
	"errors"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/internal/valueset"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

type collisionValue struct {
	err   error
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

func (v collisionValue) Equal(_ context.Context, other runtime.Value) (bool, error) {
	if v.err != nil {
		return false, v.err
	}

	o, ok := other.(collisionValue)
	if !ok {
		return false, nil
	}

	return v.label == o.label, nil
}

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
