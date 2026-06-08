package vm

import (
	"context"
	"errors"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

func TestIncDecRejectNonNumericOperands(t *testing.T) {
	for _, op := range []bytecode.Opcode{bytecode.OpIncr, bytecode.OpDecr} {
		t.Run(op.String(), func(t *testing.T) {
			program := newTestProgram(
				1,
				[]runtime.Value{runtime.NewString("3")},
				bytecode.NewInstruction(bytecode.OpLoadConst, bytecode.NewRegister(0), bytecode.NewConstant(0)),
				bytecode.NewInstruction(op, bytecode.NewRegister(0)),
				bytecode.NewInstruction(bytecode.OpReturn, bytecode.NewRegister(0)),
			)

			_, err := mustNewVM(t, program).Run(context.Background(), NewDefaultEnvironment())
			if err == nil {
				t.Fatal("expected runtime error")
			}

			var runtimeErr *RuntimeError
			if !errors.As(err, &runtimeErr) {
				t.Fatalf("expected runtime error, got %T", err)
			}
			if runtimeErr.Message != "invalid type" {
				t.Fatalf("unexpected runtime error message: got %q, want %q", runtimeErr.Message, "invalid type")
			}
		})
	}
}
