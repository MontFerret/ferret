// Code generated from antlr/FqlParser.g4 by ANTLR 4.13.2. DO NOT EDIT.

package fql // FqlParser
import "github.com/antlr4-go/antlr/v4"

// BaseFqlParserListener is a complete listener for a parse tree produced by FqlParser.
type BaseFqlParserListener struct{}

var _ FqlParserListener = &BaseFqlParserListener{}

// VisitTerminal is called when a terminal node is visited.
func (s *BaseFqlParserListener) VisitTerminal(node antlr.TerminalNode) {}

// VisitErrorNode is called when an error node is visited.
func (s *BaseFqlParserListener) VisitErrorNode(node antlr.ErrorNode) {}

// EnterEveryRule is called when any rule is entered.
func (s *BaseFqlParserListener) EnterEveryRule(ctx antlr.ParserRuleContext) {}

// ExitEveryRule is called when any rule is exited.
func (s *BaseFqlParserListener) ExitEveryRule(ctx antlr.ParserRuleContext) {}

// EnterProgram is called when production program is entered.
func (s *BaseFqlParserListener) EnterProgram(ctx *ProgramContext) {}

// ExitProgram is called when production program is exited.
func (s *BaseFqlParserListener) ExitProgram(ctx *ProgramContext) {}

// EnterHead is called when production head is entered.
func (s *BaseFqlParserListener) EnterHead(ctx *HeadContext) {}

// ExitHead is called when production head is exited.
func (s *BaseFqlParserListener) ExitHead(ctx *HeadContext) {}

// EnterUseExpression is called when production useExpression is entered.
func (s *BaseFqlParserListener) EnterUseExpression(ctx *UseExpressionContext) {}

// ExitUseExpression is called when production useExpression is exited.
func (s *BaseFqlParserListener) ExitUseExpression(ctx *UseExpressionContext) {}

// EnterUse is called when production use is entered.
func (s *BaseFqlParserListener) EnterUse(ctx *UseContext) {}

// ExitUse is called when production use is exited.
func (s *BaseFqlParserListener) ExitUse(ctx *UseContext) {}

// EnterBody is called when production body is entered.
func (s *BaseFqlParserListener) EnterBody(ctx *BodyContext) {}

// ExitBody is called when production body is exited.
func (s *BaseFqlParserListener) ExitBody(ctx *BodyContext) {}

// EnterBodyStatement is called when production bodyStatement is entered.
func (s *BaseFqlParserListener) EnterBodyStatement(ctx *BodyStatementContext) {}

// ExitBodyStatement is called when production bodyStatement is exited.
func (s *BaseFqlParserListener) ExitBodyStatement(ctx *BodyStatementContext) {}

// EnterBodyExpression is called when production bodyExpression is entered.
func (s *BaseFqlParserListener) EnterBodyExpression(ctx *BodyExpressionContext) {}

// ExitBodyExpression is called when production bodyExpression is exited.
func (s *BaseFqlParserListener) ExitBodyExpression(ctx *BodyExpressionContext) {}

// EnterVariableDeclaration is called when production variableDeclaration is entered.
func (s *BaseFqlParserListener) EnterVariableDeclaration(ctx *VariableDeclarationContext) {}

// ExitVariableDeclaration is called when production variableDeclaration is exited.
func (s *BaseFqlParserListener) ExitVariableDeclaration(ctx *VariableDeclarationContext) {}

// EnterFunctionDeclaration is called when production functionDeclaration is entered.
func (s *BaseFqlParserListener) EnterFunctionDeclaration(ctx *FunctionDeclarationContext) {}

// ExitFunctionDeclaration is called when production functionDeclaration is exited.
func (s *BaseFqlParserListener) ExitFunctionDeclaration(ctx *FunctionDeclarationContext) {}

// EnterFunctionParameterList is called when production functionParameterList is entered.
func (s *BaseFqlParserListener) EnterFunctionParameterList(ctx *FunctionParameterListContext) {}

// ExitFunctionParameterList is called when production functionParameterList is exited.
func (s *BaseFqlParserListener) ExitFunctionParameterList(ctx *FunctionParameterListContext) {}

// EnterFunctionParameter is called when production functionParameter is entered.
func (s *BaseFqlParserListener) EnterFunctionParameter(ctx *FunctionParameterContext) {}

// ExitFunctionParameter is called when production functionParameter is exited.
func (s *BaseFqlParserListener) ExitFunctionParameter(ctx *FunctionParameterContext) {}

// EnterFunctionBody is called when production functionBody is entered.
func (s *BaseFqlParserListener) EnterFunctionBody(ctx *FunctionBodyContext) {}

// ExitFunctionBody is called when production functionBody is exited.
func (s *BaseFqlParserListener) ExitFunctionBody(ctx *FunctionBodyContext) {}

// EnterFunctionArrow is called when production functionArrow is entered.
func (s *BaseFqlParserListener) EnterFunctionArrow(ctx *FunctionArrowContext) {}

// ExitFunctionArrow is called when production functionArrow is exited.
func (s *BaseFqlParserListener) ExitFunctionArrow(ctx *FunctionArrowContext) {}

// EnterFunctionBlock is called when production functionBlock is entered.
func (s *BaseFqlParserListener) EnterFunctionBlock(ctx *FunctionBlockContext) {}

// ExitFunctionBlock is called when production functionBlock is exited.
func (s *BaseFqlParserListener) ExitFunctionBlock(ctx *FunctionBlockContext) {}

// EnterFunctionStatement is called when production functionStatement is entered.
func (s *BaseFqlParserListener) EnterFunctionStatement(ctx *FunctionStatementContext) {}

// ExitFunctionStatement is called when production functionStatement is exited.
func (s *BaseFqlParserListener) ExitFunctionStatement(ctx *FunctionStatementContext) {}

// EnterExpressionStatement is called when production expressionStatement is entered.
func (s *BaseFqlParserListener) EnterExpressionStatement(ctx *ExpressionStatementContext) {}

// ExitExpressionStatement is called when production expressionStatement is exited.
func (s *BaseFqlParserListener) ExitExpressionStatement(ctx *ExpressionStatementContext) {}

