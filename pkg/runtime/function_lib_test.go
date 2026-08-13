package runtime

import (
	"context"
	"slices"
	"strings"
	"testing"
)

func TestNamespaceRegisterFunctionsNested(t *testing.T) {
	root := NewLibrary()
	nested := root.Namespace("foo").Namespace("bar")

	nested.Function().A0().
		Add("baz", func(ctx context.Context) (Value, error) {
			return None, nil
		})

	funcs, err := root.Build()
	if err != nil {
		t.Fatalf("build functions: %v", err)
	}

	names := funcs.List()
	if !slices.Contains(names, "foo::bar::baz") {
		t.Fatalf("expected fully qualified name in root, got %v", names)
	}
}

func TestNamespaceRegisterFunctionsDuplicate(t *testing.T) {
	root := NewLibrary()
	ns := root.Namespace("foo")

	ns.Function().A0().
		Add("bar", func(ctx context.Context) (Value, error) {
			return None, nil
		}).
		Add("bar", func(ctx context.Context) (Value, error) {
			return None, nil
		})

	if _, err := root.Build(); err == nil {
		t.Fatal("expected duplicate registration error")
	} else if got := strings.Count(err.Error(), "already exists"); got != 1 {
		t.Fatalf("expected shared builder error to be reported once, got %d: %v", got, err)
	}
}

func TestNamespaceNewNamespaceQualifiedNames(t *testing.T) {
	ns := NewNamespace("Foo::BAR")

	ns.Function().A0().
		Add("BAZ", func(ctx context.Context) (Value, error) {
			return None, nil
		})

	funcs, err := ns.(*library).Build()
	if err != nil {
		t.Fatalf("build functions: %v", err)
	}

	names := funcs.List()
	if !slices.Contains(names, "Foo::BAR::BAZ") {
		t.Fatalf("expected qualified name in namespace, got %v", names)
	}
}

func TestNamespaceCasingVariantsMergeAndPreserveFirstDeclaredSpelling(t *testing.T) {
	root := NewLibrary()

	root.Namespace("DB").Namespace("Postgres").Function().A0().
		Add("Query", func(ctx context.Context) (Value, error) {
			return NewString("upper"), nil
		})

	root.Namespace("db").Namespace("POSTGRES").Function().A0().
		Add("Health", func(ctx context.Context) (Value, error) {
			return NewString("lower"), nil
		})

	funcs, err := root.Build()
	if err != nil {
		t.Fatalf("build functions: %v", err)
	}

	names := funcs.List()
	if !slices.Equal(names, []string{"DB::Postgres::Health", "DB::Postgres::Query"}) {
		t.Fatalf("expected first-declared namespace casing, got %v", names)
	}

	for _, name := range []string{
		"db::postgres::query",
		"DB::POSTGRES::QUERY",
		"Db::Postgres::Query",
	} {
		fn, ok := funcs.A0().Get(name)
		if !ok {
			t.Fatalf("expected %q lookup to succeed, got %v", name, names)
		}

		value, callErr := fn(t.Context())
		if callErr != nil || value != NewString("upper") {
			t.Fatalf("call %q = %v, %v", name, value, callErr)
		}
	}

	if !funcs.Has("dB::pOsTgReS::hEaLtH") {
		t.Fatalf("expected mixed-case nested namespace lookup to succeed, got %v", names)
	}
}

func TestNamespaceCasingVariantsRejectDuplicateNormalizedIdentity(t *testing.T) {
	root := NewLibrary()
	fn := func(context.Context) (Value, error) { return None, nil }

	root.Namespace("Foo").Namespace("Bar").Function().A0().Add("Baz", fn)
	root.Namespace("fOO").Namespace("bAR").Function().A0().Add("bAZ", fn)

	if _, err := root.Build(); err == nil {
		t.Fatal("expected duplicate normalized namespace member to fail")
	}
}

func TestNamespaceRemovalUsesCaseInsensitiveQualifiedName(t *testing.T) {
	root := NewLibrary()
	root.Namespace("DB").Namespace("Postgres").Function().A0().Add("Query", func(context.Context) (Value, error) {
		return None, nil
	})

	root.Namespace("db").Namespace("POSTGRES").Function().A0().Remove("qUeRy")

	functions, err := root.Build()
	if err != nil {
		t.Fatalf("build functions: %v", err)
	}

	if functions.Size() != 0 {
		t.Fatalf("expected mixed-case namespace removal to remove member, got %v", functions.List())
	}
}

func TestNamespaceCasingDoesNotAffectRegistryHash(t *testing.T) {
	fn := func(context.Context) (Value, error) { return None, nil }
	build := func(namespace, nested, name string) *Functions {
		root := NewLibrary()
		root.Namespace(namespace).Namespace(nested).Function().A0().Add(name, fn)

		functions, err := root.Build()
		if err != nil {
			t.Fatalf("build functions: %v", err)
		}

		return functions
	}

	upper := build("DB", "POSTGRES", "QUERY")
	mixed := build("Db", "Postgres", "Query")
	if upper.Hash() != mixed.Hash() {
		t.Fatalf("equivalent namespace casing produced different hashes: %d != %d", upper.Hash(), mixed.Hash())
	}

	if !slices.Equal(upper.List(), []string{"DB::POSTGRES::QUERY"}) {
		t.Fatalf("uppercase declaration was not preserved: %v", upper.List())
	}

	if !slices.Equal(mixed.List(), []string{"Db::Postgres::Query"}) {
		t.Fatalf("mixed-case declaration was not preserved: %v", mixed.List())
	}
}

func TestNamespaceAllowsFunctionOverloading(t *testing.T) {
	root := NewLibrary()

	root.Namespace("foo").Function().A0().
		Add("bar", func(ctx context.Context) (Value, error) {
			return NewString("zero"), nil
		})

	root.Namespace("foo").Function().A1().
		Add("bar", func(ctx context.Context, arg Value) (Value, error) {
			return NewString("one"), nil
		})

	root.Namespace("foo").Function().Var().
		Add("bar", func(ctx context.Context, args ...Value) (Value, error) {
			return NewString("var"), nil
		})

	funcs, err := root.Build()
	if err != nil {
		t.Fatalf("build functions: %v", err)
	}

	names := funcs.List()
	if !slices.Contains(names, "foo::bar") {
		t.Fatalf("expected qualified name in namespace, got %v", names)
	}

	if _, ok := funcs.A0().Get("foo::bar"); !ok {
		t.Fatalf("expected A0 lookup to succeed")
	}

	if _, ok := funcs.A1().Get("foo::bar"); !ok {
		t.Fatalf("expected A1 lookup to succeed")
	}

	if _, ok := funcs.Var().Get("foo::bar"); !ok {
		t.Fatalf("expected Var lookup to succeed")
	}
}

func TestNamespaceDisallowsDuplicateFunctionOverloading(t *testing.T) {
	root := NewLibrary()

	root.Namespace("foo").Function().A0().
		Add("bar", func(ctx context.Context) (Value, error) {
			return NewString("zero"), nil
		})

	root.Namespace("foo").Function().A1().
		Add("bar", func(ctx context.Context, arg Value) (Value, error) {
			return NewString("one"), nil
		})

	root.Namespace("foo").Function().A1().
		Add("bar", func(ctx context.Context, arg Value) (Value, error) {
			return NewString("one"), nil
		})

	if _, err := root.Build(); err == nil {
		t.Fatal("expected duplicate registration error")
	}
}
