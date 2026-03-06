package runtime

import "testing"

func TestParamsGetSetHasDelete(t *testing.T) {
	var params Params

	if value, exists := params.Get("missing"); exists || value != None {
		t.Fatalf("expected missing value to return (None, false), got (%v, %t)", value, exists)
	}

	if params.Has("missing") {
		t.Fatal("expected missing key to not exist")
	}

	fallback := NewString("fallback")
	if value := params.GetOr("missing", fallback); value != fallback {
		t.Fatalf("expected fallback value, got %v", value)
	}

	if value := params.GetOr("missing", nil); value != None {
		t.Fatalf("expected None fallback for nil default, got %v", value)
	}

	params.Set("foo", NewInt(42))
	params.Set("nullable", nil)

	if !params.Has("foo") {
		t.Fatal("expected key foo to exist")
	}

	if value, exists := params.Get("foo"); !exists || value != NewInt(42) {
		t.Fatalf("expected (42, true), got (%v, %t)", value, exists)
	}

	if value, exists := params.Get("nullable"); !exists || value != None {
		t.Fatalf("expected (None, true) for nil assignment, got (%v, %t)", value, exists)
	}

	params.Delete("foo")

	if params.Has("foo") {
		t.Fatal("expected key foo to be deleted")
	}

	if value, exists := params.Get("foo"); exists || value != None {
		t.Fatalf("expected deleted key to return (None, false), got (%v, %t)", value, exists)
	}
}

func TestParamsMustGet(t *testing.T) {
	params := Params{}
	params.Set("foo", NewString("bar"))

	if value := params.MustGet("foo"); value != NewString("bar") {
		t.Fatalf("expected value bar, got %v", value)
	}

	defer func() {
		if recover() == nil {
			t.Fatal("expected panic for missing key")
		}
	}()

	_ = params.MustGet("missing")
}

func TestParamsClone(t *testing.T) {
	original := Params{}
	original.Set("foo", NewInt(1))
	original.Set("bar", nil)

	cloned := original.Clone()

	if value, exists := cloned.Get("foo"); !exists || value != NewInt(1) {
		t.Fatalf("expected clone to contain foo=1, got (%v, %t)", value, exists)
	}

	if value, exists := cloned.Get("bar"); !exists || value != None {
		t.Fatalf("expected clone to normalize nil value to None, got (%v, %t)", value, exists)
	}

	cloned.Set("foo", NewInt(2))
	cloned.Set("baz", NewInt(3))

	if value := original.MustGet("foo"); value != NewInt(1) {
		t.Fatalf("expected original foo to remain 1, got %v", value)
	}

	if original.Has("baz") {
		t.Fatal("expected clone mutations to not affect original map keys")
	}
}

func TestParamsMerge(t *testing.T) {
	left := Params{}
	left.Set("foo", NewInt(1))
	left.Set("bar", NewInt(2))

	right := Params{}
	right.Set("bar", NewInt(20))
	right.Set("baz", nil)

	left.Merge(&right)

	if value := left.MustGet("foo"); value != NewInt(1) {
		t.Fatalf("expected foo to stay 1, got %v", value)
	}

	if value := left.MustGet("bar"); value != NewInt(20) {
		t.Fatalf("expected bar to be overridden to 20, got %v", value)
	}

	if value, exists := left.Get("baz"); !exists || value != None {
		t.Fatalf("expected baz to merge as None, got (%v, %t)", value, exists)
	}

	var zero Params
	zero.Merge(&right)

	if value := zero.MustGet("bar"); value != NewInt(20) {
		t.Fatalf("expected merge into zero value params to set bar=20, got %v", value)
	}
}
