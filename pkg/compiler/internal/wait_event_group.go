package internal

import (
	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/compiler/internal/core"
	"github.com/MontFerret/ferret/v2/pkg/parser/fql"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/source"
)

type (
	waitEventGroupArm struct {
		filters []fql.IEventFilterClauseContext
		span    source.Span
	}

	waitEventGroupCompileState struct {
		arms            []waitEventGroupArm
		matchedRegs     []bytecode.Operand
		valueRegs       []bytecode.Operand
		span            source.Span
		sourcesReg      bytecode.Operand
		namesReg        bytecode.Operand
		optionsReg      bytecode.Operand
		timeoutReg      bytecode.Operand
		synchronization waitForSynchronization
	}
)

func (c *WaitCompiler) compileEventGroup(ctx fql.IWaitForEventGroupExpressionContext) bytecode.Operand {
	if waitForEventGroupTriggerClause(ctx) != nil {
		return c.compileEventGroupWithTriggerCleanup(ctx)
	}

	state, ok := c.buildWaitEventGroupState(ctx)
	if !ok {
		return bytecode.NoopOperand
	}

	streamReg := c.ctx.Function.Registers.Allocate()
	resultReg := c.ctx.Function.Registers.Allocate()
	c.ctx.Program.Emitter.EmitLoadNone(resultReg)
	c.emitWaitEventGroupStreamSetup(state, streamReg)
	c.compileWaitEventGroupTrigger(ctx)

	start := c.ctx.Program.Emitter.NewLabel()
	done := c.ctx.Program.Emitter.NewLabel()
	c.ctx.Program.Emitter.MarkLabel(start)
	c.emitWaitEventGroupIteration(state, streamReg, resultReg, bytecode.NoopOperand, start, done)
	c.ctx.Program.Emitter.MarkLabel(done)
	c.emitWaitEventGroupCleanup(state, streamReg)

	return resultReg
}

func (c *WaitCompiler) compileEventGroupWithTimeoutRecovery(
	ctx fql.IWaitForEventGroupExpressionContext,
	timeoutLabel, endLabel core.Label,
) bytecode.Operand {
	if waitForEventGroupTriggerClause(ctx) != nil {
		return c.compileEventGroupWithTriggerTimeoutCleanup(ctx, timeoutLabel, endLabel)
	}

	state, ok := c.buildWaitEventGroupState(ctx)
	if !ok {
		return bytecode.NoopOperand
	}

	streamReg := c.ctx.Function.Registers.Allocate()
	resultReg := c.ctx.Function.Registers.Allocate()
	timeoutStateReg := c.ctx.Function.Registers.Allocate()
	c.ctx.Program.Emitter.EmitLoadNone(resultReg)
	c.ctx.Program.Emitter.EmitBoolean(timeoutStateReg, false)
	c.emitWaitEventGroupStreamSetup(state, streamReg)
	c.compileWaitEventGroupTrigger(ctx)

	start := c.ctx.Program.Emitter.NewLabel()
	done := c.ctx.Program.Emitter.NewLabel()
	cleanup := c.ctx.Program.Emitter.NewLabel()
	c.ctx.Program.Emitter.MarkLabel(start)
	c.emitWaitEventGroupIteration(state, streamReg, resultReg, timeoutStateReg, start, done)
	c.ctx.Program.Emitter.EmitJump(cleanup)
	c.ctx.Program.Emitter.MarkLabel(done)
	c.ctx.Program.Emitter.EmitJump(cleanup)
	c.ctx.Program.Emitter.MarkLabel(cleanup)
	c.emitWaitEventGroupCleanup(state, streamReg)
	c.ctx.Program.Emitter.EmitJumpIfTrue(timeoutStateReg, timeoutLabel)
	c.ctx.Program.Emitter.EmitJump(endLabel)

	return resultReg
}

