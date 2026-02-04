package cfg

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

	// Initialize dominators
	dominators := make(map[int]map[int]bool)
	for _, block := range a.cfg.Blocks {
		dominators[block.ID] = make(map[int]bool)
		// Initially, every block is dominated by all blocks
		for _, b := range a.cfg.Blocks {
			dominators[block.ID][b.ID] = true
		}
	}

	// Entry is only dominated by itself
	dominators[a.cfg.Entry.ID] = map[int]bool{a.cfg.Entry.ID: true}

	// Iteratively compute dominators
	changed := true
	for changed {
		changed = false
		for _, block := range a.cfg.Blocks {
			if block == a.cfg.Entry {
				continue
			}

			// Compute intersection of dominators of all predecessors
			newDom := make(map[int]bool)
			first := true
			for _, pred := range block.Predecessors {
				if first {
					for id := range dominators[pred.ID] {
						newDom[id] = true
					}
					first = false
				} else {
					// Intersection
					for id := range newDom {
						if !dominators[pred.ID][id] {
							delete(newDom, id)
						}
					}
				}
			}

			// Add self
			newDom[block.ID] = true

			// Check if changed
			if len(newDom) != len(dominators[block.ID]) {
				changed = true
				dominators[block.ID] = newDom
			} else {
				for id := range newDom {
					if !dominators[block.ID][id] {
						changed = true
						break
					}
				}
				if changed {
					dominators[block.ID] = newDom
				}
			}
		}
	}

	// Find immediate dominator (closest dominator that is not the node itself)
	immediateDominators := make(map[int]*BasicBlock)
	for _, block := range a.cfg.Blocks {
		if block == a.cfg.Entry {
			continue
		}

		// Find the dominator with the largest ID (closest to this block in dominator tree)
		maxID := -1
		for id := range dominators[block.ID] {
			if id != block.ID && id > maxID {
				// Additional check: this must dominate all other dominators
				isDominatedByOthers := false
				for otherID := range dominators[block.ID] {
					if otherID != block.ID && otherID != id {
						if !dominators[id][otherID] {
							// This dominator is not dominated by another dominator of block
							isDominatedByOthers = true
							break
						}
					}
				}
				if !isDominatedByOthers {
					maxID = id
				}
			}
		}

		if maxID >= 0 {
			for _, b := range a.cfg.Blocks {
				if b.ID == maxID {
					immediateDominators[block.ID] = b
					break
				}
			}
		}
	}

	return immediateDominators
}

// ToDOT converts the CFG to Graphviz DOT format for visualization
func (cfg *ControlFlowGraph) ToDOT() string {
	var sb strings.Builder

	sb.WriteString("digraph CFG {\n")
	sb.WriteString("  node [shape=box];\n")

	// Write nodes
	for _, block := range cfg.Blocks {
		if block == cfg.Exit {
			sb.WriteString(fmt.Sprintf("  block%d [label=\"Exit\", shape=ellipse];\n", block.ID))
		} else {
			label := fmt.Sprintf("Block %d\\n[%d:%d]", block.ID, block.Start, block.End)
			sb.WriteString(fmt.Sprintf("  block%d [label=\"%s\"];\n", block.ID, label))
		}
	}

	// Write edges
	for _, block := range cfg.Blocks {
		for _, succ := range block.Successors {
			sb.WriteString(fmt.Sprintf("  block%d -> block%d;\n", block.ID, succ.ID))
		}
	}

	sb.WriteString("}\n")

	return sb.String()
}

// String returns a human-readable representation of the CFG
func (cfg *ControlFlowGraph) String() string {
	var sb strings.Builder

	sb.WriteString("Control Flow Graph:\n")
	sb.WriteString(fmt.Sprintf("  Entry: Block %d\n", cfg.Entry.ID))
	sb.WriteString(fmt.Sprintf("  Exit: Block %d\n", cfg.Exit.ID))
	sb.WriteString(fmt.Sprintf("  Blocks: %d\n\n", len(cfg.Blocks)))

	for _, block := range cfg.Blocks {
		if block != cfg.Exit {
			sb.WriteString(block.String())
			sb.WriteString("\n")
		}
	}

	return sb.String()
}
