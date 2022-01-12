// Code generated from antlr/FqlParser.g4 by ANTLR 4.9.3. DO NOT EDIT.

package fql // FqlParser
import "github.com/antlr/antlr4/runtime/Go/antlr"

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

// EnterCollectSelector is called when production collectSelector is entered.
func (s *BaseFqlParserListener) EnterCollectSelector(ctx *CollectSelectorContext) {}

// ExitCollectSelector is called when production collectSelector is exited.
func (s *BaseFqlParserListener) ExitCollectSelector(ctx *CollectSelectorContext) {}

// EnterCollectGrouping is called when production collectGrouping is entered.
func (s *BaseFqlParserListener) EnterCollectGrouping(ctx *CollectGroupingContext) {}

// ExitCollectGrouping is called when production collectGrouping is exited.
func (s *BaseFqlParserListener) ExitCollectGrouping(ctx *CollectGroupingContext) {}

// EnterCollectAggregator is called when production collectAggregator is entered.
func (s *BaseFqlParserListener) EnterCollectAggregator(ctx *CollectAggregatorContext) {}

// ExitCollectAggregator is called when production collectAggregator is exited.
func (s *BaseFqlParserListener) ExitCollectAggregator(ctx *CollectAggregatorContext) {}

// EnterCollectAggregateSelector is called when production collectAggregateSelector is entered.
func (s *BaseFqlParserListener) EnterCollectAggregateSelector(ctx *CollectAggregateSelectorContext) {}

// ExitCollectAggregateSelector is called when production collectAggregateSelector is exited.
func (s *BaseFqlParserListener) ExitCollectAggregateSelector(ctx *CollectAggregateSelectorContext) {}

// EnterCollectGroupVariable is called when production collectGroupVariable is entered.
func (s *BaseFqlParserListener) EnterCollectGroupVariable(ctx *CollectGroupVariableContext) {}

// ExitCollectGroupVariable is called when production collectGroupVariable is exited.
func (s *BaseFqlParserListener) ExitCollectGroupVariable(ctx *CollectGroupVariableContext) {}

// EnterCollectCounter is called when production collectCounter is entered.
func (s *BaseFqlParserListener) EnterCollectCounter(ctx *CollectCounterContext) {}

// ExitCollectCounter is called when production collectCounter is exited.
func (s *BaseFqlParserListener) ExitCollectCounter(ctx *CollectCounterContext) {}

// EnterWaitForExpression is called when production waitForExpression is entered.
func (s *BaseFqlParserListener) EnterWaitForExpression(ctx *WaitForExpressionContext) {}

// ExitWaitForExpression is called when production waitForExpression is exited.
func (s *BaseFqlParserListener) ExitWaitForExpression(ctx *WaitForExpressionContext) {}

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

// EnterSafeReservedWord is called when production safeReservedWord is entered.
func (s *BaseFqlParserListener) EnterSafeReservedWord(ctx *SafeReservedWordContext) {}

// ExitSafeReservedWord is called when production safeReservedWord is exited.
func (s *BaseFqlParserListener) ExitSafeReservedWord(ctx *SafeReservedWordContext) {}

// EnterUnsafeReservedWord is called when production unsafeReservedWord is entered.
func (s *BaseFqlParserListener) EnterUnsafeReservedWord(ctx *UnsafeReservedWordContext) {}

// ExitUnsafeReservedWord is called when production unsafeReservedWord is exited.
func (s *BaseFqlParserListener) ExitUnsafeReservedWord(ctx *UnsafeReservedWordContext) {}

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
