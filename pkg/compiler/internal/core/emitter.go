package core

import (
	"fmt"
	"strings"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/file"
)

type Emitter struct {
	labels           map[labelID]Label
	patches          map[labelID][]labelRef
	matchFailPatches map[labelID][]int
	instructions     []bytecode.Instruction
	selectorSlots    []int
	matchFailTargets []int
	spans            []file.Span
	currentSpan      file.Span
	nextLabelID      labelID
}

func NewEmitter() *Emitter {
	return &Emitter{
		instructions: make([]bytecode.Instruction, 0, 8),
		currentSpan:  file.Span{Start: -1, End: -1},
	}
}

func (e *Emitter) Bytecode() []bytecode.Instruction {
	return e.instructions
}

func (e *Emitter) Spans() []file.Span {
	if len(e.spans) == 0 {
		return nil
	}

	out := make([]file.Span, len(e.spans))
	copy(out, e.spans)

	return out
}

func (e *Emitter) AggregateSelectorSlots() []int {
	if len(e.selectorSlots) == 0 {
		return nil
	}

	out := make([]int, len(e.selectorSlots))
	copy(out, e.selectorSlots)

	return out
}

func (e *Emitter) MatchFailTargets() []int {
	if len(e.matchFailTargets) == 0 {
		return nil
	}

	out := make([]int, len(e.matchFailTargets))
	copy(out, e.matchFailTargets)

	return out
}

// WithSpan sets a span for emitted instructions within fn.
func (e *Emitter) WithSpan(span file.Span, fn func()) {
	if fn == nil {
		return
	}

	prev := e.currentSpan
	e.currentSpan = span
	fn()
	e.currentSpan = prev
}

// Size returns the number of instructions currently stored in the Emitter.
func (e *Emitter) Size() int {
	return len(e.instructions)
}

func (e *Emitter) Labels() map[int]string {
	if e.labels == nil {
		return nil
	}

	labels := make(map[int]string, len(e.labels))

	for _, def := range e.labels {
		labels[def.addr] = def.name
	}

	return labels
}

func (e *Emitter) NewLabel(name ...string) Label {
	l := e.nextLabelID
	e.nextLabelID++

	var labelName string

	if len(name) > 0 {
		if len(name) == 1 {
			labelName = name[0]
		} else {
			labelName = strings.Join(name, ".")
		}
	}

	return Label{
		id:   l,
		name: labelName,
	}
}

func (e *Emitter) MarkLabel(label Label) {
	if e.labels == nil {
		e.labels = make(map[labelID]Label)
	}

	e.labels[label.id] = Label{
		id:   label.id,
		name: label.name,
		addr: len(e.instructions),
	}

	// Back-patch any prior references to this label
	if refs, ok := e.patches[label.id]; ok {
		for _, ref := range refs {
			e.patchOperand(ref.pos, ref.field, len(e.instructions))
		}
	}

	if refs, ok := e.matchFailPatches[label.id]; ok {
		for _, pos := range refs {
			e.patchMatchFailTarget(pos, len(e.instructions))
		}
	}
}

func (e *Emitter) LabelPosition(label Label) (int, bool) {
	def, ok := e.labels[label.id]

	if !ok {
		return -1, false
	}

	return def.addr, true
}

func (e *Emitter) Position() int {
	return len(e.instructions) - 1
}

// Emit emits an opcode with no arguments.
func (e *Emitter) Emit(op bytecode.Opcode) {
	e.EmitABC(op, 0, 0, 0)
}

// EmitA emits an opcode with a single destination value argument.
func (e *Emitter) EmitA(op bytecode.Opcode, dest bytecode.Operand) {
	e.EmitABC(op, dest, 0, 0)
}

// EmitAB emits an opcode with a destination value and a single source value argument.
func (e *Emitter) EmitAB(op bytecode.Opcode, dest, src1 bytecode.Operand) {
	e.EmitABC(op, dest, src1, 0)
}

// EmitAb emits an opcode with a destination value and a boolean argument.
func (e *Emitter) EmitAb(op bytecode.Opcode, dest bytecode.Operand, arg bool) {
	var src1 bytecode.Operand

	if arg {
		src1 = 1
	}

	e.EmitABC(op, dest, src1, 0)
}

// EmitAx emits an opcode with a destination value and a custom argument.
func (e *Emitter) EmitAx(op bytecode.Opcode, dest bytecode.Operand, arg int) {
	e.EmitABC(op, dest, bytecode.Operand(arg), 0)
}

// EmitAxy emits an instruction with the given opcode, destination operand, and two integer arguments converted to operands.
func (e *Emitter) EmitAxy(op bytecode.Opcode, dest bytecode.Operand, arg1, agr2 int) {
	e.EmitABC(op, dest, bytecode.Operand(arg1), bytecode.Operand(agr2))
}

