package compiler

import (
	"errors"
)

// RegisterAllocator manages register allocation
type RegisterAllocator struct {
	usedRegisters map[Register]bool
	variables     map[string]*Variable
	currentInstr  int

	// Register pools for different purposes
	pools map[VarType][]Register
}

func NewRegisterAllocator(initialSize int) *RegisterAllocator {
	ra := &RegisterAllocator{
		usedRegisters: make(map[Register]bool),
		variables:     make(map[string]*Variable),
	}

	// Initialize register pools
	ra.initializePools(initialSize)
	return ra
}

func (ra *RegisterAllocator) initializePools(initialSize int) {
	// Allocate registers to different pools
	// Example distribution (adjustable based on query patterns):
	// - 20% for iterators
	// - 40% for temporaries
	// - 30% for collect operations
	// - 10% for results

	iteratorSize := initialSize * 2 / 10
	tempSize := initialSize * 4 / 10
	collectSize := initialSize * 3 / 10
	resultSize := initialSize - iteratorSize - tempSize - collectSize

	current := Register(0)

	// Iterator pool (R0-R19 in a 100-register machine)
	for i := 0; i < iteratorSize; i++ {
		ra.iteratorPool = append(ra.iteratorPool, current)
		current++
	}

	// Temporary pool (R20-R59)
	for i := 0; i < tempSize; i++ {
		ra.tempPool = append(ra.tempPool, current)
		current++
	}

	//// Collect pool (R60-R89)
	//for i := 0; i < collectSize; i++ {
	//	ra.collectPool = append(ra.collectPool, current)
	//	current++
	//}

	// Result pool (R90-R99)
	for i := 0; i < resultSize; i++ {
		ra.resultPool = append(ra.resultPool, current)
		current++
	}
}

// AllocateRegister assigns a register based on variable type
func (ra *RegisterAllocator) AllocateRegister(varType VarType) Register {
	var pool []Register

	switch varType {
	case VarIterator:
		pool = ra.iteratorPool
	case VarTemporary:
		pool = ra.tempPool
	case VarResult:
		pool = ra.resultPool
	}

	// Find first free register in appropriate pool
	for _, reg := range pool {
		if !ra.usedRegisters[reg] {
			ra.usedRegisters[reg] = true
			return reg
		}
	}

	// If no registers available in the pool, expand the pool
	nextPool := make([]Register, len(pool)*2)
	copy(nextPool, pool)

	// TODO: Create a custom error type
	return -1, errors.New("no registers available")
}

// LivenessAnalysis performs liveness analysis on variables
func (ra *RegisterAllocator) LivenessAnalysis(instructions []Instruction) {
	// Forward pass to record first use
	for i, instr := range instructions {
		for _, varName := range instr.VarRefs {
			if v, exists := ra.variables[varName]; exists {
				if v.FirstUse == -1 {
					v.FirstUse = i
				}
				v.LastUse = i
			}
		}
	}

	// Backward pass to optimize register allocation
	for i := len(instructions) - 1; i >= 0; i-- {
		for _, v := range ra.variables {
			if v.FirstUse <= i && i <= v.LastUse {
				v.IsLive = true
			} else if i < v.FirstUse || i > v.LastUse {
				v.IsLive = false
				// Free register if variable is dead
				if v.Register >= 0 {
					ra.usedRegisters[v.Register] = false
				}
			}
		}
	}
}
