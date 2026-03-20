package internal

import (
	"fmt"

	"github.com/antlr4-go/antlr/v4"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/compiler/internal/core"
	parserd "github.com/MontFerret/ferret/v2/pkg/parser/diagnostics"
	"github.com/MontFerret/ferret/v2/pkg/parser/fql"
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
func (c *StmtCompiler) CompileBodyStatement(ctx fql.IBodyStatementContext) {
	if ctx == nil {
		return
	}

	// Handle variable declarations (e.g., LET x = 1)
	if vd := ctx.VariableDeclaration(); vd != nil {
		c.CompileVariableDeclaration(vd)
	} else if as := ctx.AssignmentStatement(); as != nil {
		c.CompileAssignmentStatement(as)
	} else if fd := ctx.FunctionDeclaration(); fd != nil {
		// Function declarations are compiled separately.
		return
	} else if fce := ctx.FunctionCallExpression(); fce != nil {
		// Handle function calls (e.g., WAIT(1000))
		c.CompileFunctionCall(fce)
	} else if wfe := ctx.WaitForExpression(); wfe != nil {
		// Handle wait expressions (e.g., WAIT FOR x RETURN y)
		c.ctx.WaitCompiler.Compile(wfe)
	} else if de := ctx.DispatchExpression(); de != nil {
		c.ctx.DispatchCompiler.Compile(de)
	}
}

// CompileBodyExpression processes the main expression in the body of an FQL query.
// This is typically either a FOR loop or a RETURN statement, which forms the primary
// operation of the query and determines what data is returned.
// Parameters:
//   - ctx: The body expression context from the AST
func (c *StmtCompiler) CompileBodyExpression(ctx fql.IBodyExpressionContext) {
	if ctx == nil {
		return
	}

	// Handle FOR expressions (e.g., FOR x IN y RETURN z)
	if fe := ctx.ForExpression(); fe != nil {
		// Compile the FOR loop and get the destination register
		out := c.ctx.LoopCompiler.Compile(fe)

		// Emit a return instruction with the loop result
		c.ctx.Emitter.EmitA(bytecode.OpReturn, out)
	} else if re := ctx.ReturnExpression(); re != nil {
		// Handle RETURN expressions (e.g., RETURN x)
		// Compile and normalize into a register because RETURN expects a register operand.
		valReg := c.ctx.ExprCompiler.ensureRegister(c.ctx.ExprCompiler.Compile(re.Expression()))

		// Emit a return instruction with the expression result
		c.ctx.Emitter.EmitA(bytecode.OpReturn, valReg)
	}
}

// CompileFunctionStatement processes a statement inside a UDF body.
// It supports variable declarations, nested function declarations, expression statements,
// function calls, and other statements allowed in the main body.
func (c *StmtCompiler) CompileFunctionStatement(ctx fql.IFunctionStatementContext) {
	if ctx == nil {
		return
	}

	stmt, ok := ctx.(*fql.FunctionStatementContext)
	if !ok || stmt == nil {
		return
	}

	switch {
	case stmt.VariableDeclaration() != nil:
		c.CompileVariableDeclaration(stmt.VariableDeclaration())
	case stmt.AssignmentStatement() != nil:
		c.CompileAssignmentStatement(stmt.AssignmentStatement())
	case stmt.FunctionDeclaration() != nil:
		// Nested function declarations are compiled separately.
		return
	case stmt.FunctionCallExpression() != nil:
		c.CompileFunctionCall(stmt.FunctionCallExpression())
	case stmt.WaitForExpression() != nil:
		c.ctx.WaitCompiler.Compile(stmt.WaitForExpression())
	case stmt.DispatchExpression() != nil:
		c.ctx.DispatchCompiler.Compile(stmt.DispatchExpression())
	case stmt.ExpressionStatement() != nil:
		c.CompileExpressionStatement(stmt.ExpressionStatement())
	}
}

// CompileExpressionStatement evaluates an expression for its side effects and discards the result.
func (c *StmtCompiler) CompileExpressionStatement(ctx fql.IExpressionStatementContext) {
	if ctx == nil {
		return
	}

	stmt, ok := ctx.(*fql.ExpressionStatementContext)
	if !ok || stmt == nil {
		return
	}

	if expr := stmt.Expression(); expr != nil {
		c.ctx.ExprCompiler.Compile(expr)
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
func (c *StmtCompiler) CompileVariableDeclaration(ctx fql.IVariableDeclarationContext) bytecode.Operand {
	if ctx == nil {
		return bytecode.NoopOperand
	}

	name := variableDeclarationName(ctx)
	if name == "" {
		name = core.IgnorePseudoVariable
	}

	mutable := isMutableDeclaration(ctx)
	storage := declarationStorage(c.ctx, ctx.(antlr.ParserRuleContext), mutable)

	// Compile the expression that provides the variable's value
	src := c.ctx.ExprCompiler.Compile(ctx.Expression())
	srcType := operandType(c.ctx, src)

	// If this is a real variable (not the ignore pseudo-variable)
	if name != core.IgnorePseudoVariable {
		opts := core.BindingOptions{
			Mutable: mutable,
			Storage: storage,
		}

		if storage == core.BindingStorageCell {
			src = c.ctx.ExprCompiler.ensureRegister(src)

			var (
				dest bytecode.Operand
				ok   bool
			)

			if c.ctx.Symbols.Scope() == 0 {
				dest, ok = c.ctx.Symbols.DeclareGlobalWithOptions(name, srcType, opts)
			} else {
				dest, ok = c.ctx.Symbols.DeclareLocalWithOptions(name, srcType, opts)
			}

			if !ok {
				c.ctx.Errors.VariableNotUnique(ctx.(antlr.ParserRuleContext), name)
				return bytecode.NoopOperand
			}

			c.ctx.Emitter.EmitMakeCell(dest, src)
			c.ctx.Types.Set(dest, core.TypeAny)

			return dest
		}

		if mutable && src.IsConstant() {
			var (
				dest bytecode.Operand
				ok   bool
			)

			if c.ctx.Symbols.Scope() == 0 {
				dest, ok = c.ctx.Symbols.DeclareGlobalWithOptions(name, srcType, opts)
			} else {
				dest, ok = c.ctx.Symbols.DeclareLocalWithOptions(name, srcType, opts)
			}

			if !ok {
				c.ctx.Errors.VariableNotUnique(ctx.(antlr.ParserRuleContext), name)
				return bytecode.NoopOperand
			}

			c.ctx.Emitter.EmitLoadConst(dest, src)
			c.ctx.Types.Set(dest, srcType)

			return dest
		}

		if src.IsConstant() {
			var (
				dest bytecode.Operand
				ok   bool
			)

			if c.ctx.Symbols.Scope() == 0 {
				dest, ok = c.ctx.Symbols.DeclareGlobalWithOptions(name, srcType, opts)
			} else {
				dest, ok = c.ctx.Symbols.DeclareLocalWithOptions(name, srcType, opts)
			}

			if !ok {
				c.ctx.Errors.VariableNotUnique(ctx.(antlr.ParserRuleContext), name)
				return bytecode.NoopOperand
			}

			c.ctx.Emitter.EmitLoadConst(dest, src)
			c.ctx.Types.Set(dest, srcType)

			src = dest
		} else if c.ctx.Symbols.Scope() == 0 {
			if ok := c.ctx.Symbols.AssignGlobalWithOptions(name, srcType, src, opts); !ok {
				c.ctx.Errors.VariableNotUnique(ctx.(antlr.ParserRuleContext), name)
				return bytecode.NoopOperand
			}
		} else {
			if ok := c.ctx.Symbols.AssignLocalWithOptions(name, srcType, src, opts); !ok {
				c.ctx.Errors.VariableNotUnique(ctx.(antlr.ParserRuleContext), name)
				return bytecode.NoopOperand
			}
		}

		c.ctx.Types.Set(src, srcType)
		// Return the register containing the variable's value
		return src
	}

	// For ignored variables, return a no-op operand
	return bytecode.NoopOperand
}

// CompileFunctionCall processes a function call expression in an FQL query.
// It delegates the compilation to the ExprCompiler, which handles the details
// of compiling function calls with their arguments and return values.
// Parameters:
//   - ctx: The function call expression context from the AST
//
// Returns:
//   - An operand representing the register where the function call result is stored
func (c *StmtCompiler) CompileFunctionCall(ctx fql.IFunctionCallExpressionContext) bytecode.Operand {
	if ctx == nil {
		return bytecode.NoopOperand
	}

	// Delegate to the expression compiler for function call compilation
	return c.ctx.ExprCompiler.CompileFunctionCallExpression(ctx)
}

func (c *StmtCompiler) CompileAssignmentStatement(ctx fql.IAssignmentStatementContext) bytecode.Operand {
	if ctx == nil {
		return bytecode.NoopOperand
	}

	stmt, ok := ctx.(*fql.AssignmentStatementContext)
	if !ok || stmt == nil {
		return bytecode.NoopOperand
	}

	target := stmt.AssignmentTarget()
	if target == nil {
		return bytecode.NoopOperand
	}

	if member := target.MemberExpression(); member != nil {
		c.reportInvalidAssignmentTarget(member.(antlr.ParserRuleContext))
		return bytecode.NoopOperand
	}

	name := textOfBindingIdentifier(target.BindingIdentifier())
	if name == "" || name == core.IgnorePseudoVariable {
		c.reportInvalidAssignmentTarget(stmt)
		return bytecode.NoopOperand
	}

	binding, found := c.ctx.Symbols.ResolveBinding(name)
	if !found {
		c.ctx.Errors.VariableNotFound(stmt.GetStart(), name)
		return bytecode.NoopOperand
	}

	if !binding.Mutable {
		err := c.ctx.Errors.Create(parserd.SemanticError, stmt, fmt.Sprintf("Variable '%s' cannot be reassigned", name))
		err.Hint = "Declare it with VAR if you need to update it."
		c.ctx.Errors.Add(err)
		return bytecode.NoopOperand
	}

	operator := assignmentOperatorText(stmt)
	src := bytecode.NoopOperand

	if operator == "=" {
		src = c.ctx.ExprCompiler.Compile(stmt.Expression())
	} else if operator == "+=" && binding.Type == core.TypeString {
		left := c.snapshotBindingValue(binding)
		parts := append([]concatOperandSegment{{operand: left}}, buildConcatOperandSegmentsFromExpression(c.ctx.ExprCompiler, stmt.Expression())...)
		src = emitConcatOperandSegments(c.ctx, parts)
	} else {
		op, ok := resolveArithmeticBinaryOperator(operator)
		if !ok {
			return bytecode.NoopOperand
		}

		left := c.snapshotBindingValue(binding)
		right := c.ctx.ExprCompiler.Compile(stmt.Expression())
		src = emitBinaryOperation(c.ctx, stmt, op, left, right)
	}

	srcType := operandType(c.ctx, src)
	publishedType := srcType

	if c.ctx.Loops.Depth() > 0 {
		publishedType = core.JoinValueTypes(binding.Type, srcType)
	}

	binding.Type = publishedType

	return c.storeBindingValue(binding, src, publishedType)
}

func assignmentOperatorText(ctx *fql.AssignmentStatementContext) string {
	if ctx == nil || ctx.AssignmentOperator() == nil {
		return ""
	}

	return ctx.AssignmentOperator().GetText()
}

func (c *StmtCompiler) snapshotBindingValue(binding *core.Variable) bytecode.Operand {
	if c == nil || c.ctx == nil || binding == nil {
		return bytecode.NoopOperand
	}

	if binding.Storage == core.BindingStorageCell {
		return loadBindingValue(c.ctx, binding)
	}

	snapshot := c.ctx.Registers.Allocate()
	c.ctx.EmitMoveAuto(snapshot, binding.Register)

	return snapshot
}

func (c *StmtCompiler) storeBindingValue(binding *core.Variable, src bytecode.Operand, publishedType core.ValueType) bytecode.Operand {
	if c == nil || c.ctx == nil || binding == nil {
		return bytecode.NoopOperand
	}

	if binding.Storage == core.BindingStorageCell {
		src = c.ctx.ExprCompiler.ensureRegister(src)
		c.ctx.Emitter.EmitStoreCell(binding.Register, src)
		return binding.Register
	}

	if src.IsConstant() {
		c.ctx.Emitter.EmitLoadConst(binding.Register, src)
	} else {
		c.ctx.EmitMoveAuto(binding.Register, src)
	}

	c.ctx.Types.Set(binding.Register, publishedType)

	return binding.Register
}

func (c *StmtCompiler) reportInvalidAssignmentTarget(ctx antlr.ParserRuleContext) {
	if ctx == nil {
		return
	}

	err := c.ctx.Errors.Create(parserd.SyntaxError, ctx, "Assignment target must be a local variable name")
	err.Hint = "Property and index assignment are not supported. Use UPDATE for structural changes."
	c.ctx.Errors.Add(err)
}
