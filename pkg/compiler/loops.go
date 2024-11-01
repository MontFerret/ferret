package compiler

import "github.com/MontFerret/ferret/pkg/runtime"

type (
	Loop struct {
		PassThrough bool
		Distinct    bool
		Result      runtime.Operand
		Iterator    runtime.Operand
		Allocated   bool
		Position    int
	}

	LoopTable struct {
		loops     []*Loop
		registers *RegisterAllocator
	}
)

func NewLoopTable(registers *RegisterAllocator) *LoopTable {
	return &LoopTable{
		loops:     make([]*Loop, 0),
		registers: registers,
	}
}

func (lt *LoopTable) EnterLoop(position int, passThrough, distinct bool) {
	var allocate bool
	var state runtime.Operand

	// top loop
	if len(lt.loops) == 0 {
		allocate = true
		state = lt.registers.Allocate(Result)
	} else if !passThrough {
		// nested with explicit RETURN expression
		prev := lt.loops[len(lt.loops)-1]
		// if the loop above does not do pass through
		// we allocate a new state for this loop
		allocate = !prev.PassThrough
		state = prev.Result
	}

	lt.loops = append(lt.loops, &Loop{
		PassThrough: passThrough,
		Distinct:    distinct,
		Result:      state,
		Allocated:   allocate,
		Position:    position,
	})
}

func (lt *LoopTable) Loop() *Loop {
	if len(lt.loops) == 0 {
		return nil
	}

	return lt.loops[len(lt.loops)-1]
}

func (lt *LoopTable) ExitLoop() {
	lt.loops = lt.loops[:len(lt.loops)-1]
}
