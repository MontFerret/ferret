// $antlr-format off <-- used by VS Code Antlr extension

parser grammar FqlParser;

options { tokenVocab=FqlLexer; }

program
    : head* body EOF
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

namespaceIdentifier
    : namespace Identifier
    ;

body
    : bodyStatement* bodyExpression
    ;

bodyStatement
    : variableDeclaration
    | functionCallExpression
    | waitForExpression
    ;

bodyExpression
    : returnExpression
    | forExpression
    ;

variableDeclaration
    : Let Identifier Assign expression
    | Let IgnoreIdentifier Assign expression
    ;

returnExpression
    : Return Distinct? expression
    ;

inlineHighLevelExpression
    : OpenParen highLevelExpression CloseParen errorOperator?
    ;

highLevelExpression
    : forExpression
    | waitForExpression
    ;

forExpression
    : For Identifier (Comma Identifier)? In forExpressionSource
     forExpressionBody*
      forExpressionReturn
    | For Identifier Do? While expression
     forExpressionBody*
      forExpressionReturn
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

optionsClause
    : Options objectLiteral
    | Options variable
    ;

waitForExpression
    : Waitfor Event waitForEventName In waitForEventSource (optionsClause)? (waitForTimeout)?
    ;

waitForTimeout
    : integerLiteral
    | variable
    | param
    ;

waitForEventName
    : stringLiteral
    | variable
    | param
    | functionCallExpression
    | memberExpression
    ;

waitForEventSource
    : functionCallExpression
    | variable
    | memberExpression
    ;

rangeOperator
    : (integerLiteral | variable | param) Range (integerLiteral | variable | param)
    ;

arrayLiteral
    : OpenBracket (expression (Comma expression)* Comma?)? CloseBracket
    ;

objectLiteral
    : OpenBrace (propertyAssignment (Comma propertyAssignment)* Comma?)? CloseBrace
    ;

propertyAssignment
    : propertyName Colon expression
    | computedPropertyName Colon expression
    | variable
    ;

computedPropertyName
    : OpenBracket expression CloseBracket
    ;

propertyName
    : Identifier
    | stringLiteral
    | param
    ;

booleanLiteral
    : BooleanLiteral
    ;

stringLiteral
    : StringLiteral
    ;

floatLiteral
    : FloatLiteral
    ;

integerLiteral
    : IntegerLiteral
    ;

noneLiteral
    : Null
    | None
    ;

expression
    : unaryOperator expression
    | expression multiplicativeOperator expression
    | expression additiveOperator expression
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
    | floatLiteral
    | integerLiteral
    | booleanLiteral
    | arrayLiteral
    | objectLiteral
    | memberExpression
    | functionCallExpression
    | param
    | variable
    | noneLiteral
    | expressionGroup
    | inlineHighLevelExpression
    ;

expressionGroup
    : OpenParen expression CloseParen errorOperator?
    ;

memberExpression
    : memberExpressionSource memberExpressionPath+
    ;

memberExpressionSource
    : variable
    | param
    | arrayLiteral
    | objectLiteral
    | functionCall
    ;

functionCall
    : namespace functionIdentifier arguments
    ;

functionCallExpression
    : functionCall errorOperator?
    ;

memberExpressionPath
    : QuestionMark? Dot propertyName
    | (QuestionMark Dot)? computedPropertyName
    ;

errorOperator
    : QuestionMark
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

namespace
    : NamespaceSegment*
    ;

arguments
    : OpenParen (expression (Comma expression)*)? CloseParen
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

param
    : Param Identifier
    ;

variable
    : Identifier
    ;