package optimization

import (
	"github.com/MontFerret/ferret/v2/pkg/bytecode"
)

// Builder constructs a control flow graph from bytecode
type Builder struct {
	program *bytecode.Program
}

// NewBuilder creates a new CFG builder for the given program
func NewBuilder(program *bytecode.Program) *Builder {
	return &Builder{
		program: program,
	}
}

// Build constructs the control flow graph for the program
func (b *Builder) Build() (*ControlFlowGraph, error) {
	instructions := b.program.Bytecode

	if len(instructions) == 0 {
		return &ControlFlowGraph{
			Entry:  nil,
			Exit:   NewBasicBlock(0, -1), // Virtual exit block
			Blocks: []*BasicBlock{},
		}, nil
	}

	// Step 1: Identify leaders (start of basic blocks)
	leaders := b.identifyLeaders(instructions)

	// Step 2: Build basic blocks
	blocks := b.buildBasicBlocks(instructions, leaders)

	// Step 3: Create control flow edges
	b.createEdges(instructions, blocks)

	// Step 4: Create virtual exit block and connect return statements
	exit := NewBasicBlock(len(blocks), len(instructions))

	for _, block := range blocks {
		if block.IsTerminator() && len(block.Instructions) > 0 {
			lastInst := block.Instructions[len(block.Instructions)-1]

			if lastInst.Opcode == bytecode.OpReturn {
				block.AddSuccessor(exit)
			}
		}
	}

	blocks = append(blocks, exit)

	cfg := &ControlFlowGraph{
		Entry:  blocks[0],
		Exit:   exit,
		Blocks: blocks,
	}

	return cfg, nil
}

// identifyLeaders finds all instruction indices that start a new basic block
func (b *Builder) identifyLeaders(instructions []bytecode.Instruction) map[int]bool {
	leaders := make(map[int]bool)

	// The first instruction is always a leader
	leaders[0] = true

	// Scan through bytecode to find leaders
	for i, inst := range instructions {
		op := inst.Opcode

		switch op {
		case bytecode.OpJump:
			// Target of jump is a leader
			target := int(inst.Operands[0])

			if target >= 0 && target < len(instructions) {
				leaders[target] = true
			}

			// Instruction after jump is a leader (unreachable, but still a block start)
			if i+1 < len(instructions) {
				leaders[i+1] = true
			}
		case bytecode.OpJumpIfFalse, bytecode.OpJumpIfTrue, bytecode.OpJumpIfNone, bytecode.OpJumpIfNe, bytecode.OpJumpIfNeConst:
			// Target of conditional jump is a leader
			target := int(inst.Operands[0])

			if target >= 0 && target < len(instructions) {
				leaders[target] = true
			}

			// Instruction after conditional jump is a leader (fall-through path)
			if i+1 < len(instructions) {
				leaders[i+1] = true
			}
		case bytecode.OpIterNext:
			// OpIterNext is like a conditional jump: when iterator is done, jumps to dst
			target := int(inst.Operands[0])

			if target >= 0 && target < len(instructions) {
				leaders[target] = true
			}

			// Instruction after OpIterNext is a leader (fall-through when iterator has more)
			if i+1 < len(instructions) {
				leaders[i+1] = true
			}
		case bytecode.OpReturn:
			// Instruction after return is a leader (if it exists, it's unreachable but still a block)
			if i+1 < len(instructions) {
				leaders[i+1] = true
			}
		}
	}

	return leaders
}

// buildBasicBlocks creates basic blocks from the identified leaders
func (b *Builder) buildBasicBlocks(instructions []bytecode.Instruction, leaders map[int]bool) []*BasicBlock {
	blocks := make([]*BasicBlock, 0)
	blockMap := make(map[int]*BasicBlock) // Maps start index to block

	// Create blocks for each leader
	leaderIndices := make([]int, 0, len(leaders))
	for idx := range leaders {
		leaderIndices = append(leaderIndices, idx)
	}

	// Sort leader indices
	for i := 0; i < len(leaderIndices); i++ {
		for j := i + 1; j < len(leaderIndices); j++ {
			if leaderIndices[i] > leaderIndices[j] {
				leaderIndices[i], leaderIndices[j] = leaderIndices[j], leaderIndices[i]
			}
		}
	}

	// Create basic blocks
	for i, start := range leaderIndices {
		block := NewBasicBlock(i, start)
		blockMap[start] = block
		blocks = append(blocks, block)

		// Determine the end of this block
		end := len(instructions)
		if i+1 < len(leaderIndices) {
			end = leaderIndices[i+1]
		}

		// Add instructions to the block
		for j := start; j < end; j++ {
			block.AddInstruction(instructions[j])
		}
	}

	return blocks
}

// createEdges creates control flow edges between basic blocks
func (b *Builder) createEdges(instructions []bytecode.Instruction, blocks []*BasicBlock) {
	// Create a map from instruction index to block for quick lookup
	indexToBlock := make(map[int]*BasicBlock)
	for _, block := range blocks {
		indexToBlock[block.Start] = block
	}

	// For each block, determine its successors
	for _, block := range blocks {
		if len(block.Instructions) == 0 {
			continue
		}

		lastInst := block.Instructions[len(block.Instructions)-1]
		lastIdx := block.End

		switch lastInst.Opcode {
		case bytecode.OpJump:
			// Unconditional jump to target
			target := int(lastInst.Operands[0])
			if targetBlock, ok := indexToBlock[target]; ok {
				block.AddSuccessor(targetBlock)
			}

		case bytecode.OpJumpIfFalse, bytecode.OpJumpIfTrue, bytecode.OpJumpIfNone, bytecode.OpJumpIfNe, bytecode.OpJumpIfNeConst:
			// Conditional jump: has two successors
			// 1. Jump target
			target := int(lastInst.Operands[0])
			if targetBlock, ok := indexToBlock[target]; ok {
				block.AddSuccessor(targetBlock)
			}
			// 2. Fall-through to next instruction
			if lastIdx+1 < len(instructions) {
				if nextBlock, ok := indexToBlock[lastIdx+1]; ok {
					block.AddSuccessor(nextBlock)
				}
			}

		case bytecode.OpIterNext:
			// OpIterNext is like a conditional jump
			// 1. Jump target (when iterator is done)
			target := int(lastInst.Operands[0])
			if targetBlock, ok := indexToBlock[target]; ok {
				block.AddSuccessor(targetBlock)
			}

			// 2. Fall-through (when iterator has more elements)
			if lastIdx+1 < len(instructions) {
				if nextBlock, ok := indexToBlock[lastIdx+1]; ok {
					block.AddSuccessor(nextBlock)
				}
			}

		case bytecode.OpReturn:
			// Return doesn't add regular successors; handled by exit block

		default:
			// Regular instruction: fall through to next block
			if lastIdx+1 < len(instructions) {
				if nextBlock, ok := indexToBlock[lastIdx+1]; ok {
					block.AddSuccessor(nextBlock)
				}
			}
		}
	}
}
