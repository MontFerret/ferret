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

// EmitJump emits a jump opcode.
func (e *Emitter) EmitJump(op runtime.Opcode, pos int) int {
	e.EmitA(op, runtime.Operand(pos))

	return len(e.instructions) - 1
}

// EmitJumpAB emits a jump opcode with a state and an argument.
func (e *Emitter) EmitJumpAB(op runtime.Opcode, state, cond runtime.Operand, pos int) int {
	e.EmitABC(op, state, cond, runtime.Operand(pos))

	return len(e.instructions) - 1
}

// EmitJumpc emits a conditional jump opcode.
func (e *Emitter) EmitJumpc(op runtime.Opcode, pos int, reg runtime.Operand) int {
	e.EmitAB(op, runtime.Operand(pos), reg)

	return len(e.instructions) - 1
}

// PatchJump patches a jump opcode.
func (e *Emitter) PatchJump(instr int) {
	e.instructions[instr].Operands[0] = runtime.Operand(len(e.instructions) - 1)
}

// PatchJumpAB patches a jump opcode with a new destination.
func (e *Emitter) PatchJumpAB(inst int) {
	e.instructions[inst].Operands[2] = runtime.Operand(len(e.instructions) - 1)
}

// PatchJumpNextAB patches a jump instruction to jump over a current position.
func (e *Emitter) PatchJumpNextAB(instr int) {
	e.instructions[instr].Operands[2] = runtime.Operand(len(e.instructions))
}

// PatchJumpNext patches a jump instruction to jump over a current position.
func (e *Emitter) PatchJumpNext(instr int) {
	e.instructions[instr].Operands[0] = runtime.Operand(len(e.instructions))
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

// EmitAb emits an opcode with a destination register and a boolean argument.
func (e *Emitter) EmitAb(op runtime.Opcode, dest runtime.Operand, arg bool) {
	var src1 runtime.Operand

	if arg {
		src1 = 1
	}

	e.EmitABC(op, dest, src1, 0)
}

// EmitAx emits an opcode with a destination register and a custom argument.
func (e *Emitter) EmitAx(op runtime.Opcode, dest runtime.Operand, arg int) {
	e.EmitABC(op, dest, runtime.Operand(arg), 0)
}

// EmitAs emits an opcode with a destination register and a sequence of registers.
func (e *Emitter) EmitAs(op runtime.Opcode, dest runtime.Operand, seq *RegisterSequence) {
	if seq != nil {
		src1 := seq.Registers[0]
		src2 := seq.Registers[len(seq.Registers)-1]
		e.EmitABC(op, dest, src1, src2)
	} else {
		e.EmitA(op, dest)
	}
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
