package internal

import (
	"math"
	"strconv"
	"strings"

	"github.com/antlr4-go/antlr/v4"

	"github.com/MontFerret/ferret/v2/pkg/compiler/internal/core"
	"github.com/MontFerret/ferret/v2/pkg/compiler/internal/diagnostics"
	"github.com/MontFerret/ferret/v2/pkg/file"
	"github.com/MontFerret/ferret/v2/pkg/parser/fql"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/vm"
)

// WaitCompiler handles the compilation of WAITFOR expressions in FQL queries.
// It transforms wait operations into VM instructions for event streaming and polling.
type WaitCompiler struct {
	ctx *CompilerContext
}

// NewWaitCompiler creates a new instance of WaitCompiler with the given compiler context.
func NewWaitCompiler(ctx *CompilerContext) *WaitCompiler {
	return &WaitCompiler{
		ctx: ctx,
	}
}

type waitForPredicateMode int

const (
	waitForPredicateModeBool waitForPredicateMode = iota
	waitForPredicateModeExists
	waitForPredicateModeNotExists
	waitForPredicateModeValue
)

type waitForBackoff int

const (
	waitForBackoffNone waitForBackoff = iota
	waitForBackoffLinear
	waitForBackoffExponential
)

const waitForDefaultEveryMs = 100

type durationClause interface {
	DurationLiteral() fql.IDurationLiteralContext
	IntegerLiteral() fql.IIntegerLiteralContext
	FloatLiteral() fql.IFloatLiteralContext
	Variable() fql.IVariableContext
	Param() fql.IParamContext
	MemberExpression() fql.IMemberExpressionContext
	FunctionCall() fql.IFunctionCallContext
}

// Compile processes a WAITFOR expression from the FQL AST and generates the appropriate VM instructions.
func (c *WaitCompiler) Compile(ctx fql.IWaitForExpressionContext) vm.Operand {
	if ctx == nil {
		return vm.NoopOperand
	}

	c.ctx.Symbols.EnterScope()
	defer c.ctx.Symbols.ExitScope()

	if orThrow := ctx.WaitForOrThrowClause(); orThrow != nil {
		if prc, ok := orThrow.(antlr.ParserRuleContext); ok {
			err := c.ctx.Errors.Create(diagnostics.SemanticError, prc, "OR THROW is not supported")
			err.Hint = "Remove OR THROW and handle timeouts explicitly."
			c.ctx.Errors.Add(err)
		}
	}

	if ev := ctx.WaitForEventExpression(); ev != nil {
		return c.compileEvent(ev)
	}

	if pred := ctx.WaitForPredicateExpression(); pred != nil {
		return c.compilePredicate(pred)
	}

	return vm.NoopOperand
}

func (c *WaitCompiler) compileEvent(ctx fql.IWaitForEventExpressionContext) vm.Operand {
	srcReg := c.CompileWaitForEventSource(ctx.WaitForEventSource())
	eventReg := c.CompileWaitForEventName(ctx.WaitForEventName())

	var optsReg vm.Operand
	if opts := ctx.OptionsClause(); opts != nil {
		optsReg = c.CompileOptionsClause(opts)
	}

	var timeoutReg vm.Operand
	if timeout := ctx.TimeoutClause(); timeout != nil {
		timeoutReg = c.CompileTimeoutClauseContext(timeout)
	}

	streamReg := c.ctx.Registers.Allocate()
	resultReg := c.ctx.Registers.Allocate()

	c.ctx.Emitter.EmitLoadNone(resultReg)

	span := waitForSpan(ctx.WaitForEventSource(), ctx)

	c.ctx.Emitter.WithSpan(span, func() {
		c.ctx.Emitter.EmitMove(streamReg, srcReg)
		c.ctx.Emitter.EmitABC(vm.OpStream, streamReg, eventReg, optsReg)
		c.ctx.Emitter.EmitAB(vm.OpStreamIter, streamReg, timeoutReg)
	})

	start := c.ctx.Emitter.NewLabel()
	end := c.ctx.Emitter.NewLabel()

	c.ctx.Emitter.MarkLabel(start)
	c.ctx.Emitter.WithSpan(span, func() {
		c.ctx.Emitter.EmitIterNext(streamReg, end)
	})

	if filter := ctx.FilterClause(); filter != nil {
		eventValReg, _ := c.ctx.Symbols.DeclareLocal(core.PseudoVariable, core.TypeUnknown)
		c.ctx.Emitter.WithSpan(span, func() {
			c.ctx.Emitter.EmitAB(vm.OpIterValue, eventValReg, streamReg)
		})

		cond := c.ctx.ExprCompiler.Compile(filter.Expression())
		c.ctx.Emitter.EmitJumpIfFalse(cond, start)
	}

	c.ctx.Emitter.MarkLabel(end)
	c.ctx.Emitter.WithSpan(span, func() {
		c.ctx.Emitter.EmitA(vm.OpClose, streamReg)
	})

	return resultReg
}

