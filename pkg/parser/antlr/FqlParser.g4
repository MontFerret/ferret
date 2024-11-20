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

variableDeclaration
    : Let id=(Identifier | IgnoreIdentifier) Assign expression
    | Let safeReservedWord Assign expression
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
    : integerLiteral
    | param
    | variable
    | functionCallExpression
    | memberExpression
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

param
    : Param Identifier
    | Param safeReservedWord
    ;

variable
    : Identifier
    | safeReservedWord
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
    | safeReservedWord
    | unsafeReservedWord
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
    : namespace functionName OpenParen argumentList? CloseParen
    ;

functionName
    : Identifier
    | safeReservedWord
    | unsafeReservedWord
    ;

argumentList
    : expression (Comma expression)* Comma?
    ;

memberExpressionPath
    : errorOperator? Dot propertyName
    | (errorOperator Dot)? computedPropertyName
    ;

safeReservedWord
    : And
    | Or
    | Distinct
    | Filter
    | Sort
    | Limit
    | Collect
    | SortDirection
    | Into
    | Keep
    | With
    | Count
    | All
    | Any
    | Aggregate
    | Event
    | Timeout
    | Options
    | Current
    ;

unsafeReservedWord
    : Return
    | None
    | Null
    | Let
    | Use
    | Waitfor
    | While
    | Do
    | In
    | Like
    | Not
    | For
    | BooleanLiteral
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
    | left=expression logicalAndOperator right=expression
    | left=expression logicalOrOperator right=expression
    | condition=expression ternaryOperator=QuestionMark onTrue=expression? Colon onFalse=expression
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