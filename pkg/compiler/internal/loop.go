package internal

import (
	"github.com/antlr4-go/antlr/v4"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/compiler/internal/core"
	parser "github.com/MontFerret/ferret/v2/pkg/parser/diagnostics"
	"github.com/MontFerret/ferret/v2/pkg/parser/fql"
	"github.com/MontFerret/ferret/v2/pkg/source"
)

type (
	// LoopCompiler handles the compilation of FOR loop expressions in FQL queries.
	// It transforms loop operations into VM instructions for iteration, filtering, and data manipulation.
	LoopCompiler struct {
		ctx      *CompilationSession
		bindings *BindingCompiler
		collects *CollectCompiler
		exprs    *ExprCompiler
		literals *LiteralCompiler
		recovery *RecoveryCompiler
		sorts    *LoopSortCompiler
		facts    *TypeFacts
	}

	loopOperandKind int

	loopResultKind uint8

	loopResultSpec struct {
		returnCtx      fql.IReturnExpressionContext
		passThroughCtx fql.IForExpressionContext
		kind           loopResultKind
		distinct       bool
		// bodyEnd is the exclusive end of ordinary body entries compiled before finalization.
		bodyEnd int
	}

	loopOperandContext struct {
		param                  fql.IParamContext
		integerLiteral         fql.IIntegerLiteralContext
		variable               fql.IVariableContext
		memberExpression       fql.IMemberExpressionContext
		implicitCurrentExpr    fql.IImplicitCurrentExpressionContext
		implicitMemberExpr     fql.IImplicitMemberExpressionContext
		functionCallExpression fql.IFunctionCallExpressionContext
		rangeOperator          fql.IRangeOperatorContext
		arrayLiteral           fql.IArrayLiteralContext
		objectLiteral          fql.IObjectLiteralContext
	}
)

const (
	loopOperandParam loopOperandKind = iota
	loopOperandIntegerLiteral
	loopOperandVariable
	loopOperandMemberExpression
	loopOperandImplicitCurrent
	loopOperandImplicitMember
	loopOperandFunctionCallExpression
	loopOperandRangeOperator
	loopOperandArrayLiteral
	loopOperandObjectLiteral
)

const (
	loopResultEffectOnly loopResultKind = iota
	loopResultCollecting
	loopResultPassThrough
)

// NewLoopCompiler creates a new instance of LoopCompiler with the given compiler context.
func NewLoopCompiler(ctx *CompilationSession) *LoopCompiler {
	return &LoopCompiler{ctx: ctx}
}

func (c *LoopCompiler) bind(
	bindings *BindingCompiler,
	collects *CollectCompiler,
	exprs *ExprCompiler,
	literals *LiteralCompiler,
	recovery *RecoveryCompiler,
	sorts *LoopSortCompiler,
	facts *TypeFacts,
) {
	if c == nil {
		return
	}

	c.bindings = bindings
	c.collects = collects
	c.exprs = exprs
	c.literals = literals
	c.recovery = recovery
	c.sorts = sorts
	c.facts = facts
}

// Compile processes a FOR expression from the FQL AST and generates the appropriate VM instructions.
// It determines whether to compile a FOR IN loop (iteration over a collection)
// or a FOR WHILE/DO WHILE loop.
// Returns an operand representing the destination of the loop results.
func (c *LoopCompiler) Compile(ctx fql.IForExpressionContext) bytecode.Operand {
	return c.compileWithResultUse(ctx, resultRequired)
}

// CompileDiscarded evaluates a FOR expression without retaining its result.
func (c *LoopCompiler) CompileDiscarded(ctx fql.IForExpressionContext) bytecode.Operand {
	return c.compileWithResultUse(ctx, resultDiscarded)
}

// CompileWithOuterRecoveryPlan is the supported cross-compiler entrypoint for
// FOR expressions that need their recovery tails merged with an outer plan.
func (c *LoopCompiler) CompileWithOuterRecoveryPlan(ctx fql.IForExpressionContext, outerPlan core.RecoveryPlan) bytecode.Operand {
	return c.compileWithResultUseAndOuterRecoveryPlan(ctx, resultRequired, outerPlan)
}

func (c *LoopCompiler) compileWithResultUse(ctx fql.IForExpressionContext, use resultUse) bytecode.Operand {
	return c.compileWithResultUseAndOuterRecoveryPlan(ctx, use, core.RecoveryPlan{})
}

