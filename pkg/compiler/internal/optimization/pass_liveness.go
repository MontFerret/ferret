package optimization

import (
	"github.com/MontFerret/ferret/v2/pkg/bytecode"
)

const (
	LivenessAnalysisPassName = "liveness-analysis"
	LivenessAnalysis         = "liveness"
)

type (
	// LivenessAnalysisPass performs liveness analysis to determine which registers
	// are live (potentially used) at each program point
	LivenessAnalysisPass struct{}

	// LivenessInfo contains liveness information for a basic block
	LivenessInfo struct {
		LiveIn  map[int]bool // Registers live at block entry
		LiveOut map[int]bool // Registers live at block exit
		Use     map[int]bool // Registers used in block before definition
		Def     map[int]bool // Registers defined in block
	}

	useDefCollector struct {
		uses []int
		defs []int
	}
)

// NewLivenessAnalysisPass creates a new liveness analysis pass
func NewLivenessAnalysisPass() Pass {
	return &LivenessAnalysisPass{}
}

func (p *LivenessAnalysisPass) Requires() []string {
	return []string{}
}

// Name returns the pass name
func (p *LivenessAnalysisPass) Name() string {
	return LivenessAnalysisPassName
}

// Run executes liveness analysis on the program
func (p *LivenessAnalysisPass) Run(c *PassContext) (*PassResult, error) {
	liveness := computeLiveness(c.CFG)

	return &PassResult{
		Modified: false,
		Metadata: map[string]interface{}{
			"liveness": liveness,
		},
	}, nil
}

// computeLiveness performs dataflow analysis to compute liveness information
func computeLiveness(cfg *ControlFlowGraph) map[int]*LivenessInfo {
	info := make(map[int]*LivenessInfo)

	// Initialize liveness info for each block
	for _, block := range cfg.Blocks {
		info[block.ID] = &LivenessInfo{
			LiveIn:  make(map[int]bool),
			LiveOut: make(map[int]bool),
			Use:     make(map[int]bool),
			Def:     make(map[int]bool),
		}

		// Compute Use and Def sets for this block
		for _, inst := range block.Instructions {
			computeUseDefForInstruction(inst, info[block.ID])
		}
	}

	// Iteratively compute LiveIn and LiveOut until convergence
	changed := true
	for changed {
		changed = false

		// Process blocks in reverse order (better for backward dataflow)
		for i := len(cfg.Blocks) - 1; i >= 0; i-- {
			block := cfg.Blocks[i]
			blockInfo := info[block.ID]

			// LiveOut[B] = Union of LiveIn[S] for all successors S of B
			newLiveOut := make(map[int]bool)

			for _, succ := range block.Successors {
				succInfo := info[succ.ID]

				for reg := range succInfo.LiveIn {
					newLiveOut[reg] = true
				}
			}

			// LiveIn[B] = Use[B] Union (LiveOut[B] - Def[B])
			newLiveIn := make(map[int]bool)

			for reg := range blockInfo.Use {
				newLiveIn[reg] = true
			}

			for reg := range newLiveOut {
				if !blockInfo.Def[reg] {
					newLiveIn[reg] = true
				}
			}

			// Check if anything changed
			if !mapsEqual(blockInfo.LiveIn, newLiveIn) || !mapsEqual(blockInfo.LiveOut, newLiveOut) {
				changed = true
				blockInfo.LiveIn = newLiveIn
				blockInfo.LiveOut = newLiveOut
			}
		}
	}

	return info
}

// instructionUseDef returns registers used and defined by an instruction.
// Only registers (non-constants) are included.
func instructionUseDef(inst bytecode.Instruction) (uses []int, defs []int) {
	collector := useDefCollector{}
	opcode := inst.Opcode
	dst, src1, src2 := inst.Operands[0], inst.Operands[1], inst.Operands[2]

	if applyTerminalUseDef(opcode, dst, &collector) ||
		applyMoveLoadUseDef(opcode, dst, src1, src2, &collector) ||
		applyBinaryUseDef(opcode, dst, src1, src2, &collector) ||
		applyUnaryUseDef(opcode, dst, src1, &collector) ||
		applyControlFlowUseDef(opcode, src1, src2, &collector) ||
		applyDatasetUseDef(opcode, dst, src1, src2, &collector) ||
		applyIteratorUseDef(opcode, dst, src1, src2, &collector) ||
		applyCallUseDef(opcode, dst, src1, src2, &collector) ||
		applyStreamUseDef(opcode, dst, src1, src2, &collector) ||
		applyUtilityUseDef(opcode, dst, &collector) {
		return collector.uses, collector.defs
	}

	return collector.uses, collector.defs
}

