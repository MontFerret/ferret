package internal

import (
	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/compiler/internal/core"
	parserd "github.com/MontFerret/ferret/v2/pkg/parser/diagnostics"
	"github.com/MontFerret/ferret/v2/pkg/parser/fql"
)

type (
	waitForSynchronization int

	waitPredicateGroupArm struct {
		predExpr  fql.IExpressionContext
		whenExprs []fql.IExpressionContext
	}

	waitPredicateGroupCompileConfig struct {
		arms []waitPredicateGroupArm
		waitPredicateScheduleConfig
		mode            waitForPredicateMode
		synchronization waitForSynchronization
	}
)

const (
	waitForSynchronizationAny waitForSynchronization = iota
	waitForSynchronizationAll
)

func (c *WaitCompiler) compilePredicateGroup(ctx fql.IWaitForPredicateGroupExpressionContext) bytecode.Operand {
	config, ok := c.prepareWaitPredicateGroupConfig(ctx)
	if !ok {
		return bytecode.NoopOperand
	}

	state := c.initWaitPredicatePollState(config.mode, config.waitPredicateScheduleConfig)
	c.emitWaitPredicateGroupPollLoop(config, state)

	return state.resultReg
}

func (c *WaitCompiler) compilePredicateGroupWithTimeoutRecovery(
	ctx fql.IWaitForPredicateGroupExpressionContext,
	timeoutLabel, endLabel core.Label,
) bytecode.Operand {
	config, ok := c.prepareWaitPredicateGroupConfig(ctx)
	if !ok {
		return bytecode.NoopOperand
	}

	state := c.initWaitPredicatePollState(config.mode, config.waitPredicateScheduleConfig)
	c.emitWaitPredicateGroupPollLoopWithRecovery(config, state, timeoutLabel, endLabel)

	return state.resultReg
}

func (c *WaitCompiler) prepareWaitPredicateGroupConfig(
	ctx fql.IWaitForPredicateGroupExpressionContext,
) (waitPredicateGroupCompileConfig, bool) {
	if ctx == nil {
		return waitPredicateGroupCompileConfig{}, false
	}

	entries := ctx.AllWaitForPredicateGroupEntry()
	if len(entries) == 0 {
		return waitPredicateGroupCompileConfig{}, false
	}

	arms := make([]waitPredicateGroupArm, 0, len(entries))
	for _, entry := range entries {
		if entry == nil || entry.Expression() == nil {
			return waitPredicateGroupCompileConfig{}, false
		}

		if legacy := legacyWaitForOrThrowNode(entry.Expression()); legacy != nil {
			c.ctx.Program.Errors.Add(c.ctx.Program.Errors.Create(parserd.SyntaxError, legacy, "Unexpected THROW after OR in WAITFOR predicate"))
			return waitPredicateGroupCompileConfig{}, false
		}

		arms = append(arms, waitPredicateGroupArm{
			predExpr:  entry.Expression(),
			whenExprs: waitPredicateWhenExpressions(entry.AllWaitForPredicateWhenClause()),
		})
	}

	schedule, ok := c.buildWaitPredicateScheduleConfig(
		ctx.TimeoutClause(),
		ctx.EveryClause(),
		ctx.BackoffClause(),
		ctx.JitterClause(),
	)
	if !ok {
		return waitPredicateGroupCompileConfig{}, false
	}

	c.normalizeWaitPredicateScheduleConfig(&schedule)

	modeCtx := ctx.WaitForPredicateGroupMode()
	mode := resolveWaitPredicateMode(
		modeCtx != nil && modeCtx.Value() != nil,
		modeCtx != nil && modeCtx.Exists() != nil,
		modeCtx != nil && modeCtx.Not() != nil,
	)

	return waitPredicateGroupCompileConfig{
		waitPredicateScheduleConfig: schedule,
		arms:                        arms,
		mode:                        mode,
		synchronization:             resolveWaitForSynchronization(ctx.WaitForSynchronization()),
	}, true
}

