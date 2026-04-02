package internal

import (
	"math"
	"strconv"
	"strings"

	"github.com/antlr4-go/antlr/v4"

	"github.com/MontFerret/ferret/v2/pkg/source"

	parser "github.com/MontFerret/ferret/v2/pkg/parser/diagnostics"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"

	"github.com/MontFerret/ferret/v2/pkg/compiler/internal/core"
	"github.com/MontFerret/ferret/v2/pkg/parser/fql"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

type (
	// WaitCompiler handles the compilation of WAITFOR expressions in FQL queries.
	// It transforms wait operations into VM instructions for event streaming and polling.
	WaitCompiler struct {
		ctx *CompilerContext
	}

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

	waitPredicatePollState struct {
		baseEveryReg bytecode.Operand
		pollReg      bytecode.Operand
		intervalReg  bytecode.Operand
		resultReg    bytecode.Operand
		startReg     bytecode.Operand
		unitReg      bytecode.Operand
	}
)

const (
	waitForPredicateModeBool waitForPredicateMode = iota
	waitForPredicateModeExists
	waitForPredicateModeNotExists
	waitForPredicateModeValue
)

const waitForDefaultEveryMs = 100

// NewWaitCompiler creates a new instance of WaitCompiler with the given compiler context.
func NewWaitCompiler(ctx *CompilerContext) *WaitCompiler {
	return &WaitCompiler{
		ctx: ctx,
	}
}

// Compile processes a WAITFOR expression from the FQL AST and generates the appropriate VM instructions.
func (c *WaitCompiler) Compile(ctx fql.IWaitForExpressionContext) bytecode.Operand {
	return c.compileWithOuterRecovery(ctx, core.RecoveryPlan{})
}

func (c *WaitCompiler) compileWithOuterRecovery(ctx fql.IWaitForExpressionContext, outerPlan core.RecoveryPlan) bytecode.Operand {
	if ctx == nil {
		return bytecode.NoopOperand
	}

	plan := mergeRecoveryPlans(c.ctx, collectRecoveryPlan(c.ctx, ctx, core.RecoveryPlanOptions{
		AllowTimeout: true,
		HasTimeout:   waitForHasExplicitTimeoutClause(ctx),
	}), outerPlan)
	c.ctx.Symbols.EnterScope()
	defer c.ctx.Symbols.ExitScope()

	if ev := ctx.WaitForEventExpression(); ev != nil {
		return c.compileEventWithPlan(ev, plan)
	}

	if pred := ctx.WaitForPredicateExpression(); pred != nil {
		return c.compilePredicateWithPlan(pred, plan)
	}

	return bytecode.NoopOperand
}

func waitForHasExplicitTimeoutClause(ctx fql.IWaitForExpressionContext) bool {
	if ctx == nil {
		return false
	}

	if ev := ctx.WaitForEventExpression(); ev != nil && ev.TimeoutClause() != nil {
		return true
	}

	if pred := ctx.WaitForPredicateExpression(); pred != nil && pred.TimeoutClause() != nil {
		return true
	}

	return false
}

func normalizeWaitRecoveryPlan(plan core.RecoveryPlan) core.RecoveryPlan {
	normalized := plan

	if recoveryHandlerRetries(normalized.OnError) && normalized.OnError.Retry != nil &&
		normalized.OnError.Retry.FinalActionKind != core.RecoveryActionReturn &&
		normalized.OnError.Retry.Count <= 0 {
		normalized.OnError = nil
	}

	return normalized
}

func (c *WaitCompiler) compileEventWithPlan(ctx fql.IWaitForEventExpressionContext, plan core.RecoveryPlan) bytecode.Operand {
	plan = normalizeWaitRecoveryPlan(plan)

	if plan.OnTimeout == nil && (plan.OnError == nil || plan.OnError.ActionKind == core.RecoveryActionFail) {
		return c.compileEvent(ctx)
	}

	return widenRecoveryResultType(c.ctx, c.compileEventWithRecovery(ctx, plan), plan)
}

func (c *WaitCompiler) compilePredicateWithPlan(ctx fql.IWaitForPredicateExpressionContext, plan core.RecoveryPlan) bytecode.Operand {
	plan = normalizeWaitRecoveryPlan(plan)

	if plan.OnTimeout == nil {
		errorPlan := core.RecoveryPlan{OnError: plan.OnError}

		return c.ctx.PolicyCompiler.CompileWithRecoveryPlan(errorPlan, core.CatchJumpModeNone, func() bytecode.Operand {
			return c.compilePredicate(ctx)
		})
	}

	return widenRecoveryResultType(c.ctx, c.compilePredicateWithRecovery(ctx, plan), plan)
}

func (c *WaitCompiler) compileEvent(ctx fql.IWaitForEventExpressionContext) bytecode.Operand {
	srcReg := c.CompileWaitForEventSource(ctx.WaitForEventSource())
	eventReg := c.CompileWaitForEventName(ctx.WaitForEventName())

	var optsReg bytecode.Operand
	if opts := ctx.OptionsClause(); opts != nil {
		optsReg = c.CompileOptionsClause(opts)
	}

	var timeoutReg bytecode.Operand
	if timeout := ctx.TimeoutClause(); timeout != nil {
		timeoutReg = c.CompileTimeoutClauseContext(timeout)
	}

	streamReg := c.ctx.Registers.Allocate()
	resultReg := c.ctx.Registers.Allocate()

	c.ctx.Emitter.EmitLoadNone(resultReg)

	span := waitForSpan(ctx.WaitForEventSource(), ctx)

	c.ctx.Emitter.WithSpan(span, func() {
		c.ctx.Emitter.EmitMove(streamReg, srcReg)
		c.ctx.Emitter.EmitABC(bytecode.OpStream, streamReg, eventReg, optsReg)
		c.ctx.Emitter.EmitABC(bytecode.OpStreamIter, streamReg, streamReg, timeoutReg)
	})

	start := c.ctx.Emitter.NewLabel()
	end := c.ctx.Emitter.NewLabel()

	c.ctx.Emitter.MarkLabel(start)
	c.ctx.Emitter.WithSpan(span, func() {
		c.ctx.Emitter.EmitIterNext(streamReg, end)
	})

	if filter := ctx.EventFilterClause(); filter != nil {
		eventValReg, _ := c.ctx.Symbols.DeclareLocal(core.PseudoVariable, core.TypeUnknown)

		c.ctx.Emitter.WithSpan(span, func() {
			c.ctx.Emitter.EmitAB(bytecode.OpIterValue, eventValReg, streamReg)
		})

		cond := c.ctx.ExprCompiler.compileWithImplicitCurrent(filter.Expression())
		c.ctx.Emitter.EmitJumpIfFalse(cond, start)
	}

	c.ctx.Emitter.MarkLabel(end)
	c.ctx.Emitter.WithSpan(span, func() {
		c.ctx.Emitter.EmitA(bytecode.OpClose, streamReg)
	})

	return resultReg
}

func (c *WaitCompiler) compileEventWithRecovery(ctx fql.IWaitForEventExpressionContext, plan core.RecoveryPlan) bytecode.Operand {
	span := waitForSpan(ctx.WaitForEventSource(), ctx)
	streamReg := c.ctx.Registers.Allocate()
	resultReg := c.ctx.Registers.Allocate()
	timeoutStateReg := c.ctx.Registers.Allocate()
	errorStateReg := c.ctx.Registers.Allocate()

	c.ctx.Emitter.EmitLoadNone(resultReg)
	c.ctx.Emitter.EmitBoolean(timeoutStateReg, false)
	c.ctx.Emitter.EmitBoolean(errorStateReg, false)

	var (
		zeroReg             bytecode.Operand
		retriesRemainingReg bytecode.Operand
		retryDelayState     core.RetryDelayState
		retryStartLabel     core.Label
		finalAttemptLabel   core.Label
	)

	if recoveryHandlerRetries(plan.OnError) {
		zeroReg = loadConstant(c.ctx, runtime.ZeroInt)
		retriesRemainingReg = loadConstant(c.ctx, runtime.NewInt(plan.OnError.Retry.Count))
		retryDelayState = initRetryDelayState(c.ctx, plan.OnError.Retry)
		retryStartLabel = c.ctx.Emitter.NewLabel("waitfor", "event", "retry")
		if plan.OnError.Retry.FinalActionKind != core.RecoveryActionReturn {
			finalAttemptLabel = c.ctx.Emitter.NewLabel("waitfor", "event", "final")
		}
		c.ctx.Emitter.MarkLabel(retryStartLabel)
	}

	startCatch := c.ctx.Emitter.Size()

	srcReg := c.CompileWaitForEventSource(ctx.WaitForEventSource())
	eventReg := c.CompileWaitForEventName(ctx.WaitForEventName())

	var optsReg bytecode.Operand
	if opts := ctx.OptionsClause(); opts != nil {
		optsReg = c.CompileOptionsClause(opts)
	}

	var timeoutReg bytecode.Operand
	if timeout := ctx.TimeoutClause(); timeout != nil {
		timeoutReg = c.CompileTimeoutClauseContext(timeout)
	}

	c.ctx.Emitter.WithSpan(span, func() {
		c.ctx.Emitter.EmitMove(streamReg, srcReg)
		c.ctx.Emitter.EmitABC(bytecode.OpStream, streamReg, eventReg, optsReg)
		c.ctx.Emitter.EmitABC(bytecode.OpStreamIter, streamReg, streamReg, timeoutReg)
	})

	start := c.ctx.Emitter.NewLabel()
	iterationDone := c.ctx.Emitter.NewLabel()
	cleanup := c.ctx.Emitter.NewLabel()
	timeoutHandler := c.ctx.Emitter.NewLabel("waitfor", "event", "timeout")
	errorHandler := c.ctx.Emitter.NewLabel("waitfor", "event", "error")
	end := c.ctx.Emitter.NewLabel("waitfor", "event", "end")

	c.ctx.Emitter.MarkLabel(start)
	c.ctx.Emitter.WithSpan(span, func() {
		c.ctx.Emitter.EmitIterNextTimeout(streamReg, timeoutStateReg, iterationDone)
	})

	if filter := ctx.EventFilterClause(); filter != nil {
		eventValReg, _ := c.ctx.Symbols.DeclareLocal(core.PseudoVariable, core.TypeUnknown)

		c.ctx.Emitter.WithSpan(span, func() {
			c.ctx.Emitter.EmitAB(bytecode.OpIterValue, eventValReg, streamReg)
		})

		cond := c.ctx.ExprCompiler.compileWithImplicitCurrent(filter.Expression())
		c.ctx.Emitter.EmitJumpIfFalse(cond, start)
	}

	c.ctx.Emitter.EmitJump(cleanup)
	c.ctx.Emitter.MarkLabel(iterationDone)
	c.ctx.Emitter.EmitJump(cleanup)

	c.ctx.Emitter.MarkLabel(cleanup)
	c.ctx.Emitter.WithSpan(span, func() {
		c.ctx.Emitter.EmitA(bytecode.OpClose, streamReg)
	})

	endCatchExclusive := c.ctx.Emitter.Size()

	c.ctx.Emitter.EmitJumpIfTrue(timeoutStateReg, timeoutHandler)
	c.ctx.Emitter.EmitJumpIfTrue(errorStateReg, errorHandler)
	c.ctx.Emitter.EmitJump(end)

	errorPreludePC := -1
	if plan.OnError != nil && (plan.OnError.ActionKind == core.RecoveryActionReturn || plan.OnError.ActionKind == core.RecoveryActionRetry) {
		errorPreludePC = c.ctx.Emitter.Size()
		c.ctx.Emitter.EmitBoolean(errorStateReg, true)
		c.ctx.Emitter.EmitBoolean(timeoutStateReg, false)
		c.ctx.Emitter.EmitJump(cleanup)
	}

	c.ctx.Emitter.MarkLabel(timeoutHandler)
	switch {
	case plan.OnTimeout != nil && plan.OnTimeout.ActionKind == core.RecoveryActionReturn:
		fallback := c.ctx.ExprCompiler.Compile(plan.OnTimeout.Expr)
		emitMoveAuto(c.ctx, resultReg, ensureRecoveryRegister(c.ctx, fallback))
		c.ctx.Emitter.EmitJump(end)
	default:
		c.ctx.Emitter.Emit(bytecode.OpFailTimeout)
	}

	c.ctx.Emitter.MarkLabel(errorHandler)
	switch {
	case plan.OnError != nil && plan.OnError.ActionKind == core.RecoveryActionReturn:
		fallback := c.ctx.ExprCompiler.Compile(plan.OnError.Expr)
		emitMoveAuto(c.ctx, resultReg, ensureRecoveryRegister(c.ctx, fallback))
		c.ctx.Emitter.EmitJump(end)
	case recoveryHandlerRetries(plan.OnError):
		c.ctx.Emitter.EmitBoolean(errorStateReg, false)
		c.ctx.Emitter.EmitBoolean(timeoutStateReg, false)

		retriesAvailableReg := c.ctx.Registers.Allocate()
		c.ctx.Emitter.EmitGt(retriesAvailableReg, retriesRemainingReg, zeroReg)

		onExhausted := c.ctx.Emitter.NewLabel("waitfor", "event", "exhausted")
		c.ctx.Emitter.EmitJumpIfFalse(retriesAvailableReg, onExhausted)
		c.ctx.Emitter.EmitA(bytecode.OpDecr, retriesRemainingReg)
		emitRecoveryRetryDelay(c.ctx, plan.OnError.Retry, retryDelayState)

		if plan.OnError.Retry.FinalActionKind == core.RecoveryActionReturn {
			c.ctx.Emitter.EmitJump(retryStartLabel)
		} else {
			moreProtectedReg := c.ctx.Registers.Allocate()
			c.ctx.Emitter.EmitGt(moreProtectedReg, retriesRemainingReg, zeroReg)
			c.ctx.Emitter.EmitJumpIfTrue(moreProtectedReg, retryStartLabel)
			c.ctx.Emitter.EmitJump(finalAttemptLabel)
		}

		c.ctx.Emitter.MarkLabel(onExhausted)
		if plan.OnError.Retry.FinalActionKind == core.RecoveryActionReturn {
			fallback := c.ctx.ExprCompiler.Compile(plan.OnError.Retry.FinalExpr)
			emitMoveAuto(c.ctx, resultReg, ensureRecoveryRegister(c.ctx, fallback))
			c.ctx.Emitter.EmitJump(end)
		} else {
			c.ctx.Emitter.EmitJump(finalAttemptLabel)
		}
	}

	if errorPreludePC >= 0 && endCatchExclusive > startCatch {
		c.ctx.CatchTable.Push(startCatch, endCatchExclusive-1, errorPreludePC)
	}

	if recoveryHandlerRetries(plan.OnError) && plan.OnError.Retry.FinalActionKind != core.RecoveryActionReturn {
		c.ctx.Emitter.MarkLabel(finalAttemptLabel)
		finalOut := c.compileEventWithPlan(ctx, core.RecoveryPlan{OnTimeout: plan.OnTimeout})
		if finalOut != bytecode.NoopOperand && finalOut != resultReg {
			emitMoveAuto(c.ctx, resultReg, ensureRecoveryRegister(c.ctx, finalOut))
		}
	}

	c.ctx.Emitter.MarkLabel(end)

	return resultReg
}

func (c *WaitCompiler) compilePredicate(ctx fql.IWaitForPredicateExpressionContext) bytecode.Operand {
	predicate := ctx.WaitForPredicate()
	if predicate == nil {
		return bytecode.NoopOperand
	}

	predExpr := predicate.Expression()
	if predExpr == nil {
		return bytecode.NoopOperand
	}

	if legacy := legacyWaitForOrThrowNode(predExpr); legacy != nil {
		c.ctx.Errors.Add(c.ctx.Errors.Create(parser.SyntaxError, legacy, "Unexpected THROW after OR in WAITFOR predicate"))
		return bytecode.NoopOperand
	}

	config := c.buildWaitPredicateConfig(ctx, predicate, predExpr)
	c.normalizeWaitPredicateConfig(&config)

	if result, ok := c.tryCompileWaitPredicateFastPath(config); ok {
		return result
	}

	state := c.initWaitPredicatePollState(config)
	c.emitWaitPredicatePollLoop(config, state)

	return state.resultReg
}

func (c *WaitCompiler) compilePredicateWithRecovery(ctx fql.IWaitForPredicateExpressionContext, plan core.RecoveryPlan) bytecode.Operand {
	predicate := ctx.WaitForPredicate()
	if predicate == nil {
		return bytecode.NoopOperand
	}

	predExpr := predicate.Expression()
	if predExpr == nil {
		return bytecode.NoopOperand
	}

	if legacy := legacyWaitForOrThrowNode(predExpr); legacy != nil {
		c.ctx.Errors.Add(c.ctx.Errors.Create(parser.SyntaxError, legacy, "Unexpected THROW after OR in WAITFOR predicate"))
		return bytecode.NoopOperand
	}

	config := c.buildWaitPredicateConfig(ctx, predicate, predExpr)
	c.normalizeWaitPredicateConfig(&config)

	if recoveryHandlerRetries(plan.OnError) {
		return c.compilePredicateWithRetry(config, ctx, plan)
	}

	state := c.initWaitPredicatePollState(config)

	return c.emitWaitPredicatePollLoopWithRecovery(config, state, plan)
}

func (c *WaitCompiler) compilePredicateWithRetry(config waitPredicateCompileConfig, ctx fql.IWaitForPredicateExpressionContext, plan core.RecoveryPlan) bytecode.Operand {
	zeroReg := loadConstant(c.ctx, runtime.ZeroInt)
	retriesRemainingReg := loadConstant(c.ctx, runtime.NewInt(plan.OnError.Retry.Count))
	retryDelayState := initRetryDelayState(c.ctx, plan.OnError.Retry)
	retryStart := c.ctx.Emitter.NewLabel("waitfor", "predicate", "retry")
	finalAttemptLabel := c.ctx.Emitter.NewLabel("waitfor", "predicate", "final")
	end := c.ctx.Emitter.NewLabel("waitfor", "predicate", "end")

	c.ctx.Emitter.MarkLabel(retryStart)
	state := c.initWaitPredicatePollState(config)

	start := c.ctx.Emitter.NewLabel()
	success := c.ctx.Emitter.NewLabel()
	protectedTimeout := c.ctx.Emitter.NewLabel()
	timeoutHandler := c.ctx.Emitter.NewLabel("waitfor", "predicate", "timeout")

	startCatch := c.ctx.Emitter.Size()

	c.ctx.Emitter.MarkLabel(start)

	valueReg := c.ctx.ExprCompiler.Compile(config.predExpr)
	condReg := c.emitWaitPredicateCondition(config.mode, valueReg)
	c.ctx.Emitter.EmitJumpIfTrue(condReg, success)

	elapsedReg := c.emitWaitPredicateTimeoutCheck(config.timeoutReg, state.startReg, state.unitReg, protectedTimeout)
	sleepIntervalReg := c.prepareWaitSleepInterval(config, state.pollReg)
	c.emitWaitSleep(sleepIntervalReg, config.timeoutReg, elapsedReg)

	if config.backoff != core.RetryBackoffNone {
		c.emitBackoffUpdate(config.backoff, state.intervalReg, state.baseEveryReg)
		if config.capEveryReg != bytecode.NoopOperand {
			c.emitClampMax(state.intervalReg, config.capEveryReg)
		}
	}

	c.ctx.Emitter.EmitJump(start)
	c.ctx.Emitter.MarkLabel(success)
	c.emitWaitSuccessResult(config.mode, state.resultReg, valueReg)
	c.ctx.Emitter.EmitJump(end)

	c.ctx.Emitter.MarkLabel(protectedTimeout)
	c.ctx.Emitter.EmitJump(timeoutHandler)

	endCatchExclusive := c.ctx.Emitter.Size()

	errorHandlerPC := c.ctx.Emitter.Size()
	retriesAvailableReg := c.ctx.Registers.Allocate()
	c.ctx.Emitter.EmitGt(retriesAvailableReg, retriesRemainingReg, zeroReg)

	onExhausted := c.ctx.Emitter.NewLabel("waitfor", "predicate", "exhausted")
	c.ctx.Emitter.EmitJumpIfFalse(retriesAvailableReg, onExhausted)
	c.ctx.Emitter.EmitA(bytecode.OpDecr, retriesRemainingReg)
	emitRecoveryRetryDelay(c.ctx, plan.OnError.Retry, retryDelayState)

	if plan.OnError.Retry.FinalActionKind == core.RecoveryActionReturn {
		c.ctx.Emitter.EmitJump(retryStart)
	} else {
		moreProtectedReg := c.ctx.Registers.Allocate()
		c.ctx.Emitter.EmitGt(moreProtectedReg, retriesRemainingReg, zeroReg)
		c.ctx.Emitter.EmitJumpIfTrue(moreProtectedReg, retryStart)
		c.ctx.Emitter.EmitJump(finalAttemptLabel)
	}

	c.ctx.Emitter.MarkLabel(onExhausted)
	if plan.OnError.Retry.FinalActionKind == core.RecoveryActionReturn {
		fallback := c.ctx.ExprCompiler.Compile(plan.OnError.Retry.FinalExpr)
		emitMoveAuto(c.ctx, state.resultReg, ensureRecoveryRegister(c.ctx, fallback))
		c.ctx.Emitter.EmitJump(end)
	} else {
		c.ctx.Emitter.EmitJump(finalAttemptLabel)
	}

	c.ctx.Emitter.MarkLabel(timeoutHandler)
	switch {
	case plan.OnTimeout != nil && plan.OnTimeout.ActionKind == core.RecoveryActionReturn:
		fallback := c.ctx.ExprCompiler.Compile(plan.OnTimeout.Expr)
		emitMoveAuto(c.ctx, state.resultReg, ensureRecoveryRegister(c.ctx, fallback))
		c.ctx.Emitter.EmitJump(end)
	case plan.OnTimeout != nil && plan.OnTimeout.ActionKind == core.RecoveryActionFail:
		c.ctx.Emitter.Emit(bytecode.OpFailTimeout)
	default:
		c.emitWaitTimeoutResult(config.mode, state.resultReg)
		c.ctx.Emitter.EmitJump(end)
	}

	if endCatchExclusive > startCatch {
		c.ctx.CatchTable.Push(startCatch, endCatchExclusive-1, errorHandlerPC)
	}

	if plan.OnError.Retry.FinalActionKind != core.RecoveryActionReturn {
		c.ctx.Emitter.MarkLabel(finalAttemptLabel)
		finalOut := c.compilePredicateWithPlan(ctx, core.RecoveryPlan{OnTimeout: plan.OnTimeout})
		if finalOut != bytecode.NoopOperand && finalOut != state.resultReg {
			emitMoveAuto(c.ctx, state.resultReg, ensureRecoveryRegister(c.ctx, finalOut))
		}
	}

	c.ctx.Emitter.MarkLabel(end)

	return state.resultReg
}

func legacyWaitForOrThrowNode(expr fql.IExpressionContext) antlr.ParserRuleContext {
	if expr == nil || expr.LogicalOrOperator() == nil {
		return nil
	}

	return bareThrowExpressionNode(expr.GetRight())
}

func bareThrowExpressionNode(expr fql.IExpressionContext) antlr.ParserRuleContext {
	if expr == nil {
		return nil
	}

	if expr.UnaryOperator() != nil || expr.LogicalAndOperator() != nil || expr.LogicalOrOperator() != nil || expr.GetTernaryOperator() != nil {
		return nil
	}

	return bareThrowPredicateNode(expr.Predicate())
}

func bareThrowPredicateNode(pred fql.IPredicateContext) antlr.ParserRuleContext {
	if pred == nil {
		return nil
	}

	if pred.EqualityOperator() != nil || pred.ArrayOperator() != nil || pred.InOperator() != nil || pred.LikeOperator() != nil {
		return nil
	}

	return bareThrowAtomNode(pred.ExpressionAtom())
}

func bareThrowAtomNode(atom fql.IExpressionAtomContext) antlr.ParserRuleContext {
	if atom == nil {
		return nil
	}

	if atom.MultiplicativeOperator() != nil || atom.AdditiveOperator() != nil || atom.RegexpOperator() != nil {
		return nil
	}

	variable := atom.Variable()
	if variable == nil || !strings.EqualFold(matchVariableName(variable), "THROW") {
		return nil
	}

	node, ok := variable.(antlr.ParserRuleContext)
	if !ok {
		return nil
	}

	return node
}

func resolveWaitPredicateMode(hasValue, hasExists, hasNot bool) waitForPredicateMode {
	if hasValue {
		return waitForPredicateModeValue
	}

	if hasExists {
		if hasNot {
			return waitForPredicateModeNotExists
		}

		return waitForPredicateModeExists
	}

	return waitForPredicateModeBool
}

func (c *WaitCompiler) buildWaitPredicateConfig(ctx fql.IWaitForPredicateExpressionContext, predicate fql.IWaitForPredicateContext, predExpr fql.IExpressionContext) waitPredicateCompileConfig {
	everyReg, capEveryReg := c.compileEveryClause(ctx.EveryClause())
	jitterReg, jitterLiteral, hasJitter := c.compileJitterClause(ctx.JitterClause())

	return waitPredicateCompileConfig{
		mode:          resolveWaitPredicateMode(predicate.Value() != nil, predicate.Exists() != nil, predicate.Not() != nil),
		predExpr:      predExpr,
		timeoutReg:    c.compileDurationClause(ctx.TimeoutClause()),
		everyReg:      everyReg,
		capEveryReg:   capEveryReg,
		backoff:       c.compileBackoffClause(ctx.BackoffClause()),
		jitterReg:     jitterReg,
		jitterLiteral: jitterLiteral,
		hasJitter:     hasJitter,
	}
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
		c.emitClampRange(config.jitterReg, loadConstant(c.ctx, runtime.NewFloat(0)), loadConstant(c.ctx, runtime.NewFloat(1)))
	}
}

