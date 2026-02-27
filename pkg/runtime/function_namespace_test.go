package runtime

import (
	"context"
	"errors"
	"slices"
	"testing"
)

func TestNamespaceRegisterFunctionsNested(t *testing.T) {
	root := NewRootNamespace()
	nested := root.Namespace("foo").Namespace("bar")

	funcs := NewFunctionsBuilder().
		Set0("baz", func(ctx context.Context) (Value, error) {
			return None, nil
		}).
		Build()

	if err := nested.RegisterFunctions(funcs); err != nil {
		t.Fatalf("register nested functions: %v", err)
	}

	names := root.Function().Build().Names()
	if !slices.Contains(names, "FOO::BAR::BAZ") {
		t.Fatalf("expected fully qualified name in root, got %v", names)
	}
}

func TestNamespaceRegisterFunctionsDuplicate(t *testing.T) {
	ns := NewRootNamespace().Namespace("foo")

	funcs := NewFunctionsBuilder().
		Set0("bar", func(ctx context.Context) (Value, error) {
			return None, nil
		}).
		Build()

	if err := ns.RegisterFunctions(funcs); err != nil {
		t.Fatalf("register functions: %v", err)
	}

	if err := ns.RegisterFunctions(funcs); err == nil {
		t.Fatal("expected duplicate registration error")
	} else if !errors.Is(err, ErrNotUnique) {
		t.Fatalf("expected ErrNotUnique, got %v", err)
	}
}

func TestNamespaceNewNamespaceQualifiedNames(t *testing.T) {
	ns := NewNamespace("foo")

	funcs := NewFunctionsBuilder().
		Set0("bar", func(ctx context.Context) (Value, error) {
			return None, nil
		}).
		Build()

	if err := ns.RegisterFunctions(funcs); err != nil {
		t.Fatalf("register functions: %v", err)
	}

	names := ns.Function().Build().Names()
	if !slices.Contains(names, "FOO::BAR") {
		t.Fatalf("expected qualified name in namespace, got %v", names)
	}
}
