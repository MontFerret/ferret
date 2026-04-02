package internal

import (
	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/compiler/internal/core"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

type (
	RecoveryCompiler struct {
		ctx   *CompilationSession
		front *CompilationFrontend
	}

	// ProtectedRecoverySpec describes an operation that owns its protected region
	// shape but delegates generic ON ERROR recovery orchestration to RecoveryCompiler.
	ProtectedRecoverySpec struct {
		Plan                core.RecoveryPlan
		BuildProtected      func(recoveryLabel, endLabel core.Label) ProtectedRecoveryRegion
		CompileFinalAttempt func() bytecode.Operand
	}

	// ProtectedRecoveryRegion describes the guarded portion of an operation.
	// CatchHandlerPC == -1 means the catch table should jump directly to the
	// generic recovery handler emitted by RecoveryCompiler.
	ProtectedRecoveryRegion struct {
		Result            bytecode.Operand
		StartCatch        int
		EndCatchExclusive int
		CatchHandlerPC    int
	}
)

func NewRecoveryCompiler(ctx *CompilationSession) *RecoveryCompiler {
	return &RecoveryCompiler{ctx: ctx}
}

func (c *RecoveryCompiler) CompileWithErrorPolicy(policy core.ErrorPolicy, jumpMode core.CatchJumpMode, compile func() bytecode.Operand) bytecode.Operand {
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

func (c *RecoveryCompiler) CompileWithRecoveryPlan(
	plan core.RecoveryPlan,
	jumpMode core.CatchJumpMode,
	compile func() bytecode.Operand,
) bytecode.Operand {
	plan = c.NormalizePlan(plan)

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
			return c.WidenResultType(out, plan)
		}

		return c.compileWithErrorReturn(plan, c.directProtectedRegionBuilder(compile))
	case core.RecoveryActionRetry:
		return c.compileWithErrorRetry(plan, c.directProtectedRegionBuilder(compile), compile)
	default:
		return compile()
	}
}

func (c *RecoveryCompiler) CompileWithProtectedRecovery(spec ProtectedRecoverySpec) bytecode.Operand {
	if spec.BuildProtected == nil {
		return bytecode.NoopOperand
	}

	plan := c.NormalizePlan(spec.Plan)
	if plan.OnError == nil || plan.OnError.ActionKind == core.RecoveryActionFail {
		if spec.CompileFinalAttempt != nil {
			return spec.CompileFinalAttempt()
		}

		return bytecode.NoopOperand
	}

	switch plan.OnError.ActionKind {
	case core.RecoveryActionReturn:
		return c.compileWithErrorReturn(plan, spec.BuildProtected)
	case core.RecoveryActionRetry:
		return c.compileWithErrorRetry(plan, spec.BuildProtected, spec.CompileFinalAttempt)
	default:
		if spec.CompileFinalAttempt != nil {
			return spec.CompileFinalAttempt()
		}

		return bytecode.NoopOperand
	}
}

func (c *RecoveryCompiler) CompileWithRecoveryHandler(
	handler *core.RecoveryHandler,
	compile func() bytecode.Operand,
) bytecode.Operand {
	if compile == nil || !recoveryHandlerRetries(handler) {
		return bytecode.NoopOperand
	}

	return c.compileWithErrorRetry(core.RecoveryPlan{OnError: handler}, c.directProtectedRegionBuilder(compile), compile)
}

func (c *RecoveryCompiler) directProtectedRegionBuilder(compile func() bytecode.Operand) func(recoveryLabel, endLabel core.Label) ProtectedRecoveryRegion {
	return func(_ core.Label, endLabel core.Label) ProtectedRecoveryRegion {
		startCatch := c.ctx.Emitter.Size()
		out := c.EnsureRegister(compile())
		endCatchExclusive := c.ctx.Emitter.Size()

		if out == bytecode.NoopOperand || endCatchExclusive <= startCatch {
			return ProtectedRecoveryRegion{
				Result:            out,
				StartCatch:        startCatch,
				EndCatchExclusive: endCatchExclusive,
				CatchHandlerPC:    -1,
			}
		}

		c.ctx.Emitter.EmitJump(endLabel)

		return ProtectedRecoveryRegion{
			Result:            out,
			StartCatch:        startCatch,
			EndCatchExclusive: endCatchExclusive,
			CatchHandlerPC:    -1,
		}
	}
}

