// Code generated from antlr/FqlParser.g4 by ANTLR 4.8. DO NOT EDIT.

package fql // FqlParser
import "github.com/antlr/antlr4/runtime/Go/antlr"

type BaseFqlParserVisitor struct {
	*antlr.BaseParseTreeVisitor
}

func (v *BaseFqlParserVisitor) VisitProgram(ctx *ProgramContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseFqlParserVisitor) VisitHead(ctx *HeadContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseFqlParserVisitor) VisitUseExpression(ctx *UseExpressionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseFqlParserVisitor) VisitUse(ctx *UseContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseFqlParserVisitor) VisitBody(ctx *BodyContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseFqlParserVisitor) VisitBodyStatement(ctx *BodyStatementContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseFqlParserVisitor) VisitBodyExpression(ctx *BodyExpressionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseFqlParserVisitor) VisitReturnExpression(ctx *ReturnExpressionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseFqlParserVisitor) VisitForExpression(ctx *ForExpressionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseFqlParserVisitor) VisitForExpressionValueVariable(ctx *ForExpressionValueVariableContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseFqlParserVisitor) VisitForExpressionKeyVariable(ctx *ForExpressionKeyVariableContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseFqlParserVisitor) VisitForExpressionSource(ctx *ForExpressionSourceContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseFqlParserVisitor) VisitForExpressionClause(ctx *ForExpressionClauseContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseFqlParserVisitor) VisitForExpressionStatement(ctx *ForExpressionStatementContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseFqlParserVisitor) VisitForExpressionBody(ctx *ForExpressionBodyContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseFqlParserVisitor) VisitForExpressionReturn(ctx *ForExpressionReturnContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseFqlParserVisitor) VisitFilterClause(ctx *FilterClauseContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseFqlParserVisitor) VisitLimitClause(ctx *LimitClauseContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseFqlParserVisitor) VisitLimitClauseValue(ctx *LimitClauseValueContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseFqlParserVisitor) VisitSortClause(ctx *SortClauseContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseFqlParserVisitor) VisitSortClauseExpression(ctx *SortClauseExpressionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseFqlParserVisitor) VisitCollectClause(ctx *CollectClauseContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseFqlParserVisitor) VisitCollectSelector(ctx *CollectSelectorContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseFqlParserVisitor) VisitCollectGrouping(ctx *CollectGroupingContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseFqlParserVisitor) VisitCollectAggregator(ctx *CollectAggregatorContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseFqlParserVisitor) VisitCollectAggregateSelector(ctx *CollectAggregateSelectorContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseFqlParserVisitor) VisitCollectGroupVariable(ctx *CollectGroupVariableContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseFqlParserVisitor) VisitCollectCounter(ctx *CollectCounterContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseFqlParserVisitor) VisitVariableDeclaration(ctx *VariableDeclarationContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseFqlParserVisitor) VisitParam(ctx *ParamContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseFqlParserVisitor) VisitVariable(ctx *VariableContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseFqlParserVisitor) VisitRangeOperator(ctx *RangeOperatorContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseFqlParserVisitor) VisitArrayLiteral(ctx *ArrayLiteralContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseFqlParserVisitor) VisitObjectLiteral(ctx *ObjectLiteralContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseFqlParserVisitor) VisitBooleanLiteral(ctx *BooleanLiteralContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseFqlParserVisitor) VisitStringLiteral(ctx *StringLiteralContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseFqlParserVisitor) VisitIntegerLiteral(ctx *IntegerLiteralContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseFqlParserVisitor) VisitFloatLiteral(ctx *FloatLiteralContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseFqlParserVisitor) VisitNoneLiteral(ctx *NoneLiteralContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseFqlParserVisitor) VisitArrayElementList(ctx *ArrayElementListContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseFqlParserVisitor) VisitPropertyAssignment(ctx *PropertyAssignmentContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseFqlParserVisitor) VisitShorthandPropertyName(ctx *ShorthandPropertyNameContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseFqlParserVisitor) VisitComputedPropertyName(ctx *ComputedPropertyNameContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseFqlParserVisitor) VisitPropertyName(ctx *PropertyNameContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseFqlParserVisitor) VisitExpressionGroup(ctx *ExpressionGroupContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseFqlParserVisitor) VisitNamespaceIdentifier(ctx *NamespaceIdentifierContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseFqlParserVisitor) VisitNamespace(ctx *NamespaceContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseFqlParserVisitor) VisitFunctionIdentifier(ctx *FunctionIdentifierContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseFqlParserVisitor) VisitFunctionCallExpression(ctx *FunctionCallExpressionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseFqlParserVisitor) VisitMember(ctx *MemberContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseFqlParserVisitor) VisitMemberPath(ctx *MemberPathContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseFqlParserVisitor) VisitMemberExpression(ctx *MemberExpressionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseFqlParserVisitor) VisitArguments(ctx *ArgumentsContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseFqlParserVisitor) VisitExpression(ctx *ExpressionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseFqlParserVisitor) VisitForTernaryExpression(ctx *ForTernaryExpressionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseFqlParserVisitor) VisitArrayOperator(ctx *ArrayOperatorContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseFqlParserVisitor) VisitInOperator(ctx *InOperatorContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseFqlParserVisitor) VisitEqualityOperator(ctx *EqualityOperatorContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseFqlParserVisitor) VisitRegexpOperator(ctx *RegexpOperatorContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseFqlParserVisitor) VisitLogicalAndOperator(ctx *LogicalAndOperatorContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseFqlParserVisitor) VisitLogicalOrOperator(ctx *LogicalOrOperatorContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseFqlParserVisitor) VisitMultiplicativeOperator(ctx *MultiplicativeOperatorContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseFqlParserVisitor) VisitAdditiveOperator(ctx *AdditiveOperatorContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseFqlParserVisitor) VisitUnaryOperator(ctx *UnaryOperatorContext) interface{} {
	return v.VisitChildren(ctx)
}
