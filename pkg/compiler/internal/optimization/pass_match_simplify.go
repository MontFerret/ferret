package optimization

import (
	"sort"
	"strconv"
	"strings"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

const MatchSimplifyPassName = "match-simplify"

type MatchSimplifyPass struct{}

func NewMatchSimplifyPass() Pass {
	return &MatchSimplifyPass{}
}

func (p *MatchSimplifyPass) Name() string {
	return MatchSimplifyPassName
}

func (p *MatchSimplifyPass) Requires() []string {
	return []string{}
}

type matchInfo struct {
	id    int
	start int
	end   int
	next  map[int]int
}

func (p *MatchSimplifyPass) Run(ctx *PassContext) (*PassResult, error) {
	if ctx == nil || ctx.Program == nil || ctx.CFG == nil {
		return &PassResult{Modified: false}, nil
	}

	prog := ctx.Program
	code := prog.Bytecode

	if len(code) == 0 || len(prog.Metadata.Labels) == 0 {
		return &PassResult{Modified: false}, nil
	}

	matches := collectMatchInfos(prog.Metadata.Labels)
	if len(matches) == 0 {
		return &PassResult{Modified: false}, nil
	}

	modified := false
	var keep []bool

	for _, info := range matches {
		if info.start < 0 || info.end < 0 || len(info.next) == 0 {
			continue
		}

		armIndices := make([]int, 0, len(info.next))
		for idx := range info.next {
			armIndices = append(armIndices, idx)
		}
		sort.Ints(armIndices)

		nextPositions := make([]int, len(armIndices))
		for i, idx := range armIndices {
			nextPositions[i] = info.next[idx]
		}

		if info.start >= nextPositions[0] {
			continue
		}
		if nextPositions[len(nextPositions)-1] >= info.end {
			continue
		}
		valid := true
		for i := 1; i < len(nextPositions); i++ {
			if nextPositions[i-1] >= nextPositions[i] {
				valid = false
				break
			}
		}
		if !valid {
			continue
		}

		armCount := len(nextPositions)
		armStarts := make([]int, armCount)
		armNexts := make([]int, armCount)
		for i := 0; i < armCount; i++ {
			if i == 0 {
				armStarts[i] = info.start
			} else {
				armStarts[i] = nextPositions[i-1]
			}
			armNexts[i] = nextPositions[i]
		}

		if armStarts[0] < 0 || armStarts[0] >= len(code) {
			continue
		}

		firstInst := code[armStarts[0]]
		if firstInst.Opcode != bytecode.OpJumpIfNeConst || !firstInst.Operands[1].IsRegister() {
			continue
		}
		if int(firstInst.Operands[0]) != armNexts[0] {
			continue
		}

		scrReg := firstInst.Operands[1].Register()
		scrVal, ok := findConstDefValue(code, prog, armStarts[0], scrReg)
		if !ok {
			continue
		}

		foldable := true
		for i := 0; i < armCount; i++ {
			start := armStarts[i]
			next := armNexts[i]
			if start < 0 || start >= len(code) {
				foldable = false
				break
			}
			inst := code[start]
			if inst.Opcode != bytecode.OpJumpIfNeConst || !inst.Operands[1].IsRegister() || inst.Operands[1].Register() != scrReg {
				foldable = false
				break
			}
			if int(inst.Operands[0]) != next {
				foldable = false
				break
			}
			if hasJumpToTarget(code, start+1, next, next) {
				foldable = false
				break
			}
		}

		if !foldable {
			continue
		}

		matchIndex := -1
		for i := 0; i < armCount; i++ {
			inst := code[armStarts[i]]
			if !inst.Operands[2].IsConstant() {
				continue
			}
			val := prog.Constants[inst.Operands[2].Constant()]
			if runtime.CompareValues(scrVal, val) == 0 {
				matchIndex = i
				break
			}
		}

		for i := 0; i < armCount; i++ {
			start := armStarts[i]
			inst := code[start]
			if matchIndex == i {
				regOp := bytecode.NewRegister(scrReg)
				inst.Opcode = bytecode.OpMove
				inst.Operands[0] = regOp
				inst.Operands[1] = regOp
				inst.Operands[2] = 0
				code[start] = inst
				modified = true
				continue
			}

			if matchIndex == -1 || i < matchIndex {
				inst.Opcode = bytecode.OpJump
				inst.Operands[0] = bytecode.Operand(armNexts[i])
				inst.Operands[1] = 0
				inst.Operands[2] = 0
				code[start] = inst
				modified = true
			}
		}

		if matchIndex < 0 {
			continue
		}

		targets := collectJumpTargets(code, prog.CatchTable)
		removeRanges := make([][2]int, 0, armCount)
		redirectTargets := make(map[int]bool)

		for i := matchIndex + 1; i < armCount; i++ {
			start := armStarts[i]
			end := armNexts[i]
			removeRanges = append(removeRanges, [2]int{start, end})
			redirectTargets[start] = true
		}

		unsafe := false
		for _, r := range removeRanges {
			for target := range targets {
				if target >= r[0] && target < r[1] && target != r[0] {
					unsafe = true
					break
				}
			}
			if unsafe {
				break
			}
		}

		if unsafe || len(removeRanges) == 0 {
			continue
		}

		for i, inst := range code {
			if !isJumpOpcode(inst.Opcode) {
				continue
			}
			target := int(inst.Operands[0])
			if redirectTargets[target] {
				inst.Operands[0] = bytecode.Operand(info.end)
				code[i] = inst
				modified = true
			}
		}

		if keep == nil {
			keep = make([]bool, len(code))
			for i := range keep {
				keep[i] = true
			}
		}

		for _, r := range removeRanges {
			if r[0] < 0 || r[0] >= len(keep) || r[1] <= r[0] {
				continue
			}
			end := r[1]
			if end > len(keep) {
				end = len(keep)
			}
			for i := r[0]; i < end; i++ {
				keep[i] = false
			}
		}
	}

	if keep == nil {
		return &PassResult{Modified: modified}, nil
	}

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

	for i := range newCode {
		inst := newCode[i]
		if !isJumpOpcode(inst.Opcode) {
			continue
		}
		oldTarget := int(inst.Operands[0])
		if oldTarget < 0 || oldTarget >= len(indexMap) {
			continue
		}
		newTarget := remapIndexForward(indexMap, keep, oldTarget)
		if newTarget >= 0 && newTarget != oldTarget {
			inst.Operands[0] = bytecode.Operand(newTarget)
			newCode[i] = inst
		}
	}

	prog.Bytecode = newCode
	remapDebugSpans(prog, keep)
	remapLabels(prog, indexMap)
	remapCatchTable(prog, indexMap, keep)

	return &PassResult{Modified: true}, nil
}

func collectMatchInfos(labels map[int]string) []matchInfo {
	if len(labels) == 0 {
		return nil
	}

	byID := make(map[int]*matchInfo)

	for pos, name := range labels {
		id, kind, arm, ok := parseMatchLabel(name)
		if !ok {
			continue
		}
		info := byID[id]
		if info == nil {
			info = &matchInfo{id: id, start: -1, end: -1, next: make(map[int]int)}
			byID[id] = info
		}
		switch kind {
		case "start":
			info.start = pos
		case "end":
			info.end = pos
		case "next":
			info.next[arm] = pos
		}
	}

	ids := make([]int, 0, len(byID))
	for id := range byID {
		ids = append(ids, id)
	}
	sort.Ints(ids)

	out := make([]matchInfo, 0, len(ids))
	for _, id := range ids {
		out = append(out, *byID[id])
	}
	return out
}

func parseMatchLabel(name string) (int, string, int, bool) {
	parts := strings.Split(name, ".")
	if len(parts) < 3 || parts[0] != "match" {
		return 0, "", 0, false
	}
	id, err := strconv.Atoi(parts[1])
	if err != nil {
		return 0, "", 0, false
	}
	kind := parts[2]
	if kind == "start" || kind == "end" {
		if len(parts) != 3 {
			return 0, "", 0, false
		}
		return id, kind, 0, true
	}
	if kind == "next" {
		if len(parts) != 4 {
			return 0, "", 0, false
		}
		arm, err := strconv.Atoi(parts[3])
		if err != nil {
			return 0, "", 0, false
		}
		return id, kind, arm, true
	}
	return 0, "", 0, false
}

func findConstDefValue(code []bytecode.Instruction, prog *bytecode.Program, before int, reg int) (runtime.Value, bool) {
	for i := before - 1; i >= 0; i-- {
		_, defs := instructionUseDef(code[i])
		if !regIn(defs, reg) {
			continue
		}
		switch code[i].Opcode {
		case bytecode.OpLoadConst:
			op := code[i].Operands[1]
			if !op.IsConstant() || prog == nil {
				return nil, false
			}
			idx := op.Constant()
			if idx < 0 || idx >= len(prog.Constants) {
				return nil, false
			}
			return prog.Constants[idx], true
		case bytecode.OpLoadNone:
			return runtime.None, true
		case bytecode.OpLoadBool:
			if code[i].Operands[1] == 1 {
				return runtime.True, true
			}
			return runtime.False, true
		case bytecode.OpLoadZero:
			return runtime.ZeroInt, true
		default:
			return nil, false
		}
	}
	return nil, false
}

func hasJumpToTarget(code []bytecode.Instruction, start, end, target int) bool {
	if start < 0 {
		start = 0
	}
	if end > len(code) {
		end = len(code)
	}
	for i := start; i < end; i++ {
		inst := code[i]
		if !isJumpOpcode(inst.Opcode) {
			continue
		}
		if int(inst.Operands[0]) == target {
			return true
		}
	}
	return false
}
