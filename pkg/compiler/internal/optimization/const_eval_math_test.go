package optimization

import (
	"context"
	"math"
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

func TestFoldBinaryDivisionByZeroMatchesRuntimeSemantics(t *testing.T) {
	tests := []struct {
		left       runtime.Value
		right      runtime.Value
		name       string
		shouldFold bool
	}{
		{name: "int by int zero", left: runtime.NewInt(1), right: runtime.ZeroInt, shouldFold: false},
		{name: "float by int zero", left: runtime.NewFloat(1), right: runtime.ZeroInt, shouldFold: true},
		{name: "int by float zero", left: runtime.NewInt(1), right: runtime.ZeroFloat, shouldFold: true},
		{name: "float by float zero", left: runtime.NewFloat(1), right: runtime.ZeroFloat, shouldFold: true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			out, ok := foldBinary(bytecode.OpDiv, test.left, test.right, context.Background())
			if ok != test.shouldFold {
				t.Fatalf("expected fold=%v, got value=%v fold=%v", test.shouldFold, out, ok)
			}

			if !test.shouldFold {
				if out != nil {
					t.Fatalf("expected no folded value, got %v", out)
				}

				return
			}

			value, ok := out.(runtime.Float)
			if !ok || !math.IsInf(float64(value), 1) {
				t.Fatalf("expected positive infinity, got %T(%v)", out, out)
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
