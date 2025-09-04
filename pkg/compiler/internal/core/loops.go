package core

import (
	"fmt"

	"github.com/MontFerret/ferret/pkg/vm"
)

type LoopTable struct {
	stack     []*Loop
	registers *RegisterAllocator
	ordinals  map[int]int
}

func NewLoopTable(registers *RegisterAllocator) *LoopTable {
	return &LoopTable{
		stack:     make([]*Loop, 0),
		registers: registers,
		ordinals:  make(map[int]int),
	}
}

func (lt *LoopTable) NewForInLoop(loopType LoopType, distinct bool) *Loop {
	return lt.NewLoop(ForInLoop, loopType, distinct)
}

func (lt *LoopTable) NewForWhileLoop(loopType LoopType, distinct bool) *Loop {
	return lt.NewLoop(ForWhileLoop, loopType, distinct)
}

func (lt *LoopTable) NewLoop(kind LoopKind, loopType LoopType, distinct bool) *Loop {
	parent := lt.Current()
	allocate := parent == nil || parent.Type != PassThroughLoop
	result := vm.NoopOperand

	if loopType != TemporalLoop {
		if allocate {
			result = lt.registers.Allocate(Result)
		} else {
			result = parent.Dst
		}
	}

	return &Loop{
		Type:     loopType,
		Kind:     kind,
		Distinct: distinct,
		Dst:      result,
		Allocate: allocate,
	}
}

func (lt *LoopTable) Push(loop *Loop) {
	lt.stack = append(lt.stack, loop)
	loop.LabelBase = lt.NextBase()
}

func (lt *LoopTable) Pop() *Loop {
	if len(lt.stack) == 0 {
		return nil
	}
	top := lt.stack[len(lt.stack)-1]
	lt.stack = lt.stack[:len(lt.stack)-1]
	return top
}

func (lt *LoopTable) FindParent(pos int) *Loop {
	for i := pos - 1; i >= 0; i-- {
		loop := lt.stack[i]

		if loop.Allocate {
			return loop
		}
	}

	return nil
}

func (lt *LoopTable) Current() *Loop {
	if len(lt.stack) == 0 {
		return nil
	}
	return lt.stack[len(lt.stack)-1]
}

func (lt *LoopTable) Depth() int {
	return len(lt.stack)
}

func (lt *LoopTable) NextBase() string {
	depth := lt.Depth()
	lt.ordinals[depth]++
	return fmt.Sprintf("%d.%d", depth, lt.ordinals[depth])
}
