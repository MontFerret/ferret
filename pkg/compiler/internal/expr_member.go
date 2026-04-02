package internal

import (
	"github.com/antlr4-go/antlr/v4"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/compiler/internal/core"
	"github.com/MontFerret/ferret/v2/pkg/parser/diagnostics"
	"github.com/MontFerret/ferret/v2/pkg/parser/fql"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/source"
)

// CompileImplicitMemberExpression processes an implicit member expression (e.g., .name, .[0], ?.name).
// The implicit current resolution is centralized in resolveImplicitCurrent.
func (c *ExprCompiler) CompileImplicitMemberExpression(ctx fql.IImplicitMemberExpressionContext) bytecode.Operand {
	start := ctx.ImplicitMemberExpressionStart()
	if start == nil {
		return bytecode.NoopOperand
	}

	src, ok := c.resolveImplicitMemberSource(start)
	if !ok {
		return bytecode.NoopOperand
	}

	segments := ctx.AllMemberExpressionPath()
	if arrayResult, handled := c.compileImplicitMemberArrayOperator(src, start, segments); handled {
		return arrayResult
	}

	return c.compileImplicitMemberPath(src, start, segments)
}

func (c *ExprCompiler) resolveImplicitMemberSource(start fql.IImplicitMemberExpressionStartContext) (bytecode.Operand, bool) {
	if start == nil {
		return bytecode.NoopOperand, false
	}

	if c.implicitCurrentDepth == 0 {
		c.resolveImplicitCurrent(getImplicitToken(start))
		return bytecode.NoopOperand, false
	}

	return c.resolveImplicitCurrent(getImplicitToken(start))
}

func (c *ExprCompiler) compileImplicitMemberArrayOperator(src bytecode.Operand, start fql.IImplicitMemberExpressionStartContext, segments []fql.IMemberExpressionPathContext) (bytecode.Operand, bool) {
	if expansion := start.ArrayExpansion(); expansion != nil {
		return c.compileArrayExpansionChain(src, expansion, segments), true
	}

	if contraction := start.ArrayContraction(); contraction != nil {
		inlineTail, restTail := splitArrayOperatorTail(segments)
		result := c.compileArrayContraction(src, contraction, inlineTail)
		return c.continueImplicitMemberArrayResult(result, restTail), true
	}

	if question := start.ArrayQuestionMark(); question != nil {
		inlineTail, restTail := splitArrayOperatorTail(segments)
		result := c.compileArrayQuestionMark(src, question, inlineTail)
		return c.continueImplicitMemberArrayResult(result, restTail), true
	}

	if apply := start.ArrayApply(); apply != nil {
		return c.compileArrayApply(src, apply, segments), true
	}

	return bytecode.NoopOperand, false
}

func (c *ExprCompiler) continueImplicitMemberArrayResult(result bytecode.Operand, tail []fql.IMemberExpressionPathContext) bytecode.Operand {
	if len(tail) == 0 {
		return result
	}

	return c.compileMemberExpressionSegments(result, tail)
}

func (c *ExprCompiler) compileImplicitMemberPath(src bytecode.Operand, start fql.IImplicitMemberExpressionStartContext, segments []fql.IMemberExpressionPathContext) bytecode.Operand {
	if isSimpleMemberPathChain(segments) {
		return c.compileImplicitSimpleMemberExpressionSegments(src, start, segments)
	}

	return c.compileImplicitGenericMemberExpression(src, start, segments)
}

func (c *ExprCompiler) compileImplicitGenericMemberExpression(src bytecode.Operand, start fql.IImplicitMemberExpressionStartContext, segments []fql.IMemberExpressionPathContext) bytecode.Operand {
	dst, ok := c.emitImplicitMemberStartLoad(src, start)
	if !ok {
		return bytecode.NoopOperand
	}

	if len(segments) == 0 {
		return dst
	}

	return c.compileMemberExpressionSegments(dst, segments)
}

