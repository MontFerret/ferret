package optimization

import (
	"fmt"
	"strings"
)

// Analyzer provides analysis capabilities for control flow graphs
type Analyzer struct {
	cfg *ControlFlowGraph
}

// NewAnalyzer creates a new CFG analyzer
func NewAnalyzer(cfg *ControlFlowGraph) *Analyzer {
	return &Analyzer{cfg: cfg}
}

// FindReachableBlocks returns all blocks reachable from the entry block
func (a *Analyzer) FindReachableBlocks() []*BasicBlock {
	if a.cfg.Entry == nil {
		return []*BasicBlock{}
	}

	reachable := make(map[int]bool)
	visited := make(map[int]bool)
	queue := []*BasicBlock{a.cfg.Entry}

	for len(queue) > 0 {
		block := queue[0]
		queue = queue[1:]

		if visited[block.ID] {
			continue
		}

		visited[block.ID] = true
		reachable[block.ID] = true

		for _, succ := range block.Successors {
			if !visited[succ.ID] {
				queue = append(queue, succ)
			}
		}
	}

	result := make([]*BasicBlock, 0, len(reachable))
	for _, block := range a.cfg.Blocks {
		if reachable[block.ID] {
			result = append(result, block)
		}
	}

	return result
}

// FindUnreachableBlocks returns all blocks not reachable from the entry block
func (a *Analyzer) FindUnreachableBlocks() []*BasicBlock {
	reachable := a.FindReachableBlocks()
	reachableMap := make(map[int]bool)

	for _, block := range reachable {
		reachableMap[block.ID] = true
	}

	unreachable := make([]*BasicBlock, 0)

	for _, block := range a.cfg.Blocks {
		if !reachableMap[block.ID] && block != a.cfg.Exit {
			unreachable = append(unreachable, block)
		}
	}

	return unreachable
}

// FindBackEdges identifies back edges in the CFG (edges that create loops)
// A back edge is an edge from a node to one of its ancestors in a DFS traversal
func (a *Analyzer) FindBackEdges() [][2]*BasicBlock {
	if a.cfg.Entry == nil {
		return [][2]*BasicBlock{}
	}

	backEdges := make([][2]*BasicBlock, 0)
	visited := make(map[int]bool)
	inStack := make(map[int]bool)

	var dfs func(*BasicBlock)
	dfs = func(block *BasicBlock) {
		visited[block.ID] = true
		inStack[block.ID] = true

		for _, succ := range block.Successors {
			if !visited[succ.ID] {
				dfs(succ)
			} else if inStack[succ.ID] {
				// Found a back edge
				backEdges = append(backEdges, [2]*BasicBlock{block, succ})
			}
		}

		inStack[block.ID] = false
	}

	dfs(a.cfg.Entry)

	return backEdges
}

// CalculateDominators computes the dominator tree for the CFG
// A block A dominates block B if all paths from entry to B go through A
func (a *Analyzer) CalculateDominators() map[int]*BasicBlock {
	if a.cfg.Entry == nil {
		return map[int]*BasicBlock{}
	}

	dominators := initializeDominators(a.cfg)
	computeDominatorsFixedPoint(a.cfg, dominators)

	return buildImmediateDominators(a.cfg, dominators)
}

func initializeDominators(cfg *ControlFlowGraph) map[int]map[int]bool {
	dominators := make(map[int]map[int]bool, len(cfg.Blocks))

	for _, block := range cfg.Blocks {
		dominators[block.ID] = make(map[int]bool, len(cfg.Blocks))

		for _, other := range cfg.Blocks {
			dominators[block.ID][other.ID] = true
		}
	}

	dominators[cfg.Entry.ID] = map[int]bool{cfg.Entry.ID: true}
	return dominators
}

func computeDominatorsFixedPoint(cfg *ControlFlowGraph, dominators map[int]map[int]bool) {
	changed := true

	for changed {
		changed = false

		for _, block := range cfg.Blocks {
			if block == cfg.Entry {
				continue
			}

			newSet := intersectPredecessorDominators(block, dominators)
			newSet[block.ID] = true

			if updateDominatorSet(dominators, block.ID, newSet) {
				changed = true
			}
		}
	}
}

