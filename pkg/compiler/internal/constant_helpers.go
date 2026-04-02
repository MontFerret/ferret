package internal

import (
	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

func loadConstant(ctx *CompilerContext, value runtime.Value) bytecode.Operand {
	reg := ctx.Registers.Allocate()
	ctx.Emitter.EmitLoadConst(reg, ctx.Symbols.AddConstant(value))
	ctx.Types.Set(reg, valueTypeFromRuntime(value))

	return reg
}
