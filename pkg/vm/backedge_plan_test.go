package vm

import (
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
)

func TestExecInstructionResolvesJumpBackedges(t *testing.T) {
	opcodes := []bytecode.Opcode{
		bytecode.OpJump,
		bytecode.OpJumpIfFalse,
		bytecode.OpJumpIfTrue,
		bytecode.OpJumpIfNone,
		bytecode.OpJumpIfNe,
		bytecode.OpJumpIfNeConst,
		bytecode.OpJumpIfEq,
		bytecode.OpJumpIfEqConst,
		bytecode.OpJumpIfMissingProperty,
		bytecode.OpJumpIfMissingPropertyConst,
		bytecode.OpIterNext,
		bytecode.OpIterNextTimeout,
		bytecode.OpIterSkip,
		bytecode.OpIterLimit,
	}
	directions := []struct {
		name     string
		target   bytecode.Operand
		expected bool
	}{
		{name: "forward", target: 3, expected: false},
		{name: "backward", target: 1, expected: true},
		{name: "self", target: 2, expected: true},
	}

	for _, op := range opcodes {
		for _, direction := range directions {
			t.Run(op.String()+"/"+direction.name, func(t *testing.T) {
				program := &bytecode.Program{Bytecode: []bytecode.Instruction{
					bytecode.NewInstruction(bytecode.OpLoadZero, bytecode.NewRegister(0)),
					bytecode.NewInstruction(bytecode.OpLoadZero, bytecode.NewRegister(1)),
					bytecode.NewInstruction(op, direction.target, bytecode.NewRegister(0), bytecode.NewRegister(1)),
					bytecode.NewInstruction(bytecode.OpReturn, bytecode.NewRegister(0)),
				}}

				plan, err := buildExecPlan(program)
				if err != nil {
					t.Fatal(err)
				}
				if got := plan.instructions[2].isBackedge(2); got != direction.expected {
					t.Fatalf("isBackedge = %t, want %t", got, direction.expected)
				}
			})
		}
	}
}

func TestExecInstructionResolvesMatchBackedges(t *testing.T) {
	for _, test := range []struct {
		name     string
		target   int
		expected bool
	}{
		{name: "forward", target: 3, expected: false},
		{name: "backward", target: 1, expected: true},
		{name: "self", target: 2, expected: true},
	} {
		t.Run(test.name, func(t *testing.T) {
			program := &bytecode.Program{
				Bytecode: []bytecode.Instruction{
					bytecode.NewInstruction(bytecode.OpLoadZero, bytecode.NewRegister(0)),
					bytecode.NewInstruction(bytecode.OpLoadZero, bytecode.NewRegister(1)),
					bytecode.NewInstruction(bytecode.OpMatchLoadPropertyConst, bytecode.NewRegister(1), bytecode.NewRegister(0), bytecode.NewConstant(0)),
					bytecode.NewInstruction(bytecode.OpReturn, bytecode.NewRegister(0)),
				},
				Metadata: bytecode.Metadata{MatchFailTargets: []int{-1, -1, test.target, -1}},
			}

			plan, err := buildExecPlan(program)
			if err != nil {
				t.Fatal(err)
			}
			if got := plan.instructions[2].isBackedge(2); got != test.expected {
				t.Fatalf("isBackedge = %t, want %t", got, test.expected)
			}
			if got := plan.instructions[2].InlineSlot; got != test.target {
				t.Fatalf("match fail target = %d, want %d", got, test.target)
			}
		})
	}
}
