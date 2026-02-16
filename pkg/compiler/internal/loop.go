package internal

import (
	"fmt"

	"github.com/antlr4-go/antlr/v4"

	parser "github.com/MontFerret/ferret/v2/pkg/parser/diagnostics"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"

	"github.com/MontFerret/ferret/v2/pkg/compiler/internal/core"
	"github.com/MontFerret/ferret/v2/pkg/file"
	"github.com/MontFerret/ferret/v2/pkg/parser/fql"
)

// LoopCompiler handles the compilation of FOR loop expressions in FQL queries.
// It transforms loop operations into VM instructions for iteration, filtering, and data manipulation.
type LoopCompiler struct {
	ctx *CompilerContext
}

// NewLoopCompiler creates a new instance of LoopCompiler with the given compiler context.
func NewLoopCompiler(ctx *CompilerContext) *LoopCompiler {
	return &LoopCompiler{ctx: ctx}
}

// Compile processes a FOR expression from the FQL AST and generates the appropriate VM instructions.
// It determines whether to compile a FOR IN loop (iteration over a collection), a FOR WHILE loop (while condition),
// or a FOR STEP loop (with init, condition, and increment expressions).
// Returns an operand representing the destination of the loop results.
func (c *LoopCompiler) Compile(ctx fql.IForExpressionContext) bytecode.Operand {
	var returnRuleCtx antlr.RuleContext

	if ctx.In() != nil {
		returnRuleCtx = c.compileInitialization(ctx, core.ForInLoop)
	} else if ctx.Step() != nil {
		returnRuleCtx = c.compileInitialization(ctx, core.ForStepLoop)
	} else if ctx.Do() == nil {
		returnRuleCtx = c.compileInitialization(ctx, core.WhileLoop)
	} else {
		returnRuleCtx = c.compileInitialization(ctx, core.DoWhileLoop)
	}

	// Probably, a syntax error happened and no return rule context was created.
	if returnRuleCtx == nil {
		return bytecode.NoopOperand
	}

	// Compile the loop body (statements and clauses)
	if body := ctx.AllForExpressionBody(); len(body) > 0 {
		for _, b := range body {
			if ec := b.ForExpressionStatement(); ec != nil {
				c.compileForExpressionStatement(ec)
			} else if ec := b.ForExpressionClause(); ec != nil {
				c.compileForExpressionClause(ec)
			}
		}
	}

	// Finalize the loop and return the destination operand
	return c.compileFinalization(returnRuleCtx)
}

