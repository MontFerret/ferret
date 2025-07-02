package internal

import (
	"github.com/MontFerret/ferret/pkg/compiler/internal/core"
	"github.com/MontFerret/ferret/pkg/parser/fql"
	"github.com/MontFerret/ferret/pkg/vm"
)

type StmtCompiler struct {
	ctx *CompilerContext
}

func NewStmtCompiler(ctx *CompilerContext) *StmtCompiler {
	return &StmtCompiler{
		ctx: ctx,
	}
}

func (sc *StmtCompiler) Compile(ctx fql.IBodyContext) {
	for _, statement := range ctx.AllBodyStatement() {
		sc.CompileBodyStatement(statement)
	}

	sc.CompileBodyExpression(ctx.BodyExpression())
}

func (sc *StmtCompiler) CompileBodyStatement(ctx fql.IBodyStatementContext) {
	if c := ctx.VariableDeclaration(); c != nil {
		sc.CompileVariableDeclaration(c)
	} else if c := ctx.FunctionCallExpression(); c != nil {
		sc.CompileFunctionCall(c)
	} else if c := ctx.WaitForExpression(); c != nil {
		sc.ctx.WaitCompiler.Compile(c)
	}
}

func (sc *StmtCompiler) CompileBodyExpression(ctx fql.IBodyExpressionContext) {
	if c := ctx.ForExpression(); c != nil {
		out := sc.ctx.LoopCompiler.Compile(c)

		sc.ctx.Emitter.EmitA(vm.OpReturn, out)
	} else if c := ctx.ReturnExpression(); c != nil {
		valReg := sc.ctx.ExprCompiler.Compile(c.Expression())

		if valReg.IsConstant() {
			valC := valReg
			valReg = sc.ctx.Registers.Allocate(core.Temp)

			sc.ctx.Emitter.EmitAB(vm.OpLoadGlobal, valReg, valC)
		}

		sc.ctx.Emitter.EmitA(vm.OpReturn, valReg)
	}
}

func (sc *StmtCompiler) CompileVariableDeclaration(ctx fql.IVariableDeclarationContext) vm.Operand {
	name := core.IgnorePseudoVariable

	if id := ctx.Identifier(); id != nil {
		name = id.GetText()
	} else if reserved := ctx.SafeReservedWord(); reserved != nil {
		name = reserved.GetText()
	}

	src := sc.ctx.ExprCompiler.Compile(ctx.Expression())

	if name != core.IgnorePseudoVariable {
		var dest vm.Operand

		if src.IsConstant() {
			dest = sc.ctx.Symbols.DeclareGlobal(name)
			tmp := sc.ctx.Registers.Allocate(core.Temp)
			sc.ctx.Emitter.EmitAB(vm.OpLoadConst, tmp, src)
			sc.ctx.Emitter.EmitAB(vm.OpStoreGlobal, dest, tmp)
		} else if sc.ctx.Symbols.Scope() == 0 {
			dest = sc.ctx.Symbols.DeclareGlobal(name)
			sc.ctx.Emitter.EmitAB(vm.OpStoreGlobal, dest, src)
		} else {
			dest = sc.ctx.Symbols.DeclareLocal(name)
			sc.ctx.Emitter.EmitAB(vm.OpMove, dest, src)
		}

		return dest
	}

	return vm.NoopOperand
}

func (sc *StmtCompiler) CompileFunctionCall(ctx fql.IFunctionCallExpressionContext) vm.Operand {
	return sc.ctx.ExprCompiler.CompileFunctionCallExpression(ctx)
}