// EnterFunctionReturn is called when production functionReturn is entered.
func (s *BaseFqlParserListener) EnterFunctionReturn(ctx *FunctionReturnContext) {}

// ExitFunctionReturn is called when production functionReturn is exited.
func (s *BaseFqlParserListener) ExitFunctionReturn(ctx *FunctionReturnContext) {}

// EnterReturnExpression is called when production returnExpression is entered.
func (s *BaseFqlParserListener) EnterReturnExpression(ctx *ReturnExpressionContext) {}

// ExitReturnExpression is called when production returnExpression is exited.
func (s *BaseFqlParserListener) ExitReturnExpression(ctx *ReturnExpressionContext) {}

// EnterForExpression is called when production forExpression is entered.
func (s *BaseFqlParserListener) EnterForExpression(ctx *ForExpressionContext) {}

// ExitForExpression is called when production forExpression is exited.
func (s *BaseFqlParserListener) ExitForExpression(ctx *ForExpressionContext) {}

// EnterForExpressionSource is called when production forExpressionSource is entered.
func (s *BaseFqlParserListener) EnterForExpressionSource(ctx *ForExpressionSourceContext) {}

// ExitForExpressionSource is called when production forExpressionSource is exited.
func (s *BaseFqlParserListener) ExitForExpressionSource(ctx *ForExpressionSourceContext) {}

// EnterForExpressionClause is called when production forExpressionClause is entered.
func (s *BaseFqlParserListener) EnterForExpressionClause(ctx *ForExpressionClauseContext) {}

// ExitForExpressionClause is called when production forExpressionClause is exited.
func (s *BaseFqlParserListener) ExitForExpressionClause(ctx *ForExpressionClauseContext) {}

// EnterForExpressionStatement is called when production forExpressionStatement is entered.
func (s *BaseFqlParserListener) EnterForExpressionStatement(ctx *ForExpressionStatementContext) {}

// ExitForExpressionStatement is called when production forExpressionStatement is exited.
func (s *BaseFqlParserListener) ExitForExpressionStatement(ctx *ForExpressionStatementContext) {}

// EnterForExpressionBody is called when production forExpressionBody is entered.
func (s *BaseFqlParserListener) EnterForExpressionBody(ctx *ForExpressionBodyContext) {}

// ExitForExpressionBody is called when production forExpressionBody is exited.
func (s *BaseFqlParserListener) ExitForExpressionBody(ctx *ForExpressionBodyContext) {}

// EnterForExpressionReturn is called when production forExpressionReturn is entered.
func (s *BaseFqlParserListener) EnterForExpressionReturn(ctx *ForExpressionReturnContext) {}

// ExitForExpressionReturn is called when production forExpressionReturn is exited.
func (s *BaseFqlParserListener) ExitForExpressionReturn(ctx *ForExpressionReturnContext) {}

// EnterFilterClause is called when production filterClause is entered.
func (s *BaseFqlParserListener) EnterFilterClause(ctx *FilterClauseContext) {}

// ExitFilterClause is called when production filterClause is exited.
func (s *BaseFqlParserListener) ExitFilterClause(ctx *FilterClauseContext) {}

// EnterEventFilterClause is called when production eventFilterClause is entered.
func (s *BaseFqlParserListener) EnterEventFilterClause(ctx *EventFilterClauseContext) {}

// ExitEventFilterClause is called when production eventFilterClause is exited.
func (s *BaseFqlParserListener) ExitEventFilterClause(ctx *EventFilterClauseContext) {}

// EnterLimitClause is called when production limitClause is entered.
func (s *BaseFqlParserListener) EnterLimitClause(ctx *LimitClauseContext) {}

// ExitLimitClause is called when production limitClause is exited.
func (s *BaseFqlParserListener) ExitLimitClause(ctx *LimitClauseContext) {}

// EnterLimitClauseValue is called when production limitClauseValue is entered.
func (s *BaseFqlParserListener) EnterLimitClauseValue(ctx *LimitClauseValueContext) {}

// ExitLimitClauseValue is called when production limitClauseValue is exited.
func (s *BaseFqlParserListener) ExitLimitClauseValue(ctx *LimitClauseValueContext) {}

// EnterSortClause is called when production sortClause is entered.
func (s *BaseFqlParserListener) EnterSortClause(ctx *SortClauseContext) {}

// ExitSortClause is called when production sortClause is exited.
func (s *BaseFqlParserListener) ExitSortClause(ctx *SortClauseContext) {}

// EnterSortClauseExpression is called when production sortClauseExpression is entered.
func (s *BaseFqlParserListener) EnterSortClauseExpression(ctx *SortClauseExpressionContext) {}

// ExitSortClauseExpression is called when production sortClauseExpression is exited.
func (s *BaseFqlParserListener) ExitSortClauseExpression(ctx *SortClauseExpressionContext) {}

// EnterCollectClause is called when production collectClause is entered.
func (s *BaseFqlParserListener) EnterCollectClause(ctx *CollectClauseContext) {}

// ExitCollectClause is called when production collectClause is exited.
func (s *BaseFqlParserListener) ExitCollectClause(ctx *CollectClauseContext) {}

// EnterBindingIdentifier is called when production bindingIdentifier is entered.
func (s *BaseFqlParserListener) EnterBindingIdentifier(ctx *BindingIdentifierContext) {}

// ExitBindingIdentifier is called when production bindingIdentifier is exited.
func (s *BaseFqlParserListener) ExitBindingIdentifier(ctx *BindingIdentifierContext) {}

// EnterLoopVariable is called when production loopVariable is entered.
func (s *BaseFqlParserListener) EnterLoopVariable(ctx *LoopVariableContext) {}

// ExitLoopVariable is called when production loopVariable is exited.
func (s *BaseFqlParserListener) ExitLoopVariable(ctx *LoopVariableContext) {}

// EnterCollectSelector is called when production collectSelector is entered.
func (s *BaseFqlParserListener) EnterCollectSelector(ctx *CollectSelectorContext) {}

// ExitCollectSelector is called when production collectSelector is exited.
func (s *BaseFqlParserListener) ExitCollectSelector(ctx *CollectSelectorContext) {}

