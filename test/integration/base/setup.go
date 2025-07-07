package base

import (
	"context"
	"testing"

	"github.com/MontFerret/ferret/pkg/compiler"
	"github.com/MontFerret/ferret/pkg/vm"
)

func RunBenchmarkWith(b *testing.B, c *compiler.Compiler, expression string, opts ...vm.EnvironmentOption) {
	prog, err := c.Compile(expression)

	if err != nil {
		panic(err)
	}

	options := []vm.EnvironmentOption{
		vm.WithFunctions(Stdlib()),
	}
	options = append(options, opts...)

	ctx := context.Background()
	instance := vm.New(prog)

	b.ResetTimer()

	for b.Loop() {
		_, err := instance.Run(ctx, opts)

		if err != nil {
			panic(err)
		}
	}
}

func RunBenchmark(b *testing.B, expression string, opts ...vm.EnvironmentOption) {
	RunBenchmarkWith(b, compiler.New(), expression, opts...)
}
