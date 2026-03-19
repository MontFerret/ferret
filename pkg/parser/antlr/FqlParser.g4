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

	func (p *FqlParser) pushImplicitCurrent() {
		p.implicitCurrentDepth++
	}

	func (p *FqlParser) popImplicitCurrent() {
		if p.implicitCurrentDepth > 0 {
			p.implicitCurrentDepth--
		}
	}

	func (p *FqlParser) allowImplicitCurrent() bool {
		return p.implicitCurrentDepth > 0
	}

	func (p *FqlParser) isSafeReservedWordToken(token int) bool {
		switch token {
		case FqlParserAnd, FqlParserOr, FqlParserAs, FqlParserDistinct, FqlParserFilter, FqlParserSort,
			FqlParserLimit, FqlParserCollect, FqlParserSortDirection, FqlParserInto, FqlParserKeep, FqlParserWith,
			FqlParserAll, FqlParserAny, FqlParserAt, FqlParserLeast, FqlParserAggregate, FqlParserEvent, FqlParserTimeout,
			FqlParserOptions, FqlParserEvery, FqlParserBackoff, FqlParserJitter, FqlParserExists, FqlParserValue, FqlParserStep,
			FqlParserOne, FqlParserCount:
			return true
		default:
			return false
		}
	}

	func (p *FqlParser) isQueryModifierAhead() bool {
		la1 := p.GetTokenStream().LA(1)
		switch la1 {
		case FqlParserExists, FqlParserAny, FqlParserValue, FqlParserCount, FqlParserOne:
			return true
		case FqlParserIdentifier:
			tok := p.GetTokenStream().LT(1)
			if tok == nil {
				return false
			}

			return p.isQueryModifierText(tok.GetText())
		default:
			return false
		}
	}

	func (p *FqlParser) isQueryModifierText(text string) bool {
		return equalsFoldAscii(text, "EXISTS") ||
			equalsFoldAscii(text, "COUNT") ||
			equalsFoldAscii(text, "ANY") ||
			equalsFoldAscii(text, "VALUE") ||
			equalsFoldAscii(text, "ONE")
	}

	func (p *FqlParser) isCurrentIdentifierText(expected string) bool {
		tok := p.GetTokenStream().LT(1)
		if tok == nil || tok.GetTokenType() != FqlParserIdentifier {
			return false
		}

		return equalsFoldAscii(tok.GetText(), expected)
	}

	func equalsFoldAscii(actual, expected string) bool {
		if len(actual) != len(expected) {
			return false
		}

		for i := 0; i < len(actual); i++ {
			a := actual[i]
			e := expected[i]

			if a >= 'a' && a <= 'z' {
				a -= 'a' - 'A'
			}

			if e >= 'a' && e <= 'z' {
				e -= 'a' - 'A'
			}

			if a != e {
				return false
			}
		}

		return true
	}

	func (p *FqlParser) isUnsafeReservedWordToken(token int) bool {
		switch token {
		case FqlParserReturn, FqlParserDispatch, FqlParserQuery, FqlParserUsing, FqlParserNone,
			FqlParserNull, FqlParserLet, FqlParserVar, FqlParserUse, FqlParserWaitfor, FqlParserWhile, FqlParserDo, FqlParserIn,
			FqlParserLike, FqlParserNot, FqlParserFor, FqlParserBooleanLiteral, FqlParserThrow, FqlParserMatch, FqlParserWhen,
			FqlParserFunc:
			return true
		default:
			return false
		}
	}

	func (p *FqlParser) isImplicitCurrentValue() bool {
		la1 := p.GetTokenStream().LA(2)
		switch la1 {
		case FqlParserIdentifier, FqlParserStringLiteral, FqlParserParam, FqlParserOpenBracket, FqlParserBacktickOpen:
			return false
		}

		if p.isSafeReservedWordToken(la1) || p.isUnsafeReservedWordToken(la1) {
			return false
		}

		return true
	}
}