func (c *ExprCompiler) emitImplicitMemberStartLoad(src bytecode.Operand, start fql.IImplicitMemberExpressionStartContext) (bytecode.Operand, bool) {
	operand, constOperand := c.compileImplicitMemberStartOperand(start)
	if operand == bytecode.NoopOperand {
		return bytecode.NoopOperand, false
	}

	dst := c.ctx.Registers.Allocate()
	span := diagnostics.SpanFromRuleContext(start.(antlr.ParserRuleContext))
	optional := start.ErrorOperator() != nil

	c.ctx.Emitter.WithSpan(span, func() {
		op := memberLoadOpcode(operandType(c.ctx, src), constOperand, optional)
		c.ctx.Emitter.EmitABC(op, dst, src, operand)
	})

	return dst, true
}

func (c *ExprCompiler) compileWithImplicitCurrent(expr fql.IExpressionContext) bytecode.Operand {
	if expr == nil {
		return bytecode.NoopOperand
	}

	c.implicitCurrentDepth++
	defer func() {
		c.implicitCurrentDepth--
	}()

	return c.Compile(expr)
}

func (c *ExprCompiler) withImplicitCurrent(fn func()) {
	c.implicitCurrentDepth++
	defer func() {
		c.implicitCurrentDepth--
	}()

	fn()
}

func (c *ExprCompiler) CompileMemberExpression(ctx fql.IMemberExpressionContext) bytecode.Operand {
	if ctx == nil {
		return bytecode.NoopOperand
	}

	plan := collectRecoveryPlan(c.ctx, ctx, core.RecoveryPlanOptions{})
	return c.ctx.PolicyCompiler.CompileWithRecoveryPlan(plan, core.CatchJumpModeNone, func() bytecode.Operand {
		mes := ctx.MemberExpressionSource()
		segments := ctx.AllMemberExpressionPath()

		if mes == nil || len(segments) == 0 {
			return bytecode.NoopOperand
		}

		src := c.compileMemberExpressionSource(mes, segments)

		if src == bytecode.NoopOperand {
			return src
		}

		return c.compileMemberExpressionSegments(src, segments)
	})
}

func (c *ExprCompiler) compileMemberExpressionSource(mes fql.IMemberExpressionSourceContext, segments []fql.IMemberExpressionPathContext) bytecode.Operand {
	if v := mes.Variable(); v != nil {
		return c.CompileVariable(v)
	}

	if p := mes.Param(); p != nil {
		return c.CompileParam(p)
	}

	if ol := mes.ObjectLiteral(); ol != nil {
		return c.ctx.LiteralCompiler.CompileObjectLiteral(ol)
	}

	if al := mes.ArrayLiteral(); al != nil {
		return c.ctx.LiteralCompiler.CompileArrayLiteral(al)
	}

	if fc := mes.FunctionCall(); fc != nil {
		segment := segments[0]
		return c.CompileFunctionCall(fc, segment.ErrorOperator() != nil)
	}

	if fe := mes.ForExpression(); fe != nil {
		return c.ctx.LoopCompiler.Compile(fe)
	}

	if wfe := mes.WaitForExpression(); wfe != nil {
		return c.ctx.WaitCompiler.Compile(wfe)
	}

	if e := mes.Expression(); e != nil {
		return c.Compile(e)
	}

	return bytecode.NoopOperand
}

// CompileImplicitCurrentExpression processes a bare implicit current shorthand (e.g., .).
func (c *ExprCompiler) CompileImplicitCurrentExpression(ctx fql.IImplicitCurrentExpressionContext) bytecode.Operand {
	if ctx == nil {
		return bytecode.NoopOperand
	}

	src, ok := c.resolveImplicitCurrent(getImplicitToken(ctx))
	if !ok {
		return bytecode.NoopOperand
	}

	return src
}

// resolveImplicitCurrent centralizes implicit current resolution and error reporting.
// It returns a register operand and true on success.
func (c *ExprCompiler) resolveImplicitCurrent(token antlr.Token) (bytecode.Operand, bool) {
	if c.implicitCurrentDepth == 0 {
		c.ctx.Errors.VariableNotFound(token, core.PseudoVariable)
		return bytecode.NoopOperand, false
	}

	binding, found := c.ctx.Symbols.ResolveBinding(core.PseudoVariable)
	if !found {
		c.ctx.Errors.VariableNotFound(token, core.PseudoVariable)
		return bytecode.NoopOperand, false
	}

	src := c.ctx.BindingCompiler.LoadBindingValue(binding)
	src = c.ensureRegister(src)

	return src, true
}

