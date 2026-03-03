package asm

import (
	"bytes"
	"fmt"
	"text/tabwriter"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

const Version = 1

// Disassemble returns a human-readable disassembly of the given program.
func Disassemble(p *bytecode.Program, options ...DisassemblerOption) (string, error) {
	if p == nil {
		return "", ErrInvalidProgram
	}

	newDisassemblerOptions(options...)

	labels := collectLabels(p.Bytecode, p.Metadata.Labels)

	var buf bytes.Buffer
	w := tabwriter.NewWriter(&buf, 0, 4, 2, ' ', 0)

	// Header: version and source info
	_, _ = fmt.Fprintf(w, ".isa %s\n", formatVersionNum(p.ISAVersion))
	_, _ = fmt.Fprintf(w, ".asm %s\n", formatVersionNum(Version))
	_, _ = fmt.Fprintf(w, ".meta compiler %s\n", formatVersion(p.Metadata.CompilerVersion))
	_, _ = fmt.Fprintf(w, ".meta opt O%d\n", p.Metadata.OptimizationLevel)

	_, _ = fmt.Fprintln(w)

	// Header: program info
	_, _ = fmt.Fprintln(w, formatProgram(p))

	// Header: functions
	for name, args := range p.Functions.Host {
		_, _ = fmt.Fprintln(w, formatFunctionHeader(name, args))
	}

	// Header: UDFs
	for id, udf := range p.Functions.UserDefined {
		_, _ = fmt.Fprintln(w, formatUdfHeader(id, udf))
	}

	// Header: params
	for _, name := range p.Params {
		_, _ = fmt.Fprintln(w, formatParamHeader(name))
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

		var prev *bytecode.Instruction
		if ip > 0 {
			prev = &p.Bytecode[ip-1]
		}

		_, _ = fmt.Fprintf(w, "\t%s\n", disasmLine(ip, instr, p, labels, prev))
	}

	_ = w.Flush()

	return buf.String(), nil
}

// collectLabels identifies jump targets and assigns symbolic labels to them.
func collectLabels(instructions []bytecode.Instruction, names map[int]string) map[int]string {
	labels := make(map[int]string)
	counter := 0

	// Iterate through the labels in the program to initialize the labels map
	for target, name := range names {
		labels[target] = fmt.Sprintf("@%s", name)
	}

	// Collect unmarked jump targets
	for _, instr := range instructions {
		switch instr.Opcode {
		case bytecode.OpJump, bytecode.OpJumpIfFalse, bytecode.OpJumpIfTrue, bytecode.OpJumpIfNone, bytecode.OpJumpIfNe, bytecode.OpJumpIfNeConst, bytecode.OpJumpIfEq, bytecode.OpJumpIfEqConst, bytecode.OpJumpIfMissingProperty, bytecode.OpJumpIfMissingPropertyConst, bytecode.OpIterNext, bytecode.OpIterSkip, bytecode.OpIterLimit:
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
func disasmLine(ip int, instr bytecode.Instruction, p *bytecode.Program, labels map[int]string, prev *bytecode.Instruction) string {
	ops := instr.Operands
	var out string

	opcode := instr.Opcode

	switch opcode {
	// Op Jmp
	case bytecode.OpJump:
		out = fmt.Sprintf("%d: %s %s", ip, opcode, labelOrAddr(int(ops[0]), labels))

	// Op Jmp R
	case bytecode.OpJumpIfTrue, bytecode.OpJumpIfFalse, bytecode.OpJumpIfNone, bytecode.OpIterNext:
		out = fmt.Sprintf("%d: %s %s %s", ip, opcode, labelOrAddr(int(ops[0]), labels), formatOperand(ops[1]))

	// Op Jmp R, R
	case bytecode.OpJumpIfNe, bytecode.OpJumpIfNeConst, bytecode.OpJumpIfEq, bytecode.OpJumpIfEqConst, bytecode.OpJumpIfMissingProperty, bytecode.OpJumpIfMissingPropertyConst, bytecode.OpIterSkip, bytecode.OpIterLimit:
		out = fmt.Sprintf("%d: %s %s %s %s", ip, opcode, labelOrAddr(int(ops[0]), labels), formatOperand(ops[1]), formatOperand(ops[2]))

	// Op R
	case bytecode.OpLoadNone, bytecode.OpLoadZero,
		bytecode.OpClose, bytecode.OpSleep, bytecode.OpRand, bytecode.OpIncr, bytecode.OpDecr, bytecode.OpReturn:
		out = fmt.Sprintf("%d: %s %s", ip, opcode, formatOperand(ops[0]))

	// Op R Arg
	case bytecode.OpLoadBool, bytecode.OpLoadArray, bytecode.OpLoadObject, bytecode.OpDataSet, bytecode.OpDataSetCollector, bytecode.OpDataSetSorter:
		out = fmt.Sprintf("%d: %s %s %s", ip, opcode, formatOperand(ops[0]), formatArgument(ops[1]))

	// Op R R Arg
	case bytecode.OpConcat:
		out = fmt.Sprintf("%d: %s %s %s %s", ip, opcode, formatOperand(ops[0]), formatOperand(ops[1]), formatArgument(ops[2]))

	case bytecode.OpAddConst:
		if ops[2].IsConstant() {
			cIdx := ops[2].Constant()
			comment := constValue(p, cIdx)
			out = fmt.Sprintf("%d: %s %s %s %s %s", ip, opcode, formatOperand(ops[0]), formatOperand(ops[1]), formatOperand(ops[2]), formatComment(comment))
		} else {
			out = fmt.Sprintf("%d: %s %s %s %s", ip, opcode, formatOperand(ops[0]), formatOperand(ops[1]), formatOperand(ops[2]))
		}

	// Op R C
	case bytecode.OpLoadConst, bytecode.OpLoadParam:
		cIdx := ops[1].Constant()
		comment := constValue(p, cIdx)
		out = fmt.Sprintf("%d: %s %s %s %s", ip, opcode, formatOperand(ops[0]), formatOperand(ops[1]), formatComment(comment))

	// Op R R
	case bytecode.OpMove, bytecode.OpLength, bytecode.OpType, bytecode.OpExists,
		bytecode.OpIter, bytecode.OpIterValue, bytecode.OpIterKey, bytecode.OpPush, bytecode.OpArrayPush:
		out = fmt.Sprintf("%d: %s %s %s", ip, opcode, formatOperand(ops[0]), formatOperand(ops[1]))

	case bytecode.OpHCall, bytecode.OpProtectedHCall, bytecode.OpCall, bytecode.OpProtectedCall, bytecode.OpTailCall:
		out = fmt.Sprintf("%d: %s %s", ip, opcode, formatOperand(ops[0]))

		if ops[1] != bytecode.NoopOperand || ops[2] != bytecode.NoopOperand {
			out += fmt.Sprintf(" %s %s", formatOperand(ops[1]), formatOperand(ops[2]))
		}

	// Op R R R
	default:
		out = fmt.Sprintf("%d: %s %s %s %s", ip, opcode, formatOperand(ops[0]), formatOperand(ops[1]), formatOperand(ops[2]))
	}

	if isUdfCallOpcode(opcode) {
		if comment := udfCallComment(p, instr, prev); comment != "" {
			out += formatComment(comment)
		}
	}

	return out
}

func isUdfCallOpcode(op bytecode.Opcode) bool {
	switch op {
	case bytecode.OpCall, bytecode.OpProtectedCall, bytecode.OpTailCall:
		return true
	default:
		return false
	}
}

func udfCallComment(p *bytecode.Program, instr bytecode.Instruction, prev *bytecode.Instruction) string {
	if p == nil || prev == nil {
		return ""
	}

	if prev.Opcode != bytecode.OpLoadConst {
		return ""
	}

	if prev.Operands[0] != instr.Operands[0] {
		return ""
	}

	if !prev.Operands[1].IsConstant() {
		return ""
	}

	idx := prev.Operands[1].Constant()
	if idx < 0 || idx >= len(p.Constants) {
		return ""
	}

	idVal, ok := p.Constants[idx].(runtime.Int)
	if !ok {
		return ""
	}

	id := int(idVal)
	if id < 0 || id >= len(p.Functions.UserDefined) {
		return ""
	}

	return fmt.Sprintf("udf %s", p.Functions.UserDefined[id].Name)
}
