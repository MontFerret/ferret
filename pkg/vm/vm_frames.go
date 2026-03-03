package vm

import "github.com/MontFerret/ferret/v2/pkg/runtime"

func (vm *VM) unwindToProtected() bool {
	registers, pc, ok := vm.frames.UnwindToProtectedFrame(vm.registers.Values)
	if !ok {
		return false
	}

	vm.registers.Values = registers
	vm.pc = pc
	return true
}

func (vm *VM) returnToCaller(retVal runtime.Value) bool {
	registers, pc, ok := vm.frames.ReturnToCaller(vm.registers.Values, retVal)
	if !ok {
		return false
	}

	vm.registers.Values = registers
	vm.pc = pc
	return true
}
