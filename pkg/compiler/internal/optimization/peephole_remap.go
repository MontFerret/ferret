package optimization

import (
	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/file"
)

func applyPeepholeCompactionAndRemap(state *peepholeRunState) {
	newCode, indexMap := compactPeepholeInstructions(state.prog.Bytecode, state.keep)
	remapPeepholeJumps(newCode, indexMap, state.keep, state.bytecodeLen)

	state.prog.Bytecode = newCode
	remapAggregateSelectorSlots(state.prog, state.keep)
	remapMatchFailTargets(state.prog, indexMap, state.keep)
	remapDebugSpans(state.prog, state.keep)
	remapLabels(state.prog, indexMap)
	remapUdfEntries(state.prog, indexMap, state.keep)
	remapCatchTable(state.prog, indexMap, state.keep)
}

func compactPeepholeInstructions(code []bytecode.Instruction, keep []bool) ([]bytecode.Instruction, []int) {
	indexMap := make([]int, len(code))
	newCode := make([]bytecode.Instruction, 0, len(code))

	for i, inst := range code {
		if !keep[i] {
			indexMap[i] = -1
			continue
		}

		indexMap[i] = len(newCode)
		newCode = append(newCode, inst)
	}

	return newCode, indexMap
}

func remapPeepholeJumps(code []bytecode.Instruction, indexMap []int, keep []bool, oldCodeLen int) {
	for i := range code {
		inst := code[i]
		if !isJumpOpcode(inst.Opcode) {
			continue
		}
		if inst.Opcode == bytecode.OpMatchLoadPropertyConst {
			continue
		}

		oldTarget := int(inst.Operands[0])
		if oldTarget < 0 || oldTarget >= oldCodeLen {
			continue
		}

		newTarget := remapIndexForward(indexMap, keep, oldTarget)
		if newTarget >= 0 && newTarget != oldTarget {
			inst.Operands[0] = bytecode.Operand(newTarget)
			code[i] = inst
		}
	}
}

func remapIndexForward(indexMap []int, keep []bool, old int) int {
	if old < 0 || old >= len(indexMap) {
		return -1
	}
	if indexMap[old] >= 0 {
		return indexMap[old]
	}
	for i := old + 1; i < len(indexMap); i++ {
		if keep[i] {
			return indexMap[i]
		}
	}
	return -1
}

func remapIndexBackward(indexMap []int, keep []bool, old int) int {
	if old < 0 || old >= len(indexMap) {
		return -1
	}
	if indexMap[old] >= 0 {
		return indexMap[old]
	}
	for i := old - 1; i >= 0; i-- {
		if keep[i] {
			return indexMap[i]
		}
	}
	return -1
}

func remapCatchTable(prog *bytecode.Program, indexMap []int, keep []bool) {
	if prog == nil || len(prog.CatchTable) == 0 {
		return
	}

	updated := make([]bytecode.Catch, 0, len(prog.CatchTable))
	for _, entry := range prog.CatchTable {
		start := remapIndexForward(indexMap, keep, entry[0])
		end := remapIndexBackward(indexMap, keep, entry[1])
		if start < 0 || end < 0 || start > end {
			continue
		}
		jump := entry[2]
		if jump >= 0 {
			jump = remapIndexForward(indexMap, keep, jump)
		}
		updated = append(updated, bytecode.Catch{start, end, jump})
	}

	prog.CatchTable = updated
}

func remapDebugSpans(prog *bytecode.Program, keep []bool) {
	if prog == nil || len(prog.Metadata.DebugSpans) == 0 {
		return
	}
	if len(prog.Metadata.DebugSpans) != len(keep) {
		return
	}
	updated := make([]file.Span, 0, len(prog.Metadata.DebugSpans))
	for i, span := range prog.Metadata.DebugSpans {
		if keep[i] {
			updated = append(updated, span)
		}
	}
	prog.Metadata.DebugSpans = updated
}

func remapAggregateSelectorSlots(prog *bytecode.Program, keep []bool) {
	if prog == nil || len(prog.Metadata.AggregateSelectorSlots) == 0 {
		return
	}

	if len(prog.Metadata.AggregateSelectorSlots) != len(keep) {
		return
	}

	updated := make([]int, 0, len(prog.Metadata.AggregateSelectorSlots))

	for i, slot := range prog.Metadata.AggregateSelectorSlots {
		if keep[i] {
			updated = append(updated, slot)
		}
	}

	prog.Metadata.AggregateSelectorSlots = updated
}

func remapMatchFailTargets(prog *bytecode.Program, indexMap []int, keep []bool) {
	if prog == nil || len(prog.Metadata.MatchFailTargets) == 0 {
		return
	}

	if len(prog.Metadata.MatchFailTargets) != len(keep) {
		return
	}

	updated := make([]int, 0, len(prog.Metadata.MatchFailTargets))

	for i, target := range prog.Metadata.MatchFailTargets {
		if !keep[i] {
			continue
		}

		if target >= 0 {
			target = remapIndexForward(indexMap, keep, target)
		}

		updated = append(updated, target)
	}

	prog.Metadata.MatchFailTargets = updated
}

func remapLabels(prog *bytecode.Program, indexMap []int) {
	if prog == nil || len(prog.Metadata.Labels) == 0 {
		return
	}
	updated := make(map[int]string, len(prog.Metadata.Labels))
	for old, name := range prog.Metadata.Labels {
		if old < 0 || old >= len(indexMap) {
			continue
		}
		if idx := indexMap[old]; idx >= 0 {
			updated[idx] = name
		}
	}
	prog.Metadata.Labels = updated
}

func remapUdfEntries(prog *bytecode.Program, indexMap []int, keep []bool) {
	if prog == nil || len(prog.Functions.UserDefined) == 0 {
		return
	}

	for i := range prog.Functions.UserDefined {
		entry := prog.Functions.UserDefined[i].Entry
		if entry < 0 || entry >= len(indexMap) {
			continue
		}

		newEntry := indexMap[entry]
		if newEntry < 0 {
			newEntry = remapIndexForward(indexMap, keep, entry)
		}

		if newEntry >= 0 {
			prog.Functions.UserDefined[i].Entry = newEntry
		}
	}
}
