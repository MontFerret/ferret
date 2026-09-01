package vm

import (
	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

type callDescriptor struct {
	DisplayName      string
	PC               int
	Dst              bytecode.Operand
	ID               bytecode.FunctionID
	ArgCount         int
	ArgStart         int
	RecoveryBoundary bool
	CallSitePC       int
}

func callArgCount(src1, src2 bytecode.Operand) int {
	argCount := 0

	if src1.IsRegister() && src2.IsRegister() {
		start := src1.Register()
		end := src2.Register()

		if start > 0 && end >= start {
			argCount = end - start + 1
		}
	}

	return argCount
}

func getUDFID(val runtime.Value) (bytecode.FunctionID, error) {
	idVal, ok := val.(runtime.Int)
	if !ok {
		return bytecode.NoFunction, ErrInvalidFunctionName
	}

	return bytecode.FunctionID(idVal), nil
}
