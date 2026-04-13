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

	for _, name := range []string{"var", "zero", "one"} {
		if !funcs.Has(name) {
			t.Fatalf("expected function %q to exist", name)
		}
	}

	names := funcs.List()
	for _, name := range []string{"var", "zero", "one"} {
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
	f1Builder.A0().Add("A", fn0)
	f1, err := f1Builder.Build()
	if err != nil {
		t.Fatalf("build functions: %v", err)
	}

	f2Builder := NewFunctionsBuilder()
	f2Builder.A0().Add("a", fn0)
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

	for _, name := range []string{"A", "a"} {
		if !merged.Has(name) {
			t.Fatalf("expected merged function %q to exist", name)
		}
	}

	fromMap, err := NewFunctionsFromMap(map[string]Function{
		"Foo": func(ctx context.Context, args ...Value) (Value, error) {
			return None, nil
		},
	})
	if err != nil {
		t.Fatalf("functions from map: %v", err)
	}

	if !fromMap.Has("Foo") {
		t.Fatal("expected functions from map to include Foo")
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

func TestFunctionLookupIsCaseSensitive(t *testing.T) {
	fooUpper := func(context.Context) (Value, error) {
		return NewString("upper"), nil
	}
	fooLower := func(context.Context) (Value, error) {
		return NewString("lower"), nil
	}

	builder := NewFunctionsBuilder()
	builder.A0().Add("Foo", fooUpper)
	builder.A0().Add("foo", fooLower)

	funcs, err := builder.Build()
	if err != nil {
		t.Fatalf("build functions: %v", err)
	}

	if funcs.Size() != 2 {
		t.Fatalf("expected 2 functions, got %d", funcs.Size())
	}

	if !funcs.Has("Foo") || !funcs.Has("foo") {
		t.Fatalf("expected exact-case host functions to exist, got %v", funcs.List())
	}

	if funcs.Has("FOO") {
		t.Fatalf("expected wrong-case host name to be absent, got %v", funcs.List())
	}

	upper, ok := funcs.A0().Get("Foo")
	if !ok {
		t.Fatal("expected Foo lookup to succeed")
	}

	lower, ok := funcs.A0().Get("foo")
	if !ok {
		t.Fatal("expected foo lookup to succeed")
	}

	if _, ok := funcs.A0().Get("FOO"); ok {
		t.Fatal("expected wrong-case lookup to fail")
	}

	if got, _ := upper(context.Background()); got != NewString("upper") {
		t.Fatalf("unexpected Foo result: %v", got)
	}

	if got, _ := lower(context.Background()); got != NewString("lower") {
		t.Fatalf("unexpected foo result: %v", got)
	}
}

func assertEmptyFunctionCollection[T FunctionConstraint](t *testing.T, name string, col FunctionCollection[T]) {
	t.Helper()

	if col.Size() != 0 {
		t.Fatalf("expected %s collection to be empty, got size %d", name, col.Size())
	}

	if col.Has("one") {
		t.Fatalf("expected %s collection to miss registered A1 function", name)
	}

	if _, ok := col.Get("missing"); ok {
		t.Fatalf("expected %s collection lookup to miss", name)
	}
}

func TestEmptyFunctionAccessorsRemainSparse(t *testing.T) {
	builder := NewFunctionsBuilder()
	builder.A1().Add("one", func(_ context.Context, arg Value) (Value, error) {
		return arg, nil
	})

	funcs, err := builder.Build()
	if err != nil {
		t.Fatalf("build functions: %v", err)
	}

	if funcs.a1 == nil {
		t.Fatal("expected populated A1 collection")
	}

	if funcs.av != nil || funcs.a0 != nil || funcs.a2 != nil || funcs.a3 != nil || funcs.a4 != nil {
		t.Fatalf("expected sparse registry, got av=%v a0=%v a2=%v a3=%v a4=%v", funcs.av, funcs.a0, funcs.a2, funcs.a3, funcs.a4)
	}

	assertEmptyFunctionCollection(t, "var", funcs.Var())
	assertEmptyFunctionCollection(t, "a0", funcs.A0())
	assertEmptyFunctionCollection(t, "a2", funcs.A2())
	assertEmptyFunctionCollection(t, "a3", funcs.A3())
	assertEmptyFunctionCollection(t, "a4", funcs.A4())

	if funcs.av != nil || funcs.a0 != nil || funcs.a2 != nil || funcs.a3 != nil || funcs.a4 != nil {
		t.Fatalf("expected empty accessors to preserve sparse registry, got av=%v a0=%v a2=%v a3=%v a4=%v", funcs.av, funcs.a0, funcs.a2, funcs.a3, funcs.a4)
	}
}
