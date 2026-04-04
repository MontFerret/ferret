package internal

import (
	"github.com/antlr4-go/antlr/v4"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/compiler/internal/core"
	"github.com/MontFerret/ferret/v2/pkg/parser/fql"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

func (c *ExprCompiler) ensureRegister(op bytecode.Operand) bytecode.Operand {
	return c.callCompiler.ensureRegister(op)
}

func (c *ExprCompiler) emitComparison(op bytecode.Opcode, left, right bytecode.Operand) bytecode.Operand {
	dst := c.ctx.Function.Registers.Allocate()
	c.ctx.Program.Emitter.EmitABC(op, dst, left, right)

	if dst.IsRegister() {
		c.ctx.Function.Types.Set(dst, core.TypeBool)
	}

	return dst
}

func (c *ExprCompiler) emitBooleanAnd(left, right bytecode.Operand) bytecode.Operand {
	dst := c.ctx.Function.Registers.Allocate()
	skip := c.ctx.Program.Emitter.NewLabel("and.false")
	done := c.ctx.Program.Emitter.NewLabel("and.done")

	c.ctx.Program.Emitter.EmitJumpIfFalse(left, skip)
	c.ctx.Program.Emitter.EmitJumpIfFalse(right, skip)
	c.ctx.Program.Emitter.EmitAb(bytecode.OpLoadBool, dst, true)
	c.ctx.Program.Emitter.EmitJump(done)

	c.ctx.Program.Emitter.MarkLabel(skip)
	c.ctx.Program.Emitter.EmitAb(bytecode.OpLoadBool, dst, false)
	c.ctx.Program.Emitter.MarkLabel(done)

	if dst.IsRegister() {
		c.ctx.Function.Types.Set(dst, core.TypeBool)
	}

	return dst
}

func (c *ExprCompiler) CompileVariable(ctx fql.IVariableContext) bytecode.Operand {
	return c.callCompiler.compileVariable(ctx)
}

func (c *ExprCompiler) CompileParam(ctx fql.IParamContext) bytecode.Operand {
	return c.callCompiler.compileParam(ctx)
}

func (c *ExprCompiler) CompileFunctionCallExpression(ctx fql.IFunctionCallExpressionContext) bytecode.Operand {
	return c.callCompiler.compileFunctionCallExpression(ctx)
}

func (c *ExprCompiler) CompileFunctionCall(ctx fql.IFunctionCallContext, protected bool) bytecode.Operand {
	return c.callCompiler.compileFunctionCall(ctx, protected)
}

func (c *ExprCompiler) CompileFunctionCallWith(ctx fql.IFunctionCallContext, protected bool, seq core.RegisterSequence) bytecode.Operand {
	return c.callCompiler.compileFunctionCallWith(ctx, protected, seq)
}

func (c *ExprCompiler) CompileFunctionCallByNameWith(ctx fql.IFunctionCallContext, name runtime.String, protected bool, seq core.RegisterSequence) bytecode.Operand {
	return c.callCompiler.compileFunctionCallByNameWith(ctx, name, protected, seq)
}

func (c *ExprCompiler) EmitUdfTailCall(fn *core.UDFInfo, seq core.RegisterSequence, callCtx antlr.ParserRuleContext) {
	c.callCompiler.emitUdfTailCall(fn, seq, callCtx)
}

func (c *ExprCompiler) CompileArgumentList(ctx fql.IArgumentListContext) core.RegisterSequence {
	return c.callCompiler.compileArgumentList(ctx)
}

func (c *ExprCompiler) CompileRangeOperator(ctx fql.IRangeOperatorContext) bytecode.Operand {
	return c.callCompiler.compileRangeOperator(ctx)
}
