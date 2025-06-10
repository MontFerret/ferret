package internal

import "github.com/MontFerret/ferret/pkg/vm"

type FuncContext struct {
	Emitter    *Emitter
	Registers  *RegisterAllocator
	Symbols    *SymbolTable
	Loops      *LoopTable
	CatchTable []vm.Catch
}

func NewFuncContext() *FuncContext {
	registers := NewRegisterAllocator()
	return &FuncContext{
		Emitter:    NewEmitter(),
		Registers:  registers,
		Symbols:    NewSymbolTable(registers),
		Loops:      NewLoopTable(registers),
		CatchTable: make([]vm.Catch, 0),
	}
}
