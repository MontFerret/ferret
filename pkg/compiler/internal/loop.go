package internal

import "github.com/MontFerret/ferret/pkg/vm"

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

func (l *Loop) DeclareKeyVar(name string, st *SymbolTable) {
	if l.canBindVar(name) {
		l.KeyName = name
		l.Key = st.DeclareLocal(name)
	}
}

func (l *Loop) DeclareValueVar(name string, st *SymbolTable) {
	if l.canBindVar(name) {
		l.ValueName = name
		l.Value = st.DeclareLocal(name)
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

func (l *Loop) canBindVar(name string) bool {
	return name != "" && name != ignorePseudoVariable
}
