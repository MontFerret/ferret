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
	builder.A0().Add("overloaded", fn0)
	builder.A1().Add("overloaded", fn1)

	funcs, err := builder.Build()
	if err != nil {
		t.Fatalf("build functions: %v", err)
	}

	// Size counts registered definitions, including overloads.
	if funcs.Size() != 5 {
		t.Fatalf("expected 5 function definitions, got %d", funcs.Size())
	}

	if funcs.size != funcs.Size() {
		t.Fatalf("expected cached size %d, got %d", funcs.Size(), funcs.size)
	}

	// Names contains unique logical function names.
	if len(funcs.names) != 4 {
		t.Fatalf("expected 4 unique function names, got %d", len(funcs.names))
	}

	for _, name := range []string{"var", "zero", "one", "overloaded"} {
		if !funcs.Has(name) {
			t.Fatalf("expected function %q to exist", name)
		}
	}

	if !funcs.A0().Has("overloaded") {
		t.Fatal("expected overloaded/0 to exist")
	}

	if !funcs.A1().Has("overloaded") {
		t.Fatal("expected overloaded/1 to exist")
	}

	names := funcs.List()
	if !slices.IsSorted(names) {
		t.Fatalf("expected sorted names, got %v", names)
	}

	for _, name := range []string{"var", "zero", "one", "overloaded"} {
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

	// Equivalent registries must hash identically regardless of insertion order.
	equivalentBuilder := NewFunctionsBuilder()
	equivalentBuilder.A1().Add("overloaded", fn1)
	equivalentBuilder.A0().Add("overloaded", fn0)
	equivalentBuilder.A1().Add("one", fn1)
	equivalentBuilder.A0().Add("zero", fn0)
	equivalentBuilder.Var().Add("var", varFn)

	equivalent, err := equivalentBuilder.Build()
	if err != nil {
		t.Fatalf("build equivalent functions: %v", err)
	}

	if funcs.Hash() != equivalent.Hash() {
		t.Fatalf(
			"expected equivalent registries to have the same hash, got %d and %d",
			funcs.Hash(),
			equivalent.Hash(),
		)
	}

	// The same logical name registered at a different arity must change the hash.
	differentArityBuilder := NewFunctionsBuilder()
	differentArityBuilder.Var().Add("var", varFn)
	differentArityBuilder.A0().Add("zero", fn0)
	differentArityBuilder.A1().Add("one", fn1)
	differentArityBuilder.A0().Add("overloaded", fn0)

	differentArity, err := differentArityBuilder.Build()
	if err != nil {
		t.Fatalf("build functions with different arity set: %v", err)
	}

	if funcs.Hash() == differentArity.Hash() {
		t.Fatal("expected different arity sets to produce different hashes")
	}

	variadicOnlyBuilder := NewFunctionsBuilder()
	variadicOnlyBuilder.Var().Add("same", varFn)
	variadicOnly, err := variadicOnlyBuilder.Build()
	if err != nil {
		t.Fatalf("build variadic functions: %v", err)
	}

	fixedOnlyBuilder := NewFunctionsBuilder()
	fixedOnlyBuilder.A0().Add("same", fn0)
	fixedOnly, err := fixedOnlyBuilder.Build()
	if err != nil {
		t.Fatalf("build fixed functions: %v", err)
	}

	if variadicOnly.Hash() == fixedOnly.Hash() {
		t.Fatal("expected the variadic marker to differ from fixed arity zero")
	}
}

