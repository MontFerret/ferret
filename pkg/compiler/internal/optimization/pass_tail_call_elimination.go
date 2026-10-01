package optimization

import (
	"sort"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
)

const TailCallEliminationPassName = "tail-call-elimination"

// TailCallEliminationPass replaces proven-safe UDF calls in return position.
// It retains trailing returns so all instruction positions remain stable.
type TailCallEliminationPass struct{}

func NewTailCallEliminationPass() Pass {
	return &TailCallEliminationPass{}
}

func (p *TailCallEliminationPass) Name() string {
	return TailCallEliminationPassName
}

func (p *TailCallEliminationPass) Requires() []string {
	return nil
}

func (p *TailCallEliminationPass) Run(ctx *PassContext) (*PassResult, error) {
	result := &PassResult{}
	if ctx == nil || ctx.Program == nil || ctx.CFG == nil || len(ctx.Program.Functions.UserDefined) == 0 {
		return result, nil
	}

	program := ctx.Program
	code := program.Bytecode
	entries := make([]int, 0, len(program.Functions.UserDefined))
	for _, fn := range program.Functions.UserDefined {
		if fn.Entry >= 0 && fn.Entry < len(code) {
			entries = append(entries, fn.Entry)
		}
	}

	sort.Ints(entries)
	blocks := buildInstructionBlockMap(ctx.CFG, len(code))

	for i, start := range entries {
		end := len(code)
		if i+1 < len(entries) {
			end = entries[i+1]
		}

		cells := p.localCellAliases(code[start:end])

		for pc := start; pc+1 < end; pc++ {
			call, ret := code[pc], code[pc+1]
			if call.Opcode != bytecode.OpCall || ret.Opcode != bytecode.OpReturn || !call.Operands[0].IsRegister() || ret.Operands[0] != call.Operands[0] {
				continue
			}

			if blocks[pc] == nil || blocks[pc] != blocks[pc+1] || p.coveredByCatch(program, pc, pc+1) || p.forwardsLocalCell(call, cells) {
				continue
			}

			code[pc].Opcode = bytecode.OpTailCall
			result.Modified = true
		}
	}

	return result, nil
}

func (p *TailCallEliminationPass) coveredByCatch(program *bytecode.Program, start, end int) bool {
	for _, catch := range program.CatchTable {
		if catch[0] <= end && start <= catch[1] {
			return true
		}
	}

	return false
}

func (p *TailCallEliminationPass) forwardsLocalCell(call bytecode.Instruction, cells map[int]bool) bool {
	if len(cells) == 0 {
		return false
	}

	unsafe := false
	bytecode.VisitCallArgumentRegisters(call.Opcode, call.Operands[1], call.Operands[2], func(reg int) {
		unsafe = unsafe || cells[reg]
	})

	return unsafe
}

func (p *TailCallEliminationPass) localCellAliases(code []bytecode.Instruction) map[int]bool {
	var cells map[int]bool
	var pending []int
	for _, inst := range code {
		if inst.Opcode == bytecode.OpMakeCell && inst.Operands[0].IsRegister() {
			if cells == nil {
				cells = make(map[int]bool)
			}

			reg := inst.Operands[0].Register()
			if !cells[reg] {
				cells[reg] = true
				pending = append(pending, reg)
			}
		}
	}

	if len(pending) == 0 {
		return nil
	}

	// Frame replacement deletes locally owned cells. Track possible handle
	// aliases across all paths, without trusting instruction order or killing
	// aliases on overwrite. Outer-frame captures have no local MakeCell seed.
	aliases := make(map[int][]int)
	for _, inst := range code {
		if (inst.Opcode == bytecode.OpMove || inst.Opcode == bytecode.OpMoveTracked) && inst.Operands[0].IsRegister() && inst.Operands[1].IsRegister() {
			src := inst.Operands[1].Register()
			aliases[src] = append(aliases[src], inst.Operands[0].Register())
		}
	}

	for len(pending) > 0 {
		reg := pending[len(pending)-1]
		pending = pending[:len(pending)-1]
		for _, alias := range aliases[reg] {
			if !cells[alias] {
				cells[alias] = true
				pending = append(pending, alias)
			}
		}
	}

	return cells
}