func (c *WaitCompiler) tryCompileWaitPredicateFastPath(config waitPredicateCompileConfig) (bytecode.Operand, bool) {
	switch config.mode {
	case waitForPredicateModeBool:
		truth, ok := literalTruthinessFromExpression(config.predExpr)
		if !ok {
			return bytecode.NoopOperand, false
		}

		if truth {
			return c.emitImmediateWaitBool(true), true
		}

		if config.timeoutReg != bytecode.NoopOperand {
			c.ctx.Emitter.EmitA(bytecode.OpSleep, config.timeoutReg)
			return c.emitImmediateWaitBool(false), true
		}

		return bytecode.NoopOperand, false
	default:
		exists, ok := literalExistsFromExpression(config.predExpr)
		if !ok {
			return bytecode.NoopOperand, false
		}

		cond := exists
		if config.mode == waitForPredicateModeNotExists {
			cond = !exists
		}

		if cond {
			if config.mode == waitForPredicateModeValue {
				return c.ctx.ExprCompiler.Compile(config.predExpr), true
			}

			return c.emitImmediateWaitBool(true), true
		}

		if config.timeoutReg != bytecode.NoopOperand {
			c.ctx.Emitter.EmitA(bytecode.OpSleep, config.timeoutReg)
			if config.mode == waitForPredicateModeValue {
				return c.emitImmediateWaitNone(), true
			}

			return c.emitImmediateWaitBool(false), true
		}

		return bytecode.NoopOperand, false
	}
}

