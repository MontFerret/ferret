package internal

import (
	"strings"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/compiler/internal/core"
	"github.com/MontFerret/ferret/v2/pkg/parser/diagnostics"
	"github.com/MontFerret/ferret/v2/pkg/parser/fql"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/source"
)

type (
	exprQueryCallbacks struct {
		compileExpr     func(fql.IExpressionContext) bytecode.Operand
		compileParam    func(fql.IParamContext) bytecode.Operand
		compileVariable func(fql.IVariableContext) bytecode.Operand
	}

	exprQueryCompiler struct {
		ctx       *CompilationSession
		literals  *LiteralCompiler
		recovery  *RecoveryCompiler
		facts     *TypeFacts
		callbacks exprQueryCallbacks
	}
)

func newExprQueryCompiler(ctx *CompilationSession, callbacks exprQueryCallbacks) *exprQueryCompiler {
	return &exprQueryCompiler{
		ctx:       ctx,
		callbacks: callbacks,
	}
}

func (c *exprQueryCompiler) bind(literals *LiteralCompiler, recovery *RecoveryCompiler, facts *TypeFacts) {
	if c == nil {
		return
	}

	c.literals = literals
	c.recovery = recovery
	c.facts = facts
}

func (c *exprQueryCompiler) compileQueryExpression(ctx fql.IQueryExpressionContext) bytecode.Operand {
	if ctx == nil {
		return bytecode.NoopOperand
	}

	return c.recovery.CompileOperation(OperationRecoverySpec{
		Owner:    ctx,
		JumpMode: core.CatchJumpModeNone,
		CompilePlain: func() bytecode.Operand {
			if ctx == nil {
				return bytecode.NoopOperand
			}

			src, ok := c.compileQueryExpressionSource(ctx)
			if !ok {
				return bytecode.NoopOperand
			}

			span := diagnostics.SpanFromRuleContext(ctx)
			modifier := queryModifierName(ctx.QueryModifier())
			queryReg := c.emitQueryEnvelope(ctx, span)
			dst := c.emitApplyQuery(span, src, queryReg, queryOpcodeForModifier(modifier))

			if dst.IsRegister() {
				c.ctx.Function.Types.Set(dst, queryResultTypeForModifier(modifier))
			}

			return dst
		},
	})
}

func (c *exprQueryCompiler) compileQueryExpressionSource(ctx fql.IQueryExpressionContext) (bytecode.Operand, bool) {
	if ctx == nil {
		return bytecode.NoopOperand, false
	}

	sourceExpr := ctx.Expression()
	if sourceExpr == nil {
		return bytecode.NoopOperand, false
	}

	source := c.callbacks.compileExpr(sourceExpr)

	return source, source != bytecode.NoopOperand
}

func (c *exprQueryCompiler) emitQueryEnvelope(ctx fql.IQueryExpressionContext, span source.Span) bytecode.Operand {
	queryReg := c.ctx.Function.Registers.Allocate()

	c.ctx.Program.Emitter.WithSpan(span, func() {
		c.ctx.Program.Emitter.EmitArray(queryReg, 4)
	})

	kind := c.compileQueryKindOperand(ctx)
	c.emitQueryEnvelopeOperand(span, queryReg, kind)

	expression := c.compileQueryExpressionOperand(ctx.QueryPayload())
	c.emitQueryEnvelopeOperand(span, queryReg, expression)

	params := c.compileQueryParamsOperand(ctx.QueryWithOpt())
	c.emitQueryEnvelopeOperand(span, queryReg, params)

	options := c.compileQueryOptionsOperand(ctx.QueryOptionsOpt())
	c.emitQueryEnvelopeOperand(span, queryReg, options)

	return queryReg
}

func (c *exprQueryCompiler) compileQueryKindOperand(ctx fql.IQueryExpressionContext) bytecode.Operand {
	kind := ""
	if ident := ctx.GetDialect(); ident != nil {
		kind = strings.ToLower(ident.GetText())
	}

	return c.facts.LoadConstant(runtime.NewString(kind))
}

