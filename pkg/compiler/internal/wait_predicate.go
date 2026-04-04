package internal

import (
	"strconv"
	"strings"

	"github.com/antlr4-go/antlr/v4"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/compiler/internal/core"
	parserd "github.com/MontFerret/ferret/v2/pkg/parser/diagnostics"
	"github.com/MontFerret/ferret/v2/pkg/parser/fql"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

type (
	waitForPredicateMode int

	waitPredicateCompileConfig struct {
		predExpr      fql.IExpressionContext
		jitterLiteral *float64
		mode          waitForPredicateMode
		timeoutReg    bytecode.Operand
		everyReg      bytecode.Operand
		capEveryReg   bytecode.Operand
		backoff       core.RetryBackoff
		jitterReg     bytecode.Operand
		hasJitter     bool
	}
)

const (
	waitForPredicateModeBool waitForPredicateMode = iota
	waitForPredicateModeExists
	waitForPredicateModeNotExists
	waitForPredicateModeValue
)

func (c *WaitCompiler) compilePredicate(ctx fql.IWaitForPredicateExpressionContext) bytecode.Operand {
	config, ok := c.prepareWaitPredicateConfig(ctx)
	if !ok {
		return bytecode.NoopOperand
	}

	if result, fastPath := c.tryCompileWaitPredicateFastPath(config); fastPath {
		return result
	}

	state := c.initWaitPredicatePollState(config)
	c.emitWaitPredicatePollLoop(config, state)

	return state.resultReg
}

func (c *WaitCompiler) compilePredicateWithTimeoutRecovery(
	ctx fql.IWaitForPredicateExpressionContext,
	timeoutLabel, endLabel core.Label,
) bytecode.Operand {
	config, ok := c.prepareWaitPredicateConfig(ctx)
	if !ok {
		return bytecode.NoopOperand
	}

	state := c.initWaitPredicatePollState(config)

	return c.emitWaitPredicatePollLoopWithRecovery(config, state, timeoutLabel, endLabel)
}

func (c *WaitCompiler) prepareWaitPredicateConfig(ctx fql.IWaitForPredicateExpressionContext) (waitPredicateCompileConfig, bool) {
	predicate := ctx.WaitForPredicate()
	if predicate == nil {
		return waitPredicateCompileConfig{}, false
	}

	predExpr := predicate.Expression()
	if predExpr == nil {
		return waitPredicateCompileConfig{}, false
	}

	if legacy := legacyWaitForOrThrowNode(predExpr); legacy != nil {
		c.ctx.Program.Errors.Add(c.ctx.Program.Errors.Create(parserd.SyntaxError, legacy, "Unexpected THROW after OR in WAITFOR predicate"))
		return waitPredicateCompileConfig{}, false
	}

	config, ok := c.buildWaitPredicateConfig(ctx, predicate, predExpr)
	if !ok {
		return waitPredicateCompileConfig{}, false
	}

	c.normalizeWaitPredicateConfig(&config)

	return config, true
}

func (c *WaitCompiler) buildWaitPredicateConfig(
	ctx fql.IWaitForPredicateExpressionContext,
	predicate fql.IWaitForPredicateContext,
	predExpr fql.IExpressionContext,
) (waitPredicateCompileConfig, bool) {
	everyReg, capEveryReg, ok := c.compileEveryClause(ctx.EveryClause())
	if !ok {
		return waitPredicateCompileConfig{}, false
	}

	jitterReg, jitterLiteral, hasJitter := c.compileJitterClause(ctx.JitterClause())
	timeoutReg := bytecode.NoopOperand

	if timeout := ctx.TimeoutClause(); timeout != nil {
		timeoutReg = c.recovery.CompileDurationOperand(timeout)
		if timeoutReg == bytecode.NoopOperand {
			return waitPredicateCompileConfig{}, false
		}
	}

	return waitPredicateCompileConfig{
		mode:          resolveWaitPredicateMode(predicate.Value() != nil, predicate.Exists() != nil, predicate.Not() != nil),
		predExpr:      predExpr,
		timeoutReg:    timeoutReg,
		everyReg:      everyReg,
		capEveryReg:   capEveryReg,
		backoff:       c.compileBackoffClause(ctx.BackoffClause()),
		jitterReg:     jitterReg,
		jitterLiteral: jitterLiteral,
		hasJitter:     hasJitter,
	}, true
}

func (c *WaitCompiler) normalizeWaitPredicateConfig(config *waitPredicateCompileConfig) {
	if !config.hasJitter {
		return
	}

	if config.jitterLiteral != nil && *config.jitterLiteral == 0 {
		config.hasJitter = false
		return
	}

	if config.jitterLiteral == nil {
		c.emitClampRange(config.jitterReg, c.facts.LoadConstant(runtime.NewFloat(0)), c.facts.LoadConstant(runtime.NewFloat(1)))
	}
}

func (c *WaitCompiler) compileEveryClause(ctx fql.IEveryClauseContext) (bytecode.Operand, bytecode.Operand, bool) {
	if ctx == nil {
		return bytecode.NoopOperand, bytecode.NoopOperand, true
	}

	values := ctx.AllEveryClauseValue()
	if len(values) == 0 {
		return bytecode.NoopOperand, bytecode.NoopOperand, true
	}

	base := c.recovery.CompileDurationOperand(values[0])
	if base == bytecode.NoopOperand {
		return bytecode.NoopOperand, bytecode.NoopOperand, false
	}

	if len(values) > 1 {
		cap := c.recovery.CompileDurationOperand(values[1])
		if cap == bytecode.NoopOperand {
			return bytecode.NoopOperand, bytecode.NoopOperand, false
		}

		return base, cap, true
	}

	return base, bytecode.NoopOperand, true
}

func (c *WaitCompiler) compileJitterClause(ctx fql.IJitterClauseContext) (bytecode.Operand, *float64, bool) {
	if ctx == nil {
		return bytecode.NoopOperand, nil, false
	}

	valueCtx := ctx.JitterClauseValue()
	if valueCtx == nil {
		return bytecode.NoopOperand, nil, false
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
			err := c.ctx.Program.Errors.Create(parserd.SemanticError, prc, "JITTER must be between 0 and 1")
			err.Hint = "Use a value between 0 and 1, e.g. JITTER 0.2."
			c.ctx.Program.Errors.Add(err)
		}
	}

	return c.compileJitterClauseValue(valueCtx), literal, true
}

