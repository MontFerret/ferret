package optimization

import (
	"fmt"
	"strings"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
)

// BasicBlock represents a sequence of instructions with a single entry and exit point
type BasicBlock struct {
	Instructions []bytecode.Instruction
	Successors   []*BasicBlock
	Predecessors []*BasicBlock
	ID           int
	Start        int
	End          int
}

// NewBasicBlock creates a new basic block with the given ID and start position
func NewBasicBlock(id, start int) *BasicBlock {
	return &BasicBlock{
		ID:           id,
		Start:        start,
		End:          start,
		Instructions: make([]bytecode.Instruction, 0),
		Successors:   make([]*BasicBlock, 0),
		Predecessors: make([]*BasicBlock, 0),
	}
}

// AddInstruction adds an instruction to the basic block
func (bb *BasicBlock) AddInstruction(inst bytecode.Instruction) {
	bb.Instructions = append(bb.Instructions, inst)
	bb.End = bb.Start + len(bb.Instructions) - 1
}

// AddSuccessor adds a successor block to this block
func (bb *BasicBlock) AddSuccessor(succ *BasicBlock) {
	if !bb.hasSuccessor(succ) {
		bb.Successors = append(bb.Successors, succ)
	}

	if !succ.hasPredecessor(bb) {
		succ.Predecessors = append(succ.Predecessors, bb)
	}
}

// hasSuccessor checks if the given block is already a successor
func (bb *BasicBlock) hasSuccessor(block *BasicBlock) bool {
	for _, succ := range bb.Successors {
		if succ.ID == block.ID {
			return true
		}
	}

	return false
}

// hasPredecessor checks if the given block is already a predecessor
func (bb *BasicBlock) hasPredecessor(block *BasicBlock) bool {
	for _, pred := range bb.Predecessors {
		if pred.ID == block.ID {
			return true
		}
	}

	return false
}

// IsTerminator returns true if the block ends with a terminator instruction
func (bb *BasicBlock) IsTerminator() bool {
	if len(bb.Instructions) == 0 {
		return false
	}

	lastInst := bb.Instructions[len(bb.Instructions)-1]
	role := bytecode.OpcodeInfoOf(lastInst.Opcode).ControlFlow

	return role == bytecode.ControlFlowTerminator ||
		role == bytecode.ControlFlowJumpUnconditional ||
		role == bytecode.ControlFlowJumpConditional
}

// String returns a string representation of the basic block
func (bb *BasicBlock) String() string {
	var sb strings.Builder

	fmt.Fprintf(&sb, "Block %d [%d:%d]:\n", bb.ID, bb.Start, bb.End)

	for i, inst := range bb.Instructions {
		fmt.Fprintf(&sb, "  %d: %s", bb.Start+i, inst.Opcode.String())

		if inst.Operands[0] != 0 || inst.Operands[1] != 0 || inst.Operands[2] != 0 {
			fmt.Fprintf(&sb, " %d", inst.Operands[0])

			if inst.Operands[1] != 0 || inst.Operands[2] != 0 {
				fmt.Fprintf(&sb, " %d", inst.Operands[1])

				if inst.Operands[2] != 0 {
					fmt.Fprintf(&sb, " %d", inst.Operands[2])
				}
			}
		}

		sb.WriteString("\n")
	}

	sb.WriteString("  Successors: ")

	for i, succ := range bb.Successors {
		if i > 0 {
			sb.WriteString(", ")
		}

		fmt.Fprintf(&sb, "%d", succ.ID)
	}

	sb.WriteString("\n")

	return sb.String()
}
