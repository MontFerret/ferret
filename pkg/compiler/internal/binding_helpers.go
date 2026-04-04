package internal

import "github.com/MontFerret/ferret/v2/pkg/parser/fql"

func assignmentOperatorText(ctx *fql.AssignmentStatementContext) string {
	if ctx == nil || ctx.AssignmentOperator() == nil {
		return ""
	}

	return ctx.AssignmentOperator().GetText()
}
