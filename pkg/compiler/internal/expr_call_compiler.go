package internal

import (
	"fmt"

	"github.com/antlr4-go/antlr/v4"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/compiler/internal/core"
	"github.com/MontFerret/ferret/v2/pkg/parser/diagnostics"
	"github.com/MontFerret/ferret/v2/pkg/parser/fql"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/source"
)

type (
	exprCallCallbacks struct {
		compileExpr            func(fql.IExpressionContext) bytecode.Operand
		compileMember          func(fql.IMemberExpressionContext) bytecode.Operand
		compileImplicitCurrent func(fql.IImplicitCurrentExpressionContext) bytecode.Operand
		compileImplicitMember  func(fql.IImplicitMemberExpressionContext) bytecode.Operand
	}

	exprCallCompiler struct {
		ctx       *CompilationSession
		bindings  *BindingCompiler
		calls     *CallResolver
		literals  *LiteralCompiler
		recovery  *RecoveryCompiler
		facts     *TypeFacts
		callbacks exprCallCallbacks
	}
)

func newExprCallCompiler(ctx *CompilationSession, callbacks exprCallCallbacks) *exprCallCompiler {
	return &exprCallCompiler{
		ctx:       ctx,
		callbacks: callbacks,
	}
}

func (c *exprCallCompiler) bind(
	bindings *BindingCompiler,
	calls *CallResolver,
	literals *LiteralCompiler,
	recovery *RecoveryCompiler,
	facts *TypeFacts,
) {
	if c == nil {
		return
	}

	c.bindings = bindings
	c.calls = calls
	c.literals = literals
	c.recovery = recovery
	c.facts = facts
}

func (c *exprCallCompiler) ensureRegister(op bytecode.Operand) bytecode.Operand {
	return ensureOperandRegister(c.ctx, c.facts, op)
}

func (c *exprCallCompiler) compileVariable(ctx fql.IVariableContext) bytecode.Operand {
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

	binding, found := c.ctx.Function.Symbols.ResolveBinding(name)
	if !found {
		reportVariableNotFound(c.ctx, token, name)

		return bytecode.NoopOperand
	}

	if c.ctx.Program.Semantics != nil {
		c.ctx.Program.Semantics.RecordBindingReference(binding.ID, diagnostics.SpanFromRuleContext(ctx.(antlr.ParserRuleContext)))
	}

	op := c.bindings.LoadBindingValue(binding)

	if op.IsRegister() {
		return op
	}

	return c.ensureRegister(op)
}

func (c *exprCallCompiler) compileParam(ctx fql.IParamContext) bytecode.Operand {
	var name string

	if id := ctx.Identifier(); id != nil {
		name = id.GetText()
	} else if srw := ctx.SafeReservedWord(); srw != nil {
		name = srw.GetText()
	} else {
		return bytecode.NoopOperand
	}

	reg := c.ctx.Function.Registers.Allocate()
	span := diagnostics.SpanFromRuleContext(ctx)
	c.ctx.Program.Emitter.WithSpan(span, func() {
		c.ctx.Program.Emitter.EmitLoadParam(reg, c.ctx.Program.HostParams.Bind(name))
	})
	c.ctx.Function.Types.Set(reg, core.TypeAny)

	if c.ctx.Program.Semantics != nil {
		c.ctx.Program.Semantics.RecordBindParameter(name, span)
	}

	return reg
}

func (c *exprCallCompiler) compileFunctionCallExpression(ctx fql.IFunctionCallExpressionContext) bytecode.Operand {
	if ctx == nil {
		return bytecode.NoopOperand
	}

	call := ctx.FunctionCall()
	if ctx.ErrorOperator() != nil {
		return c.recovery.CompileWithErrorPolicy(core.ErrorPolicySuppress, core.CatchJumpModeNone, func() bytecode.Operand {
			return c.compileFunctionCall(call, true)
		})
	}

	return c.recovery.CompileOperation(OperationRecoverySpec{
		Owner:    ctx,
		JumpMode: core.CatchJumpModeNone,
		CompileSuppressed: func() bytecode.Operand {
			return c.compileFunctionCall(call, true)
		},
		CompilePlain: func() bytecode.Operand {
			return c.compileFunctionCall(call, false)
		},
	})
}

