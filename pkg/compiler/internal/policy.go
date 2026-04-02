package internal

import (
	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/compiler/internal/core"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

type OperationPolicyCompiler struct {
	ctx *CompilerContext
}

func NewOperationPolicyCompiler(ctx *CompilerContext) *OperationPolicyCompiler {
	return &OperationPolicyCompiler{ctx: ctx}
}

func (c *OperationPolicyCompiler) CompileWithErrorPolicy(policy core.ErrorPolicy, jumpMode core.CatchJumpMode, compile func() bytecode.Operand) bytecode.Operand {
	if compile == nil || policy != core.ErrorPolicySuppress {
		return compile()
	}

	startCatch := c.ctx.Emitter.Size()
	out := compile()
	endCatchExclusive := c.ctx.Emitter.Size()
	if endCatchExclusive <= startCatch {
		return out
	}

	endCatch := endCatchExclusive - 1
	jump := -1

	if jumpMode == core.CatchJumpModeEnd {
		jump = endCatch
	} else {
		endLabel := c.ctx.Emitter.NewLabel("error", "suppress", "end")
		c.ctx.Emitter.EmitJump(endLabel)
		jump = c.ctx.Emitter.Size()
		c.ctx.Emitter.EmitJump(endLabel)
		c.ctx.Emitter.MarkLabel(endLabel)
	}

	c.ctx.CatchTable.Push(startCatch, endCatch, jump)

	return out
}

func (c *OperationPolicyCompiler) CompileWithRecoveryPlan(
	plan core.RecoveryPlan,
	jumpMode core.CatchJumpMode,
	compile func() bytecode.Operand,
) bytecode.Operand {
	if compile == nil {
		return bytecode.NoopOperand
	}

	if plan.OnError == nil || plan.OnError.ActionKind == core.RecoveryActionFail {
		return compile()
	}

	switch plan.OnError.ActionKind {
	case core.RecoveryActionReturn:
		if hasErrorReturnNoneHandler(plan) {
			out := c.CompileWithErrorPolicy(core.ErrorPolicySuppress, jumpMode, compile)
			return widenRecoveryResultType(c.ctx, out, plan)
		}

		startCatch := c.ctx.Emitter.Size()
		out := ensureRecoveryRegister(c.ctx, compile())
		endCatchExclusive := c.ctx.Emitter.Size()

		if out == bytecode.NoopOperand || endCatchExclusive <= startCatch {
			return out
		}

		endCatch := endCatchExclusive - 1
		endLabel := c.ctx.Emitter.NewLabel("recovery", "end")

		c.ctx.Emitter.EmitJump(endLabel)
		handlerPC := c.ctx.Emitter.Size()

		fallback := c.ctx.ExprCompiler.Compile(plan.OnError.Expr)
		emitMoveAuto(c.ctx, out, ensureRecoveryRegister(c.ctx, fallback))
		c.ctx.Emitter.MarkLabel(endLabel)

		c.ctx.CatchTable.Push(startCatch, endCatch, handlerPC)

		return widenRecoveryResultType(c.ctx, out, plan)
	case core.RecoveryActionRetry:
		return c.CompileWithRecoveryHandler(plan.OnError, compile)
	default:
		return compile()
	}
}

func (c *OperationPolicyCompiler) CompileWithRecoveryHandler(
	handler *core.RecoveryHandler,
	compile func() bytecode.Operand,
) bytecode.Operand {
	if compile == nil || !recoveryHandlerRetries(handler) {
		return bytecode.NoopOperand
	}

	retry := handler.Retry
	if retry == nil {
		return compile()
	}

	if retry.FinalActionKind != core.RecoveryActionReturn && retry.Count <= 0 {
		return compile()
	}

	resultReg := bytecode.NoopOperand
	zeroReg := loadConstant(c.ctx, runtime.ZeroInt)
	retriesRemainingReg := loadConstant(c.ctx, runtime.NewInt(retry.Count))

	state := initRetryDelayState(c.ctx, retry)
	retryStart := c.ctx.Emitter.NewLabel("recovery", "retry", "start")
	endLabel := c.ctx.Emitter.NewLabel("recovery", "retry", "end")
	finalAttemptLabel := c.ctx.Emitter.NewLabel("recovery", "retry", "final")

	c.ctx.Emitter.MarkLabel(retryStart)
	startCatch := c.ctx.Emitter.Size()
	protectedOut := ensureRecoveryRegister(c.ctx, compile())
	endCatchExclusive := c.ctx.Emitter.Size()

	if protectedOut == bytecode.NoopOperand || endCatchExclusive <= startCatch {
		return protectedOut
	}

	resultReg = protectedOut
	endCatch := endCatchExclusive - 1

	c.ctx.Emitter.EmitJump(endLabel)
	handlerPC := c.ctx.Emitter.Size()

	retriesAvailableReg := c.ctx.Registers.Allocate()
	c.ctx.Emitter.EmitGt(retriesAvailableReg, retriesRemainingReg, zeroReg)

	onExhausted := c.ctx.Emitter.NewLabel("recovery", "retry", "exhausted")
	c.ctx.Emitter.EmitJumpIfFalse(retriesAvailableReg, onExhausted)
	c.ctx.Emitter.EmitA(bytecode.OpDecr, retriesRemainingReg)
	emitRecoveryRetryDelay(c.ctx, retry, state)

	if retry.FinalActionKind == core.RecoveryActionReturn {
		c.ctx.Emitter.EmitJump(retryStart)
	} else {
		moreProtectedReg := c.ctx.Registers.Allocate()
		c.ctx.Emitter.EmitGt(moreProtectedReg, retriesRemainingReg, zeroReg)
		c.ctx.Emitter.EmitJumpIfTrue(moreProtectedReg, retryStart)
		c.ctx.Emitter.EmitJump(finalAttemptLabel)
	}

	c.ctx.Emitter.MarkLabel(onExhausted)
	if retry.FinalActionKind == core.RecoveryActionReturn {
		fallback := c.ctx.ExprCompiler.Compile(retry.FinalExpr)
		emitMoveAuto(c.ctx, resultReg, ensureRecoveryRegister(c.ctx, fallback))
		c.ctx.Emitter.EmitJump(endLabel)
	} else {
		c.ctx.Emitter.EmitJump(finalAttemptLabel)
	}

	c.ctx.CatchTable.Push(startCatch, endCatch, handlerPC)

	if retry.FinalActionKind != core.RecoveryActionReturn {
		c.ctx.Emitter.MarkLabel(finalAttemptLabel)
		finalOut := ensureRecoveryRegister(c.ctx, compile())
		if finalOut != bytecode.NoopOperand && finalOut != resultReg {
			emitMoveAuto(c.ctx, resultReg, finalOut)
		}
	}

	c.ctx.Emitter.MarkLabel(endLabel)

	return widenRecoveryResultType(c.ctx, resultReg, core.RecoveryPlan{OnError: handler})
}
