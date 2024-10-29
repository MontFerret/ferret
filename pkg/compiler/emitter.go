package compiler

import "github.com/MontFerret/ferret/pkg/runtime"

type Emitter struct {
	instructions []runtime.Instruction
}

func NewEmitter() *Emitter {
	return &Emitter{
		instructions: make([]runtime.Instruction, 0, 8),
	}
}

func (e *Emitter) Size() int {
	return len(e.instructions)
}

func (e *Emitter) EmitJump(op runtime.Opcode, reg runtime.Operand) int {
	e.EmitA(op, reg)

	return len(e.instructions) - 1
}

func (e *Emitter) PatchJump(dest int) {
	e.instructions[dest].Operands[1] = runtime.Operand(len(e.instructions) - dest - 1)
}

// Emit emits an opcode with no arguments.
func (e *Emitter) Emit(op runtime.Opcode) {
	e.EmitABC(op, 0, 0, 0)
}

// EmitA emits an opcode with a single destination register argument.
func (e *Emitter) EmitA(op runtime.Opcode, dest runtime.Operand) {
	e.EmitABC(op, dest, 0, 0)
}

// EmitAB emits an opcode with a destination register and a single source register argument.
func (e *Emitter) EmitAB(op runtime.Opcode, dest, src1 runtime.Operand) {
	e.EmitABC(op, dest, src1, 0)
}

// EmitAx emits an opcode with a destination register and a custom argument.
func (e *Emitter) EmitAx(op runtime.Opcode, dest runtime.Operand, arg int) {
	e.EmitABC(op, dest, runtime.Operand(arg), 0)
}

// EmitAs emits an opcode with a destination register and a sequence of registers.
func (e *Emitter) EmitAs(op runtime.Opcode, dest runtime.Operand, seq *RegisterSequence) {
	src1 := seq.Registers[0]
	src2 := seq.Registers[len(seq.Registers)-1]
	e.EmitABC(op, dest, src1, src2)
}

// EmitABx emits an opcode with a destination and source register and a custom argument.
func (e *Emitter) EmitABx(op runtime.Opcode, dest runtime.Operand, src runtime.Operand, arg int) {
	e.EmitABC(op, dest, src, runtime.Operand(arg))
}

// EmitABC emits an opcode with a destination register and two source register arguments.
func (e *Emitter) EmitABC(op runtime.Opcode, dest, src1, src2 runtime.Operand) {
	e.instructions = append(e.instructions, runtime.Instruction{
		Opcode:   op,
		Operands: [3]runtime.Operand{dest, src1, src2},
	})
}
