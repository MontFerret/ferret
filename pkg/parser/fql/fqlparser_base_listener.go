// Code generated from antlr/FqlParser.g4 by ANTLR 4.7.1. DO NOT EDIT.

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

// EnterReturnExpression is called when production returnExpression is entered.
func (s *BaseFqlParserListener) EnterReturnExpression(ctx *ReturnExpressionContext) {}

// ExitReturnExpression is called when production returnExpression is exited.
func (s *BaseFqlParserListener) ExitReturnExpression(ctx *ReturnExpressionContext) {}

// EnterForExpression is called when production forExpression is entered.
func (s *BaseFqlParserListener) EnterForExpression(ctx *ForExpressionContext) {}

// ExitForExpression is called when production forExpression is exited.
func (s *BaseFqlParserListener) ExitForExpression(ctx *ForExpressionContext) {}

// EnterForExpressionValueVariable is called when production forExpressionValueVariable is entered.
func (s *BaseFqlParserListener) EnterForExpressionValueVariable(ctx *ForExpressionValueVariableContext) {
}

// ExitForExpressionValueVariable is called when production forExpressionValueVariable is exited.
func (s *BaseFqlParserListener) ExitForExpressionValueVariable(ctx *ForExpressionValueVariableContext) {
}

// EnterForExpressionKeyVariable is called when production forExpressionKeyVariable is entered.
func (s *BaseFqlParserListener) EnterForExpressionKeyVariable(ctx *ForExpressionKeyVariableContext) {}

// ExitForExpressionKeyVariable is called when production forExpressionKeyVariable is exited.
func (s *BaseFqlParserListener) ExitForExpressionKeyVariable(ctx *ForExpressionKeyVariableContext) {}

// EnterForExpressionSource is called when production forExpressionSource is entered.
func (s *BaseFqlParserListener) EnterForExpressionSource(ctx *ForExpressionSourceContext) {}

// ExitForExpressionSource is called when production forExpressionSource is exited.
func (s *BaseFqlParserListener) ExitForExpressionSource(ctx *ForExpressionSourceContext) {}

// EnterForExpressionClause is called when production forExpressionClause is entered.
func (s *BaseFqlParserListener) EnterForExpressionClause(ctx *ForExpressionClauseContext) {}

// ExitForExpressionClause is called when production forExpressionClause is exited.
func (s *BaseFqlParserListener) ExitForExpressionClause(ctx *ForExpressionClauseContext) {}

// EnterFilterClause is called when production filterClause is entered.
func (s *BaseFqlParserListener) EnterFilterClause(ctx *FilterClauseContext) {}

// ExitFilterClause is called when production filterClause is exited.
func (s *BaseFqlParserListener) ExitFilterClause(ctx *FilterClauseContext) {}

// EnterLimitClause is called when production limitClause is entered.
func (s *BaseFqlParserListener) EnterLimitClause(ctx *LimitClauseContext) {}

// ExitLimitClause is called when production limitClause is exited.
func (s *BaseFqlParserListener) ExitLimitClause(ctx *LimitClauseContext) {}

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

// EnterForExpressionBody is called when production forExpressionBody is entered.
func (s *BaseFqlParserListener) EnterForExpressionBody(ctx *ForExpressionBodyContext) {}

// ExitForExpressionBody is called when production forExpressionBody is exited.
func (s *BaseFqlParserListener) ExitForExpressionBody(ctx *ForExpressionBodyContext) {}

// EnterForExpressionReturn is called when production forExpressionReturn is entered.
func (s *BaseFqlParserListener) EnterForExpressionReturn(ctx *ForExpressionReturnContext) {}

// ExitForExpressionReturn is called when production forExpressionReturn is exited.
func (s *BaseFqlParserListener) ExitForExpressionReturn(ctx *ForExpressionReturnContext) {}

// EnterVariableDeclaration is called when production variableDeclaration is entered.
func (s *BaseFqlParserListener) EnterVariableDeclaration(ctx *VariableDeclarationContext) {}

