package internal

import (
	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

func compileWithRetryRecoveryPlan(
	ctx *CompilerContext,
	handler *recoveryHandler,
	compile func() bytecode.Operand,
) bytecode.Operand {
	if ctx == nil || compile == nil || !recoveryHandlerRetries(handler) {
		return bytecode.NoopOperand
	}

	retry := handler.retry
	if retry == nil {
		return compile()
	}

	if retry.finalActionKind != recoveryActionReturn && retry.count <= 0 {
		return compile()
	}

	resultReg := bytecode.NoopOperand
	zeroReg := loadConstant(ctx, runtime.ZeroInt)
	retriesRemainingReg := loadConstant(ctx, runtime.NewInt(retry.count))

	state := initRecoveryRetryDelayState(ctx, retry)
	retryStart := ctx.Emitter.NewLabel("recovery", "retry", "start")
	endLabel := ctx.Emitter.NewLabel("recovery", "retry", "end")
	finalAttemptLabel := ctx.Emitter.NewLabel("recovery", "retry", "final")

	ctx.Emitter.MarkLabel(retryStart)
	startCatch := ctx.Emitter.Size()
	protectedOut := ensureRecoveryRegister(ctx, compile())
	endCatchExclusive := ctx.Emitter.Size()

	if protectedOut == bytecode.NoopOperand || endCatchExclusive <= startCatch {
		return protectedOut
	}

	resultReg = protectedOut
	endCatch := endCatchExclusive - 1

	ctx.Emitter.EmitJump(endLabel)
	handlerPC := ctx.Emitter.Size()

	retriesAvailableReg := ctx.Registers.Allocate()
	ctx.Emitter.EmitGt(retriesAvailableReg, retriesRemainingReg, zeroReg)

	onExhausted := ctx.Emitter.NewLabel("recovery", "retry", "exhausted")
	ctx.Emitter.EmitJumpIfFalse(retriesAvailableReg, onExhausted)
	ctx.Emitter.EmitA(bytecode.OpDecr, retriesRemainingReg)
	emitRecoveryRetryDelay(ctx, retry, state)

	if retry.finalActionKind == recoveryActionReturn {
		ctx.Emitter.EmitJump(retryStart)
	} else {
		moreProtectedReg := ctx.Registers.Allocate()
		ctx.Emitter.EmitGt(moreProtectedReg, retriesRemainingReg, zeroReg)
		ctx.Emitter.EmitJumpIfTrue(moreProtectedReg, retryStart)
		ctx.Emitter.EmitJump(finalAttemptLabel)
	}

	ctx.Emitter.MarkLabel(onExhausted)
	if retry.finalActionKind == recoveryActionReturn {
		fallback := ctx.ExprCompiler.Compile(retry.finalExpr)
		ctx.EmitMoveAuto(resultReg, ensureRecoveryRegister(ctx, fallback))
		ctx.Emitter.EmitJump(endLabel)
	} else {
		ctx.Emitter.EmitJump(finalAttemptLabel)
	}

	ctx.CatchTable.Push(startCatch, endCatch, handlerPC)

	if retry.finalActionKind != recoveryActionReturn {
		ctx.Emitter.MarkLabel(finalAttemptLabel)
		finalOut := ensureRecoveryRegister(ctx, compile())
		if finalOut != bytecode.NoopOperand && finalOut != resultReg {
			ctx.EmitMoveAuto(resultReg, finalOut)
		}
	}

	ctx.Emitter.MarkLabel(endLabel)

	return widenRecoveryResultType(ctx, resultReg, recoveryPlan{onError: handler})
}

type recoveryRetryDelayState struct {
	baseReg    bytecode.Operand
	currentReg bytecode.Operand
	readyReg   bytecode.Operand
}

func initRecoveryRetryDelayState(ctx *CompilerContext, retry *recoveryRetryPlan) recoveryRetryDelayState {
	if ctx == nil || retry == nil || !retry.hasDelay {
		return recoveryRetryDelayState{}
	}

	state := recoveryRetryDelayState{
		baseReg:    ctx.Registers.Allocate(),
		currentReg: ctx.Registers.Allocate(),
		readyReg:   ctx.Registers.Allocate(),
	}

	ctx.Emitter.EmitBoolean(state.readyReg, false)

	return state
}

func emitRecoveryRetryDelay(ctx *CompilerContext, retry *recoveryRetryPlan, state recoveryRetryDelayState) {
	if ctx == nil || retry == nil || !retry.hasDelay {
		return
	}

	delayReady := ctx.Emitter.NewLabel("recovery", "retry", "delay", "ready")
	ctx.Emitter.EmitJumpIfTrue(state.readyReg, delayReady)

	delayValue := ensureRecoveryRegister(ctx, ctx.WaitCompiler.compileDurationClause(retry.delay))
	ctx.EmitMoveAuto(state.baseReg, delayValue)
	ctx.EmitMoveAuto(state.currentReg, state.baseReg)
	ctx.Emitter.EmitBoolean(state.readyReg, true)
	ctx.Emitter.MarkLabel(delayReady)

	ctx.Emitter.EmitA(bytecode.OpSleep, state.currentReg)

	if retry.backoff != waitForBackoffNone {
		ctx.WaitCompiler.emitBackoffUpdate(retry.backoff, state.currentReg, state.baseReg)
	}
}