func (c *LoopCompiler) compileWithResultUseAndOuterRecoveryPlan(
	ctx fql.IForExpressionContext,
	use resultUse,
	outerPlan core.RecoveryPlan,
) bytecode.Operand {
	if ctx == nil {
		return bytecode.NoopOperand
	}

	resultSpec := c.resolveLoopResultSpec(ctx)
	if resultSpec.kind == loopResultEffectOnly && use == resultRequired {
		c.reportReturnlessForExpression(ctx)
		use = resultDiscarded
	}

	if outerPlan.OnError == nil && outerPlan.OnTimeout == nil {
		return c.compilePlain(ctx, use, resultSpec)
	}

	return c.recovery.CompileOperation(c.newLoopOperationRecoverySpec(ctx, use, resultSpec, outerPlan))
}

func (c *LoopCompiler) compilePlain(ctx fql.IForExpressionContext, use resultUse, resultSpec loopResultSpec) bytecode.Operand {
	var initialized bool

	if ctx.In() != nil {
		initialized = c.compileInitialization(ctx, core.ForInLoop, use, resultSpec)
	} else if ctx.Do() == nil {
		initialized = c.compileInitialization(ctx, core.WhileLoop, use, resultSpec)
	} else {
		initialized = c.compileInitialization(ctx, core.DoWhileLoop, use, resultSpec)
	}

	if !initialized {
		return bytecode.NoopOperand
	}

	c.compileLoopBody(ctx, resultSpec)

	return c.compileFinalization(resultSpec)
}

func (c *LoopCompiler) newLoopOperationRecoverySpec(
	ctx fql.IForExpressionContext,
	use resultUse,
	resultSpec loopResultSpec,
	outerPlan core.RecoveryPlan,
) OperationRecoverySpec {
	recoverySpec := OperationRecoverySpec{
		OuterPlan: outerPlan,
		CompilePlain: func() bytecode.Operand {
			return c.compileRecoveryAttempt(ctx, use, resultSpec)
		},
		CompileSuppressed: func() bytecode.Operand {
			return c.compilePlain(ctx, use, resultSpec)
		},
		CompileFinalAttempt: func(plan core.RecoveryPlan) bytecode.Operand {
			out := c.compilePlain(ctx, use, resultSpec)
			if recoveryPlanHasReturnHandler(plan) {
				return c.ensureRecoveryResult(out, use)
			}

			return out
		},
	}

	if ctx != nil && ctx.In() != nil {
		recoverySpec.BuildProtected = func(recoveryLabel, timeoutLabel, endLabel core.Label) ProtectedRecoveryRegion {
			return c.buildProtectedForInRecovery(ctx, use, resultSpec, recoveryLabel, timeoutLabel, endLabel)
		}
	}

	return recoverySpec
}

func (c *LoopCompiler) compileRecoveryAttempt(
	ctx fql.IForExpressionContext,
	use resultUse,
	resultSpec loopResultSpec,
) bytecode.Operand {
	out := c.compilePlain(ctx, use, resultSpec)

	return c.ensureRecoveryResult(out, use)
}

func (c *LoopCompiler) ensureRecoveryResult(out bytecode.Operand, use resultUse) bytecode.Operand {
	if use == resultRequired || out != bytecode.NoopOperand {
		return out
	}

	result := c.ctx.Function.Registers.Allocate()
	c.ctx.Program.Emitter.EmitA(bytecode.OpLoadNone, result)
	c.ctx.Function.Types.Set(result, core.TypeNone)

	return result
}

func (c *LoopCompiler) buildProtectedForInRecovery(
	ctx fql.IForExpressionContext,
	use resultUse,
	resultSpec loopResultSpec,
	recoveryLabel, _ core.Label, endLabel core.Label,
) ProtectedRecoveryRegion {
	errorStateReg := c.ctx.Function.Registers.Allocate()
	c.ctx.Program.Emitter.EmitBoolean(errorStateReg, false)

	startCatch := c.ctx.Program.Emitter.Size()
	if !c.compileInitialization(ctx, core.ForInLoop, use, resultSpec) {
		return ProtectedRecoveryRegion{Result: bytecode.NoopOperand}
	}

	loop := c.ctx.Function.Loops.Current()
	breakLabel := loop.BreakLabel()

	c.compileLoopBody(ctx, resultSpec)

	out := c.ensureRecoveryResult(c.compileFinalization(resultSpec), use)
	endCatchExclusive := c.ctx.Program.Emitter.Size()

	routeRecovery := c.ctx.Program.Emitter.NewLabel("recovery", "for", "route")
	c.ctx.Program.Emitter.EmitJumpIfTrue(errorStateReg, routeRecovery)
	c.ctx.Program.Emitter.EmitJump(endLabel)

	errorPreludePC := c.ctx.Program.Emitter.Size()
	c.ctx.Program.Emitter.EmitBoolean(errorStateReg, true)
	c.ctx.Program.Emitter.EmitJump(breakLabel)

	c.ctx.Program.Emitter.MarkLabel(routeRecovery)
	c.ctx.Program.Emitter.EmitBoolean(errorStateReg, false)
	c.ctx.Program.Emitter.EmitJump(recoveryLabel)

	return ProtectedRecoveryRegion{
		Result:            out,
		StartCatch:        startCatch,
		EndCatchExclusive: endCatchExclusive,
		CatchHandlerPC:    errorPreludePC,
		HasTimeout:        false,
	}
}