func (c *exprCallCompiler) compileFunctionCallNoRecoveryExpression(ctx fql.IFunctionCallNoRecoveryExpressionContext) bytecode.Operand {
	if ctx == nil {
		return bytecode.NoopOperand
	}

	call := ctx.FunctionCall()
	if ctx.ErrorOperator() != nil {
		return c.recovery.CompileWithErrorPolicy(core.ErrorPolicySuppress, core.CatchJumpModeNone, func() bytecode.Operand {
			return c.compileFunctionCall(call, true)
		})
	}

	return c.compileFunctionCall(call, false)
}

func (c *exprCallCompiler) compileFunctionCall(ctx fql.IFunctionCallContext, protected bool) bytecode.Operand {
	return c.compileFunctionCallWith(ctx, protected, c.compileArgumentList(ctx.ArgumentList()))
}

func (c *exprCallCompiler) compileFunctionCallWith(ctx fql.IFunctionCallContext, protected bool, seq core.RegisterSequence) bytecode.Operand {
	name := c.calls.ResolveFunctionName(ctx)
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

	c.ctx.Program.Emitter.WithSpan(span, func() {
		out = c.compileFunctionCallByNameWith(ctx, name, protected, seq)
	})

	return out
}

func (c *exprCallCompiler) compileFunctionCallByNameWith(ctx fql.IFunctionCallContext, name runtime.String, protected bool, seq core.RegisterSequence) bytecode.Operand {
	argSpans := c.argumentSpansFromList(nil)
	if ctx != nil {
		argSpans = c.argumentSpansFromList(ctx.ArgumentList())
	}

	var callCtx antlr.ParserRuleContext
	if ctx != nil {
		if prc, ok := ctx.(antlr.ParserRuleContext); ok {
			callCtx = prc
		}
	}

	resolved := c.calls.resolveCall(ctx, name)
	var out bytecode.Operand

	switch resolved.Kind {
	case resolvedCallUDF:
		out = c.compileUdfCallWith(resolved.Function, protected, seq, callCtx, argSpans)
	case resolvedCallBuiltin:
		switch resolved.Name {
		case runtimeLength:
			dst := c.ctx.Function.Registers.Allocate()

			if len(seq) != 1 {
				out = c.reportFunctionArityError(callCtx, resolved.Name, 1, len(seq))

				break
			}

			c.ctx.Program.Emitter.EmitAB(bytecode.OpLength, dst, seq[0])
			out = dst
		case runtimeTypename:
			dst := c.ctx.Function.Registers.Allocate()

			if len(seq) != 1 {
				out = c.reportFunctionArityError(callCtx, resolved.Name, 1, len(seq))

				break
			}

			c.ctx.Program.Emitter.EmitAB(bytecode.OpType, dst, seq[0])
			out = dst
		case runtimeWait:
			if len(seq) != 1 {
				out = c.reportFunctionArityError(callCtx, resolved.Name, 1, len(seq))

				break
			}

			c.ctx.Program.Emitter.EmitA(bytecode.OpSleep, seq[0])
			out = seq[0]
		}
	case resolvedCallHost:
		out = c.compileHostFunctionCallWith(runtime.NewString(resolved.Name), protected, seq, argSpans)
	}

	if c.ctx.Program.Semantics != nil {
		c.recordCall(ctx, resolved, out, argSpans)
	}

	return out
}

func (c *exprCallCompiler) reportFunctionArityError(ctx antlr.ParserRuleContext, name string, expected, got int) bytecode.Operand {
	if c == nil || c.ctx == nil || c.ctx.Program.Errors == nil || ctx == nil {
		core.PanicInvariantf("cannot report arity error for function %q", name)
	}

	c.ctx.Program.Errors.Add(c.ctx.Program.Errors.Create(
		diagnostics.NameError,
		ctx,
		fmt.Sprintf("Function '%s' expects %d arguments, got %d", name, expected, got),
	))

	return bytecode.NoopOperand
}

func (c *exprCallCompiler) compileHostFunctionCallWith(name runtime.String, protected bool, seq core.RegisterSequence, argSpans []source.Span) bytecode.Operand {
	dest := c.ctx.Function.Registers.Allocate()
	bindingID := c.ctx.Program.HostFunctions.Bind(name.String(), len(seq))
	c.ctx.Program.Emitter.EmitLoadConst(dest, c.ctx.Function.Symbols.AddConstant(runtime.NewInt(bindingID)))

	opcode := bytecode.OpHCall
	if protected {
		opcode = bytecode.OpProtectedHCall
	}

	c.ctx.Program.Emitter.EmitAsWithCallArgumentSpans(opcode, dest, seq, argSpans)
	c.ctx.Function.Types.Set(dest, core.TypeAny)

	return dest
}