// compileInitialization handles the setup of a loop, including determining its type,
// compiling its source, declaring variables, and emitting initialization instructions.
// Parameters:
//   - ctx: The FOR expression context from the AST
//   - kind: The kind of loop (ForInLoop, WhileLoop, or ForStepLoop)
//
// Returns the rule context for the return expression or nested FOR expression.
func (c *LoopCompiler) compileInitialization(ctx fql.IForExpressionContext, kind core.LoopKind) antlr.RuleContext {
	var distinct bool
	var returnRuleCtx antlr.RuleContext
	var loopType core.LoopType
	returnCtx := ctx.ForExpressionReturn()

	if returnCtx == nil {
		return nil
	}

	// Determine the loop type and whether it should use distinct values
	if re := returnCtx.ReturnExpression(); re != nil {
		returnRuleCtx = re
		distinct = re.Distinct() != nil
		loopType = core.NormalLoop
	} else if fe := returnCtx.ForExpression(); fe != nil {
		returnRuleCtx = fe
		loopType = core.PassThroughLoop
	}

	// Create a new loop with the determined properties
	loop := c.ctx.Loops.NewLoop(kind, loopType, distinct)
	if loop.Dst.IsRegister() {
		c.ctx.Types.Set(loop.Dst, core.TypeList)
	}

	// Set up the loop source based on the loop kind
	switch kind {
	case core.ForInLoop:
		// For IN loops, compile the collection to iterate over
		loop.Src = c.compileForExpressionSource(ctx.ForExpressionSource())
	case core.ForStepLoop:
		// For STEP loops, set up functions to evaluate init, condition, and increment expressions
		loop.InitFn = func() bytecode.Operand {
			return c.ctx.ExprCompiler.Compile(ctx.GetStepInit())
		}
		loop.ConditionFn = func() bytecode.Operand {
			return c.ctx.ExprCompiler.Compile(ctx.GetStepCondition())
		}
		loop.UpdateFn = func() bytecode.Operand {
			if exp := ctx.GetStepUpdateExp(); exp != nil {
				// If an increment expression is provided, use it
				return c.ctx.ExprCompiler.Compile(exp)
			}

			variable, _, found := c.ctx.Symbols.Resolve(ctx.GetStepVariable().GetText())

			if !found {
				c.ctx.Errors.VariableNotFound(ctx.GetStepVariable(), ctx.GetStepVariable().GetText())
				return bytecode.NoopOperand
			}

			inc := ctx.GetStepUpdate()

			return c.ctx.ExprCompiler.CompileIncDec(inc, variable)
		}
	default:
		// For WHILE/DO-WHILE loops, set up a function to evaluate the condition
		loop.ConditionFn = func() bytecode.Operand {
			return c.ctx.ExprCompiler.Compile(ctx.Expression(0))
		}
	}

	// Push the loop onto the stack and enter a new symbol scope
	c.ctx.Loops.Push(loop)
	c.ctx.Symbols.EnterScope()

	// Declare variables for the loop value and counter if specified
	var valueType core.ValueType
	var keyType core.ValueType

	switch kind {
	case core.ForInLoop:
		valueType, keyType = c.inferForInTypes(ctx.ForExpressionSource(), loop.Src)
	case core.ForStepLoop:
		valueType = c.inferExpressionType(ctx.GetStepInit())
	case core.WhileLoop, core.DoWhileLoop:
		valueType = core.TypeInt
	}

	if val := ctx.GetValueVariable(); val != nil {
		varName := val.GetText()
		loop.DeclareValueVar(varName, c.ctx.Symbols, valueType)

		if loop.Value.IsRegister() {
			c.ctx.Types.Set(loop.Value, valueType)
		}

		if loop.Kind == core.ForStepLoop {
			stepVar := ctx.GetStepVariable()

			if stepVar != nil && varName != stepVar.GetText() {
				if _, _, found := c.ctx.Symbols.Resolve(stepVar.GetText()); found {
					ce := c.ctx.Errors.Create(parser.SemanticError, ctx, fmt.Sprintf("step variable missmatch: expected '%s' but got '%s'", varName, stepVar.GetText()))
					ce.Hint = "Make sure the same variable is used in all parts of the STEP loop"
					c.ctx.Errors.Add(ce)
				}
			}
		}
	}

	if ctr := ctx.GetCounterVariable(); ctr != nil {
		loop.DeclareKeyVar(ctr.GetText(), c.ctx.Symbols, keyType)

		if loop.Key.IsRegister() {
			c.ctx.Types.Set(loop.Key, keyType)
		}
	}

	// Emit VM instructions for loop initialization
	span := file.Span{Start: -1, End: -1}

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

	// Handle distinct values if needed
	if !loop.Allocate {
		// If the current loop must push distinct items, we must patch the dest dataset
		if loop.Distinct {
			parent := c.ctx.Loops.FindParent(c.ctx.Loops.Depth())

			if parent == nil {
				panic("parent loop not found in loop table")
			}

			c.ctx.Emitter.Patchx(parent.StartLabel(), 1)
		}
	}

	return returnRuleCtx
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
		expReg := c.ctx.ExprCompiler.Compile(re.Expression())

		span := file.Span{Start: -1, End: -1}

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
	// Handle function call expressions (e.g., FOR x IN getUsers())
	if fce := ctx.FunctionCallExpression(); fce != nil {
		return c.ctx.ExprCompiler.CompileFunctionCallExpression(fce)
	}

	// Handle member expressions (e.g., FOR x IN users.active)
	if me := ctx.MemberExpression(); me != nil {
		return c.ctx.ExprCompiler.CompileMemberExpression(me)
	}

	// Handle variables (e.g., FOR x IN users)
	if v := ctx.Variable(); v != nil {
		return c.ctx.ExprCompiler.CompileVariable(v)
	}

	// Handle parameters (e.g., FOR x IN @users)
	if p := ctx.Param(); p != nil {
		return c.ctx.ExprCompiler.CompileParam(p)
	}

	// Handle range operators (e.g., FOR x IN 1..10)
	if ro := ctx.RangeOperator(); ro != nil {
		return c.ctx.ExprCompiler.CompileRangeOperator(ro)
	}

	// Handle array literals (e.g., FOR x IN [1, 2, 3])
	if al := ctx.ArrayLiteral(); al != nil {
		return c.ctx.LiteralCompiler.CompileArrayLiteral(al)
	}

	// Handle object literals (e.g., FOR x IN {a: 1, b: 2})
	if ol := ctx.ObjectLiteral(); ol != nil {
		return c.ctx.LiteralCompiler.CompileObjectLiteral(ol)
	}

	return bytecode.NoopOperand
}

