package internal

import (
	"github.com/antlr4-go/antlr/v4"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/compiler/internal/core"
	parserd "github.com/MontFerret/ferret/v2/pkg/parser/diagnostics"
	"github.com/MontFerret/ferret/v2/pkg/parser/fql"
)

// StatementCompiler handles the compilation of FQL statements.
// It transforms statement operations from the AST into VM instructions.
type StatementCompiler struct {
	ctx      *CompilationSession
	bindings *BindingCompiler
	dispatch *DispatchCompiler
	exprs    *ExprCompiler
	loops    *LoopCompiler
	facts    *TypeFacts
	wait     *WaitCompiler
}

// NewStatementCompiler creates a new instance of StatementCompiler with the given compiler context.
func NewStatementCompiler(ctx *CompilationSession) *StatementCompiler {
	return &StatementCompiler{
		ctx: ctx,
	}
}

func (c *StatementCompiler) bind(bindings *BindingCompiler, dispatch *DispatchCompiler, exprs *ExprCompiler, loops *LoopCompiler, facts *TypeFacts, wait *WaitCompiler) {
	if c == nil {
		return
	}

	c.bindings = bindings
	c.dispatch = dispatch
	c.exprs = exprs
	c.loops = loops
	c.facts = facts
	c.wait = wait
}

// Compile emits an implicit NONE return when the script has no explicit RETURN.
func (c *StatementCompiler) Compile(ctx fql.IBodyContext) {
	if ctx == nil {
		return
	}

	// Process all statements in the body
	for _, statement := range ctx.AllBodyStatement() {
		c.CompileBodyStatement(statement)
	}

	if expr := ctx.BodyExpression(); expr != nil {
		c.CompileBodyExpression(expr)
	} else {
		c.emitImplicitNoneReturn()
	}
}

// CompileBodyStatement compiles a value-discarding top-level statement.
func (c *StatementCompiler) CompileBodyStatement(ctx fql.IBodyStatementContext) {
	if ctx == nil {
		return
	}

	rule, _ := ctx.(antlr.ParserRuleContext)
	c.ctx.WithDebugPoint(rule, func() {
		c.compileBodyStatement(ctx)
	})
}

func (c *StatementCompiler) compileBodyStatement(ctx fql.IBodyStatementContext) {
	if vd := ctx.VariableDeclaration(); vd != nil {
		c.bindings.CompileVariableDeclaration(vd)
	} else if as := ctx.AssignmentStatement(); as != nil {
		c.bindings.CompileAssignmentStatement(as)
	} else if ds := ctx.DeleteStatement(); ds != nil {
		c.bindings.CompileDeleteStatement(ds)
	} else if fd := ctx.FunctionDeclaration(); fd != nil {
		// Function declarations are compiled separately.
		return
	} else if fce := ctx.FunctionCallExpression(); fce != nil {
		// Handle function calls (e.g., WAIT(1s))
		c.CompileFunctionCall(fce)
	} else if wfe := ctx.WaitForExpression(); wfe != nil {
		// Handle wait expressions (e.g., WAIT FOR x RETURN y)
		c.wait.Compile(wfe)
	} else if de := ctx.DispatchExpression(); de != nil {
		c.dispatch.Compile(de)
	} else if fe := ctx.ForExpression(); fe != nil {
		c.loops.CompileDiscarded(fe)
	}
}

// CompileBodyExpression compiles the script's explicit RETURN.
func (c *StatementCompiler) CompileBodyExpression(ctx fql.IBodyExpressionContext) {
	if ctx == nil {
		return
	}

	rule, _ := ctx.(antlr.ParserRuleContext)
	c.ctx.WithDebugPointKind(rule, bytecode.DebugPointReturn, func() {
		c.compileBodyExpression(ctx)
	})
}

func (c *StatementCompiler) compileBodyExpression(ctx fql.IBodyExpressionContext) {
	if re := ctx.ReturnExpression(); re != nil {
		valReg := c.compileReturnOperand(re.ReturnValue(), re.Distinct() != nil)

		c.ctx.Program.Emitter.EmitA(bytecode.OpReturn, valReg)
	}
}

func (c *StatementCompiler) emitImplicitNoneReturn() {
	c.ctx.Program.Emitter.EmitA(bytecode.OpReturn, bytecode.NoopOperand)
}

