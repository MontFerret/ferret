package core

import (
	"fmt"
	"strings"

	"github.com/MontFerret/ferret/pkg/vm"
)

type Emitter struct {
	instructions []vm.Instruction
	labels       map[labelID]Label
	patches      map[labelID][]labelRef
	nextLabelID  labelID
}

func NewEmitter() *Emitter {
	return &Emitter{
		instructions: make([]vm.Instruction, 0, 8),
	}
}

func (e *Emitter) Bytecode() []vm.Instruction {
	return e.instructions
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
	//
	//// Back-patch any prior references to this label (keep them for future retargeting)
	//if refs, ok := e.patches[label.id]; ok {
	//	for _, ref := range refs {
	//		e.patchOperand(ref.pos, ref.field, len(e.instructions))
	//	}
	//}
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
func (e *Emitter) Emit(op vm.Opcode) {
	e.EmitABC(op, 0, 0, 0)
}

// EmitA emits an opcode with a single destination value argument.
func (e *Emitter) EmitA(op vm.Opcode, dest vm.Operand) {
	e.EmitABC(op, dest, 0, 0)
}

// EmitAB emits an opcode with a destination value and a single source value argument.
func (e *Emitter) EmitAB(op vm.Opcode, dest, src1 vm.Operand) {
	e.EmitABC(op, dest, src1, 0)
}

// EmitAb emits an opcode with a destination value and a boolean argument.
func (e *Emitter) EmitAb(op vm.Opcode, dest vm.Operand, arg bool) {
	var src1 vm.Operand

	if arg {
		src1 = 1
	}

	e.EmitABC(op, dest, src1, 0)
}

// EmitAx emits an opcode with a destination value and a custom argument.
func (e *Emitter) EmitAx(op vm.Opcode, dest vm.Operand, arg int) {
	e.EmitABC(op, dest, vm.Operand(arg), 0)
}

// EmitAxy emits an instruction with the given opcode, destination operand, and two integer arguments converted to operands.
func (e *Emitter) EmitAxy(op vm.Opcode, dest vm.Operand, arg1, agr2 int) {
	e.EmitABC(op, dest, vm.Operand(arg1), vm.Operand(agr2))
}

// EmitAs emits an opcode with a destination value and a sequence of registers.
func (e *Emitter) EmitAs(op vm.Opcode, dest vm.Operand, seq RegisterSequence) {
	if seq != nil {
		src1 := seq[0]
		src2 := seq[len(seq)-1]
		e.EmitABC(op, dest, src1, src2)
	} else {
		e.EmitA(op, dest)
	}
}

// EmitABx emits an opcode with a destination and source value and a custom argument.
func (e *Emitter) EmitABx(op vm.Opcode, dest vm.Operand, src vm.Operand, arg int) {
	e.EmitABC(op, dest, src, vm.Operand(arg))
}

// EmitABC emits an opcode with a destination value and two source value arguments.
func (e *Emitter) EmitABC(op vm.Opcode, dest, src1, src2 vm.Operand) {
	e.instructions = append(e.instructions, vm.Instruction{
		Opcode:   op,
		Operands: [3]vm.Operand{dest, src1, src2},
	})
}

// SwapAB modifies an instruction at the given position to swap operands and update its operation and destination.
func (e *Emitter) SwapAB(label Label, op vm.Opcode, dst, src1 vm.Operand) {
	e.swapInstruction(label, vm.Instruction{
		Opcode:   op,
		Operands: [3]vm.Operand{dst, src1, vm.NoopOperand},
	})
}

// SwapAx modifies an existing instruction at the specified position with a new opcode, destination, and argument.
func (e *Emitter) SwapAx(label Label, op vm.Opcode, dst vm.Operand, arg int) {
	e.swapInstruction(label, vm.Instruction{
		Opcode:   op,
		Operands: [3]vm.Operand{dst, vm.Operand(arg), vm.NoopOperand},
	})
}

// SwapAxy replaces an instruction at the specified position with a new one using the provided opcode and operands.
func (e *Emitter) SwapAxy(label Label, op vm.Opcode, dst vm.Operand, arg1, agr2 int) {
	e.swapInstruction(label, vm.Instruction{
		Opcode:   op,
		Operands: [3]vm.Operand{dst, vm.Operand(arg1), vm.Operand(agr2)},
	})
}

// SwapAs replaces an instruction at the specified position with a new instruction using the given opcode and operands.
func (e *Emitter) SwapAs(label Label, op vm.Opcode, dst vm.Operand, seq RegisterSequence) {
	e.swapInstruction(label, vm.Instruction{
		Opcode:   op,
		Operands: [3]vm.Operand{dst, seq[0], seq[len(seq)-1]},
	})
}

// InsertAx inserts a new instruction at a specific position in the instructions slice, shifting elements to the right.
// The inserted instruction includes an opcode and operands, where the third operand is set to a no-op by default.
func (e *Emitter) InsertAx(label Label, op vm.Opcode, dst vm.Operand, arg int) {
	e.insertInstruction(label, vm.Instruction{
		Opcode:   op,
		Operands: [3]vm.Operand{dst, vm.Operand(arg), vm.NoopOperand},
	})
}

// InsertAxy inserts an instruction at the specified position in the instruction list, shifting existing elements to the right.
func (e *Emitter) InsertAxy(label Label, op vm.Opcode, dst vm.Operand, arg1, arg2 int) {
	e.insertInstruction(label, vm.Instruction{
		Opcode:   op,
		Operands: [3]vm.Operand{dst, vm.Operand(arg1), vm.Operand(arg2)},
	})
}

func (e *Emitter) Patchx(label Label, arg int) {
	pos, ok := e.LabelPosition(label)

	if !ok {
		panic(fmt.Errorf("label not marked: %s", label))
	}

	current := e.instructions[pos]
	e.instructions[pos] = vm.Instruction{
		Opcode: current.Opcode,
		Operands: [3]vm.Operand{
			current.Operands[0],
			vm.Operand(arg),
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

	//// If the label is already marked, patch now as well
	//
	//if def, ok := e.labels[label.id]; ok {
	//	e.patchOperand(pos, field, def.addr)
	//}
}

// patchOperand modifies the operand at the specified position and field in the instruction set.
func (e *Emitter) patchOperand(pos int, field int, val int) {
	ins := e.instructions[pos]
	ins.Operands[field] = vm.Operand(val)
	e.instructions[pos] = ins
}

// swapInstruction swaps the operands of an instruction at a given position.
func (e *Emitter) swapInstruction(label Label, ins vm.Instruction) {
	pos, ok := e.LabelPosition(label)

	if !ok {
		panic(fmt.Errorf("label not marked: %s", label))
	}

	e.instructions[pos] = ins
}

// swapInstruction swaps the operands of an instruction at a given position.
func (e *Emitter) insertInstruction(label Label, ins vm.Instruction) {
	pos, ok := e.LabelPosition(label)

	if !ok {
		panic(fmt.Errorf("label not marked: %s", label))
	}

	// Insert instruction at position
	e.instructions = append(e.instructions[:pos],
		append([]vm.Instruction{ins}, e.instructions[pos:]...)...,
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

	// Re-patch any references that target labels whose address has moved
	for l, newAddr := range moved {
		if refs, ok := e.patches[l]; ok {
			for _, ref := range refs {
				e.patchOperand(ref.pos, ref.field, newAddr)
			}
		}
	}
}