// ExitVariableDeclaration is called when production variableDeclaration is exited.
func (s *BaseFqlParserListener) ExitVariableDeclaration(ctx *VariableDeclarationContext) {}

// EnterParam is called when production param is entered.
func (s *BaseFqlParserListener) EnterParam(ctx *ParamContext) {}

// ExitParam is called when production param is exited.
func (s *BaseFqlParserListener) ExitParam(ctx *ParamContext) {}

// EnterVariable is called when production variable is entered.
func (s *BaseFqlParserListener) EnterVariable(ctx *VariableContext) {}

// ExitVariable is called when production variable is exited.
func (s *BaseFqlParserListener) ExitVariable(ctx *VariableContext) {}

// EnterRangeOperator is called when production rangeOperator is entered.
func (s *BaseFqlParserListener) EnterRangeOperator(ctx *RangeOperatorContext) {}

// ExitRangeOperator is called when production rangeOperator is exited.
func (s *BaseFqlParserListener) ExitRangeOperator(ctx *RangeOperatorContext) {}

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

// EnterIntegerLiteral is called when production integerLiteral is entered.
func (s *BaseFqlParserListener) EnterIntegerLiteral(ctx *IntegerLiteralContext) {}

// ExitIntegerLiteral is called when production integerLiteral is exited.
func (s *BaseFqlParserListener) ExitIntegerLiteral(ctx *IntegerLiteralContext) {}

// EnterFloatLiteral is called when production floatLiteral is entered.
func (s *BaseFqlParserListener) EnterFloatLiteral(ctx *FloatLiteralContext) {}

// ExitFloatLiteral is called when production floatLiteral is exited.
func (s *BaseFqlParserListener) ExitFloatLiteral(ctx *FloatLiteralContext) {}

// EnterNoneLiteral is called when production noneLiteral is entered.
func (s *BaseFqlParserListener) EnterNoneLiteral(ctx *NoneLiteralContext) {}

// ExitNoneLiteral is called when production noneLiteral is exited.
func (s *BaseFqlParserListener) ExitNoneLiteral(ctx *NoneLiteralContext) {}

// EnterArrayElementList is called when production arrayElementList is entered.
func (s *BaseFqlParserListener) EnterArrayElementList(ctx *ArrayElementListContext) {}

// ExitArrayElementList is called when production arrayElementList is exited.
func (s *BaseFqlParserListener) ExitArrayElementList(ctx *ArrayElementListContext) {}

// EnterPropertyAssignment is called when production propertyAssignment is entered.
func (s *BaseFqlParserListener) EnterPropertyAssignment(ctx *PropertyAssignmentContext) {}

// ExitPropertyAssignment is called when production propertyAssignment is exited.
func (s *BaseFqlParserListener) ExitPropertyAssignment(ctx *PropertyAssignmentContext) {}

// EnterMemberExpression is called when production memberExpression is entered.
func (s *BaseFqlParserListener) EnterMemberExpression(ctx *MemberExpressionContext) {}

// ExitMemberExpression is called when production memberExpression is exited.
func (s *BaseFqlParserListener) ExitMemberExpression(ctx *MemberExpressionContext) {}

// EnterShorthandPropertyName is called when production shorthandPropertyName is entered.
func (s *BaseFqlParserListener) EnterShorthandPropertyName(ctx *ShorthandPropertyNameContext) {}

// ExitShorthandPropertyName is called when production shorthandPropertyName is exited.
func (s *BaseFqlParserListener) ExitShorthandPropertyName(ctx *ShorthandPropertyNameContext) {}

// EnterComputedPropertyName is called when production computedPropertyName is entered.
func (s *BaseFqlParserListener) EnterComputedPropertyName(ctx *ComputedPropertyNameContext) {}

// ExitComputedPropertyName is called when production computedPropertyName is exited.
func (s *BaseFqlParserListener) ExitComputedPropertyName(ctx *ComputedPropertyNameContext) {}