// compileReturnOperand compiles an explicit RETURN operand, including a
// directly returned FOR expression.
func (c *StatementCompiler) compileReturnOperand(ctx fql.IReturnValueContext, distinct bool) bytecode.Operand {
	if ctx == nil {
		return bytecode.NoopOperand
	}

	if expr := ctx.Expression(); expr != nil {
		return c.CompileReturnValue(expr, distinct)
	}

	loop := ctx.ForExpression()
	if loop == nil {
		return bytecode.NoopOperand
	}

	valReg := ensureOperandRegister(c.ctx, c.facts, c.loops.Compile(loop))

	return c.compileReturnDistinct(valReg, loop.(antlr.ParserRuleContext), distinct)
}

// CompileReturnValue compiles a non-loop return expression and applies its
// statement-level result modifiers.
func (c *StatementCompiler) CompileReturnValue(expr fql.IExpressionContext, distinct bool) bytecode.Operand {
	if expr == nil {
		return bytecode.NoopOperand
	}

	valReg := ensureOperandRegister(c.ctx, c.facts, c.exprs.Compile(expr))

	return c.compileReturnDistinct(valReg, expr.(antlr.ParserRuleContext), distinct)
}

func (c *StatementCompiler) compileReturnDistinct(valReg bytecode.Operand, ctx antlr.ParserRuleContext, distinct bool) bytecode.Operand {
	if !distinct {
		return valReg
	}

	switch c.facts.OperandType(valReg) {
	case core.TypeArray, core.TypeList, core.TypeUnknown, core.TypeAny:
	default:
		err := c.ctx.Program.Errors.Create(parserd.SemanticError, ctx, "RETURN DISTINCT requires an array expression")
		c.ctx.Program.Errors.Add(err)
		return valReg
	}

	dst := c.ctx.Function.Registers.Allocate()
	c.ctx.Program.Emitter.WithSpan(parserd.SpanFromRuleContext(ctx), func() {
		c.ctx.Program.Emitter.EmitAB(bytecode.OpDistinct, dst, valReg)
	})
	c.ctx.Function.Types.Set(dst, core.TypeArray)

	return dst
}

// CompileFunctionStatement compiles a value-discarding UDF block statement.
func (c *StatementCompiler) CompileFunctionStatement(ctx fql.IFunctionStatementContext) {
	if ctx == nil {
		return
	}

	stmt, ok := ctx.(*fql.FunctionStatementContext)
	if !ok || stmt == nil {
		return
	}

	c.ctx.WithDebugPoint(stmt, func() {
		c.compileFunctionStatement(stmt)
	})
}

func (c *StatementCompiler) compileFunctionStatement(stmt *fql.FunctionStatementContext) {
	switch {
	case stmt.VariableDeclaration() != nil:
		c.bindings.CompileVariableDeclaration(stmt.VariableDeclaration())
	case stmt.AssignmentStatement() != nil:
		c.bindings.CompileAssignmentStatement(stmt.AssignmentStatement())
	case stmt.DeleteStatement() != nil:
		c.bindings.CompileDeleteStatement(stmt.DeleteStatement())
	case stmt.FunctionDeclaration() != nil:
		// Nested function declarations are compiled separately.
		return
	case stmt.FunctionCallExpression() != nil:
		c.CompileFunctionCall(stmt.FunctionCallExpression())
	case stmt.WaitForExpression() != nil:
		c.wait.Compile(stmt.WaitForExpression())
	case stmt.DispatchExpression() != nil:
		c.dispatch.Compile(stmt.DispatchExpression())
	case stmt.ForExpression() != nil:
		c.loops.CompileDiscarded(stmt.ForExpression())
	case stmt.ExpressionStatement() != nil:
		c.CompileExpressionStatement(stmt.ExpressionStatement())
	}
}

// CompileExpressionStatement evaluates an expression for its side effects and discards the result.
func (c *StatementCompiler) CompileExpressionStatement(ctx fql.IExpressionStatementContext) {
	if ctx == nil {
		return
	}

	stmt, ok := ctx.(*fql.ExpressionStatementContext)
	if !ok || stmt == nil {
		return
	}

	if expr := stmt.Expression(); expr != nil {
		c.exprs.CompileDiscarded(expr)
	}
}

// CompileFunctionCall processes a function call expression in an FQL query.
// It delegates the compilation to the ExprCompiler, which handles the details
// of compiling function calls with their arguments and return values.
// Parameters:
//   - ctx: The function call expression context from the AST
//
// Returns:
//   - An operand representing the register where the function call result is stored
func (c *StatementCompiler) CompileFunctionCall(ctx fql.IFunctionCallExpressionContext) bytecode.Operand {
	if ctx == nil {
		return bytecode.NoopOperand
	}

	// Delegate to the expression compiler for function call compilation
	return c.exprs.CompileFunctionCallExpression(ctx)
}
