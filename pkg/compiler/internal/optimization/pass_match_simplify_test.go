package optimization

import (
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

func TestMatchSimplify_ConstantScrutineeRemovesArms(t *testing.T) {
	prog := &bytecode.Program{
		Registers: 3,
		Constants: []runtime.Value{
			runtime.NewInt(1),  // 0 scrutinee
			runtime.NewInt(1),  // 1 pattern 1
			runtime.NewInt(2),  // 2 pattern 2
			runtime.NewInt(10), // 3 result 1
			runtime.NewInt(20), // 4 result 2
			runtime.NewInt(30), // 5 default
		},
		Bytecode: []bytecode.Instruction{
			bytecode.NewInstruction(bytecode.OpLoadConst, 1, bytecode.NewConstant(0)),
			bytecode.NewInstruction(bytecode.OpJumpIfNeConst, bytecode.Operand(4), bytecode.NewRegister(1), bytecode.NewConstant(1)),
			bytecode.NewInstruction(bytecode.OpLoadConst, 2, bytecode.NewConstant(3)),
			bytecode.NewInstruction(bytecode.OpJump, bytecode.Operand(7)),
			bytecode.NewInstruction(bytecode.OpJumpIfNeConst, bytecode.Operand(6), bytecode.NewRegister(1), bytecode.NewConstant(2)),
			bytecode.NewInstruction(bytecode.OpLoadConst, 2, bytecode.NewConstant(4)),
			bytecode.NewInstruction(bytecode.OpLoadConst, 2, bytecode.NewConstant(5)),
			bytecode.NewInstruction(bytecode.OpReturn, bytecode.NewRegister(2)),
		},
		Metadata: bytecode.Metadata{
			Labels: map[int]string{
				1: "match.0.start",
				4: "match.0.next.0",
				6: "match.0.next.1",
				7: "match.0.end",
			},
		},
	}

	builder := NewBuilder(prog)
	cfg, err := builder.Build()
	if err != nil {
		t.Fatalf("failed to build CFG: %v", err)
	}

	pass := NewMatchSimplifyPass()
	result, err := pass.Run(&PassContext{
		Program:  prog,
		CFG:      cfg,
		Metadata: map[string]any{},
	})
	if err != nil {
		t.Fatalf("pass failed: %v", err)
	}
	if !result.Modified {
		t.Fatalf("expected pass to modify program")
	}
	if len(prog.Bytecode) >= 8 {
		t.Fatalf("expected bytecode to shrink, got len=%d", len(prog.Bytecode))
	}

	start := labelIndex(prog.Metadata.Labels, "match.0.start")
	if start < 0 || start >= len(prog.Bytecode) {
		t.Fatalf("missing match start label")
	}
	if prog.Bytecode[start].Opcode != bytecode.OpMove {
		t.Fatalf("expected first arm compare to be replaced by MOVE, got %s", prog.Bytecode[start].Opcode)
	}
	for _, inst := range prog.Bytecode {
		if inst.Opcode == bytecode.OpJumpIfNeConst {
			t.Fatalf("expected no JumpIfNeConst after simplification")
		}
	}
}

func TestMatchSimplify_NoOpOnGuardedArm(t *testing.T) {
	prog := &bytecode.Program{
		Registers: 4,
		Constants: []runtime.Value{
			runtime.NewInt(1),  // 0 scrutinee
			runtime.NewInt(1),  // 1 pattern 1
			runtime.NewInt(10), // 2 result 1
			runtime.NewInt(30), // 3 default
		},
		Bytecode: []bytecode.Instruction{
			bytecode.NewInstruction(bytecode.OpLoadConst, 1, bytecode.NewConstant(0)),
			bytecode.NewInstruction(bytecode.OpJumpIfNeConst, bytecode.Operand(5), bytecode.NewRegister(1), bytecode.NewConstant(1)),
			bytecode.NewInstruction(bytecode.OpJumpIfFalse, bytecode.Operand(5), bytecode.NewRegister(3)),
			bytecode.NewInstruction(bytecode.OpLoadConst, 2, bytecode.NewConstant(2)),
			bytecode.NewInstruction(bytecode.OpJump, bytecode.Operand(6)),
			bytecode.NewInstruction(bytecode.OpLoadConst, 2, bytecode.NewConstant(3)),
			bytecode.NewInstruction(bytecode.OpReturn, bytecode.NewRegister(2)),
		},
		Metadata: bytecode.Metadata{
			Labels: map[int]string{
				1: "match.1.start",
				5: "match.1.next.0",
				6: "match.1.end",
			},
		},
	}

	builder := NewBuilder(prog)
	cfg, err := builder.Build()
	if err != nil {
		t.Fatalf("failed to build CFG: %v", err)
	}

	pass := NewMatchSimplifyPass()
	result, err := pass.Run(&PassContext{
		Program:  prog,
		CFG:      cfg,
		Metadata: map[string]any{},
	})
	if err != nil {
		t.Fatalf("pass failed: %v", err)
	}
	if result.Modified {
		t.Fatalf("expected pass to be no-op for guarded arm")
	}
	if prog.Bytecode[1].Opcode != bytecode.OpJumpIfNeConst {
		t.Fatalf("expected guard arm to remain unchanged")
	}
}

func labelIndex(labels map[int]string, name string) int {
	for idx, lbl := range labels {
		if lbl == name {
			return idx
		}
	}
	return -1
}
