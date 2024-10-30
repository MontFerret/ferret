package compiler

import "github.com/MontFerret/ferret/pkg/runtime"

type (
	Loop struct {
		PassThrough bool
		Register    runtime.Operand
		Position    int
	}

	LoopTable struct {
		loops   []*Loop
		symbols *SymbolTable
	}
)

func NewLoopTable(symbols *SymbolTable) *LoopTable {
	return &LoopTable{
		loops:   make([]*Loop, 0),
		symbols: symbols,
	}
}

func (lt *LoopTable) EnterLoop() {
	lt.loops = append(lt.loops, &Loop{})
}

func (lt *LoopTable) ExitLoop() {
	lt.loops = lt.loops[:len(lt.loops)-1]
}
