package compiler_test

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/MontFerret/ferret/v2/compat/compiler"
	"github.com/MontFerret/ferret/v2/compat/runtime/core"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

func TestCompiler_Compile(t *testing.T) {
	c := compiler.New()

	prog, err := c.Compile(`RETURN "ok"`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if prog == nil {
		t.Fatal("expected non-nil Program")
	}
}

func TestCompiler_MustCompile_panic(t *testing.T) {
	c := compiler.New()

	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic on invalid query, got nil")
		}
	}()

	c.MustCompile(`THIS IS NOT VALID FQL !!!`)
}

func TestCompiler_RegisterFunction(t *testing.T) {
	c := compiler.New()

	err := c.RegisterFunction("HELLO", func(ctx context.Context, args ...core.Value) (core.Value, error) {
		return core.WrapValue(runtime.NewString("hello")), nil
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	prog, err := c.Compile(`RETURN HELLO()`)
	if err != nil {
		t.Fatalf("compile error: %v", err)
	}

	out, err := prog.Run(context.Background())
	if err != nil {
		t.Fatalf("run error: %v", err)
	}

	var result string
	if err := json.Unmarshal(out, &result); err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}

	if result != "hello" {
		t.Fatalf("expected \"hello\", got %q", result)
	}
}

func TestCompiler_RegisteredFunctions(t *testing.T) {
	c := compiler.New()

	_ = c.RegisterFunction("FUNC_A", func(ctx context.Context, args ...core.Value) (core.Value, error) {
		return core.WrapValue(runtime.None), nil
	})
	_ = c.RegisterFunction("FUNC_B", func(ctx context.Context, args ...core.Value) (core.Value, error) {
		return core.WrapValue(runtime.None), nil
	})

	names := c.RegisteredFunctions()
	found := map[string]bool{}
	for _, n := range names {
		found[n] = true
	}

	if !found["FUNC_A"] || !found["FUNC_B"] {
		t.Fatalf("expected FUNC_A and FUNC_B, got %v", names)
	}
}

func TestCompiler_RegisterFunctions_duplicate(t *testing.T) {
	c := compiler.New()

	fns := core.NewFunctionsFromMap(map[string]core.Function{
		"DUP_FUNC": func(ctx context.Context, args ...core.Value) (core.Value, error) {
			return core.WrapValue(runtime.NewString("dup")), nil
		},
	})

	// First registration must succeed.
	if err := c.RegisterFunctions(fns); err != nil {
		t.Fatalf("first RegisterFunctions error: %v", err)
	}

	// Second registration of the same set should be silently skipped (no error).
	if err := c.RegisterFunctions(fns); err != nil {
		t.Fatalf("second RegisterFunctions should silently skip duplicates, got: %v", err)
	}

	// Most importantly: Compile must succeed — no latent builder error must have
	// been recorded in the underlying FunctionsBuilder.
	prog, err := c.Compile(`RETURN DUP_FUNC()`)
	if err != nil {
		t.Fatalf("Compile after duplicate RegisterFunctions failed: %v", err)
	}

	out, err := prog.Run(context.Background())
	if err != nil {
		t.Fatalf("Run after duplicate RegisterFunctions failed: %v", err)
	}

	var result string
	if err := json.Unmarshal(out, &result); err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}

	if result != "dup" {
		t.Fatalf("expected \"dup\", got %q", result)
	}
}

func TestCompiler_RemoveFunction(t *testing.T) {
	c := compiler.New()

	_ = c.RegisterFunction("TEMP", func(ctx context.Context, args ...core.Value) (core.Value, error) {
		return core.WrapValue(runtime.None), nil
	})

	c.RemoveFunction("TEMP")

	// After removal, running a query using TEMP should fail (function not found at runtime).
	prog, err := c.Compile(`RETURN TEMP()`)
	if err != nil {
		// compile-time detection is also acceptable
		return
	}

	_, err = prog.Run(context.Background())
	if err == nil {
		t.Fatal("expected error when calling removed function, got nil")
	}
}