// EnterCollectGrouping is called when production collectGrouping is entered.
func (s *BaseFqlParserListener) EnterCollectGrouping(ctx *CollectGroupingContext) {}

// ExitCollectGrouping is called when production collectGrouping is exited.
func (s *BaseFqlParserListener) ExitCollectGrouping(ctx *CollectGroupingContext) {}

// EnterCollectAggregateSelector is called when production collectAggregateSelector is entered.
func (s *BaseFqlParserListener) EnterCollectAggregateSelector(ctx *CollectAggregateSelectorContext) {}

// ExitCollectAggregateSelector is called when production collectAggregateSelector is exited.
func (s *BaseFqlParserListener) ExitCollectAggregateSelector(ctx *CollectAggregateSelectorContext) {}

// EnterCollectAggregator is called when production collectAggregator is entered.
func (s *BaseFqlParserListener) EnterCollectAggregator(ctx *CollectAggregatorContext) {}

// ExitCollectAggregator is called when production collectAggregator is exited.
func (s *BaseFqlParserListener) ExitCollectAggregator(ctx *CollectAggregatorContext) {}

// EnterCollectGroupProjection is called when production collectGroupProjection is entered.
func (s *BaseFqlParserListener) EnterCollectGroupProjection(ctx *CollectGroupProjectionContext) {}

// ExitCollectGroupProjection is called when production collectGroupProjection is exited.
func (s *BaseFqlParserListener) ExitCollectGroupProjection(ctx *CollectGroupProjectionContext) {}

// EnterCollectGroupProjectionFilter is called when production collectGroupProjectionFilter is entered.
func (s *BaseFqlParserListener) EnterCollectGroupProjectionFilter(ctx *CollectGroupProjectionFilterContext) {
}

// ExitCollectGroupProjectionFilter is called when production collectGroupProjectionFilter is exited.
func (s *BaseFqlParserListener) ExitCollectGroupProjectionFilter(ctx *CollectGroupProjectionFilterContext) {
}

// EnterCollectCounter is called when production collectCounter is entered.
func (s *BaseFqlParserListener) EnterCollectCounter(ctx *CollectCounterContext) {}

// ExitCollectCounter is called when production collectCounter is exited.
func (s *BaseFqlParserListener) ExitCollectCounter(ctx *CollectCounterContext) {}

// EnterWaitForExpression is called when production waitForExpression is entered.
func (s *BaseFqlParserListener) EnterWaitForExpression(ctx *WaitForExpressionContext) {}

// ExitWaitForExpression is called when production waitForExpression is exited.
func (s *BaseFqlParserListener) ExitWaitForExpression(ctx *WaitForExpressionContext) {}

// EnterDispatchExpression is called when production dispatchExpression is entered.
func (s *BaseFqlParserListener) EnterDispatchExpression(ctx *DispatchExpressionContext) {}

// ExitDispatchExpression is called when production dispatchExpression is exited.
func (s *BaseFqlParserListener) ExitDispatchExpression(ctx *DispatchExpressionContext) {}

// EnterDispatchEventName is called when production dispatchEventName is entered.
func (s *BaseFqlParserListener) EnterDispatchEventName(ctx *DispatchEventNameContext) {}

// ExitDispatchEventName is called when production dispatchEventName is exited.
func (s *BaseFqlParserListener) ExitDispatchEventName(ctx *DispatchEventNameContext) {}

// EnterDispatchTarget is called when production dispatchTarget is entered.
func (s *BaseFqlParserListener) EnterDispatchTarget(ctx *DispatchTargetContext) {}

// ExitDispatchTarget is called when production dispatchTarget is exited.
func (s *BaseFqlParserListener) ExitDispatchTarget(ctx *DispatchTargetContext) {}

// EnterDispatchWithClause is called when production dispatchWithClause is entered.
func (s *BaseFqlParserListener) EnterDispatchWithClause(ctx *DispatchWithClauseContext) {}

// ExitDispatchWithClause is called when production dispatchWithClause is exited.
func (s *BaseFqlParserListener) ExitDispatchWithClause(ctx *DispatchWithClauseContext) {}

// EnterDispatchOptionsClause is called when production dispatchOptionsClause is entered.
func (s *BaseFqlParserListener) EnterDispatchOptionsClause(ctx *DispatchOptionsClauseContext) {}

// ExitDispatchOptionsClause is called when production dispatchOptionsClause is exited.
func (s *BaseFqlParserListener) ExitDispatchOptionsClause(ctx *DispatchOptionsClauseContext) {}

// EnterWaitForEventExpression is called when production waitForEventExpression is entered.
func (s *BaseFqlParserListener) EnterWaitForEventExpression(ctx *WaitForEventExpressionContext) {}

// ExitWaitForEventExpression is called when production waitForEventExpression is exited.
func (s *BaseFqlParserListener) ExitWaitForEventExpression(ctx *WaitForEventExpressionContext) {}

// EnterWaitForPredicateExpression is called when production waitForPredicateExpression is entered.
func (s *BaseFqlParserListener) EnterWaitForPredicateExpression(ctx *WaitForPredicateExpressionContext) {
}

// ExitWaitForPredicateExpression is called when production waitForPredicateExpression is exited.
func (s *BaseFqlParserListener) ExitWaitForPredicateExpression(ctx *WaitForPredicateExpressionContext) {
}

// EnterWaitForPredicate is called when production waitForPredicate is entered.
func (s *BaseFqlParserListener) EnterWaitForPredicate(ctx *WaitForPredicateContext) {}

// ExitWaitForPredicate is called when production waitForPredicate is exited.
func (s *BaseFqlParserListener) ExitWaitForPredicate(ctx *WaitForPredicateContext) {}

// EnterWaitForEventName is called when production waitForEventName is entered.
func (s *BaseFqlParserListener) EnterWaitForEventName(ctx *WaitForEventNameContext) {}

// ExitWaitForEventName is called when production waitForEventName is exited.
func (s *BaseFqlParserListener) ExitWaitForEventName(ctx *WaitForEventNameContext) {}

