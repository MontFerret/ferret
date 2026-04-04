package internal

import (
	"github.com/antlr4-go/antlr/v4"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/compiler/internal/core"
	parserd "github.com/MontFerret/ferret/v2/pkg/parser/diagnostics"
	"github.com/MontFerret/ferret/v2/pkg/parser/fql"
	"github.com/MontFerret/ferret/v2/pkg/source"
)

type waitEventCompileState struct {
	span       source.Span
	srcReg     bytecode.Operand
	eventReg   bytecode.Operand
	optsReg    bytecode.Operand
	timeoutReg bytecode.Operand
}

func (c *WaitCompiler) compileEvent(ctx fql.IWaitForEventExpressionContext) bytecode.Operand {
	state, ok := c.buildWaitEventState(ctx)
	if !ok {
		return bytecode.NoopOperand
	}

	streamReg := c.ctx.Function.Registers.Allocate()
	resultReg := c.ctx.Function.Registers.Allocate()

	c.ctx.Program.Emitter.EmitLoadNone(resultReg)
	c.emitWaitEventStreamSetup(state, streamReg)

	start := c.ctx.Program.Emitter.NewLabel()
	end := c.ctx.Program.Emitter.NewLabel()

	c.ctx.Program.Emitter.MarkLabel(start)
	c.emitWaitEventIteration(ctx, state, streamReg, bytecode.NoopOperand, start, end)

	c.ctx.Program.Emitter.MarkLabel(end)
	c.emitWaitEventCleanup(state, streamReg)

	return resultReg
}

func (c *WaitCompiler) compileEventWithTimeoutRecovery(
	ctx fql.IWaitForEventExpressionContext,
	timeoutLabel, endLabel core.Label,
) bytecode.Operand {
	streamReg := c.ctx.Function.Registers.Allocate()
	resultReg := c.ctx.Function.Registers.Allocate()
	timeoutStateReg := c.ctx.Function.Registers.Allocate()

	c.ctx.Program.Emitter.EmitLoadNone(resultReg)
	c.ctx.Program.Emitter.EmitBoolean(timeoutStateReg, false)

	state, ok := c.buildWaitEventState(ctx)
	if !ok {
		return bytecode.NoopOperand
	}

	c.emitWaitEventStreamSetup(state, streamReg)

	start := c.ctx.Program.Emitter.NewLabel()
	iterationDone := c.ctx.Program.Emitter.NewLabel()
	cleanup := c.ctx.Program.Emitter.NewLabel()

	c.ctx.Program.Emitter.MarkLabel(start)
	c.emitWaitEventIteration(ctx, state, streamReg, timeoutStateReg, start, iterationDone)

	c.ctx.Program.Emitter.EmitJump(cleanup)
	c.ctx.Program.Emitter.MarkLabel(iterationDone)
	c.ctx.Program.Emitter.EmitJump(cleanup)

	c.ctx.Program.Emitter.MarkLabel(cleanup)
	c.emitWaitEventCleanup(state, streamReg)

	c.ctx.Program.Emitter.EmitJumpIfTrue(timeoutStateReg, timeoutLabel)
	c.ctx.Program.Emitter.EmitJump(endLabel)

	return resultReg
}

func (c *WaitCompiler) buildWaitEventState(ctx fql.IWaitForEventExpressionContext) (waitEventCompileState, bool) {
	state := waitEventCompileState{
		span:     waitForSpan(ctx.WaitForEventSource(), ctx),
		srcReg:   c.CompileWaitForEventSource(ctx.WaitForEventSource()),
		eventReg: c.CompileWaitForEventName(ctx.WaitForEventName()),
	}

	if opts := ctx.OptionsClause(); opts != nil {
		state.optsReg = c.CompileOptionsClause(opts)
	}

	if timeout := ctx.TimeoutClause(); timeout != nil {
		state.timeoutReg = c.recovery.CompileDurationOperand(timeout)
		if state.timeoutReg == bytecode.NoopOperand {
			return waitEventCompileState{}, false
		}
	}

	return state, true
}

func (c *WaitCompiler) emitWaitEventStreamSetup(state waitEventCompileState, streamReg bytecode.Operand) {
	c.ctx.Program.Emitter.WithSpan(state.span, func() {
		c.ctx.Program.Emitter.EmitMove(streamReg, state.srcReg)
		c.ctx.Program.Emitter.EmitABC(bytecode.OpStream, streamReg, state.eventReg, state.optsReg)
		c.ctx.Program.Emitter.EmitABC(bytecode.OpStreamIter, streamReg, streamReg, state.timeoutReg)
	})
}

