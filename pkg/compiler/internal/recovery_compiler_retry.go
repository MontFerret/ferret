package internal

import (
	"strings"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/compiler/internal/core"
	"github.com/MontFerret/ferret/v2/pkg/parser/fql"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

func (c *RecoveryCompiler) CompileDurationOperand(clause core.DurationClause) bytecode.Operand {
	if c == nil || c.ctx == nil || clause == nil {
		return bytecode.NoopOperand
	}

	if dl := clause.DurationLiteral(); dl != nil {
		val, err := parseDurationLiteral(dl.GetText())
		if err != nil {
			panic(err)
		}

		return c.front.TypeFacts.LoadConstant(val)
	}

	il := clause.IntegerLiteral()
	fl := clause.FloatLiteral()
	v := clause.Variable()
	p := clause.Param()
	me := clause.MemberExpression()
	fc := clause.FunctionCall()

	return compileFirstOperand(
		newOperandBranch(il != nil, func() bytecode.Operand { return c.front.Literals.CompileIntegerLiteral(il) }),
		newOperandBranch(fl != nil, func() bytecode.Operand { return c.front.Literals.CompileFloatLiteral(fl) }),
		newOperandBranch(v != nil, func() bytecode.Operand { return c.front.Expressions.CompileVariable(v) }),
		newOperandBranch(p != nil, func() bytecode.Operand { return c.front.Expressions.CompileParam(p) }),
		newOperandBranch(me != nil, func() bytecode.Operand { return c.front.Expressions.CompileMemberExpression(me) }),
		newOperandBranch(fc != nil, func() bytecode.Operand { return c.front.Expressions.CompileFunctionCall(fc, false) }),
	)
}

func (c *RecoveryCompiler) EmitRetryDelay(retry *core.RecoveryRetryPlan, state core.RetryDelayState) {
	if c == nil || c.ctx == nil || retry == nil || !retry.HasDelay {
		return
	}

	delayReady := c.ctx.Emitter.NewLabel("recovery", "retry", "delay", "ready")
	c.ctx.Emitter.EmitJumpIfTrue(state.ReadyReg, delayReady)

	delayValue := c.EnsureRegister(c.CompileDurationOperand(retry.Delay))
	c.front.TypeFacts.EmitMoveAuto(state.BaseReg, delayValue)
	c.front.TypeFacts.EmitMoveAuto(state.CurrentReg, state.BaseReg)
	c.ctx.Emitter.EmitBoolean(state.ReadyReg, true)
	c.ctx.Emitter.MarkLabel(delayReady)

	c.ctx.Emitter.EmitA(bytecode.OpSleep, state.CurrentReg)

	if retry.Backoff != core.RetryBackoffNone {
		c.emitBackoffUpdate(retry.Backoff, state.CurrentReg, state.BaseReg)
	}
}

func (c *RecoveryCompiler) initRetryDelayState(retry *core.RecoveryRetryPlan) core.RetryDelayState {
	if c == nil || c.ctx == nil || retry == nil || !retry.HasDelay {
		return core.RetryDelayState{}
	}

	state := core.RetryDelayState{
		BaseReg:    c.ctx.Registers.Allocate(),
		CurrentReg: c.ctx.Registers.Allocate(),
		ReadyReg:   c.ctx.Registers.Allocate(),
	}

	c.ctx.Emitter.EmitBoolean(state.ReadyReg, false)

	return state
}

func (c *RecoveryCompiler) emitBackoffUpdate(strategy core.RetryBackoff, intervalReg, baseEveryReg bytecode.Operand) {
	switch strategy {
	case core.RetryBackoffLinear:
		c.ctx.Emitter.EmitABC(bytecode.OpAdd, intervalReg, intervalReg, baseEveryReg)
	case core.RetryBackoffExponential:
		twoReg := c.front.TypeFacts.LoadConstant(runtime.NewInt(2))
		c.ctx.Emitter.EmitABC(bytecode.OpMul, intervalReg, intervalReg, twoReg)
	}
}

func (c *RecoveryCompiler) resolveRetryBackoff(clause fql.IRecoveryRetryBackoffClauseContext) (core.RetryBackoff, bool) {
	if clause == nil {
		return core.RetryBackoffNone, true
	}

	kind := clause.RecoveryRetryBackoffKind()
	if kind == nil {
		c.reportInvalidTail(clause, "Expected backoff kind after 'BACKOFF' in retry policy", "Use BACKOFF CONSTANT, BACKOFF LINEAR, or BACKOFF EXPONENTIAL.")
		return core.RetryBackoffNone, false
	}

	raw := ""

	switch {
	case kind.Identifier() != nil:
		raw = kind.Identifier().GetText()
	case kind.StringLiteral() != nil:
		if parsed, ok := parseStringLiteralConst(kind.StringLiteral()); ok {
			raw = parsed.String()
		}
	case kind.None() != nil:
		raw = kind.None().GetText()
	}

	switch strings.ToUpper(strings.TrimSpace(raw)) {
	case "CONSTANT":
		return core.RetryBackoffNone, true
	case "LINEAR":
		return core.RetryBackoffLinear, true
	case "EXPONENTIAL":
		return core.RetryBackoffExponential, true
	default:
		c.reportInvalidTail(kind, "Unknown BACKOFF strategy", "Use one of: CONSTANT, LINEAR, EXPONENTIAL.")
		return core.RetryBackoffNone, false
	}
}