// EnterWaitForEventSource is called when production waitForEventSource is entered.
func (s *BaseFqlParserListener) EnterWaitForEventSource(ctx *WaitForEventSourceContext) {}

// ExitWaitForEventSource is called when production waitForEventSource is exited.
func (s *BaseFqlParserListener) ExitWaitForEventSource(ctx *WaitForEventSourceContext) {}

// EnterOptionsClause is called when production optionsClause is entered.
func (s *BaseFqlParserListener) EnterOptionsClause(ctx *OptionsClauseContext) {}

// ExitOptionsClause is called when production optionsClause is exited.
func (s *BaseFqlParserListener) ExitOptionsClause(ctx *OptionsClauseContext) {}

// EnterTimeoutClause is called when production timeoutClause is entered.
func (s *BaseFqlParserListener) EnterTimeoutClause(ctx *TimeoutClauseContext) {}

// ExitTimeoutClause is called when production timeoutClause is exited.
func (s *BaseFqlParserListener) ExitTimeoutClause(ctx *TimeoutClauseContext) {}

// EnterEveryClause is called when production everyClause is entered.
func (s *BaseFqlParserListener) EnterEveryClause(ctx *EveryClauseContext) {}

// ExitEveryClause is called when production everyClause is exited.
func (s *BaseFqlParserListener) ExitEveryClause(ctx *EveryClauseContext) {}

// EnterEveryClauseValue is called when production everyClauseValue is entered.
func (s *BaseFqlParserListener) EnterEveryClauseValue(ctx *EveryClauseValueContext) {}

// ExitEveryClauseValue is called when production everyClauseValue is exited.
func (s *BaseFqlParserListener) ExitEveryClauseValue(ctx *EveryClauseValueContext) {}

// EnterBackoffClause is called when production backoffClause is entered.
func (s *BaseFqlParserListener) EnterBackoffClause(ctx *BackoffClauseContext) {}

// ExitBackoffClause is called when production backoffClause is exited.
func (s *BaseFqlParserListener) ExitBackoffClause(ctx *BackoffClauseContext) {}

// EnterJitterClause is called when production jitterClause is entered.
func (s *BaseFqlParserListener) EnterJitterClause(ctx *JitterClauseContext) {}

// ExitJitterClause is called when production jitterClause is exited.
func (s *BaseFqlParserListener) ExitJitterClause(ctx *JitterClauseContext) {}

// EnterJitterClauseValue is called when production jitterClauseValue is entered.
func (s *BaseFqlParserListener) EnterJitterClauseValue(ctx *JitterClauseValueContext) {}

// ExitJitterClauseValue is called when production jitterClauseValue is exited.
func (s *BaseFqlParserListener) ExitJitterClauseValue(ctx *JitterClauseValueContext) {}

// EnterBackoffStrategy is called when production backoffStrategy is entered.
func (s *BaseFqlParserListener) EnterBackoffStrategy(ctx *BackoffStrategyContext) {}

// ExitBackoffStrategy is called when production backoffStrategy is exited.
func (s *BaseFqlParserListener) ExitBackoffStrategy(ctx *BackoffStrategyContext) {}

// EnterWaitForOrThrowClause is called when production waitForOrThrowClause is entered.
func (s *BaseFqlParserListener) EnterWaitForOrThrowClause(ctx *WaitForOrThrowClauseContext) {}

// ExitWaitForOrThrowClause is called when production waitForOrThrowClause is exited.
func (s *BaseFqlParserListener) ExitWaitForOrThrowClause(ctx *WaitForOrThrowClauseContext) {}

// EnterParam is called when production param is entered.
func (s *BaseFqlParserListener) EnterParam(ctx *ParamContext) {}

// ExitParam is called when production param is exited.
func (s *BaseFqlParserListener) ExitParam(ctx *ParamContext) {}

// EnterVariable is called when production variable is entered.
func (s *BaseFqlParserListener) EnterVariable(ctx *VariableContext) {}

// ExitVariable is called when production variable is exited.
func (s *BaseFqlParserListener) ExitVariable(ctx *VariableContext) {}

// EnterLiteral is called when production literal is entered.
func (s *BaseFqlParserListener) EnterLiteral(ctx *LiteralContext) {}

// ExitLiteral is called when production literal is exited.
func (s *BaseFqlParserListener) ExitLiteral(ctx *LiteralContext) {}

// EnterArrayLiteral is called when production arrayLiteral is entered.
func (s *BaseFqlParserListener) EnterArrayLiteral(ctx *ArrayLiteralContext) {}

// ExitArrayLiteral is called when production arrayLiteral is exited.
func (s *BaseFqlParserListener) ExitArrayLiteral(ctx *ArrayLiteralContext) {}

// EnterObjectLiteral is called when production objectLiteral is entered.
func (s *BaseFqlParserListener) EnterObjectLiteral(ctx *ObjectLiteralContext) {}

// ExitObjectLiteral is called when production objectLiteral is exited.
func (s *BaseFqlParserListener) ExitObjectLiteral(ctx *ObjectLiteralContext) {}

// EnterBooleanLiteral is called when production booleanLiteral is entered.
func (s *BaseFqlParserListener) EnterBooleanLiteral(ctx *BooleanLiteralContext) {}

// ExitBooleanLiteral is called when production booleanLiteral is exited.
func (s *BaseFqlParserListener) ExitBooleanLiteral(ctx *BooleanLiteralContext) {}

// EnterStringLiteral is called when production stringLiteral is entered.
func (s *BaseFqlParserListener) EnterStringLiteral(ctx *StringLiteralContext) {}

// ExitStringLiteral is called when production stringLiteral is exited.
func (s *BaseFqlParserListener) ExitStringLiteral(ctx *StringLiteralContext) {}

// EnterTemplateLiteral is called when production templateLiteral is entered.
func (s *BaseFqlParserListener) EnterTemplateLiteral(ctx *TemplateLiteralContext) {}

// ExitTemplateLiteral is called when production templateLiteral is exited.
func (s *BaseFqlParserListener) ExitTemplateLiteral(ctx *TemplateLiteralContext) {}

// EnterTemplateElement is called when production templateElement is entered.
func (s *BaseFqlParserListener) EnterTemplateElement(ctx *TemplateElementContext) {}

