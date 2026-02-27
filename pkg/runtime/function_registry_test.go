package runtime

import (
	"context"
	"slices"
	"testing"
)

func TestFunctionsBuilderBuildAndHash(t *testing.T) {
	varFn := func(ctx context.Context, args ...Value) (Value, error) {
		return None, nil
	}
	fn0 := func(ctx context.Context) (Value, error) {
		return None, nil
	}
	fn1 := func(ctx context.Context, arg Value) (Value, error) {
		return arg, nil
	}

	builder := NewFunctionsBuilder()
	builder.Var().Add("var", varFn)
	builder.A0().Add("zero", fn0)
	builder.A1().Add("one", fn1)

	funcs, err := builder.Build()
	if err != nil {
		t.Fatalf("build functions: %v", err)
	}

	if funcs.Size() != 3 {
		t.Fatalf("expected 3 functions, got %d", funcs.Size())
	}

	if funcs.size != funcs.Size() {
		t.Fatalf("expected cached size %d, got %d", funcs.Size(), funcs.size)
	}

	if len(funcs.names) != funcs.Size() {
		t.Fatalf("expected cached names length %d, got %d", funcs.Size(), len(funcs.names))
	}

	for _, name := range []string{"VAR", "ZERO", "ONE"} {
		if !funcs.Has(name) {
			t.Fatalf("expected function %q to exist", name)
		}
	}

	names := funcs.List()
	for _, name := range []string{"VAR", "ZERO", "ONE"} {
		if !slices.Contains(names, name) {
			t.Fatalf("expected names to include %q, got %v", name, names)
		}
	}

	if len(names) > 0 {
		names[0] = "MUTATED"
		if funcs.names[0] == "MUTATED" {
			t.Fatal("expected List to return a copy of cached names")
		}
	}

	if funcs.Hash() == 0 {
		t.Fatal("expected non-zero hash for non-empty functions")
	}

	if funcs.Hash() != functionsHash(funcs) {
		t.Fatalf("expected hash to match functionsHash, got %d vs %d", funcs.Hash(), functionsHash(funcs))
	}
}

func TestNewFunctionsFromAndFromMap(t *testing.T) {
	fn0 := func(ctx context.Context) (Value, error) {
		return None, nil
	}

	f1Builder := NewFunctionsBuilder()
	f1Builder.A0().Add("a", fn0)
	f1, err := f1Builder.Build()
	if err != nil {
		t.Fatalf("build functions: %v", err)
	}

	f2Builder := NewFunctionsBuilder()
	f2Builder.A0().Add("b", fn0)
	f2, err := f2Builder.Build()
	if err != nil {
		t.Fatalf("build functions: %v", err)
	}

	merged, err := NewFunctionsFrom(f1, f2)
	if err != nil {
		t.Fatalf("merge functions: %v", err)
	}

	if merged.Size() != 2 {
		t.Fatalf("expected 2 merged functions, got %d", merged.Size())
	}

	if merged.size != merged.Size() {
		t.Fatalf("expected cached size %d, got %d", merged.Size(), merged.size)
	}

	if len(merged.names) != merged.Size() {
		t.Fatalf("expected cached names length %d, got %d", merged.Size(), len(merged.names))
	}

	for _, name := range []string{"A", "B"} {
		if !merged.Has(name) {
			t.Fatalf("expected merged function %q to exist", name)
		}
	}

	fromMap, err := NewFunctionsFromMap(map[string]Function{
		"FOO": func(ctx context.Context, args ...Value) (Value, error) {
			return None, nil
		},
	})
	if err != nil {
		t.Fatalf("functions from map: %v", err)
	}

	if !fromMap.Has("FOO") {
		t.Fatal("expected functions from map to include FOO")
	}

	if fromMap.size != fromMap.Size() {
		t.Fatalf("expected cached size %d, got %d", fromMap.Size(), fromMap.size)
	}

	if len(fromMap.names) != fromMap.Size() {
		t.Fatalf("expected cached names length %d, got %d", fromMap.Size(), len(fromMap.names))
	}

	if fromMap.Hash() != functionsHash(fromMap) {
		t.Fatalf("expected map hash to match functionsHash, got %d vs %d", fromMap.Hash(), functionsHash(fromMap))
	}
}