func (c *WaitCompiler) compilePredicate(ctx fql.IWaitForPredicateExpressionContext) vm.Operand {
	predicate := ctx.WaitForPredicate()
	if predicate == nil {
		return vm.NoopOperand
	}

	mode := waitForPredicateModeBool
	if predicate.Value() != nil {
		mode = waitForPredicateModeValue
	} else if predicate.Exists() != nil {
		mode = waitForPredicateModeExists
		if predicate.Not() != nil {
			mode = waitForPredicateModeNotExists
		}
	}

	predExpr := predicate.Expression()
	if predExpr == nil {
		return vm.NoopOperand
	}

	timeoutReg := c.compileDurationClause(ctx.TimeoutClause())
	everyReg, capEveryReg := c.compileEveryClause(ctx.EveryClause())
	backoff := c.compileBackoffClause(ctx.BackoffClause())
	jitterReg, jitterLiteral, hasJitter := c.compileJitterClause(ctx.JitterClause())
	if hasJitter {
		if jitterLiteral != nil && *jitterLiteral == 0 {
			hasJitter = false
		} else if jitterLiteral == nil {
			c.emitClampRange(jitterReg, loadConstant(c.ctx, runtime.NewFloat(0)), loadConstant(c.ctx, runtime.NewFloat(1)))
		}
	}

	switch mode {
	case waitForPredicateModeBool:
		if truth, ok := literalTruthinessFromExpression(predExpr); ok {
			if truth {
				resultReg := c.ctx.Registers.Allocate()
				c.ctx.Emitter.EmitBoolean(resultReg, true)
				return resultReg
			}
			if timeoutReg != vm.NoopOperand {
				c.ctx.Emitter.EmitA(vm.OpSleep, timeoutReg)
				resultReg := c.ctx.Registers.Allocate()
				c.ctx.Emitter.EmitBoolean(resultReg, false)
				return resultReg
			}
		}
	default:
		if exists, ok := literalExistsFromExpression(predExpr); ok {
			cond := exists
			if mode == waitForPredicateModeNotExists {
				cond = !exists
			}
			if cond {
				if mode == waitForPredicateModeValue {
					return c.ctx.ExprCompiler.Compile(predExpr)
				}
				resultReg := c.ctx.Registers.Allocate()
				c.ctx.Emitter.EmitBoolean(resultReg, true)
				return resultReg
			}
			if timeoutReg != vm.NoopOperand {
				c.ctx.Emitter.EmitA(vm.OpSleep, timeoutReg)
				resultReg := c.ctx.Registers.Allocate()
				if mode == waitForPredicateModeValue {
					c.ctx.Emitter.EmitLoadNone(resultReg)
				} else {
					c.ctx.Emitter.EmitBoolean(resultReg, false)
				}
				return resultReg
			}
		}
	}

	baseEveryReg := c.ctx.Registers.Allocate()
	if everyReg != vm.NoopOperand {
		c.ctx.Emitter.EmitMove(baseEveryReg, everyReg)
	} else {
		c.ctx.Emitter.EmitLoadConst(baseEveryReg, c.ctx.Symbols.AddConstant(runtime.NewInt(waitForDefaultEveryMs)))
	}

	pollReg := baseEveryReg
	var intervalReg vm.Operand
	if backoff != waitForBackoffNone {
		intervalReg = c.ctx.Registers.Allocate()
		c.ctx.Emitter.EmitMove(intervalReg, baseEveryReg)
		pollReg = intervalReg
	}

	resultReg := c.ctx.Registers.Allocate()
	if mode == waitForPredicateModeValue {
		c.ctx.Emitter.EmitLoadNone(resultReg)
	} else {
		c.ctx.Emitter.EmitBoolean(resultReg, false)
	}

	var startReg vm.Operand
	var unitReg vm.Operand
	if timeoutReg != vm.NoopOperand {
		startReg = c.emitNow()
		unitReg = loadConstant(c.ctx, runtime.NewString("f"))
	}

	start := c.ctx.Emitter.NewLabel()
	success := c.ctx.Emitter.NewLabel()
	timeoutLabel := c.ctx.Emitter.NewLabel()
	end := c.ctx.Emitter.NewLabel()

	c.ctx.Emitter.MarkLabel(start)

	valueReg := c.ctx.ExprCompiler.Compile(predExpr)

	var condReg vm.Operand
	switch mode {
	case waitForPredicateModeValue:
		condReg = c.emitExistsCheck(valueReg)
	case waitForPredicateModeExists:
		condReg = c.emitExistsCheck(valueReg)
	case waitForPredicateModeNotExists:
		existsReg := c.emitExistsCheck(valueReg)
		condReg = c.ctx.Registers.Allocate()
		c.ctx.Emitter.EmitAB(vm.OpNot, condReg, existsReg)
	default:
		condReg = c.ctx.Registers.Allocate()
		c.ctx.Emitter.EmitAB(vm.OpCastBool, condReg, valueReg)
	}

	c.ctx.Emitter.EmitJumpIfTrue(condReg, success)

	var elapsedReg vm.Operand
	if timeoutReg != vm.NoopOperand {
		nowReg := c.emitNow()
		elapsedReg = c.emitDateDiff(startReg, nowReg, unitReg)
		reachedReg := c.ctx.Registers.Allocate()
		c.ctx.Emitter.EmitGte(reachedReg, elapsedReg, timeoutReg)
		c.ctx.Emitter.EmitJumpIfTrue(reachedReg, timeoutLabel)
	}

	sleepIntervalReg := pollReg
	if hasJitter || capEveryReg != vm.NoopOperand {
		sleepIntervalReg = c.ctx.Registers.Allocate()
		c.ctx.Emitter.EmitMove(sleepIntervalReg, pollReg)
		if hasJitter {
			c.emitApplyJitter(sleepIntervalReg, jitterReg)
		}
		if capEveryReg != vm.NoopOperand {
			c.emitClampMax(sleepIntervalReg, capEveryReg)
		}
	}
	c.emitWaitSleep(sleepIntervalReg, timeoutReg, elapsedReg)
	if backoff != waitForBackoffNone {
		c.emitBackoffUpdate(backoff, intervalReg, baseEveryReg)
		if capEveryReg != vm.NoopOperand {
			c.emitClampMax(intervalReg, capEveryReg)
		}
	}
	c.ctx.Emitter.EmitJump(start)

	c.ctx.Emitter.MarkLabel(success)
	if mode == waitForPredicateModeValue {
		c.ctx.Emitter.EmitMove(resultReg, valueReg)
	} else {
		c.ctx.Emitter.EmitBoolean(resultReg, true)
	}
	c.ctx.Emitter.EmitJump(end)

	c.ctx.Emitter.MarkLabel(timeoutLabel)
	if mode == waitForPredicateModeValue {
		c.ctx.Emitter.EmitLoadNone(resultReg)
	} else {
		c.ctx.Emitter.EmitBoolean(resultReg, false)
	}
	c.ctx.Emitter.MarkLabel(end)

	return resultReg
}

