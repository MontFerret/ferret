package asm

import (
	"bytes"
	"fmt"
	"sort"
	"strings"
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
	udfLabels := collectUdfEntryLabels(p)

	var buf bytes.Buffer
	w := tabwriter.NewWriter(&buf, 0, 4, 2, ' ', 0)

	// Header: version and source info
	_, _ = fmt.Fprintf(w, ".isa %s\n", formatVersionNum(p.ISAVersion))
	_, _ = fmt.Fprintf(w, ".asm %s\n", formatVersionNum(Version))

	_, _ = fmt.Fprintln(w)

	_, _ = fmt.Fprintln(w, ".meta")
	_, _ = fmt.Fprintf(w, "  %s\n", formatMetaCompilerRow(p.Metadata.CompilerVersion))
	_, _ = fmt.Fprintf(w, "  %s\n", formatMetaOptimizationRow(p.Metadata.OptimizationLevel))

	_, _ = fmt.Fprintln(w)

	// Header: program info
	_, _ = fmt.Fprintln(w, formatProgram(p))

	writeSection := func(name string, rows []string) {
		if len(rows) == 0 {
			return
		}

		_, _ = fmt.Fprintln(w)
		_, _ = fmt.Fprintln(w, name)

		for _, row := range rows {
			_, _ = fmt.Fprintf(w, "  %s\n", row)
		}
	}

	paramRows := make([]string, 0, len(p.Params))
	for _, name := range p.Params {
		paramRows = append(paramRows, formatParamRow(name))
	}

	constRows := make([]string, 0, len(p.Constants))
	for _, constant := range p.Constants {
		constRows = append(constRows, formatConstantRow(constant))
	}

	udfRows := make([]string, 0, len(p.Functions.UserDefined))
	for id, udf := range p.Functions.UserDefined {
		udfRows = append(udfRows, formatUdfRow(id, udf))
	}

	funcRows := make([]string, 0, len(p.Functions.Host))
	if len(p.Functions.Host) > 0 {
		names := make([]string, 0, len(p.Functions.Host))
		for name := range p.Functions.Host {
			names = append(names, name)
		}

		sort.Strings(names)

		for _, name := range names {
			funcRows = append(funcRows, formatFunctionRow(name, p.Functions.Host[name]))
		}
	}

	writeSection(".params", paramRows)
	writeSection(".consts", constRows)
	writeSection(".udf", udfRows)
	writeSection(".func", funcRows)

	if len(p.Bytecode) > 0 {
		_, _ = fmt.Fprintln(w)
		_, _ = fmt.Fprintln(w, ".entry")
	}

	// Body: disassembly
	bodyStarted := false
	for ip, instr := range p.Bytecode {
		emitted := make(map[string]struct{}, 4)
		ipLabels := make([]string, 0, 4)

		if label, ok := labels[ip]; ok {
			ipLabels = append(ipLabels, label)
			emitted[label] = struct{}{}
		}

		for _, label := range udfLabels.entries[ip] {
			if _, ok := emitted[label]; ok {
				continue
			}

			ipLabels = append(ipLabels, label)
			emitted[label] = struct{}{}
		}

		if len(ipLabels) > 0 {
			if bodyStarted {
				_, _ = fmt.Fprintln(w)
			}

			formatted := make([]string, 0, len(ipLabels))
			for _, label := range ipLabels {
				formatted = append(formatted, formatLabelDefinition(label))
			}

			_, _ = fmt.Fprintf(w, "%s\n", strings.Join(formatted, ", "))
		}

		var prev *bytecode.Instruction
		if ip > 0 {
			prev = &p.Bytecode[ip-1]
		}

		_, _ = fmt.Fprintf(w, "\t%s\n", disasmLine(ip, instr, p, labels, prev))
		bodyStarted = true
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

// formatLabelDefinition strips the label-reference prefix from label definitions.
func formatLabelDefinition(label string) string {
	return strings.TrimPrefix(label, "@")
}

type udfEntryLabels struct {
	entries map[int][]string
}

func collectUdfEntryLabels(p *bytecode.Program) udfEntryLabels {
	result := udfEntryLabels{
		entries: make(map[int][]string),
	}

	if p == nil || len(p.Bytecode) == 0 || len(p.Functions.UserDefined) == 0 {
		return result
	}

	type udfEntry struct {
		id    int
		name  string
		entry int
	}

	entries := make([]udfEntry, 0, len(p.Functions.UserDefined))
	for id, udf := range p.Functions.UserDefined {
		if udf.Entry < 0 || udf.Entry >= len(p.Bytecode) {
			continue
		}

		entries = append(entries, udfEntry{
			id:    id,
			name:  udf.Name,
			entry: udf.Entry,
		})
	}

	if len(entries) == 0 {
		return result
	}

	sort.Slice(entries, func(i, j int) bool {
		if entries[i].entry == entries[j].entry {
			return entries[i].id < entries[j].id
		}

		return entries[i].entry < entries[j].entry
	})

	for i := range entries {
		cur := entries[i]
		entryLabel := fmt.Sprintf("@udf.%d.%s", cur.id, cur.name)
		result.entries[cur.entry] = append(result.entries[cur.entry], entryLabel)
	}

	return result
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
	case bytecode.OpLoadConst:
		cIdx := ops[1].Constant()
		comment := constValue(p, cIdx)
		out = fmt.Sprintf("%d: %s %s %s %s", ip, opcode, formatOperand(ops[0]), formatOperand(ops[1]), formatComment(comment))
	case bytecode.OpLoadParam:
		slot := int(ops[1]) - 1
		comment := "<invalid>"

		if slot >= 0 && slot < len(p.Params) {
			comment = "@" + p.Params[slot]
		}

		out = fmt.Sprintf("%d: %s %s %s %s", ip, opcode, formatOperand(ops[0]), formatArgument(ops[1]), formatComment(comment))

	case bytecode.OpFail:
		out = fmt.Sprintf("%d: %s %s", ip, opcode, formatOperand(ops[0]))
		if ops[0].IsConstant() {
			out += formatComment(constValue(p, ops[0].Constant()))
		}

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

	if isHostCallOpcode(opcode) {
		if comment := hostCallComment(p, instr, prev); comment != "" {
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

func isHostCallOpcode(op bytecode.Opcode) bool {
	switch op {
	case bytecode.OpHCall, bytecode.OpProtectedHCall:
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

func hostCallComment(p *bytecode.Program, instr bytecode.Instruction, prev *bytecode.Instruction) string {
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

	name, ok := p.Constants[idx].(runtime.String)
	if !ok {
		return ""
	}

	return fmt.Sprintf("host %s", name)
}