func (c *WaitCompiler) compileEventGroupWithTriggerCleanup(
	ctx fql.IWaitForEventGroupExpressionContext,
) bytecode.Operand {
	streamReg := c.ctx.Function.Registers.Allocate()
	resultReg := c.ctx.Function.Registers.Allocate()
	streamReadyReg := c.ctx.Function.Registers.Allocate()
	c.ctx.Program.Emitter.EmitLoadNone(resultReg)
	c.ctx.Program.Emitter.EmitBoolean(streamReadyReg, false)

	startCatch := c.ctx.Program.Emitter.Size()
	state, ok := c.buildWaitEventGroupState(ctx)
	if !ok {
		return bytecode.NoopOperand
	}

	c.emitWaitEventGroupStreamSetupWithReady(state, streamReg, streamReadyReg)
	c.compileWaitEventGroupTrigger(ctx)

	start := c.ctx.Program.Emitter.NewLabel()
	done := c.ctx.Program.Emitter.NewLabel()
	exit := c.ctx.Program.Emitter.NewLabel("waitfor", "event", "group", "trigger", "exit")
	c.ctx.Program.Emitter.MarkLabel(start)
	c.emitWaitEventGroupIteration(state, streamReg, resultReg, bytecode.NoopOperand, start, done)
	c.ctx.Program.Emitter.MarkLabel(done)
	c.emitWaitEventGroupCleanupIfReady(state, streamReg, streamReadyReg)
	c.ctx.Program.Emitter.EmitJump(exit)

	endCatchExclusive := c.ctx.Program.Emitter.Size()
	errorHandlerPC := c.ctx.Program.Emitter.Size()
	c.emitWaitEventGroupCleanupIfReady(state, streamReg, streamReadyReg)
	c.ctx.Program.Emitter.Emit(bytecode.OpRethrow)
	c.ctx.Program.Emitter.MarkLabel(exit)
	c.ctx.Program.Emitter.EmitAB(bytecode.OpMove, resultReg, resultReg)
	c.ctx.Program.CatchTable.Push(startCatch, endCatchExclusive-1, errorHandlerPC)

	return resultReg
}

func (c *WaitCompiler) compileEventGroupWithTriggerTimeoutCleanup(
	ctx fql.IWaitForEventGroupExpressionContext,
	timeoutLabel, endLabel core.Label,
) bytecode.Operand {
	streamReg := c.ctx.Function.Registers.Allocate()
	resultReg := c.ctx.Function.Registers.Allocate()
	timeoutStateReg := c.ctx.Function.Registers.Allocate()
	streamReadyReg := c.ctx.Function.Registers.Allocate()
	c.ctx.Program.Emitter.EmitLoadNone(resultReg)
	c.ctx.Program.Emitter.EmitBoolean(timeoutStateReg, false)
	c.ctx.Program.Emitter.EmitBoolean(streamReadyReg, false)

	startCatch := c.ctx.Program.Emitter.Size()
	state, ok := c.buildWaitEventGroupState(ctx)
	if !ok {
		return bytecode.NoopOperand
	}

	c.emitWaitEventGroupStreamSetupWithReady(state, streamReg, streamReadyReg)
	c.compileWaitEventGroupTrigger(ctx)

	start := c.ctx.Program.Emitter.NewLabel()
	done := c.ctx.Program.Emitter.NewLabel()
	cleanup := c.ctx.Program.Emitter.NewLabel()
	c.ctx.Program.Emitter.MarkLabel(start)
	c.emitWaitEventGroupIteration(state, streamReg, resultReg, timeoutStateReg, start, done)
	c.ctx.Program.Emitter.EmitJump(cleanup)
	c.ctx.Program.Emitter.MarkLabel(done)
	c.ctx.Program.Emitter.EmitJump(cleanup)
	c.ctx.Program.Emitter.MarkLabel(cleanup)
	c.emitWaitEventGroupCleanupIfReady(state, streamReg, streamReadyReg)
	c.ctx.Program.Emitter.EmitJumpIfTrue(timeoutStateReg, timeoutLabel)
	c.ctx.Program.Emitter.EmitJump(endLabel)

	endCatchExclusive := c.ctx.Program.Emitter.Size()
	errorHandlerPC := c.ctx.Program.Emitter.Size()
	c.emitWaitEventGroupCleanupIfReady(state, streamReg, streamReadyReg)
	c.ctx.Program.Emitter.Emit(bytecode.OpRethrow)
	c.ctx.Program.CatchTable.Push(startCatch, endCatchExclusive-1, errorHandlerPC)

	return resultReg
}

