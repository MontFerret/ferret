package internal

import (
	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

// ensureOperandRegister materializes a constant operand into a register while
// preserving tracked type facts. Noop operands remain untouched.
func ensureOperandRegister(ctx *CompilationSession, facts *TypeFacts, op bytecode.Operand) bytecode.Operand {
	if ctx == nil || op == bytecode.NoopOperand || op.IsRegister() {
		return op
	}

	reg := ctx.Function.Registers.Allocate()
	ctx.Program.Emitter.EmitLoadConst(reg, op)
	ctx.Function.Types.Set(reg, facts.OperandType(op))

	return reg
}

// ensureDispatchOperandRegister normalizes missing DISPATCH operands to NONE
// before applying the standard register-materialization path.
func ensureDispatchOperandRegister(ctx *CompilationSession, facts *TypeFacts, op bytecode.Operand) bytecode.Operand {
	if op == bytecode.NoopOperand {
		op = facts.LoadConstant(runtime.None)
	}

	return ensureOperandRegister(ctx, facts, op)
}
