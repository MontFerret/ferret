package vm

import (
	"strings"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/vm/internal/data"
)

func normalizeValue(val runtime.Value) runtime.Value {
	if val == nil {
		return runtime.None
	}

	return val
}

func readOperandValue(reg []runtime.Value, constants []runtime.Value, operand bytecode.Operand) runtime.Value {
	if operand.IsConstant() {
		return constants[operand.Constant()]
	}

	return reg[operand]
}

func concatStrings(reg []runtime.Value, dst, src1, src2 bytecode.Operand) {
	start := int(src1)
	count := int(src2)

	if count <= 0 {
		reg[dst] = runtime.EmptyString
		return
	}

	if count == 1 {
		reg[dst] = runtime.ToString(reg[start])
		return
	}

	if count == 2 {
		s1 := runtime.ToString(reg[start])
		s2 := runtime.ToString(reg[start+1])

		if s1 == runtime.EmptyString {
			reg[dst] = s2
			return
		}

		if s2 == runtime.EmptyString {
			reg[dst] = s1
			return
		}

		reg[dst] = runtime.NewString(string(s1) + string(s2))
		return
	}

	if count == 3 {
		s1 := runtime.ToString(reg[start])
		s2 := runtime.ToString(reg[start+1])
		s3 := runtime.ToString(reg[start+2])

		if s1 == runtime.EmptyString {
			if s2 == runtime.EmptyString {
				reg[dst] = s3
				return
			}
			if s3 == runtime.EmptyString {
				reg[dst] = s2
				return
			}
		} else if s2 == runtime.EmptyString {
			if s3 == runtime.EmptyString {
				reg[dst] = s1
				return
			}
		}

		reg[dst] = runtime.NewString(string(s1) + string(s2) + string(s3))
		return
	}

	parts := make([]runtime.String, count)
	totalLen := 0

	for i := 0; i < count; i++ {
		s := runtime.ToString(reg[start+i])
		parts[i] = s
		totalLen += len(s)
	}

	if totalLen == 0 {
		reg[dst] = runtime.EmptyString
		return
	}

	var b strings.Builder
	b.Grow(totalLen)

	for i := 0; i < count; i++ {
		if parts[i] == runtime.EmptyString {
			continue
		}

		b.WriteString(string(parts[i]))
	}

	reg[dst] = runtime.NewString(b.String())
}

func buildCatchByPC(bytecodeLen int, catches []bytecode.Catch) []int {
	if bytecodeLen <= 0 {
		return nil
	}

	catchByPC := make([]int, bytecodeLen)

	for i := range catchByPC {
		catchByPC[i] = -1
	}

	for i, pair := range catches {
		start, end := pair[0], pair[1]

		if start < 0 {
			start = 0
		}

		if end >= bytecodeLen {
			end = bytecodeLen - 1
		}

		for pc := start; pc <= end; pc++ {
			if catchByPC[pc] == -1 {
				catchByPC[pc] = i
			}
		}
	}

	return catchByPC
}

func buildExecInstructions(code []bytecode.Instruction) []data.ExecInstruction {
	instructions := make([]data.ExecInstruction, len(code))

	for i := range code {
		instructions[i] = data.ExecInstruction{
			Instruction: code[i],
		}
	}

	return instructions
}

func maxUDFRegisters(udfs []bytecode.UDF) int {
	maxUDFRegs := 0

	for i := range udfs {
		if udfs[i].Registers > maxUDFRegs {
			maxUDFRegs = udfs[i].Registers
		}
	}

	return maxUDFRegs
}
