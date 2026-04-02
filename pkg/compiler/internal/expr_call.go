package internal

import (
	"fmt"
	"strings"

	"github.com/antlr4-go/antlr/v4"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/compiler/internal/core"
	"github.com/MontFerret/ferret/v2/pkg/parser/diagnostics"
	"github.com/MontFerret/ferret/v2/pkg/parser/fql"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/source"
)

func (c *ExprCompiler) ensureRegister(op bytecode.Operand) bytecode.Operand {
	if op == bytecode.NoopOperand {
		return op
	}

	if op.IsRegister() {
		return op
	}

	reg := c.ctx.Registers.Allocate()
	c.ctx.Emitter.EmitLoadConst(reg, op)
	c.ctx.Types.Set(reg, c.front.TypeFacts.OperandType(op))

	return reg
}

func (c *ExprCompiler) emitComparison(op bytecode.Opcode, left, right bytecode.Operand) bytecode.Operand {
	dst := c.ctx.Registers.Allocate()
	c.ctx.Emitter.EmitABC(op, dst, left, right)

	if dst.IsRegister() {
		c.ctx.Types.Set(dst, core.TypeBool)
	}

	return dst
}

func (c *ExprCompiler) emitBooleanAnd(left, right bytecode.Operand) bytecode.Operand {
	dst := c.ctx.Registers.Allocate()
	skip := c.ctx.Emitter.NewLabel("and.false")
	done := c.ctx.Emitter.NewLabel("and.done")

	c.ctx.Emitter.EmitJumpIfFalse(left, skip)
	c.ctx.Emitter.EmitJumpIfFalse(right, skip)
	c.ctx.Emitter.EmitAb(bytecode.OpLoadBool, dst, true)
	c.ctx.Emitter.EmitJump(done)

	c.ctx.Emitter.MarkLabel(skip)
	c.ctx.Emitter.EmitAb(bytecode.OpLoadBool, dst, false)
	c.ctx.Emitter.MarkLabel(done)

	if dst.IsRegister() {
		c.ctx.Types.Set(dst, core.TypeBool)
	}

	return dst
}

func (c *ExprCompiler) CompileVariable(ctx fql.IVariableContext) bytecode.Operand {
	var name string
	token := ctx.GetStart()

	if id := ctx.Identifier(); id != nil {
		name = id.GetText()
		token = id.GetSymbol()
	} else if srw := ctx.SafeReservedWord(); srw != nil {
		name = srw.GetText()
	} else {
		return bytecode.NoopOperand
	}

	binding, found := c.ctx.Symbols.ResolveBinding(name)
	if !found {
		c.ctx.Errors.VariableNotFound(token, name)

		return bytecode.NoopOperand
	}

	op := c.front.Bindings.LoadBindingValue(binding)

	if op.IsRegister() {
		return op
	}

	return c.ensureRegister(op)
}

func (c *ExprCompiler) CompileParam(ctx fql.IParamContext) bytecode.Operand {
	var name string

	if id := ctx.Identifier(); id != nil {
		name = id.GetText()
	} else if srw := ctx.SafeReservedWord(); srw != nil {
		name = srw.GetText()
	} else {
		return bytecode.NoopOperand
	}

	reg := c.ctx.Registers.Allocate()
	c.ctx.Emitter.EmitLoadParam(reg, c.ctx.Symbols.BindParam(name))
	c.ctx.Types.Set(reg, core.TypeAny)

	return reg
}

func (c *ExprCompiler) CompileFunctionCallExpression(ctx fql.IFunctionCallExpressionContext) bytecode.Operand {
	if ctx == nil {
		return bytecode.NoopOperand
	}

	call := ctx.FunctionCall()
	if ctx.ErrorOperator() != nil {
		return c.front.Recovery.CompileWithErrorPolicy(core.ErrorPolicySuppress, core.CatchJumpModeNone, func() bytecode.Operand {
			return c.CompileFunctionCall(call, true)
		})
	}

	plan := c.front.Recovery.CollectPlan(ctx, core.RecoveryPlanOptions{})
	if hasErrorReturnNoneHandler(plan) {
		out := c.front.Recovery.CompileWithErrorPolicy(core.ErrorPolicySuppress, core.CatchJumpModeNone, func() bytecode.Operand {
			return c.CompileFunctionCall(call, true)
		})

		return c.front.Recovery.WidenResultType(out, plan)
	}

	return c.front.Recovery.CompileWithRecoveryPlan(plan, core.CatchJumpModeNone, func() bytecode.Operand {
		return c.CompileFunctionCall(call, false)
	})
}

func (c *ExprCompiler) CompileFunctionCall(ctx fql.IFunctionCallContext, protected bool) bytecode.Operand {
	return c.CompileFunctionCallWith(ctx, protected, c.CompileArgumentList(ctx.ArgumentList()))
}

