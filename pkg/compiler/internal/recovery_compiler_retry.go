package internal

import (
	"errors"
	"strconv"
	"strings"

	"github.com/antlr4-go/antlr/v4"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/compiler/internal/core"
	parserd "github.com/MontFerret/ferret/v2/pkg/parser/diagnostics"
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
			c.reportInvalidDurationLiteral(dl, err)
			return bytecode.NoopOperand
		}

		return c.facts.LoadConstant(val)
	}

	if il := clause.IntegerLiteral(); il != nil {
		return c.literals.CompileIntegerLiteral(il)
	}

	if fl := clause.FloatLiteral(); fl != nil {
		return c.literals.CompileFloatLiteral(fl)
	}

	if v := clause.Variable(); v != nil {
		return c.exprs.CompileVariable(v)
	}

	if p := clause.Param(); p != nil {
		return c.exprs.CompileParam(p)
	}

	if me := clause.MemberExpression(); me != nil {
		return c.exprs.CompileMemberExpression(me)
	}

	if fc := clause.FunctionCall(); fc != nil {
		return c.exprs.CompileFunctionCall(fc, false)
	}

	return bytecode.NoopOperand
}

func (c *RecoveryCompiler) EmitRetryDelay(retry *core.RecoveryRetryPlan, state core.RetryDelayState) {
	if c == nil || c.ctx == nil || retry == nil || !retry.HasDelay {
		return
	}

	delayReady := c.ctx.Program.Emitter.NewLabel("recovery", "retry", "delay", "ready")
	c.ctx.Program.Emitter.EmitJumpIfTrue(state.ReadyReg, delayReady)

	delayValue := ensureOperandRegister(c.ctx, c.facts, c.CompileDurationOperand(retry.Delay))
	if delayValue == bytecode.NoopOperand {
		return
	}

	c.facts.EmitMoveAuto(state.BaseReg, delayValue)
	c.facts.EmitMoveAuto(state.CurrentReg, state.BaseReg)
	c.ctx.Program.Emitter.EmitBoolean(state.ReadyReg, true)
	c.ctx.Program.Emitter.MarkLabel(delayReady)

	c.ctx.Program.Emitter.EmitA(bytecode.OpSleep, state.CurrentReg)

	if retry.Backoff != core.RetryBackoffNone {
		c.emitBackoffUpdate(retry.Backoff, state.CurrentReg, state.BaseReg)
	}
}

func (c *RecoveryCompiler) initRetryDelayState(retry *core.RecoveryRetryPlan) core.RetryDelayState {
	if c == nil || c.ctx == nil || retry == nil || !retry.HasDelay {
		return core.RetryDelayState{}
	}

	state := core.RetryDelayState{
		BaseReg:    c.ctx.Function.Registers.Allocate(),
		CurrentReg: c.ctx.Function.Registers.Allocate(),
		ReadyReg:   c.ctx.Function.Registers.Allocate(),
	}

	c.ctx.Program.Emitter.EmitBoolean(state.ReadyReg, false)

	return state
}

func (c *RecoveryCompiler) emitBackoffUpdate(strategy core.RetryBackoff, intervalReg, baseEveryReg bytecode.Operand) {
	switch strategy {
	case core.RetryBackoffLinear:
		c.ctx.Program.Emitter.EmitABC(bytecode.OpAdd, intervalReg, intervalReg, baseEveryReg)
	case core.RetryBackoffExponential:
		twoReg := c.facts.LoadConstant(runtime.NewInt(2))
		c.ctx.Program.Emitter.EmitABC(bytecode.OpMul, intervalReg, intervalReg, twoReg)
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

func (c *RecoveryCompiler) reportInvalidDurationLiteral(ctx antlr.ParserRuleContext, err error) {
	if c == nil || c.ctx == nil || c.ctx.Program.Errors == nil || ctx == nil {
		core.PanicInvariant("cannot report invalid duration literal")
	}

	message := "Invalid duration literal"
	hint := "Use a valid duration, e.g. 100ms, 2s, or 1.5m."

	if errors.Is(err, strconv.ErrRange) {
		message = "Duration literal is out of range"
		hint = "Use a duration value that stays within the supported range, e.g. 100ms, 2s, or 1.5m."
	}

	diag := c.ctx.Program.Errors.Create(parserd.SyntaxError, ctx, message)
	diag.Hint = hint
	c.ctx.Program.Errors.Add(diag)
}
