package vm

import (
	"context"
	"errors"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/compiler"
	"github.com/MontFerret/ferret/v2/pkg/file"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

func compileProgram(t *testing.T, source string) *bytecode.Program {
	t.Helper()

	c := compiler.New()
	program, err := c.Compile(file.NewSource("test", source))
	if err != nil {
		t.Fatalf("compile failed: %v", err)
	}

	return program
}

func TestRunSafetyStrictRecoversPanics(t *testing.T) {
	program := compileProgram(t, "RETURN PANIC_FN()")

	instance := New(program)
	env, err := NewEnvironment([]EnvironmentOption{
		WithFunction("PANIC_FN", func(ctx context.Context, args ...runtime.Value) (runtime.Value, error) {
			panic("panic in host function")
		}),
	})
	if err != nil {
		t.Fatalf("environment build failed: %v", err)
	}

	_, err = instance.Run(context.Background(), env)
	if err == nil {
		t.Fatal("expected runtime error in strict mode")
	}

	var rtErr *RuntimeError
	if !errors.As(err, &rtErr) {
		t.Fatalf("expected wrapped runtime error, got %T", err)
	}
}

func TestRunSafetyFastPropagatesPanics(t *testing.T) {
	program := compileProgram(t, "RETURN PANIC_FN()")

	instance := NewWith(program, WithPanicPolicy(PanicPropagate))
	env, err := NewEnvironment([]EnvironmentOption{
		WithFunction("PANIC_FN", func(ctx context.Context, args ...runtime.Value) (runtime.Value, error) {
			panic("panic in host function")
		}),
	})
	if err != nil {
		t.Fatalf("environment build failed: %v", err)
	}

	defer func() {
		if recovered := recover(); recovered == nil {
			t.Fatal("expected panic to propagate in fast mode")
		}
	}()

	_, _ = instance.Run(context.Background(), env)
}

func TestRunSafetyFastStillWrapsReturnedErrors(t *testing.T) {
	program := compileProgram(t, "RETURN FAIL_FN()")

	instance := NewWith(program, WithPanicPolicy(PanicPropagate))
	env, err := NewEnvironment([]EnvironmentOption{
		WithFunction("FAIL_FN", func(ctx context.Context, args ...runtime.Value) (runtime.Value, error) {
			return runtime.None, errors.New("boom")
		}),
	})
	if err != nil {
		t.Fatalf("environment build failed: %v", err)
	}

	_, err = instance.Run(context.Background(), env)
	if err == nil {
		t.Fatal("expected runtime error")
	}

	var rtErr *RuntimeError
	if !errors.As(err, &rtErr) {
		t.Fatalf("expected wrapped runtime error, got %T", err)
	}
}