func (c *WaitCompiler) emitImmediateWaitBool(value bool) bytecode.Operand {
	resultReg := c.ctx.Registers.Allocate()
	c.ctx.Emitter.EmitBoolean(resultReg, value)

	return resultReg
}

func (c *WaitCompiler) emitImmediateWaitNone() bytecode.Operand {
	resultReg := c.ctx.Registers.Allocate()
	c.ctx.Emitter.EmitLoadNone(resultReg)

	return resultReg
}

func (c *WaitCompiler) initWaitPredicatePollState(config waitPredicateCompileConfig) waitPredicatePollState {
	state := waitPredicatePollState{
		baseEveryReg: c.ctx.Registers.Allocate(),
	}

	if config.everyReg != bytecode.NoopOperand {
		c.ctx.Emitter.EmitMove(state.baseEveryReg, config.everyReg)
	} else {
		c.ctx.Emitter.EmitLoadConst(state.baseEveryReg, c.ctx.Symbols.AddConstant(runtime.NewInt(waitForDefaultEveryMs)))
	}

	state.pollReg = state.baseEveryReg
	if config.backoff != core.RetryBackoffNone {
		state.intervalReg = c.ctx.Registers.Allocate()
		c.ctx.Emitter.EmitMove(state.intervalReg, state.baseEveryReg)
		state.pollReg = state.intervalReg
	}

	state.resultReg = c.ctx.Registers.Allocate()
	if config.mode == waitForPredicateModeValue {
		c.ctx.Emitter.EmitLoadNone(state.resultReg)
	} else {
		c.ctx.Emitter.EmitBoolean(state.resultReg, false)
	}

	if config.timeoutReg != bytecode.NoopOperand {
		state.startReg = c.emitNow()
		state.unitReg = loadConstant(c.ctx, runtime.NewString("f"))
	}

	return state
}