// compileInitialization handles loop setup and enters the loop's symbol scope.
func (c *LoopCompiler) compileInitialization(
	ctx fql.IForExpressionContext,
	kind core.LoopKind,
	use resultUse,
	resultSpec loopResultSpec,
) bool {
	if !c.validateLoopBindingPattern(ctx) {
		return false
	}

	loopType := core.NormalLoop
	if resultSpec.kind == loopResultPassThrough {
		loopType = core.PassThroughLoop
	}

	loop := c.ctx.Function.Loops.NewLoop(kind, loopType, resultSpec.distinct, use == resultRequired)
	c.setLoopDestinationType(loop)

	c.configureLoopRuntime(loop, ctx, kind)

	// Push the loop onto the stack and enter a new symbol scope
	c.ctx.Function.Loops.Push(loop)
	c.ctx.Function.Symbols.EnterScope()

	if c.ctx.Program.Semantics != nil {
		c.ctx.Program.Semantics.EnterScope(parser.SpanFromRuleContext(ctx.(antlr.ParserRuleContext)))
	}

	valueType, keyType := c.inferLoopVariableTypes(ctx, loop, kind)
	c.declareLoopVariables(ctx, loop, valueType, keyType)
	c.emitLoopInitialization(ctx, loop)
	c.patchDistinctLoopDestination(loop)

	return true
}

// resolveLoopResultSpec classifies syntax independently of whether the caller
// retains the result. Braced terminal loops parse as body statements, so the
// final one is peeled off to preserve legacy pass-through semantics.
func (c *LoopCompiler) resolveLoopResultSpec(ctx fql.IForExpressionContext) loopResultSpec {
	bodies := ctx.AllForExpressionBody()
	result := loopResultSpec{
		kind:    loopResultEffectOnly,
		bodyEnd: len(bodies),
	}

	if re := ctx.ReturnExpression(); re != nil {
		result.returnCtx = re
		result.kind = loopResultCollecting
		result.distinct = re.Distinct() != nil

		return result
	}

	if terminal := ctx.ForExpressionReturn(); terminal != nil {
		if re := terminal.ReturnExpression(); re != nil {
			result.returnCtx = re
			result.kind = loopResultCollecting
			result.distinct = re.Distinct() != nil

			return result
		}

		if nested := terminal.ForExpression(); nested != nil {
			result.passThroughCtx = nested
			result.kind = loopResultPassThrough

			return result
		}
	}

	if ctx.OpenBrace() != nil && len(bodies) > 0 {
		last := bodies[len(bodies)-1]
		if stmt := last.ForExpressionStatement(); stmt != nil {
			if nested := stmt.ForExpression(); nested != nil {
				result.passThroughCtx = nested
				result.kind = loopResultPassThrough
				result.bodyEnd = len(bodies) - 1
			}
		}
	}

	return result
}

func (c *LoopCompiler) reportReturnlessForExpression(ctx fql.IForExpressionContext) {
	rule, ok := ctx.(antlr.ParserRuleContext)
	if !ok {
		return
	}

	err := c.ctx.Program.Errors.Create(
		parser.SemanticError,
		rule,
		"A FOR loop used as an expression must return a value.",
	)
	err.Hint = "Add RETURN to the loop body, or use the loop as a statement."
	if len(err.Spans) > 0 {
		err.Spans[0].Label = "returnless FOR expression"
	}

	c.ctx.Program.Errors.Add(err)
}

func (c *LoopCompiler) setLoopDestinationType(loop *core.Loop) {
	if loop != nil && loop.Dst.IsRegister() {
		c.ctx.Function.Types.Set(loop.Dst, core.TypeList)
	}
}

func (c *LoopCompiler) configureLoopRuntime(loop *core.Loop, ctx fql.IForExpressionContext, kind core.LoopKind) {
	switch kind {
	case core.ForInLoop:
		loop.Src = c.compileForExpressionSource(ctx.ForExpressionSource())
	default:
		loop.ConditionFn = func() bytecode.Operand {
			return c.exprs.Compile(ctx.Expression())
		}
	}
}

