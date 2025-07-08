package asm

import (
	"bytes"
	"fmt"
	"text/tabwriter"

	"github.com/MontFerret/ferret/pkg/vm"
)

// Disassemble returns a human-readable disassembly of the given program.
func Disassemble(p *vm.Program, options ...DisassemblerOption) (string, error) {
	if p == nil {
		return "", ErrInvalidProgram
	}

	newDisassemblerOptions(options...)

	labels := collectLabels(p.Bytecode, p.Labels)

	var buf bytes.Buffer
	w := tabwriter.NewWriter(&buf, 0, 4, 2, ' ', 0)

	// Header: params
	for _, name := range p.Params {
		_, _ = fmt.Fprintln(w, formatParam(name))
	}

	// Header: functions
	for name, args := range p.Functions {
		_, _ = fmt.Fprintln(w, formatFunction(name, args))
	}

	// Header: constants
	for _, constant := range p.Constants {
		_, _ = fmt.Fprintln(w, formatConstant(constant))
	}

	if buf.Len() > 0 {
		_, _ = fmt.Fprintln(w)
	}

	// Body: disassembly
	for ip, instr := range p.Bytecode {
		if label, ok := labels[ip]; ok {
			_, _ = fmt.Fprintf(w, "%s:\n", label)
		}

		_, _ = fmt.Fprintf(w, "\t%s\n", disasmLine(ip, instr, p, labels))
	}

	_ = w.Flush()

	return buf.String(), nil
}

// collectLabels identifies jump targets and assigns symbolic labels to them.
func collectLabels(bytecode []vm.Instruction, names map[int]string) map[int]string {
	labels := make(map[int]string)
	counter := 0

	// Iterate through the labels in the program to initialize the labels map
	for target, name := range names {
		labels[target] = fmt.Sprintf("@%s", name)
	}

	// Collect unmarked jump targets
	for _, instr := range bytecode {
		switch instr.Opcode {
		case vm.OpJump, vm.OpJumpIfFalse, vm.OpJumpIfTrue, vm.OpIterNext, vm.OpIterSkip, vm.OpIterLimit:
			target := int(instr.Operands[0])

			if name, ok := names[target]; !ok || name == "" {
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
	// Op Jmp
	case vm.OpJump:
		out = fmt.Sprintf("%d: %s %s", ip, opcode, labelOrAddr(int(ops[0]), labels))

	// Op Jmp R
	case vm.OpJumpIfTrue, vm.OpJumpIfFalse, vm.OpIterNext:
		out = fmt.Sprintf("%d: %s %s %s", ip, opcode, labelOrAddr(int(ops[0]), labels), formatOperand(ops[1]))

	// Op Jmp R, R
	case vm.OpIterSkip, vm.OpIterLimit:
		out = fmt.Sprintf("%d: %s %s %s %s", ip, opcode, labelOrAddr(int(ops[0]), labels), formatOperand(ops[1]), formatOperand(ops[2]))

	// Op R
	case vm.OpReturn:
		out = fmt.Sprintf("%d: %s %s", ip, opcode, formatOperand(ops[0]))

	// Op R Arg
	case vm.OpDataSet, vm.OpDataSetCollector, vm.OpDataSetSorter:
		out = fmt.Sprintf("%d: %s %s %s", ip, opcode, formatOperand(ops[0]), formatArgument(ops[1]))

	// Op R C
	case vm.OpLoadConst:
		cIdx := ops[1].Constant()
		comment := constValue(p, cIdx)
		out = fmt.Sprintf("%d: %s %s %s ; %s", ip, opcode, formatOperand(ops[0]), formatOperand(ops[1]), comment)

	// Op R R
	case vm.OpIter, vm.OpIterValue, vm.OpIterKey, vm.OpMove, vm.OpPush:
		out = fmt.Sprintf("%d: %s %s %s", ip, opcode, formatOperand(ops[0]), formatOperand(ops[1]))

	// Op R R R
	default:
		out = fmt.Sprintf("%d: %s %s %s %s", ip, opcode, formatOperand(ops[0]), formatOperand(ops[1]), formatOperand(ops[2]))
	}

	if loc := formatLocation(p, ip); loc != "" {
		out += " " + loc
	}

	return out
}
