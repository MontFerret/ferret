package core

import (
	"fmt"

	"github.com/MontFerret/ferret/pkg/vm"
)

type RegisterType int

const (
	Temp RegisterType = iota
	Var
	State
	Result
)

type RegisterSequence []vm.Operand

type registerInfo struct {
	typ       RegisterType
	allocated bool
}

type RegisterAllocator struct {
	next     vm.Operand
	freelist map[RegisterType][]vm.Operand
	all      map[vm.Operand]*registerInfo
}

func NewRegisterAllocator() *RegisterAllocator {
	return &RegisterAllocator{
		next:     vm.NoopOperand + 1,
		freelist: make(map[RegisterType][]vm.Operand),
		all:      make(map[vm.Operand]*registerInfo),
	}
}

func (ra *RegisterAllocator) Allocate(typ RegisterType) vm.Operand {
	// Reuse from freelist
	if regs, ok := ra.freelist[typ]; ok && len(regs) > 0 {
		last := len(regs) - 1
		reg := regs[last]
		ra.freelist[typ] = regs[:last]
		ra.all[reg].allocated = true
		return reg
	}

	// New value
	reg := ra.next
	ra.next++

	ra.all[reg] = &registerInfo{
		typ:       typ,
		allocated: true,
	}

	return reg
}

func (ra *RegisterAllocator) Free(reg vm.Operand) {
	//info, ok := ra.all[state]
	//if !ok || !info.allocated {
	//	return // double-free or unknown
	//}
	//
	//info.allocated = false
	//ra.freelist[info.typ] = append(ra.freelist[info.typ], state)
}

func (ra *RegisterAllocator) AllocateSequence(count int) RegisterSequence {
	// Try to find a contiguous free block
	start, found := ra.findContiguousFree(count)

	if found {
		seq := make(RegisterSequence, count)
		for i := 0; i < count; i++ {
			reg := start + vm.Operand(i)
			info, ok := ra.all[reg]
			if !ok {
				info = &registerInfo{}
				ra.all[reg] = info
			}
			info.typ = Temp
			info.allocated = true
			seq[i] = reg
		}
		return seq
	}

	// Otherwise, allocate new block
	start = ra.next
	seq := make(RegisterSequence, count)

	for i := 0; i < count; i++ {
		reg := start + vm.Operand(i)
		ra.all[reg] = &registerInfo{
			typ:       Temp,
			allocated: true,
		}
		seq[i] = reg
	}

	ra.next += vm.Operand(count)
	return seq
}

func (ra *RegisterAllocator) FreeSequence(seq RegisterSequence) {
	for _, reg := range seq {
		ra.Free(reg)
	}
}

func (ra *RegisterAllocator) findContiguousFree(count int) (vm.Operand, bool) {
	if count <= 0 {
		return 0, false
	}

	limit := ra.next
	for start := vm.NoopOperand + 1; start+vm.Operand(count) <= limit; start++ {
		ok := true
		for i := 0; i < count; i++ {
			reg := start + vm.Operand(i)
			info := ra.all[reg]
			if info == nil || info.allocated {
				ok = false
				break
			}
		}
		if ok {
			return start, true
		}
	}

	return 0, false
}

func (ra *RegisterAllocator) Size() int {
	return int(ra.next)
}

func (ra *RegisterAllocator) DebugView() string {
	out := ""
	for reg := vm.NoopOperand + 1; reg < ra.next; reg++ {
		info := ra.all[reg]
		status := "FREE"
		typ := "?"
		if info != nil {
			if info.allocated {
				status = "USED"
			}
			typ = fmt.Sprintf("%v", info.typ)
		}
		out += fmt.Sprintf("R%d: %s [%s]\n", reg, status, typ)
	}
	return out
}