func intersectPredecessorDominators(block *BasicBlock, dominators map[int]map[int]bool) map[int]bool {
	intersection := make(map[int]bool)
	firstPred := true

	for _, pred := range block.Predecessors {
		if firstPred {
			for id := range dominators[pred.ID] {
				intersection[id] = true
			}
			firstPred = false
			continue
		}

		for id := range intersection {
			if !dominators[pred.ID][id] {
				delete(intersection, id)
			}
		}
	}

	return intersection
}

func updateDominatorSet(dominators map[int]map[int]bool, blockID int, newSet map[int]bool) bool {
	current := dominators[blockID]
	if mapsEqual(current, newSet) {
		return false
	}

	dominators[blockID] = newSet
	return true
}

func buildImmediateDominators(cfg *ControlFlowGraph, dominators map[int]map[int]bool) map[int]*BasicBlock {
	immediate := make(map[int]*BasicBlock)
	blockMap := make(map[int]*BasicBlock, len(cfg.Blocks))

	for _, block := range cfg.Blocks {
		blockMap[block.ID] = block
	}

	for _, block := range cfg.Blocks {
		if block == cfg.Entry {
			continue
		}

		idomID, ok := findImmediateDominator(block.ID, dominators)
		if !ok {
			continue
		}

		if idom := blockMap[idomID]; idom != nil {
			immediate[block.ID] = idom
		}
	}

	return immediate
}

func findImmediateDominator(blockID int, dominators map[int]map[int]bool) (int, bool) {
	blockDominators := dominators[blockID]

	for domID := range blockDominators {
		if domID == blockID {
			continue
		}

		if !isDominatedByOtherCandidate(blockID, domID, blockDominators, dominators) {
			return domID, true
		}
	}

	return 0, false
}

func isDominatedByOtherCandidate(blockID, candidateID int, blockDominators map[int]bool, dominators map[int]map[int]bool) bool {
	for otherID := range blockDominators {
		if otherID == blockID || otherID == candidateID {
			continue
		}

		if dominators[candidateID][otherID] {
			return true
		}
	}

	return false
}

// ToDOT converts the CFG to Graphviz DOT format for visualization
func (cfg *ControlFlowGraph) ToDOT() string {
	var sb strings.Builder

	sb.WriteString("digraph CFG {\n")
	sb.WriteString("  node [shape=box];\n")

	// Write nodes
	for _, block := range cfg.Blocks {
		if block == cfg.Exit {
			fmt.Fprintf(&sb, "  block%d [label=\"Exit\", shape=ellipse];\n", block.ID)
		} else {
			label := fmt.Sprintf("Block %d\\n[%d:%d]", block.ID, block.Start, block.End)
			fmt.Fprintf(&sb, "  block%d [label=\"%s\"];\n", block.ID, label)
		}
	}

	// Write edges
	for _, block := range cfg.Blocks {
		for _, succ := range block.Successors {
			fmt.Fprintf(&sb, "  block%d -> block%d;\n", block.ID, succ.ID)
		}
	}

	sb.WriteString("}\n")

	return sb.String()
}

// String returns a human-readable representation of the CFG
func (cfg *ControlFlowGraph) String() string {
	var sb strings.Builder

	sb.WriteString("Control Flow Graph:\n")
	fmt.Fprintf(&sb, "  Entry: Block %d\n", cfg.Entry.ID)
	fmt.Fprintf(&sb, "  Exit: Block %d\n", cfg.Exit.ID)
	fmt.Fprintf(&sb, "  Blocks: %d\n\n", len(cfg.Blocks))

	for _, block := range cfg.Blocks {
		if block != cfg.Exit {
			sb.WriteString(block.String())
			sb.WriteString("\n")
		}
	}

	return sb.String()
}
