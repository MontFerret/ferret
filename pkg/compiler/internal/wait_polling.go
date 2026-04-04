package internal

import (
	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/compiler/internal/core"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

type waitPredicatePollState struct {
	baseEveryReg bytecode.Operand
	pollReg      bytecode.Operand
	intervalReg  bytecode.Operand
	resultReg    bytecode.Operand
	startReg     bytecode.Operand
	unitReg      bytecode.Operand
}

const waitForDefaultEveryMs = 100

func (c *WaitCompiler) tryCompileWaitPredicateFastPath(config waitPredicateCompileConfig) (bytecode.Operand, bool) {
	switch config.mode {
	case waitForPredicateModeBool:
		truth, ok := literalTruthinessFromExpression(config.predExpr)
		if !ok {
			return bytecode.NoopOperand, false
		}

		if truth {
			return c.emitImmediateWaitBool(true), true
		}

		if config.timeoutReg != bytecode.NoopOperand {
			c.ctx.Program.Emitter.EmitA(bytecode.OpSleep, config.timeoutReg)
			return c.emitImmediateWaitBool(false), true
		}

		return bytecode.NoopOperand, false
	default:
		exists, ok := literalExistsFromExpression(config.predExpr)
		if !ok {
			return bytecode.NoopOperand, false
		}

		cond := exists
		if config.mode == waitForPredicateModeNotExists {
			cond = !exists
		}

		if cond {
			if config.mode == waitForPredicateModeValue {
				return c.exprs.Compile(config.predExpr), true
			}

			return c.emitImmediateWaitBool(true), true
		}

		if config.timeoutReg != bytecode.NoopOperand {
			c.ctx.Program.Emitter.EmitA(bytecode.OpSleep, config.timeoutReg)
			if config.mode == waitForPredicateModeValue {
				return c.emitImmediateWaitNone(), true
			}

			return c.emitImmediateWaitBool(false), true
		}

		return bytecode.NoopOperand, false
	}
}

func (c *WaitCompiler) emitImmediateWaitBool(value bool) bytecode.Operand {
	resultReg := c.ctx.Function.Registers.Allocate()
	c.ctx.Program.Emitter.EmitBoolean(resultReg, value)

	return resultReg
}

func (c *WaitCompiler) emitImmediateWaitNone() bytecode.Operand {
	resultReg := c.ctx.Function.Registers.Allocate()
	c.ctx.Program.Emitter.EmitLoadNone(resultReg)

	return resultReg
}

func (c *WaitCompiler) initWaitPredicatePollState(config waitPredicateCompileConfig) waitPredicatePollState {
	state := waitPredicatePollState{
		baseEveryReg: c.ctx.Function.Registers.Allocate(),
	}

	if config.everyReg != bytecode.NoopOperand {
		c.ctx.Program.Emitter.EmitMove(state.baseEveryReg, config.everyReg)
	} else {
		c.ctx.Program.Emitter.EmitLoadConst(state.baseEveryReg, c.ctx.Function.Symbols.AddConstant(runtime.NewInt(waitForDefaultEveryMs)))
	}

	state.pollReg = state.baseEveryReg
	if config.backoff != core.RetryBackoffNone {
		state.intervalReg = c.ctx.Function.Registers.Allocate()
		c.ctx.Program.Emitter.EmitMove(state.intervalReg, state.baseEveryReg)
		state.pollReg = state.intervalReg
	}

	state.resultReg = c.ctx.Function.Registers.Allocate()
	if config.mode == waitForPredicateModeValue {
		c.ctx.Program.Emitter.EmitLoadNone(state.resultReg)
	} else {
		c.ctx.Program.Emitter.EmitBoolean(state.resultReg, false)
	}

	if config.timeoutReg != bytecode.NoopOperand {
		state.startReg = c.emitNow()
		state.unitReg = c.facts.LoadConstant(runtime.NewString("f"))
	}

	return state
}

func (c *WaitCompiler) emitWaitPredicatePollLoop(config waitPredicateCompileConfig, state waitPredicatePollState) {
	start := c.ctx.Program.Emitter.NewLabel()
	success := c.ctx.Program.Emitter.NewLabel()
	timeoutLabel := c.ctx.Program.Emitter.NewLabel()
	end := c.ctx.Program.Emitter.NewLabel()

	c.ctx.Program.Emitter.MarkLabel(start)
	valueReg := c.emitWaitPredicatePollIteration(config, state, start, success, timeoutLabel)

	c.ctx.Program.Emitter.MarkLabel(success)
	c.emitWaitSuccessResult(config.mode, state.resultReg, valueReg)
	c.ctx.Program.Emitter.EmitJump(end)

	c.ctx.Program.Emitter.MarkLabel(timeoutLabel)
	c.emitWaitTimeoutResult(config.mode, state.resultReg)
	c.ctx.Program.Emitter.MarkLabel(end)
}

func (c *WaitCompiler) emitWaitPredicatePollLoopWithRecovery(
	config waitPredicateCompileConfig,
	state waitPredicatePollState,
	timeoutLabel, endLabel core.Label,
) bytecode.Operand {
	start := c.ctx.Program.Emitter.NewLabel()
	success := c.ctx.Program.Emitter.NewLabel()

	c.ctx.Program.Emitter.MarkLabel(start)
	valueReg := c.emitWaitPredicatePollIteration(config, state, start, success, timeoutLabel)

	c.ctx.Program.Emitter.MarkLabel(success)
	c.emitWaitSuccessResult(config.mode, state.resultReg, valueReg)
	c.ctx.Program.Emitter.EmitJump(endLabel)

	return state.resultReg
}

func (c *WaitCompiler) emitWaitPredicatePollIteration(
	config waitPredicateCompileConfig,
	state waitPredicatePollState,
	startLabel, successLabel, timeoutLabel core.Label,
) bytecode.Operand {
	valueReg := c.exprs.Compile(config.predExpr)
	condReg := c.emitWaitPredicateCondition(config.mode, valueReg)
	c.ctx.Program.Emitter.EmitJumpIfTrue(condReg, successLabel)

	elapsedReg := c.emitWaitPredicateTimeoutCheck(config.timeoutReg, state.startReg, state.unitReg, timeoutLabel)
	sleepIntervalReg := c.prepareWaitSleepInterval(config, state.pollReg)
	c.emitWaitSleep(sleepIntervalReg, config.timeoutReg, elapsedReg)

	if config.backoff != core.RetryBackoffNone {
		c.recovery.emitBackoffUpdate(config.backoff, state.intervalReg, state.baseEveryReg)
		if config.capEveryReg != bytecode.NoopOperand {
			c.emitClampMax(state.intervalReg, config.capEveryReg)
		}
	}

	c.ctx.Program.Emitter.EmitJump(startLabel)

	return valueReg
}

func (c *WaitCompiler) emitWaitPredicateCondition(mode waitForPredicateMode, valueReg bytecode.Operand) bytecode.Operand {
	switch mode {
	case waitForPredicateModeValue, waitForPredicateModeExists:
		return c.emitExistsCheck(valueReg)
	case waitForPredicateModeNotExists:
		existsReg := c.emitExistsCheck(valueReg)
		condReg := c.ctx.Function.Registers.Allocate()
		c.ctx.Program.Emitter.EmitAB(bytecode.OpNot, condReg, existsReg)

		return condReg
	default:
		condReg := c.ctx.Function.Registers.Allocate()
		c.ctx.Program.Emitter.EmitAB(bytecode.OpCastBool, condReg, valueReg)

		return condReg
	}
}

func (c *WaitCompiler) emitWaitSuccessResult(mode waitForPredicateMode, resultReg, valueReg bytecode.Operand) {
	if mode == waitForPredicateModeValue {
		c.ctx.Program.Emitter.EmitMove(resultReg, valueReg)
		return
	}

	c.ctx.Program.Emitter.EmitBoolean(resultReg, true)
}

func (c *WaitCompiler) emitWaitTimeoutResult(mode waitForPredicateMode, resultReg bytecode.Operand) {
	if mode == waitForPredicateModeValue {
		c.ctx.Program.Emitter.EmitLoadNone(resultReg)
		return
	}

	c.ctx.Program.Emitter.EmitBoolean(resultReg, false)
}

func (c *WaitCompiler) emitExistsCheck(val bytecode.Operand) bytecode.Operand {
	dst := c.ctx.Function.Registers.Allocate()
	c.ctx.Program.Emitter.EmitAB(bytecode.OpExists, dst, val)
	c.ctx.Function.Types.Set(dst, core.TypeBool)

	return dst
}
