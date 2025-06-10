package internal

import (
	"fmt"
	"strings"

	"github.com/MontFerret/ferret/pkg/vm"
)

type LoopType int

const (
	NormalLoop LoopType = iota
	PassThroughLoop
	TemporalLoop
)

type LoopKind int

const (
	ForLoop LoopKind = iota
	WhileLoop
	DoWhileLoop
)

type CollectorType int

const (
	CollectorTypeCounter CollectorType = iota
	CollectorTypeKey
	CollectorTypeKeyCounter
	CollectorTypeKeyGroup
)

type Loop struct {
	Type       LoopType
	Kind       LoopKind
	Distinct   bool
	Allocate   bool
	Jump       int
	JumpOffset int

	Src      vm.Operand
	Iterator vm.Operand

	ValueName string
	Value     vm.Operand
	KeyName   string
	Key       vm.Operand

	Result    vm.Operand
	ResultPos int
}

func (l *Loop) BindKeyVar(name string, st *SymbolTable) {
	l.KeyName = name
	l.Key = st.DeclareLocal(name)
}

func (l *Loop) BindValueVar(name string, st *SymbolTable) {
	l.ValueName = name
	l.Value = st.DeclareLocal(name)
}

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

func (lt *LoopTable) NewLoop(loopType LoopType, kind LoopKind, distinct bool) *Loop {
	parent := lt.Current()
	allocate := parent == nil || parent.Type != PassThroughLoop
	result := vm.NoopOperand

	if allocate && loopType != TemporalLoop {
		result = lt.registers.Allocate(Result)
	} else if parent != nil {
		result = parent.Result
	}

	return &Loop{
		Type:     loopType,
		Kind:     kind,
		Distinct: distinct,
		Result:   result,
		Allocate: allocate,
	}
}

func (lt *LoopTable) DebugView() string {
	var out strings.Builder
	for i, loop := range lt.stack {
		fmt.Fprintf(&out, "Loop[%d]: Type=%v Kind=%v Result=R%d\n", i, loop.Type, loop.Kind, loop.Result)
	}
	return out.String()
}
