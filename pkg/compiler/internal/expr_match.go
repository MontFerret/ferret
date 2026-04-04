package internal

import (
	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/parser/fql"
)

func (c *ExprCompiler) compileMatchExpression(ctx fql.IMatchExpressionContext) bytecode.Operand {
	return c.matchCompiler.compileMatchExpression(ctx)
}
