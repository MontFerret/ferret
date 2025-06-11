package internal

import (
	"github.com/MontFerret/ferret/pkg/vm"
)

type Emitter struct {
	instructions []vm.Instruction
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

// PatchSwapAB modifies an instruction at the given position to swap operands and update its operation and destination.
func (e *Emitter) PatchSwapAB(pos int, op vm.Opcode, dst, src1 vm.Operand) {
	e.instructions[pos] = vm.Instruction{
		Opcode:   op,
		Operands: [3]vm.Operand{dst, src1, vm.NoopOperand},
	}
}

// PatchSwapAx modifies an existing instruction at the specified position with a new opcode, destination, and argument.
func (e *Emitter) PatchSwapAx(pos int, op vm.Opcode, dst vm.Operand, arg int) {
	e.instructions[pos] = vm.Instruction{
		Opcode:   op,
		Operands: [3]vm.Operand{dst, vm.Operand(arg), vm.NoopOperand},
	}
}

// PatchSwapAxy replaces an instruction at the specified position with a new one using the provided opcode and operands.
func (e *Emitter) PatchSwapAxy(pos int, op vm.Opcode, dst vm.Operand, arg1, agr2 int) {
	e.instructions[pos] = vm.Instruction{
		Opcode:   op,
		Operands: [3]vm.Operand{dst, vm.Operand(arg1), vm.Operand(agr2)},
	}
}

// PatchSwapAs replaces an instruction at the specified position with a new instruction using the given opcode and operands.
func (e *Emitter) PatchSwapAs(pos int, op vm.Opcode, dst vm.Operand, seq RegisterSequence) {
	e.instructions[pos] = vm.Instruction{
		Opcode:   op,
		Operands: [3]vm.Operand{dst, seq[0], seq[len(seq)-1]},
	}
}

// PatchJump patches a jump opcode.
func (e *Emitter) PatchJump(instr int) {
	e.instructions[instr].Operands[0] = vm.Operand(len(e.instructions) - 1)
}

// PatchJumpAB patches a jump opcode with a new destination.
func (e *Emitter) PatchJumpAB(inst int) {
	e.instructions[inst].Operands[2] = vm.Operand(len(e.instructions) - 1)
}

// PatchJumpNextAB patches a jump instruction to jump over a current position.
func (e *Emitter) PatchJumpNextAB(instr int) {
	e.instructions[instr].Operands[2] = vm.Operand(len(e.instructions))
}

// PatchJumpNext patches a jump instruction to jump over a current position.
func (e *Emitter) PatchJumpNext(instr int) {
	e.instructions[instr].Operands[0] = vm.Operand(len(e.instructions))
}

// Emit emits an opcode with no arguments.
func (e *Emitter) Emit(op vm.Opcode) {
	e.EmitABC(op, 0, 0, 0)
}

// EmitA emits an opcode with a single destination register argument.
func (e *Emitter) EmitA(op vm.Opcode, dest vm.Operand) {
	e.EmitABC(op, dest, 0, 0)
}

// EmitAB emits an opcode with a destination register and a single source register argument.
func (e *Emitter) EmitAB(op vm.Opcode, dest, src1 vm.Operand) {
	e.EmitABC(op, dest, src1, 0)
}

// EmitAb emits an opcode with a destination register and a boolean argument.
func (e *Emitter) EmitAb(op vm.Opcode, dest vm.Operand, arg bool) {
	var src1 vm.Operand

	if arg {
		src1 = 1
	}

	e.EmitABC(op, dest, src1, 0)
}

// EmitAx emits an opcode with a destination register and a custom argument.
func (e *Emitter) EmitAx(op vm.Opcode, dest vm.Operand, arg int) {
	e.EmitABC(op, dest, vm.Operand(arg), 0)
}

// EmitAs emits an opcode with a destination register and a sequence of registers.
func (e *Emitter) EmitAs(op vm.Opcode, dest vm.Operand, seq RegisterSequence) {
	if seq != nil {
		src1 := seq[0]
		src2 := seq[len(seq)-1]
		e.EmitABC(op, dest, src1, src2)
	} else {
		e.EmitA(op, dest)
	}
}

// EmitABx emits an opcode with a destination and source register and a custom argument.
func (e *Emitter) EmitABx(op vm.Opcode, dest vm.Operand, src vm.Operand, arg int) {
	e.EmitABC(op, dest, src, vm.Operand(arg))
}

// EmitABC emits an opcode with a destination register and two source register arguments.
func (e *Emitter) EmitABC(op vm.Opcode, dest, src1, src2 vm.Operand) {
	e.instructions = append(e.instructions, vm.Instruction{
		Opcode:   op,
		Operands: [3]vm.Operand{dest, src1, src2},
	})
}