func (c *LoopCompiler) inferLoopVariableTypes(ctx fql.IForExpressionContext, loop *core.Loop, kind core.LoopKind) (core.ValueType, core.ValueType) {
	switch kind {
	case core.ForInLoop:
		return c.inferForInTypes(ctx.ForExpressionSource(), loop.Src)
	case core.WhileLoop, core.DoWhileLoop:
		return core.TypeInt, core.TypeUnknown
	default:
		return core.TypeUnknown, core.TypeUnknown
	}
}

func (c *LoopCompiler) declareLoopVariables(ctx fql.IForExpressionContext, loop *core.Loop, valueType, keyType core.ValueType) {
	c.declareLoopValueVariable(ctx, loop, valueType)
	c.declareLoopCounterVariable(ctx, loop, keyType)
}

func (c *LoopCompiler) declareLoopValueVariable(ctx fql.IForExpressionContext, loop *core.Loop, valueType core.ValueType) {
	if pattern := ctx.GetValuePattern(); pattern != nil {
		loop.Value = c.ctx.Function.Registers.Allocate()
		loop.Destructured = true
		c.ctx.Function.Types.Set(loop.Value, valueType)

		return
	}

	val := ctx.GetValueVariable()
	if val == nil {
		return
	}

	varName := textOfLoopVariable(val)
	valueCtx := val.(antlr.ParserRuleContext)
	bindingID := bindingIDFromRule(valueCtx)
	declared := loop.DeclareValueVarWithOptions(varName, c.ctx.Function.Symbols, valueType, core.BindingOptions{ID: bindingID})

	if declared && c.ctx.Program.Semantics != nil {
		span := parser.SpanFromRuleContext(valueCtx)
		activation := span.End
		if src := ctx.ForExpressionSource(); src != nil {
			activation = parser.SpanFromRuleContext(src.(antlr.ParserRuleContext)).End
		}

		c.ctx.Program.Semantics.RecordBinding(
			bindingID,
			varName,
			SemanticSymbolLoopBinding,
			span,
			span,
			false,
			valueType,
			activation,
			c.ctx.Program.Semantics.CurrentFunctionSymbol(),
			0,
		)
	}

	if loop.Value.IsRegister() {
		c.ctx.Function.Types.Set(loop.Value, valueType)
	}
}

func (c *LoopCompiler) declareLoopCounterVariable(ctx fql.IForExpressionContext, loop *core.Loop, keyType core.ValueType) {
	ctr := ctx.GetCounterVariable()
	if ctr == nil {
		return
	}

	counterCtx := ctr.(antlr.ParserRuleContext)
	bindingID := bindingIDFromRule(counterCtx)
	name := textOfBindingIdentifier(ctr)
	declared := loop.DeclareKeyVarWithOptions(name, c.ctx.Function.Symbols, keyType, core.BindingOptions{ID: bindingID})

	if declared && c.ctx.Program.Semantics != nil {
		span := parser.SpanFromRuleContext(counterCtx)
		activation := span.End

		if src := ctx.ForExpressionSource(); src != nil {
			activation = parser.SpanFromRuleContext(src.(antlr.ParserRuleContext)).End
		}

		c.ctx.Program.Semantics.RecordBinding(
			bindingID,
			name,
			SemanticSymbolLoopBinding,
			span,
			span,
			false,
			keyType,
			activation,
			c.ctx.Program.Semantics.CurrentFunctionSymbol(),
			0,
		)
	}
	if loop.Key.IsRegister() {
		c.ctx.Function.Types.Set(loop.Key, keyType)
	}
}

func (c *LoopCompiler) emitLoopInitialization(ctx fql.IForExpressionContext, loop *core.Loop) {
	span := source.Span{Start: -1, End: -1}

	if srcCtx := ctx.ForExpressionSource(); srcCtx != nil {
		if prc, ok := srcCtx.(antlr.ParserRuleContext); ok {
			span = parser.SpanFromRuleContext(prc)
		}
	} else if prc, ok := ctx.(antlr.ParserRuleContext); ok {
		span = parser.SpanFromRuleContext(prc)
	}

	c.ctx.Program.Emitter.WithSpan(span, func() {
		loop.EmitInitialization(c.ctx.Function.Registers, c.ctx.Program.Emitter)
	})

	if loop.Destructured {
		pattern := ctx.GetValuePattern()
		declaration := source.Span{}
		activation := 0

		if c.ctx.Program.Semantics != nil {
			declaration = parser.SpanFromRuleContext(pattern.(antlr.ParserRuleContext))
			activation = declaration.End

			if src := ctx.ForExpressionSource(); src != nil {
				activation = parser.SpanFromRuleContext(src.(antlr.ParserRuleContext)).End
			}
		}

		c.bindings.compileDestructuringPattern(
			pattern,
			loop.Value,
			false,
			SemanticSymbolLoopBinding,
			declaration,
			activation,
		)
	}
}

