package runtime

import (
	"context"
	"slices"
	"testing"
)

func TestNamespaceRegisterFunctionsNested(t *testing.T) {
	root := NewRootNamespace()
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
	if !slices.Contains(names, "FOO::BAR::BAZ") {
		t.Fatalf("expected fully qualified name in root, got %v", names)
	}
}

func TestNamespaceRegisterFunctionsDuplicate(t *testing.T) {
	root := NewRootNamespace()
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

	funcs, err := ns.(*defaultNamespace).Build()
	if err != nil {
		t.Fatalf("build functions: %v", err)
	}

	names := funcs.List()
	if !slices.Contains(names, "FOO::BAR") {
		t.Fatalf("expected qualified name in namespace, got %v", names)
	}
}