func (c *exprCallCompiler) compileUdfCallWith(fn *core.UDFInfo, protected bool, seq core.RegisterSequence, callCtx antlr.ParserRuleContext, argSpans []source.Span) bytecode.Operand {
	args := c.prepareUdfCallArgs(fn, seq, callCtx)
	dest := c.ctx.Function.Registers.Allocate()
	c.ctx.Program.Emitter.EmitLoadConst(dest, c.ctx.Function.Symbols.AddConstant(runtime.NewInt(fn.ID)))

	opcode := bytecode.OpCall
	if protected {
		opcode = bytecode.OpProtectedCall
	}

	c.ctx.Program.Emitter.EmitAsWithCallArgumentSpans(opcode, dest, args, argSpans)
	c.ctx.Function.Types.Set(dest, core.TypeAny)

	return dest
}

func (c *exprCallCompiler) emitUdfTailCall(fn *core.UDFInfo, seq core.RegisterSequence, callCtx antlr.ParserRuleContext) {
	args := c.prepareUdfCallArgs(fn, seq, callCtx)
	argSpans := c.argumentSpansFromCall(callCtx)

	dest := c.ctx.Function.Registers.Allocate()
	c.ctx.Program.Emitter.EmitLoadConst(dest, c.ctx.Function.Symbols.AddConstant(runtime.NewInt(fn.ID)))
	c.ctx.Program.Emitter.EmitAsWithCallArgumentSpans(bytecode.OpTailCall, dest, args, argSpans)

	if c.ctx.Program.Semantics != nil {
		if call, ok := callCtx.(fql.IFunctionCallContext); ok {
			c.recordCall(call, resolvedCall{Function: fn, Name: fn.DisplayName, Kind: resolvedCallUDF}, bytecode.NoopOperand, argSpans)
		}
	}
}

func (c *exprCallCompiler) recordCall(
	ctx fql.IFunctionCallContext,
	resolved resolvedCall,
	out bytecode.Operand,
	argumentSpans []source.Span,
) {
	if c == nil || c.ctx == nil || c.ctx.Program.Semantics == nil || ctx == nil {
		return
	}

	callCtx := ctx.(antlr.ParserRuleContext)
	calleeSpan := diagnostics.SpanFromRuleContext(callCtx)

	if fn := ctx.FunctionName(); fn != nil {
		calleeSpan = diagnostics.SpanFromRuleContext(fn.(antlr.ParserRuleContext))
		if ns := ctx.Namespace(); ns != nil && ns.GetStart() != nil {
			calleeSpan.Start = ns.GetStart().GetStart()
		}
	}

	call := SemanticCall{
		Name:          resolved.Name,
		Identity:      resolved.Identity,
		Span:          diagnostics.SpanFromRuleContext(callCtx),
		CalleeSpan:    calleeSpan,
		ArgumentSpans: argumentSpans,
		Type:          c.facts.OperandType(out),
	}

	switch resolved.Kind {
	case resolvedCallUDF:
		call.Kind = SemanticCallUserFunction
		call.Target = c.ctx.Program.Semantics.UserFunctionSymbol(resolved.Function)
		call.Type = core.TypeAny
		c.ctx.Program.Semantics.RecordUserFunctionReference(resolved.Function, calleeSpan)
	case resolvedCallBuiltin:
		call.Kind = SemanticCallBuiltin
	case resolvedCallHost:
		call.Kind = SemanticCallHost
	}

	c.ctx.Program.Semantics.RecordCall(call)
}

func (c *exprCallCompiler) prepareUdfCallArgs(fn *core.UDFInfo, seq core.RegisterSequence, callCtx antlr.ParserRuleContext) core.RegisterSequence {
	if fn == nil {
		return seq
	}

	if len(seq) != len(fn.Params) && c.ctx.Program.Errors != nil {
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

			c.ctx.Program.Errors.Add(c.ctx.Program.Errors.Create(diagnostics.NameError, ctx, fmt.Sprintf("Function '%s' expects %d arguments, got %d", name, len(fn.Params), len(seq))))
		}
	}

	if len(fn.Captures) == 0 {
		return seq
	}

	total := len(seq) + len(fn.Captures)
	args := c.ctx.Function.Registers.AllocateSequence(total)

	for i, src := range seq {
		c.ctx.Program.Emitter.EmitMove(args[i], src)
		c.ctx.Function.Types.Set(args[i], c.facts.OperandType(src))
	}

	for i, capture := range fn.Captures {
		binding, ok := c.ctx.Function.Symbols.ResolveBindingByID(capture.ID)
		if !ok {
			if callCtx != nil {
				c.ctx.Program.Errors.VariableNotFound(callCtx.GetStart(), capture.Name)
			}

			continue
		}

		dst := args[len(seq)+i]

		if capture.Storage == core.BindingStorageCell {
			c.ctx.Program.Emitter.EmitPlainMove(dst, binding.Register)
			c.ctx.Function.Types.Set(dst, core.TypeAny)

			continue
		}

		src := c.bindings.LoadBindingValue(binding)
		c.facts.EmitMoveAuto(dst, src)
		c.ctx.Function.Types.Set(dst, c.facts.OperandType(src))
	}

	return args
}

