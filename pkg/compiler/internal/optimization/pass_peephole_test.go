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

	builder := NewBuilder(program)
	cfg, err := builder.Build()
	if err != nil {
		return nil, err
	}

	pass := NewPeepholePass()
	return pass.Run(&PassContext{Program: program, CFG: cfg})
}

func TestPeephole_AddConstRewriteCases(t *testing.T) {
	testCases := []struct {
		name         string
		program      *bytecode.Program
		wantModified bool
		wantOpcodes  []bytecode.Opcode
	}{
		{
			name: "rewrites when temp register dies and dst equals temp",
			program: &bytecode.Program{
				Constants: []runtime.Value{
					runtime.NewInt(7),
				},
				Bytecode: []bytecode.Instruction{
					bytecode.NewInstruction(bytecode.OpLoadConst, bytecode.NewRegister(2), bytecode.NewConstant(0)),
					bytecode.NewInstruction(bytecode.OpAdd, bytecode.NewRegister(2), bytecode.NewRegister(1), bytecode.NewRegister(2)),
					bytecode.NewInstruction(bytecode.OpReturn, bytecode.NewRegister(2)),
				},
			},
			wantModified: true,
			wantOpcodes: []bytecode.Opcode{
				bytecode.OpAddConst,
				bytecode.OpReturn,
			},
		},
		{
			name: "does not rewrite when temporary is the left operand",
			program: &bytecode.Program{
				Constants: []runtime.Value{
					runtime.NewInt(7),
				},
				Bytecode: []bytecode.Instruction{
					bytecode.NewInstruction(bytecode.OpLoadConst, bytecode.NewRegister(2), bytecode.NewConstant(0)),
					bytecode.NewInstruction(bytecode.OpAdd, bytecode.NewRegister(3), bytecode.NewRegister(2), bytecode.NewRegister(1)),
					bytecode.NewInstruction(bytecode.OpReturn, bytecode.NewRegister(3)),
				},
			},
			wantModified: false,
			wantOpcodes: []bytecode.Opcode{
				bytecode.OpLoadConst,
				bytecode.OpAdd,
				bytecode.OpReturn,
			},
		},
		{
			name: "does not rewrite when load is jump target",
			program: &bytecode.Program{
				Constants: []runtime.Value{
					runtime.NewInt(5),
				},
				Bytecode: []bytecode.Instruction{
					bytecode.NewInstruction(bytecode.OpJump, bytecode.Operand(1)),
					bytecode.NewInstruction(bytecode.OpLoadConst, bytecode.NewRegister(2), bytecode.NewConstant(0)),
					bytecode.NewInstruction(bytecode.OpAdd, bytecode.NewRegister(3), bytecode.NewRegister(1), bytecode.NewRegister(2)),
					bytecode.NewInstruction(bytecode.OpReturn, bytecode.NewRegister(3)),
				},
			},
			wantModified: false,
			wantOpcodes: []bytecode.Opcode{
				bytecode.OpJump,
				bytecode.OpLoadConst,
				bytecode.OpAdd,
				bytecode.OpReturn,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			res, err := runPeephole(t, tc.program)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if res.Modified != tc.wantModified {
				t.Fatalf("unexpected modified flag: got %v, want %v", res.Modified, tc.wantModified)
			}

			if len(tc.program.Bytecode) != len(tc.wantOpcodes) {
				t.Fatalf("unexpected instruction count: got %d, want %d", len(tc.program.Bytecode), len(tc.wantOpcodes))
			}

			for i, op := range tc.wantOpcodes {
				if tc.program.Bytecode[i].Opcode != op {
					t.Fatalf("unexpected opcode at %d: got %s, want %s", i, tc.program.Bytecode[i].Opcode, op)
				}
			}
		})
	}
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

func TestPeephole_RewritesEqJumpToJumpIfNe(t *testing.T) {
	program := &bytecode.Program{
		Bytecode: []bytecode.Instruction{
			bytecode.NewInstruction(bytecode.OpEq, bytecode.NewRegister(3), bytecode.NewRegister(1), bytecode.NewRegister(2)),
			bytecode.NewInstruction(bytecode.OpJumpIfFalse, bytecode.Operand(3), bytecode.NewRegister(3)),
			bytecode.NewInstruction(bytecode.OpReturn, bytecode.NewRegister(1)),
			bytecode.NewInstruction(bytecode.OpReturn, bytecode.NewRegister(2)),
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

	inst := program.Bytecode[0]
	if inst.Opcode != bytecode.OpJumpIfNe {
		t.Fatalf("expected first instruction to be JMPNE, got %s", inst.Opcode)
	}
	if inst.Operands[0] != bytecode.Operand(2) {
		t.Fatalf("unexpected jump target: got %d, want %d", inst.Operands[0], 2)
	}
	if inst.Operands[1] != bytecode.NewRegister(1) || inst.Operands[2] != bytecode.NewRegister(2) {
		t.Fatalf("unexpected operands: got %v %v", inst.Operands[1], inst.Operands[2])
	}
}

func TestPeephole_RewritesEqConstJumpToJumpIfNeConst(t *testing.T) {
	program := &bytecode.Program{
		Constants: []runtime.Value{
			runtime.NewInt(7),
		},
		Bytecode: []bytecode.Instruction{
			bytecode.NewInstruction(bytecode.OpLoadConst, bytecode.NewRegister(3), bytecode.NewConstant(0)),
			bytecode.NewInstruction(bytecode.OpEq, bytecode.NewRegister(4), bytecode.NewRegister(1), bytecode.NewRegister(3)),
			bytecode.NewInstruction(bytecode.OpJumpIfFalse, bytecode.Operand(4), bytecode.NewRegister(4)),
			bytecode.NewInstruction(bytecode.OpReturn, bytecode.NewRegister(1)),
			bytecode.NewInstruction(bytecode.OpReturn, bytecode.NewRegister(2)),
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

	inst := program.Bytecode[0]
	if inst.Opcode != bytecode.OpJumpIfNeConst {
		t.Fatalf("expected first instruction to be JMPNEC, got %s", inst.Opcode)
	}
	if inst.Operands[0] != bytecode.Operand(2) {
		t.Fatalf("unexpected jump target: got %d, want %d", inst.Operands[0], 2)
	}
	if inst.Operands[1] != bytecode.NewRegister(1) || inst.Operands[2] != bytecode.NewConstant(0) {
		t.Fatalf("unexpected operands: got %v %v", inst.Operands[1], inst.Operands[2])
	}
}

func TestPeephole_RewritesEqJumpTrueToJumpIfEq(t *testing.T) {
	program := &bytecode.Program{
		Bytecode: []bytecode.Instruction{
			bytecode.NewInstruction(bytecode.OpEq, bytecode.NewRegister(3), bytecode.NewRegister(1), bytecode.NewRegister(2)),
			bytecode.NewInstruction(bytecode.OpJumpIfTrue, bytecode.Operand(3), bytecode.NewRegister(3)),
			bytecode.NewInstruction(bytecode.OpReturn, bytecode.NewRegister(1)),
			bytecode.NewInstruction(bytecode.OpReturn, bytecode.NewRegister(2)),
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

	inst := program.Bytecode[0]
	if inst.Opcode != bytecode.OpJumpIfEq {
		t.Fatalf("expected first instruction to be JMPEQ, got %s", inst.Opcode)
	}
	if inst.Operands[0] != bytecode.Operand(2) {
		t.Fatalf("unexpected jump target: got %d, want %d", inst.Operands[0], 2)
	}
	if inst.Operands[1] != bytecode.NewRegister(1) || inst.Operands[2] != bytecode.NewRegister(2) {
		t.Fatalf("unexpected operands: got %v %v", inst.Operands[1], inst.Operands[2])
	}
}

func TestPeephole_RewritesNeJumpFalseToJumpIfEq(t *testing.T) {
	program := &bytecode.Program{
		Bytecode: []bytecode.Instruction{
			bytecode.NewInstruction(bytecode.OpNe, bytecode.NewRegister(3), bytecode.NewRegister(1), bytecode.NewRegister(2)),
			bytecode.NewInstruction(bytecode.OpJumpIfFalse, bytecode.Operand(3), bytecode.NewRegister(3)),
			bytecode.NewInstruction(bytecode.OpReturn, bytecode.NewRegister(1)),
			bytecode.NewInstruction(bytecode.OpReturn, bytecode.NewRegister(2)),
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

	inst := program.Bytecode[0]
	if inst.Opcode != bytecode.OpJumpIfEq {
		t.Fatalf("expected first instruction to be JMPEQ, got %s", inst.Opcode)
	}
	if inst.Operands[0] != bytecode.Operand(2) {
		t.Fatalf("unexpected jump target: got %d, want %d", inst.Operands[0], 2)
	}
	if inst.Operands[1] != bytecode.NewRegister(1) || inst.Operands[2] != bytecode.NewRegister(2) {
		t.Fatalf("unexpected operands: got %v %v", inst.Operands[1], inst.Operands[2])
	}
}

func TestPeephole_RewritesNeJumpTrueToJumpIfNe(t *testing.T) {
	program := &bytecode.Program{
		Bytecode: []bytecode.Instruction{
			bytecode.NewInstruction(bytecode.OpNe, bytecode.NewRegister(3), bytecode.NewRegister(1), bytecode.NewRegister(2)),
			bytecode.NewInstruction(bytecode.OpJumpIfTrue, bytecode.Operand(3), bytecode.NewRegister(3)),
			bytecode.NewInstruction(bytecode.OpReturn, bytecode.NewRegister(1)),
			bytecode.NewInstruction(bytecode.OpReturn, bytecode.NewRegister(2)),
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

	inst := program.Bytecode[0]
	if inst.Opcode != bytecode.OpJumpIfNe {
		t.Fatalf("expected first instruction to be JMPNE, got %s", inst.Opcode)
	}
	if inst.Operands[0] != bytecode.Operand(2) {
		t.Fatalf("unexpected jump target: got %d, want %d", inst.Operands[0], 2)
	}
	if inst.Operands[1] != bytecode.NewRegister(1) || inst.Operands[2] != bytecode.NewRegister(2) {
		t.Fatalf("unexpected operands: got %v %v", inst.Operands[1], inst.Operands[2])
	}
}

func TestPeephole_RewritesEqJumpTrueConstToJumpIfEqConst(t *testing.T) {
	program := &bytecode.Program{
		Constants: []runtime.Value{
			runtime.NewInt(9),
		},
		Bytecode: []bytecode.Instruction{
			bytecode.NewInstruction(bytecode.OpLoadConst, bytecode.NewRegister(3), bytecode.NewConstant(0)),
			bytecode.NewInstruction(bytecode.OpEq, bytecode.NewRegister(4), bytecode.NewRegister(1), bytecode.NewRegister(3)),
			bytecode.NewInstruction(bytecode.OpJumpIfTrue, bytecode.Operand(4), bytecode.NewRegister(4)),
			bytecode.NewInstruction(bytecode.OpReturn, bytecode.NewRegister(1)),
			bytecode.NewInstruction(bytecode.OpReturn, bytecode.NewRegister(2)),
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

	inst := program.Bytecode[0]
	if inst.Opcode != bytecode.OpJumpIfEqConst {
		t.Fatalf("expected first instruction to be JMPEQC, got %s", inst.Opcode)
	}
	if inst.Operands[0] != bytecode.Operand(2) {
		t.Fatalf("unexpected jump target: got %d, want %d", inst.Operands[0], 2)
	}
	if inst.Operands[1] != bytecode.NewRegister(1) || inst.Operands[2] != bytecode.NewConstant(0) {
		t.Fatalf("unexpected operands: got %v %v", inst.Operands[1], inst.Operands[2])
	}
}

func TestPeephole_RewritesNeJumpFalseConstToJumpIfEqConst(t *testing.T) {
	program := &bytecode.Program{
		Constants: []runtime.Value{
			runtime.NewInt(9),
		},
		Bytecode: []bytecode.Instruction{
			bytecode.NewInstruction(bytecode.OpLoadConst, bytecode.NewRegister(3), bytecode.NewConstant(0)),
			bytecode.NewInstruction(bytecode.OpNe, bytecode.NewRegister(4), bytecode.NewRegister(1), bytecode.NewRegister(3)),
			bytecode.NewInstruction(bytecode.OpJumpIfFalse, bytecode.Operand(4), bytecode.NewRegister(4)),
			bytecode.NewInstruction(bytecode.OpReturn, bytecode.NewRegister(1)),
			bytecode.NewInstruction(bytecode.OpReturn, bytecode.NewRegister(2)),
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

	inst := program.Bytecode[0]
	if inst.Opcode != bytecode.OpJumpIfEqConst {
		t.Fatalf("expected first instruction to be JMPEQC, got %s", inst.Opcode)
	}
	if inst.Operands[0] != bytecode.Operand(2) {
		t.Fatalf("unexpected jump target: got %d, want %d", inst.Operands[0], 2)
	}
	if inst.Operands[1] != bytecode.NewRegister(1) || inst.Operands[2] != bytecode.NewConstant(0) {
		t.Fatalf("unexpected operands: got %v %v", inst.Operands[1], inst.Operands[2])
	}
}

func TestPeephole_DoesNotRewriteEqWhenResultIsLive(t *testing.T) {
	program := &bytecode.Program{
		Bytecode: []bytecode.Instruction{
			bytecode.NewInstruction(bytecode.OpEq, bytecode.NewRegister(3), bytecode.NewRegister(1), bytecode.NewRegister(2)),
			bytecode.NewInstruction(bytecode.OpJumpIfFalse, bytecode.Operand(3), bytecode.NewRegister(3)),
			bytecode.NewInstruction(bytecode.OpReturn, bytecode.NewRegister(3)),
			bytecode.NewInstruction(bytecode.OpReturn, bytecode.NewRegister(2)),
		},
	}

	res, err := runPeephole(t, program)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Modified {
		t.Fatalf("expected peephole pass to keep program unchanged")
	}
	if program.Bytecode[0].Opcode != bytecode.OpEq {
		t.Fatalf("expected first instruction to remain EQ, got %s", program.Bytecode[0].Opcode)
	}
}

func TestPeephole_DoesNotRewriteEqWhenTargeted(t *testing.T) {
	program := &bytecode.Program{
		Bytecode: []bytecode.Instruction{
			bytecode.NewInstruction(bytecode.OpJump, bytecode.Operand(1)),
			bytecode.NewInstruction(bytecode.OpEq, bytecode.NewRegister(3), bytecode.NewRegister(1), bytecode.NewRegister(2)),
			bytecode.NewInstruction(bytecode.OpJumpIfFalse, bytecode.Operand(4), bytecode.NewRegister(3)),
			bytecode.NewInstruction(bytecode.OpReturn, bytecode.NewRegister(1)),
			bytecode.NewInstruction(bytecode.OpReturn, bytecode.NewRegister(2)),
		},
	}

	res, err := runPeephole(t, program)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Modified {
		t.Fatalf("expected peephole pass to keep program unchanged")
	}
	if program.Bytecode[1].Opcode != bytecode.OpEq {
		t.Fatalf("expected EQ to remain, got %s", program.Bytecode[1].Opcode)
	}
}

func TestPeephole_DoesNotRewriteWhenLoadConstIsTarget(t *testing.T) {
	program := &bytecode.Program{
		Constants: []runtime.Value{
			runtime.NewInt(3),
		},
		Bytecode: []bytecode.Instruction{
			bytecode.NewInstruction(bytecode.OpJump, bytecode.Operand(1)),
			bytecode.NewInstruction(bytecode.OpLoadConst, bytecode.NewRegister(3), bytecode.NewConstant(0)),
			bytecode.NewInstruction(bytecode.OpEq, bytecode.NewRegister(4), bytecode.NewRegister(1), bytecode.NewRegister(3)),
			bytecode.NewInstruction(bytecode.OpJumpIfFalse, bytecode.Operand(5), bytecode.NewRegister(4)),
			bytecode.NewInstruction(bytecode.OpReturn, bytecode.NewRegister(1)),
			bytecode.NewInstruction(bytecode.OpReturn, bytecode.NewRegister(2)),
		},
	}

	res, err := runPeephole(t, program)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Modified {
		t.Fatalf("expected peephole pass to keep program unchanged")
	}
	if program.Bytecode[1].Opcode != bytecode.OpLoadConst {
		t.Fatalf("expected LOADC to remain, got %s", program.Bytecode[1].Opcode)
	}
}

func TestPeephole_SkipsAddConstRewriteWhenValueLiveAcrossJump(t *testing.T) {
	program := &bytecode.Program{
		Constants: []runtime.Value{
			runtime.NewInt(2),
		},
		Bytecode: []bytecode.Instruction{
			bytecode.NewInstruction(bytecode.OpLoadConst, bytecode.NewRegister(2), bytecode.NewConstant(0)),
			bytecode.NewInstruction(bytecode.OpAdd, bytecode.NewRegister(1), bytecode.NewRegister(1), bytecode.NewRegister(2)),
			bytecode.NewInstruction(bytecode.OpJump, bytecode.Operand(5)),
			bytecode.NewInstruction(bytecode.OpLoadConst, bytecode.NewRegister(2), bytecode.NewConstant(0)),
			bytecode.NewInstruction(bytecode.OpReturn, bytecode.NewRegister(1)),
			bytecode.NewInstruction(bytecode.OpReturn, bytecode.NewRegister(2)),
		},
	}

	res, err := runPeephole(t, program)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Modified {
		t.Fatalf("expected peephole pass to keep live temp register load")
	}
	if len(program.Bytecode) != 6 {
		t.Fatalf("expected 6 instructions, got %d", len(program.Bytecode))
	}
	if program.Bytecode[0].Opcode != bytecode.OpLoadConst {
		t.Fatalf("expected first instruction to remain LOADC, got %s", program.Bytecode[0].Opcode)
	}
	if program.Bytecode[1].Opcode != bytecode.OpAdd {
		t.Fatalf("expected second instruction to remain ADD, got %s", program.Bytecode[1].Opcode)
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
