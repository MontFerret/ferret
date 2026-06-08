package optimization

import (
	"context"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

func TestFoldBinaryRejectsNonNumericMathOperands(t *testing.T) {
	for _, op := range []bytecode.Opcode{bytecode.OpSub, bytecode.OpMul, bytecode.OpDiv, bytecode.OpMod} {
		t.Run(op.String(), func(t *testing.T) {
			if out, ok := foldBinary(op, runtime.NewString("3"), runtime.NewInt(2), context.Background()); ok || out != nil {
				t.Fatalf("expected %s with non-numeric operand not to fold, got value=%v ok=%v", op, out, ok)
			}
		})
	}
}

func TestFoldUnaryRejectsNonNumericMathOperands(t *testing.T) {
	for _, op := range []bytecode.Opcode{bytecode.OpFlipPositive, bytecode.OpFlipNegative} {
		t.Run(op.String(), func(t *testing.T) {
			if out, ok := foldUnary(op, runtime.NewString("3"), context.Background()); ok || out != nil {
				t.Fatalf("expected %s with non-numeric operand not to fold, got value=%v ok=%v", op, out, ok)
			}
		})
	}
}

func TestFoldIncDecRejectsNonNumericOperands(t *testing.T) {
	for _, op := range []bytecode.Opcode{bytecode.OpIncr, bytecode.OpDecr} {
		t.Run(op.String(), func(t *testing.T) {
			if out, ok := foldIncDecValue(op, runtime.NewString("3"), context.Background()); ok || out != nil {
				t.Fatalf("expected %s with non-numeric operand not to fold, got value=%v ok=%v", op, out, ok)
			}
		})
	}
}
