package vm

import (
	"context"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/rnd"
	vmtest "github.com/MontFerret/ferret/v2/pkg/vm/test"
)

func BenchmarkRandOpcode(b *testing.B) {
	program := newTestProgram(1, nil,
		bytecode.NewInstruction(bytecode.OpRand, bytecode.NewRegister(0)),
		bytecode.NewInstruction(bytecode.OpReturn, bytecode.NewRegister(0)),
	)
	instance, err := NewWith(program, WithTesting(vmtest.WithBenchmarkMode()))
	if err != nil {
		b.Fatal(err)
	}
	b.Cleanup(func() { _ = instance.Close() })
	ctx := rnd.WithContext(context.Background(), rnd.NewSeed(0))
	env := NewDefaultEnvironment()
	b.ReportAllocs()
	b.ResetTimer()

	for b.Loop() {
		result, err := instance.Run(ctx, env)
		if err != nil {
			b.Fatal(err)
		}

		if err := result.Close(); err != nil {
			b.Fatal(err)
		}
	}
}