func (c *WaitCompiler) emitWaitPredicatePollLoop(config waitPredicateCompileConfig, state waitPredicatePollState) {
	start := c.ctx.Emitter.NewLabel()
	success := c.ctx.Emitter.NewLabel()
	timeoutLabel := c.ctx.Emitter.NewLabel()
	end := c.ctx.Emitter.NewLabel()

	c.ctx.Emitter.MarkLabel(start)

	valueReg := c.ctx.ExprCompiler.Compile(config.predExpr)
	condReg := c.emitWaitPredicateCondition(config.mode, valueReg)
	c.ctx.Emitter.EmitJumpIfTrue(condReg, success)

	elapsedReg := c.emitWaitPredicateTimeoutCheck(config.timeoutReg, state.startReg, state.unitReg, timeoutLabel)
	sleepIntervalReg := c.prepareWaitSleepInterval(config, state.pollReg)
	c.emitWaitSleep(sleepIntervalReg, config.timeoutReg, elapsedReg)

	if config.backoff != core.RetryBackoffNone {
		c.emitBackoffUpdate(config.backoff, state.intervalReg, state.baseEveryReg)
		if config.capEveryReg != bytecode.NoopOperand {
			c.emitClampMax(state.intervalReg, config.capEveryReg)
		}
	}

	c.ctx.Emitter.EmitJump(start)
	c.ctx.Emitter.MarkLabel(success)
	c.emitWaitSuccessResult(config.mode, state.resultReg, valueReg)
	c.ctx.Emitter.EmitJump(end)

	c.ctx.Emitter.MarkLabel(timeoutLabel)
	c.emitWaitTimeoutResult(config.mode, state.resultReg)
	c.ctx.Emitter.MarkLabel(end)
}

