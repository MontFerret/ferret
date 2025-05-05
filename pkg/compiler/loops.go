package compiler

import (
	"github.com/MontFerret/ferret/pkg/vm"
)

type (
	LoopType int

	LoopKind int

	Loop struct {
		Type       LoopType
		Kind       LoopKind
		Distinct   bool
		Allocate   bool
		Jump       int
		JumpOffset int
		Src        vm.Operand
		Iterator   vm.Operand
		ValueName  string
		Value      vm.Operand
		KeyName    string
		Key        vm.Operand
		Result     vm.Operand
	}

	LoopTable struct {
		loops     []*Loop
		registers *RegisterAllocator
	}
)

const (
	NormalLoop LoopType = iota
	PassThroughLoop
	TemporalLoop
)

const (
	ForLoop LoopKind = iota
	WhileLoop
	DoWhileLoop
)

func NewLoopTable(registers *RegisterAllocator) *LoopTable {
	return &LoopTable{
		loops:     make([]*Loop, 0),
		registers: registers,
	}
}

func (lt *LoopTable) EnterLoop(loopType LoopType, kind LoopKind, distinct bool) *Loop {
	var allocate bool
	var state vm.Operand

	// top loop
	if len(lt.loops) == 0 {
		allocate = true
	} else if loopType != PassThroughLoop {
		// nested with explicit RETURN expression
		prev := lt.loops[len(lt.loops)-1]
		// if the loop above does not do pass through
		// we allocate a new state for this loop
		allocate = prev.Type != PassThroughLoop
		state = prev.Result
	} else {
		// nested with implicit RETURN expression
		// we reuse the state of the loop above
		state = lt.loops[len(lt.loops)-1].Result
	}

	if allocate {
		state = lt.registers.Allocate(Result)
	}

	lt.loops = append(lt.loops, &Loop{
		Type:     loopType,
		Kind:     kind,
		Distinct: distinct,
		Result:   state,
		Allocate: allocate,
	})

	return lt.loops[len(lt.loops)-1]
}

//func (lt *LoopTable) Fork() *Loop {}

func (lt *LoopTable) Loop() *Loop {
	if len(lt.loops) == 0 {
		return nil
	}

	return lt.loops[len(lt.loops)-1]
}

func (lt *LoopTable) ExitLoop() {
	lt.loops = lt.loops[:len(lt.loops)-1]
}
