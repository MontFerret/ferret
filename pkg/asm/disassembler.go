package asm

import (
	"bytes"
	"fmt"
	"text/tabwriter"

	"github.com/MontFerret/ferret/pkg/vm"
)

// Disassemble returns a human-readable disassembly of the given program.
func Disassemble(p *vm.Program) string {
	labels := collectLabels(p.Bytecode)

	var buf bytes.Buffer
	w := tabwriter.NewWriter(&buf, 0, 4, 2, ' ', 0)

	// Header: params
	for _, line := range formatParams(p) {
		fmt.Fprintln(w, line)
	}

	// Body: disassembly
	for ip, instr := range p.Bytecode {
		if label, ok := labels[ip]; ok {
			fmt.Fprintf(w, "%s:\n", label)
		}

		fmt.Fprintln(w, disasmLine(ip, instr, p, labels))
	}

	w.Flush()

	return buf.String()
}

// collectLabels identifies jump targets and assigns symbolic labels to them.
func collectLabels(bytecode []vm.Instruction) map[int]string {
	labels := make(map[int]string)
	counter := 0

	for _, instr := range bytecode {
		switch instr.Opcode {
		case vm.OpJump, vm.OpJumpIfFalse, vm.OpJumpIfTrue, vm.OpIterNext, vm.OpIterSkip, vm.OpIterLimit:
			target := int(instr.Operands[0])
			if _, ok := labels[target]; !ok {
				labels[target] = fmt.Sprintf("@L%d", counter)
				counter++
			}
		default:
			// Do nothing for other opcodes
		}
	}

	return labels
}

// disasmLine renders a single instruction into text, with optional constants and location info.
func disasmLine(ip int, instr vm.Instruction, p *vm.Program, labels map[int]string) string {
	ops := instr.Operands
	var out string

	opcode := instr.Opcode

	switch opcode {
	case vm.OpLoadConst:
		cIdx := ops[1].Constant()
		comment := constValue(p, cIdx)
		out = fmt.Sprintf("%d: %s R%d C%d ; %s", ip, opcode, ops[0], cIdx, comment)

	case vm.OpMove:
		out = fmt.Sprintf("%d: %s R%d R%d", ip, opcode, ops[0], ops[1])

	case vm.OpAdd:
		out = fmt.Sprintf("%d: %s R%d R%d R%d", ip, opcode, ops[0], ops[1], ops[2])

	case vm.OpJump:
		out = fmt.Sprintf("%d: %s %s", ip, opcode, labelOrAddr(int(ops[0]), labels))

	case vm.OpJumpIfTrue, vm.OpJumpIfFalse, vm.OpIterNext:
		out = fmt.Sprintf("%d: %s %s %s", ip, opcode, labelOrAddr(int(ops[0]), labels), ops[1])

	case vm.OpIterSkip, vm.OpIterLimit:
		out = fmt.Sprintf("%d: %s %s %s %s", ip, opcode, labelOrAddr(int(ops[0]), labels), ops[1], ops[2])

	case vm.OpReturn:
		out = fmt.Sprintf("%d: %s R%d", ip, opcode, ops[0])

	default:
		out = fmt.Sprintf("%d: %s %v", ip, opcode, ops)
	}

	if loc := formatLocation(p, ip); loc != "" {
		out += " " + loc
	}

	return out
}