// EmitAs emits an opcode with a destination value and a sequence of registers.
func (e *Emitter) EmitAs(op bytecode.Opcode, dest bytecode.Operand, seq RegisterSequence) {
	if len(seq) > 0 {
		src1 := seq[0]
		src2 := seq[len(seq)-1]
		e.EmitABC(op, dest, src1, src2)
	} else {
		e.EmitA(op, dest)
	}
}

// EmitABx emits an opcode with a destination and source value and a custom argument.
func (e *Emitter) EmitABx(op bytecode.Opcode, dest bytecode.Operand, src bytecode.Operand, arg int) {
	e.EmitABC(op, dest, src, bytecode.Operand(arg))
}

// EmitABC emits an opcode with a destination value and two source value arguments.
func (e *Emitter) EmitABC(op bytecode.Opcode, dest, src1, src2 bytecode.Operand) {
	e.emitInstruction(bytecode.Instruction{
		Opcode:   op,
		Operands: [3]bytecode.Operand{dest, src1, src2},
	})
}

func (e *Emitter) emitInstruction(ins bytecode.Instruction) {
	e.emitInstructionWithMetadata(ins, -1, -1)
}

func (e *Emitter) emitInstructionWithSelectorSlot(ins bytecode.Instruction, selectorSlot int) {
	e.emitInstructionWithMetadata(ins, selectorSlot, -1)
}

func (e *Emitter) emitInstructionWithMatchFailTarget(ins bytecode.Instruction, matchFailTarget int) {
	e.emitInstructionWithMetadata(ins, -1, matchFailTarget)
}

func (e *Emitter) emitInstructionWithMetadata(ins bytecode.Instruction, selectorSlot, matchFailTarget int) {
	e.instructions = append(e.instructions, ins)
	e.selectorSlots = append(e.selectorSlots, selectorSlot)
	e.matchFailTargets = append(e.matchFailTargets, matchFailTarget)
	e.spans = append(e.spans, e.currentSpan)
}

// SwapAB modifies an instruction at the given position to swap operands and update its operation and destination.
func (e *Emitter) SwapAB(label Label, op bytecode.Opcode, dst, src1 bytecode.Operand) {
	e.swapInstruction(label, bytecode.Instruction{
		Opcode:   op,
		Operands: [3]bytecode.Operand{dst, src1, bytecode.NoopOperand},
	})
}

// SwapAx modifies an existing instruction at the specified position with a new opcode, destination, and argument.
func (e *Emitter) SwapAx(label Label, op bytecode.Opcode, dst bytecode.Operand, arg int) {
	e.swapInstruction(label, bytecode.Instruction{
		Opcode:   op,
		Operands: [3]bytecode.Operand{dst, bytecode.Operand(arg), bytecode.NoopOperand},
	})
}

// SwapAxy replaces an instruction at the specified position with a new one using the provided opcode and operands.
func (e *Emitter) SwapAxy(label Label, op bytecode.Opcode, dst bytecode.Operand, arg1, agr2 int) {
	e.swapInstruction(label, bytecode.Instruction{
		Opcode:   op,
		Operands: [3]bytecode.Operand{dst, bytecode.Operand(arg1), bytecode.Operand(agr2)},
	})
}

// SwapAs replaces an instruction at the specified position with a new instruction using the given opcode and operands.
func (e *Emitter) SwapAs(label Label, op bytecode.Opcode, dst bytecode.Operand, seq RegisterSequence) {
	e.swapInstruction(label, bytecode.Instruction{
		Opcode:   op,
		Operands: [3]bytecode.Operand{dst, seq[0], seq[len(seq)-1]},
	})
}

// InsertAx inserts a new instruction at a specific position in the instructions slice, shifting elements to the right.
// The inserted instruction includes an opcode and operands, where the third operand is set to a no-op by default.
func (e *Emitter) InsertAx(label Label, op bytecode.Opcode, dst bytecode.Operand, arg int) {
	e.insertInstruction(label, bytecode.Instruction{
		Opcode:   op,
		Operands: [3]bytecode.Operand{dst, bytecode.Operand(arg), bytecode.NoopOperand},
	})
}

// InsertAxy inserts an instruction at the specified position in the instruction list, shifting existing elements to the right.
func (e *Emitter) InsertAxy(label Label, op bytecode.Opcode, dst bytecode.Operand, arg1, arg2 int) {
	e.insertInstruction(label, bytecode.Instruction{
		Opcode:   op,
		Operands: [3]bytecode.Operand{dst, bytecode.Operand(arg1), bytecode.Operand(arg2)},
	})
}

