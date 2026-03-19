package optimization

import "fmt"

const RegisterCoalescingPassName = "register-coalescing"

type (
	// RegisterCoalescingPass performs register coalescing to reduce register usage
	// by renumbering registers based on liveness intervals
	RegisterCoalescingPass struct{}

	liveInterval struct {
		reg   int
		start int
		end   int
	}
)

// NewRegisterCoalescingPass creates a new register coalescing pass
func NewRegisterCoalescingPass() Pass {
	return &RegisterCoalescingPass{}
}

// Name returns the pass name
func (p *RegisterCoalescingPass) Name() string {
	return RegisterCoalescingPassName
}

func (p *RegisterCoalescingPass) Requires() []string {
	return []string{LivenessAnalysisPassName}
}

// Run executes register coalescing on the program
func (p *RegisterCoalescingPass) Run(ctx *PassContext) (*PassResult, error) {
	meta, ok := ctx.Metadata[LivenessAnalysisPassName].(map[string]any)

	if !ok {
		return nil, fmt.Errorf("%w: pass %q requires %q metadata", ErrMissingDependency, RegisterCoalescingPassName, LivenessAnalysisPassName)
	}
	raw, ok := meta[LivenessAnalysis]

	if !ok {
		return nil, fmt.Errorf("%w: pass %q requires %q key in %q metadata", ErrMissingDependency, RegisterCoalescingPassName, LivenessAnalysis, LivenessAnalysisPassName)
	}

	liveness, ok := raw.(map[int]*LivenessInfo)

	if !ok {
		return nil, fmt.Errorf("%w: pass %q received unexpected liveness data type from %q", ErrMissingDependency, RegisterCoalescingPassName, LivenessAnalysisPassName)
	}

	unsafeRegs := collectPinnedRegs(ctx.Program)
	cfg := ctx.CFG

	// Fold trivial move chains into their defining instruction to shorten live ranges.
	folded := foldMovesIntoDefs(ctx.Program, cfg, liveness, unsafeRegs)
	if folded {
		// Recompute liveness after rewriting registers to keep analyses consistent.
		liveness = computeLiveness(cfg)
	}

	// Build interference graph
	interferenceGraph := buildInterferenceGraph(cfg, liveness, ctx.Program.Registers)

	// Count register definitions to keep move coalescing safe.
	defCounts := countRegisterDefs(ctx.Program)

	// First, try move-based coalescing
	coalesceMap := findCoalesceCandidates(ctx.Program, cfg, interferenceGraph, unsafeRegs, defCounts)

	// Apply move-based coalescing first
	moveCoalesced := applyCoalescing(ctx.Program, coalesceMap, unsafeRegs)

	if moveCoalesced {
		// Rebuild CFG and analyses after changing register assignments.
		builder := NewBuilder(ctx.Program)
		newCFG, err := builder.Build()

		if err != nil {
			return nil, ErrCFGBuildFailed
		}

		cfg = newCFG
		liveness = computeLiveness(cfg)
		unsafeRegs = collectPinnedRegs(ctx.Program)
		interferenceGraph = buildInterferenceGraph(cfg, liveness, ctx.Program.Registers)
	}

	// Then, perform register renumbering based on liveness intervals
	var renumberMap map[int]int
	var renumbered bool

	if isLinearCFG(cfg) {
		renumbered = applyLinearRenumbering(ctx.Program, unsafeRegs)
	} else {
		renumberMap = computeRegisterRenumbering(cfg, liveness, ctx.Program.Registers, unsafeRegs, interferenceGraph)
		renumbered = applyRenumbering(ctx.Program, renumberMap, unsafeRegs)
	}

	modified := folded || moveCoalesced || renumbered

	return &PassResult{
		Modified: modified,
		Metadata: map[string]interface{}{
			"registers_coalesced": len(coalesceMap),
			"coalesce_map":        coalesceMap,
			"renumber_map":        renumberMap,
		},
	}, nil
}

func isLinearCFG(cfg *ControlFlowGraph) bool {
	if cfg.Entry == nil {
		return true
	}

	blocks := cfg.Blocks

	for i, block := range blocks {
		if block == cfg.Exit {
			continue
		}

		if len(block.Predecessors) > 1 || len(block.Successors) > 1 {
			return false
		}

		if len(block.Successors) == 1 {
			next := blocks[i+1]
			if block.Successors[0] != next && block.Successors[0] != cfg.Exit {
				return false
			}
		}
	}

	return true
}
