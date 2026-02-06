package optimization

import (
	"reflect"
	"testing"

	"github.com/MontFerret/ferret/pkg/runtime"
	"github.com/MontFerret/ferret/pkg/vm"
)

func TestRegisterCoalescing_CoalescesWhenSrcDies(t *testing.T) {
	program := &vm.Program{
		Constants: []runtime.Value{
			runtime.Int(10),
		},
		Registers: 4,
		Bytecode: []vm.Instruction{
			vm.NewInstruction(vm.OpLoadConst, vm.NewRegister(1), vm.NewConstant(0)),
			vm.NewInstruction(vm.OpMove, vm.NewRegister(2), vm.NewRegister(1)),
			vm.NewInstruction(vm.OpAdd, vm.NewRegister(3), vm.NewRegister(2), vm.NewRegister(2)),
			vm.NewInstruction(vm.OpReturn, vm.NewRegister(3)),
		},
	}

	runCoalescing(t, program)

	expected := []vm.Instruction{
		vm.NewInstruction(vm.OpLoadConst, vm.NewRegister(1), vm.NewConstant(0)),
		vm.NewInstruction(vm.OpMove, vm.NewRegister(1), vm.NewRegister(1)),
		vm.NewInstruction(vm.OpAdd, vm.NewRegister(1), vm.NewRegister(1), vm.NewRegister(1)),
		vm.NewInstruction(vm.OpReturn, vm.NewRegister(1)),
	}

	assertBytecodeEqual(t, program.Bytecode, expected)
}

func TestRegisterCoalescing_NoCoalesceWhenInterfering(t *testing.T) {
	program := &vm.Program{
		Constants: []runtime.Value{
			runtime.Int(1),
		},
		Registers: 4,
		Bytecode: []vm.Instruction{
			vm.NewInstruction(vm.OpLoadConst, vm.NewRegister(1), vm.NewConstant(0)),
			vm.NewInstruction(vm.OpMove, vm.NewRegister(2), vm.NewRegister(1)),
			vm.NewInstruction(vm.OpIncr, vm.NewRegister(2)),
			vm.NewInstruction(vm.OpAdd, vm.NewRegister(3), vm.NewRegister(2), vm.NewRegister(1)),
			vm.NewInstruction(vm.OpReturn, vm.NewRegister(3)),
		},
	}

	runCoalescing(t, program)

	expected := []vm.Instruction{
		vm.NewInstruction(vm.OpLoadConst, vm.NewRegister(1), vm.NewConstant(0)),
		vm.NewInstruction(vm.OpMove, vm.NewRegister(2), vm.NewRegister(1)),
		vm.NewInstruction(vm.OpIncr, vm.NewRegister(2)),
		vm.NewInstruction(vm.OpAdd, vm.NewRegister(1), vm.NewRegister(2), vm.NewRegister(1)),
		vm.NewInstruction(vm.OpReturn, vm.NewRegister(1)),
	}

	assertBytecodeEqual(t, program.Bytecode, expected)
}

func TestRegisterCoalescing_NoCoalesceForRangeSensitiveRegs(t *testing.T) {
	program := &vm.Program{
		Constants: []runtime.Value{
			runtime.Int(1),
			runtime.Int(2),
		},
		Registers: 6,
		Bytecode: []vm.Instruction{
			vm.NewInstruction(vm.OpLoadConst, vm.NewRegister(1), vm.NewConstant(0)),
			vm.NewInstruction(vm.OpLoadConst, vm.NewRegister(2), vm.NewConstant(1)),
			vm.NewInstruction(vm.OpMove, vm.NewRegister(3), vm.NewRegister(1)),
			vm.NewInstruction(vm.OpLoadObject, vm.NewRegister(4), vm.NewRegister(1), vm.NewRegister(2)),
			vm.NewInstruction(vm.OpAdd, vm.NewRegister(5), vm.NewRegister(3), vm.NewRegister(1)),
			vm.NewInstruction(vm.OpReturn, vm.NewRegister(5)),
		},
	}

	runCoalescing(t, program)

	expected := []vm.Instruction{
		vm.NewInstruction(vm.OpLoadConst, vm.NewRegister(1), vm.NewConstant(0)),
		vm.NewInstruction(vm.OpLoadConst, vm.NewRegister(2), vm.NewConstant(1)),
		vm.NewInstruction(vm.OpMove, vm.NewRegister(3), vm.NewRegister(1)),
		vm.NewInstruction(vm.OpLoadObject, vm.NewRegister(4), vm.NewRegister(1), vm.NewRegister(2)),
		vm.NewInstruction(vm.OpAdd, vm.NewRegister(3), vm.NewRegister(3), vm.NewRegister(1)),
		vm.NewInstruction(vm.OpReturn, vm.NewRegister(3)),
	}

	assertBytecodeEqual(t, program.Bytecode, expected)
}

func runCoalescing(t *testing.T, program *vm.Program) {
	t.Helper()

	p := NewPipeline()
	p.Add(NewLivenessAnalysisPass())
	p.Add(NewRegisterCoalescingPass())

	if _, err := p.Run(program); err != nil {
		t.Fatalf("coalescing pipeline failed: %v", err)
	}
}

func assertBytecodeEqual(t *testing.T, got, want []vm.Instruction) {
	t.Helper()

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("bytecode mismatch:\nexpected: %#v\ngot: %#v", want, got)
	}
}
