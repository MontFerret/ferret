package optimization

import (
	"reflect"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/file"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

func runPeephole(t *testing.T, program *bytecode.Program) (*PassResult, error) {
	t.Helper()
	pass := NewPeepholePass()
	return pass.Run(&PassContext{Program: program})
}

func TestPeephole_RemovesRedundantLoads(t *testing.T) {
	program := &bytecode.Program{
		Constants: []runtime.Value{
			runtime.NewString("a"),
		},
		Bytecode: []bytecode.Instruction{
			bytecode.NewInstruction(bytecode.OpLoadConst, bytecode.NewRegister(1), bytecode.NewConstant(0)),
			bytecode.NewInstruction(bytecode.OpLoadConst, bytecode.NewRegister(1), bytecode.NewConstant(0)),
			bytecode.NewInstruction(bytecode.OpReturn, bytecode.NewRegister(1)),
		},
	}

	res, err := runPeephole(t, program)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !res.Modified {
		t.Fatalf("expected peephole pass to modify program")
	}
	if len(program.Bytecode) != 2 {
		t.Fatalf("expected 2 instructions, got %d", len(program.Bytecode))
	}
	if program.Bytecode[0].Opcode != bytecode.OpLoadConst {
		t.Fatalf("expected first instruction to be LOADC, got %s", program.Bytecode[0].Opcode)
	}
	if program.Bytecode[1].Opcode != bytecode.OpReturn {
		t.Fatalf("expected last instruction to be RETURN, got %s", program.Bytecode[1].Opcode)
	}
}

func TestPeephole_RewritesAddConstRight(t *testing.T) {
	program := &bytecode.Program{
		Constants: []runtime.Value{
			runtime.NewInt(2),
			runtime.NewInt(3),
		},
		Bytecode: []bytecode.Instruction{
			bytecode.NewInstruction(bytecode.OpLoadConst, bytecode.NewRegister(1), bytecode.NewConstant(0)),
			bytecode.NewInstruction(bytecode.OpLoadConst, bytecode.NewRegister(2), bytecode.NewConstant(1)),
			bytecode.NewInstruction(bytecode.OpAdd, bytecode.NewRegister(3), bytecode.NewRegister(1), bytecode.NewRegister(2)),
			bytecode.NewInstruction(bytecode.OpReturn, bytecode.NewRegister(3)),
		},
	}

	res, err := runPeephole(t, program)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !res.Modified {
		t.Fatalf("expected peephole pass to modify program")
	}
	if len(program.Bytecode) != 3 {
		t.Fatalf("expected 3 instructions, got %d", len(program.Bytecode))
	}
	add := program.Bytecode[1]
	if add.Opcode != bytecode.OpAddConst {
		t.Fatalf("expected ADDC, got %s", add.Opcode)
	}
	if !add.Operands[2].IsConstant() || add.Operands[2].Constant() != 1 {
		t.Fatalf("expected constant operand C1, got %s", add.Operands[2])
	}
}

func TestPeephole_RemovesSelfMove(t *testing.T) {
	program := &bytecode.Program{
		Constants: []runtime.Value{
			runtime.NewString("a"),
		},
		Bytecode: []bytecode.Instruction{
			bytecode.NewInstruction(bytecode.OpLoadConst, bytecode.NewRegister(1), bytecode.NewConstant(0)),
			bytecode.NewInstruction(bytecode.OpMove, bytecode.NewRegister(1), bytecode.NewRegister(1)),
			bytecode.NewInstruction(bytecode.OpReturn, bytecode.NewRegister(1)),
		},
	}

	res, err := runPeephole(t, program)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !res.Modified {
		t.Fatalf("expected peephole pass to modify program")
	}
	if len(program.Bytecode) != 2 {
		t.Fatalf("expected 2 instructions, got %d", len(program.Bytecode))
	}
	if program.Bytecode[0].Opcode != bytecode.OpLoadConst {
		t.Fatalf("expected first instruction to be LOADC, got %s", program.Bytecode[0].Opcode)
	}
	if program.Bytecode[1].Opcode != bytecode.OpReturn {
		t.Fatalf("expected last instruction to be RETURN, got %s", program.Bytecode[1].Opcode)
	}
}

func TestPeephole_UpdatesJumpTargets(t *testing.T) {
	program := &bytecode.Program{
		Constants: []runtime.Value{
			runtime.NewString("a"),
			runtime.NewString("b"),
		},
		Bytecode: []bytecode.Instruction{
			bytecode.NewInstruction(bytecode.OpLoadConst, bytecode.NewRegister(1), bytecode.NewConstant(0)),
			bytecode.NewInstruction(bytecode.OpLoadConst, bytecode.NewRegister(1), bytecode.NewConstant(1)),
			bytecode.NewInstruction(bytecode.OpJump, bytecode.Operand(3)),
			bytecode.NewInstruction(bytecode.OpReturn, bytecode.NewRegister(1)),
		},
	}

	res, err := runPeephole(t, program)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !res.Modified {
		t.Fatalf("expected peephole pass to modify program")
	}
	if len(program.Bytecode) != 3 {
		t.Fatalf("expected 3 instructions, got %d", len(program.Bytecode))
	}
	jump := program.Bytecode[1]
	if jump.Opcode != bytecode.OpJump {
		t.Fatalf("expected instruction 1 to be JUMP, got %s", jump.Opcode)
	}
	if jump.Operands[0] != bytecode.Operand(2) {
		t.Fatalf("expected jump target to be 2, got %d", jump.Operands[0])
	}
}

func TestPeephole_SkipsWhenNextUsesDef(t *testing.T) {
	program := &bytecode.Program{
		Constants: []runtime.Value{
			runtime.NewInt(1),
			runtime.NewInt(2),
		},
		Bytecode: []bytecode.Instruction{
			bytecode.NewInstruction(bytecode.OpLoadConst, bytecode.NewRegister(2), bytecode.NewConstant(0)),
			bytecode.NewInstruction(bytecode.OpLoadConst, bytecode.NewRegister(1), bytecode.NewConstant(1)),
			bytecode.NewInstruction(bytecode.OpAdd, bytecode.NewRegister(1), bytecode.NewRegister(1), bytecode.NewRegister(2)),
			bytecode.NewInstruction(bytecode.OpReturn, bytecode.NewRegister(1)),
		},
	}

	res, err := runPeephole(t, program)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Modified {
		t.Fatalf("expected peephole pass to keep program intact")
	}
	if len(program.Bytecode) != 4 {
		t.Fatalf("expected 4 instructions, got %d", len(program.Bytecode))
	}
}

func TestPeephole_SkipsRemovalWhenTargetedByJump(t *testing.T) {
	program := &bytecode.Program{
		Constants: []runtime.Value{
			runtime.NewString("a"),
		},
		Bytecode: []bytecode.Instruction{
			bytecode.NewInstruction(bytecode.OpJump, bytecode.Operand(1)),
			bytecode.NewInstruction(bytecode.OpMove, bytecode.NewRegister(1), bytecode.NewRegister(1)),
			bytecode.NewInstruction(bytecode.OpLoadConst, bytecode.NewRegister(1), bytecode.NewConstant(0)),
			bytecode.NewInstruction(bytecode.OpReturn, bytecode.NewRegister(1)),
		},
	}

	res, err := runPeephole(t, program)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Modified {
		t.Fatalf("expected peephole pass to keep jump target instruction")
	}
	if len(program.Bytecode) != 4 {
		t.Fatalf("expected 4 instructions, got %d", len(program.Bytecode))
	}
}

func TestPeephole_RemapsCatchDebugSpansAndLabels(t *testing.T) {
	program := &bytecode.Program{
		Constants: []runtime.Value{
			runtime.NewString("a"),
			runtime.NewString("b"),
		},
		Bytecode: []bytecode.Instruction{
			bytecode.NewInstruction(bytecode.OpLoadConst, bytecode.NewRegister(2), bytecode.NewConstant(0)),
			bytecode.NewInstruction(bytecode.OpLoadConst, bytecode.NewRegister(1), bytecode.NewConstant(1)),
			bytecode.NewInstruction(bytecode.OpLoadNone, bytecode.NewRegister(1)),
			bytecode.NewInstruction(bytecode.OpReturn, bytecode.NewRegister(1)),
		},
		CatchTable: []bytecode.Catch{
			{1, 2, 3},
		},
		Metadata: bytecode.Metadata{
			DebugSpans: []file.Span{
				{Start: 0, End: 1},
				{Start: 2, End: 3},
				{Start: 4, End: 5},
				{Start: 6, End: 7},
			},
			Labels: map[int]string{
				3: "end",
			},
		},
	}

	res, err := runPeephole(t, program)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !res.Modified {
		t.Fatalf("expected peephole pass to modify program")
	}
	if len(program.Bytecode) != 3 {
		t.Fatalf("expected 3 instructions, got %d", len(program.Bytecode))
	}

	expectedCatch := []bytecode.Catch{
		{1, 1, 2},
	}
	if !reflect.DeepEqual(program.CatchTable, expectedCatch) {
		t.Fatalf("unexpected catch table: %#v", program.CatchTable)
	}

	if len(program.Metadata.DebugSpans) != 3 {
		t.Fatalf("expected 3 debug spans, got %d", len(program.Metadata.DebugSpans))
	}

	if label, ok := program.Metadata.Labels[2]; !ok || label != "end" {
		t.Fatalf("expected label 'end' at index 2, got %v", program.Metadata.Labels)
	}
}
