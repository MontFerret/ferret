package internal

import (
	"github.com/antlr4-go/antlr/v4"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/compiler/internal/core"
	"github.com/MontFerret/ferret/v2/pkg/parser/fql"
)

func isMutableDeclaration(ctx fql.IVariableDeclarationContext) bool {
	if ctx == nil {
		return false
	}

	decl, ok := ctx.(*fql.VariableDeclarationContext)

	return ok && decl.Var() != nil
}

func declarationStorage(ctx *CompilerContext, decl antlr.ParserRuleContext, mutable bool) core.BindingStorage {
	if ctx == nil || decl == nil || !mutable {
		return core.BindingStorageValue
	}

	if _, ok := ctx.PromotedBindings[decl]; ok {
		return core.BindingStorageCell
	}

	return core.BindingStorageValue
}

func loadBindingValue(ctx *CompilerContext, binding *core.Variable) bytecode.Operand {
	if ctx == nil || binding == nil {
		return bytecode.NoopOperand
	}

	if binding.Storage != core.BindingStorageCell {
		return binding.Register
	}

	dst := ctx.Registers.Allocate()
	ctx.Emitter.EmitLoadCell(dst, binding.Register)
	ctx.Types.Set(dst, binding.Type)

	return dst
}

func variableDeclarationName(ctx fql.IVariableDeclarationContext) string {
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
