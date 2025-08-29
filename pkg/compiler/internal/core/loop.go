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
	ForStepLoop
)

type Loop struct {
	Kind     LoopKind
	Type     LoopType
	Distinct bool
	Allocate bool

	StartLabel Label
	CondLabel  Label
	IncrLabel  Label
	EndLabel   Label

	Src   vm.Operand
	State vm.Operand

	ValueName string
	Value     vm.Operand
	KeyName   string
	Key       vm.Operand

	// For WHILE/STEP loops
	InitFn      func() vm.Operand
	ConditionFn func() vm.Operand
	IncrementFn func() vm.Operand

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
	l.CondLabel = emitter.NewLabel("loop", name, "cond")
	l.EndLabel = emitter.NewLabel("loop", name, "end")

	if l.Kind != ForInLoop {
		l.IncrLabel = emitter.NewLabel("loop", name, "incr")
	}

	emitter.MarkLabel(l.StartLabel)

	if l.Allocate {
		emitter.EmitAb(vm.OpDataSet, l.Dst, l.Distinct)
	}

	switch l.Kind {
	case ForInLoop:
		l.emitForInLoopIteration(alloc, emitter)
	case ForStepLoop:
		l.emitForStepLoopIteration(alloc, emitter)
	default:
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
	// For WHILE/STEP loops, the value is already in the destination register
	// No additional emission needed as the variable is directly assigned
	if l.Kind == ForInLoop {
		emitter.EmitIterValue(dst, l.State)
	}
}

func (l *Loop) EmitKey(dst vm.Operand, emitter *Emitter) {
	if l.Kind == ForInLoop {
		emitter.EmitIterKey(dst, l.State)
	}
}

func (l *Loop) EmitFinalization(emitter *Emitter) {
	if l.Kind == ForInLoop {
		emitter.EmitJump(l.CondLabel)
	} else {
		emitter.EmitJump(l.IncrLabel)
	}

	emitter.MarkLabel(l.EndLabel)

	if l.Kind == ForInLoop {
		emitter.EmitA(vm.OpClose, l.State)
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
	if l.State == vm.NoopOperand {
		l.State = alloc.Allocate(Temp)
	}

	emitter.EmitIter(l.State, l.Src)
	emitter.MarkLabel(l.CondLabel)
	emitter.EmitJumpc(vm.OpIterNext, l.State, l.EndLabel)
}

func (l *Loop) emitForWhileLoopIteration(_ *RegisterAllocator, emitter *Emitter) {
	if l.ConditionFn == nil {
		panic("condition function must be defined for while loop")
	}

	if l.Value != vm.NoopOperand {
		// Initialize the loop variable
		emitter.EmitA(vm.OpLoadZero, l.Value)
	}

	// Jump to the initial condition check (skipping the increment)
	emitter.EmitJump(l.CondLabel)

	// Placeholder for the loop condition
	emitter.MarkLabel(l.IncrLabel)

	if l.Value != vm.NoopOperand {
		emitter.EmitA(vm.OpIncr, l.Value)
	}

	// Mark the continue label (initial condition check point)
	emitter.MarkLabel(l.CondLabel)

	// Evaluate the condition
	condition := l.ConditionFn()
	emitter.EmitJumpIfFalse(condition, l.EndLabel)
}

func (l *Loop) emitForStepLoopIteration(_ *RegisterAllocator, emitter *Emitter) {
	if l.InitFn == nil || l.ConditionFn == nil || l.IncrementFn == nil {
		panic("step functions must be defined for step loop")
	}

	// Initialize the loop variable
	initValue := l.InitFn()

	if l.Value != vm.NoopOperand {
		emitter.EmitAB(vm.OpMove, l.Value, initValue)
	}

	// Jump to the initial condition check (skipping the increment)
	emitter.EmitJump(l.IncrLabel)

	// Mark the jump target for loop iterations (increment + condition check)
	emitter.MarkLabel(l.CondLabel)

	// Execute increment (this happens on every loop-back, but not on first iteration)
	if l.Value != vm.NoopOperand {
		incrementValue := l.IncrementFn()
		emitter.EmitAB(vm.OpMove, l.Value, incrementValue)
	}

	// Mark the continue label (initial condition check point)
	emitter.MarkLabel(l.IncrLabel)

	// Evaluate the condition
	condition := l.ConditionFn()
	emitter.EmitJumpIfFalse(condition, l.EndLabel)
}

func (l *Loop) emitStepIncrement(emitter *Emitter) {
	if l.Kind == ForStepLoop && l.IncrementFn != nil {
		// Execute the increment expression and assign it to the loop variable
		incrementValue := l.IncrementFn()

		if l.Value != vm.NoopOperand {
			emitter.EmitAB(vm.OpMove, l.Value, incrementValue)
		}
	}
}

func (l *Loop) canDeclareVar(name string) bool {
	return name != "" && name != IgnorePseudoVariable
}

func (l *Loop) canBindVar(op vm.Operand) bool {
	return op != vm.NoopOperand
}