// getImplicitToken picks the span anchor for implicit-current diagnostics.
func getImplicitToken(ctx antlr.ParserRuleContext) antlr.Token {
	switch v := ctx.(type) {
	case fql.IImplicitMemberExpressionStartContext:
		if dot := v.Dot(); dot != nil {
			return dot.GetSymbol()
		}
	case fql.IImplicitCurrentExpressionContext:
		if dot := v.Dot(); dot != nil {
			return dot.GetSymbol()
		}
	}

	return ctx.GetStart()
}

func (c *ExprCompiler) compileMemberExpressionSegments(src bytecode.Operand, segments []fql.IMemberExpressionPathContext) bytecode.Operand {
	if len(segments) == 0 {
		return src
	}

	if isSimpleMemberPathChain(segments) {
		return c.compileSimpleMemberExpressionSegments(src, segments)
	}

	for idx, segment := range segments {
		p := segment.(*fql.MemberExpressionPathContext)

		if contraction := p.ArrayContraction(); contraction != nil {
			inlineTail, restTail := splitArrayOperatorTail(segments[idx+1:])
			result := c.compileArrayContraction(src, contraction, inlineTail)

			if len(restTail) == 0 {
				return result
			}

			return c.compileMemberExpressionSegments(result, restTail)
		}

		if expansion := p.ArrayExpansion(); expansion != nil {
			return c.compileArrayExpansionChain(src, expansion, segments[idx+1:])
		}

		if question := p.ArrayQuestionMark(); question != nil {
			inlineTail, restTail := splitArrayOperatorTail(segments[idx+1:])
			result := c.compileArrayQuestionMark(src, question, inlineTail)

			if len(restTail) == 0 {
				return result
			}

			return c.compileMemberExpressionSegments(result, restTail)
		}

		if apply := p.ArrayApply(); apply != nil {
			return c.compileArrayApply(src, apply, segments[idx+1:])
		}

		src2, constOperand := c.compileMemberPathOperand(p)
		dst := c.ctx.Registers.Allocate()
		span := diagnostics.SpanFromRuleContext(p)

		c.ctx.Emitter.WithSpan(span, func() {
			optional := p.ErrorOperator() != nil
			op := memberLoadOpcode(operandType(c.ctx, src), constOperand, optional)

			c.ctx.Emitter.EmitABC(op, dst, src, src2)
		})

		src = dst
	}

	return src
}

func (c *ExprCompiler) compileImplicitSimpleMemberExpressionSegments(src bytecode.Operand, start fql.IImplicitMemberExpressionStartContext, segments []fql.IMemberExpressionPathContext) bytecode.Operand {
	startOp, startConst := c.compileImplicitMemberStartOperand(start)
	if startOp == bytecode.NoopOperand {
		return bytecode.NoopOperand
	}

	state := &optionalMemberChainState{}
	startSpan := diagnostics.SpanFromRuleContext(start.(antlr.ParserRuleContext))
	result := c.emitOptionalMemberLoadSegment(startSpan, src, startOp, startConst, start.ErrorOperator() != nil, state)

	for _, segment := range segments {
		p := segment.(*fql.MemberExpressionPathContext)
		segmentOp, constOperand := c.compileMemberPathOperand(p)
		span := diagnostics.SpanFromRuleContext(p)
		result = c.emitOptionalMemberLoadSegment(span, result, segmentOp, constOperand, p.ErrorOperator() != nil, state)
	}

	c.finalizeOptionalMemberChain(state)
	return result
}

func (c *ExprCompiler) emitOptionalMemberLoadSegment(span source.Span, src, segmentOp bytecode.Operand, constOperand, optional bool, state *optionalMemberChainState) bytecode.Operand {
	dst := c.allocateOptionalMemberDestination(src, state)

	c.ctx.Emitter.WithSpan(span, func() {
		op := memberLoadOpcode(operandType(c.ctx, src), constOperand, optional)
		c.ctx.Emitter.EmitABC(op, dst, src, segmentOp)

		if optional {
			c.ctx.Emitter.EmitJumpIfNone(dst, c.optionalMemberEndLabel(state))
		}
	})

	if optional {
		state.stickyDst = true
	}

	return dst
}

