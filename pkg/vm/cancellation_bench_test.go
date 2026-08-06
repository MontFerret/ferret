package vm

import (
	"context"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
	vmtest "github.com/MontFerret/ferret/v2/pkg/vm/test"
)

func BenchmarkCancellationSafepoint(b *testing.B) {
	b.Run("StraightLine", benchmarkStraightLineCancellationSafepoint)
	b.Run("Backedge", benchmarkBackedgeCancellationSafepoint)
}

func benchmarkStraightLineCancellationSafepoint(b *testing.B) {
	const instructionCount = 1024

	instructions := make([]bytecode.Instruction, 0, instructionCount+1)
	for range instructionCount {
		instructions = append(instructions, bytecode.NewInstruction(
			bytecode.OpLoadConst,
			bytecode.NewRegister(0),
			bytecode.NewConstant(0),
		))
	}
	instructions = append(instructions, bytecode.NewInstruction(bytecode.OpReturn, bytecode.NewRegister(0)))

	program := newTestProgram(1, []runtime.Value{runtime.ZeroInt}, instructions...)

	benchmarkCancellationContexts(b, program)
}

func benchmarkBackedgeCancellationSafepoint(b *testing.B) {
	const limit = 1024

	program := newTestProgram(
		3,
		[]runtime.Value{runtime.NewInt(limit)},
		bytecode.NewInstruction(bytecode.OpLoadZero, bytecode.NewRegister(0)),
		bytecode.NewInstruction(bytecode.OpLoadConst, bytecode.NewRegister(1), bytecode.NewConstant(0)),
		bytecode.NewInstruction(bytecode.OpLt, bytecode.NewRegister(2), bytecode.NewRegister(0), bytecode.NewRegister(1)),
		bytecode.NewInstruction(bytecode.OpJumpIfFalse, bytecode.Operand(6), bytecode.NewRegister(2)),
		bytecode.NewInstruction(bytecode.OpIncr, bytecode.NewRegister(0)),
		bytecode.NewInstruction(bytecode.OpJump, bytecode.Operand(2)),
		bytecode.NewInstruction(bytecode.OpReturn, bytecode.NewRegister(0)),
	)

	benchmarkCancellationContexts(b, program)
}

func benchmarkCancellationContexts(b *testing.B, program *bytecode.Program) {
	b.Run("Background", func(b *testing.B) {
		benchmarkCancellationSafepoint(b, program, context.Background())
	})

	b.Run("Cancelable", func(b *testing.B) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		benchmarkCancellationSafepoint(b, program, ctx)
	})
}

func benchmarkCancellationSafepoint(b *testing.B, program *bytecode.Program, ctx context.Context) {
	instance, err := NewWith(program, WithTesting(vmtest.WithBenchmarkMode()))
	if err != nil {
		b.Fatalf("vm init failed: %v", err)
	}
	defer func() { _ = instance.Close() }()

	env := NewDefaultEnvironment()
	b.ReportAllocs()
	b.ResetTimer()

	for range b.N {
		result, runErr := instance.Run(ctx, env)
		if runErr != nil {
			b.Fatalf("run failed: %v", runErr)
		}
		if closeErr := result.Close(); closeErr != nil {
			b.Fatalf("result close failed: %v", closeErr)
		}
	}
}
