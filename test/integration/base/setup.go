package base

import (
	"context"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/asm"
	"github.com/MontFerret/ferret/v2/pkg/file"

	"github.com/MontFerret/ferret/v2/pkg/compiler"
	"github.com/MontFerret/ferret/v2/pkg/vm"
)

func RunBenchmarkWith(b *testing.B, c *compiler.Compiler, expression string, opts ...vm.EnvironmentOption) {
	prog, err := c.Compile(file.NewSource("benchmark", expression))

	if err != nil {
		panic(err)
	}

	options := []vm.EnvironmentOption{
		vm.WithNamespace(Stdlib()),
	}
	options = append(options, opts...)

	ctx := context.Background()
	instance := vm.New(prog)
	env, err := vm.NewEnvironment(options)
	if err != nil {
		panic(err)
	}

	if testing.Verbose() {
		println("Query:")
		println(expression)
		println("Bytecode:")
		println(asm.Disassemble(prog))
	}

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		_, err := instance.Run(ctx, env)

		if err != nil {
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