func (c *useDefCollector) addUse(op bytecode.Operand) {
	if op != bytecode.NoopOperand && op.IsRegister() {
		c.uses = append(c.uses, op.Register())
	}
}

func (c *useDefCollector) addDef(op bytecode.Operand) {
	if op != bytecode.NoopOperand && op.IsRegister() {
		c.defs = append(c.defs, op.Register())
	}
}

func (c *useDefCollector) addRangeUses(start, count int) {
	if count <= 0 || start <= 0 {
		return
	}

	for reg := start; reg < start+count; reg++ {
		c.uses = append(c.uses, reg)
	}
}

func applyTerminalUseDef(opcode bytecode.Opcode, dst bytecode.Operand, collector *useDefCollector) bool {
	if opcode != bytecode.OpReturn {
		return false
	}

	collector.addUse(dst)
	return true
}

func applyMoveLoadUseDef(opcode bytecode.Opcode, dst, src1, src2 bytecode.Operand, collector *useDefCollector) bool {
	switch opcode {
	case bytecode.OpMove, bytecode.OpMoveTracked:
		collector.addUse(src1)
		collector.addDef(dst)
		return true
	case bytecode.OpLoadConst, bytecode.OpLoadParam, bytecode.OpLoadNone, bytecode.OpLoadBool, bytecode.OpLoadZero, bytecode.OpRand,
		bytecode.OpLoadArray, bytecode.OpLoadObject:
		collector.addDef(dst)
		return true
	case bytecode.OpLoadRange:
		collector.addUse(src1)
		collector.addUse(src2)
		collector.addDef(dst)
		return true
	case bytecode.OpLoadAggregateKey:
		collector.addUse(src1)
		collector.addUse(src2)
		collector.addDef(dst)
		return true
	case bytecode.OpMatchLoadPropertyConst:
		collector.addUse(src1)
		collector.addDef(dst)
		return true
	case bytecode.OpAggregateUpdate:
		collector.addUse(dst)
		collector.addUse(src1)
		return true
	case bytecode.OpAggregateGroupUpdate:
		collector.addUse(dst)
		collector.addUse(src1)
		collector.addUse(src2)
		return true
	case bytecode.OpConcat:
		collector.addDef(dst)
		if !src1.IsRegister() {
			return true
		}
		collector.addRangeUses(src1.Register(), int(src2))
		return true
	default:
		return false
	}
}

func applyBinaryUseDef(opcode bytecode.Opcode, dst, src1, src2 bytecode.Operand, collector *useDefCollector) bool {
	if !isBinaryUseDefOpcode(opcode) {
		return false
	}

	collector.addUse(src1)
	collector.addUse(src2)
	collector.addDef(dst)
	return true
}

func isBinaryUseDefOpcode(opcode bytecode.Opcode) bool {
	switch opcode {
	case bytecode.OpAdd, bytecode.OpAddConst, bytecode.OpSub, bytecode.OpMul, bytecode.OpDiv, bytecode.OpMod,
		bytecode.OpCmp,
		bytecode.OpEq, bytecode.OpNe, bytecode.OpGt, bytecode.OpLt, bytecode.OpGte, bytecode.OpLte,
		bytecode.OpAnyEq, bytecode.OpAnyNe, bytecode.OpAnyGt, bytecode.OpAnyGte, bytecode.OpAnyLt, bytecode.OpAnyLte,
		bytecode.OpAnyIn,
		bytecode.OpNoneEq, bytecode.OpNoneNe, bytecode.OpNoneGt, bytecode.OpNoneGte, bytecode.OpNoneLt, bytecode.OpNoneLte,
		bytecode.OpNoneIn,
		bytecode.OpAllEq, bytecode.OpAllNe, bytecode.OpAllGt, bytecode.OpAllGte, bytecode.OpAllLt, bytecode.OpAllLte,
		bytecode.OpAllIn,
		bytecode.OpIn, bytecode.OpLike, bytecode.OpRegexp,
		bytecode.OpLoadIndex, bytecode.OpLoadIndexOptional, bytecode.OpLoadIndexConst, bytecode.OpLoadIndexOptionalConst,
		bytecode.OpLoadKey, bytecode.OpLoadKeyOptional, bytecode.OpLoadKeyConst, bytecode.OpLoadKeyOptionalConst,
		bytecode.OpLoadProperty, bytecode.OpLoadPropertyOptional, bytecode.OpLoadPropertyConst, bytecode.OpLoadPropertyOptionalConst,
		bytecode.OpQuery:
		return true
	default:
		return false
	}
}

