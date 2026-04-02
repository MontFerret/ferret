package internal

import (
	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/compiler/internal/core"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

func initRetryDelayState(ctx *CompilerContext, retry *core.RecoveryRetryPlan) core.RetryDelayState {
	if ctx == nil || retry == nil || !retry.HasDelay {
		return core.RetryDelayState{}
	}

	state := core.RetryDelayState{
		BaseReg:    ctx.Registers.Allocate(),
		CurrentReg: ctx.Registers.Allocate(),
		ReadyReg:   ctx.Registers.Allocate(),
	}

	ctx.Emitter.EmitBoolean(state.ReadyReg, false)

	return state
}

func compileDurationOperand(ctx *CompilerContext, clause core.DurationClause) bytecode.Operand {
	if ctx == nil || clause == nil {
		return bytecode.NoopOperand
	}

	if dl := clause.DurationLiteral(); dl != nil {
		val, err := parseDurationLiteral(dl.GetText())
		if err != nil {
			panic(err)
		}

		return loadConstant(ctx, val)
	}

	il := clause.IntegerLiteral()
	fl := clause.FloatLiteral()
	v := clause.Variable()
	p := clause.Param()
	me := clause.MemberExpression()
	fc := clause.FunctionCall()

	return compileFirstOperand(
		newOperandBranch(il != nil, func() bytecode.Operand { return ctx.LiteralCompiler.CompileIntegerLiteral(il) }),
		newOperandBranch(fl != nil, func() bytecode.Operand { return ctx.LiteralCompiler.CompileFloatLiteral(fl) }),
		newOperandBranch(v != nil, func() bytecode.Operand { return ctx.ExprCompiler.CompileVariable(v) }),
		newOperandBranch(p != nil, func() bytecode.Operand { return ctx.ExprCompiler.CompileParam(p) }),
		newOperandBranch(me != nil, func() bytecode.Operand { return ctx.ExprCompiler.CompileMemberExpression(me) }),
		newOperandBranch(fc != nil, func() bytecode.Operand { return ctx.ExprCompiler.CompileFunctionCall(fc, false) }),
	)
}

func emitBackoffUpdate(ctx *CompilerContext, strategy core.RetryBackoff, intervalReg, baseEveryReg bytecode.Operand) {
	switch strategy {
	case core.RetryBackoffLinear:
		ctx.Emitter.EmitABC(bytecode.OpAdd, intervalReg, intervalReg, baseEveryReg)
	case core.RetryBackoffExponential:
		twoReg := loadConstant(ctx, runtime.NewInt(2))
		ctx.Emitter.EmitABC(bytecode.OpMul, intervalReg, intervalReg, twoReg)
	default:
		return
	}
}