func (c *WaitCompiler) buildWaitEventGroupState(
	ctx fql.IWaitForEventGroupExpressionContext,
) (waitEventGroupCompileState, bool) {
	if ctx == nil {
		return waitEventGroupCompileState{}, false
	}

	entries := ctx.AllWaitForEventGroupEntry()
	if len(entries) == 0 {
		return waitEventGroupCompileState{}, false
	}

	state := waitEventGroupCompileState{
		span:            waitForSpan(ctx, nil),
		arms:            make([]waitEventGroupArm, 0, len(entries)),
		synchronization: resolveWaitForSynchronization(ctx.WaitForSynchronization()),
		sourcesReg:      c.ctx.Function.Registers.Allocate(),
		namesReg:        c.ctx.Function.Registers.Allocate(),
		optionsReg:      c.ctx.Function.Registers.Allocate(),
	}
	c.ctx.Program.Emitter.EmitArray(state.sourcesReg, len(entries))
	c.ctx.Program.Emitter.EmitArray(state.namesReg, len(entries))
	c.ctx.Program.Emitter.EmitArray(state.optionsReg, len(entries))

	for _, entry := range entries {
		sourceReg := c.CompileWaitForEventSource(entry.WaitForEventSource())
		nameReg := c.CompileWaitForEventName(entry.WaitForEventName())

		if sourceReg == bytecode.NoopOperand || nameReg == bytecode.NoopOperand {
			return waitEventGroupCompileState{}, false
		}

		optionsReg := bytecode.NoopOperand
		if options := entry.OptionsClause(); options != nil {
			optionsReg = c.CompileOptionsClause(options)
		} else {
			optionsReg = c.ctx.Function.Registers.Allocate()
			c.ctx.Program.Emitter.EmitLoadNone(optionsReg)
		}

		c.ctx.Program.Emitter.EmitArrayPush(state.sourcesReg, sourceReg)
		c.ctx.Program.Emitter.EmitArrayPush(state.namesReg, nameReg)
		c.ctx.Program.Emitter.EmitArrayPush(state.optionsReg, optionsReg)
		state.arms = append(state.arms, waitEventGroupArm{
			span:    waitForSpan(entry, ctx),
			filters: entry.AllEventFilterClause(),
		})
	}

	if timeout := waitForEventGroupTimeoutClause(ctx); timeout != nil {
		state.timeoutReg = c.recovery.CompileDurationExpression(timeout.Expression())

		if state.timeoutReg == bytecode.NoopOperand {
			return waitEventGroupCompileState{}, false
		}
	}

	if state.synchronization == waitForSynchronizationAll {
		state.matchedRegs = make([]bytecode.Operand, len(entries))
		state.valueRegs = make([]bytecode.Operand, len(entries))

		for idx := range entries {
			state.matchedRegs[idx] = c.ctx.Function.Registers.Allocate()
			state.valueRegs[idx] = c.ctx.Function.Registers.Allocate()
			c.ctx.Program.Emitter.EmitBoolean(state.matchedRegs[idx], false)
			c.ctx.Program.Emitter.EmitLoadNone(state.valueRegs[idx])
		}
	}

	return state, true
}

func (c *WaitCompiler) emitWaitEventGroupStreamSetup(
	state waitEventGroupCompileState,
	streamReg bytecode.Operand,
) {
	c.emitWaitEventGroupStreamSetupWithReady(state, streamReg, bytecode.NoopOperand)
}

func (c *WaitCompiler) emitWaitEventGroupStreamSetupWithReady(
	state waitEventGroupCompileState,
	streamReg, streamReadyReg bytecode.Operand,
) {
	c.ctx.Program.Emitter.WithSpan(state.span, func() {
		c.ctx.Program.Emitter.EmitMove(streamReg, state.sourcesReg)
		c.ctx.Program.Emitter.EmitABC(bytecode.OpStreamGroup, streamReg, state.namesReg, state.optionsReg)
		if streamReadyReg != bytecode.NoopOperand {
			c.ctx.Program.Emitter.EmitBoolean(streamReadyReg, true)
		}
		c.ctx.Program.Emitter.EmitABC(bytecode.OpStreamIter, streamReg, streamReg, state.timeoutReg)
	})
}

func (c *WaitCompiler) compileWaitEventGroupTrigger(ctx fql.IWaitForEventGroupExpressionContext) {
	if trigger := waitForEventGroupTriggerClause(ctx); trigger != nil {
		c.compileWaitForTriggerClause(trigger)
	}
}

func (c *WaitCompiler) emitWaitEventGroupIteration(
	state waitEventGroupCompileState,
	streamReg, resultReg, timeoutStateReg bytecode.Operand,
	restartLabel, doneLabel core.Label,
) {
	c.ctx.Program.Emitter.WithSpan(state.span, func() {
		if timeoutStateReg != bytecode.NoopOperand {
			c.ctx.Program.Emitter.EmitIterNextTimeout(streamReg, timeoutStateReg, doneLabel)
		} else {
			c.ctx.Program.Emitter.EmitIterNext(streamReg, doneLabel)
		}
	})

	eventReg := c.ctx.Function.Registers.Allocate()
	armIndexReg := c.ctx.Function.Registers.Allocate()
	c.ctx.Program.Emitter.EmitIterValue(eventReg, streamReg)
	c.ctx.Program.Emitter.EmitIterKey(armIndexReg, streamReg)

	if state.synchronization == waitForSynchronizationAny {
		c.emitWaitEventAnyDispatch(state, streamReg, eventReg, armIndexReg, resultReg, restartLabel, doneLabel)
		return
	}

	c.emitWaitEventAllDispatch(state, streamReg, eventReg, armIndexReg, resultReg, restartLabel, doneLabel)
}

