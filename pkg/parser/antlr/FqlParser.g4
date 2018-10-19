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
    : Collect collectCounter
    | Collect collectAggregator
    | Collect collectGrouping collectAggregator
    | Collect collectGrouping collectGroupVariable
    | Collect collectGrouping collectCounter
    | Collect collectGrouping
    ;

collectSelector
    : Identifier Assign expression
    ;

collectGrouping
    : collectSelector (Comma collectSelector)*
    ;

collectAggregator
    : Aggregate collectSelector (Comma collectSelector)*
    ;

collectGroupVariable
    : Into collectSelector
    | Into Identifier (Keep Identifier)?
    ;

collectCounter
    : With Count Into Identifier
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
    | TemplateStringLiteral
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
    : unaryOperator expression
    | expression equalityOperator expression
    | expression logicalOperator expression
    | expression mathOperator expression
    | functionCallExpression
    | OpenParen expressionSequence CloseParen
    | expression arrayOperator (inOperator | equalityOperator) expression
    | expression inOperator expression
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

arrayOperator
    : All
    | Any
    | None
    ;

inOperator
    : In
    | Not In
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