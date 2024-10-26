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

// EmitABx emits an opcode with a destination register and a custom argument.
func (e *Emitter) EmitABx(op runtime.Opcode, dest runtime.Operand, arg int) {
	e.EmitABC(op, dest, runtime.Operand(arg), 0)
}

// EmitABC emits an opcode with a destination register and two source register arguments.
func (e *Emitter) EmitABC(op runtime.Opcode, dest, src1, src2 runtime.Operand) {
	e.instructions = append(e.instructions, runtime.Instruction{
		Opcode:   op,
		Operands: [3]runtime.Operand{dest, src1, src2},
	})
}
