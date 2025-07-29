package internal

import (
	"github.com/MontFerret/ferret/pkg/compiler/internal/core"
	"github.com/MontFerret/ferret/pkg/parser/fql"
	"github.com/MontFerret/ferret/pkg/vm"
)

// StmtCompiler handles the compilation of FQL statements.
// It transforms statement operations from the AST into VM instructions.
type StmtCompiler struct {
	ctx *CompilerContext
}

// NewStmtCompiler creates a new instance of StmtCompiler with the given compiler context.
func NewStmtCompiler(ctx *CompilerContext) *StmtCompiler {
	return &StmtCompiler{
		ctx: ctx,
	}
}

// Compile processes the main body of an FQL query.
// It first compiles all body statements (variable declarations, function calls, etc.),
// and then compiles the body expression (FOR loops or RETURN statements).
// Parameters:
//   - ctx: The body context from the AST
func (c *StmtCompiler) Compile(ctx fql.IBodyContext) {
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
func (c *StmtCompiler) CompileBodyStatement(ctx fql.IBodyStatementContext) {
	// Handle variable declarations (e.g., LET x = 1)
	if vd := ctx.VariableDeclaration(); vd != nil {
		c.CompileVariableDeclaration(vd)
	} else if fce := ctx.FunctionCallExpression(); fce != nil {
		// Handle function calls (e.g., WAIT(1000))
		c.CompileFunctionCall(fce)
	} else if wfe := ctx.WaitForExpression(); wfe != nil {
		// Handle wait expressions (e.g., WAIT FOR x RETURN y)
		c.ctx.WaitCompiler.Compile(wfe)
	}
}

// CompileBodyExpression processes the main expression in the body of an FQL query.
// This is typically either a FOR loop or a RETURN statement, which forms the primary
// operation of the query and determines what data is returned.
// Parameters:
//   - ctx: The body expression context from the AST
func (c *StmtCompiler) CompileBodyExpression(ctx fql.IBodyExpressionContext) {
	// Handle FOR expressions (e.g., FOR x IN y RETURN z)
	if fe := ctx.ForExpression(); fe != nil {
		// Compile the FOR loop and get the destination register
		out := c.ctx.LoopCompiler.Compile(fe)

		// Emit a return instruction with the loop result
		c.ctx.Emitter.EmitA(vm.OpReturn, out)
	} else if re := ctx.ReturnExpression(); re != nil {
		// Handle RETURN expressions (e.g., RETURN x)
		// Compile the expression to return
		valReg := c.ctx.ExprCompiler.Compile(re.Expression())

		// If the result is a constant, we need to move it to a temporary register
		// because constants cannot be directly returned
		if valReg.IsConstant() {
			valC := valReg
			valReg = c.ctx.Registers.Allocate(core.Temp)

			// Move the constant value to the temporary register
			c.ctx.Emitter.EmitMove(valReg, valC)
		}

		// Emit a return instruction with the expression result
		c.ctx.Emitter.EmitA(vm.OpReturn, valReg)
	}
}

// CompileVariableDeclaration processes a variable declaration statement in an FQL query.
// It handles both regular identifiers and safe reserved words as variable names,
// and manages the assignment of values to either global or local variables based on scope.
// Parameters:
//   - ctx: The variable declaration context from the AST
//
// Returns:
//   - An operand representing the register where the variable value is stored,
//     or NoopOperand if the variable is ignored
func (c *StmtCompiler) CompileVariableDeclaration(ctx fql.IVariableDeclarationContext) vm.Operand {
	// Start with the ignore pseudo-variable as the default name
	name := core.IgnorePseudoVariable

	// Extract the variable name from either an identifier or a safe reserved word
	if id := ctx.Identifier(); id != nil {
		name = id.GetText()
	} else if reserved := ctx.SafeReservedWord(); reserved != nil {
		name = reserved.GetText()
	}

	// Compile the expression that provides the variable's value
	src := c.ctx.ExprCompiler.Compile(ctx.Expression())

	// If this is a real variable (not the ignore pseudo-variable)
	if name != core.IgnorePseudoVariable {
		// Handle constant values differently - they need to be loaded into a register
		if src.IsConstant() {
			// Declare a global variable and load the constant into it
			dest, ok := c.ctx.Symbols.DeclareGlobal(name, core.TypeUnknown)

			if !ok {
				c.ctx.Errors.VariableNotUnique(ctx, name)

				return vm.NoopOperand
			}

			c.ctx.Emitter.EmitAB(vm.OpLoadConst, dest, src)
			c.ctx.Registers.Free(src)

			src = dest
		} else if c.ctx.Symbols.Scope() == 0 {
			// If we're in the global scope, assign as a global variable
			if ok := c.ctx.Symbols.AssignGlobal(name, core.TypeUnknown, src); !ok {
				c.ctx.Errors.VariableNotUnique(ctx, name)

				return vm.NoopOperand
			}
		} else {
			// Otherwise, assign as a local variable in the current scope
			if ok := c.ctx.Symbols.AssignLocal(name, core.TypeUnknown, src); !ok {
				c.ctx.Errors.VariableNotUnique(ctx, name)

				return vm.NoopOperand
			}
		}

		// Return the register containing the variable's value
		return src
	}

	// For ignored variables, return a no-op operand
	return vm.NoopOperand
}

// CompileFunctionCall processes a function call expression in an FQL query.
// It delegates the compilation to the ExprCompiler, which handles the details
// of compiling function calls with their arguments and return values.
// Parameters:
//   - ctx: The function call expression context from the AST
//
// Returns:
//   - An operand representing the register where the function call result is stored
func (c *StmtCompiler) CompileFunctionCall(ctx fql.IFunctionCallExpressionContext) vm.Operand {
	// Delegate to the expression compiler for function call compilation
	return c.ctx.ExprCompiler.CompileFunctionCallExpression(ctx)
}
