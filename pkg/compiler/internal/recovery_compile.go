package internal

import (
	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/compiler/internal/core"
)

func compileWithRecoveryPlan(
	ctx *CompilerContext,
	plan recoveryPlan,
	jumpMode catchJumpMode,
	compile func() bytecode.Operand,
) bytecode.Operand {
	if ctx == nil || compile == nil {
		return bytecode.NoopOperand
	}

	if plan.onError == nil || plan.onError.actionKind == recoveryActionFail {
		return compile()
	}

	switch plan.onError.actionKind {
	case recoveryActionReturn:
		if hasErrorReturnNoneHandler(plan) {
			out := compileWithErrorPolicy(ctx, errorPolicySuppress, jumpMode, compile)
			return widenRecoveryResultType(ctx, out, plan)
		}

		startCatch := ctx.Emitter.Size()
		out := ensureRecoveryRegister(ctx, compile())
		endCatchExclusive := ctx.Emitter.Size()

		if out == bytecode.NoopOperand || endCatchExclusive <= startCatch {
			return out
		}

		endCatch := endCatchExclusive - 1
		endLabel := ctx.Emitter.NewLabel("recovery", "end")

		ctx.Emitter.EmitJump(endLabel)
		handlerPC := ctx.Emitter.Size()

		fallback := ctx.ExprCompiler.Compile(plan.onError.expr)
		ctx.EmitMoveAuto(out, ensureRecoveryRegister(ctx, fallback))
		ctx.Emitter.MarkLabel(endLabel)

		ctx.CatchTable.Push(startCatch, endCatch, handlerPC)

		return widenRecoveryResultType(ctx, out, plan)
	case recoveryActionRetry:
		return compileWithRetryRecoveryPlan(ctx, plan.onError, compile)
	default:
		return compile()
	}
}

func ensureRecoveryRegister(ctx *CompilerContext, op bytecode.Operand) bytecode.Operand {
	if ctx == nil || op == bytecode.NoopOperand || op.IsRegister() {
		return op
	}

	dst := ctx.Registers.Allocate()
	ctx.Emitter.EmitLoadConst(dst, op)
	ctx.Types.Set(dst, operandType(ctx, op))

	return dst
}

func widenRecoveryResultType(ctx *CompilerContext, out bytecode.Operand, plan recoveryPlan) bytecode.Operand {
	if ctx == nil || !out.IsRegister() || !recoveryPlanHasReturnHandler(plan) {
		return out
	}

	ctx.Types.Set(out, core.TypeAny)

	return out
}