func (c *RecoveryCompiler) compileWithErrorReturn(
	plan core.RecoveryPlan,
	buildProtected func(recoveryLabel, endLabel core.Label) ProtectedRecoveryRegion,
) bytecode.Operand {
	recoveryLabel := c.ctx.Emitter.NewLabel("recovery", "error", "handle")
	endLabel := c.ctx.Emitter.NewLabel("recovery", "error", "end")
	region := buildProtected(recoveryLabel, endLabel)

	if region.Result == bytecode.NoopOperand || region.EndCatchExclusive <= region.StartCatch {
		return region.Result
	}

	handlerPC := c.ctx.Emitter.Size()
	c.ctx.Emitter.MarkLabel(recoveryLabel)

	fallback := c.front.Expressions.Compile(plan.OnError.Expr)
	c.front.TypeFacts.EmitMoveAuto(c.EnsureRegister(region.Result), c.EnsureRegister(fallback))
	c.ctx.Emitter.MarkLabel(endLabel)

	c.pushProtectedCatch(region, handlerPC)

	return c.WidenResultType(region.Result, plan)
}

func (c *RecoveryCompiler) compileWithErrorRetry(
	plan core.RecoveryPlan,
	buildProtected func(recoveryLabel, endLabel core.Label) ProtectedRecoveryRegion,
	compileFinalAttempt func() bytecode.Operand,
) bytecode.Operand {
	handler := plan.OnError
	if !recoveryHandlerRetries(handler) || handler.Retry == nil {
		if compileFinalAttempt != nil {
			return compileFinalAttempt()
		}

		return bytecode.NoopOperand
	}

	retry := handler.Retry
	if retry.FinalActionKind != core.RecoveryActionReturn && retry.Count <= 0 {
		if compileFinalAttempt != nil {
			return compileFinalAttempt()
		}

		return bytecode.NoopOperand
	}

	resultReg := bytecode.NoopOperand
	zeroReg := c.front.TypeFacts.LoadConstant(runtime.ZeroInt)
	retriesRemainingReg := c.front.TypeFacts.LoadConstant(runtime.NewInt(retry.Count))

	state := c.initRetryDelayState(retry)
	retryStart := c.ctx.Emitter.NewLabel("recovery", "retry", "start")
	recoveryLabel := c.ctx.Emitter.NewLabel("recovery", "retry", "handle")
	endLabel := c.ctx.Emitter.NewLabel("recovery", "retry", "end")
	var finalAttemptLabel core.Label

	if retry.FinalActionKind != core.RecoveryActionReturn {
		finalAttemptLabel = c.ctx.Emitter.NewLabel("recovery", "retry", "final")
	}

	c.ctx.Emitter.MarkLabel(retryStart)
	region := buildProtected(recoveryLabel, endLabel)

	if region.Result == bytecode.NoopOperand || region.EndCatchExclusive <= region.StartCatch {
		return region.Result
	}

	resultReg = c.EnsureRegister(region.Result)

	handlerPC := c.ctx.Emitter.Size()
	c.ctx.Emitter.MarkLabel(recoveryLabel)

	retriesAvailableReg := c.ctx.Registers.Allocate()
	c.ctx.Emitter.EmitGt(retriesAvailableReg, retriesRemainingReg, zeroReg)

	onExhausted := c.ctx.Emitter.NewLabel("recovery", "retry", "exhausted")
	c.ctx.Emitter.EmitJumpIfFalse(retriesAvailableReg, onExhausted)
	c.ctx.Emitter.EmitA(bytecode.OpDecr, retriesRemainingReg)
	c.EmitRetryDelay(retry, state)

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
		fallback := c.front.Expressions.Compile(retry.FinalExpr)
		c.front.TypeFacts.EmitMoveAuto(resultReg, c.EnsureRegister(fallback))
		c.ctx.Emitter.EmitJump(endLabel)
	} else {
		c.ctx.Emitter.EmitJump(finalAttemptLabel)
	}

	c.pushProtectedCatch(region, handlerPC)

	if retry.FinalActionKind != core.RecoveryActionReturn {
		c.ctx.Emitter.MarkLabel(finalAttemptLabel)

		if compileFinalAttempt != nil {
			finalOut := c.EnsureRegister(compileFinalAttempt())
			if finalOut != bytecode.NoopOperand && finalOut != resultReg {
				c.front.TypeFacts.EmitMoveAuto(resultReg, finalOut)
			}
		}
	}

	c.ctx.Emitter.MarkLabel(endLabel)

	return c.WidenResultType(resultReg, core.RecoveryPlan{OnError: handler})
}

func (c *RecoveryCompiler) pushProtectedCatch(region ProtectedRecoveryRegion, handlerPC int) {
	if region.EndCatchExclusive <= region.StartCatch {
		return
	}

	catchHandlerPC := handlerPC
	if region.CatchHandlerPC >= 0 {
		catchHandlerPC = region.CatchHandlerPC
	}

	c.ctx.CatchTable.Push(region.StartCatch, region.EndCatchExclusive-1, catchHandlerPC)
}