// ExitTemplateElement is called when production templateElement is exited.
func (s *BaseFqlParserListener) ExitTemplateElement(ctx *TemplateElementContext) {}

// EnterFloatLiteral is called when production floatLiteral is entered.
func (s *BaseFqlParserListener) EnterFloatLiteral(ctx *FloatLiteralContext) {}

// ExitFloatLiteral is called when production floatLiteral is exited.
func (s *BaseFqlParserListener) ExitFloatLiteral(ctx *FloatLiteralContext) {}

// EnterIntegerLiteral is called when production integerLiteral is entered.
func (s *BaseFqlParserListener) EnterIntegerLiteral(ctx *IntegerLiteralContext) {}

// ExitIntegerLiteral is called when production integerLiteral is exited.
func (s *BaseFqlParserListener) ExitIntegerLiteral(ctx *IntegerLiteralContext) {}

// EnterNoneLiteral is called when production noneLiteral is entered.
func (s *BaseFqlParserListener) EnterNoneLiteral(ctx *NoneLiteralContext) {}

// ExitNoneLiteral is called when production noneLiteral is exited.
func (s *BaseFqlParserListener) ExitNoneLiteral(ctx *NoneLiteralContext) {}

// EnterPropertyAssignment is called when production propertyAssignment is entered.
func (s *BaseFqlParserListener) EnterPropertyAssignment(ctx *PropertyAssignmentContext) {}

// ExitPropertyAssignment is called when production propertyAssignment is exited.
func (s *BaseFqlParserListener) ExitPropertyAssignment(ctx *PropertyAssignmentContext) {}

// EnterComputedPropertyName is called when production computedPropertyName is entered.
func (s *BaseFqlParserListener) EnterComputedPropertyName(ctx *ComputedPropertyNameContext) {}

// ExitComputedPropertyName is called when production computedPropertyName is exited.
func (s *BaseFqlParserListener) ExitComputedPropertyName(ctx *ComputedPropertyNameContext) {}

// EnterPropertyName is called when production propertyName is entered.
func (s *BaseFqlParserListener) EnterPropertyName(ctx *PropertyNameContext) {}

// ExitPropertyName is called when production propertyName is exited.
func (s *BaseFqlParserListener) ExitPropertyName(ctx *PropertyNameContext) {}

// EnterNamespaceIdentifier is called when production namespaceIdentifier is entered.
func (s *BaseFqlParserListener) EnterNamespaceIdentifier(ctx *NamespaceIdentifierContext) {}

// ExitNamespaceIdentifier is called when production namespaceIdentifier is exited.
func (s *BaseFqlParserListener) ExitNamespaceIdentifier(ctx *NamespaceIdentifierContext) {}

// EnterNamespace is called when production namespace is entered.
func (s *BaseFqlParserListener) EnterNamespace(ctx *NamespaceContext) {}

// ExitNamespace is called when production namespace is exited.
func (s *BaseFqlParserListener) ExitNamespace(ctx *NamespaceContext) {}

// EnterMemberExpression is called when production memberExpression is entered.
func (s *BaseFqlParserListener) EnterMemberExpression(ctx *MemberExpressionContext) {}

// ExitMemberExpression is called when production memberExpression is exited.
func (s *BaseFqlParserListener) ExitMemberExpression(ctx *MemberExpressionContext) {}

// EnterMemberExpressionSource is called when production memberExpressionSource is entered.
func (s *BaseFqlParserListener) EnterMemberExpressionSource(ctx *MemberExpressionSourceContext) {}

// ExitMemberExpressionSource is called when production memberExpressionSource is exited.
func (s *BaseFqlParserListener) ExitMemberExpressionSource(ctx *MemberExpressionSourceContext) {}

// EnterFunctionCallExpression is called when production functionCallExpression is entered.
func (s *BaseFqlParserListener) EnterFunctionCallExpression(ctx *FunctionCallExpressionContext) {}

// ExitFunctionCallExpression is called when production functionCallExpression is exited.
func (s *BaseFqlParserListener) ExitFunctionCallExpression(ctx *FunctionCallExpressionContext) {}

// EnterFunctionCall is called when production functionCall is entered.
func (s *BaseFqlParserListener) EnterFunctionCall(ctx *FunctionCallContext) {}

// ExitFunctionCall is called when production functionCall is exited.
func (s *BaseFqlParserListener) ExitFunctionCall(ctx *FunctionCallContext) {}

// EnterFunctionName is called when production functionName is entered.
func (s *BaseFqlParserListener) EnterFunctionName(ctx *FunctionNameContext) {}

// ExitFunctionName is called when production functionName is exited.
func (s *BaseFqlParserListener) ExitFunctionName(ctx *FunctionNameContext) {}

// EnterArgumentList is called when production argumentList is entered.
func (s *BaseFqlParserListener) EnterArgumentList(ctx *ArgumentListContext) {}

// ExitArgumentList is called when production argumentList is exited.
func (s *BaseFqlParserListener) ExitArgumentList(ctx *ArgumentListContext) {}

// EnterMemberExpressionPath is called when production memberExpressionPath is entered.
func (s *BaseFqlParserListener) EnterMemberExpressionPath(ctx *MemberExpressionPathContext) {}

// ExitMemberExpressionPath is called when production memberExpressionPath is exited.
func (s *BaseFqlParserListener) ExitMemberExpressionPath(ctx *MemberExpressionPathContext) {}

// EnterArrayExpansion is called when production arrayExpansion is entered.
func (s *BaseFqlParserListener) EnterArrayExpansion(ctx *ArrayExpansionContext) {}

// ExitArrayExpansion is called when production arrayExpansion is exited.
func (s *BaseFqlParserListener) ExitArrayExpansion(ctx *ArrayExpansionContext) {}

// EnterArrayContraction is called when production arrayContraction is entered.
func (s *BaseFqlParserListener) EnterArrayContraction(ctx *ArrayContractionContext) {}

// ExitArrayContraction is called when production arrayContraction is exited.
func (s *BaseFqlParserListener) ExitArrayContraction(ctx *ArrayContractionContext) {}

