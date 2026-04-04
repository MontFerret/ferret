package core

import (
	"github.com/MontFerret/ferret/v2/pkg/bytecode"
)

type (
	LoopType int

	LoopKind int

	Loop struct {
		ConditionFn   func() bytecode.Operand
		ValueName     string
		LabelBase     string
		KeyName       string
		resetRegs     []bytecode.Operand
		continueLabel Label
		endLabel      Label
		bodyLabel     Label
		condLabel     Label
		startLabel    Label
		Type          LoopType
		Dst           bytecode.Operand
		Src           bytecode.Operand
		Kind          LoopKind
		Value         bytecode.Operand
		Key           bytecode.Operand
		State         bytecode.Operand
		Distinct      bool
		Allocate      bool
	}
)

const (
	NormalLoop LoopType = iota
	PassThroughLoop
	TemporalLoop
)

const (
	ForInLoop LoopKind = iota
	WhileLoop
	DoWhileLoop
)

func (l *Loop) StartLabel() Label {
	return l.startLabel
}

func (l *Loop) ContinueLabel() Label {
	if l.Kind == ForInLoop {
		return l.condLabel
	}

	return l.continueLabel
}

func (l *Loop) BreakLabel() Label {
	return l.endLabel
}

func (l *Loop) DeclareKeyVar(name string, st *SymbolTable, typ ValueType) bool {
	if l.canDeclareVar(name) {
		reg, ok := st.DeclareLocal(name, typ)

		if !ok {
			return false
		}

		l.Key = reg
		l.KeyName = name
	}

	return true
}

func (l *Loop) DeclareValueVar(name string, st *SymbolTable, typ ValueType) bool {
	if l.canDeclareVar(name) {
		reg, ok := st.DeclareLocal(name, typ)

		if !ok {
			return false
		}

		l.Value = reg
		l.ValueName = name
	}

	return true
}

func (l *Loop) EmitInitialization(alloc *RegisterAllocator, emitter *Emitter) {
	name := l.LabelBase
	l.startLabel = emitter.NewLabel("loop", name, "start")
	l.condLabel = emitter.NewLabel("loop", name, "cond")
	l.continueLabel = emitter.NewLabel("loop", l.LabelBase, "continue")
	l.bodyLabel = emitter.NewLabel("loop", name, "body")
	l.endLabel = emitter.NewLabel("loop", name, "end")

	emitter.MarkLabel(l.startLabel)

	if l.Allocate {
		emitter.EmitAb(bytecode.OpDataSet, l.Dst, l.Distinct)
	}

	switch l.Kind {
	case ForInLoop:
		l.emitForInLoopIteration(alloc, emitter)
	case WhileLoop:
		l.emitForWhileLoopIteration(alloc, emitter)
	default:
		l.emitForDoWhileLoopIteration(alloc, emitter)
	}

	if l.canBindVar(l.Value) {
		l.EmitValue(l.Value, emitter)
	}

	if l.canBindVar(l.Key) {
		l.EmitKey(l.Key, emitter)
	}

	emitter.MarkLabel(l.bodyLabel)
}

func (l *Loop) RegisterReset(reg bytecode.Operand) {
	if reg == bytecode.NoopOperand {
		return
	}

	l.resetRegs = append(l.resetRegs, reg)
}

func (l *Loop) EmitValue(dst bytecode.Operand, emitter *Emitter) {
	// For WHILE loops, the value is already in the destination register
	// No additional emission needed as the variable is directly assigned
	if l.Kind == ForInLoop {
		emitter.EmitIterValue(dst, l.State)
	}
}

func (l *Loop) EmitKey(dst bytecode.Operand, emitter *Emitter) {
	if l.Kind == ForInLoop {
		emitter.EmitIterKey(dst, l.State)
	}
}

func (l *Loop) EmitFinalization(emitter *Emitter) {
	emitter.EmitJump(l.ContinueLabel())
	emitter.MarkLabel(l.endLabel)

	if l.Kind == ForInLoop {
		emitter.EmitA(bytecode.OpClose, l.State)
	}

	for _, reg := range l.resetRegs {
		emitter.EmitA(bytecode.OpLoadZero, reg)
	}
}

func (l *Loop) PatchDestinationAx(alloc *RegisterAllocator, emitter *Emitter, op bytecode.Opcode, arg int) bytecode.Operand {
	if l.Allocate {
		emitter.SwapAx(l.startLabel, op, l.Dst, arg)

		return l.Dst
	}

	tmp := alloc.Allocate()
	emitter.InsertAx(l.startLabel, op, tmp, arg)
	return tmp
}

func (l *Loop) PatchDestinationAxy(alloc *RegisterAllocator, emitter *Emitter, op bytecode.Opcode, arg1, arg2 int) bytecode.Operand {
	if l.Allocate {
		emitter.SwapAxy(l.startLabel, op, l.Dst, arg1, arg2)

		return l.Dst
	}

	tmp := alloc.Allocate()
	emitter.InsertAxy(l.startLabel, op, tmp, arg1, arg2)
	return tmp
}

func (l *Loop) emitForInLoopIteration(alloc *RegisterAllocator, emitter *Emitter) {
	if l.State == bytecode.NoopOperand {
		l.State = alloc.Allocate()
	}

	emitter.EmitIter(l.State, l.Src)
	emitter.MarkLabel(l.condLabel)
	emitter.EmitJumpc(bytecode.OpIterNext, l.State, l.endLabel)
}

func (l *Loop) emitForWhileLoopIteration(_ *RegisterAllocator, emitter *Emitter) {
	if l.ConditionFn == nil {
		PanicInvariant("condition function must be defined for while loop")
	}

	if l.Value != bytecode.NoopOperand {
		// Initialize the loop variable
		emitter.EmitA(bytecode.OpLoadZero, l.Value)
	}

	// Jump to the initial condition check (skipping the increment)
	emitter.EmitJump(l.condLabel)

	// Mark the loop-back target even when the loop variable is ignored.
	emitter.MarkLabel(l.continueLabel)

	if l.Value != bytecode.NoopOperand {
		emitter.EmitA(bytecode.OpIncr, l.Value)
	}

	// Mark the continue label (initial condition check point)
	emitter.MarkLabel(l.condLabel)

	// Evaluate the condition
	condition := l.ConditionFn()
	emitter.EmitJumpIfFalse(condition, l.endLabel)
}

func (l *Loop) emitForDoWhileLoopIteration(_ *RegisterAllocator, emitter *Emitter) {
	if l.ConditionFn == nil {
		PanicInvariant("condition function must be defined for while loop")
	}

	if l.Value != bytecode.NoopOperand {
		// Initialize the loop variable
		emitter.EmitA(bytecode.OpLoadZero, l.Value)
	}

	// Jump to the loop body first
	emitter.EmitJump(l.bodyLabel)

	// Mark the loop-back target even when the loop variable is ignored.
	emitter.MarkLabel(l.continueLabel)

	if l.Value != bytecode.NoopOperand {
		emitter.EmitA(bytecode.OpIncr, l.Value)
	}

	// Mark the continue label (initial condition check point)
	emitter.MarkLabel(l.condLabel)

	// Evaluate the condition
	condition := l.ConditionFn()
	emitter.EmitJumpIfFalse(condition, l.endLabel)
}

func (l *Loop) canDeclareVar(name string) bool {
	return name != "" && name != IgnorePseudoVariable
}

func (l *Loop) canBindVar(op bytecode.Operand) bool {
	return op != bytecode.NoopOperand
}
