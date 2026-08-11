package vm_test

import (
	"context"
	"errors"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/vm"
	"github.com/MontFerret/ferret/v2/test/spec"
	. "github.com/MontFerret/ferret/v2/test/spec/exec"
)

func TestCoalesceOperator(t *testing.T) {
	RunSpecs(t, []spec.Spec{
		S("RETURN NONE ?? 42", 42),
		S("RETURN NULL ?? 42", 42),
		S("RETURN 10 ?? 42", 10),
		S("RETURN false ?? true", false),
		S("RETURN 0 ?? 100", 0),
		S(`RETURN "" ?? "fallback"`, ""),
		Array("RETURN [] ?? [1, 2, 3]", []any{}),
		Object("RETURN {} ?? { fallback: true }", map[string]any{}),
		S("RETURN NONE ?? NONE", nil),
		S("RETURN NONE ?? NONE ?? 3", 3),
		S("RETURN NONE ?? 0 ?? 3", 0),
		S(`RETURN NONE ?? "unknown"`, "unknown"),
		S(`RETURN {}?.missing ?? "fallback"`, "fallback"),
		S(`RETURN { items: [] }?.items?.[0] ?? "fallback"`, "fallback"),
		S("RETURN 0 ?? 1 OR 2", 0),
		S("RETURN TRUE ? NONE ?? 3 : 4", 3),
		S("RETURN FALSE ? 1 : NONE ?? 4", 4),
		S("RETURN FALSE ? : NONE ?? 4", 4),
	})
}

func TestCoalesceOperatorShortCircuitsFallback(t *testing.T) {
	calls := 0

	RunSpecs(t, []spec.Spec{
		S("RETURN 1 ?? FALLBACK()", 1),
	}, vm.WithFunction("FALLBACK", func(context.Context, ...runtime.Value) (runtime.Value, error) {
		calls++

		return runtime.NewInt(42), nil
	}))

	if calls != 0 {
		t.Fatalf("expected present value to skip fallback, got %d calls", calls)
	}

	RunSpecs(t, []spec.Spec{
		S("RETURN NONE ?? FALLBACK()", 42),
	}, vm.WithFunction("FALLBACK", func(context.Context, ...runtime.Value) (runtime.Value, error) {
		calls++

		return runtime.NewInt(42), nil
	}))

	if calls != 2 {
		t.Fatalf("expected one fallback call at each optimization level, got %d calls", calls)
	}
}

func TestCoalesceOperatorPropagatesErrors(t *testing.T) {
	boom := errors.New("boom")

	RunSpecs(t, []spec.Spec{
		Error("RETURN FAIL() ?? 1"),
		Error("RETURN NONE ?? FAIL()"),
		S("RETURN FAIL()? ?? 1", 1),
	}, vm.WithFunction("FAIL", func(context.Context, ...runtime.Value) (runtime.Value, error) {
		return runtime.None, boom
	}))
}