func (c *WaitCompiler) emitWaitPredicatePollLoopWithRecovery(config waitPredicateCompileConfig, state waitPredicatePollState, plan core.RecoveryPlan) bytecode.Operand {
	start := c.ctx.Emitter.NewLabel()
	success := c.ctx.Emitter.NewLabel()
	protectedTimeout := c.ctx.Emitter.NewLabel()
	timeoutHandler := c.ctx.Emitter.NewLabel("waitfor", "predicate", "timeout")
	end := c.ctx.Emitter.NewLabel("waitfor", "predicate", "end")

	startCatch := c.ctx.Emitter.Size()

	c.ctx.Emitter.MarkLabel(start)

	valueReg := c.ctx.ExprCompiler.Compile(config.predExpr)
	condReg := c.emitWaitPredicateCondition(config.mode, valueReg)
	c.ctx.Emitter.EmitJumpIfTrue(condReg, success)

	elapsedReg := c.emitWaitPredicateTimeoutCheck(config.timeoutReg, state.startReg, state.unitReg, protectedTimeout)
	sleepIntervalReg := c.prepareWaitSleepInterval(config, state.pollReg)
	c.emitWaitSleep(sleepIntervalReg, config.timeoutReg, elapsedReg)

	if config.backoff != core.RetryBackoffNone {
		c.emitBackoffUpdate(config.backoff, state.intervalReg, state.baseEveryReg)
		if config.capEveryReg != bytecode.NoopOperand {
			c.emitClampMax(state.intervalReg, config.capEveryReg)
		}
	}

	c.ctx.Emitter.EmitJump(start)
	c.ctx.Emitter.MarkLabel(success)
	c.emitWaitSuccessResult(config.mode, state.resultReg, valueReg)
	c.ctx.Emitter.EmitJump(end)

	c.ctx.Emitter.MarkLabel(protectedTimeout)
	c.ctx.Emitter.EmitJump(timeoutHandler)

	endCatchExclusive := c.ctx.Emitter.Size()

	errorHandlerPC := -1
	if plan.OnError != nil && plan.OnError.ActionKind == core.RecoveryActionReturn {
		errorHandlerPC = c.ctx.Emitter.Size()
		fallback := c.ctx.ExprCompiler.Compile(plan.OnError.Expr)
		emitMoveAuto(c.ctx, state.resultReg, ensureRecoveryRegister(c.ctx, fallback))
		c.ctx.Emitter.EmitJump(end)
	}

	c.ctx.Emitter.MarkLabel(timeoutHandler)
	switch {
	case plan.OnTimeout != nil && plan.OnTimeout.ActionKind == core.RecoveryActionReturn:
		fallback := c.ctx.ExprCompiler.Compile(plan.OnTimeout.Expr)
		emitMoveAuto(c.ctx, state.resultReg, ensureRecoveryRegister(c.ctx, fallback))
		c.ctx.Emitter.EmitJump(end)
	case plan.OnTimeout != nil && plan.OnTimeout.ActionKind == core.RecoveryActionFail:
		c.ctx.Emitter.Emit(bytecode.OpFailTimeout)
	default:
		c.emitWaitTimeoutResult(config.mode, state.resultReg)
		c.ctx.Emitter.EmitJump(end)
	}

	c.ctx.Emitter.MarkLabel(end)

	if errorHandlerPC >= 0 && endCatchExclusive > startCatch {
		c.ctx.CatchTable.Push(startCatch, endCatchExclusive-1, errorHandlerPC)
	}

	return state.resultReg
}

