// $antlr-format off <-- used by VS Code Antlr extension

parser grammar FqlParser;

options { tokenVocab=FqlLexer; }

@parser::members {
	func (p *FqlParser) isWaitForPredicateStart() bool {
		la1 := p.GetTokenStream().LA(1)
		if la1 == FqlParserExists || la1 == FqlParserValue {
			return true
		}
		return la1 == FqlParserNot && p.GetTokenStream().LA(2) == FqlParserExists
	}
}

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
    : Use namespaceIdentifier As alias=Identifier
    ;

body
    : bodyStatement* bodyExpression
    ;

bodyStatement
    : variableDeclaration
    | functionCallExpression
    | waitForExpression
    | dispatchExpression
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
    | For valueVariable=(Identifier | IgnoreIdentifier) Assign stepInit=expression While stepCondition=expression Step stepVariable=(Identifier | IgnoreIdentifier) stepUpdate=(Increment | Decrement)
        forExpressionBody*
        forExpressionReturn
    | For valueVariable=(Identifier | IgnoreIdentifier) Assign stepInit=expression While stepCondition=expression Step stepVariable=(Identifier | IgnoreIdentifier) Assign stepUpdateExp=expression
        forExpressionBody*
        forExpressionReturn
    | For valueVariable=(Identifier | IgnoreIdentifier) Do? While expression
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
  : Collect collectGrouping collectAggregator collectGroupProjection
  | Collect collectGrouping collectAggregator
  | Collect collectGrouping collectGroupProjection
  | Collect collectGrouping collectCounter
  | Collect collectGrouping
  | Collect collectAggregator collectGroupProjection
  | Collect collectAggregator
  | Collect collectCounter
  ;


collectSelector
    : Identifier Assign expression
    ;

collectGrouping
    : collectSelector (Comma collectSelector)*
    ;

collectAggregateSelector
    : Identifier Assign functionCallExpression
    ;

collectAggregator
    : Aggregate collectAggregateSelector (Comma collectAggregateSelector)*
    ;

collectGroupProjection
    : Into collectSelector
    | Into Identifier (collectGroupProjectionFilter)?
    ;

collectGroupProjectionFilter
    : Keep Identifier (Comma Identifier)*
    ;

collectCounter
    : With Identifier Into Identifier
    ;

waitForExpression
    : Waitfor waitForEventExpression waitForOrThrowClause?
    | Waitfor waitForPredicateExpression waitForOrThrowClause?
    ;

dispatchExpression
    : Dispatch dispatchEventName In dispatchTarget (dispatchWithClause)? (dispatchOptionsClause)?
    ;

dispatchEventName
    : stringLiteral
    | variable
    | param
    | memberExpression
    | functionCall
    ;

dispatchTarget
    : functionCallExpression
    | variable
    | param
    | memberExpression
    ;

dispatchWithClause
    : With expression
    ;

dispatchOptionsClause
    : Options expression
    ;

waitForEventExpression
    : Event waitForEventName In waitForEventSource (optionsClause)? (filterClause)? (timeoutClause)?
    ;

waitForPredicateExpression
    : waitForPredicate (timeoutClause)? (everyClause)? (backoffClause)? (jitterClause)?
    ;

waitForPredicate
    : Exists expression
    | Not Exists expression
    | Value expression
    | {!p.isWaitForPredicateStart()}? expression
    ;

waitForEventName
    : stringLiteral
    | variable
    | param
    | memberExpression
    | functionCall
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
    : Timeout (durationLiteral | integerLiteral | floatLiteral | variable | param | memberExpression | functionCall)
    ;

everyClause
    : Every everyClauseValue (Comma everyClauseValue)?
    ;

everyClauseValue
    : durationLiteral
    | integerLiteral
    | floatLiteral
    | variable
    | param
    | memberExpression
    | functionCall
    ;

backoffClause
    : Backoff backoffStrategy
    ;

jitterClause
    : Jitter jitterClauseValue
    ;

jitterClauseValue
    : floatLiteral
    | integerLiteral
    | variable
    | param
    | memberExpression
    | functionCall
    ;

backoffStrategy
    : Identifier
    | stringLiteral
    | None
    ;

waitForOrThrowClause
    : Or Throw
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
    | templateLiteral
    ;

templateLiteral
    : BacktickOpen templateElement* BacktickClose
    ;

templateElement
    : TemplateChars
    | TemplateExprStart expression TemplateExprEnd
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
    | OpenParen (forExpression | waitForExpression | expression) CloseParen
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
    | arrayContraction
    | arrayExpansion
    | arrayQuestionMark
    | arrayApply
    ;

arrayExpansion
    : OpenBracket star=Multi inlineExpression? CloseBracket
    ;

arrayContraction
    : OpenBracket stars+=Multi stars+=Multi+ inlineExpression? CloseBracket
    ;

arrayQuestionMark
    : OpenBracket QuestionMark (Filter expression | arrayQuestionQuantifier Filter expression)? CloseBracket
    ;

arrayQuestionQuantifier
    : Any
    | All
    | None
    | At Least OpenParen arrayQuestionQuantifierValue CloseParen
    | arrayQuestionQuantifierValue Range arrayQuestionQuantifierValue
    | arrayQuestionQuantifierValue
    ;

arrayQuestionQuantifierValue
    : integerLiteral
    | param
    ;

arrayApply
    : OpenBracket Tilde queryLiteral CloseBracket
    ;

inlineExpression
    : inlineFilter inlineLimit? inlineReturn?
    | inlineLimit inlineReturn?
    | inlineReturn
    ;

inlineFilter
    : Filter expression
    ;

inlineLimit
    : Limit limitClauseValue (Comma limitClauseValue)?
    ;

inlineReturn
    : Return expression
    ;

safeReservedWord
    : And
    | Or
    | As
    | Distinct
    | Filter
    | Sort
    | Limit
    | Collect
    | SortDirection
    | Into
    | Keep
    | With
    | All
    | Any
    | At
    | Least
    | Aggregate
    | Event
    | Timeout
    | Options
    | Every
    | Backoff
    | Jitter
    | Exists
    | Value
    | Current
    | Step
    ;

unsafeReservedWord
    : Return
    | Dispatch
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
    | Throw
    ;

durationLiteral
    : DurationLiteral
    ;

rangeOperator
    : left=rangeOperand Range right=rangeOperand
    ;

rangeOperand
    : integerLiteral
    | variable
    | param
    | functionCallExpression
    | implicitMemberExpression
    | memberExpression
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
    | implicitMemberExpression
    | memberExpression
    | param
    | dispatchExpression
    | waitForExpression
    | OpenParen (forExpression | waitForExpression | expression) CloseParen errorOperator?
    ;

implicitMemberExpression
    : implicitMemberExpressionStart memberExpressionPath*
    ;

implicitMemberExpressionStart
    : errorOperator? Dot propertyName
    | errorOperator? Dot computedPropertyName
    ;

queryLiteral
    : Identifier (stringLiteral (OpenParen expression CloseParen)?)?
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
