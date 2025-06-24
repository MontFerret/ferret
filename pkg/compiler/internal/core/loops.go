package core

import (
	"fmt"
	"strings"

	"github.com/MontFerret/ferret/pkg/vm"
)

type LoopTable struct {
	stack     []*Loop
	registers *RegisterAllocator
}

func NewLoopTable(registers *RegisterAllocator) *LoopTable {
	return &LoopTable{
		stack:     make([]*Loop, 0),
		registers: registers,
	}
}

func (lt *LoopTable) CreateFor(loopType LoopType, src vm.Operand, distinct bool) *Loop {
	parent := lt.Current()
	allocate := parent == nil || parent.Type != PassThroughLoop
	result := vm.NoopOperand

	if loopType != TemporalLoop {
		if allocate {
			result = lt.registers.Allocate(Result)
		} else if parent != nil {
			result = parent.Dst
		}
	}

	return &Loop{
		Type:     loopType,
		Kind:     ForLoop,
		Distinct: distinct,
		Src:      src,
		Dst:      result,
		Allocate: allocate,
	}
}

func (lt *LoopTable) Push(loop *Loop) {
	lt.stack = append(lt.stack, loop)
}

func (lt *LoopTable) Pop() *Loop {
	if len(lt.stack) == 0 {
		return nil
	}
	top := lt.stack[len(lt.stack)-1]
	lt.stack = lt.stack[:len(lt.stack)-1]
	return top
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

func (lt *LoopTable) DebugView() string {
	var out strings.Builder
	for i, loop := range lt.stack {
		fmt.Fprintf(&out, "Loop[%d]: Type=%v Kind=%v Dst=R%d\n", i, loop.Type, loop.Kind, loop.Dst)
	}
	return out.String()
}