func (c *LoopCompiler) validateLoopBindingPattern(ctx fql.IForExpressionContext) bool {
	if ctx == nil || ctx.GetValuePattern() == nil {
		return true
	}

	leaves := structuredBindingPatternLeaves(ctx.GetValuePattern())
	duplicate, first, ok := duplicateBindingPatternLeaf(leaves)
	if !ok {
		return true
	}

	c.ctx.Program.Errors.DuplicateDestructuringBinding(duplicate.Context, first.Context, duplicate.Name)

	return false
}

func (c *LoopCompiler) patchDistinctLoopDestination(loop *core.Loop) {
	if !loop.CollectResult || loop.Allocate || !loop.Distinct {
		return
	}

	parent := c.ctx.Function.Loops.RequiredParent(c.ctx.Function.Loops.Depth())
	c.ctx.Program.Emitter.Patchx(parent.StartLabel(), 1)
}

// compileFinalization evaluates any result-producing terminal and exits the loop scope.
func (c *LoopCompiler) compileFinalization(resultSpec loopResultSpec) bytecode.Operand {
	loop := c.ctx.Function.Loops.Current()

	switch resultSpec.kind {
	case loopResultCollecting:
		// Normal loops always evaluate the return expression, but only retained
		// results are appended to the loop destination.
		re := resultSpec.returnCtx
		returnUse := resultDiscarded

		if loop.CollectResult {
			returnUse = resultRequired
		}

		compileReturn := func() {
			value := re.ReturnValue()
			if value == nil {
				return
			}

			var (
				result    bytecode.Operand
				resultCtx antlr.ParserRuleContext
			)

			if expr := value.Expression(); expr != nil {
				result = c.exprs.compileWithResultUse(expr, returnUse)
				resultCtx = expr.(antlr.ParserRuleContext)
			} else if nested := value.ForExpression(); nested != nil {
				result = c.compileWithResultUse(nested, returnUse)
				resultCtx = nested.(antlr.ParserRuleContext)
			}

			span := parser.SpanFromRuleContext(re)
			if resultCtx != nil {
				span = parser.SpanFromRuleContext(resultCtx)
			}

			if loop.CollectResult {
				c.ctx.Program.Emitter.WithSpan(span, func() {
					c.ctx.Program.Emitter.EmitAB(bytecode.OpPush, loop.Dst, result)
				})
			}
		}

		if loop.CollectResult {
			c.ctx.WithDebugPointKind(re, bytecode.DebugPointReturn, compileReturn)
		} else {
			c.ctx.WithRetainedDebugPointKind(re, bytecode.DebugPointReturn, compileReturn)
		}
	case loopResultPassThrough:
		use := resultDiscarded
		if loop.CollectResult {
			use = resultRequired
		}

		c.compileWithResultUse(resultSpec.passThroughCtx, use)
	}

	// Emit VM instructions for loop finalization
	loop.EmitFinalization(c.ctx.Program.Emitter)

	// Clean up the symbol scope and pop the loop from the stack
	c.ctx.Function.Symbols.ExitScope()

	if c.ctx.Program.Semantics != nil {
		c.ctx.Program.Semantics.ExitScope()
	}

	c.ctx.Function.Loops.Pop()

	return loop.Dst
}

// compileForExpressionSource processes the source expression for a FOR IN loop.
func (c *LoopCompiler) compileForExpressionSource(ctx fql.IForExpressionSourceContext) bytecode.Operand {
	if ctx == nil {
		return bytecode.NoopOperand
	}

	return c.exprs.Compile(ctx.Expression())
}

func (c *LoopCompiler) compileLoopBody(ctx fql.IForExpressionContext, resultSpec loopResultSpec) {
	if ctx == nil {
		return
	}

	body := ctx.AllForExpressionBody()

	for _, entry := range body[:resultSpec.bodyEnd] {
		if statement := entry.ForExpressionStatement(); statement != nil {
			c.compileForExpressionStatement(statement)
		} else if clause := entry.ForExpressionClause(); clause != nil {
			c.compileForExpressionClause(clause)
		}
	}
}

