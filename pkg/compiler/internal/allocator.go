package internal

import (
	"github.com/MontFerret/ferret/pkg/vm"
)

type (
	// RegisterType represents the type of a register
	RegisterType int

	// RegisterStatus tracks register usage
	RegisterStatus struct {
		IsAllocated bool
		Type        RegisterType // Type of variable stored
	}

	RegisterLifetime struct {
		Start int // Instruction number where register becomes live
		End   int // Instruction number where register dies
	}

	RegisterSequence struct {
		Registers []vm.Operand
	}

	// RegisterAllocator manages register allocation
	RegisterAllocator struct {
		registers    map[vm.Operand]*RegisterStatus
		nextRegister vm.Operand
		currentInstr int
	}
)

const (
	Temp   RegisterType = iota // Short-lived intermediate results
	Var                        // Local variables
	State                      // FOR loop state
	Result                     // FOR loop result
)

func NewRegisterSequence(registers ...vm.Operand) *RegisterSequence {
	return &RegisterSequence{
		Registers: registers,
	}
}

func NewRegisterAllocator() *RegisterAllocator {
	return &RegisterAllocator{
		registers:    make(map[vm.Operand]*RegisterStatus),
		nextRegister: vm.NoopOperand + 1, // we start at 1 to avoid NoopOperand
	}
}

func (ra *RegisterAllocator) Size() int {
	return int(ra.nextRegister)
}

// Allocate assigns a register based on variable type
func (ra *RegisterAllocator) Allocate(regType RegisterType) vm.Operand {
	// Try to find a free register first
	reg, found := ra.findFreeRegister()

	if found {
		return reg
	}

	// If no free registers, create a new one
	newReg := ra.nextRegister
	ra.nextRegister++

	// Initialize register status
	ra.registers[newReg] = &RegisterStatus{
		IsAllocated: true,
		Type:        regType,
	}

	return newReg
}

// Free marks a register as available
func (ra *RegisterAllocator) Free(reg vm.Operand) {
	//if status, exists := ra.registers[reg]; exists {
	//status.IsAllocated = false
	//status.Lifetime.End = ra.currentInstr

	//// Clean up interference graph
	//delete(ra.usageGraph, reg)
	//
	//for _, edges := range ra.usageGraph {
	//	delete(edges, reg)
	//}
	//}
}

// findFreeRegister looks for an unused register
func (ra *RegisterAllocator) findFreeRegister() (vm.Operand, bool) {
	// First, try to find a completely free register
	for reg, status := range ra.registers {
		if !status.IsAllocated {
			return reg, true
		}
	}

	return 0, false
}

// AllocateSequence allocates a sequence of registers for function arguments or similar uses
func (ra *RegisterAllocator) AllocateSequence(count int) *RegisterSequence {
	sequence := &RegisterSequence{
		Registers: make([]vm.Operand, count),
	}

	// First pass: try to find contiguous free registers
	startReg, found := ra.findContiguousRegisters(count)

	if found {
		// Use contiguous block
		for i := 0; i < count; i++ {
			reg := startReg + vm.Operand(i)
			sequence.Registers[i] = reg

			// Initialize or update register status
			ra.registers[reg] = &RegisterStatus{
				IsAllocated: true,
				Type:        Temp,
			}
		}
	} else {
		// Allocate registers individually if contiguous block not available
		for i := 0; i < count; i++ {
			reg := ra.Allocate(Temp)
			sequence.Registers[i] = reg
		}
	}

	return sequence
}

// findContiguousRegisters attempts to find a block of consecutive free registers
func (ra *RegisterAllocator) findContiguousRegisters(count int) (vm.Operand, bool) {
	if count <= 0 {
		return 0, false
	}

	// First, try to find a contiguous block in existing registers
	maxReg := ra.nextRegister
	for start := vm.NoopOperand + 1; start < maxReg; start++ {
		if ra.isContiguousBlockFree(start, count) {
			return start, true
		}
	}

	// If no existing contiguous block found, allocate new block
	startReg := ra.nextRegister
	ra.nextRegister += vm.Operand(count)

	return startReg, true
}

// isContiguousBlockFree checks if a block of registers is available
func (ra *RegisterAllocator) isContiguousBlockFree(start vm.Operand, count int) bool {
	for i := 0; i < count; i++ {
		reg := start + vm.Operand(i)

		if status := ra.registers[reg]; status != nil && status.IsAllocated {
			return false
		}
	}
	return true
}

// FreeSequence frees all registers in a sequence
func (ra *RegisterAllocator) FreeSequence(seq *RegisterSequence) {
	for _, reg := range seq.Registers {
		ra.Free(reg)
	}
}
