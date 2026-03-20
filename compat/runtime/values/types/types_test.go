package types_test

import (
	"testing"

	"github.com/MontFerret/ferret/v2/compat/runtime/core"
	"github.com/MontFerret/ferret/v2/compat/runtime/values/types"
)

func TestTypeConstants_NotNil(t *testing.T) {
	all := []struct {
		typ  core.Type
		name string
	}{
		{"None", types.None},
		{"Boolean", types.Boolean},
		{"Int", types.Int},
		{"Float", types.Float},
		{"String", types.String},
		{"DateTime", types.DateTime},
		{"Array", types.Array},
		{"Object", types.Object},
		{"Binary", types.Binary},
	}

	for _, tc := range all {
		if tc.typ == nil {
			t.Fatalf("%s type constant is nil", tc.name)
		}

		if tc.typ.String() != tc.name {
			t.Fatalf("expected %q, got %q", tc.name, tc.typ.String())
		}
	}
}

func TestTypeConstants_Equals(t *testing.T) {
	if !types.Int.Equals(types.Int) {
		t.Fatal("Int should equal itself")
	}

	if types.Int.Equals(types.Float) {
		t.Fatal("Int should not equal Float")
	}

	if types.None.Equals(nil) {
		t.Fatal("None should not equal nil")
	}
}

func TestTypeConstants_ID(t *testing.T) {
	// IDs should be non-zero and consistent.
	id1 := types.Int.ID()
	id2 := types.Int.ID()

	if id1 != id2 {
		t.Fatalf("expected consistent ID, got %d and %d", id1, id2)
	}

	if id1 == 0 {
		t.Fatal("expected non-zero ID for Int")
	}

	// Different types should have different IDs.
	if types.Int.ID() == types.Float.ID() {
		t.Fatal("Int and Float should have different IDs")
	}
}

func TestCompare(t *testing.T) {
	tests := []struct {
		a        core.Type
		b        core.Type
		name     string
		expected int64
	}{
		{"same type", types.Int, types.Int, 0},
		{"None < Boolean", types.None, types.Boolean, -1},
		{"Boolean < Int", types.Boolean, types.Int, -1},
		{"Int < Float", types.Int, types.Float, -1},
		{"Float < String", types.Float, types.String, -1},
		{"String < DateTime", types.String, types.DateTime, -1},
		{"DateTime < Binary", types.DateTime, types.Binary, -1},
		{"Binary < Array", types.Binary, types.Array, -1},
		{"Array < Object", types.Array, types.Object, -1},
		{"Object > None", types.Object, types.None, 1},
		{"Float > Boolean", types.Float, types.Boolean, 1},
	}

	for _, tc := range tests {
		result := types.Compare(tc.a, tc.b)
		if result != tc.expected {
			t.Fatalf("%s: expected %d, got %d", tc.name, tc.expected, result)
		}
	}
}

func TestCompare_nil(t *testing.T) {
	// nil should be treated as rank 0 (same as None)
	if types.Compare(nil, types.None) != 0 {
		t.Fatal("expected nil to equal None rank")
	}

	if types.Compare(nil, types.Int) != -1 {
		t.Fatal("expected nil < Int")
	}
}
