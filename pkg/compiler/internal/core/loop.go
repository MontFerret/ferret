package core

import (
	"strconv"

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
	ForInLoop LoopKind = iota
	ForWhileLoop
	DoWhileLoop
)

type Loop struct {
	Kind     LoopKind
	Type     LoopType
	Distinct bool
	Allocate bool

	StartLabel Label
	JumpLabel  Label
	EndLabel   Label

	Src      vm.Operand
	SrcFn    func() vm.Operand
	Iterator vm.Operand

	ValueName string
	Value     vm.Operand
	KeyName   string
	Key       vm.Operand

	Dst vm.Operand
}

func (l *Loop) DeclareKeyVar(name string, st *SymbolTable) bool {
	if l.canDeclareVar(name) {
		reg, ok := st.DeclareLocal(name, TypeUnknown)

		if !ok {
			return false
		}

		l.Key = reg
		l.KeyName = name
	}

	return true
}

func (l *Loop) DeclareValueVar(name string, st *SymbolTable) bool {
	if l.canDeclareVar(name) {
		reg, ok := st.DeclareLocal(name, TypeUnknown)

		if !ok {
			return false
		}

		l.Value = reg
		l.ValueName = name
	}

	return true
}

func (l *Loop) EmitInitialization(alloc *RegisterAllocator, emitter *Emitter, depth int) {
	name := strconv.Itoa(depth)
	l.StartLabel = emitter.NewLabel("loop", name, "start")
	l.JumpLabel = emitter.NewLabel("loop", name, "jump")
	l.EndLabel = emitter.NewLabel("loop", name, "end")

	emitter.MarkLabel(l.StartLabel)

	if l.Allocate {
		emitter.EmitAb(vm.OpDataSet, l.Dst, l.Distinct)
	}

	if l.Kind == ForInLoop {
		l.emitForInLoopIteration(alloc, emitter)
	} else {
		l.emitForWhileLoopIteration(alloc, emitter)
	}

	if l.canBindVar(l.Value) {
		l.EmitValue(l.Value, emitter)
	}

	if l.canBindVar(l.Key) {
		l.EmitKey(l.Key, emitter)
	}
}

func (l *Loop) EmitValue(dst vm.Operand, emitter *Emitter) {
	if l.Kind == ForInLoop {
		emitter.EmitIterValue(dst, l.Iterator)
	}
}

func (l *Loop) EmitKey(dst vm.Operand, emitter *Emitter) {
	if l.Kind == ForInLoop {
		emitter.EmitIterKey(dst, l.Iterator)
	} else {
		emitter.EmitAB(vm.OpMove, dst, l.Iterator)
	}
}

func (l *Loop) EmitFinalization(emitter *Emitter) {
	emitter.EmitJump(l.JumpLabel)
	emitter.MarkLabel(l.EndLabel)

	if l.Kind == ForInLoop {
		emitter.EmitA(vm.OpClose, l.Iterator)
	}
}

func (l *Loop) PatchDestinationAx(alloc *RegisterAllocator, emitter *Emitter, op vm.Opcode, arg int) vm.Operand {
	if l.Allocate {
		emitter.SwapAx(l.StartLabel, op, l.Dst, arg)

		return l.Dst
	}

	tmp := alloc.Allocate(Temp)
	emitter.InsertAx(l.StartLabel, op, tmp, arg)
	return tmp
}

func (l *Loop) PatchDestinationAxy(alloc *RegisterAllocator, emitter *Emitter, op vm.Opcode, arg1, arg2 int) vm.Operand {
	if l.Allocate {
		emitter.SwapAxy(l.StartLabel, op, l.Dst, arg1, arg2)

		return l.Dst
	}

	tmp := alloc.Allocate(Temp)
	emitter.InsertAxy(l.StartLabel, op, tmp, arg1, arg2)
	return tmp
}

func (l *Loop) emitForInLoopIteration(alloc *RegisterAllocator, emitter *Emitter) {
	if l.Iterator == vm.NoopOperand {
		l.Iterator = alloc.Allocate(Temp)
	}

	emitter.EmitIter(l.Iterator, l.Src)
	emitter.MarkLabel(l.JumpLabel)
	emitter.EmitJumpc(vm.OpIterNext, l.Iterator, l.EndLabel)
}

func (l *Loop) emitForWhileLoopIteration(alloc *RegisterAllocator, emitter *Emitter) {
	l.Iterator = alloc.Allocate(Temp)
	emitter.EmitA(vm.OpLoadZero, l.Iterator)
	emitter.EmitA(vm.OpDecr, l.Iterator)

	// Placeholder for the loop condition
	emitter.MarkLabel(l.JumpLabel)

	if l.SrcFn == nil {
		panic("source function must be defined for while loop")
	}

	emitter.EmitA(vm.OpIncr, l.Iterator)

	l.Src = l.SrcFn()

	emitter.EmitJumpIfFalse(l.Src, l.EndLabel)
}

func (l *Loop) canDeclareVar(name string) bool {
	return name != "" && name != IgnorePseudoVariable
}

func (l *Loop) canBindVar(op vm.Operand) bool {
	return op != vm.NoopOperand
}