// compileForExpressionStatement processes statements within a FOR loop body.
// These can be declarations, assignments, deletes, expressions, or nested loops.
// The results of these statements are not used directly in the loop result.
func (c *LoopCompiler) compileForExpressionStatement(ctx fql.IForExpressionStatementContext) {
	rule, _ := ctx.(antlr.ParserRuleContext)
	c.ctx.WithDebugPoint(rule, func() {
		c.compileForExpressionStatementInner(ctx)
	})
}

func (c *LoopCompiler) compileForExpressionStatementInner(ctx fql.IForExpressionStatementContext) {
	// Handle variable declarations (e.g., LET x = 1)
	if vd := ctx.VariableDeclaration(); vd != nil {
		_ = c.bindings.CompileVariableDeclaration(vd)
	} else if as := ctx.AssignmentStatement(); as != nil {
		_ = c.bindings.CompileAssignmentStatement(as)
	} else if ds := ctx.DeleteStatement(); ds != nil {
		_ = c.bindings.CompileDeleteStatement(ds)
	} else if fe := ctx.ForExpression(); fe != nil {
		_ = c.CompileDiscarded(fe)
	} else if es := ctx.ExpressionStatement(); es != nil {
		_ = c.exprs.CompileDiscarded(es.Expression())
	}
}

// compileForExpressionClause processes clauses within a FOR loop body.
// These can be LIMIT, FILTER, SORT, or COLLECT clauses that modify the loop behavior.
// Each clause type is delegated to a specific compilation method.
func (c *LoopCompiler) compileForExpressionClause(ctx fql.IForExpressionClauseContext) {
	rule, _ := ctx.(antlr.ParserRuleContext)
	c.ctx.WithDebugPoint(rule, func() {
		c.compileForExpressionClauseInner(ctx)
	})
}

func (c *LoopCompiler) compileForExpressionClauseInner(ctx fql.IForExpressionClauseContext) {
	// Handle LIMIT clause (e.g., LIMIT 10)
	if lc := ctx.LimitClause(); lc != nil {
		c.compileLimitClause(lc)
	} else if fc := ctx.FilterClause(); fc != nil {
		// Handle FILTER clause (e.g., FILTER x > 5)
		c.compileFilterClause(fc)
	} else if sc := ctx.SortClause(); sc != nil {
		// Handle SORT clause (e.g., SORT x DESC)
		c.compileSortClause(sc)
	} else if cc := ctx.CollectClause(); cc != nil {
		// Handle COLLECT clause (e.g., COLLECT x = y)
		c.compileCollectClause(cc)
	}
}

// compileLimitClause processes a LIMIT clause in a FOR loop.
// It handles both simple LIMIT clauses and LIMIT with OFFSET clauses.
// For a single value, it's treated as a limit. For two values, the first is offset and the second is limit.
func (c *LoopCompiler) compileLimitClause(ctx fql.ILimitClauseContext) {
	clauses := ctx.AllLimitClauseValue()

	if len(clauses) == 1 {
		// Simple LIMIT clause (e.g., LIMIT 10)
		c.compileLimit(c.compileLimitClauseValue(clauses[0]))
	} else {
		// LIMIT with OFFSET clause (e.g., LIMIT 5, 10 - offset 5, limit 10)
		c.compileOffset(c.compileLimitClauseValue(clauses[0]))
		c.compileLimit(c.compileLimitClauseValue(clauses[1]))
	}
}

// compileLimitClauseValue processes a value in a LIMIT clause.
// It handles various types of expressions that can be used as limit or offset values,
// such as parameters, integer literals, variables, member expressions, and function calls.
// Returns an operand representing the compiled limit/offset value.
func (c *LoopCompiler) compileLimitClauseValue(ctx fql.ILimitClauseValueContext) bytecode.Operand {
	return c.compileLoopOperand(
		loopOperandContext{
			param:                  ctx.Param(),
			integerLiteral:         ctx.IntegerLiteral(),
			variable:               ctx.Variable(),
			memberExpression:       ctx.MemberExpression(),
			implicitCurrentExpr:    ctx.ImplicitCurrentExpression(),
			implicitMemberExpr:     ctx.ImplicitMemberExpression(),
			functionCallExpression: ctx.FunctionCallExpression(),
		},
		loopOperandParam,
		loopOperandIntegerLiteral,
		loopOperandVariable,
		loopOperandMemberExpression,
		loopOperandImplicitCurrent,
		loopOperandImplicitMember,
		loopOperandFunctionCallExpression,
	)
}

