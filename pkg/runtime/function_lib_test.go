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

func TestNamespacePreservesNamespaceCaseAndCanonicalizesFunctionNames(t *testing.T) {
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
	if !slices.Contains(names, "Foo::bar") || !slices.Contains(names, "foo::bar") {
		t.Fatalf("expected exact-case namespaces and canonical function names, got %v", names)
	}

	if _, ok := funcs.A0().Get("Foo::BAR"); !ok {
		t.Fatalf("expected mixed-case terminal lookup to succeed, got %v", names)
	}

	if _, ok := funcs.A0().Get("FOO::bar"); ok {
		t.Fatalf("expected wrong-case namespace lookup to fail, got %v", names)
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
