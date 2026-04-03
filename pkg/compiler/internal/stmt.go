package internal

import (
	"github.com/MontFerret/ferret/v2/pkg/bytecode"
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

// Compile processes the main body of an FQL query.
// It first compiles all body statements (variable declarations, function calls, etc.),
// and then compiles the body expression (FOR loops or RETURN statements).
// Parameters:
//   - ctx: The body context from the AST
func (c *StatementCompiler) Compile(ctx fql.IBodyContext) {
	if ctx == nil {
		return
	}

	// Process all statements in the body
	for _, statement := range ctx.AllBodyStatement() {
		c.CompileBodyStatement(statement)
	}

	// Process the final expression in the body
	c.CompileBodyExpression(ctx.BodyExpression())
}

// CompileBodyStatement processes a single statement in the body of an FQL query.
// It determines the type of statement (variable declaration, function call, or wait expression)
// and delegates to the appropriate compilation method.
// Parameters:
//   - ctx: The body statement context from the AST
func (c *StatementCompiler) CompileBodyStatement(ctx fql.IBodyStatementContext) {
	if ctx == nil {
		return
	}

	if vd := ctx.VariableDeclaration(); vd != nil {
		c.bindings.CompileVariableDeclaration(vd)
	} else if as := ctx.AssignmentStatement(); as != nil {
		c.bindings.CompileAssignmentStatement(as)
	} else if fd := ctx.FunctionDeclaration(); fd != nil {
		// Function declarations are compiled separately.
		return
	} else if fce := ctx.FunctionCallExpression(); fce != nil {
		// Handle function calls (e.g., WAIT(1000))
		c.CompileFunctionCall(fce)
	} else if wfe := ctx.WaitForExpression(); wfe != nil {
		// Handle wait expressions (e.g., WAIT FOR x RETURN y)
		c.wait.Compile(wfe)
	} else if de := ctx.DispatchExpression(); de != nil {
		c.dispatch.Compile(de)
	}
}

// CompileBodyExpression processes the main expression in the body of an FQL query.
// This is typically either a FOR loop or a RETURN statement, which forms the primary
// operation of the query and determines what data is returned.
// Parameters:
//   - ctx: The body expression context from the AST
func (c *StatementCompiler) CompileBodyExpression(ctx fql.IBodyExpressionContext) {
	if ctx == nil {
		return
	}

	// Handle FOR expressions (e.g., FOR x IN y RETURN z)
	if fe := ctx.ForExpression(); fe != nil {
		// Compile the FOR loop and get the destination register
		out := c.loops.Compile(fe)

		// Emit a return instruction with the loop result
		c.ctx.Emitter.EmitA(bytecode.OpReturn, out)
	} else if re := ctx.ReturnExpression(); re != nil {
		// Handle RETURN expressions (e.g., RETURN x)
		// Compile and normalize into a register because RETURN expects a register operand.
		valReg := ensureOperandRegister(c.ctx, c.facts, c.exprs.Compile(re.Expression()))

		// Emit a return instruction with the expression result
		c.ctx.Emitter.EmitA(bytecode.OpReturn, valReg)
	}
}

// CompileFunctionStatement processes a statement inside a UDF body.
// It supports variable declarations, nested function declarations, expression statements,
// function calls, and other statements allowed in the main body.
func (c *StatementCompiler) CompileFunctionStatement(ctx fql.IFunctionStatementContext) {
	if ctx == nil {
		return
	}

	stmt, ok := ctx.(*fql.FunctionStatementContext)
	if !ok || stmt == nil {
		return
	}

	switch {
	case stmt.VariableDeclaration() != nil:
		c.bindings.CompileVariableDeclaration(stmt.VariableDeclaration())
	case stmt.AssignmentStatement() != nil:
		c.bindings.CompileAssignmentStatement(stmt.AssignmentStatement())
	case stmt.FunctionDeclaration() != nil:
		// Nested function declarations are compiled separately.
		return
	case stmt.FunctionCallExpression() != nil:
		c.CompileFunctionCall(stmt.FunctionCallExpression())
	case stmt.WaitForExpression() != nil:
		c.wait.Compile(stmt.WaitForExpression())
	case stmt.DispatchExpression() != nil:
		c.dispatch.Compile(stmt.DispatchExpression())
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
		c.exprs.Compile(expr)
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
