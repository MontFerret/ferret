parser grammar FqlParser;

options { tokenVocab=FqlLexer; }

program
    : body
    ;

body
    : (bodyStatement)* bodyExpression
    ;

bodyStatement
    : functionCallExpression
    | variableDeclaration
    ;

bodyExpression
    : returnExpression
    | forExpression
    ;

returnExpression
    : Return (Distinct)? expression
    | Return (Distinct)? OpenParen forExpression CloseParen
    | Return forTernaryExpression
    ;

forExpression
    : For forExpressionValueVariable (Comma forExpressionKeyVariable)? In forExpressionSource
     (forExpressionClause)*
     (forExpressionBody)*
      forExpressionReturn
    ;

forExpressionValueVariable
    : Identifier
    ;

forExpressionKeyVariable
    : Identifier
    ;

forExpressionSource
    : functionCallExpression
    | arrayLiteral
    | objectLiteral
    | variable
    | memberExpression
    | rangeOperator
    | param
    ;

forExpressionClause
    : limitClause
    | sortClause
    | filterClause
    | collectClause
    ;

filterClause
    : Filter expression
    ;

limitClause
    : Limit IntegerLiteral (Comma IntegerLiteral)?
    ;

sortClause
    : Sort sortClauseExpression (Comma sortClauseExpression)*
    ;

sortClauseExpression
    : expression SortDirection?
    ;

collectClause
    : Collect collectVariable Assign expression
    | Collect collectVariable Assign expression Into collectGroupVariable
    | Collect collectVariable Assign expression Into collectGroupVariable Keep collectKeepVariable
    | Collect collectVariable Assign expression With Count collectCountVariable
    | Collect collectVariable Assign expression Aggregate collectAggregateVariable Assign collectAggregateExpression
    | Collect Aggregate collectAggregateVariable Assign collectAggregateExpression
    | Collect With Count Into collectCountVariable
    ;

collectVariable
    : Identifier
    ;

collectGroupVariable
    : Identifier
    ;

collectKeepVariable
    : Identifier
    ;

collectCountVariable
    : Identifier
    ;

collectAggregateVariable
    : Identifier
    ;

collectAggregateExpression
    : expression
    ;

collectOption
    :
    ;

forExpressionBody
    : variableDeclaration
    | functionCallExpression
    ;

forExpressionReturn
    : returnExpression
    | forExpression
    ;

variableDeclaration
    : Let Identifier Assign expression
    | Let Identifier Assign OpenParen forExpression CloseParen
    | Let Identifier Assign forTernaryExpression
    ;

param
    : Param Identifier
    ;

variable
    : Identifier
    ;

rangeOperator
    : (integerLiteral | variable | param) Range (integerLiteral | variable | param)
    ;

arrayLiteral
    : OpenBracket arrayElementList? CloseBracket
    ;

objectLiteral
    : OpenBrace (propertyAssignment (Comma propertyAssignment)*)? Comma? CloseBrace
    ;

booleanLiteral
    : BooleanLiteral
    ;

stringLiteral
    : StringLiteral
    ;

integerLiteral
    : IntegerLiteral
    ;

floatLiteral
    : FloatLiteral
    ;

noneLiteral
    : Null
    | None
    ;

arrayElementList
    : expression (Comma + expression)*
    ;

propertyAssignment
    : propertyName Colon expression
    | computedPropertyName Colon expression
    | shorthandPropertyName
    ;

memberExpression
    : Identifier (Dot propertyName (computedPropertyName)*)+
    | Identifier computedPropertyName (Dot propertyName (computedPropertyName)*)* (computedPropertyName (Dot propertyName)*)*
    ;

shorthandPropertyName
    : variable
    ;

computedPropertyName
    : OpenBracket expression CloseBracket
    ;

propertyName
    : Identifier
    ;

expressionSequence
    : expression (Comma expression)*
    ;

functionCallExpression
    : Identifier arguments
    ;

arguments
    : OpenParen(expression (Comma expression)*)?CloseParen
    ;

expression
    : expression equalityOperator expression
    | expression logicalOperator expression
    | expression mathOperator expression
    | functionCallExpression
    | OpenParen expressionSequence CloseParen
    | Plus expression
    | Minus expression
    | expression (All)? (Not)? In expression
    | Not expression
    | expression QuestionMark expression? Colon expression
    | rangeOperator
    | stringLiteral
    | integerLiteral
    | floatLiteral
    | booleanLiteral
    | arrayLiteral
    | objectLiteral
    | variable
    | memberExpression
    | noneLiteral
    | param
    ;

forTernaryExpression
    : expression QuestionMark expression? Colon OpenParen forExpression CloseParen
    | expression QuestionMark OpenParen forExpression CloseParen Colon expression
    | expression QuestionMark OpenParen forExpression CloseParen Colon OpenParen forExpression CloseParen
    ;

equalityOperator
    : Gt
    | Lt
    | Eq
    | Gte
    | Lte
    | Neq
    ;

logicalOperator
    : And
    | Or
    ;

mathOperator
    : Plus
    | Minus
    | Multi
    | Div
    | Mod
    ;

unaryOperator
    : Not
    | Plus
    | Minus
    | Like
    ;