func (c *ExprCompiler) allocateOptionalMemberDestination(src bytecode.Operand, state *optionalMemberChainState) bytecode.Operand {
	if state != nil && state.stickyDst && src.IsRegister() {
		return src
	}

	return c.ctx.Registers.Allocate()
}

func (c *ExprCompiler) optionalMemberEndLabel(state *optionalMemberChainState) core.Label {
	if state == nil {
		return core.Label{}
	}

	if state.hasJump {
		return state.endLabel
	}

	state.endLabel = c.ctx.Emitter.NewLabel("member", "optional", "end")
	state.hasJump = true

	return state.endLabel
}

func (c *ExprCompiler) finalizeOptionalMemberChain(state *optionalMemberChainState) {
	if state == nil || !state.hasJump {
		return
	}

	c.ctx.Emitter.MarkLabel(state.endLabel)
}

func isSimpleMemberPathChain(segments []fql.IMemberExpressionPathContext) bool {
	for _, segment := range segments {
		p := segment.(*fql.MemberExpressionPathContext)

		if p.PropertyName() == nil && p.ComputedPropertyName() == nil {
			return false
		}
	}

	return true
}

func (c *ExprCompiler) compileSimpleMemberExpressionSegments(src bytecode.Operand, segments []fql.IMemberExpressionPathContext) bytecode.Operand {
	result := src
	stickyDst := false
	hasJump := false
	var endLabel core.Label

	for _, segment := range segments {
		p := segment.(*fql.MemberExpressionPathContext)
		src2, constOperand := c.compileMemberPathOperand(p)
		optional := p.ErrorOperator() != nil

		dst := result
		if !stickyDst || !dst.IsRegister() {
			dst = c.ctx.Registers.Allocate()
		}

		span := diagnostics.SpanFromRuleContext(p)

		c.ctx.Emitter.WithSpan(span, func() {
			op := memberLoadOpcode(operandType(c.ctx, result), constOperand, optional)
			c.ctx.Emitter.EmitABC(op, dst, result, src2)

			if optional {
				if !hasJump {
					endLabel = c.ctx.Emitter.NewLabel("member", "optional", "end")
					hasJump = true
				}
				c.ctx.Emitter.EmitJumpIfNone(dst, endLabel)
			}
		})

		if optional {
			stickyDst = true
		}

		result = dst
	}

	if hasJump {
		c.ctx.Emitter.MarkLabel(endLabel)
	}

	return result
}

func (c *ExprCompiler) compileMemberPathOperand(p *fql.MemberExpressionPathContext) (bytecode.Operand, bool) {
	if pn := p.PropertyName(); pn != nil {
		if constOp, ok := c.ctx.LiteralCompiler.CompilePropertyNameConst(pn); ok {
			return constOp, true
		}

		return c.ctx.LiteralCompiler.CompilePropertyName(pn), false
	}

	if cpn := p.ComputedPropertyName(); cpn != nil {
		if val, ok := literalValueFromExpression(cpn.Expression()); ok {
			switch val.(type) {
			case *runtime.Array, *runtime.Object:
				return c.ctx.LiteralCompiler.CompileComputedPropertyName(cpn), false
			default:
				return c.ctx.Symbols.AddConstant(val), true
			}
		}

		return c.ctx.LiteralCompiler.CompileComputedPropertyName(cpn), false
	}

	return bytecode.NoopOperand, false
}