// compileForExpressionStatement processes statements within a FOR loop body.
// These can be variable declarations or function calls.
// The results of these statements are not used directly in the loop result.
func (c *LoopCompiler) compileForExpressionStatement(ctx fql.IForExpressionStatementContext) {
	// Handle variable declarations (e.g., LET x = 1)
	if vd := ctx.VariableDeclaration(); vd != nil {
		_ = c.ctx.StmtCompiler.CompileVariableDeclaration(vd)
	} else if fce := ctx.FunctionCallExpression(); fce != nil {
		// Handle function calls (e.g., doSomething())
		_ = c.ctx.ExprCompiler.CompileFunctionCallExpression(fce)
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
	// Handle parameters (e.g., LIMIT @limit)
	if pm := ctx.Param(); pm != nil {
		return c.ctx.ExprCompiler.CompileParam(pm)
	}

	// Handle integer literals (e.g., LIMIT 10)
	if il := ctx.IntegerLiteral(); il != nil {
		return c.ctx.LiteralCompiler.CompileIntegerLiteral(il)
	}

	// Handle variables (e.g., LIMIT limit)
	if vb := ctx.Variable(); vb != nil {
		return c.ctx.ExprCompiler.CompileVariable(vb)
	}

	// Handle member expressions (e.g., LIMIT config.limit)
	if me := ctx.MemberExpression(); me != nil {
		return c.ctx.ExprCompiler.CompileMemberExpression(me)
	}

	// Handle function calls (e.g., LIMIT getLimit())
	if fce := ctx.FunctionCallExpression(); fce != nil {
		return c.ctx.ExprCompiler.CompileFunctionCallExpression(fce)
	}

	return bytecode.NoopOperand
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
	src := c.ctx.ExprCompiler.Compile(ctx.Expression())
	// Get the jump label for the current loop
	label := c.ctx.Loops.Current().ContinueLabel()
	// Emit a jump instruction that skips to the next iteration if the filter condition is false
	c.ctx.Emitter.EmitJumpIfFalse(src, label)
}

// compileSortClause processes a SORT clause in a FOR loop.
// It delegates the compilation to the specialized LoopSortCompiler.
func (c *LoopCompiler) compileSortClause(ctx fql.ISortClauseContext) {
	// Delegate to the specialized sort compiler
	c.ctx.LoopSortCompiler.Compile(ctx)
}

// compileCollectClause processes a COLLECT clause in a FOR loop.
// It delegates the compilation to the specialized LoopCollectCompiler.
func (c *LoopCompiler) compileCollectClause(ctx fql.ICollectClauseContext) {
	// Delegate to the specialized collect compiler
	c.ctx.LoopCollectCompiler.Compile(ctx)
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
		if op, _, ok := c.ctx.Symbols.Resolve(v.GetText()); ok {
			return c.inferValueKeyFromCollection(operandType(c.ctx, op))
		}
	}

	if srcCtx.Param() != nil || srcCtx.FunctionCallExpression() != nil {
		return core.TypeAny, core.TypeAny
	}

	if srcCtx.MemberExpression() != nil {
		return c.inferValueKeyFromCollection(operandType(c.ctx, src))
	}

	return c.inferValueKeyFromCollection(operandType(c.ctx, src))
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
		return literalType(lit)
	}

	if v := ctx.Variable(); v != nil {
		if op, _, ok := c.ctx.Symbols.Resolve(v.GetText()); ok {
			return operandType(c.ctx, op)
		}
		return core.TypeUnknown
	}

	if ctx.Param() != nil || ctx.FunctionCallExpression() != nil {
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