// EnterArrayQuestionMark is called when production arrayQuestionMark is entered.
func (s *BaseFqlParserListener) EnterArrayQuestionMark(ctx *ArrayQuestionMarkContext) {}

// ExitArrayQuestionMark is called when production arrayQuestionMark is exited.
func (s *BaseFqlParserListener) ExitArrayQuestionMark(ctx *ArrayQuestionMarkContext) {}

// EnterArrayQuestionQuantifier is called when production arrayQuestionQuantifier is entered.
func (s *BaseFqlParserListener) EnterArrayQuestionQuantifier(ctx *ArrayQuestionQuantifierContext) {}

// ExitArrayQuestionQuantifier is called when production arrayQuestionQuantifier is exited.
func (s *BaseFqlParserListener) ExitArrayQuestionQuantifier(ctx *ArrayQuestionQuantifierContext) {}

// EnterArrayQuestionQuantifierValue is called when production arrayQuestionQuantifierValue is entered.
func (s *BaseFqlParserListener) EnterArrayQuestionQuantifierValue(ctx *ArrayQuestionQuantifierValueContext) {
}

// ExitArrayQuestionQuantifierValue is called when production arrayQuestionQuantifierValue is exited.
func (s *BaseFqlParserListener) ExitArrayQuestionQuantifierValue(ctx *ArrayQuestionQuantifierValueContext) {
}

// EnterArrayApply is called when production arrayApply is entered.
func (s *BaseFqlParserListener) EnterArrayApply(ctx *ArrayApplyContext) {}

// ExitArrayApply is called when production arrayApply is exited.
func (s *BaseFqlParserListener) ExitArrayApply(ctx *ArrayApplyContext) {}

// EnterInlineExpression is called when production inlineExpression is entered.
func (s *BaseFqlParserListener) EnterInlineExpression(ctx *InlineExpressionContext) {}

// ExitInlineExpression is called when production inlineExpression is exited.
func (s *BaseFqlParserListener) ExitInlineExpression(ctx *InlineExpressionContext) {}

// EnterInlineFilter is called when production inlineFilter is entered.
func (s *BaseFqlParserListener) EnterInlineFilter(ctx *InlineFilterContext) {}

// ExitInlineFilter is called when production inlineFilter is exited.
func (s *BaseFqlParserListener) ExitInlineFilter(ctx *InlineFilterContext) {}

// EnterInlineLimit is called when production inlineLimit is entered.
func (s *BaseFqlParserListener) EnterInlineLimit(ctx *InlineLimitContext) {}

// ExitInlineLimit is called when production inlineLimit is exited.
func (s *BaseFqlParserListener) ExitInlineLimit(ctx *InlineLimitContext) {}

// EnterInlineReturn is called when production inlineReturn is entered.
func (s *BaseFqlParserListener) EnterInlineReturn(ctx *InlineReturnContext) {}

// ExitInlineReturn is called when production inlineReturn is exited.
func (s *BaseFqlParserListener) ExitInlineReturn(ctx *InlineReturnContext) {}

// EnterSafeReservedWord is called when production safeReservedWord is entered.
func (s *BaseFqlParserListener) EnterSafeReservedWord(ctx *SafeReservedWordContext) {}

// ExitSafeReservedWord is called when production safeReservedWord is exited.
func (s *BaseFqlParserListener) ExitSafeReservedWord(ctx *SafeReservedWordContext) {}

// EnterUnsafeReservedWord is called when production unsafeReservedWord is entered.
func (s *BaseFqlParserListener) EnterUnsafeReservedWord(ctx *UnsafeReservedWordContext) {}

// ExitUnsafeReservedWord is called when production unsafeReservedWord is exited.
func (s *BaseFqlParserListener) ExitUnsafeReservedWord(ctx *UnsafeReservedWordContext) {}

// EnterDurationLiteral is called when production durationLiteral is entered.
func (s *BaseFqlParserListener) EnterDurationLiteral(ctx *DurationLiteralContext) {}

// ExitDurationLiteral is called when production durationLiteral is exited.
func (s *BaseFqlParserListener) ExitDurationLiteral(ctx *DurationLiteralContext) {}

// EnterRangeOperator is called when production rangeOperator is entered.
func (s *BaseFqlParserListener) EnterRangeOperator(ctx *RangeOperatorContext) {}

// ExitRangeOperator is called when production rangeOperator is exited.
func (s *BaseFqlParserListener) ExitRangeOperator(ctx *RangeOperatorContext) {}

// EnterRangeOperand is called when production rangeOperand is entered.
func (s *BaseFqlParserListener) EnterRangeOperand(ctx *RangeOperandContext) {}

// ExitRangeOperand is called when production rangeOperand is exited.
func (s *BaseFqlParserListener) ExitRangeOperand(ctx *RangeOperandContext) {}

// EnterExpression is called when production expression is entered.
func (s *BaseFqlParserListener) EnterExpression(ctx *ExpressionContext) {}

// ExitExpression is called when production expression is exited.
func (s *BaseFqlParserListener) ExitExpression(ctx *ExpressionContext) {}

// EnterPredicate is called when production predicate is entered.
func (s *BaseFqlParserListener) EnterPredicate(ctx *PredicateContext) {}

// ExitPredicate is called when production predicate is exited.
func (s *BaseFqlParserListener) ExitPredicate(ctx *PredicateContext) {}

// EnterExpressionAtom is called when production expressionAtom is entered.
func (s *BaseFqlParserListener) EnterExpressionAtom(ctx *ExpressionAtomContext) {}

// ExitExpressionAtom is called when production expressionAtom is exited.
func (s *BaseFqlParserListener) ExitExpressionAtom(ctx *ExpressionAtomContext) {}

// EnterMatchExpression is called when production matchExpression is entered.
func (s *BaseFqlParserListener) EnterMatchExpression(ctx *MatchExpressionContext) {}

// ExitMatchExpression is called when production matchExpression is exited.
func (s *BaseFqlParserListener) ExitMatchExpression(ctx *MatchExpressionContext) {}

