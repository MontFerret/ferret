package optimization

import (
	"sort"

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

			if bytecode.IsTerminatorOpcode(lastInst.Opcode) {
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
		role := bytecode.OpcodeInfoOf(op).ControlFlow

		if role == bytecode.ControlFlowJumpUnconditional || role == bytecode.ControlFlowJumpConditional {
			target := b.jumpTarget(i, inst)

			if target >= 0 && target < len(instructions) {
				leaders[target] = true
			}

			// Instruction after branch is a leader (fall-through or unreachable block start).
			if i+1 < len(instructions) {
				leaders[i+1] = true
			}

			continue
		}

		if role == bytecode.ControlFlowTerminator {
			// Instruction after terminator is a leader (if it exists, it's unreachable but still a block).
			if i+1 < len(instructions) {
				leaders[i+1] = true
			}
		}
	}

	for _, entry := range b.program.CatchTable {
		start := entry[0]
		end := entry[1]
		handler := entry[2]

		if start >= 0 && start < len(instructions) {
			leaders[start] = true
		}

		if end >= 0 && end+1 < len(instructions) {
			leaders[end+1] = true
		}

		if handler >= 0 && handler < len(instructions) {
			leaders[handler] = true
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

	sort.Ints(leaderIndices)

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
		role := bytecode.OpcodeInfoOf(lastInst.Opcode).ControlFlow

		switch role {
		case bytecode.ControlFlowJumpUnconditional:
			target := b.jumpTarget(lastIdx, lastInst)
			if targetBlock, ok := indexToBlock[target]; ok {
				block.AddSuccessor(targetBlock)
			}
		case bytecode.ControlFlowJumpConditional:
			target := b.jumpTarget(lastIdx, lastInst)
			if targetBlock, ok := indexToBlock[target]; ok {
				block.AddSuccessor(targetBlock)
			}

			if lastIdx+1 < len(instructions) {
				if nextBlock, ok := indexToBlock[lastIdx+1]; ok {
					block.AddSuccessor(nextBlock)
				}
			}
		case bytecode.ControlFlowTerminator:
			// Return/tail-call doesn't add regular successors; handled by exit block.
		default:
			if lastIdx+1 < len(instructions) {
				if nextBlock, ok := indexToBlock[lastIdx+1]; ok {
					block.AddSuccessor(nextBlock)
				}
			}
		}

		b.addCatchSuccessors(block, indexToBlock)
	}
}

func (b *Builder) addCatchSuccessors(block *BasicBlock, indexToBlock map[int]*BasicBlock) {
	if block == nil || len(block.Instructions) == 0 {
		return
	}

	for _, entry := range b.program.CatchTable {
		if !blockOverlapsCatch(block, entry) {
			continue
		}

		handler := entry[2]
		if handlerBlock, ok := indexToBlock[handler]; ok {
			block.AddSuccessor(handlerBlock)
		}
	}
}

func blockOverlapsCatch(block *BasicBlock, entry bytecode.Catch) bool {
	start := entry[0]
	end := entry[1]

	return start >= 0 && end >= start && block.Start <= end && start <= block.End
}

func (b *Builder) jumpTarget(pc int, inst bytecode.Instruction) int {
	if inst.Opcode == bytecode.OpMatchLoadPropertyConst {
		targets := b.program.Metadata.MatchFailTargets
		if pc >= 0 && pc < len(targets) {
			return targets[pc]
		}

		return -1
	}

	targetIdx := bytecode.JumpTargetOperandIndex(inst.Opcode)
	if targetIdx < 0 {
		return -1
	}

	return int(inst.Operands[targetIdx])
}
