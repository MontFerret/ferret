package runtime

import (
	"context"
	"slices"
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
	}
}

func TestNamespaceNewNamespaceQualifiedNames(t *testing.T) {
	ns := NewNamespace("foo")

	ns.Function().A0().
		Add("bar", func(ctx context.Context) (Value, error) {
			return None, nil
		})

	funcs, err := ns.(*library).Build()
	if err != nil {
		t.Fatalf("build functions: %v", err)
	}

	names := funcs.List()
	if !slices.Contains(names, "foo::bar") {
		t.Fatalf("expected qualified name in namespace, got %v", names)
	}
}

func TestNamespaceAllowsCaseDistinctQualifiedNames(t *testing.T) {
	root := NewLibrary()

	root.Namespace("Foo").Function().A0().
		Add("Bar", func(ctx context.Context) (Value, error) {
			return NewString("upper"), nil
		})

	root.Namespace("foo").Function().A0().
		Add("Bar", func(ctx context.Context) (Value, error) {
			return NewString("lower"), nil
		})

	funcs, err := root.Build()
	if err != nil {
		t.Fatalf("build functions: %v", err)
	}

	names := funcs.List()
	if !slices.Contains(names, "Foo::Bar") || !slices.Contains(names, "foo::Bar") {
		t.Fatalf("expected exact-case qualified names, got %v", names)
	}

	if _, ok := funcs.A0().Get("FOO::BAR"); ok {
		t.Fatalf("expected wrong-case qualified lookup to fail, got %v", names)
	}
}