func (c *WaitCompiler) emitWaitPredicateCondition(mode waitForPredicateMode, valueReg bytecode.Operand) bytecode.Operand {
	switch mode {
	case waitForPredicateModeValue, waitForPredicateModeExists:
		return c.emitExistsCheck(valueReg)
	case waitForPredicateModeNotExists:
		existsReg := c.emitExistsCheck(valueReg)
		condReg := c.ctx.Registers.Allocate()
		c.ctx.Emitter.EmitAB(bytecode.OpNot, condReg, existsReg)

		return condReg
	default:
		condReg := c.ctx.Registers.Allocate()
		c.ctx.Emitter.EmitAB(bytecode.OpCastBool, condReg, valueReg)

		return condReg
	}
}

func (c *WaitCompiler) emitWaitPredicateTimeoutCheck(timeoutReg, startReg, unitReg bytecode.Operand, timeoutLabel core.Label) bytecode.Operand {
	if timeoutReg == bytecode.NoopOperand {
		return bytecode.NoopOperand
	}

	nowReg := c.emitNow()
	elapsedReg := c.emitDateDiff(startReg, nowReg, unitReg)
	reachedReg := c.ctx.Registers.Allocate()
	c.ctx.Emitter.EmitGte(reachedReg, elapsedReg, timeoutReg)
	c.ctx.Emitter.EmitJumpIfTrue(reachedReg, timeoutLabel)

	return elapsedReg
}

func (c *WaitCompiler) prepareWaitSleepInterval(config waitPredicateCompileConfig, pollReg bytecode.Operand) bytecode.Operand {
	if !config.hasJitter && config.capEveryReg == bytecode.NoopOperand {
		return pollReg
	}

	sleepIntervalReg := c.ctx.Registers.Allocate()
	c.ctx.Emitter.EmitMove(sleepIntervalReg, pollReg)

	if config.hasJitter {
		c.emitApplyJitter(sleepIntervalReg, config.jitterReg)
	}

	if config.capEveryReg != bytecode.NoopOperand {
		c.emitClampMax(sleepIntervalReg, config.capEveryReg)
	}

	return sleepIntervalReg
}

func (c *WaitCompiler) emitWaitSuccessResult(mode waitForPredicateMode, resultReg, valueReg bytecode.Operand) {
	if mode == waitForPredicateModeValue {
		c.ctx.Emitter.EmitMove(resultReg, valueReg)
		return
	}

	c.ctx.Emitter.EmitBoolean(resultReg, true)
}

func (c *WaitCompiler) emitWaitTimeoutResult(mode waitForPredicateMode, resultReg bytecode.Operand) {
	if mode == waitForPredicateModeValue {
		c.ctx.Emitter.EmitLoadNone(resultReg)
		return
	}

	c.ctx.Emitter.EmitBoolean(resultReg, false)
}

func (c *WaitCompiler) emitExistsCheck(val bytecode.Operand) bytecode.Operand {
	dst := c.ctx.Registers.Allocate()
	c.ctx.Emitter.EmitAB(bytecode.OpExists, dst, val)
	c.ctx.Types.Set(dst, core.TypeBool)

	return dst
}

func (c *WaitCompiler) emitNow() bytecode.Operand {
	return c.ctx.ExprCompiler.CompileFunctionCallByNameWith(nil, runtime.NewString("NOW"), false, nil)
}

func (c *WaitCompiler) emitDateDiff(start, end, unit bytecode.Operand) bytecode.Operand {
	return c.emitFunctionCall(runtime.NewString("DATE_DIFF"), start, end, unit)
}

