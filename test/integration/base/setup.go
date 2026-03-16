package base

import (
	"context"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/asm"
	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/file"

	"github.com/MontFerret/ferret/v2/pkg/compiler"
	"github.com/MontFerret/ferret/v2/pkg/vm"
)

func compileBenchmarkProgram(c *compiler.Compiler, expression string) *bytecode.Program {
	prog, err := c.Compile(file.NewSource("benchmark", expression))

	if err != nil {
		panic(err)
	}

	return prog
}

func newBenchmarkEnvironment(envOpts ...vm.EnvironmentOption) *vm.Environment {
	options := []vm.EnvironmentOption{
		vm.WithNamespace(Stdlib()),
	}
	options = append(options, envOpts...)

	env, err := vm.NewEnvironment(options)
	if err != nil {
		panic(err)
	}

	return env
}

func newBenchmarkVM(program *bytecode.Program, vmOpts ...vm.Option) *vm.VM {
	instance, err := vm.NewWith(program, vmOpts...)
	if err != nil {
		panic(err)
	}

	return instance
}

func prepareBenchmark(c *compiler.Compiler, expression string, vmOpts []vm.Option, envOpts ...vm.EnvironmentOption) (*bytecode.Program, *vm.VM, *vm.Environment) {
	prog := compileBenchmarkProgram(c, expression)

	instance := newBenchmarkVM(prog, vmOpts...)
	env := newBenchmarkEnvironment(envOpts...)

	return prog, instance, env
}

func RunBenchmarkWith(b *testing.B, c *compiler.Compiler, expression string, opts ...vm.EnvironmentOption) {
	prog, instance, env := prepareBenchmark(c, expression, []vm.Option{vm.WithBenchmarkResultMode()}, opts...)
	ctx := context.Background()

	if testing.Verbose() {
		println("Query:")
		println(expression)
		println("Bytecode:")
		println(asm.Disassemble(prog))
	}

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		result, err := instance.Run(ctx, env)

		if err != nil {
			panic(err)
		}

		if err := result.Close(); err != nil {
			panic(err)
		}
	}
}

func RunBenchmark(b *testing.B, expression string, opts ...vm.EnvironmentOption) {
	RunBenchmarkWithOptimization(b, expression, compiler.O1, opts...)
}

func RunBenchmarkWithOptimization(b *testing.B, expression string, level compiler.OptimizationLevel, opts ...vm.EnvironmentOption) {
	RunBenchmarkWith(b, compiler.New(compiler.WithOptimizationLevel(level)), expression, opts...)
}

func RunResultBenchmarkWith(b *testing.B, c *compiler.Compiler, expression string, opts ...vm.EnvironmentOption) {
	prog, instance, env := prepareBenchmark(c, expression, nil, opts...)
	ctx := context.Background()

	if testing.Verbose() {
		println("Query:")
		println(expression)
		println("Bytecode:")
		println(asm.Disassemble(prog))
	}

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		result, err := instance.Run(ctx, env)

		if err != nil {
			panic(err)
		}

		if err := result.Close(); err != nil {
			panic(err)
		}
	}
}

func RunResultBenchmark(b *testing.B, expression string, opts ...vm.EnvironmentOption) {
	RunResultBenchmarkWithOptimization(b, expression, compiler.O1, opts...)
}

func RunResultBenchmarkWithOptimization(b *testing.B, expression string, level compiler.OptimizationLevel, opts ...vm.EnvironmentOption) {
	RunResultBenchmarkWith(b, compiler.New(compiler.WithOptimizationLevel(level)), expression, opts...)
}

func RunOutputBenchmarkWith(b *testing.B, c *compiler.Compiler, expression string, opts ...vm.EnvironmentOption) {
	prog, instance, env := prepareBenchmark(c, expression, nil, opts...)
	ctx := newTestContext()

	if testing.Verbose() {
		println("Query:")
		println(expression)
		println("Bytecode:")
		println(asm.Disassemble(prog))
	}

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		result, err := instance.Run(ctx, env)

		if err != nil {
			panic(err)
		}

		if _, err := materializeJSONResult(result); err != nil {
			panic(err)
		}
	}
}

func RunOutputBenchmark(b *testing.B, expression string, opts ...vm.EnvironmentOption) {
	RunOutputBenchmarkWithOptimization(b, expression, compiler.O1, opts...)
}

func RunOutputBenchmarkWithOptimization(b *testing.B, expression string, level compiler.OptimizationLevel, opts ...vm.EnvironmentOption) {
	RunOutputBenchmarkWith(b, compiler.New(compiler.WithOptimizationLevel(level)), expression, opts...)
}
