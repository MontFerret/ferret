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

func formatComment(comment string) string {
	return fmt.Sprintf(" ; %s", comment)
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
	comment := "regs consts params"

	return fmt.Sprintf(".prog %d %d %d %s", p.Registers, len(p.Constants), len(p.Params), formatComment(comment))
}

func formatVersion(v string) string {
	if v == "" {
		return "-"
	}

	return v
}

func formatVersionNum(v int) string {
	if v <= 0 {
		return formatVersion("")
	}

	return formatVersion(fmt.Sprintf("%d", v))
}

func formatMetaCompilerRow(version string) string {
	return fmt.Sprintf("compiler %s", formatVersion(version))
}

func formatMetaOptimizationRow(level int) string {
	return fmt.Sprintf("opt O%d", level)
}

func formatParamRow(name string) string {
	return name
}

func formatFunctionRow(name string, args int) string {
	return fmt.Sprintf("%s %d ; name params", name, args)
}

func formatUdfRow(id int, udf bytecode.UDF) string {
	comment := "id name entry registers params"
	return fmt.Sprintf("%d %s %d %d %d %s", id, udf.Name, udf.Entry, udf.Registers, udf.Params, formatComment(comment))
}

func formatConstantRow(constant runtime.Value) string {
	return constantAsText(constant)
}

// formatParamHeader is a legacy single-line header helper.
func formatParamHeader(name string) string {
	return fmt.Sprintf(".param %s", formatParamRow(name))
}

// formatFunctionHeader is a legacy single-line header helper.
func formatFunctionHeader(name string, args int) string {
	return fmt.Sprintf(".func %s", formatFunctionRow(name, args))
}

// formatUdfHeader is a legacy single-line header helper.
func formatUdfHeader(id int, udf bytecode.UDF) string {
	return fmt.Sprintf(".udf %s", formatUdfRow(id, udf))
}

// formatConstant is a legacy single-line header helper.
func formatConstant(constant runtime.Value) string {
	return fmt.Sprintf(".const %s", formatConstantRow(constant))
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
