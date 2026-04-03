package internal

import (
	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/compiler/internal/core"
	"github.com/MontFerret/ferret/v2/pkg/parser/fql"
)

// WaitCompiler handles the compilation of WAITFOR expressions in FQL queries.
// It transforms wait operations into VM instructions for event streaming and polling.
type WaitCompiler struct {
	ctx   *CompilationSession
	front *CompilationFrontend
}

// NewWaitCompiler creates a new instance of WaitCompiler with the given compiler context.
func NewWaitCompiler(ctx *CompilationSession) *WaitCompiler {
	return &WaitCompiler{
		ctx: ctx,
	}
}

// Compile processes a WAITFOR expression from the FQL AST and generates the appropriate VM instructions.
func (c *WaitCompiler) Compile(ctx fql.IWaitForExpressionContext) bytecode.Operand {
	return c.compileWithOuterRecovery(ctx, core.RecoveryPlan{})
}

func (c *WaitCompiler) compileWithOuterRecovery(ctx fql.IWaitForExpressionContext, outerPlan core.RecoveryPlan) bytecode.Operand {
	if ctx == nil {
		return bytecode.NoopOperand
	}

	plan := c.front.Recovery.NormalizePlan(c.front.Recovery.MergePlans(c.front.Recovery.CollectPlan(ctx, core.RecoveryPlanOptions{
		AllowTimeout: true,
		HasTimeout:   waitForHasExplicitTimeoutClause(ctx),
	}), outerPlan))
	c.ctx.Symbols.EnterScope()
	defer c.ctx.Symbols.ExitScope()

	if ev := ctx.WaitForEventExpression(); ev != nil {
		return c.compileEventWithPlan(ev, plan)
	}

	if pred := ctx.WaitForPredicateExpression(); pred != nil {
		return c.compilePredicateWithPlan(pred, plan)
	}

	return bytecode.NoopOperand
}

func waitForHasExplicitTimeoutClause(ctx fql.IWaitForExpressionContext) bool {
	if ctx == nil {
		return false
	}

	if ev := ctx.WaitForEventExpression(); ev != nil && ev.TimeoutClause() != nil {
		return true
	}

	if pred := ctx.WaitForPredicateExpression(); pred != nil && pred.TimeoutClause() != nil {
		return true
	}

	return false
}

func (c *WaitCompiler) compileEventWithPlan(ctx fql.IWaitForEventExpressionContext, plan core.RecoveryPlan) bytecode.Operand {
	plan = c.front.Recovery.NormalizePlan(plan)

	if plan.OnTimeout == nil && (plan.OnError == nil || plan.OnError.ActionKind == core.RecoveryActionFail) {
		return c.compileEvent(ctx)
	}

	if plan.OnError == nil || plan.OnError.ActionKind == core.RecoveryActionFail {
		return c.front.Recovery.WidenResultType(c.compileEventWithTimeoutRecovery(ctx, plan), plan)
	}

	return c.front.Recovery.WidenResultType(c.front.Recovery.CompileWithProtectedRecovery(ProtectedRecoverySpec{
		Plan: plan,
		BuildProtected: func(recoveryLabel, endLabel core.Label) ProtectedRecoveryRegion {
			return c.buildProtectedEventRecovery(ctx, plan, recoveryLabel, endLabel)
		},
		CompileFinalAttempt: func() bytecode.Operand {
			return c.compileEventWithPlan(ctx, core.RecoveryPlan{OnTimeout: plan.OnTimeout})
		},
	}), plan)
}

func (c *WaitCompiler) compilePredicateWithPlan(ctx fql.IWaitForPredicateExpressionContext, plan core.RecoveryPlan) bytecode.Operand {
	plan = c.front.Recovery.NormalizePlan(plan)

	if plan.OnTimeout == nil {
		errorPlan := core.RecoveryPlan{OnError: plan.OnError}

		return c.front.Recovery.CompileWithRecoveryPlan(errorPlan, core.CatchJumpModeNone, func() bytecode.Operand {
			return c.compilePredicate(ctx)
		})
	}

	if plan.OnError == nil || plan.OnError.ActionKind == core.RecoveryActionFail {
		return c.front.Recovery.WidenResultType(c.compilePredicateWithTimeoutRecovery(ctx, plan), plan)
	}

	return c.front.Recovery.WidenResultType(c.front.Recovery.CompileWithProtectedRecovery(ProtectedRecoverySpec{
		Plan: plan,
		BuildProtected: func(recoveryLabel, endLabel core.Label) ProtectedRecoveryRegion {
			return c.buildProtectedPredicateRecovery(ctx, plan, recoveryLabel, endLabel)
		},
		CompileFinalAttempt: func() bytecode.Operand {
			return c.compilePredicateWithPlan(ctx, core.RecoveryPlan{OnTimeout: plan.OnTimeout})
		},
	}), plan)
}

