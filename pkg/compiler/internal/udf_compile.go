package internal

import (
	"github.com/antlr4-go/antlr/v4"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/compiler/internal/core"
	"github.com/MontFerret/ferret/v2/pkg/parser/fql"
)

// UDFCompiler compiles user-defined functions into bytecode.
type UDFCompiler struct {
	ctx      *CompilationSession
	calls    *CallResolver
	exprs    *ExprCompiler
	facts    *TypeFacts
	recovery *RecoveryCompiler
	stmts    *StatementCompiler
}

func NewUDFCompiler(ctx *CompilationSession) *UDFCompiler {
	return &UDFCompiler{ctx: ctx}
}

func (c *UDFCompiler) bind(calls *CallResolver, exprs *ExprCompiler, facts *TypeFacts, recovery *RecoveryCompiler, stmts *StatementCompiler) {
	if c == nil {
		return
	}

	c.calls = calls
	c.exprs = exprs
	c.facts = facts
	c.recovery = recovery
	c.stmts = stmts
}

func (c *UDFCompiler) CompileAll() {
	if c == nil || c.ctx == nil || c.ctx.Program.UDFs == nil {
		return
	}

	for _, fn := range c.ctx.Program.UDFs.Functions {
		c.compile(fn)
	}
}

func (c *UDFCompiler) compile(fn *core.UDFInfo) {
	if c == nil || c.ctx == nil || fn == nil || fn.Decl == nil {
		return
	}

	c.withFunctionCompileState(fn, func() {
		fn.Entry = c.ctx.Program.Emitter.Size()

		c.ctx.Function.Symbols.EnterScope()

		for _, name := range fn.Params {
			c.ctx.Function.Symbols.DeclareLocal(name, core.TypeAny)
		}

		for _, capture := range fn.Captures {
			c.ctx.Function.Symbols.DeclareLocalWithOptions(capture.Name, core.TypeAny, core.BindingOptions{
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
					c.stmts.CompileFunctionStatement(stmt)
				}

				c.compileReturn(block.FunctionReturn())
			}
		}

		fn.Registers = c.ctx.Function.Registers.Size()
	})
}

// withFunctionCompileState isolates function-local compiler state for a single
// UDF body compilation. It swaps in a fresh FunctionContext and restores the
// outer one when compilation finishes.
func (c *UDFCompiler) withFunctionCompileState(fn *core.UDFInfo, compile func()) {
	if c == nil || c.ctx == nil || fn == nil || compile == nil {
		return
	}

	outerFunction := c.ctx.Function

	var outerParams []string
	if outerFunction.Symbols != nil {
		outerParams = outerFunction.Symbols.Params()
	}

	// Create fresh function context for this UDF.
	localFunction := NewFunctionContext(c.ctx.Program.Constants)
	localFunction.UDFScope = fn.BodyScope

	c.ctx.Function = localFunction

	for _, name := range outerParams {
		c.ctx.Function.Symbols.BindParam(name)
	}

	defer func() {
		udfParams := c.ctx.Function.Symbols.Params()
		udfFunctions := c.ctx.Function.Symbols.Functions()

		c.ctx.Function = outerFunction

		for _, name := range udfParams {
			c.ctx.Function.Symbols.BindParam(name)
		}

		for name, args := range udfFunctions {
			c.ctx.Function.Symbols.BindFunction(name, args)
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

	if fce := directFunctionCall(expr); fce != nil && fce.ErrorOperator() == nil && allowsTailCallRecovery(c.recovery.CollectPlan(fce, core.RecoveryPlanOptions{})) {
		call := fce.FunctionCall()
		if call != nil {
			if fn, ok := c.calls.ResolveUDF(call); ok {
				seq := c.exprs.CompileArgumentList(call.ArgumentList())
				c.exprs.EmitUdfTailCall(fn, seq, call.(antlr.ParserRuleContext))
				return
			}
		}
	}

	val := ensureOperandRegister(c.ctx, c.facts, c.exprs.Compile(expr))

	c.ctx.Program.Emitter.EmitA(bytecode.OpReturn, val)
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