func (c *WaitCompiler) emitExistsCheck(val vm.Operand) vm.Operand {
	dst := c.ctx.Registers.Allocate()
	c.ctx.Emitter.EmitAB(vm.OpExists, dst, val)
	c.ctx.Types.Set(dst, core.TypeBool)
	return dst
}

func (c *WaitCompiler) emitNow() vm.Operand {
	return c.ctx.ExprCompiler.CompileFunctionCallByNameWith(runtime.NewString("NOW"), false, nil)
}

func (c *WaitCompiler) emitDateDiff(start, end, unit vm.Operand) vm.Operand {
	return c.emitFunctionCall(runtime.NewString("DATE_DIFF"), start, end, unit)
}

func (c *WaitCompiler) emitFunctionCall(name runtime.String, args ...vm.Operand) vm.Operand {
	if len(args) == 0 {
		return c.ctx.ExprCompiler.CompileFunctionCallByNameWith(name, false, nil)
	}

	seq := c.ctx.Registers.AllocateSequence(len(args))
	for i, arg := range args {
		c.ctx.Emitter.EmitMove(seq[i], arg)
		c.ctx.Types.Set(seq[i], operandType(c.ctx, arg))
	}

	return c.ctx.ExprCompiler.CompileFunctionCallByNameWith(name, false, seq)
}

