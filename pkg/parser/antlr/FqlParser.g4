// $antlr-format off <-- used by VS Code Antlr extension

parser grammar FqlParser;

options { tokenVocab=FqlLexer; }

program
    : head* body
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

returnExpression
    : Return (Distinct)? expression
    ;

forExpression
    : For valueVariable=(Identifier | IgnoreIdentifier) (Comma counterVariable=Identifier)? In forExpressionSource
     forExpressionBody*
      forExpressionReturn
    | For counterVariable=(Identifier | IgnoreIdentifier) Do? While expression
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

waitForExpression
    : Waitfor Event waitForEventName In waitForEventSource (optionsClause)? (filterClause)? (timeoutClause)?
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

optionsClause
    : Options objectLiteral
    ;

timeoutClause
    : Timeout (integerLiteral | variable | param | memberExpression | functionCall)
    ;

variableDeclaration
    : Let Identifier Assign expression
    | Let IgnoreIdentifier Assign expression
    ;

param
    : Param Identifier
    ;

variable
    : Identifier
    ;

literal
    : arrayLiteral
    | objectLiteral
    | booleanLiteral
    | stringLiteral
    | floatLiteral
    | integerLiteral
    | noneLiteral
    ;

arrayLiteral
    : OpenBracket argumentList? CloseBracket
    ;

objectLiteral
    : OpenBrace (propertyAssignment (Comma propertyAssignment)* Comma?)? CloseBrace
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

namespaceIdentifier
    : namespace Identifier
    ;

namespace
    : NamespaceSegment*
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

functionCallExpression
    : functionCall errorOperator?
    ;

functionCall
    : namespace functionIdentifier OpenParen argumentList? CloseParen
    ;

argumentList
    : expression (Comma expression)* Comma?
    ;

memberExpressionPath
    : errorOperator? Dot propertyName
    | (errorOperator Dot)? computedPropertyName
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
    | Timeout
    ;

rangeOperator
    : left=rangeOperand Range right=rangeOperand
    ;

rangeOperand
    : integerLiteral
    | variable
    | param
    ;

expression
    : unaryOperator right=expression
    | condition=expression ternaryOperator=QuestionMark onTrue=expression? Colon onFalse=expression
    | left=expression logicalAndOperator right=expression
    | left=expression logicalOrOperator right=expression
    | predicate
    ;

predicate
    : left=predicate equalityOperator right=predicate
    | left=predicate arrayOperator right=predicate
    | left=predicate inOperator right=predicate
    | left=predicate likeOperator right=predicate
    | expressionAtom
    ;

expressionAtom
    : left=expressionAtom multiplicativeOperator right=expressionAtom
    | left=expressionAtom additiveOperator right=expressionAtom
    | left=expressionAtom regexpOperator right=expressionAtom
    | functionCallExpression
    | rangeOperator
    | literal
    | variable
    | memberExpression
    | param
    | OpenParen (forExpression | waitForExpression | expression) CloseParen errorOperator?
    ;

arrayOperator
    : operator=(All | Any | None) (inOperator | equalityOperator)
    ;

equalityOperator
    : Gt
    | Lt
    | Eq
    | Gte
    | Lte
    | Neq
    ;

inOperator
    : Not? In
    ;

likeOperator
    : Not? Like
    ;

unaryOperator
    : Not
    | Plus
    | Minus
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

errorOperator
    : QuestionMark
    ;