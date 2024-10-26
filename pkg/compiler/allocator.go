package compiler

import "github.com/MontFerret/ferret/pkg/runtime"

type (
	// RegisterStatus tracks register usage
	RegisterStatus struct {
		IsAllocated bool
		LastUse     int               // Instruction number of last use
		NextUse     int               // Instruction number of next use
		VarName     string            // Associated variable name, if any
		VarType     VarType           // Type of variable stored
		Lifetime    *RegisterLifetime // Lifetime information
	}

	RegisterLifetime struct {
		Start int // Instruction number where register becomes live
		End   int // Instruction number where register dies
	}

	// RegisterAllocator manages register allocation
	RegisterAllocator struct {
		registers    map[runtime.Operand]*RegisterStatus
		nextRegister runtime.Operand
		currentInstr int
		lifetimes    map[string]*RegisterLifetime
		usageGraph   map[runtime.Operand]map[runtime.Operand]bool
	}
)

func NewRegisterAllocator() *RegisterAllocator {
	return &RegisterAllocator{
		registers:    make(map[runtime.Operand]*RegisterStatus),
		nextRegister: runtime.ResultOperand + 1, // we start at 1 to avoid ResultOperand
		lifetimes:    make(map[string]*RegisterLifetime),
		usageGraph:   make(map[runtime.Operand]map[runtime.Operand]bool),
	}
}

func (ra *RegisterAllocator) AllocateLocalVar(name string) runtime.Operand {
	// Allocate register
	reg := ra.Allocate(VarLocal)

	// Update register status
	ra.registers[reg].VarName = name

	return reg
}

func (ra *RegisterAllocator) AllocateTempVar() runtime.Operand {
	return ra.Allocate(VarTemporary)
}

// Allocate assigns a register based on variable type
func (ra *RegisterAllocator) Allocate(varType VarType) runtime.Operand {
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
		LastUse:     ra.currentInstr,
		NextUse:     -1,
		VarType:     varType,
		Lifetime:    &RegisterLifetime{Start: ra.currentInstr},
	}

	return newReg
}

// Free marks a register as available
func (ra *RegisterAllocator) Free(reg runtime.Operand) {
	if status, exists := ra.registers[reg]; exists {
		status.IsAllocated = false
		status.Lifetime.End = ra.currentInstr

		// Clean up interference graph
		delete(ra.usageGraph, reg)

		for _, edges := range ra.usageGraph {
			delete(edges, reg)
		}
	}
}

// findFreeRegister looks for an unused register
func (ra *RegisterAllocator) findFreeRegister() (runtime.Operand, bool) {
	// First, try to find a completely free register
	for reg, status := range ra.registers {
		if !status.IsAllocated {
			return reg, true
		}
	}

	// If no free registers, try to find one that's no longer needed
	var candidate runtime.Operand
	var found bool
	maxLastUse := -1

	for reg, status := range ra.registers {
		if status.NextUse == -1 && status.LastUse > maxLastUse {
			maxLastUse = status.LastUse
			candidate = reg
			found = true
		}
	}

	if found {
		// Free the candidate register
		ra.Free(candidate)

		return candidate, true
	}

	return 0, false
}

// UpdateUse updates the usage information for a register
func (ra *RegisterAllocator) UpdateUse(reg runtime.Operand) {
	status := ra.registers[reg]

	if status == nil {
		return
	}

	status.LastUse = ra.currentInstr

	// Update interference graph for simultaneously live registers
	for otherReg, otherStatus := range ra.registers {
		if otherReg != reg && otherStatus.IsAllocated &&
			ra.registersInterfere(reg, otherReg) {
			ra.addInterference(reg, otherReg)
		}
	}

	ra.currentInstr++
}

// registersInterfere checks if two registers have overlapping lifetimes
func (ra *RegisterAllocator) registersInterfere(reg1, reg2 runtime.Operand) bool {
	status1 := ra.registers[reg1]
	status2 := ra.registers[reg2]

	if status1 == nil || status2 == nil {
		return false
	}

	// Registers interfere if their lifetimes overlap
	return status1.Lifetime.Start <= status2.Lifetime.End &&
		status2.Lifetime.Start <= status1.Lifetime.End
}

// addInterference records that two registers interfere
func (ra *RegisterAllocator) addInterference(reg1, reg2 runtime.Operand) {
	if ra.usageGraph[reg1] == nil {
		ra.usageGraph[reg1] = make(map[runtime.Operand]bool)
	}
	if ra.usageGraph[reg2] == nil {
		ra.usageGraph[reg2] = make(map[runtime.Operand]bool)
	}

	ra.usageGraph[reg1][reg2] = true
	ra.usageGraph[reg2][reg1] = true
}
