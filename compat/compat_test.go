package compat_test

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/MontFerret/ferret/v2/compat"
	compatruntime "github.com/MontFerret/ferret/v2/compat/runtime"
	"github.com/MontFerret/ferret/v2/compat/runtime/core"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

func TestNew(t *testing.T) {
	inst := compat.New()
	if inst == nil {
		t.Fatal("expected non-nil Instance")
	}
}

func TestExec_simple(t *testing.T) {
	inst := compat.New()
	ctx := context.Background()

	out, err := inst.Exec(ctx, `RETURN "hello"`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var result string
	if err := json.Unmarshal(out, &result); err != nil {
		t.Fatalf("failed to unmarshal output %q: %v", out, err)
	}

	if result != "hello" {
		t.Fatalf("expected \"hello\", got %q", result)
	}
}

func TestExec_withParam(t *testing.T) {
	inst := compat.New()
	ctx := context.Background()

	out, err := inst.Exec(ctx, `RETURN @name`, compatruntime.WithParam("name", "ferret"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var result string
	if err := json.Unmarshal(out, &result); err != nil {
		t.Fatalf("failed to unmarshal output %q: %v", out, err)
	}

	if result != "ferret" {
		t.Fatalf("expected \"ferret\", got %q", result)
	}
}

func TestExec_withParams(t *testing.T) {
	inst := compat.New()
	ctx := context.Background()

	out, err := inst.Exec(ctx, `RETURN @x + @y`,
		compatruntime.WithParams(map[string]interface{}{"x": 3, "y": 4}))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var result int
	if err := json.Unmarshal(out, &result); err != nil {
		t.Fatalf("failed to unmarshal output %q: %v", out, err)
	}

	if result != 7 {
		t.Fatalf("expected 7, got %d", result)
	}
}

func TestCompile_andRun(t *testing.T) {
	inst := compat.New()
	ctx := context.Background()

	prog := inst.MustCompile(`RETURN 42`)

	if prog.Source() != `RETURN 42` {
		t.Fatalf("unexpected source: %q", prog.Source())
	}

	out, err := prog.Run(ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var result int
	if err := json.Unmarshal(out, &result); err != nil {
		t.Fatalf("failed to unmarshal output %q: %v", out, err)
	}

	if result != 42 {
		t.Fatalf("expected 42, got %d", result)
	}
}

func TestCompile_params(t *testing.T) {
	inst := compat.New()

	prog := inst.MustCompile(`RETURN @value`)
	params := prog.Params()

	if len(params) != 1 || params[0] != "value" {
		t.Fatalf("expected params=[\"value\"], got %v", params)
	}
}

func TestInstance_Run(t *testing.T) {
	inst := compat.New()
	ctx := context.Background()

	prog := inst.MustCompile(`RETURN "via Run"`)
	out, err := inst.Run(ctx, prog)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var result string
	if err := json.Unmarshal(out, &result); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	if result != "via Run" {
		t.Fatalf("expected \"via Run\", got %q", result)
	}
}

func TestWithoutStdlib(t *testing.T) {
	inst := compat.New(compat.WithoutStdlib())
	ctx := context.Background()

	// A stdlib function (e.g. CONCAT) should not be available.
	_, err := inst.Exec(ctx, `RETURN CONCAT("a", "b")`)
	if err == nil {
		t.Fatal("expected error when stdlib is disabled, got nil")
	}
}

func TestFunctions_RegisterViaNamespace(t *testing.T) {
	inst := compat.New()

	// Register a function through the compat Namespace API.
	// Note: this does NOT affect the already-built engine — it's for API compatibility only.
	ns := inst.Functions()
	err := ns.RegisterFunction("MY_FUNC", func(ctx context.Context, args ...core.Value) (core.Value, error) {
		return core.WrapValue(runtime.NewString("ok")), nil
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Double registration should return an error.
	err = ns.RegisterFunction("MY_FUNC", func(ctx context.Context, args ...core.Value) (core.Value, error) {
		return core.WrapValue(runtime.NewString("ok")), nil
	})
	if err == nil {
		t.Fatal("expected error on duplicate registration, got nil")
	}
}

func TestFunctions_RegisteredFunctions(t *testing.T) {
	inst := compat.New()
	ns := inst.Functions()

	_ = ns.RegisterFunction("CUSTOM_A", func(ctx context.Context, args ...core.Value) (core.Value, error) {
		return core.WrapValue(runtime.None), nil
	})
	_ = ns.RegisterFunction("CUSTOM_B", func(ctx context.Context, args ...core.Value) (core.Value, error) {
		return core.WrapValue(runtime.None), nil
	})

	names := ns.RegisteredFunctions()
	found := map[string]bool{}
	for _, n := range names {
		found[n] = true
	}

	if !found["CUSTOM_A"] || !found["CUSTOM_B"] {
		t.Fatalf("expected CUSTOM_A and CUSTOM_B in registered functions, got %v", names)
	}
}

func TestFunctions_RemoveFunction(t *testing.T) {
	inst := compat.New()
	ns := inst.Functions()

	_ = ns.RegisterFunction("TO_REMOVE", func(ctx context.Context, args ...core.Value) (core.Value, error) {
		return core.WrapValue(runtime.None), nil
	})

	ns.RemoveFunction("TO_REMOVE")

	names := ns.RegisteredFunctions()
	for _, n := range names {
		if n == "TO_REMOVE" {
			t.Fatal("function should have been removed")
		}
	}
}