func (c *ExprCompiler) compileImplicitMemberStartOperand(start fql.IImplicitMemberExpressionStartContext) (bytecode.Operand, bool) {
	if pn := start.PropertyName(); pn != nil {
		if constOp, ok := c.ctx.LiteralCompiler.CompilePropertyNameConst(pn); ok {
			return constOp, true
		}

		return c.ctx.LiteralCompiler.CompilePropertyName(pn), false
	}

	if cpn := start.ComputedPropertyName(); cpn != nil {
		if val, ok := literalValueFromExpression(cpn.Expression()); ok {
			switch val.(type) {
			case *runtime.Array, *runtime.Object:
				return c.ctx.LiteralCompiler.CompileComputedPropertyName(cpn), false
			default:
				return c.ctx.Symbols.AddConstant(val), true
			}
		}

		return c.ctx.LiteralCompiler.CompileComputedPropertyName(cpn), false
	}

	return bytecode.NoopOperand, false
}

func memberLoadOpcode(srcType core.ValueType, constOperand, optional bool) bytecode.Opcode {
	switch srcType {
	case core.TypeArray:
		if constOperand {
			if optional {
				return bytecode.OpLoadIndexOptionalConst
			}

			return bytecode.OpLoadIndexConst
		}

		if optional {
			return bytecode.OpLoadIndexOptional
		}

		return bytecode.OpLoadIndex
	case core.TypeObject:
		if constOperand {
			if optional {
				return bytecode.OpLoadKeyOptionalConst
			}

			return bytecode.OpLoadKeyConst
		}

		if optional {
			return bytecode.OpLoadKeyOptional
		}

		return bytecode.OpLoadKey
	default:
		if constOperand {
			if optional {
				return bytecode.OpLoadPropertyOptionalConst
			}

			return bytecode.OpLoadPropertyConst
		}

		if optional {
			return bytecode.OpLoadPropertyOptional
		}

		return bytecode.OpLoadProperty
	}
}

func splitArrayOperatorTail(segments []fql.IMemberExpressionPathContext) ([]fql.IMemberExpressionPathContext, []fql.IMemberExpressionPathContext) {
	if len(segments) > 0 {
		p := segments[0].(*fql.MemberExpressionPathContext)

		if p.ArrayContraction() != nil || p.ArrayExpansion() != nil || p.ArrayQuestionMark() != nil {
			return nil, segments
		}
	}

	return segments, nil
}

// splitTerminalArrayContractionTail hoists only the final contraction segment.
// Earlier array operators stay in the per-element tail so existing projection semantics remain unchanged.
func splitTerminalArrayContractionTail(segments []fql.IMemberExpressionPathContext) ([]fql.IMemberExpressionPathContext, fql.IArrayContractionContext) {
	if len(segments) == 0 {
		return nil, nil
	}

	last := segments[len(segments)-1].(*fql.MemberExpressionPathContext)
	contraction := last.ArrayContraction()
	if contraction == nil {
		return segments, nil
	}

	return segments[:len(segments)-1], contraction
}

func (c *ExprCompiler) compileArrayExpansionChain(src bytecode.Operand, expansion fql.IArrayExpansionContext, tail []fql.IMemberExpressionPathContext) bytecode.Operand {
	inline := expansion.InlineExpression()

	if inline == nil {
		if next, rest := nextArrayExpansion(tail); next != nil {
			return c.compileArrayExpansionChain(src, next, rest)
		}

		return c.compileArrayExpansionChainWithFilters(src, expansion, tail, nil)
	}

	if !isFilterOnlyInline(inline) {
		tail = dropIdentityExpansions(tail)

		return c.compileArrayExpansionChainWithFilters(src, expansion, tail, nil)
	}

	extraFilters, rest := collectFilterOnlyTail(tail)

	return c.compileArrayExpansionChainWithFilters(src, expansion, rest, extraFilters)
}

func (c *ExprCompiler) compileArrayExpansionWithFilters(src bytecode.Operand, expansion fql.IArrayExpansionContext, tail []fql.IMemberExpressionPathContext, extraFilters []fql.IExpressionContext) bytecode.Operand {
	span := diagnostics.SpanFromRuleContext(expansion)
	inline := expansion.InlineExpression()

	return c.compileArrayIteration(src, span, tail, inline, extraFilters)
}

func (c *ExprCompiler) compileArrayExpansionChainWithFilters(src bytecode.Operand, expansion fql.IArrayExpansionContext, tail []fql.IMemberExpressionPathContext, extraFilters []fql.IExpressionContext) bytecode.Operand {
	inlineTail, restTail := splitArrayOperatorTail(tail)
	result := c.compileArrayExpansionWithFilters(src, expansion, inlineTail, extraFilters)

	if len(restTail) == 0 {
		return result
	}

	return c.compileMemberExpressionSegments(result, restTail)
}

