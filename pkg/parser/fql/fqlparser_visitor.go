// Code generated from antlr/FqlParser.g4 by ANTLR 4.13.2. DO NOT EDIT.

package fql // FqlParser
import "github.com/antlr4-go/antlr/v4"

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

	// Visit a parse tree produced by FqlParser#eventFilterClause.
	VisitEventFilterClause(ctx *EventFilterClauseContext) interface{}

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

	// Visit a parse tree produced by FqlParser#collectAggregateSelector.
	VisitCollectAggregateSelector(ctx *CollectAggregateSelectorContext) interface{}

	// Visit a parse tree produced by FqlParser#collectAggregator.
	VisitCollectAggregator(ctx *CollectAggregatorContext) interface{}

	// Visit a parse tree produced by FqlParser#collectGroupProjection.
	VisitCollectGroupProjection(ctx *CollectGroupProjectionContext) interface{}

	// Visit a parse tree produced by FqlParser#collectGroupProjectionFilter.
	VisitCollectGroupProjectionFilter(ctx *CollectGroupProjectionFilterContext) interface{}

	// Visit a parse tree produced by FqlParser#collectCounter.
	VisitCollectCounter(ctx *CollectCounterContext) interface{}

	// Visit a parse tree produced by FqlParser#waitForExpression.
	VisitWaitForExpression(ctx *WaitForExpressionContext) interface{}

	// Visit a parse tree produced by FqlParser#dispatchExpression.
	VisitDispatchExpression(ctx *DispatchExpressionContext) interface{}

	// Visit a parse tree produced by FqlParser#dispatchEventName.
	VisitDispatchEventName(ctx *DispatchEventNameContext) interface{}

	// Visit a parse tree produced by FqlParser#dispatchTarget.
	VisitDispatchTarget(ctx *DispatchTargetContext) interface{}

	// Visit a parse tree produced by FqlParser#dispatchWithClause.
	VisitDispatchWithClause(ctx *DispatchWithClauseContext) interface{}

	// Visit a parse tree produced by FqlParser#dispatchOptionsClause.
	VisitDispatchOptionsClause(ctx *DispatchOptionsClauseContext) interface{}

	// Visit a parse tree produced by FqlParser#waitForEventExpression.
	VisitWaitForEventExpression(ctx *WaitForEventExpressionContext) interface{}

	// Visit a parse tree produced by FqlParser#waitForPredicateExpression.
	VisitWaitForPredicateExpression(ctx *WaitForPredicateExpressionContext) interface{}

	// Visit a parse tree produced by FqlParser#waitForPredicate.
	VisitWaitForPredicate(ctx *WaitForPredicateContext) interface{}

	// Visit a parse tree produced by FqlParser#waitForEventName.
	VisitWaitForEventName(ctx *WaitForEventNameContext) interface{}

	// Visit a parse tree produced by FqlParser#waitForEventSource.
	VisitWaitForEventSource(ctx *WaitForEventSourceContext) interface{}

	// Visit a parse tree produced by FqlParser#optionsClause.
	VisitOptionsClause(ctx *OptionsClauseContext) interface{}

	// Visit a parse tree produced by FqlParser#timeoutClause.
	VisitTimeoutClause(ctx *TimeoutClauseContext) interface{}

	// Visit a parse tree produced by FqlParser#everyClause.
	VisitEveryClause(ctx *EveryClauseContext) interface{}

	// Visit a parse tree produced by FqlParser#everyClauseValue.
	VisitEveryClauseValue(ctx *EveryClauseValueContext) interface{}

	// Visit a parse tree produced by FqlParser#backoffClause.
	VisitBackoffClause(ctx *BackoffClauseContext) interface{}

	// Visit a parse tree produced by FqlParser#jitterClause.
	VisitJitterClause(ctx *JitterClauseContext) interface{}

	// Visit a parse tree produced by FqlParser#jitterClauseValue.
	VisitJitterClauseValue(ctx *JitterClauseValueContext) interface{}

	// Visit a parse tree produced by FqlParser#backoffStrategy.
	VisitBackoffStrategy(ctx *BackoffStrategyContext) interface{}

	// Visit a parse tree produced by FqlParser#waitForOrThrowClause.
	VisitWaitForOrThrowClause(ctx *WaitForOrThrowClauseContext) interface{}

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

	// Visit a parse tree produced by FqlParser#templateLiteral.
	VisitTemplateLiteral(ctx *TemplateLiteralContext) interface{}

	// Visit a parse tree produced by FqlParser#templateElement.
	VisitTemplateElement(ctx *TemplateElementContext) interface{}

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

	// Visit a parse tree produced by FqlParser#arrayExpansion.
	VisitArrayExpansion(ctx *ArrayExpansionContext) interface{}

	// Visit a parse tree produced by FqlParser#arrayContraction.
	VisitArrayContraction(ctx *ArrayContractionContext) interface{}

	// Visit a parse tree produced by FqlParser#arrayQuestionMark.
	VisitArrayQuestionMark(ctx *ArrayQuestionMarkContext) interface{}

	// Visit a parse tree produced by FqlParser#arrayQuestionQuantifier.
	VisitArrayQuestionQuantifier(ctx *ArrayQuestionQuantifierContext) interface{}

	// Visit a parse tree produced by FqlParser#arrayQuestionQuantifierValue.
	VisitArrayQuestionQuantifierValue(ctx *ArrayQuestionQuantifierValueContext) interface{}

	// Visit a parse tree produced by FqlParser#arrayApply.
	VisitArrayApply(ctx *ArrayApplyContext) interface{}

	// Visit a parse tree produced by FqlParser#inlineExpression.
	VisitInlineExpression(ctx *InlineExpressionContext) interface{}

	// Visit a parse tree produced by FqlParser#inlineFilter.
	VisitInlineFilter(ctx *InlineFilterContext) interface{}

	// Visit a parse tree produced by FqlParser#inlineLimit.
	VisitInlineLimit(ctx *InlineLimitContext) interface{}

	// Visit a parse tree produced by FqlParser#inlineReturn.
	VisitInlineReturn(ctx *InlineReturnContext) interface{}

	// Visit a parse tree produced by FqlParser#safeReservedWord.
	VisitSafeReservedWord(ctx *SafeReservedWordContext) interface{}

	// Visit a parse tree produced by FqlParser#unsafeReservedWord.
	VisitUnsafeReservedWord(ctx *UnsafeReservedWordContext) interface{}

	// Visit a parse tree produced by FqlParser#durationLiteral.
	VisitDurationLiteral(ctx *DurationLiteralContext) interface{}

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

	// Visit a parse tree produced by FqlParser#implicitMemberExpression.
	VisitImplicitMemberExpression(ctx *ImplicitMemberExpressionContext) interface{}

	// Visit a parse tree produced by FqlParser#implicitCurrentExpression.
	VisitImplicitCurrentExpression(ctx *ImplicitCurrentExpressionContext) interface{}

	// Visit a parse tree produced by FqlParser#implicitMemberExpressionStart.
	VisitImplicitMemberExpressionStart(ctx *ImplicitMemberExpressionStartContext) interface{}

	// Visit a parse tree produced by FqlParser#queryLiteral.
	VisitQueryLiteral(ctx *QueryLiteralContext) interface{}

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
