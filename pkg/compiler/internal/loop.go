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
	return c.CompileWithOuterRecoveryPlan(ctx, core.RecoveryPlan{})
}

// CompileWithOuterRecoveryPlan is the supported cross-compiler entrypoint for
// FOR expressions that need their recovery tails merged with an outer plan.
func (c *LoopCompiler) CompileWithOuterRecoveryPlan(ctx fql.IForExpressionContext, outerPlan core.RecoveryPlan) bytecode.Operand {
	if ctx == nil {
		return bytecode.NoopOperand
	}

	if outerPlan.OnError == nil && outerPlan.OnTimeout == nil {
		return c.compilePlain(ctx)
	}

	return c.recovery.CompileOperation(c.newLoopOperationRecoverySpec(ctx, outerPlan))
}

func (c *LoopCompiler) compilePlain(ctx fql.IForExpressionContext) bytecode.Operand {
	var returnRuleCtx antlr.RuleContext

	if ctx.In() != nil {
		returnRuleCtx = c.compileInitialization(ctx, core.ForInLoop)
	} else if ctx.Do() == nil {
		returnRuleCtx = c.compileInitialization(ctx, core.WhileLoop)
	} else {
		returnRuleCtx = c.compileInitialization(ctx, core.DoWhileLoop)
	}

	// Probably, a syntax error happened and no return rule context was created.
	if returnRuleCtx == nil {
		return bytecode.NoopOperand
	}

	c.compileLoopBody(ctx)

	// Finalize the loop and return the destination operand
	return c.compileFinalization(returnRuleCtx)
}

func (c *LoopCompiler) newLoopOperationRecoverySpec(ctx fql.IForExpressionContext, outerPlan core.RecoveryPlan) OperationRecoverySpec {
	spec := OperationRecoverySpec{
		OuterPlan: outerPlan,
		CompilePlain: func() bytecode.Operand {
			return c.compilePlain(ctx)
		},
	}

	if ctx != nil && ctx.In() != nil {
		spec.BuildProtected = func(recoveryLabel, timeoutLabel, endLabel core.Label) ProtectedRecoveryRegion {
			return c.buildProtectedForInRecovery(ctx, recoveryLabel, timeoutLabel, endLabel)
		}
	}

	return spec
}

func (c *LoopCompiler) buildProtectedForInRecovery(
	ctx fql.IForExpressionContext,
	recoveryLabel, _ core.Label, endLabel core.Label,
) ProtectedRecoveryRegion {
	errorStateReg := c.ctx.Registers.Allocate()
	c.ctx.Emitter.EmitBoolean(errorStateReg, false)

	startCatch := c.ctx.Emitter.Size()
	returnRuleCtx := c.compileInitialization(ctx, core.ForInLoop)
	if returnRuleCtx == nil {
		return ProtectedRecoveryRegion{Result: bytecode.NoopOperand}
	}

	loop := c.ctx.Loops.Current()
	breakLabel := loop.BreakLabel()

	c.compileLoopBody(ctx)

	out := c.compileFinalization(returnRuleCtx)
	endCatchExclusive := c.ctx.Emitter.Size()

	routeRecovery := c.ctx.Emitter.NewLabel("recovery", "for", "route")
	c.ctx.Emitter.EmitJumpIfTrue(errorStateReg, routeRecovery)
	c.ctx.Emitter.EmitJump(endLabel)

	errorPreludePC := c.ctx.Emitter.Size()
	c.ctx.Emitter.EmitBoolean(errorStateReg, true)
	c.ctx.Emitter.EmitJump(breakLabel)

	c.ctx.Emitter.MarkLabel(routeRecovery)
	c.ctx.Emitter.EmitBoolean(errorStateReg, false)
	c.ctx.Emitter.EmitJump(recoveryLabel)

	return ProtectedRecoveryRegion{
		Result:            out,
		StartCatch:        startCatch,
		EndCatchExclusive: endCatchExclusive,
		CatchHandlerPC:    errorPreludePC,
		HasTimeout:        false,
	}
}

// compileInitialization handles the setup of a loop, including determining its type,
// compiling its source, declaring variables, and emitting initialization instructions.
// Parameters:
//   - ctx: The FOR expression context from the AST
//   - kind: The kind of loop (ForInLoop, WhileLoop, or DoWhileLoop)
//
// Returns the rule context for the return expression or nested FOR expression.
func (c *LoopCompiler) compileInitialization(ctx fql.IForExpressionContext, kind core.LoopKind) antlr.RuleContext {
	returnRuleCtx, distinct, loopType, ok := c.resolveLoopReturnSpec(ctx.ForExpressionReturn())
	if !ok {
		return nil
	}

	// Create a new loop with the determined properties
	loop := c.ctx.Loops.NewLoop(kind, loopType, distinct)
	c.setLoopDestinationType(loop)

	c.configureLoopRuntime(loop, ctx, kind)

	// Push the loop onto the stack and enter a new symbol scope
	c.ctx.Loops.Push(loop)
	c.ctx.Symbols.EnterScope()

	valueType, keyType := c.inferLoopVariableTypes(ctx, loop, kind)
	c.declareLoopVariables(ctx, loop, valueType, keyType)
	c.emitLoopInitialization(ctx, loop)
	c.patchDistinctLoopDestination(loop)

	return returnRuleCtx
}

