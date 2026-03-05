package optimization

import (
	"reflect"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
)

type instructionUseDefCase struct {
	name     string
	inst     bytecode.Instruction
	wantUses []int
	wantDefs []int
}

var instructionUseDefCases = []instructionUseDefCase{
	{
		name:     "move uses source and defines destination",
		inst:     bytecode.NewInstruction(bytecode.OpMove, bytecode.NewRegister(2), bytecode.NewRegister(1)),
		wantUses: []int{1},
		wantDefs: []int{2},
	},
	{
		name:     "concat uses contiguous register range",
		inst:     bytecode.NewInstruction(bytecode.OpConcat, bytecode.NewRegister(6), bytecode.NewRegister(3), bytecode.Operand(2)),
		wantUses: []int{3, 4},
		wantDefs: []int{6},
	},
	{
		name:     "add const ignores constant operand",
		inst:     bytecode.NewInstruction(bytecode.OpAddConst, bytecode.NewRegister(4), bytecode.NewRegister(1), bytecode.NewConstant(0)),
		wantUses: []int{1},
		wantDefs: []int{4},
	},
	{
		name:     "jump if ne const uses only compared register",
		inst:     bytecode.NewInstruction(bytecode.OpJumpIfNeConst, bytecode.Operand(7), bytecode.NewRegister(5), bytecode.NewConstant(1)),
		wantUses: []int{5},
		wantDefs: nil,
	},
	{
		name:     "iterator limit updates iterator register",
		inst:     bytecode.NewInstruction(bytecode.OpIterLimit, bytecode.NewRegister(9), bytecode.NewRegister(2), bytecode.NewRegister(8)),
		wantUses: []int{2, 8},
		wantDefs: []int{2},
	},
	{
		name:     "stream uses all inputs and defines destination",
		inst:     bytecode.NewInstruction(bytecode.OpStream, bytecode.NewRegister(3), bytecode.NewRegister(1), bytecode.NewRegister(2)),
		wantUses: []int{3, 1, 2},
		wantDefs: []int{3},
	},
	{
		name:     "close uses and defines same register",
		inst:     bytecode.NewInstruction(bytecode.OpClose, bytecode.NewRegister(4)),
		wantUses: []int{4},
		wantDefs: []int{4},
	},
}

func runInstructionUseDefCase(t *testing.T, tc instructionUseDefCase) {
	t.Helper()

	gotUses, gotDefs := instructionUseDef(tc.inst)
	if !reflect.DeepEqual(gotUses, tc.wantUses) {
		t.Fatalf("unexpected uses: got %v, want %v", gotUses, tc.wantUses)
	}
	if !reflect.DeepEqual(gotDefs, tc.wantDefs) {
		t.Fatalf("unexpected defs: got %v, want %v", gotDefs, tc.wantDefs)
	}
}

func TestInstructionUseDef(t *testing.T) {
	for _, tc := range instructionUseDefCases {
		t.Run(tc.name, func(t *testing.T) {
			runInstructionUseDefCase(t, tc)
		})
	}
}
