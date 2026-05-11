package internal

import (
	"github.com/MontFerret/ferret/v2/pkg/compiler/internal/core"
	"github.com/MontFerret/ferret/v2/pkg/parser/fql"
)

func bindingDeclarationName(ctx fql.IVariableDeclarationContext) string {
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

func assignmentOperatorText(ctx *fql.AssignmentStatementContext) string {
	if ctx == nil || ctx.AssignmentOperator() == nil {
		return ""
	}

	return ctx.AssignmentOperator().GetText()
}

func augmentedAssignmentKnownTypeAllowed(operator string, typ core.ValueType) bool {
	if typ == core.TypeUnknown || typ == core.TypeAny {
		return true
	}

	switch operator {
	case "+=":
		return typ == core.TypeInt || typ == core.TypeFloat || typ == core.TypeString
	case "-=", "*=", "/=":
		return typ == core.TypeInt || typ == core.TypeFloat
	default:
		return true
	}
}