func isFilterOnlyInline(inline fql.IInlineExpressionContext) bool {
	if inline == nil {
		return false
	}

	return inline.InlineFilter() != nil && inline.InlineLimit() == nil && inline.InlineReturn() == nil
}

func nextArrayExpansion(segments []fql.IMemberExpressionPathContext) (fql.IArrayExpansionContext, []fql.IMemberExpressionPathContext) {
	if len(segments) == 0 {
		return nil, segments
	}

	p := segments[0].(*fql.MemberExpressionPathContext)

	if expansion := p.ArrayExpansion(); expansion != nil {
		return expansion, segments[1:]
	}

	return nil, segments
}

func dropIdentityExpansions(segments []fql.IMemberExpressionPathContext) []fql.IMemberExpressionPathContext {
	for len(segments) > 0 {
		p := segments[0].(*fql.MemberExpressionPathContext)
		expansion := p.ArrayExpansion()

		if expansion == nil {
			break
		}

		if expansion.InlineExpression() != nil {
			break
		}

		segments = segments[1:]
	}

	return segments
}

func collectFilterOnlyTail(segments []fql.IMemberExpressionPathContext) ([]fql.IExpressionContext, []fql.IMemberExpressionPathContext) {
	extraFilters := make([]fql.IExpressionContext, 0)
	rest := segments

	for len(rest) > 0 {
		p := rest[0].(*fql.MemberExpressionPathContext)
		expansion := p.ArrayExpansion()
		if expansion == nil {
			break
		}

		inline := expansion.InlineExpression()
		if inline == nil {
			rest = rest[1:]
			continue
		}

		if !isFilterOnlyInline(inline) {
			break
		}

		filter := inline.InlineFilter()
		if filter != nil {
			extraFilters = append(extraFilters, filter.Expression())
		}

		rest = rest[1:]
	}

	return extraFilters, rest
}

func (c *ExprCompiler) compileArrayQuestionMark(src bytecode.Operand, question fql.IArrayQuestionMarkContext, tail []fql.IMemberExpressionPathContext) bytecode.Operand {
	span := diagnostics.SpanFromRuleContext(question)

	loop := &core.Loop{
		Kind:     core.ForInLoop,
		Type:     core.NormalLoop,
		Distinct: false,
		Allocate: false,
		Dst:      bytecode.NoopOperand,
		Src:      src,
	}

	c.ctx.Loops.Push(loop)
	c.ctx.Symbols.EnterScope()

	loop.DeclareValueVar(core.PseudoVariable, c.ctx.Symbols, core.TypeAny)

	if loop.Value.IsRegister() {
		c.ctx.Types.Set(loop.Value, core.TypeAny)
	}

	count := c.ctx.Registers.Allocate()
	total := c.ctx.Registers.Allocate()

	c.ctx.Emitter.WithSpan(span, func() {
		c.ctx.Emitter.EmitA(bytecode.OpLoadZero, count)
		c.ctx.Emitter.EmitA(bytecode.OpLoadZero, total)
		loop.EmitInitialization(c.ctx.Registers, c.ctx.Emitter)
	})

	c.ctx.Emitter.EmitA(bytecode.OpIncr, total)

	if filter := question.Expression(); filter != nil {
		cond := c.compileWithImplicitCurrent(filter)
		label := c.ctx.Loops.Current().ContinueLabel()
		c.ctx.Emitter.EmitJumpIfFalse(cond, label)
	}

	c.ctx.Emitter.EmitA(bytecode.OpIncr, count)

	loop.EmitFinalization(c.ctx.Emitter)

	c.ctx.Symbols.ExitScope()
	c.ctx.Loops.Pop()

	result := c.compileArrayQuestionQuantifier(question, count, total)

	if len(tail) > 0 {
		result = c.compileMemberExpressionSegments(result, tail)
	}

	if result.IsRegister() {
		c.ctx.Types.Set(result, core.TypeBool)
	}

	return result
}

