package core

import (
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

	Dst    vm.Operand
	DstPos int
}

func (l *Loop) DeclareKeyVar(name string, st *SymbolTable) {
	if l.canDeclareVar(name) {
		l.KeyName = name
		l.Key = st.DeclareLocal(name)
	}
}

func (l *Loop) DeclareValueVar(name string, st *SymbolTable) {
	if l.canDeclareVar(name) {
		l.ValueName = name
		l.Value = st.DeclareLocal(name)
	}
}

func (l *Loop) EmitInitialization(alloc *RegisterAllocator, emitter *Emitter) {
	if l.Allocate {
		emitter.EmitAb(vm.OpDataSet, l.Dst, l.Distinct)
		l.DstPos = emitter.Position()
	}

	if l.Iterator == vm.NoopOperand {
		l.Iterator = alloc.Allocate(Temp)
	}

	emitter.EmitIter(l.Iterator, l.Src)

	// JumpPlaceholder is a placeholder for the exit jump position
	l.Jump = emitter.EmitJumpc(vm.OpIterNext, JumpPlaceholder, l.Iterator)

	if l.canBindVar(l.Value) {
		l.EmitValue(l.Value, emitter)
	}

	if l.canBindVar(l.Key) {
		l.EmitKey(l.Key, emitter)
	}
}

func (l *Loop) EmitValue(dst vm.Operand, emitter *Emitter) {
	emitter.EmitIterValue(dst, l.Iterator)
}

func (l *Loop) EmitKey(dst vm.Operand, emitter *Emitter) {
	emitter.EmitIterKey(dst, l.Iterator)
}

func (l *Loop) EmitFinalization(emitter *Emitter) {
	emitter.EmitJump(l.Jump - l.JumpOffset)
	emitter.EmitA(vm.OpClose, l.Iterator)
	emitter.PatchJump(l.Jump)
}

func (l *Loop) canDeclareVar(name string) bool {
	return name != "" && name != IgnorePseudoVariable
}

func (l *Loop) canBindVar(op vm.Operand) bool {
	return op != vm.NoopOperand
}
