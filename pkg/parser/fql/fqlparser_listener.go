// Code generated from antlr/FqlParser.g4 by ANTLR 4.9.3. DO NOT EDIT.

package fql // FqlParser
import "github.com/antlr/antlr4/runtime/Go/antlr"

// FqlParserListener is a complete listener for a parse tree produced by FqlParser.
type FqlParserListener interface {
	antlr.ParseTreeListener

	// EnterProgram is called when entering the program production.
	EnterProgram(c *ProgramContext)

	// EnterHead is called when entering the head production.
	EnterHead(c *HeadContext)

	// EnterUseExpression is called when entering the useExpression production.
	EnterUseExpression(c *UseExpressionContext)

	// EnterUse is called when entering the use production.
	EnterUse(c *UseContext)

	// EnterBody is called when entering the body production.
	EnterBody(c *BodyContext)

	// EnterBodyStatement is called when entering the bodyStatement production.
	EnterBodyStatement(c *BodyStatementContext)

	// EnterBodyExpression is called when entering the bodyExpression production.
	EnterBodyExpression(c *BodyExpressionContext)

	// EnterVariableDeclaration is called when entering the variableDeclaration production.
	EnterVariableDeclaration(c *VariableDeclarationContext)

	// EnterReturnExpression is called when entering the returnExpression production.
	EnterReturnExpression(c *ReturnExpressionContext)

	// EnterForExpression is called when entering the forExpression production.
	EnterForExpression(c *ForExpressionContext)

	// EnterForExpressionSource is called when entering the forExpressionSource production.
	EnterForExpressionSource(c *ForExpressionSourceContext)

	// EnterForExpressionClause is called when entering the forExpressionClause production.
	EnterForExpressionClause(c *ForExpressionClauseContext)

	// EnterForExpressionStatement is called when entering the forExpressionStatement production.
	EnterForExpressionStatement(c *ForExpressionStatementContext)

	// EnterForExpressionBody is called when entering the forExpressionBody production.
	EnterForExpressionBody(c *ForExpressionBodyContext)

	// EnterForExpressionReturn is called when entering the forExpressionReturn production.
	EnterForExpressionReturn(c *ForExpressionReturnContext)

	// EnterFilterClause is called when entering the filterClause production.
	EnterFilterClause(c *FilterClauseContext)

	// EnterLimitClause is called when entering the limitClause production.
	EnterLimitClause(c *LimitClauseContext)

	// EnterLimitClauseValue is called when entering the limitClauseValue production.
	EnterLimitClauseValue(c *LimitClauseValueContext)

	// EnterSortClause is called when entering the sortClause production.
	EnterSortClause(c *SortClauseContext)

	// EnterSortClauseExpression is called when entering the sortClauseExpression production.
	EnterSortClauseExpression(c *SortClauseExpressionContext)

	// EnterCollectClause is called when entering the collectClause production.
	EnterCollectClause(c *CollectClauseContext)

	// EnterCollectSelector is called when entering the collectSelector production.
	EnterCollectSelector(c *CollectSelectorContext)

	// EnterCollectGrouping is called when entering the collectGrouping production.
	EnterCollectGrouping(c *CollectGroupingContext)

	// EnterCollectAggregator is called when entering the collectAggregator production.
	EnterCollectAggregator(c *CollectAggregatorContext)

	// EnterCollectAggregateSelector is called when entering the collectAggregateSelector production.
	EnterCollectAggregateSelector(c *CollectAggregateSelectorContext)

	// EnterCollectGroupVariable is called when entering the collectGroupVariable production.
	EnterCollectGroupVariable(c *CollectGroupVariableContext)

	// EnterCollectCounter is called when entering the collectCounter production.
	EnterCollectCounter(c *CollectCounterContext)

	// EnterWaitForExpression is called when entering the waitForExpression production.
	EnterWaitForExpression(c *WaitForExpressionContext)

	// EnterWaitForEventName is called when entering the waitForEventName production.
	EnterWaitForEventName(c *WaitForEventNameContext)

	// EnterWaitForEventSource is called when entering the waitForEventSource production.
	EnterWaitForEventSource(c *WaitForEventSourceContext)

	// EnterOptionsClause is called when entering the optionsClause production.
	EnterOptionsClause(c *OptionsClauseContext)

	// EnterTimeoutClause is called when entering the timeoutClause production.
	EnterTimeoutClause(c *TimeoutClauseContext)

	// EnterParam is called when entering the param production.
	EnterParam(c *ParamContext)

	// EnterVariable is called when entering the variable production.
	EnterVariable(c *VariableContext)

	// EnterLiteral is called when entering the literal production.
	EnterLiteral(c *LiteralContext)

	// EnterArrayLiteral is called when entering the arrayLiteral production.
	EnterArrayLiteral(c *ArrayLiteralContext)

	// EnterObjectLiteral is called when entering the objectLiteral production.
	EnterObjectLiteral(c *ObjectLiteralContext)

	// EnterBooleanLiteral is called when entering the booleanLiteral production.
	EnterBooleanLiteral(c *BooleanLiteralContext)

	// EnterStringLiteral is called when entering the stringLiteral production.
	EnterStringLiteral(c *StringLiteralContext)

	// EnterFloatLiteral is called when entering the floatLiteral production.
	EnterFloatLiteral(c *FloatLiteralContext)

	// EnterIntegerLiteral is called when entering the integerLiteral production.
	EnterIntegerLiteral(c *IntegerLiteralContext)

	// EnterNoneLiteral is called when entering the noneLiteral production.
	EnterNoneLiteral(c *NoneLiteralContext)

	// EnterPropertyAssignment is called when entering the propertyAssignment production.
	EnterPropertyAssignment(c *PropertyAssignmentContext)

	// EnterComputedPropertyName is called when entering the computedPropertyName production.
	EnterComputedPropertyName(c *ComputedPropertyNameContext)

	// EnterPropertyName is called when entering the propertyName production.
	EnterPropertyName(c *PropertyNameContext)

	// EnterNamespaceIdentifier is called when entering the namespaceIdentifier production.
	EnterNamespaceIdentifier(c *NamespaceIdentifierContext)

	// EnterNamespace is called when entering the namespace production.
	EnterNamespace(c *NamespaceContext)

	// EnterMemberExpression is called when entering the memberExpression production.
	EnterMemberExpression(c *MemberExpressionContext)

	// EnterMemberExpressionSource is called when entering the memberExpressionSource production.
	EnterMemberExpressionSource(c *MemberExpressionSourceContext)

	// EnterFunctionCallExpression is called when entering the functionCallExpression production.
	EnterFunctionCallExpression(c *FunctionCallExpressionContext)

	// EnterFunctionCall is called when entering the functionCall production.
	EnterFunctionCall(c *FunctionCallContext)

	// EnterFunctionName is called when entering the functionName production.
	EnterFunctionName(c *FunctionNameContext)

	// EnterArgumentList is called when entering the argumentList production.
	EnterArgumentList(c *ArgumentListContext)

	// EnterMemberExpressionPath is called when entering the memberExpressionPath production.
	EnterMemberExpressionPath(c *MemberExpressionPathContext)

	// EnterSafeReservedWord is called when entering the safeReservedWord production.
	EnterSafeReservedWord(c *SafeReservedWordContext)

	// EnterUnsafeReservedWord is called when entering the unsafeReservedWord production.
	EnterUnsafeReservedWord(c *UnsafeReservedWordContext)

	// EnterRangeOperator is called when entering the rangeOperator production.
	EnterRangeOperator(c *RangeOperatorContext)

	// EnterRangeOperand is called when entering the rangeOperand production.
	EnterRangeOperand(c *RangeOperandContext)

	// EnterExpression is called when entering the expression production.
	EnterExpression(c *ExpressionContext)

	// EnterPredicate is called when entering the predicate production.
	EnterPredicate(c *PredicateContext)

	// EnterExpressionAtom is called when entering the expressionAtom production.
	EnterExpressionAtom(c *ExpressionAtomContext)

	// EnterArrayOperator is called when entering the arrayOperator production.
	EnterArrayOperator(c *ArrayOperatorContext)

	// EnterEqualityOperator is called when entering the equalityOperator production.
	EnterEqualityOperator(c *EqualityOperatorContext)

	// EnterInOperator is called when entering the inOperator production.
	EnterInOperator(c *InOperatorContext)

	// EnterLikeOperator is called when entering the likeOperator production.
	EnterLikeOperator(c *LikeOperatorContext)

	// EnterUnaryOperator is called when entering the unaryOperator production.
	EnterUnaryOperator(c *UnaryOperatorContext)

	// EnterRegexpOperator is called when entering the regexpOperator production.
	EnterRegexpOperator(c *RegexpOperatorContext)

	// EnterLogicalAndOperator is called when entering the logicalAndOperator production.
	EnterLogicalAndOperator(c *LogicalAndOperatorContext)

	// EnterLogicalOrOperator is called when entering the logicalOrOperator production.
	EnterLogicalOrOperator(c *LogicalOrOperatorContext)

	// EnterMultiplicativeOperator is called when entering the multiplicativeOperator production.
	EnterMultiplicativeOperator(c *MultiplicativeOperatorContext)

	// EnterAdditiveOperator is called when entering the additiveOperator production.
	EnterAdditiveOperator(c *AdditiveOperatorContext)

	// EnterErrorOperator is called when entering the errorOperator production.
	EnterErrorOperator(c *ErrorOperatorContext)

	// ExitProgram is called when exiting the program production.
	ExitProgram(c *ProgramContext)

	// ExitHead is called when exiting the head production.
	ExitHead(c *HeadContext)

	// ExitUseExpression is called when exiting the useExpression production.
	ExitUseExpression(c *UseExpressionContext)

	// ExitUse is called when exiting the use production.
	ExitUse(c *UseContext)

	// ExitBody is called when exiting the body production.
	ExitBody(c *BodyContext)

	// ExitBodyStatement is called when exiting the bodyStatement production.
	ExitBodyStatement(c *BodyStatementContext)

	// ExitBodyExpression is called when exiting the bodyExpression production.
	ExitBodyExpression(c *BodyExpressionContext)

	// ExitVariableDeclaration is called when exiting the variableDeclaration production.
	ExitVariableDeclaration(c *VariableDeclarationContext)

	// ExitReturnExpression is called when exiting the returnExpression production.
	ExitReturnExpression(c *ReturnExpressionContext)

	// ExitForExpression is called when exiting the forExpression production.
	ExitForExpression(c *ForExpressionContext)

	// ExitForExpressionSource is called when exiting the forExpressionSource production.
	ExitForExpressionSource(c *ForExpressionSourceContext)

	// ExitForExpressionClause is called when exiting the forExpressionClause production.
	ExitForExpressionClause(c *ForExpressionClauseContext)

	// ExitForExpressionStatement is called when exiting the forExpressionStatement production.
	ExitForExpressionStatement(c *ForExpressionStatementContext)

	// ExitForExpressionBody is called when exiting the forExpressionBody production.
	ExitForExpressionBody(c *ForExpressionBodyContext)

	// ExitForExpressionReturn is called when exiting the forExpressionReturn production.
	ExitForExpressionReturn(c *ForExpressionReturnContext)

	// ExitFilterClause is called when exiting the filterClause production.
	ExitFilterClause(c *FilterClauseContext)

	// ExitLimitClause is called when exiting the limitClause production.
	ExitLimitClause(c *LimitClauseContext)

	// ExitLimitClauseValue is called when exiting the limitClauseValue production.
	ExitLimitClauseValue(c *LimitClauseValueContext)

	// ExitSortClause is called when exiting the sortClause production.
	ExitSortClause(c *SortClauseContext)

	// ExitSortClauseExpression is called when exiting the sortClauseExpression production.
	ExitSortClauseExpression(c *SortClauseExpressionContext)

	// ExitCollectClause is called when exiting the collectClause production.
	ExitCollectClause(c *CollectClauseContext)

	// ExitCollectSelector is called when exiting the collectSelector production.
	ExitCollectSelector(c *CollectSelectorContext)

	// ExitCollectGrouping is called when exiting the collectGrouping production.
	ExitCollectGrouping(c *CollectGroupingContext)

	// ExitCollectAggregator is called when exiting the collectAggregator production.
	ExitCollectAggregator(c *CollectAggregatorContext)

	// ExitCollectAggregateSelector is called when exiting the collectAggregateSelector production.
	ExitCollectAggregateSelector(c *CollectAggregateSelectorContext)

	// ExitCollectGroupVariable is called when exiting the collectGroupVariable production.
	ExitCollectGroupVariable(c *CollectGroupVariableContext)

	// ExitCollectCounter is called when exiting the collectCounter production.
	ExitCollectCounter(c *CollectCounterContext)

	// ExitWaitForExpression is called when exiting the waitForExpression production.
	ExitWaitForExpression(c *WaitForExpressionContext)

	// ExitWaitForEventName is called when exiting the waitForEventName production.
	ExitWaitForEventName(c *WaitForEventNameContext)

	// ExitWaitForEventSource is called when exiting the waitForEventSource production.
	ExitWaitForEventSource(c *WaitForEventSourceContext)

	// ExitOptionsClause is called when exiting the optionsClause production.
	ExitOptionsClause(c *OptionsClauseContext)

	// ExitTimeoutClause is called when exiting the timeoutClause production.
	ExitTimeoutClause(c *TimeoutClauseContext)

	// ExitParam is called when exiting the param production.
	ExitParam(c *ParamContext)

	// ExitVariable is called when exiting the variable production.
	ExitVariable(c *VariableContext)

	// ExitLiteral is called when exiting the literal production.
	ExitLiteral(c *LiteralContext)

	// ExitArrayLiteral is called when exiting the arrayLiteral production.
	ExitArrayLiteral(c *ArrayLiteralContext)

	// ExitObjectLiteral is called when exiting the objectLiteral production.
	ExitObjectLiteral(c *ObjectLiteralContext)

	// ExitBooleanLiteral is called when exiting the booleanLiteral production.
	ExitBooleanLiteral(c *BooleanLiteralContext)

	// ExitStringLiteral is called when exiting the stringLiteral production.
	ExitStringLiteral(c *StringLiteralContext)

	// ExitFloatLiteral is called when exiting the floatLiteral production.
	ExitFloatLiteral(c *FloatLiteralContext)

	// ExitIntegerLiteral is called when exiting the integerLiteral production.
	ExitIntegerLiteral(c *IntegerLiteralContext)

	// ExitNoneLiteral is called when exiting the noneLiteral production.
	ExitNoneLiteral(c *NoneLiteralContext)

	// ExitPropertyAssignment is called when exiting the propertyAssignment production.
	ExitPropertyAssignment(c *PropertyAssignmentContext)

	// ExitComputedPropertyName is called when exiting the computedPropertyName production.
	ExitComputedPropertyName(c *ComputedPropertyNameContext)

	// ExitPropertyName is called when exiting the propertyName production.
	ExitPropertyName(c *PropertyNameContext)

	// ExitNamespaceIdentifier is called when exiting the namespaceIdentifier production.
	ExitNamespaceIdentifier(c *NamespaceIdentifierContext)

	// ExitNamespace is called when exiting the namespace production.
	ExitNamespace(c *NamespaceContext)

	// ExitMemberExpression is called when exiting the memberExpression production.
	ExitMemberExpression(c *MemberExpressionContext)

	// ExitMemberExpressionSource is called when exiting the memberExpressionSource production.
	ExitMemberExpressionSource(c *MemberExpressionSourceContext)

	// ExitFunctionCallExpression is called when exiting the functionCallExpression production.
	ExitFunctionCallExpression(c *FunctionCallExpressionContext)

	// ExitFunctionCall is called when exiting the functionCall production.
	ExitFunctionCall(c *FunctionCallContext)

	// ExitFunctionName is called when exiting the functionName production.
	ExitFunctionName(c *FunctionNameContext)

	// ExitArgumentList is called when exiting the argumentList production.
	ExitArgumentList(c *ArgumentListContext)

	// ExitMemberExpressionPath is called when exiting the memberExpressionPath production.
	ExitMemberExpressionPath(c *MemberExpressionPathContext)

	// ExitSafeReservedWord is called when exiting the safeReservedWord production.
	ExitSafeReservedWord(c *SafeReservedWordContext)

	// ExitUnsafeReservedWord is called when exiting the unsafeReservedWord production.
	ExitUnsafeReservedWord(c *UnsafeReservedWordContext)

	// ExitRangeOperator is called when exiting the rangeOperator production.
	ExitRangeOperator(c *RangeOperatorContext)

	// ExitRangeOperand is called when exiting the rangeOperand production.
	ExitRangeOperand(c *RangeOperandContext)

	// ExitExpression is called when exiting the expression production.
	ExitExpression(c *ExpressionContext)

	// ExitPredicate is called when exiting the predicate production.
	ExitPredicate(c *PredicateContext)

	// ExitExpressionAtom is called when exiting the expressionAtom production.
	ExitExpressionAtom(c *ExpressionAtomContext)

	// ExitArrayOperator is called when exiting the arrayOperator production.
	ExitArrayOperator(c *ArrayOperatorContext)

	// ExitEqualityOperator is called when exiting the equalityOperator production.
	ExitEqualityOperator(c *EqualityOperatorContext)

	// ExitInOperator is called when exiting the inOperator production.
	ExitInOperator(c *InOperatorContext)

	// ExitLikeOperator is called when exiting the likeOperator production.
	ExitLikeOperator(c *LikeOperatorContext)

	// ExitUnaryOperator is called when exiting the unaryOperator production.
	ExitUnaryOperator(c *UnaryOperatorContext)

	// ExitRegexpOperator is called when exiting the regexpOperator production.
	ExitRegexpOperator(c *RegexpOperatorContext)

	// ExitLogicalAndOperator is called when exiting the logicalAndOperator production.
	ExitLogicalAndOperator(c *LogicalAndOperatorContext)

	// ExitLogicalOrOperator is called when exiting the logicalOrOperator production.
	ExitLogicalOrOperator(c *LogicalOrOperatorContext)

	// ExitMultiplicativeOperator is called when exiting the multiplicativeOperator production.
	ExitMultiplicativeOperator(c *MultiplicativeOperatorContext)

	// ExitAdditiveOperator is called when exiting the additiveOperator production.
	ExitAdditiveOperator(c *AdditiveOperatorContext)

	// ExitErrorOperator is called when exiting the errorOperator production.
	ExitErrorOperator(c *ErrorOperatorContext)
}
