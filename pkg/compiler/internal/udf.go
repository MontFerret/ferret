package internal

import (
	"github.com/antlr4-go/antlr/v4"

	"github.com/MontFerret/ferret/v2/pkg/compiler/internal/optimization"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/compiler/internal/core"
	"github.com/MontFerret/ferret/v2/pkg/parser/fql"
)

type (
	// UDFCompiler compiles user-defined functions into bytecode.
	UDFCompiler struct {
		ctx *CompilerContext
	}

	udfCompileState struct {
		registers *core.RegisterAllocator
		symbols   *core.SymbolTable
		types     *core.TypeTracker
		loops     *core.LoopTable
		udfScope  *core.UDFScope
	}
)

func NewUDFCompiler(ctx *CompilerContext) *UDFCompiler {
	return &UDFCompiler{ctx: ctx}
}

func (c *UDFCompiler) CompileAll() {
	if c == nil || c.ctx == nil || c.ctx.UDFs == nil {
		return
	}

	for _, fn := range c.ctx.UDFs.Functions {
		c.compile(fn)
	}
}

func (c *UDFCompiler) compile(fn *core.UDFInfo) {
	if c == nil || c.ctx == nil || fn == nil || fn.Decl == nil {
		return
	}

	state := udfCompileState{
		registers: c.ctx.Registers,
		symbols:   c.ctx.Symbols,
		types:     c.ctx.Types,
		loops:     c.ctx.Loops,
		udfScope:  c.ctx.UDFScope,
	}

	var outerParams []string
	if state.symbols != nil {
		outerParams = state.symbols.Params()
	}

	c.ctx.Registers = core.NewRegisterAllocator()
	c.ctx.Symbols = core.NewSymbolTable(c.ctx.Registers, c.ctx.Constants)

	// Pre-seed param bindings from the outer table so UDF-emitted LOADP slots
	// remain aligned with program-level param ordering.
	for _, name := range outerParams {
		c.ctx.Symbols.BindParam(name)
	}

	c.ctx.Types = core.NewTypeTracker()
	c.ctx.Loops = core.NewLoopTable(c.ctx.Registers)
	c.ctx.UDFScope = fn.BodyScope

	fn.Entry = c.ctx.Emitter.Size()

	c.ctx.Symbols.EnterScope()

	for _, name := range fn.Params {
		c.ctx.Symbols.DeclareLocal(name, core.TypeAny)
	}

	for _, capture := range fn.Captures {
		c.ctx.Symbols.DeclareLocalWithOptions(capture.Name, core.TypeAny, core.BindingOptions{
			Mutable: capture.Mutable,
			Storage: capture.Storage,
		})
	}

	body := fn.Decl.FunctionBody()
	if body != nil {
		if arrow := body.FunctionArrow(); arrow != nil {
			c.compileExpressionReturn(arrow.Expression())
		} else if block := body.FunctionBlock(); block != nil {
			for _, stmt := range block.AllFunctionStatement() {
				c.ctx.StmtCompiler.CompileFunctionStatement(stmt)
			}

			c.compileReturn(block.FunctionReturn())
		}
	}

	fn.Registers = c.ctx.Registers.Size()

	// Preserve metadata discovered while compiling UDF bodies. UDF compilation
	// uses a temporary symbol table pre-seeded with outer param bindings, and
	// then merges any newly discovered params/host function bindings back.
	udfParams := c.ctx.Symbols.Params()
	udfFunctions := c.ctx.Symbols.Functions()

	c.ctx.Registers = state.registers
	c.ctx.Symbols = state.symbols
	c.ctx.Types = state.types
	c.ctx.Loops = state.loops
	c.ctx.UDFScope = state.udfScope

	for _, name := range udfParams {
		c.ctx.Symbols.BindParam(name)
	}

	for name, args := range udfFunctions {
		c.ctx.Symbols.BindFunction(name, args)
	}
}

func (c *UDFCompiler) compileReturn(ctx fql.IFunctionReturnContext) {
	if ctx == nil {
		return
	}

	expr := ctx.Expression()
	c.compileExpressionReturn(expr)
}

func (c *UDFCompiler) compileExpressionReturn(expr fql.IExpressionContext) {
	if expr == nil {
		return
	}

	if fce := directFunctionCall(expr); fce != nil && fce.ErrorOperator() == nil && allowsTailCallRecovery(collectRecoveryPlan(c.ctx, fce, core.RecoveryPlanOptions{})) {
		call := fce.FunctionCall()
		if call != nil {
			if name, ok := getUDFName(call, c.ctx.UseAliases); ok {
				if fn, ok := c.ctx.UDFs.Resolve(name, c.ctx.UDFScope); ok {
					seq := c.ctx.ExprCompiler.CompileArgumentList(call.ArgumentList())
					c.ctx.ExprCompiler.EmitUdfTailCall(fn, seq, call.(antlr.ParserRuleContext))
					return
				}
			}
		}
	}

	val := c.ctx.ExprCompiler.ensureRegister(c.ctx.ExprCompiler.Compile(expr))

	c.ctx.Emitter.EmitA(bytecode.OpReturn, val)
}

func CollectUDFs(ctx *CompilerContext, program *fql.ProgramContext) *core.UDFTable {
	table := core.NewUDFTable()
	table.GlobalScope = core.NewUDFScope(nil)

	if program == nil || program.Body() == nil {
		return table
	}

	body, ok := program.Body().(*fql.BodyContext)
	if !ok {
		return table
	}

	top := collectScopeFunctionsFromBody(ctx, table, table.GlobalScope, body)

	for _, fn := range top {
		collectNestedFunctions(ctx, table, fn)
	}

	if ctx != nil && ctx.OptimizationLevel > optimization.LevelNone {
		pruneUnusedUDFs(ctx, table, body)
	}

	analyzeCaptures(ctx, table, body)

	return table
}

func directFunctionCall(expr fql.IExpressionContext) fql.IFunctionCallExpressionContext {
	if expr == nil {
		return nil
	}

	expCtx, ok := expr.(*fql.ExpressionContext)
	if !ok {
		return nil
	}

	if expCtx.GetLeft() != nil || expCtx.GetRight() != nil || expCtx.GetCondition() != nil || expCtx.GetOnTrue() != nil || expCtx.GetOnFalse() != nil {
		return nil
	}

	pred := expCtx.Predicate()
	if pred == nil {
		return nil
	}

	predCtx, ok := pred.(*fql.PredicateContext)
	if !ok {
		return nil
	}

	if predCtx.GetLeft() != nil || predCtx.GetRight() != nil {
		return nil
	}

	atom := predCtx.ExpressionAtom()
	if atom == nil {
		return nil
	}

	atomCtx, ok := atom.(*fql.ExpressionAtomContext)
	if !ok {
		return nil
	}

	if atomCtx.GetLeft() != nil || atomCtx.GetRight() != nil {
		return nil
	}

	return atomCtx.FunctionCallExpression()
}
