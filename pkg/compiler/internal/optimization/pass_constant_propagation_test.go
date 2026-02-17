package optimization

import (
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

func runConstantPropagation(t *testing.T, program *bytecode.Program) (*PassResult, error) {
	t.Helper()

	builder := NewBuilder(program)
	cfg, err := builder.Build()
	if err != nil {
		return nil, err
	}

	pass := NewConstantPropagationPass()

	return pass.Run(&PassContext{
		Program: program,
		CFG:     cfg,
	})
}

func TestConstantPropagation_FoldsConcat(t *testing.T) {
	program := &bytecode.Program{
		Constants: []runtime.Value{
			runtime.NewString("sum="),
			runtime.NewInt(3),
		},
		Bytecode: []bytecode.Instruction{
			bytecode.NewInstruction(bytecode.OpLoadConst, bytecode.NewRegister(1), bytecode.NewConstant(0)),
			bytecode.NewInstruction(bytecode.OpLoadConst, bytecode.NewRegister(2), bytecode.NewConstant(1)),
			bytecode.NewInstruction(bytecode.OpConcat, bytecode.NewRegister(3), bytecode.NewRegister(1), bytecode.Operand(2)),
			bytecode.NewInstruction(bytecode.OpReturn, bytecode.NewRegister(3)),
		},
	}

	res, err := runConstantPropagation(t, program)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !res.Modified {
		t.Fatalf("expected pass to modify program")
	}

	inst := program.Bytecode[2]
	if inst.Opcode != bytecode.OpLoadConst {
		t.Fatalf("expected CONCAT to fold into LOADC, got %s", inst.Opcode)
	}
	if !inst.Operands[1].IsConstant() {
		t.Fatalf("expected LOADC to use a constant operand, got %s", inst.Operands[1])
	}

	val := program.Constants[inst.Operands[1].Constant()]
	got, ok := val.(runtime.String)
	if !ok {
		t.Fatalf("expected folded value to be a string constant, got %T", val)
	}
	if got != runtime.NewString("sum=3") {
		t.Fatalf("unexpected folded value: %q", got)
	}
}

func TestConstantPropagation_FoldsEmptyConcat(t *testing.T) {
	program := &bytecode.Program{
		Bytecode: []bytecode.Instruction{
			bytecode.NewInstruction(bytecode.OpConcat, bytecode.NewRegister(1), bytecode.NewRegister(1), bytecode.Operand(0)),
			bytecode.NewInstruction(bytecode.OpReturn, bytecode.NewRegister(1)),
		},
	}

	res, err := runConstantPropagation(t, program)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !res.Modified {
		t.Fatalf("expected pass to modify program")
	}

	inst := program.Bytecode[0]
	if inst.Opcode != bytecode.OpLoadConst {
		t.Fatalf("expected CONCAT to fold into LOADC, got %s", inst.Opcode)
	}
	if !inst.Operands[1].IsConstant() {
		t.Fatalf("expected LOADC to use a constant operand, got %s", inst.Operands[1])
	}

	val := program.Constants[inst.Operands[1].Constant()]
	got, ok := val.(runtime.String)
	if !ok {
		t.Fatalf("expected folded value to be a string constant, got %T", val)
	}
	if got != runtime.EmptyString {
		t.Fatalf("expected empty string constant, got %q", got)
	}
}

func TestConstantPropagation_DoesNotFoldConcatWithUnknownInputs(t *testing.T) {
	program := &bytecode.Program{
		Constants: []runtime.Value{
			runtime.NewString("sum="),
		},
		Bytecode: []bytecode.Instruction{
			bytecode.NewInstruction(bytecode.OpLoadConst, bytecode.NewRegister(1), bytecode.NewConstant(0)),
			bytecode.NewInstruction(bytecode.OpMove, bytecode.NewRegister(2), bytecode.NewRegister(3)),
			bytecode.NewInstruction(bytecode.OpConcat, bytecode.NewRegister(4), bytecode.NewRegister(1), bytecode.Operand(2)),
			bytecode.NewInstruction(bytecode.OpReturn, bytecode.NewRegister(4)),
		},
	}

	res, err := runConstantPropagation(t, program)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Modified {
		t.Fatalf("expected pass to keep CONCAT when input is unknown")
	}
	if program.Bytecode[2].Opcode != bytecode.OpConcat {
		t.Fatalf("expected CONCAT to remain unchanged, got %s", program.Bytecode[2].Opcode)
	}
}
