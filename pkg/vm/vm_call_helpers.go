package vm

import "github.com/MontFerret/ferret/v2/pkg/bytecode"

func callArgRange(src1, src2 bytecode.Operand) (int, int, bool) {
	if !src1.IsRegister() || !src2.IsRegister() {
		return 0, 0, false
	}

	start := src1.Register()
	end := src2.Register()

	if start <= 0 || end < start {
		return 0, 0, false
	}

	return start, end, true
}

func callArgInfo(src1, src2 bytecode.Operand) (int, int) {
	start, end, ok := callArgRange(src1, src2)
	if !ok {
		return 0, 0
	}

	return start, end - start + 1
}
