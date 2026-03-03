package vm

func (vm *VM) unwindToProtected() bool {
	registers, pc, ok := vm.frames.UnwindToProtected(vm.registers.Values)
	if !ok {
		return false
	}

	vm.registers.Values = registers
	vm.pc = pc
	return true
}