// EnterMatchPatternArms is called when production matchPatternArms is entered.
func (s *BaseFqlParserListener) EnterMatchPatternArms(ctx *MatchPatternArmsContext) {}

// ExitMatchPatternArms is called when production matchPatternArms is exited.
func (s *BaseFqlParserListener) ExitMatchPatternArms(ctx *MatchPatternArmsContext) {}

// EnterMatchPatternArmList is called when production matchPatternArmList is entered.
func (s *BaseFqlParserListener) EnterMatchPatternArmList(ctx *MatchPatternArmListContext) {}

// ExitMatchPatternArmList is called when production matchPatternArmList is exited.
func (s *BaseFqlParserListener) ExitMatchPatternArmList(ctx *MatchPatternArmListContext) {}

// EnterMatchGuardArms is called when production matchGuardArms is entered.
func (s *BaseFqlParserListener) EnterMatchGuardArms(ctx *MatchGuardArmsContext) {}

// ExitMatchGuardArms is called when production matchGuardArms is exited.
func (s *BaseFqlParserListener) ExitMatchGuardArms(ctx *MatchGuardArmsContext) {}

// EnterMatchGuardArmList is called when production matchGuardArmList is entered.
func (s *BaseFqlParserListener) EnterMatchGuardArmList(ctx *MatchGuardArmListContext) {}

// ExitMatchGuardArmList is called when production matchGuardArmList is exited.
func (s *BaseFqlParserListener) ExitMatchGuardArmList(ctx *MatchGuardArmListContext) {}

// EnterMatchPatternArm is called when production matchPatternArm is entered.
func (s *BaseFqlParserListener) EnterMatchPatternArm(ctx *MatchPatternArmContext) {}

// ExitMatchPatternArm is called when production matchPatternArm is exited.
func (s *BaseFqlParserListener) ExitMatchPatternArm(ctx *MatchPatternArmContext) {}

// EnterMatchPatternGuard is called when production matchPatternGuard is entered.
func (s *BaseFqlParserListener) EnterMatchPatternGuard(ctx *MatchPatternGuardContext) {}

// ExitMatchPatternGuard is called when production matchPatternGuard is exited.
func (s *BaseFqlParserListener) ExitMatchPatternGuard(ctx *MatchPatternGuardContext) {}

// EnterMatchGuardArm is called when production matchGuardArm is entered.
func (s *BaseFqlParserListener) EnterMatchGuardArm(ctx *MatchGuardArmContext) {}

// ExitMatchGuardArm is called when production matchGuardArm is exited.
func (s *BaseFqlParserListener) ExitMatchGuardArm(ctx *MatchGuardArmContext) {}

// EnterMatchDefaultArm is called when production matchDefaultArm is entered.
func (s *BaseFqlParserListener) EnterMatchDefaultArm(ctx *MatchDefaultArmContext) {}

// ExitMatchDefaultArm is called when production matchDefaultArm is exited.
func (s *BaseFqlParserListener) ExitMatchDefaultArm(ctx *MatchDefaultArmContext) {}

// EnterMatchPattern is called when production matchPattern is entered.
func (s *BaseFqlParserListener) EnterMatchPattern(ctx *MatchPatternContext) {}

// ExitMatchPattern is called when production matchPattern is exited.
func (s *BaseFqlParserListener) ExitMatchPattern(ctx *MatchPatternContext) {}

// EnterMatchBindingPattern is called when production matchBindingPattern is entered.
func (s *BaseFqlParserListener) EnterMatchBindingPattern(ctx *MatchBindingPatternContext) {}

// ExitMatchBindingPattern is called when production matchBindingPattern is exited.
func (s *BaseFqlParserListener) ExitMatchBindingPattern(ctx *MatchBindingPatternContext) {}

// EnterMatchLiteralPattern is called when production matchLiteralPattern is entered.
func (s *BaseFqlParserListener) EnterMatchLiteralPattern(ctx *MatchLiteralPatternContext) {}

// ExitMatchLiteralPattern is called when production matchLiteralPattern is exited.
func (s *BaseFqlParserListener) ExitMatchLiteralPattern(ctx *MatchLiteralPatternContext) {}

// EnterMatchObjectPattern is called when production matchObjectPattern is entered.
func (s *BaseFqlParserListener) EnterMatchObjectPattern(ctx *MatchObjectPatternContext) {}

// ExitMatchObjectPattern is called when production matchObjectPattern is exited.
func (s *BaseFqlParserListener) ExitMatchObjectPattern(ctx *MatchObjectPatternContext) {}

// EnterMatchObjectPatternProperty is called when production matchObjectPatternProperty is entered.
func (s *BaseFqlParserListener) EnterMatchObjectPatternProperty(ctx *MatchObjectPatternPropertyContext) {
}

// ExitMatchObjectPatternProperty is called when production matchObjectPatternProperty is exited.
func (s *BaseFqlParserListener) ExitMatchObjectPatternProperty(ctx *MatchObjectPatternPropertyContext) {
}

// EnterMatchObjectPatternKey is called when production matchObjectPatternKey is entered.
func (s *BaseFqlParserListener) EnterMatchObjectPatternKey(ctx *MatchObjectPatternKeyContext) {}

// ExitMatchObjectPatternKey is called when production matchObjectPatternKey is exited.
func (s *BaseFqlParserListener) ExitMatchObjectPatternKey(ctx *MatchObjectPatternKeyContext) {}

// EnterImplicitMemberExpression is called when production implicitMemberExpression is entered.
func (s *BaseFqlParserListener) EnterImplicitMemberExpression(ctx *ImplicitMemberExpressionContext) {}

// ExitImplicitMemberExpression is called when production implicitMemberExpression is exited.
func (s *BaseFqlParserListener) ExitImplicitMemberExpression(ctx *ImplicitMemberExpressionContext) {}

// EnterImplicitCurrentExpression is called when production implicitCurrentExpression is entered.
func (s *BaseFqlParserListener) EnterImplicitCurrentExpression(ctx *ImplicitCurrentExpressionContext) {
}

// ExitImplicitCurrentExpression is called when production implicitCurrentExpression is exited.
func (s *BaseFqlParserListener) ExitImplicitCurrentExpression(ctx *ImplicitCurrentExpressionContext) {
}