func (c *LoopCompiler) resolveLoopReturnSpec(returnCtx fql.IForExpressionReturnContext) (antlr.RuleContext, bool, core.LoopType, bool) {
	if returnCtx == nil {
		return nil, false, core.NormalLoop, false
	}

	if re := returnCtx.ReturnExpression(); re != nil {
		return re, re.Distinct() != nil, core.NormalLoop, true
	}

	if fe := returnCtx.ForExpression(); fe != nil {
		return fe, false, core.PassThroughLoop, true
	}

	return nil, false, core.NormalLoop, false
}

func (c *LoopCompiler) setLoopDestinationType(loop *core.Loop) {
	if loop != nil && loop.Dst.IsRegister() {
		c.ctx.Types.Set(loop.Dst, core.TypeList)
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
	val := ctx.GetValueVariable()
	if val == nil {
		return
	}

	varName := textOfLoopVariable(val)
	loop.DeclareValueVar(varName, c.ctx.Symbols, valueType)

	if loop.Value.IsRegister() {
		c.ctx.Types.Set(loop.Value, valueType)
	}
}

func (c *LoopCompiler) declareLoopCounterVariable(ctx fql.IForExpressionContext, loop *core.Loop, keyType core.ValueType) {
	ctr := ctx.GetCounterVariable()
	if ctr == nil {
		return
	}

	loop.DeclareKeyVar(textOfBindingIdentifier(ctr), c.ctx.Symbols, keyType)
	if loop.Key.IsRegister() {
		c.ctx.Types.Set(loop.Key, keyType)
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

	c.ctx.Emitter.WithSpan(span, func() {
		loop.EmitInitialization(c.ctx.Registers, c.ctx.Emitter)
	})
}

func (c *LoopCompiler) patchDistinctLoopDestination(loop *core.Loop) {
	if loop.Allocate || !loop.Distinct {
		return
	}

	parent := c.ctx.Loops.RequiredParent(c.ctx.Loops.Depth())
	c.ctx.Emitter.Patchx(parent.StartLabel(), 1)
}

// compileFinalization handles the teardown of a loop, including processing the return expression,
// emitting finalization instructions, and cleaning up the symbol scope.
// Parameters:
//   - ctx: The rule context for the return expression or nested FOR expression
//
// Returns the destination operand containing the loop results.
func (c *LoopCompiler) compileFinalization(ctx antlr.RuleContext) bytecode.Operand {
	loop := c.ctx.Loops.Current()

	// Process the return expression based on the loop type
	if loop.Type != core.PassThroughLoop {
		// For normal loops, compile the return expression and push the result to the destination
		re := ctx.(*fql.ReturnExpressionContext)
		expReg := c.exprs.Compile(re.Expression())

		span := source.Span{Start: -1, End: -1}

		if exprCtx := re.Expression(); exprCtx != nil {
			if prc, ok := exprCtx.(antlr.ParserRuleContext); ok {
				span = parser.SpanFromRuleContext(prc)
			}
		} else {
			span = parser.SpanFromRuleContext(re)
		}

		c.ctx.Emitter.WithSpan(span, func() {
			c.ctx.Emitter.EmitAB(bytecode.OpPush, loop.Dst, expReg)
		})
	} else if ctx != nil {
		// For pass-through loops, recursively compile the nested FOR expression
		if fe, ok := ctx.(*fql.ForExpressionContext); ok {
			c.Compile(fe)
		}
	}

	// Emit VM instructions for loop finalization
	loop.EmitFinalization(c.ctx.Emitter)

	// Clean up the symbol scope and pop the loop from the stack
	c.ctx.Symbols.ExitScope()
	c.ctx.Loops.Pop()

	return loop.Dst
}

// compileForExpressionSource processes the source expression for a FOR IN loop.
// It handles various types of expressions that can be used as the source collection,
// such as function calls, member expressions, variables, parameters, range operators, and literals.
// Returns an operand representing the compiled source expression.
func (c *LoopCompiler) compileForExpressionSource(ctx fql.IForExpressionSourceContext) bytecode.Operand {
	return c.compileLoopOperand(
		loopOperandContext{
			param:                  ctx.Param(),
			variable:               ctx.Variable(),
			memberExpression:       ctx.MemberExpression(),
			functionCallExpression: ctx.FunctionCallExpression(),
			rangeOperator:          ctx.RangeOperator(),
			arrayLiteral:           ctx.ArrayLiteral(),
			objectLiteral:          ctx.ObjectLiteral(),
		},
		loopOperandFunctionCallExpression,
		loopOperandMemberExpression,
		loopOperandVariable,
		loopOperandParam,
		loopOperandRangeOperator,
		loopOperandArrayLiteral,
		loopOperandObjectLiteral,
	)
}

func (c *LoopCompiler) compileLoopBody(ctx fql.IForExpressionContext) {
	if ctx == nil {
		return
	}

	if body := ctx.AllForExpressionBody(); len(body) > 0 {
		for _, b := range body {
			if ec := b.ForExpressionStatement(); ec != nil {
				c.compileForExpressionStatement(ec)
			} else if ec := b.ForExpressionClause(); ec != nil {
				c.compileForExpressionClause(ec)
			}
		}
	}
}

// compileForExpressionStatement processes statements within a FOR loop body.
// These can be variable declarations or function calls.
// The results of these statements are not used directly in the loop result.
func (c *LoopCompiler) compileForExpressionStatement(ctx fql.IForExpressionStatementContext) {
	// Handle variable declarations (e.g., LET x = 1)
	if vd := ctx.VariableDeclaration(); vd != nil {
		_ = c.bindings.CompileVariableDeclaration(vd)
	} else if as := ctx.AssignmentStatement(); as != nil {
		_ = c.bindings.CompileAssignmentStatement(as)
	} else if fce := ctx.FunctionCallExpression(); fce != nil {
		// Handle function calls (e.g., doSomething())
		_ = c.exprs.CompileFunctionCallExpression(fce)
	}
}

// compileForExpressionClause processes clauses within a FOR loop body.
// These can be LIMIT, FILTER, SORT, or COLLECT clauses that modify the loop behavior.
// Each clause type is delegated to a specific compilation method.
func (c *LoopCompiler) compileForExpressionClause(ctx fql.IForExpressionClauseContext) {
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
	state := c.ctx.Registers.Allocate()
	c.ctx.Loops.Current().RegisterReset(state)
	// Emit the iterator limit instruction with the loop's end label
	c.ctx.Emitter.EmitIterLimit(state, src, c.ctx.Loops.Current().BreakLabel())
}

// compileOffset emits VM instructions to skip a number of iterations at the start of a loop.
// It allocates a state register and emits an iterator skip instruction with the loop's jump label.
func (c *LoopCompiler) compileOffset(src bytecode.Operand) {
	// Allocate a state register for the offset operation
	state := c.ctx.Registers.Allocate()
	c.ctx.Loops.Current().RegisterReset(state)
	// Emit the iterator skip instruction with the loop's jump label
	c.ctx.Emitter.EmitIterSkip(state, src, c.ctx.Loops.Current().ContinueLabel())
}

// compileFilterClause processes a FILTER clause in a FOR loop.
// It compiles the filter expression and emits a conditional jump instruction
// that skips the current iteration if the filter condition is false.
func (c *LoopCompiler) compileFilterClause(ctx fql.IFilterClauseContext) {
	// Compile the filter expression (e.g., FILTER x > 5)
	// Get the jump label for the current loop
	label := c.ctx.Loops.Current().ContinueLabel()
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

	if srcCtx.RangeOperator() != nil {
		return core.TypeInt, core.TypeInt
	}

	if al := srcCtx.ArrayLiteral(); al != nil {
		return c.inferArrayLiteralElementType(al), core.TypeInt
	}

	if srcCtx.ObjectLiteral() != nil {
		return core.TypeAny, core.TypeString
	}

	if v := srcCtx.Variable(); v != nil {
		if binding, ok := c.ctx.Symbols.ResolveBinding(v.GetText()); ok {
			return c.inferValueKeyFromCollection(binding.Type)
		}
	}

	if srcCtx.Param() != nil || srcCtx.FunctionCallExpression() != nil {
		return core.TypeAny, core.TypeAny
	}

	if srcCtx.MemberExpression() != nil {
		return c.inferValueKeyFromCollection(c.facts.OperandType(src))
	}

	return c.inferValueKeyFromCollection(c.facts.OperandType(src))
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

	args := ctx.ArgumentList()
	if args == nil {
		return core.TypeUnknown
	}

	exps := args.AllExpression()
	if len(exps) == 0 {
		return core.TypeUnknown
	}

	elemType := core.TypeUnknown

	for _, exp := range exps {
		typ := c.inferExpressionType(exp)
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
		if binding, ok := c.ctx.Symbols.ResolveBinding(v.GetText()); ok {
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