func (c *exprQueryCompiler) compileQueryExpressionOperand(ctx fql.IQueryPayloadContext) bytecode.Operand {
	if ctx == nil {
		return c.facts.LoadConstant(runtime.EmptyString)
	}

	if literal := ctx.StringLiteral(); literal != nil {
		if value, ok := parseStringLiteralConst(literal); ok {
			return c.facts.LoadConstant(value)
		}

		return c.literals.CompileStringLiteral(literal)
	}

	if param := ctx.Param(); param != nil {
		return c.callbacks.compileParam(param)
	}

	if variable := ctx.Variable(); variable != nil {
		return c.callbacks.compileVariable(variable)
	}

	return c.facts.LoadConstant(runtime.EmptyString)
}

func (c *exprQueryCompiler) compileQueryParamsOperand(ctx fql.IQueryWithOptContext) bytecode.Operand {
	if ctx == nil || ctx.Expression() == nil {
		return c.facts.LoadConstant(runtime.None)
	}

	return c.callbacks.compileExpr(ctx.Expression())
}

func (c *exprQueryCompiler) compileQueryOptionsOperand(ctx fql.IQueryOptionsOptContext) bytecode.Operand {
	if ctx == nil || ctx.Expression() == nil {
		return c.facts.LoadConstant(runtime.None)
	}

	return c.callbacks.compileExpr(ctx.Expression())
}

func (c *exprQueryCompiler) emitQueryEnvelopeOperand(span source.Span, queryReg, value bytecode.Operand) {
	c.ctx.Program.Emitter.WithSpan(span, func() {
		c.ctx.Program.Emitter.EmitArrayPush(queryReg, value)
	})
}

func (c *exprQueryCompiler) emitApplyQuery(span source.Span, src, queryReg bytecode.Operand, opcode bytecode.Opcode) bytecode.Operand {
	result := c.ctx.Function.Registers.Allocate()

	c.ctx.Program.Emitter.WithSpan(span, func() {
		c.ctx.Program.Emitter.EmitABC(opcode, result, src, queryReg)
	})

	return result
}

func (c *exprQueryCompiler) compileQueryLiteral(ctx fql.IQueryLiteralContext) bytecode.Operand {
	if ctx == nil {
		return bytecode.NoopOperand
	}

	kind := ""
	if ident := ctx.Identifier(); ident != nil {
		kind = strings.ToLower(ident.GetText())
	}

	dst := c.ctx.Function.Registers.Allocate()
	span := diagnostics.SpanFromRuleContext(ctx)

	c.ctx.Program.Emitter.WithSpan(span, func() {
		c.ctx.Program.Emitter.EmitArray(dst, 4)
	})

	kindReg := c.facts.LoadConstant(runtime.NewString(kind))

	c.ctx.Program.Emitter.WithSpan(span, func() {
		c.ctx.Program.Emitter.EmitArrayPush(dst, kindReg)
	})

	expressionReg := c.facts.LoadConstant(runtime.EmptyString)
	if str := ctx.StringLiteral(); str != nil {
		if val, ok := parseStringLiteralConst(str); ok {
			expressionReg = c.facts.LoadConstant(val)
		} else {
			expressionReg = c.literals.CompileStringLiteral(str)
		}
	}

	c.ctx.Program.Emitter.WithSpan(span, func() {
		c.ctx.Program.Emitter.EmitArrayPush(dst, expressionReg)
	})

	params := ctx.Expression()
	var paramsReg bytecode.Operand

	if params == nil {
		paramsReg = c.facts.LoadConstant(runtime.None)
	} else {
		paramsReg = c.callbacks.compileExpr(params)
	}

	optionsReg := c.facts.LoadConstant(runtime.None)

	c.ctx.Program.Emitter.WithSpan(span, func() {
		c.ctx.Program.Emitter.EmitArrayPush(dst, paramsReg)
		c.ctx.Program.Emitter.EmitArrayPush(dst, optionsReg)
	})

	if dst.IsRegister() {
		c.ctx.Function.Types.Set(dst, core.TypeAny)
	}

	return dst
}
