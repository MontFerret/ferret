package internal

import (
	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/compiler/internal/core"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

type (
	RecoveryCompiler struct {
		ctx      *CompilationSession
		exprs    *ExprCompiler
		literals *LiteralCompiler
		facts    *TypeFacts
	}

	// OperationRecoverySpec describes one operation's recovery-owned execution surface.
	// Operation compilers provide execution shapes while RecoveryCompiler owns plan
	// collection, plan merging, widening, retry/fallback orchestration, and catch setup.
	OperationRecoverySpec struct {
		Owner                core.RecoveryTailOwner
		OuterPlan            core.RecoveryPlan
		CompilePlain         func() bytecode.Operand
		CompileSuppressed    func() bytecode.Operand
		CompileTimeoutAware  func(timeoutLabel, endLabel core.Label) bytecode.Operand
		BuildProtected       func(recoveryLabel, timeoutLabel, endLabel core.Label) ProtectedRecoveryRegion
		ShouldBuildProtected func(plan core.RecoveryPlan) bool
		CompileFinalAttempt  func(plan core.RecoveryPlan) bytecode.Operand
		JumpMode             core.CatchJumpMode
		Options              core.RecoveryPlanOptions
	}

	// ProtectedRecoveryRegion describes the guarded portion of an operation.
	// CatchHandlerPC == -1 means the catch table should jump directly to the
	// generic recovery handler emitted by RecoveryCompiler.
	ProtectedRecoveryRegion struct {
		Result            bytecode.Operand
		StartCatch        int
		EndCatchExclusive int
		CatchHandlerPC    int
		HasTimeout        bool
	}
)

func NewRecoveryCompiler(ctx *CompilationSession) *RecoveryCompiler {
	return &RecoveryCompiler{ctx: ctx}
}

func (c *RecoveryCompiler) bind(exprs *ExprCompiler, literals *LiteralCompiler, facts *TypeFacts) {
	if c == nil {
		return
	}

	c.exprs = exprs
	c.literals = literals
	c.facts = facts
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

func (c *RecoveryCompiler) CompileOperation(spec OperationRecoverySpec) bytecode.Operand {
	plan := c.NormalizePlan(c.MergePlans(c.CollectPlan(spec.Owner, spec.Options), spec.OuterPlan))

	if plan.OnError != nil && plan.OnError.ActionKind != core.RecoveryActionFail {
		if spec.shouldBuildProtected(plan) {
			return c.compileProtectedOperation(plan, spec)
		}

		return c.compileDirectOperation(plan, spec)
	}

	return c.WidenResultType(c.compileOperationFinalAttempt(spec, plan), plan)
}

func (c *RecoveryCompiler) compileDirectOperation(plan core.RecoveryPlan, spec OperationRecoverySpec) bytecode.Operand {
	if spec.CompilePlain == nil {
		return bytecode.NoopOperand
	}

	if hasErrorReturnNoneHandler(plan) {
		compile := spec.CompileSuppressed
		if compile == nil {
			compile = spec.CompilePlain
		}

		out := c.CompileWithErrorPolicy(core.ErrorPolicySuppress, spec.JumpMode, compile)
		return c.WidenResultType(out, plan)
	}

	buildProtected := c.directProtectedRegionBuilder(spec.CompilePlain)

	switch plan.OnError.ActionKind {
	case core.RecoveryActionReturn:
		return c.compileOperationWithErrorReturn(plan, buildProtected)
	case core.RecoveryActionRetry:
		return c.compileOperationWithErrorRetry(plan, buildProtected, func() bytecode.Operand {
			return c.compileOperationFinalAttempt(spec, core.RecoveryPlan{OnTimeout: plan.OnTimeout})
		})
	default:
		return c.WidenResultType(c.compileOperationFinalAttempt(spec, plan), plan)
	}
}

func (c *RecoveryCompiler) compileProtectedOperation(plan core.RecoveryPlan, spec OperationRecoverySpec) bytecode.Operand {
	if spec.BuildProtected == nil {
		return c.WidenResultType(c.compileOperationFinalAttempt(spec, plan), plan)
	}

	switch plan.OnError.ActionKind {
	case core.RecoveryActionReturn:
		return c.compileOperationWithErrorReturn(plan, spec.BuildProtected)
	case core.RecoveryActionRetry:
		return c.compileOperationWithErrorRetry(plan, spec.BuildProtected, func() bytecode.Operand {
			return c.compileOperationFinalAttempt(spec, core.RecoveryPlan{OnTimeout: plan.OnTimeout})
		})
	default:
		return c.WidenResultType(c.compileOperationFinalAttempt(spec, plan), plan)
	}
}

func (c *RecoveryCompiler) compileOperationFinalAttempt(spec OperationRecoverySpec, plan core.RecoveryPlan) bytecode.Operand {
	if spec.CompileFinalAttempt != nil {
		return spec.CompileFinalAttempt(plan)
	}

	if plan.OnTimeout != nil && spec.CompileTimeoutAware != nil {
		return c.compileWithTimeoutRecovery(plan, spec.CompileTimeoutAware)
	}

	if spec.CompilePlain != nil {
		return spec.CompilePlain()
	}

	return bytecode.NoopOperand
}

func (c *RecoveryCompiler) compileWithTimeoutRecovery(
	plan core.RecoveryPlan,
	compile func(timeoutLabel, endLabel core.Label) bytecode.Operand,
) bytecode.Operand {
	if compile == nil {
		return bytecode.NoopOperand
	}

	timeoutLabel := c.ctx.Emitter.NewLabel("recovery", "timeout", "handle")
	endLabel := c.ctx.Emitter.NewLabel("recovery", "timeout", "end")
	out := ensureOperandRegister(c.ctx, c.facts, compile(timeoutLabel, endLabel))

	if out == bytecode.NoopOperand {
		return out
	}

	c.emitTimeoutHandler(out, plan, timeoutLabel, endLabel)
	c.ctx.Emitter.MarkLabel(endLabel)

	return c.WidenResultType(out, plan)
}

func (c *RecoveryCompiler) directProtectedRegionBuilder(compile func() bytecode.Operand) func(recoveryLabel, timeoutLabel, endLabel core.Label) ProtectedRecoveryRegion {
	return func(_ core.Label, _ core.Label, endLabel core.Label) ProtectedRecoveryRegion {
		startCatch := c.ctx.Emitter.Size()
		out := ensureOperandRegister(c.ctx, c.facts, compile())
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
			HasTimeout:        false,
		}
	}
}

func (c *RecoveryCompiler) compileOperationWithErrorReturn(
	plan core.RecoveryPlan,
	buildProtected func(recoveryLabel, timeoutLabel, endLabel core.Label) ProtectedRecoveryRegion,
) bytecode.Operand {
	recoveryLabel := c.ctx.Emitter.NewLabel("recovery", "error", "handle")
	var timeoutLabel core.Label
	if plan.OnTimeout != nil {
		timeoutLabel = c.ctx.Emitter.NewLabel("recovery", "timeout", "handle")
	}
	endLabel := c.ctx.Emitter.NewLabel("recovery", "error", "end")
	region := buildProtected(recoveryLabel, timeoutLabel, endLabel)

	if region.Result == bytecode.NoopOperand || region.EndCatchExclusive <= region.StartCatch {
		return region.Result
	}

	handlerPC := c.ctx.Emitter.Size()
	c.ctx.Emitter.MarkLabel(recoveryLabel)

	fallback := c.exprs.Compile(plan.OnError.Expr)
	c.facts.EmitMoveAuto(ensureOperandRegister(c.ctx, c.facts, region.Result), ensureOperandRegister(c.ctx, c.facts, fallback))
	c.ctx.Emitter.EmitJump(endLabel)

	if region.HasTimeout {
		c.emitTimeoutHandler(region.Result, plan, timeoutLabel, endLabel)
	}

	c.ctx.Emitter.MarkLabel(endLabel)

	c.pushProtectedCatch(region, handlerPC)

	return c.WidenResultType(region.Result, plan)
}

func (c *RecoveryCompiler) compileOperationWithErrorRetry(
	plan core.RecoveryPlan,
	buildProtected func(recoveryLabel, timeoutLabel, endLabel core.Label) ProtectedRecoveryRegion,
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
	zeroReg := c.facts.LoadConstant(runtime.ZeroInt)
	retriesRemainingReg := c.facts.LoadConstant(runtime.NewInt(retry.Count))

	state := c.initRetryDelayState(retry)
	retryStart := c.ctx.Emitter.NewLabel("recovery", "retry", "start")
	recoveryLabel := c.ctx.Emitter.NewLabel("recovery", "retry", "handle")
	var timeoutLabel core.Label
	if plan.OnTimeout != nil {
		timeoutLabel = c.ctx.Emitter.NewLabel("recovery", "timeout", "handle")
	}
	endLabel := c.ctx.Emitter.NewLabel("recovery", "retry", "end")
	var finalAttemptLabel core.Label

	if retry.FinalActionKind != core.RecoveryActionReturn {
		finalAttemptLabel = c.ctx.Emitter.NewLabel("recovery", "retry", "final")
	}

	c.ctx.Emitter.MarkLabel(retryStart)
	region := buildProtected(recoveryLabel, timeoutLabel, endLabel)

	if region.Result == bytecode.NoopOperand || region.EndCatchExclusive <= region.StartCatch {
		return region.Result
	}

	resultReg = ensureOperandRegister(c.ctx, c.facts, region.Result)

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
		fallback := c.exprs.Compile(retry.FinalExpr)
		c.facts.EmitMoveAuto(resultReg, ensureOperandRegister(c.ctx, c.facts, fallback))
		c.ctx.Emitter.EmitJump(endLabel)
	} else {
		c.ctx.Emitter.EmitJump(finalAttemptLabel)
	}

	if region.HasTimeout {
		c.emitTimeoutHandler(resultReg, plan, timeoutLabel, endLabel)
	}

	if retry.FinalActionKind != core.RecoveryActionReturn {
		c.ctx.Emitter.MarkLabel(finalAttemptLabel)

		if compileFinalAttempt != nil {
			finalOut := ensureOperandRegister(c.ctx, c.facts, compileFinalAttempt())
			if finalOut != bytecode.NoopOperand && finalOut != resultReg {
				c.facts.EmitMoveAuto(resultReg, finalOut)
			}
		}
	}

	c.pushProtectedCatch(region, handlerPC)
	c.ctx.Emitter.MarkLabel(endLabel)

	return c.WidenResultType(resultReg, plan)
}

func (c *RecoveryCompiler) emitTimeoutHandler(
	result bytecode.Operand,
	plan core.RecoveryPlan,
	timeoutLabel, endLabel core.Label,
) {
	c.ctx.Emitter.MarkLabel(timeoutLabel)

	switch {
	case plan.OnTimeout != nil && plan.OnTimeout.ActionKind == core.RecoveryActionReturn:
		fallback := c.exprs.Compile(plan.OnTimeout.Expr)
		c.facts.EmitMoveAuto(ensureOperandRegister(c.ctx, c.facts, result), ensureOperandRegister(c.ctx, c.facts, fallback))
		c.ctx.Emitter.EmitJump(endLabel)
	default:
		c.ctx.Emitter.Emit(bytecode.OpFailTimeout)
	}
}

func (spec OperationRecoverySpec) shouldBuildProtected(plan core.RecoveryPlan) bool {
	if spec.BuildProtected == nil {
		return false
	}

	if spec.ShouldBuildProtected != nil {
		return spec.ShouldBuildProtected(plan)
	}

	return true
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