func (c *LoopCompiler) compileLoopOperand(source loopOperandContext, order ...loopOperandKind) bytecode.Operand {
	branches := make([]operandBranch, 0, len(order))

	for _, kind := range order {
		switch kind {
		case loopOperandParam:
			branches = append(branches, newOperandBranch(source.param != nil, func() bytecode.Operand { return c.exprs.CompileParam(source.param) }))
		case loopOperandIntegerLiteral:
			branches = append(branches, newOperandBranch(source.integerLiteral != nil, func() bytecode.Operand { return c.literals.CompileIntegerLiteral(source.integerLiteral) }))
		case loopOperandVariable:
			branches = append(branches, newOperandBranch(source.variable != nil, func() bytecode.Operand { return c.exprs.CompileVariable(source.variable) }))
		case loopOperandMemberExpression:
			branches = append(branches, newOperandBranch(source.memberExpression != nil, func() bytecode.Operand { return c.exprs.CompileMemberExpression(source.memberExpression) }))
		case loopOperandImplicitCurrent:
			branches = append(branches, newOperandBranch(source.implicitCurrentExpr != nil, func() bytecode.Operand {
				return c.exprs.CompileImplicitCurrentExpression(source.implicitCurrentExpr)
			}))
		case loopOperandImplicitMember:
			branches = append(branches, newOperandBranch(source.implicitMemberExpr != nil, func() bytecode.Operand {
				return c.exprs.CompileImplicitMemberExpression(source.implicitMemberExpr)
			}))
		case loopOperandFunctionCallExpression:
			branches = append(branches, newOperandBranch(source.functionCallExpression != nil, func() bytecode.Operand {
				return c.exprs.CompileFunctionCallExpression(source.functionCallExpression)
			}))
		case loopOperandRangeOperator:
			branches = append(branches, newOperandBranch(source.rangeOperator != nil, func() bytecode.Operand { return c.exprs.CompileRangeOperator(source.rangeOperator) }))
		case loopOperandArrayLiteral:
			branches = append(branches, newOperandBranch(source.arrayLiteral != nil, func() bytecode.Operand { return c.literals.CompileArrayLiteral(source.arrayLiteral) }))
		case loopOperandObjectLiteral:
			branches = append(branches, newOperandBranch(source.objectLiteral != nil, func() bytecode.Operand { return c.literals.CompileObjectLiteral(source.objectLiteral) }))
		}
	}

	return compileFirstOperand(branches...)
}

// compileLimit emits VM instructions to limit the number of iterations in a loop.
// It allocates a state register and emits an iterator limit instruction with the loop's end label.
func (c *LoopCompiler) compileLimit(src bytecode.Operand) {
	// Allocate a state register for the limit operation
	state := c.ctx.Function.Registers.Allocate()
	c.ctx.Function.Loops.Current().RegisterReset(state)
	// Emit the iterator limit instruction with the loop's end label
	c.ctx.Program.Emitter.EmitIterLimit(state, src, c.ctx.Function.Loops.Current().BreakLabel())
}

// compileOffset emits VM instructions to skip a number of iterations at the start of a loop.
// It allocates a state register and emits an iterator skip instruction with the loop's jump label.
func (c *LoopCompiler) compileOffset(src bytecode.Operand) {
	// Allocate a state register for the offset operation
	state := c.ctx.Function.Registers.Allocate()
	c.ctx.Function.Loops.Current().RegisterReset(state)
	// Emit the iterator skip instruction with the loop's jump label
	c.ctx.Program.Emitter.EmitIterSkip(state, src, c.ctx.Function.Loops.Current().ContinueLabel())
}

// compileFilterClause processes a FILTER clause in a FOR loop.
// It compiles the filter expression and emits a conditional jump instruction
// that skips the current iteration if the filter condition is false.
func (c *LoopCompiler) compileFilterClause(ctx fql.IFilterClauseContext) {
	// Compile the filter expression (e.g., FILTER x > 5)
	// Get the jump label for the current loop
	label := c.ctx.Function.Loops.Current().ContinueLabel()
	// Emit a jump instruction that skips to the next iteration if the filter condition is false
	c.exprs.EmitConditionJump(ctx.Expression(), label, false)
}

// compileSortClause processes a SORT clause in a FOR loop.
// It delegates the compilation to the specialized LoopSortCompiler.
func (c *LoopCompiler) compileSortClause(ctx fql.ISortClauseContext) {
	// Delegate to the specialized sort compiler
	c.sorts.Compile(ctx)
}

// compileCollectClause processes a COLLECT clause in a FOR loop.
// It delegates the compilation to the specialized CollectCompiler.
func (c *LoopCompiler) compileCollectClause(ctx fql.ICollectClauseContext) {
	// Delegate to the specialized collect compiler
	c.collects.Compile(ctx)
}

