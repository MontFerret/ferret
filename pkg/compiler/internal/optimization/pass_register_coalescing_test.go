package optimization

import (
	"reflect"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

func TestRegisterCoalescing_CoalescesWhenSrcDies(t *testing.T) {
	program := &bytecode.Program{
		Constants: []runtime.Value{
			runtime.Int(10),
		},
		Registers: 4,
		Bytecode: []bytecode.Instruction{
			bytecode.NewInstruction(bytecode.OpLoadConst, bytecode.NewRegister(1), bytecode.NewConstant(0)),
			bytecode.NewInstruction(bytecode.OpMove, bytecode.NewRegister(2), bytecode.NewRegister(1)),
			bytecode.NewInstruction(bytecode.OpAdd, bytecode.NewRegister(3), bytecode.NewRegister(2), bytecode.NewRegister(2)),
			bytecode.NewInstruction(bytecode.OpReturn, bytecode.NewRegister(3)),
		},
	}

	runCoalescing(t, program)

	expected := []bytecode.Instruction{
		bytecode.NewInstruction(bytecode.OpLoadConst, bytecode.NewRegister(1), bytecode.NewConstant(0)),
		bytecode.NewInstruction(bytecode.OpMove, bytecode.NewRegister(1), bytecode.NewRegister(1)),
		bytecode.NewInstruction(bytecode.OpAdd, bytecode.NewRegister(1), bytecode.NewRegister(1), bytecode.NewRegister(1)),
		bytecode.NewInstruction(bytecode.OpReturn, bytecode.NewRegister(1)),
	}

	assertBytecodeEqual(t, program.Bytecode, expected)
}

func TestRegisterCoalescing_NoCoalesceWhenInterfering(t *testing.T) {
	program := &bytecode.Program{
		Constants: []runtime.Value{
			runtime.Int(1),
		},
		Registers: 4,
		Bytecode: []bytecode.Instruction{
			bytecode.NewInstruction(bytecode.OpLoadConst, bytecode.NewRegister(1), bytecode.NewConstant(0)),
			bytecode.NewInstruction(bytecode.OpMove, bytecode.NewRegister(2), bytecode.NewRegister(1)),
			bytecode.NewInstruction(bytecode.OpIncr, bytecode.NewRegister(2)),
			bytecode.NewInstruction(bytecode.OpAdd, bytecode.NewRegister(3), bytecode.NewRegister(2), bytecode.NewRegister(1)),
			bytecode.NewInstruction(bytecode.OpReturn, bytecode.NewRegister(3)),
		},
	}

	runCoalescing(t, program)

	expected := []bytecode.Instruction{
		bytecode.NewInstruction(bytecode.OpLoadConst, bytecode.NewRegister(1), bytecode.NewConstant(0)),
		bytecode.NewInstruction(bytecode.OpMove, bytecode.NewRegister(2), bytecode.NewRegister(1)),
		bytecode.NewInstruction(bytecode.OpIncr, bytecode.NewRegister(2)),
		bytecode.NewInstruction(bytecode.OpAdd, bytecode.NewRegister(1), bytecode.NewRegister(2), bytecode.NewRegister(1)),
		bytecode.NewInstruction(bytecode.OpReturn, bytecode.NewRegister(1)),
	}

	assertBytecodeEqual(t, program.Bytecode, expected)
}

func TestRegisterCoalescing_NoCoalesceForRangeSensitiveRegs(t *testing.T) {
	program := &bytecode.Program{
		Constants: []runtime.Value{
			runtime.Int(1),
			runtime.Int(2),
		},
		Registers: 6,
		Bytecode: []bytecode.Instruction{
			bytecode.NewInstruction(bytecode.OpLoadConst, bytecode.NewRegister(1), bytecode.NewConstant(0)),
			bytecode.NewInstruction(bytecode.OpLoadConst, bytecode.NewRegister(2), bytecode.NewConstant(1)),
			bytecode.NewInstruction(bytecode.OpMove, bytecode.NewRegister(3), bytecode.NewRegister(1)),
			bytecode.NewInstruction(bytecode.OpHCall, bytecode.NewRegister(4), bytecode.NewRegister(1), bytecode.NewRegister(2)),
			bytecode.NewInstruction(bytecode.OpAdd, bytecode.NewRegister(5), bytecode.NewRegister(3), bytecode.NewRegister(1)),
			bytecode.NewInstruction(bytecode.OpReturn, bytecode.NewRegister(5)),
		},
	}

	runCoalescing(t, program)

	expected := []bytecode.Instruction{
		bytecode.NewInstruction(bytecode.OpLoadConst, bytecode.NewRegister(1), bytecode.NewConstant(0)),
		bytecode.NewInstruction(bytecode.OpLoadConst, bytecode.NewRegister(2), bytecode.NewConstant(1)),
		bytecode.NewInstruction(bytecode.OpMove, bytecode.NewRegister(3), bytecode.NewRegister(1)),
		bytecode.NewInstruction(bytecode.OpHCall, bytecode.NewRegister(4), bytecode.NewRegister(1), bytecode.NewRegister(2)),
		bytecode.NewInstruction(bytecode.OpAdd, bytecode.NewRegister(3), bytecode.NewRegister(3), bytecode.NewRegister(1)),
		bytecode.NewInstruction(bytecode.OpReturn, bytecode.NewRegister(3)),
	}

	assertBytecodeEqual(t, program.Bytecode, expected)
}

func TestRegisterCoalescing_PinsCellHandleRegisters(t *testing.T) {
	program := &bytecode.Program{
		Constants: []runtime.Value{
			runtime.Int(1),
			runtime.Int(2),
		},
		Registers: 8,
		Bytecode: []bytecode.Instruction{
			bytecode.NewInstruction(bytecode.OpLoadConst, bytecode.NewRegister(3), bytecode.NewConstant(0)),
			bytecode.NewInstruction(bytecode.OpMakeCell, bytecode.NewRegister(5), bytecode.NewRegister(3)),
			bytecode.NewInstruction(bytecode.OpLoadConst, bytecode.NewRegister(6), bytecode.NewConstant(1)),
			bytecode.NewInstruction(bytecode.OpStoreCell, bytecode.NewRegister(5), bytecode.NewRegister(6)),
			bytecode.NewInstruction(bytecode.OpLoadCell, bytecode.NewRegister(7), bytecode.NewRegister(5)),
			bytecode.NewInstruction(bytecode.OpReturn, bytecode.NewRegister(7)),
		},
	}

	runCoalescing(t, program)

	expected := []bytecode.Instruction{
		bytecode.NewInstruction(bytecode.OpLoadConst, bytecode.NewRegister(1), bytecode.NewConstant(0)),
		bytecode.NewInstruction(bytecode.OpMakeCell, bytecode.NewRegister(5), bytecode.NewRegister(1)),
		bytecode.NewInstruction(bytecode.OpLoadConst, bytecode.NewRegister(1), bytecode.NewConstant(1)),
		bytecode.NewInstruction(bytecode.OpStoreCell, bytecode.NewRegister(5), bytecode.NewRegister(1)),
		bytecode.NewInstruction(bytecode.OpLoadCell, bytecode.NewRegister(1), bytecode.NewRegister(5)),
		bytecode.NewInstruction(bytecode.OpReturn, bytecode.NewRegister(1)),
	}

	assertBytecodeEqual(t, program.Bytecode, expected)
}

func runCoalescing(t *testing.T, program *bytecode.Program) {
	t.Helper()

	p := NewPipeline()
	p.Add(NewLivenessAnalysisPass())
	p.Add(NewRegisterCoalescingPass())

	if _, err := p.Run(program); err != nil {
		t.Fatalf("coalescing pipeline failed: %v", err)
	}
}

func assertBytecodeEqual(t *testing.T, got, want []bytecode.Instruction) {
	t.Helper()

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("bytecode mismatch:\nexpected: %#v\ngot: %#v", want, got)
	}
}
