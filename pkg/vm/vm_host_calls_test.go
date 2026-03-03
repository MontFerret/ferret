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

func TestRunReturnsUnresolvedFunctionWhenHostCacheEntryIsMissing(t *testing.T) {
	c := compiler.New()
	program, err := c.Compile(file.NewSource("test", "RETURN TEST(1)"))
	if err != nil {
		t.Fatalf("compile failed: %v", err)
	}

	instance := New(program)
	env, err := NewEnvironment([]EnvironmentOption{
		WithFunction("TEST", func(ctx context.Context, args ...runtime.Value) (runtime.Value, error) {
			return runtime.True, nil
		}),
	})
	if err != nil {
		t.Fatalf("environment build failed: %v", err)
	}

	if _, err := instance.Run(context.Background(), env); err != nil {
		t.Fatalf("first run failed: %v", err)
	}

	hostPC := -1
	for i, inst := range program.Bytecode {
		if inst.Opcode == bytecode.OpHCall || inst.Opcode == bytecode.OpProtectedHCall {
			hostPC = i
			break
		}
	}

	if hostPC < 0 {
		t.Fatal("host call opcode not found")
	}

	instance.cache.HostFunctions[hostPC] = nil

	_, err = instance.Run(context.Background(), env)
	if err == nil {
		t.Fatal("expected unresolved function error")
	}

	var rtErr *RuntimeError
	if !errors.As(err, &rtErr) {
		t.Fatalf("expected runtime error, got %T", err)
	}

	if rtErr.Message != "Unresolved function" {
		t.Fatalf("expected unresolved function message, got %q", rtErr.Message)
	}
}

func TestSetCallResult_AppliesCatchJumpZeroAndFallbackValue(t *testing.T) {
	instance := New(&bytecode.Program{
		Registers: 2,
		Bytecode: []bytecode.Instruction{
			bytecode.NewInstruction(bytecode.OpLoadNone, bytecode.NewRegister(1)),
			bytecode.NewInstruction(bytecode.OpLoadNone, bytecode.NewRegister(1)),
		},
		CatchTable: []bytecode.Catch{
			{1, 1, 0},
		},
	})

	instance.pc = 1
	instance.registers.Values[1] = runtime.True

	err := instance.setCallResult(
		bytecode.OpHCall,
		bytecode.NewRegister(1),
		runtime.True,
		errors.New("boom"),
	)

	if err != nil {
		t.Fatalf("expected caught error to be swallowed, got %v", err)
	}

	if got := instance.registers.Values[1]; got != runtime.None {
		t.Fatalf("expected destination to be reset to none, got %v", got)
	}

	if got, want := instance.pc, 0; got != want {
		t.Fatalf("expected catch jump target %d, got %d", want, got)
	}
}