func (c *WaitCompiler) emitWaitEventAnyDispatch(
	state waitEventGroupCompileState,
	_ bytecode.Operand,
	eventReg, armIndexReg, resultReg bytecode.Operand,
	restartLabel, doneLabel core.Label,
) {
	for idx, arm := range state.arms {
		nextArm := c.ctx.Program.Emitter.NewLabel()
		indexConst := c.ctx.Function.Symbols.AddConstant(runtime.NewInt(idx))
		c.ctx.Program.Emitter.EmitJumpCompare(bytecode.OpJumpIfNeConst, armIndexReg, indexConst, nextArm)
		c.emitWaitEventGroupFilters(arm, eventReg, restartLabel)
		c.ctx.Program.Emitter.EmitMove(resultReg, eventReg)
		c.ctx.Program.Emitter.EmitJump(doneLabel)
		c.ctx.Program.Emitter.MarkLabel(nextArm)
	}

	c.ctx.Program.Emitter.EmitJump(restartLabel)
}

func (c *WaitCompiler) emitWaitEventAllDispatch(
	state waitEventGroupCompileState,
	streamReg, eventReg, armIndexReg, resultReg bytecode.Operand,
	restartLabel, doneLabel core.Label,
) {
	for idx, arm := range state.arms {
		nextArm := c.ctx.Program.Emitter.NewLabel()
		indexConst := c.ctx.Function.Symbols.AddConstant(runtime.NewInt(idx))
		c.ctx.Program.Emitter.EmitJumpCompare(bytecode.OpJumpIfNeConst, armIndexReg, indexConst, nextArm)
		c.ctx.Program.Emitter.EmitJumpIfTrue(state.matchedRegs[idx], restartLabel)
		c.emitWaitEventGroupFilters(arm, eventReg, restartLabel)
		c.ctx.Program.Emitter.EmitMove(state.valueRegs[idx], eventReg)
		c.ctx.Program.Emitter.EmitBoolean(state.matchedRegs[idx], true)
		indexReg := c.facts.LoadConstant(runtime.NewInt(idx))
		c.ctx.Program.Emitter.EmitAB(bytecode.OpStreamGroupArmDone, streamReg, indexReg)

		for _, matchedReg := range state.matchedRegs {
			c.ctx.Program.Emitter.EmitJumpIfFalse(matchedReg, restartLabel)
		}

		c.ctx.Program.Emitter.EmitArray(resultReg, len(state.valueRegs))

		for _, valueReg := range state.valueRegs {
			c.ctx.Program.Emitter.EmitArrayPush(resultReg, valueReg)
		}

		c.ctx.Program.Emitter.EmitJump(doneLabel)
		c.ctx.Program.Emitter.MarkLabel(nextArm)
	}

	c.ctx.Program.Emitter.EmitJump(restartLabel)
}

func (c *WaitCompiler) emitWaitEventGroupFilters(
	arm waitEventGroupArm,
	eventReg bytecode.Operand,
	restartLabel core.Label,
) {
	if len(arm.filters) == 0 {
		return
	}

	c.ctx.Function.Symbols.EnterScope()
	defer c.ctx.Function.Symbols.ExitScope()
	c.ctx.Function.Symbols.AssignLocal(core.PseudoVariable, core.TypeUnknown, eventReg)

	for _, filter := range arm.filters {
		conditionReg := c.exprs.CompileWithImplicitCurrent(filter.Expression())
		c.ctx.Program.Emitter.EmitJumpIfFalse(conditionReg, restartLabel)
	}
}

func (c *WaitCompiler) emitWaitEventGroupCleanup(
	state waitEventGroupCompileState,
	streamReg bytecode.Operand,
) {
	c.ctx.Program.Emitter.WithSpan(state.span, func() {
		c.ctx.Program.Emitter.EmitA(bytecode.OpClose, streamReg)
	})
}

func (c *WaitCompiler) emitWaitEventGroupCleanupIfReady(
	state waitEventGroupCompileState,
	streamReg, streamReadyReg bytecode.Operand,
) {
	if streamReadyReg == bytecode.NoopOperand {
		c.emitWaitEventGroupCleanup(state, streamReg)
		return
	}

	skip := c.ctx.Program.Emitter.NewLabel("waitfor", "event", "group", "cleanup", "skip")
	c.ctx.Program.Emitter.EmitJumpIfFalse(streamReadyReg, skip)
	c.emitWaitEventGroupCleanup(state, streamReg)
	c.ctx.Program.Emitter.EmitBoolean(streamReadyReg, false)
	c.ctx.Program.Emitter.MarkLabel(skip)
}