// EnterPropertyName is called when production propertyName is entered.
func (s *BaseFqlParserListener) EnterPropertyName(ctx *PropertyNameContext) {}

// ExitPropertyName is called when production propertyName is exited.
func (s *BaseFqlParserListener) ExitPropertyName(ctx *PropertyNameContext) {}

// EnterExpressionSequence is called when production expressionSequence is entered.
func (s *BaseFqlParserListener) EnterExpressionSequence(ctx *ExpressionSequenceContext) {}

// ExitExpressionSequence is called when production expressionSequence is exited.
func (s *BaseFqlParserListener) ExitExpressionSequence(ctx *ExpressionSequenceContext) {}

// EnterFunctionCallExpression is called when production functionCallExpression is entered.
func (s *BaseFqlParserListener) EnterFunctionCallExpression(ctx *FunctionCallExpressionContext) {}

// ExitFunctionCallExpression is called when production functionCallExpression is exited.
func (s *BaseFqlParserListener) ExitFunctionCallExpression(ctx *FunctionCallExpressionContext) {}

// EnterArguments is called when production arguments is entered.
func (s *BaseFqlParserListener) EnterArguments(ctx *ArgumentsContext) {}

// ExitArguments is called when production arguments is exited.
func (s *BaseFqlParserListener) ExitArguments(ctx *ArgumentsContext) {}

// EnterExpression is called when production expression is entered.
func (s *BaseFqlParserListener) EnterExpression(ctx *ExpressionContext) {}

// ExitExpression is called when production expression is exited.
func (s *BaseFqlParserListener) ExitExpression(ctx *ExpressionContext) {}

// EnterForTernaryExpression is called when production forTernaryExpression is entered.
func (s *BaseFqlParserListener) EnterForTernaryExpression(ctx *ForTernaryExpressionContext) {}

// ExitForTernaryExpression is called when production forTernaryExpression is exited.
func (s *BaseFqlParserListener) ExitForTernaryExpression(ctx *ForTernaryExpressionContext) {}

// EnterArrayOperator is called when production arrayOperator is entered.
func (s *BaseFqlParserListener) EnterArrayOperator(ctx *ArrayOperatorContext) {}

// ExitArrayOperator is called when production arrayOperator is exited.
func (s *BaseFqlParserListener) ExitArrayOperator(ctx *ArrayOperatorContext) {}

// EnterInOperator is called when production inOperator is entered.
func (s *BaseFqlParserListener) EnterInOperator(ctx *InOperatorContext) {}

// ExitInOperator is called when production inOperator is exited.
func (s *BaseFqlParserListener) ExitInOperator(ctx *InOperatorContext) {}

// EnterEqualityOperator is called when production equalityOperator is entered.
func (s *BaseFqlParserListener) EnterEqualityOperator(ctx *EqualityOperatorContext) {}

// ExitEqualityOperator is called when production equalityOperator is exited.
func (s *BaseFqlParserListener) ExitEqualityOperator(ctx *EqualityOperatorContext) {}

// EnterLogicalOperator is called when production logicalOperator is entered.
func (s *BaseFqlParserListener) EnterLogicalOperator(ctx *LogicalOperatorContext) {}

// ExitLogicalOperator is called when production logicalOperator is exited.
func (s *BaseFqlParserListener) ExitLogicalOperator(ctx *LogicalOperatorContext) {}

// EnterMathOperator is called when production mathOperator is entered.
func (s *BaseFqlParserListener) EnterMathOperator(ctx *MathOperatorContext) {}

// ExitMathOperator is called when production mathOperator is exited.
func (s *BaseFqlParserListener) ExitMathOperator(ctx *MathOperatorContext) {}

// EnterUnaryOperator is called when production unaryOperator is entered.
func (s *BaseFqlParserListener) EnterUnaryOperator(ctx *UnaryOperatorContext) {}

// ExitUnaryOperator is called when production unaryOperator is exited.
func (s *BaseFqlParserListener) ExitUnaryOperator(ctx *UnaryOperatorContext) {}