func (c *ExprCompiler) compileArrayQuestionQuantifier(question fql.IArrayQuestionMarkContext, count, total bytecode.Operand) bytecode.Operand {
	quant := question.ArrayQuestionQuantifier()
	zero := c.ctx.Registers.Allocate()
	c.ctx.Emitter.EmitA(bytecode.OpLoadZero, zero)

	if quant == nil || quant.Any() != nil {
		return c.emitComparison(bytecode.OpGt, count, zero)
	}

	if quant.None() != nil {
		return c.emitComparison(bytecode.OpEq, count, zero)
	}

	if quant.All() != nil {
		return c.emitComparison(bytecode.OpEq, count, total)
	}

	values := quant.AllArrayQuestionQuantifierValue()

	if quant.At() != nil {
		if len(values) == 0 {
			return c.emitComparison(bytecode.OpGt, count, zero)
		}

		value := c.compileArrayQuestionQuantifierValue(values[0])

		return c.emitComparison(bytecode.OpGte, count, value)
	}

	if quant.Range() != nil && len(values) >= 2 {
		min := c.compileArrayQuestionQuantifierValue(values[0])
		max := c.compileArrayQuestionQuantifierValue(values[1])

		left := c.emitComparison(bytecode.OpGte, count, min)
		right := c.emitComparison(bytecode.OpLte, count, max)

		return c.emitBooleanAnd(left, right)
	}

	if len(values) > 0 {
		value := c.compileArrayQuestionQuantifierValue(values[0])

		return c.emitComparison(bytecode.OpEq, count, value)
	}

	return c.emitComparison(bytecode.OpGt, count, zero)
}

func (c *ExprCompiler) compileArrayQuestionQuantifierValue(ctx fql.IArrayQuestionQuantifierValueContext) bytecode.Operand {
	if ctx == nil {
		return bytecode.NoopOperand
	}

	if il := ctx.IntegerLiteral(); il != nil {
		return c.ctx.LiteralCompiler.CompileIntegerLiteral(il)
	}

	if pm := ctx.Param(); pm != nil {
		return c.CompileParam(pm)
	}

	return bytecode.NoopOperand
}

func (c *ExprCompiler) compileArrayApply(src bytecode.Operand, apply fql.IArrayApplyContext, tail []fql.IMemberExpressionPathContext) bytecode.Operand {
	if apply == nil {
		return src
	}

	query := c.compileQueryLiteral(apply.QueryLiteral())
	if query == bytecode.NoopOperand {
		return bytecode.NoopOperand
	}

	dst := c.ctx.Registers.Allocate()
	span := diagnostics.SpanFromRuleContext(apply)

	c.ctx.Emitter.WithSpan(span, func() {
		c.ctx.Emitter.EmitABC(bytecode.OpQuery, dst, src, query)
	})

	if len(tail) > 0 {
		return c.compileMemberExpressionSegments(dst, tail)
	}

	if dst.IsRegister() {
		c.ctx.Types.Set(dst, core.TypeList)
	}

	return dst
}

//lint:ignore U1000 Ignore unused method
func (c *ExprCompiler) compileArrayExpansion(src bytecode.Operand, expansion fql.IArrayExpansionContext, tail []fql.IMemberExpressionPathContext) bytecode.Operand {
	span := diagnostics.SpanFromRuleContext(expansion)
	inline := expansion.InlineExpression()

	return c.compileArrayIteration(src, span, tail, inline, nil)
}

func (c *ExprCompiler) compileArrayContraction(src bytecode.Operand, contraction fql.IArrayContractionContext, tail []fql.IMemberExpressionPathContext) bytecode.Operand {
	depth := arrayContractionDepth(contraction)

	if depth < 1 {
		depth = 1
	}

	span := diagnostics.SpanFromRuleContext(contraction)
	dst := c.ctx.Registers.Allocate()

	c.ctx.Emitter.WithSpan(span, func() {
		c.ctx.Emitter.EmitABx(bytecode.OpFlatten, dst, src, depth)
	})

	if dst.IsRegister() {
		c.ctx.Types.Set(dst, core.TypeList)
	}

	inline := contraction.InlineExpression()

	if len(tail) == 0 && inline == nil {
		return dst
	}

	return c.compileArrayIteration(dst, span, tail, inline, nil)
}

