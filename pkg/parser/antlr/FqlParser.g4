parser grammar FqlParser;

options { tokenVocab=FqlLexer; }

program
    : (head)* body
    ;

head
    : useExpression
    ;

useExpression
    : use
    ;

use
    : Use namespaceIdentifier
    ;

body
    : (bodyStatement)* bodyExpression
    ;

bodyStatement
    : functionCallExpression
    | variableDeclaration
    | waitForExpression
    ;

bodyExpression
    : returnExpression
    | forExpression
    ;

returnExpression
    : Return (Distinct)? OpenParen forExpression CloseParen
    | Return ternaryExpression
    | Return OpenParen waitForExpression CloseParen
    | Return (Distinct)? expression
    ;

forExpression
    : For forExpressionValueVariable (Comma forExpressionKeyVariable)? In forExpressionSource
     (forExpressionBody)*
      forExpressionReturn
    | For forExpressionValueVariable (Do)? While expression
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

forExpressionStatement
    : variableDeclaration
    | functionCallExpression
    ;

forExpressionBody
    : forExpressionStatement
    | forExpressionClause
    ;

forExpressionReturn
    : returnExpression
    | forExpression
    ;

filterClause
    : Filter expression
    ;

limitClause
    : Limit limitClauseValue (Comma limitClauseValue)?
    ;

limitClauseValue
    : IntegerLiteral
    | param
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
    : Aggregate collectAggregateSelector (Comma collectAggregateSelector)*
    ;

collectAggregateSelector
    : Identifier Assign functionCallExpression
    ;

collectGroupVariable
    : Into collectSelector
    | Into Identifier (Keep Identifier)?
    ;

collectCounter
    : With Count Into Identifier
    ;


waitForExpression
    : waitForEventExpression
    ;

waitForTimeout
    : integerLiteral
    | variable
    | functionCallExpression
    | memberExpression
    ;

waitForOptions
    : objectLiteral
    ;

waitForEventName
    : stringLiteral
    | variable
    | param
    | functionCallExpression
    | memberExpression
    ;

waitForEventSource
    : variable
    | functionCallExpression
    | memberExpression
    ;

waitForEventExpression
    : Waitfor Event waitForEventName In waitForEventSource (Options waitForOptions)? (waitForTimeout)?
    ;

variableDeclaration
    : Let Identifier Assign OpenParen forExpression CloseParen
    | Let Identifier Assign ternaryExpression
    | Let Identifier Assign OpenParen waitForExpression CloseParen
    | Let Identifier Assign expression
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

shorthandPropertyName
    : variable
    ;

computedPropertyName
    : OpenBracket expression CloseBracket
    ;

propertyName
    : Identifier
    | stringLiteral
    | param
    ;

expressionGroup
    : OpenParen expression CloseParen
    ;

namespaceIdentifier
    : namespace Identifier
    ;

namespace
    : (NamespaceSegment)*
    ;

functionIdentifier
    : Identifier
    | And
    | Or
    | For
    | Return
    | Distinct
    | Filter
    | Sort
    | Limit
    | Let
    | Collect
    | SortDirection
    | None
    | Null
    | BooleanLiteral
    | Use
    | Into
    | Keep
    | With
    | Count
    | All
    | Any
    | Aggregate
    | Like
    | Not
    | In
    | Waitfor
    | Event
    ;

functionCallExpression
    : namespace functionIdentifier arguments
    ;

member
    : Identifier
    | functionCallExpression
    | param
    ;

memberPath
    : (Dot propertyName (computedPropertyName)*)+
    | computedPropertyName (Dot propertyName (computedPropertyName)*)* (computedPropertyName (Dot propertyName)*)*
    ;

memberExpression
    : member memberPath
    ;

arguments
    : OpenParen(expression (Comma expression)*)?CloseParen
    ;

expression
    : unaryOperator expression
    | expression multiplicativeOperator expression
    | expression additiveOperator expression
    | functionCallExpression
    | expressionGroup
    | expression arrayOperator (inOperator | equalityOperator) expression
    | expression inOperator expression
    | expression likeOperator expression
    | expression equalityOperator expression
    | expression regexpOperator expression
    | expression logicalAndOperator expression
    | expression logicalOrOperator expression
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

ternaryExpression
    : expression QuestionMark OpenParen forExpression CloseParen Colon OpenParen forExpression CloseParen
    | expression QuestionMark OpenParen forExpression CloseParen Colon OpenParen waitForExpression CloseParen
    | expression QuestionMark OpenParen waitForExpression CloseParen Colon OpenParen waitForExpression CloseParen
    | expression QuestionMark OpenParen waitForExpression CloseParen Colon OpenParen forExpression CloseParen
    | expression QuestionMark OpenParen waitForExpression CloseParen Colon expression
    | expression QuestionMark OpenParen forExpression CloseParen Colon expression
    | expression QuestionMark expression? Colon OpenParen forExpression CloseParen
    | expression QuestionMark expression? Colon OpenParen waitForExpression CloseParen
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

likeOperator
    : Like
    | Not Like
    ;

equalityOperator
    : Gt
    | Lt
    | Eq
    | Gte
    | Lte
    | Neq
    ;

regexpOperator
    : RegexMatch
    | RegexNotMatch
    ;

logicalAndOperator
    : And
    ;

logicalOrOperator
    : Or
    ;

multiplicativeOperator
    : Multi
    | Div
    | Mod
    ;

additiveOperator
    : Plus
    | Minus
    ;

unaryOperator
    : Not
    | Plus
    | Minus
    ;