@parser::structmembers {
	implicitCurrentDepth int
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
    | assignmentStatement
    | functionDeclaration
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
    | Var bindingIdentifier Assign expression
    ;

assignmentStatement
    : assignmentTarget Assign expression
    ;

assignmentTarget
    : bindingIdentifier
    | memberExpression
    ;

functionDeclaration
    : Func functionName OpenParen functionParameterList? CloseParen functionBody
    ;

functionParameterList
    : functionParameter (Comma functionParameter)* Comma?
    ;

functionParameter
    : Identifier
    ;

functionBody
    : functionArrow
    | functionBlock
    ;

functionArrow
    : Arrow expression
    ;

functionBlock
    : OpenParen functionStatement* functionReturn CloseParen
    ;

functionStatement
    : variableDeclaration
    | assignmentStatement
    | functionDeclaration
    | functionCallExpression
    | waitForExpression
    | dispatchExpression
    | expressionStatement
    ;

expressionStatement
    : expression
    ;

functionReturn
    : Return expression
    ;

returnExpression
    : Return (Distinct)? expression
    ;

forExpression
    : For valueVariable=loopVariable (Comma counterVariable=bindingIdentifier)? In forExpressionSource
        forExpressionBody*
        forExpressionReturn
    | For valueVariable=loopVariable Do? While expression
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
    | assignmentStatement
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

eventFilterClause
    : When {p.pushImplicitCurrent()} expression {p.popImplicitCurrent()}
    ;

limitClause
    : Limit limitClauseValue (Comma limitClauseValue)?
    ;

limitClauseValue
    : integerLiteral
    | param
    | variable
    | functionCallExpression
    | implicitCurrentExpression
    | {p.allowImplicitCurrent()}? implicitMemberExpression
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

bindingIdentifier
    : Identifier
    | safeReservedWord
    ;

loopVariable
    : bindingIdentifier
    | IgnoreIdentifier
    ;

collectSelector
    : bindingIdentifier Assign expression
    ;

collectGrouping
    : collectSelector (Comma collectSelector)*
    ;

collectAggregateSelector
    : bindingIdentifier Assign functionCallExpression
    ;

collectAggregator
    : Aggregate collectAggregateSelector (Comma collectAggregateSelector)*
    ;

collectGroupProjection
    : Into collectSelector
    | Into bindingIdentifier (collectGroupProjectionFilter)?
    ;

collectGroupProjectionFilter
    : Keep Identifier (Comma Identifier)*
    ;

collectCounter
    : With Count Into bindingIdentifier
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
    : Event waitForEventName In waitForEventSource (optionsClause)? (eventFilterClause)? (timeoutClause)?
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
    : OpenBracket QuestionMark (Filter {p.pushImplicitCurrent()} expression {p.popImplicitCurrent()} | arrayQuestionQuantifier Filter {p.pushImplicitCurrent()} expression {p.popImplicitCurrent()})? CloseBracket
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
    : Filter {p.pushImplicitCurrent()} expression {p.popImplicitCurrent()}
    ;

inlineLimit
    : Limit {p.pushImplicitCurrent()} limitClauseValue (Comma limitClauseValue)? {p.popImplicitCurrent()}
    ;

inlineReturn
    : Return {p.pushImplicitCurrent()} expression {p.popImplicitCurrent()}
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
    | Count
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
    | Step
    | One
    ;

unsafeReservedWord
    : Return
    | Dispatch
    | Query
    | Using
    | None
    | Null
    | Let
    | Var
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
    | Match
    | When
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
    | implicitCurrentExpression
    | {p.allowImplicitCurrent()}? implicitMemberExpression
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
    | matchExpression
    | queryExpression
    | functionCallExpression
    | rangeOperator
    | literal
    | variable
    | implicitCurrentExpression
    | {p.allowImplicitCurrent()}? implicitMemberExpression
    | memberExpression
    | param
    | dispatchExpression
    | waitForExpression
    | OpenParen (forExpression | waitForExpression | expression) CloseParen errorOperator?
    ;

matchExpression
    : Match expression matchPatternArms
    | Match matchGuardArms
    ;

matchPatternArms
    : OpenParen matchPatternArmList? matchDefaultArm Comma? CloseParen
    ;

matchPatternArmList
    : matchPatternArm (Comma matchPatternArm)* Comma
    ;

matchGuardArms
    : OpenParen matchGuardArmList? matchDefaultArm Comma? CloseParen
    ;

matchGuardArmList
    : matchGuardArm (Comma matchGuardArm)* Comma
    ;

matchPatternArm
    : matchPattern matchPatternGuard? Arrow expression
    ;

matchPatternGuard
    : When expression
    ;

matchGuardArm
    : When expression Arrow expression
    ;

matchDefaultArm
    : IgnoreIdentifier Arrow expression
    ;

matchPattern
    : matchLiteralPattern
    | matchBindingPattern
    | matchObjectPattern
    ;

matchBindingPattern
    : Identifier
    | safeReservedWord
    ;

matchLiteralPattern
    : noneLiteral
    | booleanLiteral
    | stringLiteral
    | floatLiteral
    | integerLiteral
    ;

matchObjectPattern
    : OpenBrace (matchObjectPatternProperty (Comma matchObjectPatternProperty)* Comma?)? CloseBrace
    ;

matchObjectPatternProperty
    : matchObjectPatternKey Colon matchPattern
    ;

matchObjectPatternKey
    : Identifier
    | stringLiteral
    | safeReservedWord
    | unsafeReservedWord
    ;

implicitMemberExpression
    : implicitMemberExpressionStart memberExpressionPath*
    ;

implicitCurrentExpression
    : {p.allowImplicitCurrent() && p.isImplicitCurrentValue()}? Dot
    ;

implicitMemberExpressionStart
    : errorOperator? Dot propertyName
    | errorOperator? Dot computedPropertyName
    | Dot arrayExpansion
    | Dot arrayContraction
    | Dot arrayQuestionMark
    | Dot arrayApply
    ;

queryExpression
    : Query queryModifier queryPayload In expression Using dialect=Identifier queryWithOpt?
    | Query {!p.isQueryModifierAhead()}? queryPayload In expression Using dialect=Identifier queryWithOpt?
    ;

queryModifier
    : Exists
    | Any
    | Value
    | Count
    | One
    ;

queryPayload
    : stringLiteral
    | param
    | variable
    ;

queryWithOpt
    : With expression
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
