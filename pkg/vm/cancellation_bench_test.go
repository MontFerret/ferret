package vm

import (
	"context"
	"strings"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/compiler"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/source"
	vmtest "github.com/MontFerret/ferret/v2/pkg/vm/test"
)

func BenchmarkCancellationSafepoint(b *testing.B) {
	b.Run("StraightLine", benchmarkStraightLineCancellationSafepoint)
	b.Run("Backedge", benchmarkBackedgeCancellationSafepoint)
	b.Run("SourcePointEveryInstruction", benchmarkDenseSourcePointCancellationSafepoint)
	b.Run("SourcePointStatementInterval", benchmarkIntervalSourcePointCancellationSafepoint)
	b.Run("UDFCallLoop", benchmarkUDFCallCancellationSafepoint)
	b.Run("RecursiveUDF", benchmarkRecursiveUDFCancellationSafepoint)
	b.Run("TailRecursiveUDF", benchmarkTailRecursiveUDFCancellationSafepoint)
	b.Run("IteratorLoop", benchmarkIteratorCancellationSafepoint)
	b.Run("HostCallSequence", benchmarkHostCallCancellationSafepoint)
	b.Run("ReleaseCompiled", benchmarkReleaseCompiledCancellationSafepoint)
	b.Run("DebugCompiled", benchmarkDebugCompiledCancellationSafepoint)
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

	benchmarkCancellationContexts(b, program, NewDefaultEnvironment())
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

	benchmarkCancellationContexts(b, program, NewDefaultEnvironment())
}

func benchmarkDenseSourcePointCancellationSafepoint(b *testing.B) {
	benchmarkSourcePointCancellationSafepoint(b, 1)
}

func benchmarkIntervalSourcePointCancellationSafepoint(b *testing.B) {
	benchmarkSourcePointCancellationSafepoint(b, 16)
}

func benchmarkSourcePointCancellationSafepoint(b *testing.B, interval int) {
	const instructionCount = 1024

	instructions := make([]bytecode.Instruction, 0, instructionCount+(instructionCount/interval)+1)
	for idx := range instructionCount {
		if idx%interval == 0 {
			instructions = append(instructions, bytecode.NewInstruction(bytecode.OpSourcePoint, bytecode.Operand(idx/interval)))
		}
		instructions = append(instructions, bytecode.NewInstruction(bytecode.OpLoadConst, bytecode.NewRegister(0), bytecode.NewConstant(0)))
	}
	instructions = append(instructions, bytecode.NewInstruction(bytecode.OpReturn, bytecode.NewRegister(0)))

	benchmarkCancellationContexts(b, newTestProgram(1, []runtime.Value{runtime.ZeroInt}, instructions...), NewDefaultEnvironment())
}

func benchmarkUDFCallCancellationSafepoint(b *testing.B) {
	program := benchmarkCompileCancellationProgram(b, mustNewCompiler(b, compiler.WithOptimizationLevel(compiler.None)), `
FUNC inc(value) => value + 1
RETURN FOR value IN 1..128 RETURN inc(value)
`)
	benchmarkCancellationContexts(b, program, NewDefaultEnvironment())
}

func benchmarkRecursiveUDFCancellationSafepoint(b *testing.B) {
	program := benchmarkCompileCancellationProgram(b, mustNewCompiler(b, compiler.WithOptimizationLevel(compiler.None)), `
FUNC fact(value) {
  RETURN MATCH value {
    0 => 1,
    _ => value * fact(value - 1),
  }
}
RETURN fact(16)
`)
	benchmarkCancellationContexts(b, program, NewDefaultEnvironment())
}

func benchmarkTailRecursiveUDFCancellationSafepoint(b *testing.B) {
	program := benchmarkCompileCancellationProgram(b, mustNewCompiler(b, compiler.WithOptimizationLevel(compiler.None)), `
FUNC count(value, total) {
  RETURN MATCH value {
    0 => total,
    _ => count(value - 1, total + 1),
  }
}
RETURN count(64, 0)
`)
	benchmarkCancellationContexts(b, program, NewDefaultEnvironment())
}

func benchmarkIteratorCancellationSafepoint(b *testing.B) {
	program := benchmarkCompileCancellationProgram(b, mustNewCompiler(b, compiler.WithOptimizationLevel(compiler.None)), "FOR value IN 1..1024 RETURN value")
	benchmarkCancellationContexts(b, program, NewDefaultEnvironment())
}

func benchmarkHostCallCancellationSafepoint(b *testing.B) {
	var query strings.Builder
	for idx := range 64 {
		query.WriteString("LET value")
		query.WriteString(runtime.NewInt(idx).String())
		query.WriteString(" = PING()\n")
	}
	query.WriteString("RETURN 1")

	program := benchmarkCompileCancellationProgram(b, mustNewCompiler(b, compiler.WithOptimizationLevel(compiler.None)), query.String())
	env, err := NewEnvironment([]EnvironmentOption{WithFunction("PING", func(context.Context, ...runtime.Value) (runtime.Value, error) {
		return runtime.NewInt(1), nil
	})})
	if err != nil {
		b.Fatal(err)
	}
	benchmarkCancellationContexts(b, program, env)
}

func benchmarkReleaseCompiledCancellationSafepoint(b *testing.B) {
	program := benchmarkCompileCancellationProgram(b, mustNewCompiler(b, compiler.WithOptimizationLevel(compiler.None)), "FOR value IN 1..128 RETURN value * 2")
	benchmarkCancellationContexts(b, program, NewDefaultEnvironment())
}

func benchmarkDebugCompiledCancellationSafepoint(b *testing.B) {
	program := benchmarkCompileCancellationProgram(b, mustNewCompiler(b, compiler.WithDebugInfo()), "FOR value IN 1..128 RETURN value * 2")
	benchmarkCancellationContexts(b, program, NewDefaultEnvironment())
}

func benchmarkCompileCancellationProgram(b *testing.B, c *compiler.Compiler, query string) *bytecode.Program {
	b.Helper()

	program, err := c.Compile(source.NewAnonymous(query))
	if err != nil {
		b.Fatal(err)
	}

	return program
}

func benchmarkCancellationContexts(b *testing.B, program *bytecode.Program, env *Environment) {
	b.Run("Background", func(b *testing.B) {
		benchmarkCancellationSafepoint(b, program, env, context.Background())
	})

	b.Run("Cancelable", func(b *testing.B) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		benchmarkCancellationSafepoint(b, program, env, ctx)
	})
}

func benchmarkCancellationSafepoint(b *testing.B, program *bytecode.Program, env *Environment, ctx context.Context) {
	instance, err := NewWith(program, WithTesting(vmtest.WithBenchmarkMode()))
	if err != nil {
		b.Fatalf("vm init failed: %v", err)
	}
	defer func() { _ = instance.Close() }()

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
