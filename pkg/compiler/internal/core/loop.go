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
	JumpLabel  Label
	EndLabel   Label
	ContinueLabel Label  // For STEP loops, where clauses jump to continue

	Src      vm.Operand
	SrcFn    func() vm.Operand
	Iterator vm.Operand

	ValueName string
	Value     vm.Operand
	KeyName   string
	Key       vm.Operand

	// For STEP loops
	StepInitFn      func() vm.Operand
	StepConditionFn func() vm.Operand
	StepIncrementFn func() vm.Operand

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
	
	// For STEP loops, we need a continue label where failed clauses can jump
	if l.Kind == ForStepLoop {
		l.ContinueLabel = emitter.NewLabel("loop", name, "continue")
	}

	emitter.MarkLabel(l.StartLabel)

	if l.Allocate {
		emitter.EmitAb(vm.OpDataSet, l.Dst, l.Distinct)
	}

	if l.Kind == ForInLoop {
		l.emitForInLoopIteration(alloc, emitter)
	} else if l.Kind == ForStepLoop {
		l.emitForStepLoopIteration(alloc, emitter)
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
	} else if l.Kind == ForStepLoop {
		// For STEP loops, the value is already in the destination register
		// No additional emission needed as the variable is directly assigned
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

// Common helper methods for loop patterns

// allocateLoopIterator allocates an iterator register for WHILE and STEP loops
func (l *Loop) allocateLoopIterator(alloc *RegisterAllocator) {
	if l.Iterator == vm.NoopOperand {
		l.Iterator = alloc.Allocate(Temp)
	}
}

// markJumpLabel marks the common jump label used by all loop types
func (l *Loop) markJumpLabel(emitter *Emitter) {
	emitter.MarkLabel(l.JumpLabel)
}

// emitConditionalJumpToEnd emits a conditional jump to the loop end label
func (l *Loop) emitConditionalJumpToEnd(emitter *Emitter, condition vm.Operand) {
	emitter.EmitJumpIfFalse(condition, l.EndLabel)
}

// initializeCounterState initializes counter state for WHILE loops
func (l *Loop) initializeCounterState(emitter *Emitter) {
	emitter.EmitA(vm.OpLoadZero, l.Iterator)
	emitter.EmitA(vm.OpDecr, l.Iterator)
}

// incrementCounter increments the counter for WHILE loops
func (l *Loop) incrementCounter(emitter *Emitter) {
	emitter.EmitA(vm.OpIncr, l.Iterator)
}

func (l *Loop) emitForInLoopIteration(alloc *RegisterAllocator, emitter *Emitter) {
	l.allocateLoopIterator(alloc)

	emitter.EmitIter(l.Iterator, l.Src)
	l.markJumpLabel(emitter)
	emitter.EmitJumpc(vm.OpIterNext, l.Iterator, l.EndLabel)
}

func (l *Loop) emitForWhileLoopIteration(alloc *RegisterAllocator, emitter *Emitter) {
	l.Iterator = alloc.Allocate(Temp)
	l.initializeCounterState(emitter)

	// Placeholder for the loop condition
	l.markJumpLabel(emitter)

	if l.SrcFn == nil {
		panic("source function must be defined for while loop")
	}

	l.incrementCounter(emitter)

	l.Src = l.SrcFn()

	l.emitConditionalJumpToEnd(emitter, l.Src)
}

func (l *Loop) emitForStepLoopIteration(alloc *RegisterAllocator, emitter *Emitter) {
	if l.StepInitFn == nil || l.StepConditionFn == nil || l.StepIncrementFn == nil {
		panic("step functions must be defined for step loop")
	}

	// Initialize the loop variable
	initValue := l.StepInitFn()
	if l.Value != vm.NoopOperand {
		emitter.EmitAB(vm.OpMove, l.Value, initValue)
	}

	// Jump to the initial condition check (skipping the increment)
	emitter.EmitJump(l.ContinueLabel)

	// Mark the jump target for loop iterations (increment + condition check)
	l.markJumpLabel(emitter)

	// Execute increment (this happens on every loop-back, but not on first iteration)
	if l.Value != vm.NoopOperand {
		incrementValue := l.StepIncrementFn()
		emitter.EmitAB(vm.OpMove, l.Value, incrementValue)
	}

	// Mark the continue label (initial condition check point)
	emitter.MarkLabel(l.ContinueLabel)

	// Evaluate the condition
	condition := l.StepConditionFn()
	l.emitConditionalJumpToEnd(emitter, condition)
}

func (l *Loop) EmitStepIncrement(emitter *Emitter) {
	if l.Kind == ForStepLoop && l.StepIncrementFn != nil {
		// Execute the increment expression and assign it to the loop variable
		incrementValue := l.StepIncrementFn()
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
