package internal

import (
	"fmt"

	"github.com/antlr4-go/antlr/v4"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/compiler/internal/core"
	parserd "github.com/MontFerret/ferret/v2/pkg/parser/diagnostics"
	"github.com/MontFerret/ferret/v2/pkg/parser/fql"
)

type (
	BindingCompiler struct {
		ctx                  *CompilationSession
		exprs                *ExprCompiler
		facts                *TypeFacts
		promotedDeclarations map[antlr.ParserRuleContext]struct{}
	}
)

func NewBindingCompiler(ctx *CompilationSession) *BindingCompiler {
	return &BindingCompiler{
		ctx:                  ctx,
		promotedDeclarations: make(map[antlr.ParserRuleContext]struct{}),
	}
}

func (c *BindingCompiler) bind(exprs *ExprCompiler, facts *TypeFacts) {
	if c == nil {
		return
	}

	c.exprs = exprs
	c.facts = facts
}

// PromoteDeclaration marks a declaration that must be emitted as a cell-backed binding
// because a nested scope writes through the captured variable.
func (c *BindingCompiler) PromoteDeclaration(decl antlr.ParserRuleContext) {
	if c == nil || decl == nil {
		return
	}

	c.promotedDeclarations[decl] = struct{}{}
}

func (c *BindingCompiler) IsPromotedDeclaration(decl antlr.ParserRuleContext) bool {
	if c == nil || decl == nil {
		return false
	}

	_, ok := c.promotedDeclarations[decl]
	return ok
}

func (c *BindingCompiler) CompileVariableDeclaration(ctx fql.IVariableDeclarationContext) bytecode.Operand {
	if ctx == nil {
		return bytecode.NoopOperand
	}

	name := c.declarationName(ctx)
	if name == "" {
		name = core.IgnorePseudoVariable
	}

	decl := ctx.(antlr.ParserRuleContext)
	mutable := c.isMutableDeclaration(ctx)
	storage := c.declarationStorage(decl, mutable)
	src := c.exprs.Compile(ctx.Expression())
	srcType := c.facts.OperandType(src)

	if name == core.IgnorePseudoVariable {
		return bytecode.NoopOperand
	}

	opts := core.BindingOptions{
		Mutable: mutable,
		Storage: storage,
	}

	if storage == core.BindingStorageCell {
		src = ensureOperandRegister(c.ctx, c.facts, src)

		dest, ok := c.declareBinding(name, srcType, src, opts)
		if !ok {
			c.ctx.Program.Errors.VariableNotUnique(decl, name)
			return bytecode.NoopOperand
		}

		c.ctx.Program.Emitter.EmitMakeCell(dest, src)
		c.ctx.Function.Types.Set(dest, core.TypeAny)

		return dest
	}

	if src.IsConstant() {
		dest, ok := c.declareBinding(name, srcType, src, opts)
		if !ok {
			c.ctx.Program.Errors.VariableNotUnique(decl, name)
			return bytecode.NoopOperand
		}

		c.ctx.Program.Emitter.EmitLoadConst(dest, src)
		c.ctx.Function.Types.Set(dest, srcType)

		return dest
	}

	if !c.assignBinding(name, srcType, src, opts) {
		c.ctx.Program.Errors.VariableNotUnique(decl, name)
		return bytecode.NoopOperand
	}

	c.ctx.Function.Types.Set(src, srcType)
	return src
}

func (c *BindingCompiler) CompileAssignmentStatement(ctx fql.IAssignmentStatementContext) bytecode.Operand {
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

	binding, found := c.ctx.Function.Symbols.ResolveBinding(name)
	if !found {
		c.ctx.Program.Errors.VariableNotFound(stmt.GetStart(), name)
		return bytecode.NoopOperand
	}

	if !binding.Mutable {
		err := c.ctx.Program.Errors.Create(parserd.SemanticError, stmt, fmt.Sprintf("Variable '%s' cannot be reassigned", name))
		err.Hint = "Declare it with VAR if you need to update it."
		c.ctx.Program.Errors.Add(err)
		return bytecode.NoopOperand
	}

	operator := assignmentOperatorText(stmt)
	src := bytecode.NoopOperand

	if operator == "=" {
		src = c.exprs.Compile(stmt.Expression())
	} else if operator == "+=" && binding.Type == core.TypeString {
		left := c.snapshotBindingValue(binding)
		parts := append([]concatOperandSegment{{operand: left}}, buildConcatOperandSegmentsFromExpression(c.exprs, stmt.Expression())...)
		src = emitConcatOperandSegments(c.ctx, c.facts, parts)
	} else {
		op, ok := resolveArithmeticBinaryOperator(operator)
		if !ok {
			return bytecode.NoopOperand
		}

		left := c.snapshotBindingValue(binding)
		right := c.exprs.Compile(stmt.Expression())
		src = emitBinaryOperation(c.ctx, c.facts, stmt, op, left, right)
	}

	srcType := c.facts.OperandType(src)
	publishedType := srcType

	if c.ctx.Function.Loops.Depth() > 0 {
		publishedType = core.JoinValueTypes(binding.Type, srcType)
	}

	binding.Type = publishedType

	return c.storeBindingValue(binding, src, publishedType)
}

