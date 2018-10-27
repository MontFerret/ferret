// Code generated from antlr/FqlParser.g4 by ANTLR 4.7.1. DO NOT EDIT.

package fql // FqlParser
import "github.com/antlr/antlr4/runtime/Go/antlr"

// A complete Visitor for a parse tree produced by FqlParser.
type FqlParserVisitor interface {
	antlr.ParseTreeVisitor

	// Visit a parse tree produced by FqlParser#program.
	VisitProgram(ctx *ProgramContext) interface{}

	// Visit a parse tree produced by FqlParser#body.
	VisitBody(ctx *BodyContext) interface{}

	// Visit a parse tree produced by FqlParser#bodyStatement.
	VisitBodyStatement(ctx *BodyStatementContext) interface{}

	// Visit a parse tree produced by FqlParser#bodyExpression.
	VisitBodyExpression(ctx *BodyExpressionContext) interface{}

	// Visit a parse tree produced by FqlParser#returnExpression.
	VisitReturnExpression(ctx *ReturnExpressionContext) interface{}

	// Visit a parse tree produced by FqlParser#forExpression.
	VisitForExpression(ctx *ForExpressionContext) interface{}

	// Visit a parse tree produced by FqlParser#forExpressionValueVariable.
	VisitForExpressionValueVariable(ctx *ForExpressionValueVariableContext) interface{}

	// Visit a parse tree produced by FqlParser#forExpressionKeyVariable.
	VisitForExpressionKeyVariable(ctx *ForExpressionKeyVariableContext) interface{}

	// Visit a parse tree produced by FqlParser#forExpressionSource.
	VisitForExpressionSource(ctx *ForExpressionSourceContext) interface{}

	// Visit a parse tree produced by FqlParser#forExpressionClause.
	VisitForExpressionClause(ctx *ForExpressionClauseContext) interface{}

	// Visit a parse tree produced by FqlParser#forExpressionStatement.
	VisitForExpressionStatement(ctx *ForExpressionStatementContext) interface{}

	// Visit a parse tree produced by FqlParser#forExpressionBody.
	VisitForExpressionBody(ctx *ForExpressionBodyContext) interface{}

	// Visit a parse tree produced by FqlParser#forExpressionReturn.
	VisitForExpressionReturn(ctx *ForExpressionReturnContext) interface{}

	// Visit a parse tree produced by FqlParser#filterClause.
	VisitFilterClause(ctx *FilterClauseContext) interface{}

	// Visit a parse tree produced by FqlParser#limitClause.
	VisitLimitClause(ctx *LimitClauseContext) interface{}

	// Visit a parse tree produced by FqlParser#sortClause.
	VisitSortClause(ctx *SortClauseContext) interface{}

	// Visit a parse tree produced by FqlParser#sortClauseExpression.
	VisitSortClauseExpression(ctx *SortClauseExpressionContext) interface{}

	// Visit a parse tree produced by FqlParser#collectClause.
	VisitCollectClause(ctx *CollectClauseContext) interface{}

	// Visit a parse tree produced by FqlParser#collectSelector.
	VisitCollectSelector(ctx *CollectSelectorContext) interface{}

	// Visit a parse tree produced by FqlParser#collectGrouping.
	VisitCollectGrouping(ctx *CollectGroupingContext) interface{}

	// Visit a parse tree produced by FqlParser#collectAggregator.
	VisitCollectAggregator(ctx *CollectAggregatorContext) interface{}

	// Visit a parse tree produced by FqlParser#collectAggregateSelector.
	VisitCollectAggregateSelector(ctx *CollectAggregateSelectorContext) interface{}

	// Visit a parse tree produced by FqlParser#collectGroupVariable.
	VisitCollectGroupVariable(ctx *CollectGroupVariableContext) interface{}

	// Visit a parse tree produced by FqlParser#collectCounter.
	VisitCollectCounter(ctx *CollectCounterContext) interface{}

	// Visit a parse tree produced by FqlParser#variableDeclaration.
	VisitVariableDeclaration(ctx *VariableDeclarationContext) interface{}

	// Visit a parse tree produced by FqlParser#param.
	VisitParam(ctx *ParamContext) interface{}

	// Visit a parse tree produced by FqlParser#variable.
	VisitVariable(ctx *VariableContext) interface{}

	// Visit a parse tree produced by FqlParser#rangeOperator.
	VisitRangeOperator(ctx *RangeOperatorContext) interface{}

	// Visit a parse tree produced by FqlParser#arrayLiteral.
	VisitArrayLiteral(ctx *ArrayLiteralContext) interface{}

	// Visit a parse tree produced by FqlParser#objectLiteral.
	VisitObjectLiteral(ctx *ObjectLiteralContext) interface{}

	// Visit a parse tree produced by FqlParser#booleanLiteral.
	VisitBooleanLiteral(ctx *BooleanLiteralContext) interface{}

	// Visit a parse tree produced by FqlParser#stringLiteral.
	VisitStringLiteral(ctx *StringLiteralContext) interface{}

	// Visit a parse tree produced by FqlParser#integerLiteral.
	VisitIntegerLiteral(ctx *IntegerLiteralContext) interface{}

	// Visit a parse tree produced by FqlParser#floatLiteral.
	VisitFloatLiteral(ctx *FloatLiteralContext) interface{}

	// Visit a parse tree produced by FqlParser#noneLiteral.
	VisitNoneLiteral(ctx *NoneLiteralContext) interface{}

	// Visit a parse tree produced by FqlParser#arrayElementList.
	VisitArrayElementList(ctx *ArrayElementListContext) interface{}

	// Visit a parse tree produced by FqlParser#propertyAssignment.
	VisitPropertyAssignment(ctx *PropertyAssignmentContext) interface{}

	// Visit a parse tree produced by FqlParser#memberExpression.
	VisitMemberExpression(ctx *MemberExpressionContext) interface{}

	// Visit a parse tree produced by FqlParser#shorthandPropertyName.
	VisitShorthandPropertyName(ctx *ShorthandPropertyNameContext) interface{}

	// Visit a parse tree produced by FqlParser#computedPropertyName.
	VisitComputedPropertyName(ctx *ComputedPropertyNameContext) interface{}

	// Visit a parse tree produced by FqlParser#propertyName.
	VisitPropertyName(ctx *PropertyNameContext) interface{}

	// Visit a parse tree produced by FqlParser#expressionSequence.
	VisitExpressionSequence(ctx *ExpressionSequenceContext) interface{}

	// Visit a parse tree produced by FqlParser#functionCallExpression.
	VisitFunctionCallExpression(ctx *FunctionCallExpressionContext) interface{}

	// Visit a parse tree produced by FqlParser#arguments.
	VisitArguments(ctx *ArgumentsContext) interface{}

	// Visit a parse tree produced by FqlParser#expression.
	VisitExpression(ctx *ExpressionContext) interface{}

	// Visit a parse tree produced by FqlParser#forTernaryExpression.
	VisitForTernaryExpression(ctx *ForTernaryExpressionContext) interface{}

	// Visit a parse tree produced by FqlParser#arrayOperator.
	VisitArrayOperator(ctx *ArrayOperatorContext) interface{}

	// Visit a parse tree produced by FqlParser#inOperator.
	VisitInOperator(ctx *InOperatorContext) interface{}

	// Visit a parse tree produced by FqlParser#equalityOperator.
	VisitEqualityOperator(ctx *EqualityOperatorContext) interface{}

	// Visit a parse tree produced by FqlParser#logicalOperator.
	VisitLogicalOperator(ctx *LogicalOperatorContext) interface{}

	// Visit a parse tree produced by FqlParser#mathOperator.
	VisitMathOperator(ctx *MathOperatorContext) interface{}

	// Visit a parse tree produced by FqlParser#unaryOperator.
	VisitUnaryOperator(ctx *UnaryOperatorContext) interface{}
}
