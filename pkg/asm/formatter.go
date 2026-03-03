package asm

import (
	"fmt"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

// labelOrAddr returns a label name if one exists for the given address; otherwise just the number.
func labelOrAddr(pos int, labels map[int]string) string {
	if label, ok := labels[pos]; ok {
		return label
	}

	return fmt.Sprintf("%d", pos)
}

// constantAsText formats a constant value as a string.
func constantAsText(constant runtime.Value) string {
	if runtime.IsNumber(constant) {
		return fmt.Sprintf("%d", constant)
	}

	return fmt.Sprintf("%q", constant.String())
}

// constValue renders the constant at a given index from the program.
func constValue(p *bytecode.Program, idx int) string {
	if idx >= 0 && idx < len(p.Constants) {
		return constantAsText(p.Constants[idx])
	}

	return "<invalid>"
}

func formatProgram(p *bytecode.Program) string {
	return fmt.Sprintf(".prog %d %d %d", len(p.Bytecode), len(p.Constants), len(p.Params))
}

// formatLocation returns a line/col comment if available for the given instruction.
func formatLocation(p *bytecode.Program, ip int) string {
	//if ip < len(p.Locations) {
	//	loc := p.Locations[ip]
	//
	//	return fmt.Sprintf("; line %d col %d", loc.Line, loc.Column)
	//}

	return ""
}

// formatParam generates comments mapping register indices to parameter names.
func formatParam(name string) string {
	return fmt.Sprintf(".param %s", name)
}

// formatFunction generates comments for the functions defined in the program.
func formatFunction(name string, args int) string {
	return fmt.Sprintf(".func %s %d", name, args)
}

// formatUdf generates comments for the UDF table entries.
func formatUdf(id int, udf bytecode.UDF) string {
	return fmt.Sprintf(".udf %d %s %d %d %d", id, udf.Name, udf.Entry, udf.Registers, udf.Params)
}

// formatConstant generates a comment for a constant value in the program.
func formatConstant(constant runtime.Value) string {
	return fmt.Sprintf(".const %s", constantAsText(constant))
}

func formatOperand(op bytecode.Operand) string {
	if op.IsRegister() {
		return fmt.Sprintf("R%d", op.Register())
	}

	return fmt.Sprintf("C%d", op.Constant())
}

func formatArgument(op bytecode.Operand) string {
	return fmt.Sprintf("%d", op.Register())
}