func (c *exprCallCompiler) compileArgumentList(ctx fql.IArgumentListContext) core.RegisterSequence {
	var seq core.RegisterSequence

	if ctx == nil {
		return seq
	}

	exps := ctx.AllExpression()
	size := len(exps)

	if size > 0 {
		seq = c.ctx.Function.Registers.AllocateSequence(size)

		for i, exp := range exps {
			if val, ok := c.facts.LiteralValueFromExpression(exp); ok && (bool(runtime.IsScalar(val)) || val == runtime.None) {
				c.ctx.Program.Emitter.EmitLoadConst(seq[i], c.ctx.Function.Symbols.AddConstant(val))
				c.ctx.Function.Types.Set(seq[i], c.facts.ValueTypeFromRuntime(val))

				continue
			}

			srcReg := c.callbacks.compileExpr(exp)

			if srcReg.IsConstant() {
				c.ctx.Program.Emitter.EmitLoadConst(seq[i], srcReg)
			} else {
				c.ctx.Program.Emitter.EmitMove(seq[i], srcReg)
			}

			c.ctx.Function.Types.Set(seq[i], c.facts.OperandType(srcReg))
		}
	}

	return seq
}

func (c *exprCallCompiler) argumentSpansFromCall(ctx antlr.ParserRuleContext) []source.Span {
	call, ok := ctx.(fql.IFunctionCallContext)
	if !ok {
		return nil
	}

	return c.argumentSpansFromList(call.ArgumentList())
}

func (c *exprCallCompiler) argumentSpansFromList(ctx fql.IArgumentListContext) []source.Span {
	if ctx == nil {
		return nil
	}

	exps := ctx.AllExpression()
	if len(exps) == 0 {
		return nil
	}

	spans := make([]source.Span, len(exps))

	for i, exp := range exps {
		spans[i] = source.Span{Start: -1, End: -1}

		prc, ok := exp.(antlr.ParserRuleContext)
		if !ok {
			continue
		}

		spans[i] = diagnostics.SpanFromRuleContext(prc)
	}

	return spans
}

func (c *exprCallCompiler) compileRangeOperator(ctx fql.IRangeOperatorContext) bytecode.Operand {
	dst := c.ctx.Function.Registers.Allocate()
	start := c.compileRangeOperand(ctx.GetLeft())
	end := c.compileRangeOperand(ctx.GetRight())

	span := source.Span{Start: -1, End: -1}

	if prc, ok := ctx.(antlr.ParserRuleContext); ok {
		span = diagnostics.SpanFromRuleContext(prc)
	}

	c.ctx.Program.Emitter.WithSpan(span, func() {
		c.ctx.Program.Emitter.EmitRange(dst, start, end)
	})

	c.ctx.Function.Types.Set(dst, core.TypeList)

	return dst
}

func (c *exprCallCompiler) compileRangeOperand(ctx fql.IRangeOperandContext) bytecode.Operand {
	if v := ctx.Variable(); v != nil {
		return c.compileVariable(v)
	}

	if p := ctx.Param(); p != nil {
		return c.compileParam(p)
	}

	if il := ctx.IntegerLiteral(); il != nil {
		return c.literals.CompileIntegerLiteral(il)
	}

	if me := ctx.MemberExpression(); me != nil {
		return c.callbacks.compileMember(me)
	}

	if ice := ctx.ImplicitCurrentExpression(); ice != nil {
		return c.callbacks.compileImplicitCurrent(ice)
	}

	if ime := ctx.ImplicitMemberExpression(); ime != nil {
		return c.callbacks.compileImplicitMember(ime)
	}

	if fc := ctx.FunctionCallExpression(); fc != nil {
		return c.compileFunctionCallExpression(fc)
	}

	return bytecode.NoopOperand
}
