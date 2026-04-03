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

func (c *ExprCompiler) compileQueryExpression(ctx fql.IQueryExpressionContext) bytecode.Operand {
	if ctx == nil {
		return bytecode.NoopOperand
	}

	return c.front.Recovery.CompileOperation(OperationRecoverySpec{
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
			queryResult := c.emitApplyQuery(span, src, queryReg)
			dst := c.lowerQueryModifier(span, modifier, queryResult)

			if dst.IsRegister() {
				c.ctx.Types.Set(dst, queryResultTypeForModifier(modifier))
			}

			return dst
		},
	})
}

func (c *ExprCompiler) compileQueryExpressionSource(ctx fql.IQueryExpressionContext) (bytecode.Operand, bool) {
	if ctx == nil {
		return bytecode.NoopOperand, false
	}

	sourceExpr := ctx.Expression()
	if sourceExpr == nil {
		return bytecode.NoopOperand, false
	}

	source := c.Compile(sourceExpr)
	return source, source != bytecode.NoopOperand
}

func (c *ExprCompiler) emitQueryEnvelope(ctx fql.IQueryExpressionContext, span source.Span) bytecode.Operand {
	queryReg := c.ctx.Registers.Allocate()

	c.ctx.Emitter.WithSpan(span, func() {
		c.ctx.Emitter.EmitArray(queryReg, 3)
	})

	kind := c.compileQueryKindOperand(ctx)
	c.emitQueryEnvelopeOperand(span, queryReg, kind)

	payload := c.compileQueryPayloadOperand(ctx.QueryPayload())
	c.emitQueryEnvelopeOperand(span, queryReg, payload)

	options := c.compileQueryOptionsOperand(ctx.QueryWithOpt())
	c.emitQueryEnvelopeOperand(span, queryReg, options)

	return queryReg
}

func (c *ExprCompiler) compileQueryKindOperand(ctx fql.IQueryExpressionContext) bytecode.Operand {
	kind := ""
	if ident := ctx.GetDialect(); ident != nil {
		kind = strings.ToLower(ident.GetText())
	}

	return c.front.TypeFacts.LoadConstant(runtime.NewString(kind))
}

func (c *ExprCompiler) compileQueryPayloadOperand(ctx fql.IQueryPayloadContext) bytecode.Operand {
	if ctx == nil {
		return c.front.TypeFacts.LoadConstant(runtime.EmptyString)
	}

	if literal := ctx.StringLiteral(); literal != nil {
		if value, ok := parseStringLiteralConst(literal); ok {
			return c.front.TypeFacts.LoadConstant(value)
		}

		return c.front.Literals.CompileStringLiteral(literal)
	}

	if param := ctx.Param(); param != nil {
		return c.CompileParam(param)
	}

	if variable := ctx.Variable(); variable != nil {
		return c.CompileVariable(variable)
	}

	return c.front.TypeFacts.LoadConstant(runtime.EmptyString)
}

func (c *ExprCompiler) compileQueryOptionsOperand(ctx fql.IQueryWithOptContext) bytecode.Operand {
	if ctx == nil || ctx.Expression() == nil {
		return c.front.TypeFacts.LoadConstant(runtime.None)
	}

	return c.Compile(ctx.Expression())
}

func (c *ExprCompiler) emitQueryEnvelopeOperand(span source.Span, queryReg, value bytecode.Operand) {
	c.ctx.Emitter.WithSpan(span, func() {
		c.ctx.Emitter.EmitArrayPush(queryReg, value)
	})
}

func (c *ExprCompiler) emitApplyQuery(span source.Span, src, queryReg bytecode.Operand) bytecode.Operand {
	result := c.ctx.Registers.Allocate()

	c.ctx.Emitter.WithSpan(span, func() {
		c.ctx.Emitter.EmitABC(bytecode.OpQuery, result, src, queryReg)
	})

	return result
}

func (c *ExprCompiler) lowerQueryModifier(span source.Span, modifier queryModifier, queryResult bytecode.Operand) bytecode.Operand {
	switch modifier {
	case queryModifierExists:
		dst := c.ctx.Registers.Allocate()
		c.ctx.Emitter.WithSpan(span, func() {
			c.ctx.Emitter.EmitAB(bytecode.OpExists, dst, queryResult)
		})
		return dst
	case queryModifierCount:
		dst := c.ctx.Registers.Allocate()
		c.ctx.Emitter.WithSpan(span, func() {
			c.ctx.Emitter.EmitAB(bytecode.OpLength, dst, queryResult)
		})
		return dst
	case queryModifierAny:
		dst := c.ctx.Registers.Allocate()
		zero := c.ctx.Symbols.AddConstant(runtime.NewInt(0))
		c.ctx.Emitter.WithSpan(span, func() {
			c.ctx.Emitter.EmitABC(bytecode.OpLoadIndexOptionalConst, dst, queryResult, zero)
		})
		return dst
	case queryModifierValue:
		return c.lowerQueryModifierValue(span, queryResult)
	case queryModifierOne:
		return c.lowerQueryModifierOne(span, queryResult)
	default:
		return queryResult
	}
}

func queryModifierName(ctx fql.IQueryModifierContext) queryModifier {
	if ctx == nil {
		return queryModifierUnknown
	}

	return parseQueryModifier(ctx.GetText())
}