func (c *WaitCompiler) compileJitterClauseValue(ctx fql.IJitterClauseValueContext) bytecode.Operand {
	if ctx == nil {
		return bytecode.NoopOperand
	}

	fl := ctx.FloatLiteral()
	il := ctx.IntegerLiteral()
	v := ctx.Variable()
	p := ctx.Param()
	me := ctx.MemberExpression()
	fc := ctx.FunctionCall()

	return compileFirstOperand(
		newOperandBranch(fl != nil, func() bytecode.Operand { return c.literals.CompileFloatLiteral(fl) }),
		newOperandBranch(il != nil, func() bytecode.Operand { return c.literals.CompileIntegerLiteral(il) }),
		newOperandBranch(v != nil, func() bytecode.Operand { return c.exprs.CompileVariable(v) }),
		newOperandBranch(p != nil, func() bytecode.Operand { return c.exprs.CompileParam(p) }),
		newOperandBranch(me != nil, func() bytecode.Operand { return c.exprs.CompileMemberExpression(me) }),
		newOperandBranch(fc != nil, func() bytecode.Operand { return c.exprs.CompileFunctionCall(fc, false) }),
	)
}

func (c *WaitCompiler) compileBackoffClause(ctx fql.IBackoffClauseContext) core.RetryBackoff {
	if ctx == nil {
		return core.RetryBackoffNone
	}

	strategyCtx := ctx.BackoffStrategy()
	if strategyCtx == nil {
		return core.RetryBackoffNone
	}

	var strategy string
	switch {
	case strategyCtx.None() != nil:
		strategy = "NONE"
	case strategyCtx.Identifier() != nil:
		strategy = strategyCtx.Identifier().GetText()
	case strategyCtx.StringLiteral() != nil:
		if val, ok := parseStringLiteralConst(strategyCtx.StringLiteral()); ok {
			strategy = val.String()
		} else {
			if prc, ok := ctx.(antlr.ParserRuleContext); ok {
				err := c.ctx.Program.Errors.Create(parserd.SemanticError, prc, "BACKOFF strategy must be a constant string")
				err.Hint = "Use one of: NONE, LINEAR, EXPONENTIAL."
				c.ctx.Program.Errors.Add(err)
			}
			return core.RetryBackoffNone
		}
	default:
		return core.RetryBackoffNone
	}

	strategy = strings.ToUpper(strings.TrimSpace(strategy))

	switch strategy {
	case "", "NONE":
		return core.RetryBackoffNone
	case "LINEAR":
		return core.RetryBackoffLinear
	case "EXPONENTIAL":
		return core.RetryBackoffExponential
	default:
		if prc, ok := ctx.(antlr.ParserRuleContext); ok {
			err := c.ctx.Program.Errors.Create(parserd.SemanticError, prc, "Unknown BACKOFF strategy")
			err.Hint = "Use one of: NONE, LINEAR, EXPONENTIAL."
			c.ctx.Program.Errors.Add(err)
		}

		return core.RetryBackoffNone
	}
}
