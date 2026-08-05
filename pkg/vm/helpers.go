package vm

import (
	"strings"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
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

func concatStrings(reg []runtime.Value, src1, src2 bytecode.Operand) (runtime.Value, error) {
	start := int(src1)
	count := int(src2)

	if count <= 0 {
		return runtime.EmptyString, nil
	}

	if count == 1 {
		return runtime.ToString(reg[start]), nil
	}

	if count == 2 {
		s1 := runtime.ToString(reg[start])
		s2 := runtime.ToString(reg[start+1])

		if s1 == runtime.EmptyString {
			return s2, nil
		}

		if s2 == runtime.EmptyString {
			return s1, nil
		}

		return runtime.NewString(string(s1) + string(s2)), nil
	}

	if count == 3 {
		s1 := runtime.ToString(reg[start])
		s2 := runtime.ToString(reg[start+1])
		s3 := runtime.ToString(reg[start+2])

		if s1 == runtime.EmptyString {
			if s2 == runtime.EmptyString {
				return s3, nil
			}

			if s3 == runtime.EmptyString {
				return s2, nil
			}
		} else if s2 == runtime.EmptyString {
			if s3 == runtime.EmptyString {
				return s1, nil
			}
		}

		return runtime.NewString(string(s1) + string(s2) + string(s3)), nil
	}

	parts := make([]runtime.String, count)
	totalLen := 0

	for i := 0; i < count; i++ {
		s := runtime.ToString(reg[start+i])
		parts[i] = s
		totalLen += len(s)
	}

	if totalLen == 0 {
		return runtime.EmptyString, nil
	}

	var b strings.Builder
	b.Grow(totalLen)

	for i := 0; i < count; i++ {
		if parts[i] == runtime.EmptyString {
			continue
		}

		b.WriteString(string(parts[i]))
	}

	return runtime.NewString(b.String()), nil
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
