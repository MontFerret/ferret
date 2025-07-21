package internal

import (
	"github.com/MontFerret/ferret/pkg/compiler/internal/core"
	"github.com/MontFerret/ferret/pkg/parser/fql"
	"github.com/MontFerret/ferret/pkg/runtime"
	"github.com/MontFerret/ferret/pkg/vm"
)

type WaitCompiler struct {
	ctx *CompilerContext
}

func NewWaitCompiler(ctx *CompilerContext) *WaitCompiler {
	return &WaitCompiler{
		ctx: ctx,
	}
}

func (wc *WaitCompiler) Compile(ctx fql.IWaitForExpressionContext) vm.Operand {
	wc.ctx.Symbols.EnterScope()

	srcReg := wc.CompileWaitForEventSource(ctx.WaitForEventSource())
	eventReg := wc.CompileWaitForEventName(ctx.WaitForEventName())

	var optsReg vm.Operand

	if opts := ctx.OptionsClause(); opts != nil {
		optsReg = wc.CompileOptionsClause(opts)
	}

	var timeoutReg vm.Operand

	if timeout := ctx.TimeoutClause(); timeout != nil {
		timeoutReg = wc.CompileTimeoutClauseContext(timeout)
	}

	streamReg := wc.ctx.Registers.Allocate(core.Temp)

	// We move the source object to the stream register in order to re-use it in OpStream
	wc.ctx.Emitter.EmitMove(streamReg, srcReg)
	wc.ctx.Emitter.EmitABC(vm.OpStream, streamReg, eventReg, optsReg)
	wc.ctx.Emitter.EmitAB(vm.OpStreamIter, streamReg, timeoutReg)

	var valReg vm.Operand
	start := wc.ctx.Emitter.NewLabel()
	end := wc.ctx.Emitter.NewLabel()

	wc.ctx.Emitter.MarkLabel(start)
	// Now we start iterating over the stream
	wc.ctx.Emitter.EmitIterNext(streamReg, end)

	if filter := ctx.FilterClause(); filter != nil {
		valReg = wc.ctx.Symbols.DeclareLocal(core.PseudoVariable, core.TypeUnknown)
		wc.ctx.Emitter.EmitAB(vm.OpIterValue, valReg, streamReg)

		cond := wc.ctx.ExprCompiler.Compile(filter.Expression())

		wc.ctx.Emitter.EmitJumpIfFalse(cond, start)

		// TODO: Do we need to use timeout here too? We can really get stuck in the loop if no event satisfies the filter
	}

	wc.ctx.Emitter.MarkLabel(end)
	// Clean up the stream
	wc.ctx.Emitter.EmitA(vm.OpClose, streamReg)
	wc.ctx.Symbols.ExitScope()

	return vm.NoopOperand
}

func (wc *WaitCompiler) CompileWaitForEventName(ctx fql.IWaitForEventNameContext) vm.Operand {
	if c := ctx.StringLiteral(); c != nil {
		return wc.ctx.LiteralCompiler.CompileStringLiteral(c)
	}

	if c := ctx.Variable(); c != nil {
		return wc.ctx.ExprCompiler.CompileVariable(c)
	}

	if c := ctx.Param(); c != nil {
		return wc.ctx.ExprCompiler.CompileParam(c)
	}

	if c := ctx.MemberExpression(); c != nil {
		return wc.ctx.ExprCompiler.CompileMemberExpression(c)
	}

	if c := ctx.FunctionCallExpression(); c != nil {
		return wc.ctx.ExprCompiler.CompileFunctionCallExpression(c)
	}

	panic(runtime.Error(core.ErrUnexpectedToken, ctx.GetText()))
}

func (wc *WaitCompiler) CompileWaitForEventSource(ctx fql.IWaitForEventSourceContext) vm.Operand {
	if c := ctx.Variable(); c != nil {
		return wc.ctx.ExprCompiler.CompileVariable(c)
	}

	if c := ctx.MemberExpression(); c != nil {
		return wc.ctx.ExprCompiler.CompileMemberExpression(c)
	}

	if c := ctx.FunctionCallExpression(); c != nil {
		return wc.ctx.ExprCompiler.CompileFunctionCallExpression(c)
	}

	panic(runtime.Error(core.ErrUnexpectedToken, ctx.GetText()))
}

func (wc *WaitCompiler) CompileOptionsClause(ctx fql.IOptionsClauseContext) vm.Operand {
	if c := ctx.ObjectLiteral(); c != nil {
		return wc.ctx.LiteralCompiler.CompileObjectLiteral(c)
	}

	panic(runtime.Error(core.ErrUnexpectedToken, ctx.GetText()))
}

func (wc *WaitCompiler) CompileTimeoutClauseContext(ctx fql.ITimeoutClauseContext) vm.Operand {
	if c := ctx.IntegerLiteral(); c != nil {
		return wc.ctx.LiteralCompiler.CompileIntegerLiteral(c)
	}

	if c := ctx.Variable(); c != nil {
		return wc.ctx.ExprCompiler.CompileVariable(c)
	}

	if c := ctx.Param(); c != nil {
		return wc.ctx.ExprCompiler.CompileParam(c)
	}

	if c := ctx.MemberExpression(); c != nil {
		return wc.ctx.ExprCompiler.CompileMemberExpression(c)
	}

	if c := ctx.FunctionCall(); c != nil {
		return wc.ctx.ExprCompiler.CompileFunctionCall(c, false)
	}

	panic(runtime.Error(core.ErrUnexpectedToken, ctx.GetText()))
}
