package internal

import (
	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/compiler/internal/core"
	"github.com/MontFerret/ferret/v2/pkg/parser/fql"
)

// WaitCompiler handles the compilation of WAITFOR expressions in FQL queries.
// It transforms wait operations into VM instructions for event streaming and polling.
type WaitCompiler struct {
	ctx      *CompilationSession
	exprs    *ExprCompiler
	literals *LiteralCompiler
	recovery *RecoveryCompiler
	facts    *TypeFacts
}

// NewWaitCompiler creates a new instance of WaitCompiler with the given compiler context.
func NewWaitCompiler(ctx *CompilationSession) *WaitCompiler {
	return &WaitCompiler{
		ctx: ctx,
	}
}

func (c *WaitCompiler) bind(exprs *ExprCompiler, literals *LiteralCompiler, recovery *RecoveryCompiler, facts *TypeFacts) {
	if c == nil {
		return
	}

	c.exprs = exprs
	c.literals = literals
	c.recovery = recovery
	c.facts = facts
}

// Compile processes a WAITFOR expression from the FQL AST and generates the appropriate VM instructions.
func (c *WaitCompiler) Compile(ctx fql.IWaitForExpressionContext) bytecode.Operand {
	if ctx == nil {
		return bytecode.NoopOperand
	}

	c.ctx.Symbols.EnterScope()
	defer c.ctx.Symbols.ExitScope()

	return c.recovery.CompileOperation(c.newWaitOperationRecoverySpec(ctx, core.RecoveryPlan{}))
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

func (c *WaitCompiler) newWaitOperationRecoverySpec(ctx fql.IWaitForExpressionContext, outerPlan core.RecoveryPlan) OperationRecoverySpec {
	spec := OperationRecoverySpec{
		Owner: ctx,
		Options: core.RecoveryPlanOptions{
			AllowTimeout: true,
			HasTimeout:   waitForHasExplicitTimeoutClause(ctx),
		},
		OuterPlan: outerPlan,
	}

	if ctx == nil {
		return spec
	}

	if ev := ctx.WaitForEventExpression(); ev != nil {
		spec.CompilePlain = func() bytecode.Operand {
			return c.compileEvent(ev)
		}
		spec.BuildProtected = func(recoveryLabel, timeoutLabel, endLabel core.Label) ProtectedRecoveryRegion {
			return c.buildProtectedEventRecovery(ev, recoveryLabel, timeoutLabel, endLabel)
		}

		if ev.TimeoutClause() != nil {
			spec.CompileTimeoutAware = func(timeoutLabel, endLabel core.Label) bytecode.Operand {
				return c.compileEventWithTimeoutRecovery(ev, timeoutLabel, endLabel)
			}
		}

		return spec
	}

	if pred := ctx.WaitForPredicateExpression(); pred != nil {
		spec.CompilePlain = func() bytecode.Operand {
			return c.compilePredicate(pred)
		}

		if pred.TimeoutClause() != nil {
			spec.CompileTimeoutAware = func(timeoutLabel, endLabel core.Label) bytecode.Operand {
				return c.compilePredicateWithTimeoutRecovery(pred, timeoutLabel, endLabel)
			}
			spec.BuildProtected = func(recoveryLabel, timeoutLabel, endLabel core.Label) ProtectedRecoveryRegion {
				return c.buildProtectedPredicateRecovery(pred, recoveryLabel, timeoutLabel, endLabel)
			}
			spec.ShouldBuildProtected = func(plan core.RecoveryPlan) bool {
				return plan.OnTimeout != nil
			}
		}
	}

	return spec
}

func (c *WaitCompiler) buildProtectedEventRecovery(
	ctx fql.IWaitForEventExpressionContext,
	recoveryLabel, timeoutLabel, endLabel core.Label,
) ProtectedRecoveryRegion {
	hasTimeout := ctx != nil && ctx.TimeoutClause() != nil
	streamReg := c.ctx.Registers.Allocate()
	resultReg := c.ctx.Registers.Allocate()
	errorStateReg := c.ctx.Registers.Allocate()
	timeoutStateReg := bytecode.NoopOperand

	if hasTimeout {
		timeoutStateReg = c.ctx.Registers.Allocate()
		c.ctx.Emitter.EmitBoolean(timeoutStateReg, false)
	}

	c.ctx.Emitter.EmitLoadNone(resultReg)
	c.ctx.Emitter.EmitBoolean(errorStateReg, false)

	startCatch := c.ctx.Emitter.Size()
	state, ok := c.buildWaitEventState(ctx)
	if !ok {
		return ProtectedRecoveryRegion{Result: bytecode.NoopOperand}
	}

	c.emitWaitEventStreamSetup(state, streamReg)

	start := c.ctx.Emitter.NewLabel()
	iterationDone := c.ctx.Emitter.NewLabel()
	cleanup := c.ctx.Emitter.NewLabel()
	routeRecovery := c.ctx.Emitter.NewLabel("waitfor", "event", "recover")

	c.ctx.Emitter.MarkLabel(start)
	c.emitWaitEventIteration(ctx, state, streamReg, timeoutStateReg, start, iterationDone)

	c.ctx.Emitter.EmitJump(cleanup)
	c.ctx.Emitter.MarkLabel(iterationDone)
	c.ctx.Emitter.EmitJump(cleanup)

	c.ctx.Emitter.MarkLabel(cleanup)
	c.emitWaitEventCleanup(state, streamReg)

	endCatchExclusive := c.ctx.Emitter.Size()

	if hasTimeout {
		c.ctx.Emitter.EmitJumpIfTrue(timeoutStateReg, timeoutLabel)
	}
	c.ctx.Emitter.EmitJumpIfTrue(errorStateReg, routeRecovery)
	c.ctx.Emitter.EmitJump(endLabel)

	errorPreludePC := c.ctx.Emitter.Size()
	c.ctx.Emitter.EmitBoolean(errorStateReg, true)
	if hasTimeout {
		c.ctx.Emitter.EmitBoolean(timeoutStateReg, false)
	}
	c.ctx.Emitter.EmitJump(cleanup)

	c.ctx.Emitter.MarkLabel(routeRecovery)
	c.ctx.Emitter.EmitBoolean(errorStateReg, false)
	if hasTimeout {
		c.ctx.Emitter.EmitBoolean(timeoutStateReg, false)
	}
	c.ctx.Emitter.EmitJump(recoveryLabel)

	return ProtectedRecoveryRegion{
		Result:            resultReg,
		StartCatch:        startCatch,
		EndCatchExclusive: endCatchExclusive,
		CatchHandlerPC:    errorPreludePC,
		HasTimeout:        hasTimeout,
	}
}

func (c *WaitCompiler) buildProtectedPredicateRecovery(
	ctx fql.IWaitForPredicateExpressionContext,
	_ core.Label,
	timeoutLabel core.Label,
	endLabel core.Label,
) ProtectedRecoveryRegion {
	config, ok := c.prepareWaitPredicateConfig(ctx)
	if !ok {
		return ProtectedRecoveryRegion{Result: bytecode.NoopOperand}
	}

	state := c.initWaitPredicatePollState(config)
	hasTimeout := config.timeoutReg != bytecode.NoopOperand

	start := c.ctx.Emitter.NewLabel()
	success := c.ctx.Emitter.NewLabel()
	protectedTimeout := core.Label{}
	if hasTimeout {
		protectedTimeout = c.ctx.Emitter.NewLabel()
	}

	startCatch := c.ctx.Emitter.Size()

	c.ctx.Emitter.MarkLabel(start)
	valueReg := c.emitWaitPredicatePollIteration(config, state, start, success, protectedTimeout)

	c.ctx.Emitter.MarkLabel(success)
	c.emitWaitSuccessResult(config.mode, state.resultReg, valueReg)
	c.ctx.Emitter.EmitJump(endLabel)

	endCatchExclusive := c.ctx.Emitter.Size()

	if hasTimeout {
		c.ctx.Emitter.MarkLabel(protectedTimeout)
		c.ctx.Emitter.EmitJump(timeoutLabel)
	}

	return ProtectedRecoveryRegion{
		Result:            state.resultReg,
		StartCatch:        startCatch,
		EndCatchExclusive: endCatchExclusive,
		CatchHandlerPC:    -1,
		HasTimeout:        hasTimeout,
	}
}