func (c *WaitCompiler) emitWaitSleep(intervalReg, timeoutReg, elapsedReg vm.Operand) {
	if timeoutReg == vm.NoopOperand {
		c.ctx.Emitter.EmitA(vm.OpSleep, intervalReg)
		return
	}

	sleepReg := c.ctx.Registers.Allocate()
	c.ctx.Emitter.EmitMove(sleepReg, intervalReg)

	remainingReg := c.ctx.Registers.Allocate()
	c.ctx.Emitter.EmitABC(vm.OpSub, remainingReg, timeoutReg, elapsedReg)

	shouldTrim := c.ctx.Registers.Allocate()
	c.ctx.Emitter.EmitLt(shouldTrim, remainingReg, sleepReg)

	useRemaining := c.ctx.Emitter.NewLabel()
	continueSleep := c.ctx.Emitter.NewLabel()

	c.ctx.Emitter.EmitJumpIfTrue(shouldTrim, useRemaining)
	c.ctx.Emitter.EmitJump(continueSleep)

	c.ctx.Emitter.MarkLabel(useRemaining)
	c.ctx.Emitter.EmitMove(sleepReg, remainingReg)
	c.ctx.Emitter.MarkLabel(continueSleep)

	c.ctx.Emitter.EmitA(vm.OpSleep, sleepReg)
}

func (c *WaitCompiler) emitBackoffUpdate(strategy waitForBackoff, intervalReg, baseEveryReg vm.Operand) {
	switch strategy {
	case waitForBackoffLinear:
		c.ctx.Emitter.EmitABC(vm.OpAdd, intervalReg, intervalReg, baseEveryReg)
	case waitForBackoffExponential:
		twoReg := loadConstant(c.ctx, runtime.NewInt(2))
		c.ctx.Emitter.EmitABC(vm.OpMulti, intervalReg, intervalReg, twoReg)
	default:
		return
	}
}