func (c *WaitCompiler) emitFunctionCall(name runtime.String, args ...bytecode.Operand) bytecode.Operand {
	if len(args) == 0 {
		return c.ctx.ExprCompiler.CompileFunctionCallByNameWith(nil, name, false, nil)
	}

	seq := c.ctx.Registers.AllocateSequence(len(args))

	for i, arg := range args {
		c.ctx.Emitter.EmitMove(seq[i], arg)
		c.ctx.Types.Set(seq[i], operandType(c.ctx, arg))
	}

	return c.ctx.ExprCompiler.CompileFunctionCallByNameWith(nil, name, false, seq)
}

func (c *WaitCompiler) emitWaitSleep(intervalReg, timeoutReg, elapsedReg bytecode.Operand) {
	if timeoutReg == bytecode.NoopOperand {
		c.ctx.Emitter.EmitA(bytecode.OpSleep, intervalReg)
		return
	}

	sleepReg := c.ctx.Registers.Allocate()
	c.ctx.Emitter.EmitMove(sleepReg, intervalReg)

	remainingReg := c.ctx.Registers.Allocate()
	c.ctx.Emitter.EmitABC(bytecode.OpSub, remainingReg, timeoutReg, elapsedReg)

	shouldTrim := c.ctx.Registers.Allocate()
	c.ctx.Emitter.EmitLt(shouldTrim, remainingReg, sleepReg)

	useRemaining := c.ctx.Emitter.NewLabel()
	continueSleep := c.ctx.Emitter.NewLabel()

	c.ctx.Emitter.EmitJumpIfTrue(shouldTrim, useRemaining)
	c.ctx.Emitter.EmitJump(continueSleep)

	c.ctx.Emitter.MarkLabel(useRemaining)
	c.ctx.Emitter.EmitMove(sleepReg, remainingReg)
	c.ctx.Emitter.MarkLabel(continueSleep)

	c.ctx.Emitter.EmitA(bytecode.OpSleep, sleepReg)
}

func (c *WaitCompiler) emitBackoffUpdate(strategy core.RetryBackoff, intervalReg, baseEveryReg bytecode.Operand) {
	switch strategy {
	case core.RetryBackoffLinear:
		c.ctx.Emitter.EmitABC(bytecode.OpAdd, intervalReg, intervalReg, baseEveryReg)
	case core.RetryBackoffExponential:
		twoReg := loadConstant(c.ctx, runtime.NewInt(2))
		c.ctx.Emitter.EmitABC(bytecode.OpMul, intervalReg, intervalReg, twoReg)
	default:
		return
	}
}