func (c *WaitCompiler) buildProtectedEventRecovery(
	ctx fql.IWaitForEventExpressionContext,
	plan core.RecoveryPlan,
	recoveryLabel, endLabel core.Label,
) ProtectedRecoveryRegion {
	streamReg := c.ctx.Registers.Allocate()
	resultReg := c.ctx.Registers.Allocate()
	timeoutStateReg := c.ctx.Registers.Allocate()
	errorStateReg := c.ctx.Registers.Allocate()

	c.ctx.Emitter.EmitLoadNone(resultReg)
	c.ctx.Emitter.EmitBoolean(timeoutStateReg, false)
	c.ctx.Emitter.EmitBoolean(errorStateReg, false)

	startCatch := c.ctx.Emitter.Size()
	state := c.buildWaitEventState(ctx)

	c.emitWaitEventStreamSetup(state, streamReg)

	start := c.ctx.Emitter.NewLabel()
	iterationDone := c.ctx.Emitter.NewLabel()
	cleanup := c.ctx.Emitter.NewLabel()
	timeoutHandler := c.ctx.Emitter.NewLabel("waitfor", "event", "timeout")
	routeRecovery := c.ctx.Emitter.NewLabel("waitfor", "event", "recover")

	c.ctx.Emitter.MarkLabel(start)
	c.emitWaitEventIteration(ctx, state, streamReg, timeoutStateReg, start, iterationDone)

	c.ctx.Emitter.EmitJump(cleanup)
	c.ctx.Emitter.MarkLabel(iterationDone)
	c.ctx.Emitter.EmitJump(cleanup)

	c.ctx.Emitter.MarkLabel(cleanup)
	c.emitWaitEventCleanup(state, streamReg)

	endCatchExclusive := c.ctx.Emitter.Size()

	c.ctx.Emitter.EmitJumpIfTrue(timeoutStateReg, timeoutHandler)
	c.ctx.Emitter.EmitJumpIfTrue(errorStateReg, routeRecovery)
	c.ctx.Emitter.EmitJump(endLabel)

	errorPreludePC := c.ctx.Emitter.Size()
	c.ctx.Emitter.EmitBoolean(errorStateReg, true)
	c.ctx.Emitter.EmitBoolean(timeoutStateReg, false)
	c.ctx.Emitter.EmitJump(cleanup)

	c.ctx.Emitter.MarkLabel(timeoutHandler)
	switch {
	case plan.OnTimeout != nil && plan.OnTimeout.ActionKind == core.RecoveryActionReturn:
		fallback := c.front.Expressions.Compile(plan.OnTimeout.Expr)
		c.front.TypeFacts.EmitMoveAuto(resultReg, c.front.Recovery.EnsureRegister(fallback))
		c.ctx.Emitter.EmitJump(endLabel)
	default:
		c.ctx.Emitter.Emit(bytecode.OpFailTimeout)
	}

	c.ctx.Emitter.MarkLabel(routeRecovery)
	c.ctx.Emitter.EmitBoolean(errorStateReg, false)
	c.ctx.Emitter.EmitBoolean(timeoutStateReg, false)
	c.ctx.Emitter.EmitJump(recoveryLabel)

	return ProtectedRecoveryRegion{
		Result:            resultReg,
		StartCatch:        startCatch,
		EndCatchExclusive: endCatchExclusive,
		CatchHandlerPC:    errorPreludePC,
	}
}

func (c *WaitCompiler) buildProtectedPredicateRecovery(
	ctx fql.IWaitForPredicateExpressionContext,
	plan core.RecoveryPlan,
	_ core.Label,
	endLabel core.Label,
) ProtectedRecoveryRegion {
	predicate := ctx.WaitForPredicate()
	if predicate == nil {
		return ProtectedRecoveryRegion{Result: bytecode.NoopOperand}
	}

	predExpr := predicate.Expression()
	if predExpr == nil {
		return ProtectedRecoveryRegion{Result: bytecode.NoopOperand}
	}

	config := c.buildWaitPredicateConfig(ctx, predicate, predExpr)
	c.normalizeWaitPredicateConfig(&config)
	state := c.initWaitPredicatePollState(config)

	start := c.ctx.Emitter.NewLabel()
	success := c.ctx.Emitter.NewLabel()
	protectedTimeout := c.ctx.Emitter.NewLabel()
	timeoutHandler := c.ctx.Emitter.NewLabel("waitfor", "predicate", "timeout")

	startCatch := c.ctx.Emitter.Size()

	c.ctx.Emitter.MarkLabel(start)
	valueReg := c.emitWaitPredicatePollIteration(config, state, start, success, protectedTimeout)

	c.ctx.Emitter.MarkLabel(success)
	c.emitWaitSuccessResult(config.mode, state.resultReg, valueReg)
	c.ctx.Emitter.EmitJump(endLabel)

	c.ctx.Emitter.MarkLabel(protectedTimeout)
	c.ctx.Emitter.EmitJump(timeoutHandler)

	endCatchExclusive := c.ctx.Emitter.Size()

	c.ctx.Emitter.MarkLabel(timeoutHandler)
	switch {
	case plan.OnTimeout != nil && plan.OnTimeout.ActionKind == core.RecoveryActionReturn:
		fallback := c.front.Expressions.Compile(plan.OnTimeout.Expr)
		c.front.TypeFacts.EmitMoveAuto(state.resultReg, c.front.Recovery.EnsureRegister(fallback))
		c.ctx.Emitter.EmitJump(endLabel)
	case plan.OnTimeout != nil && plan.OnTimeout.ActionKind == core.RecoveryActionFail:
		c.ctx.Emitter.Emit(bytecode.OpFailTimeout)
	default:
		c.emitWaitTimeoutResult(config.mode, state.resultReg)
		c.ctx.Emitter.EmitJump(endLabel)
	}

	return ProtectedRecoveryRegion{
		Result:            state.resultReg,
		StartCatch:        startCatch,
		EndCatchExclusive: endCatchExclusive,
		CatchHandlerPC:    -1,
	}
}