func (e *Emitter) Patchx(label Label, arg int) {
	pos, ok := e.LabelPosition(label)

	if !ok {
		panic(fmt.Errorf("label not marked: %s", label))
	}

	current := e.instructions[pos]
	e.instructions[pos] = bytecode.Instruction{
		Opcode: current.Opcode,
		Operands: [3]bytecode.Operand{
			current.Operands[0],
			bytecode.Operand(arg),
			current.Operands[2],
		},
	}
}

// addLabelRef adds a reference to a label at a specific position and field in the instruction set.
func (e *Emitter) addLabelRef(pos int, field int, label Label) {
	if e.labels == nil {
		e.labels = make(map[labelID]Label)
	}

	if def, ok := e.labels[label.id]; ok {
		// Already marked → patch immediately
		e.patchOperand(pos, field, def.addr)
		return
	}

	if e.patches == nil {
		e.patches = make(map[labelID][]labelRef)
	}

	// Always remember the reference so we can retarget it later if needed
	e.patches[label.id] = append(e.patches[label.id], labelRef{pos: pos, field: field})
}

// patchOperand modifies the operand at the specified position and field in the instruction set.
func (e *Emitter) patchOperand(pos int, field int, val int) {
	ins := e.instructions[pos]
	ins.Operands[field] = bytecode.Operand(val)
	e.instructions[pos] = ins
}

// swapInstruction swaps the operands of an instruction at a given position.
func (e *Emitter) swapInstruction(label Label, ins bytecode.Instruction) {
	e.swapInstructionWithSelectorSlot(label, ins, -1)
}

func (e *Emitter) swapInstructionWithSelectorSlot(label Label, ins bytecode.Instruction, selectorSlot int) {
	pos, ok := e.LabelPosition(label)

	if !ok {
		panic(fmt.Errorf("label not marked: %s", label))
	}

	e.instructions[pos] = ins
	e.selectorSlots[pos] = selectorSlot
	e.matchFailTargets[pos] = -1
}

// swapInstruction swaps the operands of an instruction at a given position.
func (e *Emitter) insertInstruction(label Label, ins bytecode.Instruction) {
	e.insertInstructionWithSelectorSlot(label, ins, -1)
}

func (e *Emitter) insertInstructionWithSelectorSlot(label Label, ins bytecode.Instruction, selectorSlot int) {
	pos, ok := e.LabelPosition(label)

	if !ok {
		panic(fmt.Errorf("label not marked: %s", label))
	}

	// Insert instruction at position
	e.instructions = append(e.instructions[:pos],
		append([]bytecode.Instruction{ins}, e.instructions[pos:]...)...,
	)
	e.selectorSlots = append(e.selectorSlots[:pos],
		append([]int{selectorSlot}, e.selectorSlots[pos:]...)...,
	)
	e.matchFailTargets = append(e.matchFailTargets[:pos],
		append([]int{-1}, e.matchFailTargets[pos:]...)...,
	)
	e.spans = append(e.spans[:pos],
		append([]file.Span{e.currentSpan}, e.spans[pos:]...)...,
	)

	// Adjust all subsequent label addresses
	moved := make(map[labelID]int, 4)
	for l, d := range e.labels {
		if d.addr >= pos {
			newAddr := d.addr + 1
			e.labels[l] = Label{
				id:   d.id,
				name: d.name,
				addr: newAddr,
			}
			moved[l] = newAddr
		}
	}

	// Adjust all patch positions as well
	for l, refs := range e.patches {
		for i, ref := range refs {
			if ref.pos >= pos {
				e.patches[l][i].pos++
			}
		}
	}

	for l, refs := range e.matchFailPatches {
		for i, refPos := range refs {
			if refPos >= pos {
				e.matchFailPatches[l][i]++
			}
		}
	}

	// Re-patch any references that target labels whose address has moved
	for l, newAddr := range moved {
		if refs, ok := e.patches[l]; ok {
			for _, ref := range refs {
				e.patchOperand(ref.pos, ref.field, newAddr)
			}
		}

		if refs, ok := e.matchFailPatches[l]; ok {
			for _, refPos := range refs {
				e.patchMatchFailTarget(refPos, newAddr)
			}
		}
	}
}

func (e *Emitter) addMatchFailLabelRef(pos int, label Label) {
	if e.labels == nil {
		e.labels = make(map[labelID]Label)
	}

	if def, ok := e.labels[label.id]; ok {
		e.patchMatchFailTarget(pos, def.addr)
		return
	}

	if e.matchFailPatches == nil {
		e.matchFailPatches = make(map[labelID][]int)
	}

	e.matchFailPatches[label.id] = append(e.matchFailPatches[label.id], pos)
}

func (e *Emitter) patchMatchFailTarget(pos int, val int) {
	e.matchFailTargets[pos] = val
}
