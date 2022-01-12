// Code generated from antlr/FqlParser.g4 by ANTLR 4.9.3. DO NOT EDIT.

package fql // FqlParser
import "github.com/antlr/antlr4/runtime/Go/antlr"

// A complete Visitor for a parse tree produced by FqlParser.
type FqlParserVisitor interface {
	antlr.ParseTreeVisitor

	// Visit a parse tree produced by FqlParser#program.
	VisitProgram(ctx *ProgramContext) interface{}

	// Visit a parse tree produced by FqlParser#head.
	VisitHead(ctx *HeadContext) interface{}

	// Visit a parse tree produced by FqlParser#useExpression.
	VisitUseExpression(ctx *UseExpressionContext) interface{}

	// Visit a parse tree produced by FqlParser#use.
	VisitUse(ctx *UseContext) interface{}

	// Visit a parse tree produced by FqlParser#body.
	VisitBody(ctx *BodyContext) interface{}

	// Visit a parse tree produced by FqlParser#bodyStatement.
	VisitBodyStatement(ctx *BodyStatementContext) interface{}

	// Visit a parse tree produced by FqlParser#bodyExpression.
	VisitBodyExpression(ctx *BodyExpressionContext) interface{}

	// Visit a parse tree produced by FqlParser#variableDeclaration.
	VisitVariableDeclaration(ctx *VariableDeclarationContext) interface{}

	// Visit a parse tree produced by FqlParser#returnExpression.
	VisitReturnExpression(ctx *ReturnExpressionContext) interface{}

	// Visit a parse tree produced by FqlParser#forExpression.
	VisitForExpression(ctx *ForExpressionContext) interface{}

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

	// Visit a parse tree produced by FqlParser#limitClauseValue.
	VisitLimitClauseValue(ctx *LimitClauseValueContext) interface{}

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

	// Visit a parse tree produced by FqlParser#waitForExpression.
	VisitWaitForExpression(ctx *WaitForExpressionContext) interface{}

	// Visit a parse tree produced by FqlParser#waitForEventName.
	VisitWaitForEventName(ctx *WaitForEventNameContext) interface{}

	// Visit a parse tree produced by FqlParser#waitForEventSource.
	VisitWaitForEventSource(ctx *WaitForEventSourceContext) interface{}

	// Visit a parse tree produced by FqlParser#optionsClause.
	VisitOptionsClause(ctx *OptionsClauseContext) interface{}

	// Visit a parse tree produced by FqlParser#timeoutClause.
	VisitTimeoutClause(ctx *TimeoutClauseContext) interface{}

	// Visit a parse tree produced by FqlParser#param.
	VisitParam(ctx *ParamContext) interface{}

	// Visit a parse tree produced by FqlParser#variable.
	VisitVariable(ctx *VariableContext) interface{}

	// Visit a parse tree produced by FqlParser#literal.
	VisitLiteral(ctx *LiteralContext) interface{}

	// Visit a parse tree produced by FqlParser#arrayLiteral.
	VisitArrayLiteral(ctx *ArrayLiteralContext) interface{}

	// Visit a parse tree produced by FqlParser#objectLiteral.
	VisitObjectLiteral(ctx *ObjectLiteralContext) interface{}

	// Visit a parse tree produced by FqlParser#booleanLiteral.
	VisitBooleanLiteral(ctx *BooleanLiteralContext) interface{}

	// Visit a parse tree produced by FqlParser#stringLiteral.
	VisitStringLiteral(ctx *StringLiteralContext) interface{}

	// Visit a parse tree produced by FqlParser#floatLiteral.
	VisitFloatLiteral(ctx *FloatLiteralContext) interface{}

	// Visit a parse tree produced by FqlParser#integerLiteral.
	VisitIntegerLiteral(ctx *IntegerLiteralContext) interface{}

	// Visit a parse tree produced by FqlParser#noneLiteral.
	VisitNoneLiteral(ctx *NoneLiteralContext) interface{}

	// Visit a parse tree produced by FqlParser#propertyAssignment.
	VisitPropertyAssignment(ctx *PropertyAssignmentContext) interface{}

	// Visit a parse tree produced by FqlParser#computedPropertyName.
	VisitComputedPropertyName(ctx *ComputedPropertyNameContext) interface{}

	// Visit a parse tree produced by FqlParser#propertyName.
	VisitPropertyName(ctx *PropertyNameContext) interface{}

	// Visit a parse tree produced by FqlParser#namespaceIdentifier.
	VisitNamespaceIdentifier(ctx *NamespaceIdentifierContext) interface{}

	// Visit a parse tree produced by FqlParser#namespace.
	VisitNamespace(ctx *NamespaceContext) interface{}

	// Visit a parse tree produced by FqlParser#memberExpression.
	VisitMemberExpression(ctx *MemberExpressionContext) interface{}

	// Visit a parse tree produced by FqlParser#memberExpressionSource.
	VisitMemberExpressionSource(ctx *MemberExpressionSourceContext) interface{}

	// Visit a parse tree produced by FqlParser#functionCallExpression.
	VisitFunctionCallExpression(ctx *FunctionCallExpressionContext) interface{}

	// Visit a parse tree produced by FqlParser#functionCall.
	VisitFunctionCall(ctx *FunctionCallContext) interface{}

	// Visit a parse tree produced by FqlParser#functionName.
	VisitFunctionName(ctx *FunctionNameContext) interface{}

	// Visit a parse tree produced by FqlParser#argumentList.
	VisitArgumentList(ctx *ArgumentListContext) interface{}

	// Visit a parse tree produced by FqlParser#memberExpressionPath.
	VisitMemberExpressionPath(ctx *MemberExpressionPathContext) interface{}

	// Visit a parse tree produced by FqlParser#safeReservedWord.
	VisitSafeReservedWord(ctx *SafeReservedWordContext) interface{}

	// Visit a parse tree produced by FqlParser#unsafeReservedWord.
	VisitUnsafeReservedWord(ctx *UnsafeReservedWordContext) interface{}

	// Visit a parse tree produced by FqlParser#rangeOperator.
	VisitRangeOperator(ctx *RangeOperatorContext) interface{}

	// Visit a parse tree produced by FqlParser#rangeOperand.
	VisitRangeOperand(ctx *RangeOperandContext) interface{}

	// Visit a parse tree produced by FqlParser#expression.
	VisitExpression(ctx *ExpressionContext) interface{}

	// Visit a parse tree produced by FqlParser#predicate.
	VisitPredicate(ctx *PredicateContext) interface{}

	// Visit a parse tree produced by FqlParser#expressionAtom.
	VisitExpressionAtom(ctx *ExpressionAtomContext) interface{}

	// Visit a parse tree produced by FqlParser#arrayOperator.
	VisitArrayOperator(ctx *ArrayOperatorContext) interface{}

	// Visit a parse tree produced by FqlParser#equalityOperator.
	VisitEqualityOperator(ctx *EqualityOperatorContext) interface{}

	// Visit a parse tree produced by FqlParser#inOperator.
	VisitInOperator(ctx *InOperatorContext) interface{}

	// Visit a parse tree produced by FqlParser#likeOperator.
	VisitLikeOperator(ctx *LikeOperatorContext) interface{}

	// Visit a parse tree produced by FqlParser#unaryOperator.
	VisitUnaryOperator(ctx *UnaryOperatorContext) interface{}

	// Visit a parse tree produced by FqlParser#regexpOperator.
	VisitRegexpOperator(ctx *RegexpOperatorContext) interface{}

	// Visit a parse tree produced by FqlParser#logicalAndOperator.
	VisitLogicalAndOperator(ctx *LogicalAndOperatorContext) interface{}

	// Visit a parse tree produced by FqlParser#logicalOrOperator.
	VisitLogicalOrOperator(ctx *LogicalOrOperatorContext) interface{}

	// Visit a parse tree produced by FqlParser#multiplicativeOperator.
	VisitMultiplicativeOperator(ctx *MultiplicativeOperatorContext) interface{}

	// Visit a parse tree produced by FqlParser#additiveOperator.
	VisitAdditiveOperator(ctx *AdditiveOperatorContext) interface{}

	// Visit a parse tree produced by FqlParser#errorOperator.
	VisitErrorOperator(ctx *ErrorOperatorContext) interface{}
}