func (c *LoopCompiler) inferForInTypes(srcCtx fql.IForExpressionSourceContext, src bytecode.Operand) (core.ValueType, core.ValueType) {
	if srcCtx == nil {
		return core.TypeUnknown, core.TypeUnknown
	}

	return c.inferForInExpressionTypes(srcCtx.Expression(), src)
}

func (c *LoopCompiler) inferForInExpressionTypes(ctx fql.IExpressionContext, src bytecode.Operand) (core.ValueType, core.ValueType) {
	atom := c.sourceExpressionAtom(ctx)

	if atom == nil {
		return c.inferValueKeyFromCollection(c.facts.OperandType(src))
	}

	if atom.RangeOperator() != nil {
		return core.TypeInt, core.TypeInt
	}

	if lit := atom.Literal(); lit != nil {
		switch {
		case lit.ArrayLiteral() != nil:
			return c.inferArrayLiteralElementType(lit.ArrayLiteral()), core.TypeInt
		case lit.ObjectLiteral() != nil:
			return core.TypeAny, core.TypeString
		}
	}

	if v := atom.Variable(); v != nil {
		if binding, ok := c.ctx.Function.Symbols.ResolveBinding(v.GetText()); ok {
			return c.inferValueKeyFromCollection(binding.Type)
		}
	}

	if atom.Param() != nil || atom.FunctionCallExpression() != nil {
		return core.TypeAny, core.TypeAny
	}

	if atom.MemberExpression() != nil {
		return c.inferValueKeyFromCollection(c.facts.OperandType(src))
	}

	return c.inferValueKeyFromCollection(c.facts.OperandType(src))
}

func (c *LoopCompiler) sourceExpressionAtom(ctx fql.IExpressionContext) fql.IExpressionAtomContext {
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

	if nested := atom.Expression(); nested != nil {
		return c.sourceExpressionAtom(nested)
	}

	return atom
}

func (c *LoopCompiler) inferValueKeyFromCollection(typ core.ValueType) (core.ValueType, core.ValueType) {
	switch typ {
	case core.TypeList, core.TypeArray:
		return core.TypeAny, core.TypeInt
	case core.TypeMap, core.TypeObject:
		return core.TypeAny, core.TypeString
	case core.TypeAny:
		return core.TypeAny, core.TypeAny
	default:
		return core.TypeUnknown, core.TypeUnknown
	}
}

func (c *LoopCompiler) inferArrayLiteralElementType(ctx fql.IArrayLiteralContext) core.ValueType {
	if ctx == nil {
		return core.TypeUnknown
	}

	entries := ctx.AllArrayEntry()
	if len(entries) == 0 {
		return core.TypeUnknown
	}

	elemType := core.TypeUnknown

	for _, entry := range entries {
		if entry.SpreadEntry() != nil {
			return core.TypeAny
		}

		typ := c.inferExpressionType(entry.Expression())
		if typ == core.TypeUnknown {
			return core.TypeAny
		}

		if elemType == core.TypeUnknown {
			elemType = typ

			continue
		}

		if elemType != typ {
			return core.TypeAny
		}
	}

	return elemType
}

func (c *LoopCompiler) inferExpressionType(ctx fql.IExpressionContext) core.ValueType {
	if ctx == nil {
		return core.TypeUnknown
	}

	if p := ctx.Predicate(); p != nil {
		return c.inferPredicateType(p)
	}

	return core.TypeUnknown
}

func (c *LoopCompiler) inferPredicateType(ctx fql.IPredicateContext) core.ValueType {
	if ctx == nil {
		return core.TypeUnknown
	}

	if atom := ctx.ExpressionAtom(); atom != nil {
		return c.inferExpressionAtomType(atom)
	}

	return core.TypeUnknown
}

func (c *LoopCompiler) inferExpressionAtomType(ctx fql.IExpressionAtomContext) core.ValueType {
	if ctx == nil {
		return core.TypeUnknown
	}

	if lit := ctx.Literal(); lit != nil {
		return c.facts.LiteralType(lit)
	}

	if v := ctx.Variable(); v != nil {
		if binding, ok := c.ctx.Function.Symbols.ResolveBinding(v.GetText()); ok {
			return binding.Type
		}
		return core.TypeUnknown
	}

	if ctx.Param() != nil || ctx.FunctionCallExpression() != nil {
		return core.TypeAny
	}

	if ctx.MatchExpression() != nil {
		return core.TypeAny
	}

	if ctx.RangeOperator() != nil {
		return core.TypeList
	}

	if ctx.ForExpression() != nil {
		return core.TypeList
	}

	return core.TypeUnknown
}
