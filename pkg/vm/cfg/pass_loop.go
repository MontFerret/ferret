package cfg

import (
	"github.com/MontFerret/ferret/pkg/vm"
)

// LoopDetectionPass identifies natural loops in the CFG
type LoopDetectionPass struct{}

// NewLoopDetectionPass creates a new loop detection pass
func NewLoopDetectionPass() *LoopDetectionPass {
	return &LoopDetectionPass{}
}

// Name returns the pass name
func (p *LoopDetectionPass) Name() string {
	return "loop-detection"
}

// Run executes loop detection on the program
func (p *LoopDetectionPass) Run(program *vm.Program, cfg *ControlFlowGraph) (*PassResult, error) {
	loops := detectLoops(cfg)

	return &PassResult{
		Modified: false,
		Metadata: map[string]interface{}{
			"loops": loops,
		},
	}, nil
}

// Loop represents a natural loop in the CFG
type Loop struct {
	Header *BasicBlock   // Loop header (target of back edge)
	Blocks []*BasicBlock // All blocks in the loop
	Exits  []*BasicBlock // Exit blocks (have successors outside loop)
}

// detectLoops identifies all natural loops using back edges
func detectLoops(cfg *ControlFlowGraph) []*Loop {
	analyzer := NewAnalyzer(cfg)
	backEdges := analyzer.FindBackEdges()

	loops := make([]*Loop, 0)

	// For each back edge, find the natural loop
	for _, edge := range backEdges {
		tail, header := edge[0], edge[1]
		loop := findNaturalLoop(cfg, header, tail)
		loops = append(loops, loop)
	}

	return loops
}

// findNaturalLoop finds all blocks in the natural loop for a given back edge
func findNaturalLoop(cfg *ControlFlowGraph, header, tail *BasicBlock) *Loop {
	loopBlocks := make(map[int]*BasicBlock)
	loopBlocks[header.ID] = header

	// Use worklist algorithm to find all blocks in the loop
	worklist := []*BasicBlock{tail}
	visited := make(map[int]bool)
	visited[header.ID] = true

	for len(worklist) > 0 {
		block := worklist[0]
		worklist = worklist[1:]

		if visited[block.ID] {
			continue
		}

		visited[block.ID] = true
		loopBlocks[block.ID] = block

		// Add all predecessors to worklist
		for _, pred := range block.Predecessors {
			if !visited[pred.ID] {
				worklist = append(worklist, pred)
			}
		}
	}

	// Convert map to slice
	blocks := make([]*BasicBlock, 0, len(loopBlocks))
	for _, block := range loopBlocks {
		blocks = append(blocks, block)
	}

	// Find exit blocks (blocks in loop with successors outside loop)
	exits := make([]*BasicBlock, 0)
	for _, block := range blocks {
		for _, succ := range block.Successors {
			if loopBlocks[succ.ID] == nil {
				exits = append(exits, block)
				break
			}
		}
	}

	return &Loop{
		Header: header,
		Blocks: blocks,
		Exits:  exits,
	}
}