func (c *ExprCompiler) CompileFunctionCallWith(ctx fql.IFunctionCallContext, protected bool, seq core.RegisterSequence) bytecode.Operand {
	name := c.front.Calls.ResolveFunctionName(ctx)
	span := source.Span{Start: -1, End: -1}

	if ctx != nil {
		if fn := ctx.FunctionName(); fn != nil {
			span = diagnostics.SpanFromRuleContext(fn)

			if ns := ctx.Namespace(); ns != nil && ns.GetStart() != nil {
				span.Start = ns.GetStart().GetStart()
			}
		} else if prc, ok := ctx.(antlr.ParserRuleContext); ok {
			span = diagnostics.SpanFromRuleContext(prc)
		}
	}

	var out bytecode.Operand

	c.ctx.Emitter.WithSpan(span, func() {
		out = c.CompileFunctionCallByNameWith(ctx, name, protected, seq)
	})

	return out
}

func (c *ExprCompiler) CompileFunctionCallByNameWith(ctx fql.IFunctionCallContext, name runtime.String, protected bool, seq core.RegisterSequence) bytecode.Operand {
	nameStr := name.String()
	builtinName := strings.ToUpper(nameStr)

	namespaced := strings.Contains(nameStr, runtime.NamespaceSeparator)
	if ctx != nil {
		if ns := ctx.Namespace(); ns != nil && ns.GetText() != "" {
			namespaced = true
		}
	}

	var callCtx antlr.ParserRuleContext
	if ctx != nil {
		if prc, ok := ctx.(antlr.ParserRuleContext); ok {
			callCtx = prc
		}
	}

	if !namespaced && c.ctx.UDFs != nil && c.ctx.UDFScope != nil {
		if fn, ok := c.front.Calls.ResolveUDF(ctx); ok {
			return c.compileUdfCallWith(fn, protected, seq, callCtx)
		}
	}

	if !namespaced {
		switch builtinName {
		case runtimeLength:
			dst := c.ctx.Registers.Allocate()

			if seq == nil || len(seq) != 1 {
				panic(runtime.Error(runtime.ErrInvalidArgument, runtimeLength+": expected 1 argument"))
			}

			c.ctx.Emitter.EmitAB(bytecode.OpLength, dst, seq[0])

			return dst
		case runtimeTypename:
			dst := c.ctx.Registers.Allocate()

			if seq == nil || len(seq) != 1 {
				panic(runtime.Error(runtime.ErrInvalidArgument, runtimeTypename+": expected 1 argument"))
			}

			c.ctx.Emitter.EmitAB(bytecode.OpType, dst, seq[0])

			return dst
		case runtimeWait:
			if len(seq) != 1 {
				panic(runtime.Error(runtime.ErrInvalidArgument, runtimeWait+": expected 1 argument"))
			}

			c.ctx.Emitter.EmitA(bytecode.OpSleep, seq[0])

			return seq[0]
		}
	}

	return c.compileHostFunctionCallWith(name, protected, seq)
}

func (c *ExprCompiler) compileHostFunctionCallWith(name runtime.String, protected bool, seq core.RegisterSequence) bytecode.Operand {
	dest := c.ctx.Registers.Allocate()
	c.ctx.Emitter.EmitLoadConst(dest, c.ctx.Symbols.AddConstant(name))
	c.ctx.Symbols.BindFunction(name.String(), len(seq))

	opcode := bytecode.OpHCall
	if protected {
		opcode = bytecode.OpProtectedHCall
	}

	c.ctx.Emitter.EmitAs(opcode, dest, seq)

	c.ctx.Types.Set(dest, core.TypeAny)

	return dest
}

func (c *ExprCompiler) compileUdfCallWith(fn *core.UDFInfo, protected bool, seq core.RegisterSequence, callCtx antlr.ParserRuleContext) bytecode.Operand {
	args := c.prepareUdfCallArgs(fn, seq, callCtx)

	dest := c.ctx.Registers.Allocate()
	c.ctx.Emitter.EmitLoadConst(dest, c.ctx.Symbols.AddConstant(runtime.NewInt(fn.ID)))

	opcode := bytecode.OpCall
	if protected {
		opcode = bytecode.OpProtectedCall
	}

	c.ctx.Emitter.EmitAs(opcode, dest, args)

	c.ctx.Types.Set(dest, core.TypeAny)

	return dest
}

func (c *ExprCompiler) EmitUdfTailCall(fn *core.UDFInfo, seq core.RegisterSequence, callCtx antlr.ParserRuleContext) {
	args := c.prepareUdfCallArgs(fn, seq, callCtx)

	dest := c.ctx.Registers.Allocate()
	c.ctx.Emitter.EmitLoadConst(dest, c.ctx.Symbols.AddConstant(runtime.NewInt(fn.ID)))

	c.ctx.Emitter.EmitAs(bytecode.OpTailCall, dest, args)
}

