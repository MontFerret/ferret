package internal

import (
	"github.com/antlr4-go/antlr/v4"

	"github.com/MontFerret/ferret/v2/pkg/compiler/internal/core"
	"github.com/MontFerret/ferret/v2/pkg/compiler/internal/diagnostics"
	"github.com/MontFerret/ferret/v2/pkg/file"
	"github.com/MontFerret/ferret/v2/pkg/parser/fql"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/vm"
)

type DispatchCompiler struct {
	ctx *CompilerContext
}

func NewDispatchCompiler(ctx *CompilerContext) *DispatchCompiler {
	return &DispatchCompiler{
		ctx: ctx,
	}
}

func (c *DispatchCompiler) Compile(ctx fql.IDispatchExpressionContext) vm.Operand {
	if ctx == nil {
		return vm.NoopOperand
	}

	targetReg := c.ensureRegister(c.compileTarget(ctx.DispatchTarget()))
	eventReg := c.ensureRegister(c.compileEventName(ctx.DispatchEventName()))
	payloadReg := c.ensureRegister(c.compilePayload(ctx.DispatchWithClause()))
	optionsReg := c.ensureRegister(c.compileOptions(ctx.DispatchOptionsClause()))
	argsReg := c.buildDispatchArgs(payloadReg, optionsReg)

	dst := c.ctx.Registers.Allocate()
	span := dispatchSpan(ctx)

	c.ctx.Emitter.WithSpan(span, func() {
		c.ctx.Emitter.EmitMove(dst, targetReg)
		c.ctx.Emitter.EmitABC(vm.OpDispatch, dst, eventReg, argsReg)
	})

	c.ctx.Types.Set(dst, core.TypeAny)

	return dst
}

func (c *DispatchCompiler) compileEventName(ctx fql.IDispatchEventNameContext) vm.Operand {
	if ctx == nil {
		return vm.NoopOperand
	}

	if sl := ctx.StringLiteral(); sl != nil {
		return c.ctx.LiteralCompiler.CompileStringLiteral(sl)
	}

	if v := ctx.Variable(); v != nil {
		return c.ctx.ExprCompiler.CompileVariable(v)
	}

	if p := ctx.Param(); p != nil {
		return c.ctx.ExprCompiler.CompileParam(p)
	}

	if me := ctx.MemberExpression(); me != nil {
		return c.ctx.ExprCompiler.CompileMemberExpression(me)
	}

	if fc := ctx.FunctionCall(); fc != nil {
		return c.ctx.ExprCompiler.CompileFunctionCall(fc, false)
	}

	return vm.NoopOperand
}

func (c *DispatchCompiler) compileTarget(ctx fql.IDispatchTargetContext) vm.Operand {
	if ctx == nil {
		return vm.NoopOperand
	}

	if v := ctx.Variable(); v != nil {
		return c.ctx.ExprCompiler.CompileVariable(v)
	}

	if p := ctx.Param(); p != nil {
		return c.ctx.ExprCompiler.CompileParam(p)
	}

	if me := ctx.MemberExpression(); me != nil {
		return c.ctx.ExprCompiler.CompileMemberExpression(me)
	}

	if fc := ctx.FunctionCallExpression(); fc != nil {
		return c.ctx.ExprCompiler.CompileFunctionCallExpression(fc)
	}

	return vm.NoopOperand
}

func (c *DispatchCompiler) compilePayload(ctx fql.IDispatchWithClauseContext) vm.Operand {
	if ctx == nil || ctx.Expression() == nil {
		return loadConstant(c.ctx, runtime.None)
	}

	return c.ctx.ExprCompiler.Compile(ctx.Expression())
}

func (c *DispatchCompiler) compileOptions(ctx fql.IDispatchOptionsClauseContext) vm.Operand {
	if ctx == nil || ctx.Expression() == nil {
		return loadConstant(c.ctx, runtime.None)
	}

	return c.ctx.ExprCompiler.Compile(ctx.Expression())
}

func (c *DispatchCompiler) buildDispatchArgs(payload, options vm.Operand) vm.Operand {
	dst := c.ctx.Registers.Allocate()
	payloadKey := c.ctx.Symbols.AddConstant(runtime.NewString("payload"))
	optionsKey := c.ctx.Symbols.AddConstant(runtime.NewString("options"))

	c.ctx.Emitter.EmitObject(dst, 2)
	c.ctx.Emitter.EmitObjectSetConst(dst, payloadKey, payload)
	c.ctx.Emitter.EmitObjectSetConst(dst, optionsKey, options)
	c.ctx.Types.Set(dst, core.TypeObject)

	return dst
}

func (c *DispatchCompiler) ensureRegister(op vm.Operand) vm.Operand {
	if op == vm.NoopOperand {
		return loadConstant(c.ctx, runtime.None)
	}

	if op.IsRegister() {
		return op
	}

	dst := c.ctx.Registers.Allocate()
	c.ctx.Emitter.EmitLoadConst(dst, op)
	c.ctx.Types.Set(dst, operandType(c.ctx, op))

	return dst
}

func dispatchSpan(ctx fql.IDispatchExpressionContext) file.Span {
	if ctx == nil {
		return file.Span{Start: -1, End: -1}
	}

	if prc, ok := ctx.(antlr.ParserRuleContext); ok {
		return diagnostics.SpanFromRuleContext(prc)
	}

	return file.Span{Start: -1, End: -1}
}
