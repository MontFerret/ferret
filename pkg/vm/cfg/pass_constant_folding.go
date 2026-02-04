package cfg

import (
	"github.com/MontFerret/ferret/pkg/runtime"
	"github.com/MontFerret/ferret/pkg/vm"
)

// ConstantFoldingPass performs constant folding and propagation
type ConstantFoldingPass struct{}

// NewConstantFoldingPass creates a new constant folding pass
func NewConstantFoldingPass() *ConstantFoldingPass {
	return &ConstantFoldingPass{}
}

// Name returns the pass name
func (p *ConstantFoldingPass) Name() string {
	return "constant-folding"
}

// Run executes constant folding on the program
func (p *ConstantFoldingPass) Run(program *vm.Program, cfg *ControlFlowGraph) (*PassResult, error) {
	modified := false

	// Track known constant values for registers
	constants := make(map[int]runtime.Value)

	// Process each block
	for _, block := range cfg.Blocks {
		for i, inst := range block.Instructions {
			folded, newInst := tryFoldInstruction(inst, constants, program.Constants)
			if folded {
				block.Instructions[i] = newInst
				// Update the bytecode in the program
				bytecodeIdx := block.Start + i
				program.Bytecode[bytecodeIdx] = newInst
				modified = true
			}

			// Track constants defined by this instruction
			updateConstantTracking(inst, constants, program.Constants)
		}
	}

	return &PassResult{
		Modified: modified,
		Metadata: map[string]interface{}{
			"constants_folded": modified,
		},
	}, nil
}

// tryFoldInstruction attempts to fold an instruction with constant operands
func tryFoldInstruction(inst vm.Instruction, constants map[int]runtime.Value, programConstants []runtime.Value) (bool, vm.Instruction) {
	op := inst.Opcode
	dst, src1, src2 := inst.Operands[0], inst.Operands[1], inst.Operands[2]

	// Try to get constant values for source operands
	val1, ok1 := getConstantValue(src1, constants, programConstants)
	val2, ok2 := getConstantValue(src2, constants, programConstants)

	switch op {
	case vm.OpAdd:
		if ok1 && ok2 {
			if int1, ok := val1.(runtime.Int); ok {
				if int2, ok := val2.(runtime.Int); ok {
					result := int1 + int2
					// Store result in constants tracking
					constants[dst.Register()] = result
					// Can't modify instruction to use a new constant without modifying constant table
					// Just track it for now
					return false, inst
				}
			}
		}

	case vm.OpSub:
		if ok1 && ok2 {
			if int1, ok := val1.(runtime.Int); ok {
				if int2, ok := val2.(runtime.Int); ok {
					result := int1 - int2
					constants[dst.Register()] = result
					return false, inst
				}
			}
		}

	case vm.OpMulti:
		if ok1 && ok2 {
			if int1, ok := val1.(runtime.Int); ok {
				if int2, ok := val2.(runtime.Int); ok {
					result := int1 * int2
					constants[dst.Register()] = result
					return false, inst
				}
			}
		}

	case vm.OpMove:
		// Propagate constant through move
		if ok1 {
			constants[dst.Register()] = val1
		}
	}

	return false, inst
}

// getConstantValue retrieves the constant value for an operand if available
func getConstantValue(op vm.Operand, constants map[int]runtime.Value, programConstants []runtime.Value) (runtime.Value, bool) {
	if op.IsConstant() {
		idx := op.Constant()
		if idx >= 0 && idx < len(programConstants) {
			return programConstants[idx], true
		}
	} else if op.IsRegister() {
		reg := op.Register()
		if val, ok := constants[reg]; ok {
			return val, true
		}
	}
	return nil, false
}

// updateConstantTracking updates the constant tracking map based on instruction
func updateConstantTracking(inst vm.Instruction, constants map[int]runtime.Value, programConstants []runtime.Value) {
	op := inst.Opcode
	dst, src1 := inst.Operands[0], inst.Operands[1]

	switch op {
	case vm.OpLoadConst:
		// dst gets a constant value
		if src1.IsConstant() && dst.IsRegister() {
			idx := src1.Constant()
			if idx >= 0 && idx < len(programConstants) {
				constants[dst.Register()] = programConstants[idx]
			}
		}

	case vm.OpLoadNone:
		if dst.IsRegister() {
			constants[dst.Register()] = runtime.None
		}

	case vm.OpLoadBool:
		if dst.IsRegister() {
			constants[dst.Register()] = runtime.Boolean(src1 == 1)
		}

	case vm.OpLoadZero:
		if dst.IsRegister() {
			constants[dst.Register()] = runtime.ZeroInt
		}

	case vm.OpMove:
		// Propagate constant through move
		if src1.IsRegister() {
			if val, ok := constants[src1.Register()]; ok && dst.IsRegister() {
				constants[dst.Register()] = val
			}
		} else if src1.IsConstant() && dst.IsRegister() {
			idx := src1.Constant()
			if idx >= 0 && idx < len(programConstants) {
				constants[dst.Register()] = programConstants[idx]
			}
		}

	default:
		// Most operations invalidate constant tracking for dst
		if dst.IsRegister() {
			delete(constants, dst.Register())
		}
	}
}
