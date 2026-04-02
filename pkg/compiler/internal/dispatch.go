package internal

import (
	"github.com/antlr4-go/antlr/v4"

	"github.com/MontFerret/ferret/v2/pkg/source"

	"github.com/MontFerret/ferret/v2/pkg/parser/diagnostics"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"

	"github.com/MontFerret/ferret/v2/pkg/compiler/internal/core"
	"github.com/MontFerret/ferret/v2/pkg/parser/fql"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

type DispatchCompiler struct {
	ctx   *CompilationSession
	front *CompilationFrontend
}

func NewDispatchCompiler(ctx *CompilationSession) *DispatchCompiler {
	return &DispatchCompiler{
		ctx: ctx,
	}
}

func (c *DispatchCompiler) Compile(ctx fql.IDispatchExpressionContext) bytecode.Operand {
	if ctx == nil {
		return bytecode.NoopOperand
	}

	plan := c.front.Recovery.CollectPlan(ctx, core.RecoveryPlanOptions{})
	return c.front.Recovery.CompileWithRecoveryPlan(plan, core.CatchJumpModeNone, func() bytecode.Operand {
		targetReg := c.ensureRegister(c.compileTarget(ctx.DispatchTarget()))
		eventReg := c.ensureRegister(c.compileEventName(ctx.DispatchEventName()))
		payloadReg := c.ensureRegister(c.compilePayload(ctx.DispatchWithClause()))
		optionsReg := c.ensureRegister(c.compileOptions(ctx.DispatchOptionsClause()))
		argsReg := c.buildDispatchArgs(payloadReg, optionsReg)

		dst := c.ctx.Registers.Allocate()
		span := dispatchSpan(ctx)

		c.ctx.Emitter.WithSpan(span, func() {
			c.ctx.Emitter.EmitMove(dst, targetReg)
			c.ctx.Emitter.EmitABC(bytecode.OpDispatch, dst, eventReg, argsReg)
		})

		c.ctx.Types.Set(dst, core.TypeNone)

		return dst
	})
}

func (c *DispatchCompiler) compileEventName(ctx fql.IDispatchEventNameContext) bytecode.Operand {
	if ctx == nil {
		return bytecode.NoopOperand
	}

	if sl := ctx.StringLiteral(); sl != nil {
		return c.front.Literals.CompileStringLiteral(sl)
	}

	if v := ctx.Variable(); v != nil {
		return c.front.Expressions.CompileVariable(v)
	}

	if p := ctx.Param(); p != nil {
		return c.front.Expressions.CompileParam(p)
	}

	if me := ctx.MemberExpression(); me != nil {
		return c.front.Expressions.CompileMemberExpression(me)
	}

	if fc := ctx.FunctionCall(); fc != nil {
		return c.front.Expressions.CompileFunctionCall(fc, false)
	}

	return bytecode.NoopOperand
}

func (c *DispatchCompiler) compileTarget(ctx fql.IDispatchTargetContext) bytecode.Operand {
	if ctx == nil {
		return bytecode.NoopOperand
	}

	if v := ctx.Variable(); v != nil {
		return c.front.Expressions.CompileVariable(v)
	}

	if p := ctx.Param(); p != nil {
		return c.front.Expressions.CompileParam(p)
	}

	if me := ctx.MemberExpression(); me != nil {
		return c.front.Expressions.CompileMemberExpression(me)
	}

	if fc := ctx.FunctionCallExpression(); fc != nil {
		return c.front.Expressions.CompileFunctionCallExpression(fc)
	}

	return bytecode.NoopOperand
}

func (c *DispatchCompiler) compilePayload(ctx fql.IDispatchWithClauseContext) bytecode.Operand {
	if ctx == nil || ctx.Expression() == nil {
		return c.front.TypeFacts.LoadConstant(runtime.None)
	}

	return c.front.Expressions.Compile(ctx.Expression())
}

func (c *DispatchCompiler) compileOptions(ctx fql.IDispatchOptionsClauseContext) bytecode.Operand {
	if ctx == nil || ctx.Expression() == nil {
		return c.front.TypeFacts.LoadConstant(runtime.None)
	}

	return c.front.Expressions.Compile(ctx.Expression())
}

func (c *DispatchCompiler) buildDispatchArgs(payload, options bytecode.Operand) bytecode.Operand {
	dst := c.ctx.Registers.Allocate()
	payloadKey := c.ctx.Symbols.AddConstant(runtime.NewString("payload"))
	optionsKey := c.ctx.Symbols.AddConstant(runtime.NewString("options"))

	c.ctx.Emitter.EmitObject(dst, 2)
	c.ctx.Emitter.EmitObjectSetConst(dst, payloadKey, payload)
	c.ctx.Emitter.EmitObjectSetConst(dst, optionsKey, options)
	c.ctx.Types.Set(dst, core.TypeObject)

	return dst
}

func (c *DispatchCompiler) ensureRegister(op bytecode.Operand) bytecode.Operand {
	if op == bytecode.NoopOperand {
		return c.front.TypeFacts.LoadConstant(runtime.None)
	}

	if op.IsRegister() {
		return op
	}

	dst := c.ctx.Registers.Allocate()
	c.ctx.Emitter.EmitLoadConst(dst, op)
	c.ctx.Types.Set(dst, c.front.TypeFacts.OperandType(op))

	return dst
}

func dispatchSpan(ctx fql.IDispatchExpressionContext) source.Span {
	if ctx == nil {
		return source.Span{Start: -1, End: -1}
	}

	if prc, ok := ctx.(antlr.ParserRuleContext); ok {
		return diagnostics.SpanFromRuleContext(prc)
	}

	return source.Span{Start: -1, End: -1}
}