func (c *WaitCompiler) emitWaitEventIteration(
	ctx fql.IWaitForEventExpressionContext,
	state waitEventCompileState,
	streamReg, timeoutStateReg bytecode.Operand,
	restartLabel, doneLabel core.Label,
) {
	c.ctx.Program.Emitter.WithSpan(state.span, func() {
		if timeoutStateReg != bytecode.NoopOperand {
			c.ctx.Program.Emitter.EmitIterNextTimeout(streamReg, timeoutStateReg, doneLabel)
			return
		}

		c.ctx.Program.Emitter.EmitIterNext(streamReg, doneLabel)
	})

	if filter := ctx.EventFilterClause(); filter != nil {
		eventValReg, _ := c.ctx.Function.Symbols.DeclareLocal(core.PseudoVariable, core.TypeUnknown)

		c.ctx.Program.Emitter.WithSpan(state.span, func() {
			c.ctx.Program.Emitter.EmitAB(bytecode.OpIterValue, eventValReg, streamReg)
		})

		cond := c.exprs.CompileWithImplicitCurrent(filter.Expression())
		c.ctx.Program.Emitter.EmitJumpIfFalse(cond, restartLabel)
	}
}

func (c *WaitCompiler) emitWaitEventCleanup(state waitEventCompileState, streamReg bytecode.Operand) {
	c.ctx.Program.Emitter.WithSpan(state.span, func() {
		c.ctx.Program.Emitter.EmitA(bytecode.OpClose, streamReg)
	})
}

func waitForSpan(src antlr.RuleContext, fallback antlr.RuleContext) source.Span {
	span := source.Span{Start: -1, End: -1}

	if src != nil {
		if prc, ok := src.(antlr.ParserRuleContext); ok {
			span = parserd.SpanFromRuleContext(prc)
			return span
		}
	}

	if fallback != nil {
		if prc, ok := fallback.(antlr.ParserRuleContext); ok {
			span = parserd.SpanFromRuleContext(prc)
		}
	}

	return span
}

// CompileWaitForEventName processes the event name expression in a WAITFOR statement.
func (c *WaitCompiler) CompileWaitForEventName(ctx fql.IWaitForEventNameContext) bytecode.Operand {
	sl := ctx.StringLiteral()
	v := ctx.Variable()
	p := ctx.Param()
	me := ctx.MemberExpression()
	fce := ctx.FunctionCall()

	return compileFirstOperand(
		newOperandBranch(sl != nil, func() bytecode.Operand { return c.literals.CompileStringLiteral(sl) }),
		newOperandBranch(v != nil, func() bytecode.Operand { return c.exprs.CompileVariable(v) }),
		newOperandBranch(p != nil, func() bytecode.Operand { return c.exprs.CompileParam(p) }),
		newOperandBranch(me != nil, func() bytecode.Operand { return c.exprs.CompileMemberExpression(me) }),
		newOperandBranch(fce != nil, func() bytecode.Operand { return c.exprs.CompileFunctionCall(fce, false) }),
	)
}

// CompileWaitForEventSource processes the event source expression in a WAITFOR statement.
func (c *WaitCompiler) CompileWaitForEventSource(ctx fql.IWaitForEventSourceContext) bytecode.Operand {
	v := ctx.Variable()
	me := ctx.MemberExpression()
	fce := ctx.FunctionCallExpression()

	return compileFirstOperand(
		newOperandBranch(v != nil, func() bytecode.Operand { return c.exprs.CompileVariable(v) }),
		newOperandBranch(me != nil, func() bytecode.Operand { return c.exprs.CompileMemberExpression(me) }),
		newOperandBranch(fce != nil, func() bytecode.Operand { return c.exprs.CompileFunctionCallExpression(fce) }),
	)
}

// CompileOptionsClause processes the options clause in a WAITFOR statement.
func (c *WaitCompiler) CompileOptionsClause(ctx fql.IOptionsClauseContext) bytecode.Operand {
	if ol := ctx.ObjectLiteral(); ol != nil {
		return c.literals.CompileObjectLiteral(ol)
	}

	return bytecode.NoopOperand
}
