package vm

import (
	"context"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	vmtest "github.com/MontFerret/ferret/v2/pkg/vm/test"
)

func BenchmarkElapsedOpcode(b *testing.B) {
	program := newTestProgram(
		1,
		nil,
		bytecode.NewInstruction(bytecode.OpElapsed, bytecode.NewRegister(0)),
		bytecode.NewInstruction(bytecode.OpReturn, bytecode.NewRegister(0)),
	)
	instance, err := NewWith(program, WithTesting(vmtest.WithBenchmarkMode()))
	if err != nil {
		b.Fatalf("vm init failed: %v", err)
	}
	defer func() { _ = instance.Close() }()

	ctx := context.Background()
	env := NewDefaultEnvironment()
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		result, runErr := instance.Run(ctx, env)
		if runErr != nil {
			b.Fatalf("run failed: %v", runErr)
		}
		if closeErr := result.Close(); closeErr != nil {
			b.Fatalf("result close failed: %v", closeErr)
		}
	}
}
