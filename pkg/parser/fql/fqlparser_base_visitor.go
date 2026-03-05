// Code generated from antlr/FqlParser.g4 by ANTLR 4.13.2. DO NOT EDIT.

package fql // FqlParser
import "github.com/antlr4-go/antlr/v4"

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

func (v *BaseFqlParserVisitor) VisitVariableDeclaration(ctx *VariableDeclarationContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseFqlParserVisitor) VisitFunctionDeclaration(ctx *FunctionDeclarationContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseFqlParserVisitor) VisitFunctionParameterList(ctx *FunctionParameterListContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseFqlParserVisitor) VisitFunctionParameter(ctx *FunctionParameterContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseFqlParserVisitor) VisitFunctionBody(ctx *FunctionBodyContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseFqlParserVisitor) VisitFunctionArrow(ctx *FunctionArrowContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseFqlParserVisitor) VisitFunctionBlock(ctx *FunctionBlockContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseFqlParserVisitor) VisitFunctionStatement(ctx *FunctionStatementContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseFqlParserVisitor) VisitExpressionStatement(ctx *ExpressionStatementContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseFqlParserVisitor) VisitFunctionReturn(ctx *FunctionReturnContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseFqlParserVisitor) VisitReturnExpression(ctx *ReturnExpressionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseFqlParserVisitor) VisitForExpression(ctx *ForExpressionContext) interface{} {
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

func (v *BaseFqlParserVisitor) VisitEventFilterClause(ctx *EventFilterClauseContext) interface{} {
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

func (v *BaseFqlParserVisitor) VisitBindingIdentifier(ctx *BindingIdentifierContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseFqlParserVisitor) VisitLoopVariable(ctx *LoopVariableContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseFqlParserVisitor) VisitCollectSelector(ctx *CollectSelectorContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseFqlParserVisitor) VisitCollectGrouping(ctx *CollectGroupingContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseFqlParserVisitor) VisitCollectAggregateSelector(ctx *CollectAggregateSelectorContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseFqlParserVisitor) VisitCollectAggregator(ctx *CollectAggregatorContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseFqlParserVisitor) VisitCollectGroupProjection(ctx *CollectGroupProjectionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseFqlParserVisitor) VisitCollectGroupProjectionFilter(ctx *CollectGroupProjectionFilterContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseFqlParserVisitor) VisitCollectCounter(ctx *CollectCounterContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseFqlParserVisitor) VisitWaitForExpression(ctx *WaitForExpressionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseFqlParserVisitor) VisitDispatchExpression(ctx *DispatchExpressionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseFqlParserVisitor) VisitDispatchEventName(ctx *DispatchEventNameContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseFqlParserVisitor) VisitDispatchTarget(ctx *DispatchTargetContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseFqlParserVisitor) VisitDispatchWithClause(ctx *DispatchWithClauseContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseFqlParserVisitor) VisitDispatchOptionsClause(ctx *DispatchOptionsClauseContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseFqlParserVisitor) VisitWaitForEventExpression(ctx *WaitForEventExpressionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseFqlParserVisitor) VisitWaitForPredicateExpression(ctx *WaitForPredicateExpressionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseFqlParserVisitor) VisitWaitForPredicate(ctx *WaitForPredicateContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseFqlParserVisitor) VisitWaitForEventName(ctx *WaitForEventNameContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseFqlParserVisitor) VisitWaitForEventSource(ctx *WaitForEventSourceContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseFqlParserVisitor) VisitOptionsClause(ctx *OptionsClauseContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseFqlParserVisitor) VisitTimeoutClause(ctx *TimeoutClauseContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseFqlParserVisitor) VisitEveryClause(ctx *EveryClauseContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseFqlParserVisitor) VisitEveryClauseValue(ctx *EveryClauseValueContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseFqlParserVisitor) VisitBackoffClause(ctx *BackoffClauseContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseFqlParserVisitor) VisitJitterClause(ctx *JitterClauseContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseFqlParserVisitor) VisitJitterClauseValue(ctx *JitterClauseValueContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseFqlParserVisitor) VisitBackoffStrategy(ctx *BackoffStrategyContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseFqlParserVisitor) VisitWaitForOrThrowClause(ctx *WaitForOrThrowClauseContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseFqlParserVisitor) VisitParam(ctx *ParamContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseFqlParserVisitor) VisitVariable(ctx *VariableContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseFqlParserVisitor) VisitLiteral(ctx *LiteralContext) interface{} {
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

func (v *BaseFqlParserVisitor) VisitTemplateLiteral(ctx *TemplateLiteralContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseFqlParserVisitor) VisitTemplateElement(ctx *TemplateElementContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseFqlParserVisitor) VisitFloatLiteral(ctx *FloatLiteralContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseFqlParserVisitor) VisitIntegerLiteral(ctx *IntegerLiteralContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseFqlParserVisitor) VisitNoneLiteral(ctx *NoneLiteralContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseFqlParserVisitor) VisitPropertyAssignment(ctx *PropertyAssignmentContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseFqlParserVisitor) VisitComputedPropertyName(ctx *ComputedPropertyNameContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseFqlParserVisitor) VisitPropertyName(ctx *PropertyNameContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseFqlParserVisitor) VisitNamespaceIdentifier(ctx *NamespaceIdentifierContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseFqlParserVisitor) VisitNamespace(ctx *NamespaceContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseFqlParserVisitor) VisitMemberExpression(ctx *MemberExpressionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseFqlParserVisitor) VisitMemberExpressionSource(ctx *MemberExpressionSourceContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseFqlParserVisitor) VisitFunctionCallExpression(ctx *FunctionCallExpressionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseFqlParserVisitor) VisitFunctionCall(ctx *FunctionCallContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseFqlParserVisitor) VisitFunctionName(ctx *FunctionNameContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseFqlParserVisitor) VisitArgumentList(ctx *ArgumentListContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseFqlParserVisitor) VisitMemberExpressionPath(ctx *MemberExpressionPathContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseFqlParserVisitor) VisitArrayExpansion(ctx *ArrayExpansionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseFqlParserVisitor) VisitArrayContraction(ctx *ArrayContractionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseFqlParserVisitor) VisitArrayQuestionMark(ctx *ArrayQuestionMarkContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseFqlParserVisitor) VisitArrayQuestionQuantifier(ctx *ArrayQuestionQuantifierContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseFqlParserVisitor) VisitArrayQuestionQuantifierValue(ctx *ArrayQuestionQuantifierValueContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseFqlParserVisitor) VisitArrayApply(ctx *ArrayApplyContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseFqlParserVisitor) VisitInlineExpression(ctx *InlineExpressionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseFqlParserVisitor) VisitInlineFilter(ctx *InlineFilterContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseFqlParserVisitor) VisitInlineLimit(ctx *InlineLimitContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseFqlParserVisitor) VisitInlineReturn(ctx *InlineReturnContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseFqlParserVisitor) VisitSafeReservedWord(ctx *SafeReservedWordContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseFqlParserVisitor) VisitUnsafeReservedWord(ctx *UnsafeReservedWordContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseFqlParserVisitor) VisitDurationLiteral(ctx *DurationLiteralContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseFqlParserVisitor) VisitRangeOperator(ctx *RangeOperatorContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseFqlParserVisitor) VisitRangeOperand(ctx *RangeOperandContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseFqlParserVisitor) VisitExpression(ctx *ExpressionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseFqlParserVisitor) VisitPredicate(ctx *PredicateContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseFqlParserVisitor) VisitExpressionAtom(ctx *ExpressionAtomContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseFqlParserVisitor) VisitMatchExpression(ctx *MatchExpressionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseFqlParserVisitor) VisitMatchPatternArms(ctx *MatchPatternArmsContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseFqlParserVisitor) VisitMatchPatternArmList(ctx *MatchPatternArmListContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseFqlParserVisitor) VisitMatchGuardArms(ctx *MatchGuardArmsContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseFqlParserVisitor) VisitMatchGuardArmList(ctx *MatchGuardArmListContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseFqlParserVisitor) VisitMatchPatternArm(ctx *MatchPatternArmContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseFqlParserVisitor) VisitMatchPatternGuard(ctx *MatchPatternGuardContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseFqlParserVisitor) VisitMatchGuardArm(ctx *MatchGuardArmContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseFqlParserVisitor) VisitMatchDefaultArm(ctx *MatchDefaultArmContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseFqlParserVisitor) VisitMatchPattern(ctx *MatchPatternContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseFqlParserVisitor) VisitMatchBindingPattern(ctx *MatchBindingPatternContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseFqlParserVisitor) VisitMatchLiteralPattern(ctx *MatchLiteralPatternContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseFqlParserVisitor) VisitMatchObjectPattern(ctx *MatchObjectPatternContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseFqlParserVisitor) VisitMatchObjectPatternProperty(ctx *MatchObjectPatternPropertyContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseFqlParserVisitor) VisitMatchObjectPatternKey(ctx *MatchObjectPatternKeyContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseFqlParserVisitor) VisitImplicitMemberExpression(ctx *ImplicitMemberExpressionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseFqlParserVisitor) VisitImplicitCurrentExpression(ctx *ImplicitCurrentExpressionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseFqlParserVisitor) VisitImplicitMemberExpressionStart(ctx *ImplicitMemberExpressionStartContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseFqlParserVisitor) VisitQueryExpression(ctx *QueryExpressionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseFqlParserVisitor) VisitQueryModifier(ctx *QueryModifierContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseFqlParserVisitor) VisitQueryPayload(ctx *QueryPayloadContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseFqlParserVisitor) VisitQueryWithOpt(ctx *QueryWithOptContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseFqlParserVisitor) VisitQueryLiteral(ctx *QueryLiteralContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseFqlParserVisitor) VisitArrayOperator(ctx *ArrayOperatorContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseFqlParserVisitor) VisitEqualityOperator(ctx *EqualityOperatorContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseFqlParserVisitor) VisitInOperator(ctx *InOperatorContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseFqlParserVisitor) VisitLikeOperator(ctx *LikeOperatorContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseFqlParserVisitor) VisitUnaryOperator(ctx *UnaryOperatorContext) interface{} {
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

func (v *BaseFqlParserVisitor) VisitErrorOperator(ctx *ErrorOperatorContext) interface{} {
	return v.VisitChildren(ctx)
}