func TestNewFunctionsFromAndFromMap(t *testing.T) {
	fn0 := func(ctx context.Context) (Value, error) {
		return None, nil
	}
	fn1 := func(ctx context.Context, arg Value) (Value, error) {
		return arg, nil
	}

	f1Builder := NewFunctionsBuilder()
	f1Builder.A0().Add("A", fn0)
	f1Builder.A0().Add("overloaded", fn0)

	f1, err := f1Builder.Build()
	if err != nil {
		t.Fatalf("build functions: %v", err)
	}

	f2Builder := NewFunctionsBuilder()
	f2Builder.A0().Add("B", fn0)
	f2Builder.A1().Add("overloaded", fn1)

	f2, err := f2Builder.Build()
	if err != nil {
		t.Fatalf("build functions: %v", err)
	}

	merged, err := NewFunctionsFrom(f1, f2)
	if err != nil {
		t.Fatalf("merge functions: %v", err)
	}

	// a/0, b/0, overloaded/0, overloaded/1.
	if merged.Size() != 4 {
		t.Fatalf("expected 4 merged definitions, got %d", merged.Size())
	}

	if merged.size != merged.Size() {
		t.Fatalf("expected cached size %d, got %d", merged.Size(), merged.size)
	}

	// a, b, overloaded.
	if len(merged.names) != 3 {
		t.Fatalf(
			"expected 3 unique cached names, got %d: %v",
			len(merged.names),
			merged.names,
		)
	}

	for _, name := range []string{"a", "B", "overloaded"} {
		if !merged.Has(name) {
			t.Fatalf("expected merged function %q to exist", name)
		}
	}

	if !merged.A0().Has("overloaded") {
		t.Fatal("expected merged overloaded/0 function")
	}

	if !merged.A1().Has("overloaded") {
		t.Fatal("expected merged overloaded/1 function")
	}

	fromMap, err := NewFunctionsFromMap(map[string]Function{
		"Foo": func(ctx context.Context, args ...Value) (Value, error) {
			return None, nil
		},
	})
	if err != nil {
		t.Fatalf("functions from map: %v", err)
	}

	if !fromMap.Has("FOO") {
		t.Fatal("expected functions from map to include canonical foo identity")
	}

	if got := fromMap.List(); !slices.Equal(got, []string{"foo"}) {
		t.Fatalf("expected canonical functions from map, got %v", got)
	}

	if fromMap.size != fromMap.Size() {
		t.Fatalf(
			"expected cached size %d, got %d",
			fromMap.Size(),
			fromMap.size,
		)
	}

	if len(fromMap.names) != 1 {
		t.Fatalf(
			"expected one unique cached name, got %d: %v",
			len(fromMap.names),
			fromMap.names,
		)
	}

	if fromMap.Hash() == 0 {
		t.Fatal("expected non-zero hash for functions created from map")
	}

	equivalent, err := NewFunctionsFromMap(map[string]Function{
		"fOO": func(ctx context.Context, args ...Value) (Value, error) {
			return None, nil
		},
	})
	if err != nil {
		t.Fatalf("build equivalent functions from map: %v", err)
	}

	if fromMap.Hash() != equivalent.Hash() {
		t.Fatalf(
			"expected equivalent map registries to have the same hash, got %d and %d",
			fromMap.Hash(),
			equivalent.Hash(),
		)
	}
}

func TestNewFunctionsBuilderFromPreservesOverloadsAndRejectsSameSignature(t *testing.T) {
	fn0 := func(context.Context) (Value, error) { return None, nil }
	fn1 := func(_ context.Context, arg Value) (Value, error) { return arg, nil }
	varFn := func(context.Context, ...Value) (Value, error) { return None, nil }

	sourceBuilder := NewFunctionsBuilder()
	sourceBuilder.A0().Add("OVERLOAD", fn0)
	sourceBuilder.A1().Add("OVERLOAD", fn1)
	sourceBuilder.Var().Add("OVERLOAD", varFn)
	source, err := sourceBuilder.Build()
	if err != nil {
		t.Fatalf("build source functions: %v", err)
	}

	cloned, err := NewFunctionsBuilderFrom(source).Build()
	if err != nil {
		t.Fatalf("clone functions: %v", err)
	}
	if cloned.Size() != 3 || len(cloned.List()) != 1 {
		t.Fatalf("expected three definitions for one logical name, got size=%d list=%v", cloned.Size(), cloned.List())
	}
	if !cloned.A0().Has("OVERLOAD") || !cloned.A1().Has("OVERLOAD") || !cloned.Var().Has("OVERLOAD") {
		t.Fatal("expected all overloads to survive cloning")
	}

	duplicateBuilder := NewFunctionsBuilder()
	duplicateBuilder.A0().Add("overLOAD", fn0)
	duplicate, err := duplicateBuilder.Build()
	if err != nil {
		t.Fatalf("build duplicate functions: %v", err)
	}

	if _, err := NewFunctionsFrom(source, duplicate); err == nil {
		t.Fatal("expected merging the same name and arity to fail")
	}
}