// EnterImplicitMemberExpressionStart is called when production implicitMemberExpressionStart is entered.
func (s *BaseFqlParserListener) EnterImplicitMemberExpressionStart(ctx *ImplicitMemberExpressionStartContext) {
}

// ExitImplicitMemberExpressionStart is called when production implicitMemberExpressionStart is exited.
func (s *BaseFqlParserListener) ExitImplicitMemberExpressionStart(ctx *ImplicitMemberExpressionStartContext) {
}

// EnterQueryExpression is called when production queryExpression is entered.
func (s *BaseFqlParserListener) EnterQueryExpression(ctx *QueryExpressionContext) {}

// ExitQueryExpression is called when production queryExpression is exited.
func (s *BaseFqlParserListener) ExitQueryExpression(ctx *QueryExpressionContext) {}

// EnterQueryModifier is called when production queryModifier is entered.
func (s *BaseFqlParserListener) EnterQueryModifier(ctx *QueryModifierContext) {}

// ExitQueryModifier is called when production queryModifier is exited.
func (s *BaseFqlParserListener) ExitQueryModifier(ctx *QueryModifierContext) {}

// EnterQueryPayload is called when production queryPayload is entered.
func (s *BaseFqlParserListener) EnterQueryPayload(ctx *QueryPayloadContext) {}

// ExitQueryPayload is called when production queryPayload is exited.
func (s *BaseFqlParserListener) ExitQueryPayload(ctx *QueryPayloadContext) {}

// EnterQueryWithOpt is called when production queryWithOpt is entered.
func (s *BaseFqlParserListener) EnterQueryWithOpt(ctx *QueryWithOptContext) {}

// ExitQueryWithOpt is called when production queryWithOpt is exited.
func (s *BaseFqlParserListener) ExitQueryWithOpt(ctx *QueryWithOptContext) {}

// EnterQueryLiteral is called when production queryLiteral is entered.
func (s *BaseFqlParserListener) EnterQueryLiteral(ctx *QueryLiteralContext) {}

// ExitQueryLiteral is called when production queryLiteral is exited.
func (s *BaseFqlParserListener) ExitQueryLiteral(ctx *QueryLiteralContext) {}

// EnterArrayOperator is called when production arrayOperator is entered.
func (s *BaseFqlParserListener) EnterArrayOperator(ctx *ArrayOperatorContext) {}

// ExitArrayOperator is called when production arrayOperator is exited.
func (s *BaseFqlParserListener) ExitArrayOperator(ctx *ArrayOperatorContext) {}

// EnterEqualityOperator is called when production equalityOperator is entered.
func (s *BaseFqlParserListener) EnterEqualityOperator(ctx *EqualityOperatorContext) {}

// ExitEqualityOperator is called when production equalityOperator is exited.
func (s *BaseFqlParserListener) ExitEqualityOperator(ctx *EqualityOperatorContext) {}

// EnterInOperator is called when production inOperator is entered.
func (s *BaseFqlParserListener) EnterInOperator(ctx *InOperatorContext) {}

// ExitInOperator is called when production inOperator is exited.
func (s *BaseFqlParserListener) ExitInOperator(ctx *InOperatorContext) {}

// EnterLikeOperator is called when production likeOperator is entered.
func (s *BaseFqlParserListener) EnterLikeOperator(ctx *LikeOperatorContext) {}

// ExitLikeOperator is called when production likeOperator is exited.
func (s *BaseFqlParserListener) ExitLikeOperator(ctx *LikeOperatorContext) {}

// EnterUnaryOperator is called when production unaryOperator is entered.
func (s *BaseFqlParserListener) EnterUnaryOperator(ctx *UnaryOperatorContext) {}

// ExitUnaryOperator is called when production unaryOperator is exited.
func (s *BaseFqlParserListener) ExitUnaryOperator(ctx *UnaryOperatorContext) {}

// EnterRegexpOperator is called when production regexpOperator is entered.
func (s *BaseFqlParserListener) EnterRegexpOperator(ctx *RegexpOperatorContext) {}

// ExitRegexpOperator is called when production regexpOperator is exited.
func (s *BaseFqlParserListener) ExitRegexpOperator(ctx *RegexpOperatorContext) {}

// EnterLogicalAndOperator is called when production logicalAndOperator is entered.
func (s *BaseFqlParserListener) EnterLogicalAndOperator(ctx *LogicalAndOperatorContext) {}

// ExitLogicalAndOperator is called when production logicalAndOperator is exited.
func (s *BaseFqlParserListener) ExitLogicalAndOperator(ctx *LogicalAndOperatorContext) {}

// EnterLogicalOrOperator is called when production logicalOrOperator is entered.
func (s *BaseFqlParserListener) EnterLogicalOrOperator(ctx *LogicalOrOperatorContext) {}

// ExitLogicalOrOperator is called when production logicalOrOperator is exited.
func (s *BaseFqlParserListener) ExitLogicalOrOperator(ctx *LogicalOrOperatorContext) {}

// EnterMultiplicativeOperator is called when production multiplicativeOperator is entered.
func (s *BaseFqlParserListener) EnterMultiplicativeOperator(ctx *MultiplicativeOperatorContext) {}

// ExitMultiplicativeOperator is called when production multiplicativeOperator is exited.
func (s *BaseFqlParserListener) ExitMultiplicativeOperator(ctx *MultiplicativeOperatorContext) {}

// EnterAdditiveOperator is called when production additiveOperator is entered.
func (s *BaseFqlParserListener) EnterAdditiveOperator(ctx *AdditiveOperatorContext) {}

// ExitAdditiveOperator is called when production additiveOperator is exited.
func (s *BaseFqlParserListener) ExitAdditiveOperator(ctx *AdditiveOperatorContext) {}

// EnterErrorOperator is called when production errorOperator is entered.
func (s *BaseFqlParserListener) EnterErrorOperator(ctx *ErrorOperatorContext) {}

// ExitErrorOperator is called when production errorOperator is exited.
func (s *BaseFqlParserListener) ExitErrorOperator(ctx *ErrorOperatorContext) {}
