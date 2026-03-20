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
		{types.None, "None"},
		{types.Boolean, "Boolean"},
		{types.Int, "Int"},
		{types.Float, "Float"},
		{types.String, "String"},
		{types.DateTime, "DateTime"},
		{types.Array, "Array"},
		{types.Object, "Object"},
		{types.Binary, "Binary"},
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
		{types.Int, types.Int, "same type", 0},
		{types.None, types.Boolean, "None < Boolean", -1},
		{types.Boolean, types.Int, "Boolean < Int", -1},
		{types.Int, types.Float, "Int < Float", -1},
		{types.Float, types.String, "Float < String", -1},
		{types.String, types.DateTime, "String < DateTime", -1},
		{types.DateTime, types.Binary, "DateTime < Binary", -1},
		{types.Binary, types.Array, "Binary < Array", -1},
		{types.Array, types.Object, "Array < Object", -1},
		{types.Object, types.None, "Object > None", 1},
		{types.Float, types.Boolean, "Float > Boolean", 1},
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
