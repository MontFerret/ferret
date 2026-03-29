package optimization

import (
	"reflect"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/source"
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

type peepholeOpcodeCase struct {
	program      func() *bytecode.Program
	name         string
	wantOpcodes  []bytecode.Opcode
	wantModified bool
}

type peepholeJumpRewriteCase struct {
	program        func() *bytecode.Program
	name           string
	wantJumpTarget bytecode.Operand
	wantLeft       bytecode.Operand
	wantRight      bytecode.Operand
	wantOpcode     bytecode.Opcode
}

var addConstRewriteCases = []peepholeOpcodeCase{
	{
		name: "rewrites when temp register dies and dst equals temp",
		program: func() *bytecode.Program {
			return &bytecode.Program{
				Constants: []runtime.Value{
					runtime.NewInt(7),
				},
				Bytecode: []bytecode.Instruction{
					bytecode.NewInstruction(bytecode.OpLoadConst, bytecode.NewRegister(2), bytecode.NewConstant(0)),
					bytecode.NewInstruction(bytecode.OpAdd, bytecode.NewRegister(2), bytecode.NewRegister(1), bytecode.NewRegister(2)),
					bytecode.NewInstruction(bytecode.OpReturn, bytecode.NewRegister(2)),
				},
			}
		},
		wantModified: true,
		wantOpcodes: []bytecode.Opcode{
			bytecode.OpAddConst,
			bytecode.OpReturn,
		},
	},
	{
		name: "does not rewrite when temporary is the left operand",
		program: func() *bytecode.Program {
			return &bytecode.Program{
				Constants: []runtime.Value{
					runtime.NewInt(7),
				},
				Bytecode: []bytecode.Instruction{
					bytecode.NewInstruction(bytecode.OpLoadConst, bytecode.NewRegister(2), bytecode.NewConstant(0)),
					bytecode.NewInstruction(bytecode.OpAdd, bytecode.NewRegister(3), bytecode.NewRegister(2), bytecode.NewRegister(1)),
					bytecode.NewInstruction(bytecode.OpReturn, bytecode.NewRegister(3)),
				},
			}
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
		program: func() *bytecode.Program {
			return &bytecode.Program{
				Constants: []runtime.Value{
					runtime.NewInt(5),
				},
				Bytecode: []bytecode.Instruction{
					bytecode.NewInstruction(bytecode.OpJump, bytecode.Operand(1)),
					bytecode.NewInstruction(bytecode.OpLoadConst, bytecode.NewRegister(2), bytecode.NewConstant(0)),
					bytecode.NewInstruction(bytecode.OpAdd, bytecode.NewRegister(3), bytecode.NewRegister(1), bytecode.NewRegister(2)),
					bytecode.NewInstruction(bytecode.OpReturn, bytecode.NewRegister(3)),
				},
			}
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

var redundantDefRemovalCases = []peepholeOpcodeCase{
	{
		name: "removes redundant load const",
		program: func() *bytecode.Program {
			return &bytecode.Program{
				Constants: []runtime.Value{
					runtime.NewString("a"),
				},
				Bytecode: []bytecode.Instruction{
					bytecode.NewInstruction(bytecode.OpLoadConst, bytecode.NewRegister(1), bytecode.NewConstant(0)),
					bytecode.NewInstruction(bytecode.OpLoadConst, bytecode.NewRegister(1), bytecode.NewConstant(0)),
					bytecode.NewInstruction(bytecode.OpReturn, bytecode.NewRegister(1)),
				},
			}
		},
		wantModified: true,
		wantOpcodes: []bytecode.Opcode{
			bytecode.OpLoadConst,
			bytecode.OpReturn,
		},
	},
	{
		name: "removes self move",
		program: func() *bytecode.Program {
			return &bytecode.Program{
				Constants: []runtime.Value{
					runtime.NewString("a"),
				},
				Bytecode: []bytecode.Instruction{
					bytecode.NewInstruction(bytecode.OpLoadConst, bytecode.NewRegister(1), bytecode.NewConstant(0)),
					bytecode.NewInstruction(bytecode.OpMove, bytecode.NewRegister(1), bytecode.NewRegister(1)),
					bytecode.NewInstruction(bytecode.OpReturn, bytecode.NewRegister(1)),
				},
			}
		},
		wantModified: true,
		wantOpcodes: []bytecode.Opcode{
			bytecode.OpLoadConst,
			bytecode.OpReturn,
		},
	},
}

var comparisonJumpRewriteCases = []peepholeJumpRewriteCase{
	{
		name: "eq + jump false -> jump if ne",
		program: func() *bytecode.Program {
			return &bytecode.Program{
				Bytecode: []bytecode.Instruction{
					bytecode.NewInstruction(bytecode.OpEq, bytecode.NewRegister(3), bytecode.NewRegister(1), bytecode.NewRegister(2)),
					bytecode.NewInstruction(bytecode.OpJumpIfFalse, bytecode.Operand(3), bytecode.NewRegister(3)),
					bytecode.NewInstruction(bytecode.OpReturn, bytecode.NewRegister(1)),
					bytecode.NewInstruction(bytecode.OpReturn, bytecode.NewRegister(2)),
				},
			}
		},
		wantOpcode:     bytecode.OpJumpIfNe,
		wantJumpTarget: bytecode.Operand(2),
		wantLeft:       bytecode.NewRegister(1),
		wantRight:      bytecode.NewRegister(2),
	},
	{
		name: "eq const + jump false -> jump if ne const",
		program: func() *bytecode.Program {
			return &bytecode.Program{
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
		},
		wantOpcode:     bytecode.OpJumpIfNeConst,
		wantJumpTarget: bytecode.Operand(2),
		wantLeft:       bytecode.NewRegister(1),
		wantRight:      bytecode.NewConstant(0),
	},
	{
		name: "eq + jump true -> jump if eq",
		program: func() *bytecode.Program {
			return &bytecode.Program{
				Bytecode: []bytecode.Instruction{
					bytecode.NewInstruction(bytecode.OpEq, bytecode.NewRegister(3), bytecode.NewRegister(1), bytecode.NewRegister(2)),
					bytecode.NewInstruction(bytecode.OpJumpIfTrue, bytecode.Operand(3), bytecode.NewRegister(3)),
					bytecode.NewInstruction(bytecode.OpReturn, bytecode.NewRegister(1)),
					bytecode.NewInstruction(bytecode.OpReturn, bytecode.NewRegister(2)),
				},
			}
		},
		wantOpcode:     bytecode.OpJumpIfEq,
		wantJumpTarget: bytecode.Operand(2),
		wantLeft:       bytecode.NewRegister(1),
		wantRight:      bytecode.NewRegister(2),
	},
	{
		name: "ne + jump false -> jump if eq",
		program: func() *bytecode.Program {
			return &bytecode.Program{
				Bytecode: []bytecode.Instruction{
					bytecode.NewInstruction(bytecode.OpNe, bytecode.NewRegister(3), bytecode.NewRegister(1), bytecode.NewRegister(2)),
					bytecode.NewInstruction(bytecode.OpJumpIfFalse, bytecode.Operand(3), bytecode.NewRegister(3)),
					bytecode.NewInstruction(bytecode.OpReturn, bytecode.NewRegister(1)),
					bytecode.NewInstruction(bytecode.OpReturn, bytecode.NewRegister(2)),
				},
			}
		},
		wantOpcode:     bytecode.OpJumpIfEq,
		wantJumpTarget: bytecode.Operand(2),
		wantLeft:       bytecode.NewRegister(1),
		wantRight:      bytecode.NewRegister(2),
	},
	{
		name: "ne + jump true -> jump if ne",
		program: func() *bytecode.Program {
			return &bytecode.Program{
				Bytecode: []bytecode.Instruction{
					bytecode.NewInstruction(bytecode.OpNe, bytecode.NewRegister(3), bytecode.NewRegister(1), bytecode.NewRegister(2)),
					bytecode.NewInstruction(bytecode.OpJumpIfTrue, bytecode.Operand(3), bytecode.NewRegister(3)),
					bytecode.NewInstruction(bytecode.OpReturn, bytecode.NewRegister(1)),
					bytecode.NewInstruction(bytecode.OpReturn, bytecode.NewRegister(2)),
				},
			}
		},
		wantOpcode:     bytecode.OpJumpIfNe,
		wantJumpTarget: bytecode.Operand(2),
		wantLeft:       bytecode.NewRegister(1),
		wantRight:      bytecode.NewRegister(2),
	},
	{
		name: "eq const + jump true -> jump if eq const",
		program: func() *bytecode.Program {
			return &bytecode.Program{
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
		},
		wantOpcode:     bytecode.OpJumpIfEqConst,
		wantJumpTarget: bytecode.Operand(2),
		wantLeft:       bytecode.NewRegister(1),
		wantRight:      bytecode.NewConstant(0),
	},
	{
		name: "ne const + jump false -> jump if eq const",
		program: func() *bytecode.Program {
			return &bytecode.Program{
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
		},
		wantOpcode:     bytecode.OpJumpIfEqConst,
		wantJumpTarget: bytecode.Operand(2),
		wantLeft:       bytecode.NewRegister(1),
		wantRight:      bytecode.NewConstant(0),
	},
}

func runPeepholeOpcodeCase(t *testing.T, tc peepholeOpcodeCase) {
	t.Helper()

	program := tc.program()
	res, err := runPeephole(t, program)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if res.Modified != tc.wantModified {
		t.Fatalf("unexpected modified flag: got %v, want %v", res.Modified, tc.wantModified)
	}

	if len(program.Bytecode) != len(tc.wantOpcodes) {
		t.Fatalf("unexpected instruction count: got %d, want %d", len(program.Bytecode), len(tc.wantOpcodes))
	}

	for i, op := range tc.wantOpcodes {
		if program.Bytecode[i].Opcode != op {
			t.Fatalf("unexpected opcode at %d: got %s, want %s", i, program.Bytecode[i].Opcode, op)
		}
	}
}

func runPeepholeJumpRewriteCase(t *testing.T, tc peepholeJumpRewriteCase) {
	t.Helper()

	program := tc.program()
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
	if inst.Opcode != tc.wantOpcode {
		t.Fatalf("unexpected first opcode: got %s, want %s", inst.Opcode, tc.wantOpcode)
	}
	if inst.Operands[0] != tc.wantJumpTarget {
		t.Fatalf("unexpected jump target: got %d, want %d", inst.Operands[0], tc.wantJumpTarget)
	}
	if inst.Operands[1] != tc.wantLeft || inst.Operands[2] != tc.wantRight {
		t.Fatalf("unexpected operands: got %v %v", inst.Operands[1], inst.Operands[2])
	}
}

func TestPeephole_AddConstRewriteCases(t *testing.T) {
	for _, tc := range addConstRewriteCases {
		t.Run(tc.name, func(t *testing.T) {
			runPeepholeOpcodeCase(t, tc)
		})
	}
}

func TestPeephole_RemovesRedundantDefs(t *testing.T) {
	for _, tc := range redundantDefRemovalCases {
		t.Run(tc.name, func(t *testing.T) {
			runPeepholeOpcodeCase(t, tc)
		})
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

func TestPeephole_RewritesComparisonJumps(t *testing.T) {
	for _, tc := range comparisonJumpRewriteCases {
		t.Run(tc.name, func(t *testing.T) {
			runPeepholeJumpRewriteCase(t, tc)
		})
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

func TestPeephole_RemovesRedundantPureDefAfterLabel(t *testing.T) {
	prog := &bytecode.Program{
		Registers: 2,
		Constants: []runtime.Value{
			runtime.NewString("default"),
		},
		Bytecode: []bytecode.Instruction{
			bytecode.NewInstruction(bytecode.OpLoadConst, bytecode.NewRegister(1), bytecode.NewConstant(0)),
			bytecode.NewInstruction(bytecode.OpLoadConst, bytecode.NewRegister(1), bytecode.NewConstant(0)),
			bytecode.NewInstruction(bytecode.OpReturn, bytecode.NewRegister(1)),
		},
		Metadata: bytecode.Metadata{
			Labels: map[int]string{
				0: "match.0.next.1",
			},
		},
	}

	builder := NewBuilder(prog)
	cfg, err := builder.Build()
	if err != nil {
		t.Fatalf("failed to build CFG: %v", err)
	}

	pass := NewPeepholePass()
	result, err := pass.Run(&PassContext{
		Program:  prog,
		CFG:      cfg,
		Metadata: map[string]any{},
	})
	if err != nil {
		t.Fatalf("pass failed: %v", err)
	}
	if !result.Modified {
		t.Fatalf("expected peephole to modify program")
	}
	if len(prog.Bytecode) != 2 {
		t.Fatalf("expected redundant LOADC removal, got len=%d", len(prog.Bytecode))
	}
	if prog.Bytecode[0].Opcode != bytecode.OpLoadConst || prog.Bytecode[1].Opcode != bytecode.OpReturn {
		t.Fatalf("unexpected bytecode sequence after peephole")
	}
	if prog.Metadata.Labels[0] != "match.0.next.1" {
		t.Fatalf("expected label to remain on first instruction")
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

func TestPeephole_KeepsInstructionZeroWhenTargetedByCatchJump(t *testing.T) {
	program := &bytecode.Program{
		Bytecode: []bytecode.Instruction{
			bytecode.NewInstruction(bytecode.OpMove, bytecode.NewRegister(1), bytecode.NewRegister(1)),
			bytecode.NewInstruction(bytecode.OpLoadNone, bytecode.NewRegister(1)),
			bytecode.NewInstruction(bytecode.OpReturn, bytecode.NewRegister(1)),
		},
		CatchTable: []bytecode.Catch{
			{1, 1, 0},
		},
	}

	res, err := runPeephole(t, program)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Modified {
		t.Fatalf("expected peephole pass to keep catch jump target at pc 0")
	}
	if len(program.Bytecode) != 3 {
		t.Fatalf("expected 3 instructions, got %d", len(program.Bytecode))
	}
	if program.Bytecode[0].Opcode != bytecode.OpMove {
		t.Fatalf("expected first instruction to remain MOVE, got %s", program.Bytecode[0].Opcode)
	}

	expectedCatch := []bytecode.Catch{
		{1, 1, 0},
	}
	if !reflect.DeepEqual(program.CatchTable, expectedCatch) {
		t.Fatalf("unexpected catch table: %#v", program.CatchTable)
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
			AggregateSelectorSlots: []int{-1, 7, -1, 9},
			MatchFailTargets:       []int{3, -1, -1, -1},
			DebugSpans: []source.Span{
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

	expectedAggregateSelectorSlots := []int{-1, -1, 9}
	if !reflect.DeepEqual(program.Metadata.AggregateSelectorSlots, expectedAggregateSelectorSlots) {
		t.Fatalf("unexpected aggregate selector slots: %#v", program.Metadata.AggregateSelectorSlots)
	}

	expectedMatchFailTargets := []int{2, -1, -1}
	if !reflect.DeepEqual(program.Metadata.MatchFailTargets, expectedMatchFailTargets) {
		t.Fatalf("unexpected match fail targets: %#v", program.Metadata.MatchFailTargets)
	}

	if label, ok := program.Metadata.Labels[2]; !ok || label != "end" {
		t.Fatalf("expected label 'end' at index 2, got %v", program.Metadata.Labels)
	}
}

type peepholeJumpOpcodeMatrixCase struct {
	name         string
	compareOp    bytecode.Opcode
	jumpOp       bytecode.Opcode
	wantJumpOp   bytecode.Opcode
	wantConstOp  bytecode.Opcode
	wantResolved bool
}

var peepholeJumpOpcodeMatrixCases = []peepholeJumpOpcodeMatrixCase{
	{
		name:         "eq with jump false",
		compareOp:    bytecode.OpEq,
		jumpOp:       bytecode.OpJumpIfFalse,
		wantJumpOp:   bytecode.OpJumpIfNe,
		wantConstOp:  bytecode.OpJumpIfNeConst,
		wantResolved: true,
	},
	{
		name:         "eq with jump true",
		compareOp:    bytecode.OpEq,
		jumpOp:       bytecode.OpJumpIfTrue,
		wantJumpOp:   bytecode.OpJumpIfEq,
		wantConstOp:  bytecode.OpJumpIfEqConst,
		wantResolved: true,
	},
	{
		name:         "ne with jump false",
		compareOp:    bytecode.OpNe,
		jumpOp:       bytecode.OpJumpIfFalse,
		wantJumpOp:   bytecode.OpJumpIfEq,
		wantConstOp:  bytecode.OpJumpIfEqConst,
		wantResolved: true,
	},
	{
		name:         "ne with jump true",
		compareOp:    bytecode.OpNe,
		jumpOp:       bytecode.OpJumpIfTrue,
		wantJumpOp:   bytecode.OpJumpIfNe,
		wantConstOp:  bytecode.OpJumpIfNeConst,
		wantResolved: true,
	},
	{
		name:         "unsupported compare opcode",
		compareOp:    bytecode.OpGt,
		jumpOp:       bytecode.OpJumpIfTrue,
		wantResolved: false,
	},
	{
		name:         "unsupported jump opcode",
		compareOp:    bytecode.OpEq,
		jumpOp:       bytecode.OpJump,
		wantResolved: false,
	},
}

func runPeepholeJumpOpcodeMatrixCase(t *testing.T, tc peepholeJumpOpcodeMatrixCase) {
	t.Helper()

	gotJumpOp, gotConstOp, ok := resolveComparisonJumpOpcode(tc.compareOp, tc.jumpOp)
	if ok != tc.wantResolved {
		t.Fatalf("unexpected resolved flag: got %v, want %v", ok, tc.wantResolved)
	}

	if !tc.wantResolved {
		return
	}

	if gotJumpOp != tc.wantJumpOp {
		t.Fatalf("unexpected jump opcode: got %s, want %s", gotJumpOp, tc.wantJumpOp)
	}

	if gotConstOp != tc.wantConstOp {
		t.Fatalf("unexpected const jump opcode: got %s, want %s", gotConstOp, tc.wantConstOp)
	}
}

func TestPeephole_ResolveComparisonJumpOpcode(t *testing.T) {
	for _, tc := range peepholeJumpOpcodeMatrixCases {
		t.Run(tc.name, func(t *testing.T) {
			runPeepholeJumpOpcodeMatrixCase(t, tc)
		})
	}
}