func (c *BindingCompiler) LoadBindingValue(binding *core.Variable) bytecode.Operand {
	if c == nil || c.ctx == nil || binding == nil {
		return bytecode.NoopOperand
	}

	if binding.Storage != core.BindingStorageCell {
		return binding.Register
	}

	dst := c.ctx.Function.Registers.Allocate()
	c.ctx.Program.Emitter.EmitLoadCell(dst, binding.Register)
	c.ctx.Function.Types.Set(dst, binding.Type)

	return dst
}

func (c *BindingCompiler) captureBindingForDeclaration(ctx fql.IVariableDeclarationContext) captureBindingInfo {
	if ctx == nil {
		return captureBindingInfo{}
	}

	return captureBindingInfo{
		Name:    c.declarationName(ctx),
		Mutable: c.isMutableDeclaration(ctx),
		Decl:    ctx.(antlr.ParserRuleContext),
	}
}

func (c *BindingCompiler) declarationName(ctx fql.IVariableDeclarationContext) string {
	if ctx == nil {
		return ""
	}

	if id := ctx.BindingIdentifier(); id != nil {
		return textOfBindingIdentifier(id)
	}

	if id := ctx.Identifier(); id != nil {
		return id.GetText()
	}

	if id := ctx.SafeReservedWord(); id != nil {
		return id.GetText()
	}

	return ""
}

func (c *BindingCompiler) isMutableDeclaration(ctx fql.IVariableDeclarationContext) bool {
	if ctx == nil {
		return false
	}

	decl, ok := ctx.(*fql.VariableDeclarationContext)
	return ok && decl.Var() != nil
}

func (c *BindingCompiler) declarationStorage(decl antlr.ParserRuleContext, mutable bool) core.BindingStorage {
	if decl == nil || !mutable {
		return core.BindingStorageValue
	}

	if c.IsPromotedDeclaration(decl) {
		return core.BindingStorageCell
	}

	return core.BindingStorageValue
}

func (c *BindingCompiler) declareBinding(
	name string,
	srcType core.ValueType,
	src bytecode.Operand,
	opts core.BindingOptions,
) (bytecode.Operand, bool) {
	if c.ctx.Function.Symbols.Scope() == 0 {
		return c.ctx.Function.Symbols.DeclareGlobalWithOptions(name, srcType, opts)
	}

	return c.ctx.Function.Symbols.DeclareLocalWithOptions(name, srcType, opts)
}

func (c *BindingCompiler) assignBinding(
	name string,
	srcType core.ValueType,
	src bytecode.Operand,
	opts core.BindingOptions,
) bool {
	if c.ctx.Function.Symbols.Scope() == 0 {
		return c.ctx.Function.Symbols.AssignGlobalWithOptions(name, srcType, src, opts)
	}

	return c.ctx.Function.Symbols.AssignLocalWithOptions(name, srcType, src, opts)
}

func (c *BindingCompiler) snapshotBindingValue(binding *core.Variable) bytecode.Operand {
	if c == nil || c.ctx == nil || binding == nil {
		return bytecode.NoopOperand
	}

	if binding.Storage == core.BindingStorageCell {
		return c.LoadBindingValue(binding)
	}

	snapshot := c.ctx.Function.Registers.Allocate()
	c.facts.EmitMoveAuto(snapshot, binding.Register)

	return snapshot
}

func (c *BindingCompiler) storeBindingValue(binding *core.Variable, src bytecode.Operand, publishedType core.ValueType) bytecode.Operand {
	if c == nil || c.ctx == nil || binding == nil {
		return bytecode.NoopOperand
	}

	if binding.Storage == core.BindingStorageCell {
		src = ensureOperandRegister(c.ctx, c.facts, src)
		c.ctx.Program.Emitter.EmitStoreCell(binding.Register, src)
		return binding.Register
	}

	if src.IsConstant() {
		c.ctx.Program.Emitter.EmitLoadConst(binding.Register, src)
	} else {
		c.facts.EmitMoveAuto(binding.Register, src)
	}

	c.ctx.Function.Types.Set(binding.Register, publishedType)

	return binding.Register
}

func (c *BindingCompiler) reportInvalidAssignmentTarget(ctx antlr.ParserRuleContext) {
	if ctx == nil {
		return
	}

	err := c.ctx.Program.Errors.Create(parserd.SyntaxError, ctx, "Assignment target must be a local variable name")
	err.Hint = "Property and index assignment are not supported. Use UPDATE for structural changes."
	c.ctx.Program.Errors.Add(err)
}

func assignmentOperatorText(ctx *fql.AssignmentStatementContext) string {
	if ctx == nil || ctx.AssignmentOperator() == nil {
		return ""
	}

	return ctx.AssignmentOperator().GetText()
}