func TestFunctionLookupUsesCanonicalQualifiedName(t *testing.T) {
	foo := func(context.Context) (Value, error) {
		return NewString("foo"), nil
	}

	builder := NewFunctionsBuilder()
	builder.A0().Add("DB::Postgres::Foo", foo)

	funcs, err := builder.Build()
	if err != nil {
		t.Fatalf("build functions: %v", err)
	}

	if funcs.Size() != 1 {
		t.Fatalf("expected 1 function, got %d", funcs.Size())
	}

	if got := funcs.List(); !slices.Equal(got, []string{"db::postgres::foo"}) {
		t.Fatalf("expected canonical function metadata, got %v", got)
	}

	if got := funcs.A0().Names(); !slices.Equal(got, []string{"db::postgres::foo"}) {
		t.Fatalf("expected canonical collection names, got %v", got)
	}

	var enumerated []string
	funcs.A0().ForEach(func(_ Function0, name string) bool {
		enumerated = append(enumerated, name)

		return true
	})
	if !slices.Equal(enumerated, []string{"db::postgres::foo"}) {
		t.Fatalf("expected canonical ForEach name, got %v", enumerated)
	}

	for _, name := range []string{"db::postgres::foo", "DB::POSTGRES::FOO", "Db::Postgres::FoO"} {
		if !funcs.Has(name) {
			t.Fatalf("expected %q to resolve canonical function, got %v", name, funcs.List())
		}

		resolved, ok := funcs.A0().Get(name)
		if !ok {
			t.Fatalf("expected %q lookup to succeed", name)
		}

		if got, _ := resolved(context.Background()); got != NewString("foo") {
			t.Fatalf("unexpected %q result: %v", name, got)
		}
	}

	duplicate := NewFunctionsBuilder()
	duplicate.A0().Add("db::postgres::foo", foo)
	duplicate.A0().Add("DB::POSTGRES::FOO", foo)

	if _, err := duplicate.Build(); err == nil {
		t.Fatal("expected case-only duplicate registration to fail")
	}

	removable := NewFunctionsBuilder()
	removable.A0().Add("DB::Postgres::Foo", foo)
	removable.A0().Remove("db::POSTGRES::fOO")

	removed, err := removable.Build()
	if err != nil {
		t.Fatalf("build functions after mixed-case removal: %v", err)
	}

	if removed.Size() != 0 {
		t.Fatalf("expected mixed-case removal to remove canonical function, got %v", removed.List())
	}
}

func TestStandaloneFunctionCollectionPreservesExactMapKeys(t *testing.T) {
	fn := func(context.Context) (Value, error) { return None, nil }
	functions := NewFunctionCollectionFromMap(map[string]Function0{"Foo": fn})

	if !functions.Has("Foo") {
		t.Fatal("expected exact map key to resolve")
	}

	if functions.Has("foo") || functions.Has("FOO") {
		t.Fatal("standalone map-backed collection must preserve exact lookup behavior")
	}

	if got := functions.Names(); !slices.Equal(got, []string{"Foo"}) {
		t.Fatalf("standalone collection names = %v, want exact map key", got)
	}
}

func TestFunctionRegistrationRejectsEmptyTerminalName(t *testing.T) {
	fn := func(context.Context) (Value, error) { return None, nil }
	builder := NewFunctionsBuilder()
	builder.A0().Add("", fn)

	if _, err := builder.Build(); err == nil {
		t.Fatal("expected empty root function name to fail")
	}

	library := NewLibrary()
	library.Namespace("TEST").Function().A0().Add("", fn)

	if _, err := library.Build(); err == nil {
		t.Fatal("expected empty namespaced function name to fail")
	}
}

func TestNewFunctionsFromMapRejectsCaseOnlyDuplicate(t *testing.T) {
	fn := func(context.Context, ...Value) (Value, error) { return None, nil }

	if _, err := NewFunctionsFromMap(map[string]Function{"Foo": fn, "foo": fn}); err == nil {
		t.Fatal("expected case-only duplicate map registrations to fail")
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