func (c *WaitCompiler) emitClampMin(target, min bytecode.Operand) {
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

func (c *WaitCompiler) emitClampMax(target, max bytecode.Operand) {
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

func (c *WaitCompiler) emitClampRange(target, min, max bytecode.Operand) {
	c.emitClampMin(target, min)
	c.emitClampMax(target, max)
}

func (c *WaitCompiler) emitApplyJitter(intervalReg, jitterReg bytecode.Operand) {
	if intervalReg == bytecode.NoopOperand || jitterReg == bytecode.NoopOperand {
		return
	}

	randReg := c.ctx.Registers.Allocate()
	c.ctx.Emitter.EmitA(bytecode.OpRand, randReg)

	twoReg := loadConstant(c.ctx, runtime.NewFloat(2))
	oneReg := loadConstant(c.ctx, runtime.NewFloat(1))

	twoJitterReg := c.ctx.Registers.Allocate()
	c.ctx.Emitter.EmitABC(bytecode.OpMul, twoJitterReg, jitterReg, twoReg)

	randScaleReg := c.ctx.Registers.Allocate()
	c.ctx.Emitter.EmitABC(bytecode.OpMul, randScaleReg, randReg, twoJitterReg)

	oneMinusReg := c.ctx.Registers.Allocate()
	c.ctx.Emitter.EmitABC(bytecode.OpSub, oneMinusReg, oneReg, jitterReg)

	multiplierReg := c.ctx.Registers.Allocate()
	c.ctx.Emitter.EmitABC(bytecode.OpAdd, multiplierReg, oneMinusReg, randScaleReg)

	c.ctx.Emitter.EmitABC(bytecode.OpMul, intervalReg, intervalReg, multiplierReg)
}

func waitForSpan(src antlr.RuleContext, fallback antlr.RuleContext) source.Span {
	span := source.Span{Start: -1, End: -1}

	if src != nil {
		if prc, ok := src.(antlr.ParserRuleContext); ok {
			span = parser.SpanFromRuleContext(prc)
			return span
		}
	}

	if fallback != nil {
		if prc, ok := fallback.(antlr.ParserRuleContext); ok {
			span = parser.SpanFromRuleContext(prc)
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
		newOperandBranch(sl != nil, func() bytecode.Operand { return c.ctx.LiteralCompiler.CompileStringLiteral(sl) }),
		newOperandBranch(v != nil, func() bytecode.Operand { return c.ctx.ExprCompiler.CompileVariable(v) }),
		newOperandBranch(p != nil, func() bytecode.Operand { return c.ctx.ExprCompiler.CompileParam(p) }),
		newOperandBranch(me != nil, func() bytecode.Operand { return c.ctx.ExprCompiler.CompileMemberExpression(me) }),
		newOperandBranch(fce != nil, func() bytecode.Operand { return c.ctx.ExprCompiler.CompileFunctionCall(fce, false) }),
	)
}

// CompileWaitForEventSource processes the event source expression in a WAITFOR statement.
func (c *WaitCompiler) CompileWaitForEventSource(ctx fql.IWaitForEventSourceContext) bytecode.Operand {
	v := ctx.Variable()
	me := ctx.MemberExpression()
	fce := ctx.FunctionCallExpression()

	return compileFirstOperand(
		newOperandBranch(v != nil, func() bytecode.Operand { return c.ctx.ExprCompiler.CompileVariable(v) }),
		newOperandBranch(me != nil, func() bytecode.Operand { return c.ctx.ExprCompiler.CompileMemberExpression(me) }),
		newOperandBranch(fce != nil, func() bytecode.Operand { return c.ctx.ExprCompiler.CompileFunctionCallExpression(fce) }),
	)
}

// CompileOptionsClause processes the options clause in a WAITFOR statement.
func (c *WaitCompiler) CompileOptionsClause(ctx fql.IOptionsClauseContext) bytecode.Operand {
	if ol := ctx.ObjectLiteral(); ol != nil {
		return c.ctx.LiteralCompiler.CompileObjectLiteral(ol)
	}

	return bytecode.NoopOperand
}

// CompileTimeoutClauseContext processes the timeout clause in a WAITFOR statement.
func (c *WaitCompiler) CompileTimeoutClauseContext(ctx fql.ITimeoutClauseContext) bytecode.Operand {
	return c.compileDurationClause(ctx)
}

func (c *WaitCompiler) compileEveryClause(ctx fql.IEveryClauseContext) (bytecode.Operand, bytecode.Operand) {
	if ctx == nil {
		return bytecode.NoopOperand, bytecode.NoopOperand
	}

	values := ctx.AllEveryClauseValue()
	if len(values) == 0 {
		return bytecode.NoopOperand, bytecode.NoopOperand
	}

	base := c.compileDurationClause(values[0])
	if len(values) > 1 {
		return base, c.compileDurationClause(values[1])
	}

	return base, bytecode.NoopOperand
}

func (c *WaitCompiler) compileDurationClause(ctx core.DurationClause) bytecode.Operand {
	if ctx == nil {
		return bytecode.NoopOperand
	}

	if dl := ctx.DurationLiteral(); dl != nil {
		val, err := parseDurationLiteral(dl.GetText())
		if err != nil {
			panic(err)
		}

		return loadConstant(c.ctx, val)
	}

	il := ctx.IntegerLiteral()
	fl := ctx.FloatLiteral()
	v := ctx.Variable()
	p := ctx.Param()
	me := ctx.MemberExpression()
	fc := ctx.FunctionCall()

	return compileFirstOperand(
		newOperandBranch(il != nil, func() bytecode.Operand { return c.ctx.LiteralCompiler.CompileIntegerLiteral(il) }),
		newOperandBranch(fl != nil, func() bytecode.Operand { return c.ctx.LiteralCompiler.CompileFloatLiteral(fl) }),
		newOperandBranch(v != nil, func() bytecode.Operand { return c.ctx.ExprCompiler.CompileVariable(v) }),
		newOperandBranch(p != nil, func() bytecode.Operand { return c.ctx.ExprCompiler.CompileParam(p) }),
		newOperandBranch(me != nil, func() bytecode.Operand { return c.ctx.ExprCompiler.CompileMemberExpression(me) }),
		newOperandBranch(fc != nil, func() bytecode.Operand { return c.ctx.ExprCompiler.CompileFunctionCall(fc, false) }),
	)
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
			err := c.ctx.Errors.Create(parser.SemanticError, prc, "JITTER must be between 0 and 1")
			err.Hint = "Use a value between 0 and 1, e.g. JITTER 0.2."
			c.ctx.Errors.Add(err)
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
		newOperandBranch(fl != nil, func() bytecode.Operand { return c.ctx.LiteralCompiler.CompileFloatLiteral(fl) }),
		newOperandBranch(il != nil, func() bytecode.Operand { return c.ctx.LiteralCompiler.CompileIntegerLiteral(il) }),
		newOperandBranch(v != nil, func() bytecode.Operand { return c.ctx.ExprCompiler.CompileVariable(v) }),
		newOperandBranch(p != nil, func() bytecode.Operand { return c.ctx.ExprCompiler.CompileParam(p) }),
		newOperandBranch(me != nil, func() bytecode.Operand { return c.ctx.ExprCompiler.CompileMemberExpression(me) }),
		newOperandBranch(fc != nil, func() bytecode.Operand { return c.ctx.ExprCompiler.CompileFunctionCall(fc, false) }),
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
				err := c.ctx.Errors.Create(parser.SemanticError, prc, "BACKOFF strategy must be a constant string")
				err.Hint = "Use one of: NONE, LINEAR, EXPONENTIAL."
				c.ctx.Errors.Add(err)
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
			err := c.ctx.Errors.Create(parser.SemanticError, prc, "Unknown BACKOFF strategy")
			err.Hint = "Use one of: NONE, LINEAR, EXPONENTIAL."
			c.ctx.Errors.Add(err)
		}

		return core.RetryBackoffNone
	}
}

func parseDurationLiteral(text string) (runtime.Value, error) {
	raw := normalizeDurationLiteral(text)
	if raw == "" {
		return runtime.None, strconv.ErrSyntax
	}

	number, unit, ok := splitDurationLiteral(raw)
	if !ok || number == "" {
		return runtime.None, strconv.ErrSyntax
	}

	value, err := parseDurationLiteralNumber(number)
	if err != nil {
		return runtime.None, err
	}

	multiplier, ok := durationUnitMultiplier(unit)
	if !ok {
		return runtime.None, strconv.ErrSyntax
	}

	return durationValueFromMilliseconds(value * multiplier), nil
}

func normalizeDurationLiteral(text string) string {
	return strings.ToUpper(strings.TrimSpace(text))
}

func splitDurationLiteral(raw string) (string, string, bool) {
	switch {
	case strings.HasSuffix(raw, "MS"):
		return strings.TrimSuffix(raw, "MS"), "MS", true
	case strings.HasSuffix(raw, "S"):
		return strings.TrimSuffix(raw, "S"), "S", true
	case strings.HasSuffix(raw, "M"):
		return strings.TrimSuffix(raw, "M"), "M", true
	case strings.HasSuffix(raw, "H"):
		return strings.TrimSuffix(raw, "H"), "H", true
	case strings.HasSuffix(raw, "D"):
		return strings.TrimSuffix(raw, "D"), "D", true
	default:
		return "", "", false
	}
}

func parseDurationLiteralNumber(raw string) (float64, error) {
	return strconv.ParseFloat(raw, 64)
}

func durationUnitMultiplier(unit string) (float64, bool) {
	switch unit {
	case "MS":
		return 1, true
	case "S":
		return 1000, true
	case "M":
		return 60000, true
	case "H":
		return 3600000, true
	case "D":
		return 86400000, true
	default:
		return 0, false
	}
}

func durationValueFromMilliseconds(ms float64) runtime.Value {
	if frac := math.Mod(ms, 1); frac == 0 {
		return runtime.NewInt64(int64(ms))
	}

	return runtime.NewFloat(ms)
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
		if str, ok := parseStringLiteralConst(lit.StringLiteral()); ok {
			return str.String() != "", true
		}
		return false, false
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
		if str, ok := parseStringLiteralConst(lit.StringLiteral()); ok {
			return str.String() != "", true
		}
		return false, false
	default:
		return true, true
	}
}