func queryResultTypeForModifier(modifier queryModifier) core.ValueType {
	switch modifier {
	case queryModifierExists:
		return core.TypeBool
	case queryModifierCount:
		return core.TypeInt
	case queryModifierAny, queryModifierValue, queryModifierOne:
		return core.TypeAny
	default:
		return core.TypeList
	}
}

func parseQueryModifier(text string) queryModifier {
	switch strings.ToLower(text) {
	case string(queryModifierExists):
		return queryModifierExists
	case string(queryModifierCount):
		return queryModifierCount
	case string(queryModifierAny):
		return queryModifierAny
	case string(queryModifierValue):
		return queryModifierValue
	case string(queryModifierOne):
		return queryModifierOne
	default:
		return queryModifierUnknown
	}
}

func (c *ExprCompiler) lowerQueryModifierValue(span source.Span, queryResult bytecode.Operand) bytecode.Operand {
	dst := c.ctx.Registers.Allocate()
	cond := c.ctx.Registers.Allocate()
	zero := c.ctx.Symbols.AddConstant(runtime.NewInt(0))
	message := c.ctx.Symbols.AddConstant(runtime.NewString(queryValueFailMessage))
	success := c.ctx.Emitter.NewLabel("query", string(queryModifierValue), "ok")
	end := c.ctx.Emitter.NewLabel("query", string(queryModifierValue), "end")

	c.ctx.Emitter.WithSpan(span, func() {
		c.ctx.Emitter.EmitAB(bytecode.OpExists, cond, queryResult)
		c.ctx.Emitter.EmitJumpIfTrue(cond, success)
		c.ctx.Emitter.EmitLoadNone(dst)
		c.ctx.Emitter.EmitA(bytecode.OpFail, message)
		c.ctx.Emitter.EmitJump(end)
		c.ctx.Emitter.MarkLabel(success)
		c.ctx.Emitter.EmitABC(bytecode.OpLoadIndexConst, dst, queryResult, zero)
		c.ctx.Emitter.MarkLabel(end)
	})

	return dst
}

func (c *ExprCompiler) lowerQueryModifierOne(span source.Span, queryResult bytecode.Operand) bytecode.Operand {
	dst := c.ctx.Registers.Allocate()
	length := c.ctx.Registers.Allocate()
	one := c.ctx.Symbols.AddConstant(runtime.NewInt(1))
	zero := c.ctx.Symbols.AddConstant(runtime.NewInt(0))
	message := c.ctx.Symbols.AddConstant(runtime.NewString(queryOneFailMessage))
	success := c.ctx.Emitter.NewLabel("query", string(queryModifierOne), "ok")
	end := c.ctx.Emitter.NewLabel("query", string(queryModifierOne), "end")

	c.ctx.Emitter.WithSpan(span, func() {
		c.ctx.Emitter.EmitAB(bytecode.OpLength, length, queryResult)
		c.ctx.Emitter.EmitJumpCompare(bytecode.OpJumpIfEqConst, length, one, success)
		c.ctx.Emitter.EmitLoadNone(dst)
		c.ctx.Emitter.EmitA(bytecode.OpFail, message)
		c.ctx.Emitter.EmitJump(end)
		c.ctx.Emitter.MarkLabel(success)
		c.ctx.Emitter.EmitABC(bytecode.OpLoadIndexConst, dst, queryResult, zero)
		c.ctx.Emitter.MarkLabel(end)
	})

	return dst
}

func (c *ExprCompiler) compileQueryLiteral(ctx fql.IQueryLiteralContext) bytecode.Operand {
	if ctx == nil {
		return bytecode.NoopOperand
	}

	kind := ""
	if ident := ctx.Identifier(); ident != nil {
		kind = strings.ToLower(ident.GetText())
	}

	dst := c.ctx.Registers.Allocate()
	span := diagnostics.SpanFromRuleContext(ctx)

	c.ctx.Emitter.WithSpan(span, func() {
		c.ctx.Emitter.EmitArray(dst, 3)
	})

	kindReg := c.front.TypeFacts.LoadConstant(runtime.NewString(kind))

	c.ctx.Emitter.WithSpan(span, func() {
		c.ctx.Emitter.EmitArrayPush(dst, kindReg)
	})

	payloadReg := c.front.TypeFacts.LoadConstant(runtime.EmptyString)
	if str := ctx.StringLiteral(); str != nil {
		if val, ok := parseStringLiteralConst(str); ok {
			payloadReg = c.front.TypeFacts.LoadConstant(val)
		} else {
			payloadReg = c.front.Literals.CompileStringLiteral(str)
		}
	}

	c.ctx.Emitter.WithSpan(span, func() {
		c.ctx.Emitter.EmitArrayPush(dst, payloadReg)
	})

	params := ctx.Expression()
	var paramsReg bytecode.Operand

	if params == nil {
		paramsReg = c.front.TypeFacts.LoadConstant(runtime.None)
	} else {
		paramsReg = c.Compile(params)
	}

	c.ctx.Emitter.WithSpan(span, func() {
		c.ctx.Emitter.EmitArrayPush(dst, paramsReg)
	})

	if dst.IsRegister() {
		c.ctx.Types.Set(dst, core.TypeAny)
	}

	return dst
}