func (c *ExprCompiler) prepareUdfCallArgs(fn *core.UDFInfo, seq core.RegisterSequence, callCtx antlr.ParserRuleContext) core.RegisterSequence {
	if fn == nil {
		return seq
	}

	if len(seq) != len(fn.Params) && c.ctx.Errors != nil {
		ctx := callCtx
		if ctx == nil && fn.Decl != nil {
			if prc, ok := fn.Decl.(antlr.ParserRuleContext); ok {
				ctx = prc
			}
		}

		if ctx != nil {
			name := fn.DisplayName
			if name == "" {
				name = fn.Name
			}

			c.ctx.Errors.Add(c.ctx.Errors.Create(diagnostics.NameError, ctx, fmt.Sprintf("Function '%s' expects %d arguments, got %d", name, len(fn.Params), len(seq))))
		}
	}

	if len(fn.Captures) == 0 {
		return seq
	}

	total := len(seq) + len(fn.Captures)
	args := c.ctx.Registers.AllocateSequence(total)

	for i, src := range seq {
		c.ctx.Emitter.EmitMove(args[i], src)
		c.ctx.Types.Set(args[i], c.front.TypeFacts.OperandType(src))
	}

	for i, capture := range fn.Captures {
		binding, ok := c.ctx.Symbols.ResolveBinding(capture.Name)
		if !ok {
			if callCtx != nil {
				c.ctx.Errors.VariableNotFound(callCtx.GetStart(), capture.Name)
			}
			continue
		}

		dst := args[len(seq)+i]

		if capture.Storage == core.BindingStorageCell {
			c.ctx.Emitter.EmitPlainMove(dst, binding.Register)
			c.ctx.Types.Set(dst, core.TypeAny)
			continue
		}

		src := c.front.Bindings.LoadBindingValue(binding)
		c.front.TypeFacts.EmitMoveAuto(dst, src)
		c.ctx.Types.Set(dst, c.front.TypeFacts.OperandType(src))
	}

	return args
}

func (c *ExprCompiler) CompileArgumentList(ctx fql.IArgumentListContext) core.RegisterSequence {
	var seq core.RegisterSequence

	if ctx == nil {
		return seq
	}

	exps := ctx.AllExpression()
	size := len(exps)

	if size > 0 {
		seq = c.ctx.Registers.AllocateSequence(size)

		for i, exp := range exps {
			if val, ok := c.front.TypeFacts.LiteralValueFromExpression(exp); ok && (bool(runtime.IsScalar(val)) || val == runtime.None) {
				c.ctx.Emitter.EmitLoadConst(seq[i], c.ctx.Symbols.AddConstant(val))
				c.ctx.Types.Set(seq[i], c.front.TypeFacts.ValueTypeFromRuntime(val))
				continue
			}

			srcReg := c.Compile(exp)

			if srcReg.IsConstant() {
				c.ctx.Emitter.EmitLoadConst(seq[i], srcReg)
			} else {
				c.ctx.Emitter.EmitMove(seq[i], srcReg)
			}
			c.ctx.Types.Set(seq[i], c.front.TypeFacts.OperandType(srcReg))
		}
	}

	return seq
}

func (c *ExprCompiler) CompileRangeOperator(ctx fql.IRangeOperatorContext) bytecode.Operand {
	dst := c.ctx.Registers.Allocate()
	start := c.compileRangeOperand(ctx.GetLeft())
	end := c.compileRangeOperand(ctx.GetRight())

	span := source.Span{Start: -1, End: -1}

	if prc, ok := ctx.(antlr.ParserRuleContext); ok {
		span = diagnostics.SpanFromRuleContext(prc)
	}

	c.ctx.Emitter.WithSpan(span, func() {
		c.ctx.Emitter.EmitRange(dst, start, end)
	})

	c.ctx.Types.Set(dst, core.TypeList)

	return dst
}

func (c *ExprCompiler) compileRangeOperand(ctx fql.IRangeOperandContext) bytecode.Operand {
	if v := ctx.Variable(); v != nil {
		return c.CompileVariable(v)
	}

	if p := ctx.Param(); p != nil {
		return c.CompileParam(p)
	}

	if il := ctx.IntegerLiteral(); il != nil {
		return c.front.Literals.CompileIntegerLiteral(il)
	}

	if me := ctx.MemberExpression(); me != nil {
		return c.CompileMemberExpression(me)
	}

	if ice := ctx.ImplicitCurrentExpression(); ice != nil {
		return c.CompileImplicitCurrentExpression(ice)
	}

	if ime := ctx.ImplicitMemberExpression(); ime != nil {
		return c.CompileImplicitMemberExpression(ime)
	}

	if fc := ctx.FunctionCallExpression(); fc != nil {
		return c.CompileFunctionCallExpression(fc)
	}

	return bytecode.NoopOperand
}
