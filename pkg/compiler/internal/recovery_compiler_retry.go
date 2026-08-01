package internal

import (
	"strings"

	"github.com/antlr4-go/antlr/v4"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/compiler/internal/core"
	parserd "github.com/MontFerret/ferret/v2/pkg/parser/diagnostics"
	"github.com/MontFerret/ferret/v2/pkg/parser/fql"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

func (c *RecoveryCompiler) CompileDurationExpression(expr fql.IExpressionContext) bytecode.Operand {
	if c == nil || c.ctx == nil || expr == nil {
		return bytecode.NoopOperand
	}

	operand := c.exprs.Compile(expr)
	if value, ok := c.facts.LiteralValueFromExpression(expr); ok {
		if duration, ok := value.(runtime.Duration); ok && duration >= 0 {
			return operand
		}
	}

	return c.validateDurationOperand(expr.(antlr.ParserRuleContext), operand)
}

func (c *RecoveryCompiler) CompileRetryDelay(delay fql.IRecoveryRetryDelayValueContext) bytecode.Operand {
	if c == nil || c.ctx == nil || delay == nil || delay.Predicate() == nil {
		return bytecode.NoopOperand
	}

	predicate := delay.Predicate()
	operand := c.exprs.compilePredicate(predicate)
	unary := delay.UnaryOperator()

	if unary != nil {
		dst := c.ctx.Function.Registers.Allocate()
		var op bytecode.Opcode

		switch {
		case unary.Not() != nil:
			op = bytecode.OpNot
		case unary.Minus() != nil:
			op = bytecode.OpFlipNegative
		case unary.Plus() != nil:
			op = bytecode.OpFlipPositive
		default:
			return bytecode.NoopOperand
		}

		c.ctx.Program.Emitter.WithSpan(parserd.SpanFromRuleContext(delay.(antlr.ParserRuleContext)), func() {
			c.ctx.Program.Emitter.EmitAB(op, dst, operand)
		})
		c.ctx.Function.Types.Set(dst, c.facts.OperandType(operand))
		operand = dst
	}

	if unary == nil || unary.Plus() != nil {
		if value, ok := c.facts.LiteralValueFromPredicate(predicate); ok {
			if duration, ok := value.(runtime.Duration); ok && duration >= 0 {
				return operand
			}
		}
	}

	return c.validateDurationOperand(delay.(antlr.ParserRuleContext), operand)
}

func (c *RecoveryCompiler) validateDurationOperand(node antlr.ParserRuleContext, operand bytecode.Operand) bytecode.Operand {
	if operand == bytecode.NoopOperand {
		return bytecode.NoopOperand
	}

	dst := c.ctx.Function.Registers.Allocate()
	zero := c.facts.LoadConstant(runtime.ZeroDuration)
	valid := c.ctx.Program.Emitter.NewLabel("scheduling", "duration", "valid")
	isNegative := c.ctx.Function.Registers.Allocate()
	failure := c.ctx.Function.Symbols.AddConstant(runtime.NewString("wait duration must not be negative"))
	c.ctx.Program.Emitter.WithSpan(parserd.SpanFromRuleContext(node), func() {
		c.ctx.Program.Emitter.EmitABC(bytecode.OpAdd, dst, operand, zero)
		c.ctx.Program.Emitter.EmitLt(isNegative, dst, zero)
		c.ctx.Program.Emitter.EmitJumpIfFalse(isNegative, valid)
		c.ctx.Program.Emitter.EmitA(bytecode.OpFail, failure)
	})
	c.ctx.Program.Emitter.MarkLabel(valid)
	c.ctx.Function.Types.Set(dst, core.TypeDuration)
	c.ctx.Function.Types.Set(isNegative, core.TypeBool)

	return dst
}

func (c *RecoveryCompiler) EmitRetryDelay(retry *core.RecoveryRetryPlan, state core.RetryDelayState) {
	if c == nil || c.ctx == nil || retry == nil || !retry.HasDelay {
		return
	}

	delayReady := c.ctx.Program.Emitter.NewLabel("recovery", "retry", "delay", "ready")
	c.ctx.Program.Emitter.EmitJumpIfTrue(state.ReadyReg, delayReady)

	delayValue := ensureOperandRegister(c.ctx, c.facts, c.CompileRetryDelay(retry.Delay))
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