func (c *WaitCompiler) emitWaitPredicateGroupPollLoop(
	config waitPredicateGroupCompileConfig,
	state waitPredicatePollState,
) {
	start := c.ctx.Program.Emitter.NewLabel()
	timeoutLabel := c.ctx.Program.Emitter.NewLabel()
	end := c.ctx.Program.Emitter.NewLabel()

	c.emitWaitPredicateGroupPollCycle(config, state, start, timeoutLabel, end)

	c.ctx.Program.Emitter.MarkLabel(timeoutLabel)
	c.emitWaitTimeoutResult(config.mode, state.resultReg)
	c.ctx.Program.Emitter.MarkLabel(end)
}

func (c *WaitCompiler) emitWaitPredicateGroupPollLoopWithRecovery(
	config waitPredicateGroupCompileConfig,
	state waitPredicatePollState,
	timeoutLabel, endLabel core.Label,
) {
	start := c.ctx.Program.Emitter.NewLabel()
	c.emitWaitPredicateGroupPollCycle(config, state, start, timeoutLabel, endLabel)
}

func (c *WaitCompiler) emitWaitPredicateGroupPollCycle(
	config waitPredicateGroupCompileConfig,
	state waitPredicatePollState,
	startLabel, timeoutLabel, endLabel core.Label,
) {
	c.ctx.Program.Emitter.MarkLabel(startLabel)

	if config.synchronization == waitForSynchronizationAny {
		c.emitWaitPredicateAnyCycle(config, state.resultReg, endLabel)
	} else {
		c.emitWaitPredicateAllCycle(config, state.resultReg, endLabel)
	}

	c.emitWaitPredicatePollRetry(config.waitPredicateScheduleConfig, state, startLabel, timeoutLabel)
}

func (c *WaitCompiler) emitWaitPredicateAnyCycle(
	config waitPredicateGroupCompileConfig,
	resultReg bytecode.Operand,
	endLabel core.Label,
) {
	for _, arm := range config.arms {
		nextArm := c.ctx.Program.Emitter.NewLabel()
		valueReg := c.exprs.Compile(arm.predExpr)
		c.emitWaitPredicateGroupArmGuard(config.mode, arm.whenExprs, valueReg, nextArm)
		c.emitWaitSuccessResult(config.mode, resultReg, valueReg)
		c.ctx.Program.Emitter.EmitJump(endLabel)
		c.ctx.Program.Emitter.MarkLabel(nextArm)
	}
}

func (c *WaitCompiler) emitWaitPredicateAllCycle(
	config waitPredicateGroupCompileConfig,
	resultReg bytecode.Operand,
	endLabel core.Label,
) {
	cycleFailed := c.ctx.Program.Emitter.NewLabel()
	values := make([]bytecode.Operand, 0, len(config.arms))

	for _, arm := range config.arms {
		valueReg := c.exprs.Compile(arm.predExpr)
		values = append(values, valueReg)
		c.emitWaitPredicateGroupArmGuard(config.mode, arm.whenExprs, valueReg, cycleFailed)
	}

	if config.mode == waitForPredicateModeValue {
		c.ctx.Program.Emitter.EmitArray(resultReg, len(values))
		for _, valueReg := range values {
			c.ctx.Program.Emitter.EmitArrayPush(resultReg, valueReg)
		}
	} else {
		c.ctx.Program.Emitter.EmitBoolean(resultReg, true)
	}
	c.ctx.Program.Emitter.EmitJump(endLabel)
	c.ctx.Program.Emitter.MarkLabel(cycleFailed)
}

func (c *WaitCompiler) emitWaitPredicateGroupArmGuard(
	mode waitForPredicateMode,
	whenExprs []fql.IExpressionContext,
	valueReg bytecode.Operand,
	failureLabel core.Label,
) {
	if mode == waitForPredicateModeValue {
		c.ctx.Program.Emitter.EmitJumpIfNone(valueReg, failureLabel)
	} else {
		conditionReg := c.emitWaitPredicateCondition(mode, valueReg)
		c.ctx.Program.Emitter.EmitJumpIfFalse(conditionReg, failureLabel)
	}

	c.emitWaitPredicateWhenConditions(whenExprs, valueReg, failureLabel)
}
