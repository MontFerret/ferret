package internal

import (
	"math"
	"strconv"
	"strings"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/compiler/internal/core"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

func (c *WaitCompiler) emitWaitPredicateTimeoutCheck(
	timeoutReg, startReg, unitReg bytecode.Operand,
	timeoutLabel core.Label,
) bytecode.Operand {
	if timeoutReg == bytecode.NoopOperand {
		return bytecode.NoopOperand
	}

	nowReg := c.emitNow()
	elapsedReg := c.emitDateDiff(startReg, nowReg, unitReg)
	reachedReg := c.ctx.Registers.Allocate()
	c.ctx.Emitter.EmitGte(reachedReg, elapsedReg, timeoutReg)
	c.ctx.Emitter.EmitJumpIfTrue(reachedReg, timeoutLabel)

	return elapsedReg
}

func (c *WaitCompiler) prepareWaitSleepInterval(config waitPredicateCompileConfig, pollReg bytecode.Operand) bytecode.Operand {
	if !config.hasJitter && config.capEveryReg == bytecode.NoopOperand {
		return pollReg
	}

	sleepIntervalReg := c.ctx.Registers.Allocate()
	c.ctx.Emitter.EmitMove(sleepIntervalReg, pollReg)

	if config.hasJitter {
		c.emitApplyJitter(sleepIntervalReg, config.jitterReg)
	}

	if config.capEveryReg != bytecode.NoopOperand {
		c.emitClampMax(sleepIntervalReg, config.capEveryReg)
	}

	return sleepIntervalReg
}

func (c *WaitCompiler) emitNow() bytecode.Operand {
	return c.front.Expressions.CompileFunctionCallByNameWith(nil, runtime.NewString("NOW"), false, nil)
}

func (c *WaitCompiler) emitDateDiff(start, end, unit bytecode.Operand) bytecode.Operand {
	return c.emitFunctionCall(runtime.NewString("DATE_DIFF"), start, end, unit)
}

func (c *WaitCompiler) emitFunctionCall(name runtime.String, args ...bytecode.Operand) bytecode.Operand {
	if len(args) == 0 {
		return c.front.Expressions.CompileFunctionCallByNameWith(nil, name, false, nil)
	}

	seq := c.ctx.Registers.AllocateSequence(len(args))

	for i, arg := range args {
		c.ctx.Emitter.EmitMove(seq[i], arg)
		c.ctx.Types.Set(seq[i], c.front.TypeFacts.OperandType(arg))
	}

	return c.front.Expressions.CompileFunctionCallByNameWith(nil, name, false, seq)
}

func (c *WaitCompiler) emitWaitSleep(intervalReg, timeoutReg, elapsedReg bytecode.Operand) {
	if timeoutReg == bytecode.NoopOperand {
		c.ctx.Emitter.EmitA(bytecode.OpSleep, intervalReg)
		return
	}

	sleepReg := c.ctx.Registers.Allocate()
	c.ctx.Emitter.EmitMove(sleepReg, intervalReg)

	remainingReg := c.ctx.Registers.Allocate()
	c.ctx.Emitter.EmitABC(bytecode.OpSub, remainingReg, timeoutReg, elapsedReg)

	shouldTrim := c.ctx.Registers.Allocate()
	c.ctx.Emitter.EmitLt(shouldTrim, remainingReg, sleepReg)

	useRemaining := c.ctx.Emitter.NewLabel()
	continueSleep := c.ctx.Emitter.NewLabel()

	c.ctx.Emitter.EmitJumpIfTrue(shouldTrim, useRemaining)
	c.ctx.Emitter.EmitJump(continueSleep)

	c.ctx.Emitter.MarkLabel(useRemaining)
	c.ctx.Emitter.EmitMove(sleepReg, remainingReg)
	c.ctx.Emitter.MarkLabel(continueSleep)

	c.ctx.Emitter.EmitA(bytecode.OpSleep, sleepReg)
}

func (c *WaitCompiler) emitClampMin(target, min bytecode.Operand) {
	lessReg := c.ctx.Registers.Allocate()
	c.ctx.Emitter.EmitLt(lessReg, target, min)

	useMin := c.ctx.Emitter.NewLabel()
	end := c.ctx.Emitter.NewLabel()

	c.ctx.Emitter.EmitJumpIfTrue(lessReg, useMin)
	c.ctx.Emitter.EmitJump(end)

	c.ctx.Emitter.MarkLabel(useMin)
	c.ctx.Emitter.EmitMove(target, min)
	c.ctx.Emitter.MarkLabel(end)
}

func (c *WaitCompiler) emitClampMax(target, max bytecode.Operand) {
	greaterReg := c.ctx.Registers.Allocate()
	c.ctx.Emitter.EmitGt(greaterReg, target, max)

	useMax := c.ctx.Emitter.NewLabel()
	end := c.ctx.Emitter.NewLabel()

	c.ctx.Emitter.EmitJumpIfTrue(greaterReg, useMax)
	c.ctx.Emitter.EmitJump(end)

	c.ctx.Emitter.MarkLabel(useMax)
	c.ctx.Emitter.EmitMove(target, max)
	c.ctx.Emitter.MarkLabel(end)
}

func (c *WaitCompiler) emitClampRange(target, min, max bytecode.Operand) {
	c.emitClampMin(target, min)
	c.emitClampMax(target, max)
}

func (c *WaitCompiler) emitApplyJitter(intervalReg, jitterReg bytecode.Operand) {
	if intervalReg == bytecode.NoopOperand || jitterReg == bytecode.NoopOperand {
		return
	}

	randReg := c.ctx.Registers.Allocate()
	c.ctx.Emitter.EmitA(bytecode.OpRand, randReg)

	twoReg := c.front.TypeFacts.LoadConstant(runtime.NewFloat(2))
	oneReg := c.front.TypeFacts.LoadConstant(runtime.NewFloat(1))

	twoJitterReg := c.ctx.Registers.Allocate()
	c.ctx.Emitter.EmitABC(bytecode.OpMul, twoJitterReg, jitterReg, twoReg)

	randScaleReg := c.ctx.Registers.Allocate()
	c.ctx.Emitter.EmitABC(bytecode.OpMul, randScaleReg, randReg, twoJitterReg)

	oneMinusReg := c.ctx.Registers.Allocate()
	c.ctx.Emitter.EmitABC(bytecode.OpSub, oneMinusReg, oneReg, jitterReg)

	multiplierReg := c.ctx.Registers.Allocate()
	c.ctx.Emitter.EmitABC(bytecode.OpAdd, multiplierReg, oneMinusReg, randScaleReg)

	c.ctx.Emitter.EmitABC(bytecode.OpMul, intervalReg, intervalReg, multiplierReg)
}

func parseDurationLiteral(text string) (runtime.Value, error) {
	raw := normalizeDurationLiteral(text)
	if raw == "" {
		return runtime.None, strconv.ErrSyntax
	}

	number, unit, ok := splitDurationLiteral(raw)
	if !ok || number == "" {
		return runtime.None, strconv.ErrSyntax
	}

	value, err := parseDurationLiteralNumber(number)
	if err != nil {
		return runtime.None, err
	}

	multiplier, ok := durationUnitMultiplier(unit)
	if !ok {
		return runtime.None, strconv.ErrSyntax
	}

	ms := value * multiplier
	if math.IsNaN(ms) || math.IsInf(ms, 0) {
		return runtime.None, strconv.ErrRange
	}

	if frac := math.Mod(ms, 1); frac == 0 {
		const (
			maxInt64Float = float64(1<<63 - 1)
			minInt64Float = -float64(1 << 63)
		)

		if ms < minInt64Float || ms > maxInt64Float {
			return runtime.None, strconv.ErrRange
		}
	}

	return durationValueFromMilliseconds(ms), nil
}

func normalizeDurationLiteral(text string) string {
	return strings.ToUpper(strings.TrimSpace(text))
}

func splitDurationLiteral(raw string) (string, string, bool) {
	switch {
	case strings.HasSuffix(raw, "MS"):
		return strings.TrimSuffix(raw, "MS"), "MS", true
	case strings.HasSuffix(raw, "S"):
		return strings.TrimSuffix(raw, "S"), "S", true
	case strings.HasSuffix(raw, "M"):
		return strings.TrimSuffix(raw, "M"), "M", true
	case strings.HasSuffix(raw, "H"):
		return strings.TrimSuffix(raw, "H"), "H", true
	case strings.HasSuffix(raw, "D"):
		return strings.TrimSuffix(raw, "D"), "D", true
	default:
		return "", "", false
	}
}

func parseDurationLiteralNumber(raw string) (float64, error) {
	return strconv.ParseFloat(raw, 64)
}

func durationUnitMultiplier(unit string) (float64, bool) {
	switch unit {
	case "MS":
		return 1, true
	case "S":
		return 1000, true
	case "M":
		return 60000, true
	case "H":
		return 3600000, true
	case "D":
		return 86400000, true
	default:
		return 0, false
	}
}

func durationValueFromMilliseconds(ms float64) runtime.Value {
	if frac := math.Mod(ms, 1); frac == 0 {
		return runtime.NewInt64(int64(ms))
	}

	return runtime.NewFloat(ms)
}