func applyUnaryUseDef(opcode bytecode.Opcode, dst, src1 bytecode.Operand, collector *useDefCollector) bool {
	switch opcode {
	case bytecode.OpIncr, bytecode.OpDecr:
		collector.addUse(dst)
		collector.addDef(dst)
		return true
	case bytecode.OpCastBool, bytecode.OpNegate, bytecode.OpNot, bytecode.OpFlipPositive, bytecode.OpFlipNegative, bytecode.OpLength, bytecode.OpType, bytecode.OpFlatten, bytecode.OpExists:
		collector.addUse(src1)
		collector.addDef(dst)
		return true
	default:
		return false
	}
}

func applyControlFlowUseDef(opcode bytecode.Opcode, src1, src2 bytecode.Operand, collector *useDefCollector) bool {
	switch opcode {
	case bytecode.OpJumpIfFalse, bytecode.OpJumpIfTrue, bytecode.OpJumpIfNone,
		bytecode.OpJumpIfNeConst, bytecode.OpJumpIfEqConst, bytecode.OpJumpIfMissingPropertyConst:
		collector.addUse(src1)
		return true
	case bytecode.OpJumpIfNe, bytecode.OpJumpIfEq, bytecode.OpJumpIfMissingProperty:
		collector.addUse(src1)
		collector.addUse(src2)
		return true
	case bytecode.OpJump, bytecode.OpFail:
		return true
	default:
		return false
	}
}

func applyDatasetUseDef(opcode bytecode.Opcode, dst, src1, src2 bytecode.Operand, collector *useDefCollector) bool {
	switch opcode {
	case bytecode.OpDataSet, bytecode.OpDataSetCollector, bytecode.OpDataSetSorter, bytecode.OpDataSetMultiSorter:
		collector.addDef(dst)
		return true
	case bytecode.OpCounterInc:
		collector.addUse(dst)
		return true
	case bytecode.OpPush, bytecode.OpArrayPush:
		collector.addUse(dst)
		collector.addUse(src1)
		return true
	case bytecode.OpPushKV, bytecode.OpObjectSet, bytecode.OpObjectSetConst:
		collector.addUse(dst)
		collector.addUse(src1)
		collector.addUse(src2)
		return true
	default:
		return false
	}
}

func applyIteratorUseDef(opcode bytecode.Opcode, dst, src1, src2 bytecode.Operand, collector *useDefCollector) bool {
	switch opcode {
	case bytecode.OpIter, bytecode.OpIterValue, bytecode.OpIterKey:
		collector.addUse(src1)
		collector.addDef(dst)
		return true
	case bytecode.OpIterLimit, bytecode.OpIterSkip:
		collector.addUse(src1)
		collector.addUse(src2)
		collector.addDef(src1)
		return true
	case bytecode.OpIterNext:
		collector.addUse(src1)
		return true
	default:
		return false
	}
}

func applyCallUseDef(opcode bytecode.Opcode, dst, src1, src2 bytecode.Operand, collector *useDefCollector) bool {
	switch opcode {
	case bytecode.OpHCall, bytecode.OpProtectedHCall, bytecode.OpCall, bytecode.OpProtectedCall, bytecode.OpTailCall:
		collector.addUse(dst)
		bytecode.VisitCallArgumentRegisters(opcode, src1, src2, func(reg int) {
			collector.uses = append(collector.uses, reg)
		})
		if opcode != bytecode.OpTailCall {
			collector.addDef(dst)
		}
		return true
	default:
		return false
	}
}

func applyStreamUseDef(opcode bytecode.Opcode, dst, src1, src2 bytecode.Operand, collector *useDefCollector) bool {
	switch opcode {
	case bytecode.OpStream, bytecode.OpDispatch:
		collector.addUse(dst)
		collector.addUse(src1)
		collector.addUse(src2)
		collector.addDef(dst)
		return true
	case bytecode.OpStreamIter:
		collector.addUse(src1)
		collector.addUse(src2)
		collector.addDef(dst)
		return true
	default:
		return false
	}
}

func applyUtilityUseDef(opcode bytecode.Opcode, dst bytecode.Operand, collector *useDefCollector) bool {
	switch opcode {
	case bytecode.OpClose:
		collector.addUse(dst)
		collector.addDef(dst)
		return true
	case bytecode.OpSleep:
		collector.addUse(dst)
		return true
	default:
		return false
	}
}

// computeUseDefForInstruction updates Use and Def sets for an instruction
func computeUseDefForInstruction(inst bytecode.Instruction, info *LivenessInfo) {
	uses, defs := instructionUseDef(inst)

	// Use-before-def within the block
	for _, reg := range uses {
		if !info.Def[reg] {
			info.Use[reg] = true
		}
	}

	for _, reg := range defs {
		info.Def[reg] = true
	}
}

// mapsEqual checks if two register maps are equal
func mapsEqual(a, b map[int]bool) bool {
	if len(a) != len(b) {
		return false
	}

	for k := range a {
		if !b[k] {
			return false
		}
	}

	return true
}