func arrayContractionDepth(ctx fql.IArrayContractionContext) int {
	if ctx == nil {
		return 1
	}

	count := len(ctx.GetStars())

	if count > 1 {
		return count - 1
	}

	return 1
}

func (c *ExprCompiler) compileArrayIteration(src bytecode.Operand, span source.Span, tail []fql.IMemberExpressionPathContext, inline fql.IInlineExpressionContext, extraFilters []fql.IExpressionContext) bytecode.Operand {
	tail, postLoopContraction := splitTerminalArrayContractionTail(tail)

	loop := &core.Loop{
		Kind:     core.ForInLoop,
		Type:     core.NormalLoop,
		Distinct: false,
		Allocate: true,
		Dst:      c.ctx.Registers.Allocate(),
		Src:      src,
	}

	c.ctx.Loops.Push(loop)
	c.ctx.Symbols.EnterScope()

	loop.DeclareValueVar(core.PseudoVariable, c.ctx.Symbols, core.TypeAny)
	if loop.Value.IsRegister() {
		c.ctx.Types.Set(loop.Value, core.TypeAny)
	}

	c.ctx.Emitter.WithSpan(span, func() {
		loop.EmitInitialization(c.ctx.Registers, c.ctx.Emitter)
	})

	if inline != nil {
		c.compileInlineFilter(inline)
	}

	for _, expr := range extraFilters {
		c.compileInlineFilterExpr(expr)
	}

	if inline != nil {
		c.compileInlineLimit(inline)
	}

	projection := loop.Value

	if inline != nil {
		if ret := inline.InlineReturn(); ret != nil {
			projection = c.compileWithImplicitCurrent(ret.Expression())
		}
	}

	if len(tail) > 0 {
		projection = c.compileMemberExpressionSegments(projection, tail)
	}

	c.ctx.Emitter.EmitAB(bytecode.OpPush, loop.Dst, projection)
	loop.EmitFinalization(c.ctx.Emitter)

	c.ctx.Symbols.ExitScope()
	c.ctx.Loops.Pop()

	if loop.Dst.IsRegister() {
		c.ctx.Types.Set(loop.Dst, core.TypeList)
	}

	if postLoopContraction != nil {
		return c.compileArrayContraction(loop.Dst, postLoopContraction, nil)
	}

	return loop.Dst
}

func (c *ExprCompiler) compileInlineFilter(inline fql.IInlineExpressionContext) {
	if inline == nil {
		return
	}

	filter := inline.InlineFilter()

	if filter == nil {
		return
	}

	src := c.compileWithImplicitCurrent(filter.Expression())
	label := c.ctx.Loops.Current().ContinueLabel()
	c.ctx.Emitter.EmitJumpIfFalse(src, label)
}

func (c *ExprCompiler) compileInlineFilterExpr(expr fql.IExpressionContext) {
	if expr == nil {
		return
	}

	src := c.compileWithImplicitCurrent(expr)
	label := c.ctx.Loops.Current().ContinueLabel()
	c.ctx.Emitter.EmitJumpIfFalse(src, label)
}

func (c *ExprCompiler) compileInlineLimit(inline fql.IInlineExpressionContext) {
	if inline == nil {
		return
	}

	limit := inline.InlineLimit()
	if limit == nil {
		return
	}

	clauses := limit.AllLimitClauseValue()
	if len(clauses) == 0 {
		return
	}

	c.withImplicitCurrent(func() {
		if len(clauses) == 1 {
			c.ctx.LoopCompiler.compileLimit(c.ctx.LoopCompiler.compileLimitClauseValue(clauses[0]))
			return
		}

		c.ctx.LoopCompiler.compileOffset(c.ctx.LoopCompiler.compileLimitClauseValue(clauses[0]))
		c.ctx.LoopCompiler.compileLimit(c.ctx.LoopCompiler.compileLimitClauseValue(clauses[1]))
	})
}