func (c *WaitCompiler) emitClampMin(target, min vm.Operand) {
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

func (c *WaitCompiler) emitClampMax(target, max vm.Operand) {
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

func (c *WaitCompiler) emitClampRange(target, min, max vm.Operand) {
	c.emitClampMin(target, min)
	c.emitClampMax(target, max)
}

func (c *WaitCompiler) emitApplyJitter(intervalReg, jitterReg vm.Operand) {
	if intervalReg == vm.NoopOperand || jitterReg == vm.NoopOperand {
		return
	}

	randReg := c.ctx.Registers.Allocate()
	c.ctx.Emitter.EmitA(vm.OpRand, randReg)

	twoReg := loadConstant(c.ctx, runtime.NewFloat(2))
	oneReg := loadConstant(c.ctx, runtime.NewFloat(1))

	twoJitterReg := c.ctx.Registers.Allocate()
	c.ctx.Emitter.EmitABC(vm.OpMulti, twoJitterReg, jitterReg, twoReg)

	randScaleReg := c.ctx.Registers.Allocate()
	c.ctx.Emitter.EmitABC(vm.OpMulti, randScaleReg, randReg, twoJitterReg)

	oneMinusReg := c.ctx.Registers.Allocate()
	c.ctx.Emitter.EmitABC(vm.OpSub, oneMinusReg, oneReg, jitterReg)

	multiplierReg := c.ctx.Registers.Allocate()
	c.ctx.Emitter.EmitABC(vm.OpAdd, multiplierReg, oneMinusReg, randScaleReg)

	c.ctx.Emitter.EmitABC(vm.OpMulti, intervalReg, intervalReg, multiplierReg)
}

func waitForSpan(source antlr.RuleContext, fallback antlr.RuleContext) file.Span {
	span := file.Span{Start: -1, End: -1}

	if source != nil {
		if prc, ok := source.(antlr.ParserRuleContext); ok {
			span = diagnostics.SpanFromRuleContext(prc)
			return span
		}
	}

	if fallback != nil {
		if prc, ok := fallback.(antlr.ParserRuleContext); ok {
			span = diagnostics.SpanFromRuleContext(prc)
		}
	}

	return span
}

// CompileWaitForEventName processes the event name expression in a WAITFOR statement.
func (c *WaitCompiler) CompileWaitForEventName(ctx fql.IWaitForEventNameContext) vm.Operand {
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

	if fce := ctx.FunctionCall(); fce != nil {
		return c.ctx.ExprCompiler.CompileFunctionCall(fce, false)
	}

	return vm.NoopOperand
}

// CompileWaitForEventSource processes the event source expression in a WAITFOR statement.
func (c *WaitCompiler) CompileWaitForEventSource(ctx fql.IWaitForEventSourceContext) vm.Operand {
	if v := ctx.Variable(); v != nil {
		return c.ctx.ExprCompiler.CompileVariable(v)
	}

	if me := ctx.MemberExpression(); me != nil {
		return c.ctx.ExprCompiler.CompileMemberExpression(me)
	}

	if fce := ctx.FunctionCallExpression(); fce != nil {
		return c.ctx.ExprCompiler.CompileFunctionCallExpression(fce)
	}

	return vm.NoopOperand
}

// CompileOptionsClause processes the options clause in a WAITFOR statement.
func (c *WaitCompiler) CompileOptionsClause(ctx fql.IOptionsClauseContext) vm.Operand {
	if ol := ctx.ObjectLiteral(); ol != nil {
		return c.ctx.LiteralCompiler.CompileObjectLiteral(ol)
	}

	return vm.NoopOperand
}

// CompileTimeoutClauseContext processes the timeout clause in a WAITFOR statement.
func (c *WaitCompiler) CompileTimeoutClauseContext(ctx fql.ITimeoutClauseContext) vm.Operand {
	return c.compileDurationClause(ctx)
}

func (c *WaitCompiler) compileEveryClause(ctx fql.IEveryClauseContext) (vm.Operand, vm.Operand) {
	if ctx == nil {
		return vm.NoopOperand, vm.NoopOperand
	}

	values := ctx.AllEveryClauseValue()
	if len(values) == 0 {
		return vm.NoopOperand, vm.NoopOperand
	}

	base := c.compileDurationClause(values[0])
	if len(values) > 1 {
		return base, c.compileDurationClause(values[1])
	}

	return base, vm.NoopOperand
}

func (c *WaitCompiler) compileDurationClause(ctx durationClause) vm.Operand {
	if ctx == nil {
		return vm.NoopOperand
	}

	if dl := ctx.DurationLiteral(); dl != nil {
		val, err := parseDurationLiteral(dl.GetText())
		if err != nil {
			panic(err)
		}
		return loadConstant(c.ctx, val)
	}

	if il := ctx.IntegerLiteral(); il != nil {
		return c.ctx.LiteralCompiler.CompileIntegerLiteral(il)
	}

	if fl := ctx.FloatLiteral(); fl != nil {
		return c.ctx.LiteralCompiler.CompileFloatLiteral(fl)
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

func (c *WaitCompiler) compileJitterClause(ctx fql.IJitterClauseContext) (vm.Operand, *float64, bool) {
	if ctx == nil {
		return vm.NoopOperand, nil, false
	}

	valueCtx := ctx.JitterClauseValue()
	if valueCtx == nil {
		return vm.NoopOperand, nil, false
	}

	var literal *float64
	if fl := valueCtx.FloatLiteral(); fl != nil {
		if val, err := strconv.ParseFloat(fl.GetText(), 64); err == nil {
			literal = &val
		}
	} else if il := valueCtx.IntegerLiteral(); il != nil {
		if val, err := strconv.ParseFloat(il.GetText(), 64); err == nil {
			literal = &val
		}
	}

	if literal != nil && (*literal < 0 || *literal > 1) {
		if prc, ok := valueCtx.(antlr.ParserRuleContext); ok {
			err := c.ctx.Errors.Create(diagnostics.SemanticError, prc, "JITTER must be between 0 and 1")
			err.Hint = "Use a value between 0 and 1, e.g. JITTER 0.2."
			c.ctx.Errors.Add(err)
		}
	}

	return c.compileJitterClauseValue(valueCtx), literal, true
}

func (c *WaitCompiler) compileJitterClauseValue(ctx fql.IJitterClauseValueContext) vm.Operand {
	if ctx == nil {
		return vm.NoopOperand
	}

	if fl := ctx.FloatLiteral(); fl != nil {
		return c.ctx.LiteralCompiler.CompileFloatLiteral(fl)
	}

	if il := ctx.IntegerLiteral(); il != nil {
		return c.ctx.LiteralCompiler.CompileIntegerLiteral(il)
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

func (c *WaitCompiler) compileBackoffClause(ctx fql.IBackoffClauseContext) waitForBackoff {
	if ctx == nil {
		return waitForBackoffNone
	}

	strategyCtx := ctx.BackoffStrategy()
	if strategyCtx == nil {
		return waitForBackoffNone
	}

	var strategy string
	switch {
	case strategyCtx.None() != nil:
		strategy = "NONE"
	case strategyCtx.Identifier() != nil:
		strategy = strategyCtx.Identifier().GetText()
	case strategyCtx.StringLiteral() != nil:
		strategy = parseStringLiteral(strategyCtx.StringLiteral()).String()
	default:
		return waitForBackoffNone
	}

	strategy = strings.ToUpper(strings.TrimSpace(strategy))

	switch strategy {
	case "", "NONE":
		return waitForBackoffNone
	case "LINEAR":
		return waitForBackoffLinear
	case "EXPONENTIAL":
		return waitForBackoffExponential
	default:
		if prc, ok := ctx.(antlr.ParserRuleContext); ok {
			err := c.ctx.Errors.Create(diagnostics.SemanticError, prc, "Unknown BACKOFF strategy")
			err.Hint = "Use one of: NONE, LINEAR, EXPONENTIAL."
			c.ctx.Errors.Add(err)
		}
		return waitForBackoffNone
	}
}

func parseDurationLiteral(text string) (runtime.Value, error) {
	raw := strings.ToUpper(strings.TrimSpace(text))
	if raw == "" {
		return runtime.None, strconv.ErrSyntax
	}

	var unit string
	switch {
	case strings.HasSuffix(raw, "MS"):
		unit = "MS"
		raw = strings.TrimSuffix(raw, "MS")
	case strings.HasSuffix(raw, "S"):
		unit = "S"
		raw = strings.TrimSuffix(raw, "S")
	case strings.HasSuffix(raw, "M"):
		unit = "M"
		raw = strings.TrimSuffix(raw, "M")
	case strings.HasSuffix(raw, "H"):
		unit = "H"
		raw = strings.TrimSuffix(raw, "H")
	case strings.HasSuffix(raw, "D"):
		unit = "D"
		raw = strings.TrimSuffix(raw, "D")
	default:
		return runtime.None, strconv.ErrSyntax
	}

	if raw == "" {
		return runtime.None, strconv.ErrSyntax
	}

	value, err := strconv.ParseFloat(raw, 64)
	if err != nil {
		return runtime.None, err
	}

	multiplier := float64(1)
	switch unit {
	case "MS":
		multiplier = 1
	case "S":
		multiplier = 1000
	case "M":
		multiplier = 60000
	case "H":
		multiplier = 3600000
	case "D":
		multiplier = 86400000
	default:
		return runtime.None, strconv.ErrSyntax
	}

	ms := value * multiplier
	if frac := math.Mod(ms, 1); frac == 0 {
		return runtime.NewInt64(int64(ms)), nil
	}

	return runtime.NewFloat(ms), nil
}

func literalFromExpression(ctx fql.IExpressionContext) fql.ILiteralContext {
	if ctx == nil {
		return nil
	}

	predicate := ctx.Predicate()
	if predicate == nil {
		return nil
	}

	atom := predicate.ExpressionAtom()
	if atom == nil {
		return nil
	}

	return atom.Literal()
}

func literalExistsFromExpression(ctx fql.IExpressionContext) (bool, bool) {
	lit := literalFromExpression(ctx)
	if lit == nil {
		return false, false
	}

	switch {
	case lit.NoneLiteral() != nil:
		return false, true
	case lit.StringLiteral() != nil:
		str := parseStringLiteral(lit.StringLiteral())
		return str.String() != "", true
	case lit.ArrayLiteral() != nil:
		arr := lit.ArrayLiteral()
		return arr.ArgumentList() != nil, true
	case lit.ObjectLiteral() != nil:
		obj := lit.ObjectLiteral()
		return len(obj.AllPropertyAssignment()) > 0, true
	default:
		return true, true
	}
}

func literalTruthinessFromExpression(ctx fql.IExpressionContext) (bool, bool) {
	lit := literalFromExpression(ctx)
	if lit == nil {
		return false, false
	}

	switch {
	case lit.NoneLiteral() != nil:
		return false, true
	case lit.BooleanLiteral() != nil:
		return strings.ToLower(lit.BooleanLiteral().GetText()) == "true", true
	case lit.IntegerLiteral() != nil:
		val, err := strconv.Atoi(lit.IntegerLiteral().GetText())
		if err != nil {
			return false, false
		}
		return val != 0, true
	case lit.FloatLiteral() != nil:
		val, err := strconv.ParseFloat(lit.FloatLiteral().GetText(), 64)
		if err != nil {
			return false, false
		}
		return val != 0, true
	case lit.StringLiteral() != nil:
		str := parseStringLiteral(lit.StringLiteral())
		return str.String() != "", true
	default:
		return true, true
	}
}
