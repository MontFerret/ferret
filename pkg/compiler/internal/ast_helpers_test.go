package internal

import (
	"reflect"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

func TestCompileFirstOperand_Order(t *testing.T) {
	order := make([]int, 0)

	op := compileFirstOperand(
		newOperandBranch(true, func() bytecode.Operand {
			order = append(order, 1)
			return bytecode.NoopOperand
		}),
		newOperandBranch(false, func() bytecode.Operand {
			t.Fatal("disabled branch should not execute")
			return bytecode.NoopOperand
		}),
		newOperandBranch(true, func() bytecode.Operand {
			order = append(order, 2)
			return bytecode.NewRegister(7)
		}),
		newOperandBranch(true, func() bytecode.Operand {
			order = append(order, 3)
			return bytecode.NewRegister(8)
		}),
	)

	if op != bytecode.NewRegister(7) {
		t.Fatalf("unexpected operand: got %s", op)
	}

	if !reflect.DeepEqual(order, []int{1, 2}) {
		t.Fatalf("unexpected execution order: got %v", order)
	}
}

func TestCompileFirstOperand_NoCandidate(t *testing.T) {
	op := compileFirstOperand(
		newOperandBranch(false, func() bytecode.Operand { return bytecode.NewRegister(1) }),
		newOperandBranch(true, func() bytecode.Operand { return bytecode.NoopOperand }),
	)

	if op != bytecode.NoopOperand {
		t.Fatalf("expected noop operand, got %s", op)
	}
}

func TestLiteralBooleanValue(t *testing.T) {
	val, ok := literalBooleanValue("true")
	if !ok || val != runtime.True {
		t.Fatalf("expected true literal, got (%v, %v)", val, ok)
	}

	val, ok = literalBooleanValue("FALSE")
	if !ok || val != runtime.False {
		t.Fatalf("expected false literal, got (%v, %v)", val, ok)
	}

	if _, ok = literalBooleanValue("invalid"); ok {
		t.Fatal("expected invalid boolean literal to fail")
	}
}

func TestLiteralNumericValue(t *testing.T) {
	floatVal, ok := literalFloatValue("1.5")
	if !ok {
		t.Fatal("expected float literal to parse")
	}

	f, ok := floatVal.(runtime.Float)
	if !ok || f != runtime.NewFloat(1.5) {
		t.Fatalf("unexpected float literal result: %v", floatVal)
	}

	if _, ok = literalFloatValue("a"); ok {
		t.Fatal("expected invalid float literal to fail")
	}

	intVal, ok := literalIntValue("42")
	if !ok {
		t.Fatal("expected int literal to parse")
	}

	i, ok := intVal.(runtime.Int)
	if !ok || i != runtime.NewInt(42) {
		t.Fatalf("unexpected int literal result: %v", intVal)
	}

	if _, ok = literalIntValue("x"); ok {
		t.Fatal("expected invalid int literal to fail")
	}
}

func TestResolveWaitPredicateMode(t *testing.T) {
	tests := []struct {
		name      string
		hasValue  bool
		hasExists bool
		hasNot    bool
		expected  waitForPredicateMode
	}{
		{
			name:      "value takes precedence",
			hasValue:  true,
			hasExists: true,
			hasNot:    true,
			expected:  waitForPredicateModeValue,
		},
		{
			name:      "exists",
			hasExists: true,
			expected:  waitForPredicateModeExists,
		},
		{
			name:      "not exists",
			hasExists: true,
			hasNot:    true,
			expected:  waitForPredicateModeNotExists,
		},
		{
			name:     "bool fallback",
			expected: waitForPredicateModeBool,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mode := resolveWaitPredicateMode(tt.hasValue, tt.hasExists, tt.hasNot)
			if mode != tt.expected {
				t.Fatalf("unexpected mode: got %d want %d", mode, tt.expected)
			}
		})
	}
}
