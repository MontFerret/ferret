package internal

import (
	"github.com/antlr4-go/antlr/v4"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/compiler/internal/core"
	"github.com/MontFerret/ferret/v2/pkg/parser/fql"
)

// UDFCompiler compiles user-defined functions into bytecode.
type UDFCompiler struct {
	ctx   *CompilationSession
	front *CompilationFrontend
}

func NewUDFCompiler(ctx *CompilationSession) *UDFCompiler {
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

	c.withFunctionCompileState(fn, func() {
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
					c.front.Statements.CompileFunctionStatement(stmt)
				}

				c.compileReturn(block.FunctionReturn())
			}
		}

		fn.Registers = c.ctx.Registers.Size()
	})
}

// withFunctionCompileState is the single save/restore path for UDF-local
// session state. New UDF-local compiler fields must live on
// CompilationSession.functionCompileState so this value swap stays exhaustive.
func (c *UDFCompiler) withFunctionCompileState(fn *core.UDFInfo, compile func()) {
	if c == nil || c.ctx == nil || fn == nil || compile == nil {
		return
	}

	outerState := c.ctx.functionCompileState

	var outerParams []string
	if outerState.Symbols != nil {
		outerParams = outerState.Symbols.Params()
	}

	localState := functionCompileState{
		Registers: core.NewRegisterAllocator(),
		Types:     core.NewTypeTracker(),
		UDFScope:  fn.BodyScope,
	}
	localState.Symbols = core.NewSymbolTable(localState.Registers, c.ctx.Constants)
	localState.Loops = core.NewLoopTable(localState.Registers)

	c.ctx.functionCompileState = localState

	for _, name := range outerParams {
		c.ctx.Symbols.BindParam(name)
	}

	defer func() {
		udfParams := c.ctx.Symbols.Params()
		udfFunctions := c.ctx.Symbols.Functions()

		c.ctx.functionCompileState = outerState

		for _, name := range udfParams {
			c.ctx.Symbols.BindParam(name)
		}

		for name, args := range udfFunctions {
			c.ctx.Symbols.BindFunction(name, args)
		}
	}()

	compile()
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

	if fce := directFunctionCall(expr); fce != nil && fce.ErrorOperator() == nil && allowsTailCallRecovery(c.front.Recovery.CollectPlan(fce, core.RecoveryPlanOptions{})) {
		call := fce.FunctionCall()
		if call != nil {
			if fn, ok := c.front.Calls.ResolveUDF(call); ok {
				seq := c.front.Expressions.CompileArgumentList(call.ArgumentList())
				c.front.Expressions.EmitUdfTailCall(fn, seq, call.(antlr.ParserRuleContext))
				return
			}
		}
	}

	val := c.front.Expressions.ensureRegister(c.front.Expressions.Compile(expr))

	c.ctx.Emitter.EmitA(bytecode.OpReturn, val)
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

	if atomCtx.ExpressionAtom(0) != nil || atomCtx.ExpressionAtom(1) != nil {
		return nil
	}

	return atomCtx.FunctionCallExpression()
}
