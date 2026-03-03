package vm

import (
	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

type callFrame struct {
	returnPC   int
	returnDest bytecode.Operand
	registers  []runtime.Value
	protected  bool
	fnID       int
}

func (vm *VM) pushFrame(returnPC int, returnDest bytecode.Operand, protected bool, fnID int) {
	vm.frames = append(vm.frames, callFrame{
		returnPC:   returnPC,
		returnDest: returnDest,
		registers:  vm.registers.Values,
		protected:  protected,
		fnID:       fnID,
	})
}

func (vm *VM) popFrame() (callFrame, bool) {
	if len(vm.frames) == 0 {
		return callFrame{}, false
	}

	frame := vm.frames[len(vm.frames)-1]
	vm.frames = vm.frames[:len(vm.frames)-1]
	return frame, true
}

func (vm *VM) unwindToProtected() bool {
	for i := len(vm.frames) - 1; i >= 0; i-- {
		if !vm.frames[i].protected {
			continue
		}

		frame := vm.frames[i]
		for j := i + 1; j < len(vm.frames); j++ {
			vm.regPool.put(vm.frames[j].registers)
		}

		vm.frames = vm.frames[:i]
		vm.regPool.put(vm.registers.Values)
		vm.registers.Values = frame.registers

		if frame.returnDest.IsRegister() {
			vm.registers.Values[frame.returnDest] = runtime.None
		}

		vm.pc = frame.returnPC
		return true
	}

	return false
}
