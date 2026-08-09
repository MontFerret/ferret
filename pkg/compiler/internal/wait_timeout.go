package internal

import (
	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/compiler/internal/core"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

func (c *WaitCompiler) emitWaitPredicateTimeoutCheck(
	timeoutReg, startReg bytecode.Operand,
	timeoutLabel core.Label,
) bytecode.Operand {
	if timeoutReg == bytecode.NoopOperand {
		return bytecode.NoopOperand
	}

	currentReg := c.emitElapsed()
	elapsedReg := c.ctx.Function.Registers.Allocate()
	c.ctx.Program.Emitter.EmitABC(bytecode.OpSub, elapsedReg, currentReg, startReg)
	c.ctx.Function.Types.Set(elapsedReg, core.TypeDuration)
	reachedReg := c.ctx.Function.Registers.Allocate()
	c.ctx.Program.Emitter.EmitGte(reachedReg, elapsedReg, timeoutReg)
	c.ctx.Program.Emitter.EmitJumpIfTrue(reachedReg, timeoutLabel)

	return elapsedReg
}

func (c *WaitCompiler) prepareWaitSleepInterval(config waitPredicateScheduleConfig, pollReg bytecode.Operand) bytecode.Operand {
	if config.everyZero || (!config.hasJitter && config.capEveryReg == bytecode.NoopOperand) {
		return pollReg
	}

	sleepIntervalReg := c.ctx.Function.Registers.Allocate()
	c.ctx.Program.Emitter.EmitMove(sleepIntervalReg, pollReg)

	if config.hasJitter {
		c.emitApplyJitter(sleepIntervalReg, config.jitterReg)
	}

	if config.capEveryReg != bytecode.NoopOperand {
		c.emitClampMax(sleepIntervalReg, config.capEveryReg)
	}

	return sleepIntervalReg
}

func (c *WaitCompiler) emitElapsed() bytecode.Operand {
	dst := c.ctx.Function.Registers.Allocate()
	c.ctx.Program.Emitter.EmitA(bytecode.OpElapsed, dst)
	c.ctx.Function.Types.Set(dst, core.TypeDuration)

	return dst
}

func (c *WaitCompiler) emitWaitSleep(intervalReg, timeoutReg, elapsedReg bytecode.Operand, intervalIsZero bool) {
	if timeoutReg == bytecode.NoopOperand || intervalIsZero {
		c.ctx.Program.Emitter.EmitA(bytecode.OpSleep, intervalReg)
		return
	}

	sleepReg := c.ctx.Function.Registers.Allocate()
	c.ctx.Program.Emitter.EmitMove(sleepReg, intervalReg)

	remainingReg := c.ctx.Function.Registers.Allocate()
	c.ctx.Program.Emitter.EmitABC(bytecode.OpSub, remainingReg, timeoutReg, elapsedReg)

	shouldTrim := c.ctx.Function.Registers.Allocate()
	c.ctx.Program.Emitter.EmitLt(shouldTrim, remainingReg, sleepReg)

	useRemaining := c.ctx.Program.Emitter.NewLabel()
	continueSleep := c.ctx.Program.Emitter.NewLabel()

	c.ctx.Program.Emitter.EmitJumpIfTrue(shouldTrim, useRemaining)
	c.ctx.Program.Emitter.EmitJump(continueSleep)

	c.ctx.Program.Emitter.MarkLabel(useRemaining)
	c.ctx.Program.Emitter.EmitMove(sleepReg, remainingReg)
	c.ctx.Program.Emitter.MarkLabel(continueSleep)

	c.ctx.Program.Emitter.EmitA(bytecode.OpSleep, sleepReg)
}

func (c *WaitCompiler) emitClampMin(target, min bytecode.Operand) {
	lessReg := c.ctx.Function.Registers.Allocate()
	c.ctx.Program.Emitter.EmitLt(lessReg, target, min)

	useMin := c.ctx.Program.Emitter.NewLabel()
	end := c.ctx.Program.Emitter.NewLabel()

	c.ctx.Program.Emitter.EmitJumpIfTrue(lessReg, useMin)
	c.ctx.Program.Emitter.EmitJump(end)

	c.ctx.Program.Emitter.MarkLabel(useMin)
	c.ctx.Program.Emitter.EmitMove(target, min)
	c.ctx.Program.Emitter.MarkLabel(end)
}

func (c *WaitCompiler) emitClampMax(target, max bytecode.Operand) {
	greaterReg := c.ctx.Function.Registers.Allocate()
	c.ctx.Program.Emitter.EmitGt(greaterReg, target, max)

	useMax := c.ctx.Program.Emitter.NewLabel()
	end := c.ctx.Program.Emitter.NewLabel()

	c.ctx.Program.Emitter.EmitJumpIfTrue(greaterReg, useMax)
	c.ctx.Program.Emitter.EmitJump(end)

	c.ctx.Program.Emitter.MarkLabel(useMax)
	c.ctx.Program.Emitter.EmitMove(target, max)
	c.ctx.Program.Emitter.MarkLabel(end)
}

func (c *WaitCompiler) emitClampRange(target, min, max bytecode.Operand) {
	c.emitClampMin(target, min)
	c.emitClampMax(target, max)
}

func (c *WaitCompiler) emitApplyJitter(intervalReg, jitterReg bytecode.Operand) {
	if intervalReg == bytecode.NoopOperand || jitterReg == bytecode.NoopOperand {
		return
	}

	randReg := c.ctx.Function.Registers.Allocate()
	c.ctx.Program.Emitter.EmitA(bytecode.OpRand, randReg)

	twoReg := c.facts.LoadConstant(runtime.NewFloat(2))
	oneReg := c.facts.LoadConstant(runtime.NewFloat(1))

	twoJitterReg := c.ctx.Function.Registers.Allocate()
	c.ctx.Program.Emitter.EmitABC(bytecode.OpMul, twoJitterReg, jitterReg, twoReg)

	randScaleReg := c.ctx.Function.Registers.Allocate()
	c.ctx.Program.Emitter.EmitABC(bytecode.OpMul, randScaleReg, randReg, twoJitterReg)

	oneMinusReg := c.ctx.Function.Registers.Allocate()
	c.ctx.Program.Emitter.EmitABC(bytecode.OpSub, oneMinusReg, oneReg, jitterReg)

	multiplierReg := c.ctx.Function.Registers.Allocate()
	c.ctx.Program.Emitter.EmitABC(bytecode.OpAdd, multiplierReg, oneMinusReg, randScaleReg)

	c.ctx.Program.Emitter.EmitABC(bytecode.OpMul, intervalReg, intervalReg, multiplierReg)
}
