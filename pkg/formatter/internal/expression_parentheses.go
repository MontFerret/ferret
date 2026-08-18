package internal

import (
	"github.com/antlr4-go/antlr/v4"

	"github.com/MontFerret/ferret/v2/pkg/parser/fql"
)

type (
	expressionPrecedence uint8

	expressionAssociativity uint8

	expressionOperandSide uint8

	expressionOperation struct {
		precedence    expressionPrecedence
		associativity expressionAssociativity
		side          expressionOperandSide
	}
)

const (
	precedenceNone expressionPrecedence = iota
	precedenceTernary
	precedenceCoalesce
	precedenceLogicalOr
	precedenceLogicalAnd
	precedenceUnary
	precedenceLike
	precedenceIn
	precedenceArray
	precedenceEquality
	precedenceRegexp
	precedenceAdditive
	precedenceMultiplicative
	precedencePrimary
)

const (
	associativityNone expressionAssociativity = iota
	associativityLeft
	associativityRight
)

const (
	operandSideNone expressionOperandSide = iota
	operandSideLeft
	operandSideRight
	operandSideMiddle
)

func expressionRootPrecedence(ctx *fql.ExpressionContext) expressionPrecedence {
	if ctx == nil {
		return precedenceNone
	}

	switch {
	case ctx.GetTernaryOperator() != nil:
		return precedenceTernary
	case ctx.GetCoalesceOperator() != nil:
		return precedenceCoalesce
	case ctx.LogicalOrOperator() != nil:
		return precedenceLogicalOr
	case ctx.LogicalAndOperator() != nil:
		return precedenceLogicalAnd
	case ctx.UnaryOperator() != nil:
		return precedenceUnary
	case ctx.Predicate() != nil:
		return predicateRootPrecedence(ctx.Predicate().(*fql.PredicateContext))
	default:
		return precedenceNone
	}
}

func predicateRootPrecedence(ctx *fql.PredicateContext) expressionPrecedence {
	if ctx == nil {
		return precedenceNone
	}

	switch {
	case ctx.EqualityOperator() != nil:
		return precedenceEquality
	case ctx.ArrayOperator() != nil:
		return precedenceArray
	case ctx.InOperator() != nil:
		return precedenceIn
	case ctx.LikeOperator() != nil:
		return precedenceLike
	case ctx.ExpressionAtom() != nil:
		return expressionAtomRootPrecedence(ctx.ExpressionAtom().(*fql.ExpressionAtomContext))
	default:
		return precedenceNone
	}
}

func expressionAtomRootPrecedence(ctx *fql.ExpressionAtomContext) expressionPrecedence {
	if ctx == nil {
		return precedenceNone
	}

	switch {
	case ctx.MultiplicativeOperator() != nil:
		return precedenceMultiplicative
	case ctx.AdditiveOperator() != nil:
		return precedenceAdditive
	case ctx.RegexpOperator() != nil:
		return precedenceRegexp
	default:
		return precedencePrimary
	}
}

func canRemoveExpressionParentheses(inner *fql.ExpressionContext, outer expressionOperation) bool {
	if inner == nil || outer.precedence == precedenceNone {
		return inner != nil
	}

	innerPrecedence := expressionRootPrecedence(inner)
	if innerPrecedence > outer.precedence {
		return true
	}

	if innerPrecedence < outer.precedence {
		return false
	}

	if innerPrecedence == precedenceUnary {
		return false
	}

	if innerPrecedence == precedenceTernary {
		return outer.side == operandSideLeft || outer.side == operandSideMiddle
	}

	switch outer.associativity {
	case associativityLeft:
		return outer.side == operandSideLeft
	case associativityRight:
		return outer.side == operandSideRight
	default:
		return false
	}
}

func binaryExpressionOperation(precedence expressionPrecedence, side expressionOperandSide) expressionOperation {
	return expressionOperation{
		precedence:    precedence,
		associativity: associativityLeft,
		side:          side,
	}
}

func coalesceExpressionOperation(side expressionOperandSide) expressionOperation {
	return expressionOperation{
		precedence:    precedenceCoalesce,
		associativity: associativityRight,
		side:          side,
	}
}

func unaryExpressionOperation() expressionOperation {
	return expressionOperation{
		precedence:    precedenceUnary,
		associativity: associativityRight,
		side:          operandSideRight,
	}
}

func ternaryExpressionOperation(side expressionOperandSide) expressionOperation {
	return expressionOperation{
		precedence: precedenceTernary,
		side:       side,
	}
}

func expressionPrimaryAtom(ctx *fql.ExpressionContext) *fql.ExpressionAtomContext {
	if ctx == nil || ctx.UnaryOperator() != nil || ctx.LogicalAndOperator() != nil || ctx.LogicalOrOperator() != nil ||
		ctx.GetCoalesceOperator() != nil || ctx.GetTernaryOperator() != nil {
		return nil
	}

	predicate, ok := ctx.Predicate().(*fql.PredicateContext)
	if !ok || predicate == nil || predicate.EqualityOperator() != nil || predicate.ArrayOperator() != nil ||
		predicate.InOperator() != nil || predicate.LikeOperator() != nil {
		return nil
	}

	atom, ok := predicate.ExpressionAtom().(*fql.ExpressionAtomContext)
	if !ok || atom == nil || atom.MultiplicativeOperator() != nil || atom.AdditiveOperator() != nil || atom.RegexpOperator() != nil {
		return nil
	}

	return atom
}

func parenthesizedExpressionInner(ctx *fql.ExpressionAtomContext) antlr.ParserRuleContext {
	if ctx == nil {
		return nil
	}

	if loop, ok := ctx.ForExpression().(antlr.ParserRuleContext); ok {
		return loop
	}

	if wait, ok := ctx.WaitForExpression().(antlr.ParserRuleContext); ok {
		return wait
	}

	if expr, ok := ctx.Expression().(antlr.ParserRuleContext); ok {
		return expr
	}

	return nil
}

func waitForEventSourceNeedsParentheses(ctx *fql.ExpressionContext) bool {
	atom := expressionPrimaryAtom(ctx)
	if atom == nil || atom.OpenParen() == nil {
		return false
	}

	inner, ok := atom.Expression().(*fql.ExpressionContext)
	if !ok {
		return false
	}

	return expressionExposesInOperator(inner, expressionOperation{})
}

// expressionExposesInOperator follows only expression paths whose formatted
// tokens can reach the outer WAITFOR delimiter. A ternary true branch is
// bounded by its closing colon; condition and false branches are not.
func expressionExposesInOperator(ctx *fql.ExpressionContext, outer expressionOperation) bool {
	if ctx == nil {
		return false
	}

	switch {
	case ctx.UnaryOperator() != nil:
		return expressionChildExposesInOperator(ctx.GetRight(), unaryExpressionOperation())
	case ctx.GetCoalesceOperator() != nil:
		return expressionChildExposesInOperator(
			ctx.GetLeft(),
			coalesceExpressionOperation(operandSideLeft),
		) || expressionChildExposesInOperator(
			ctx.GetRight(),
			coalesceExpressionOperation(operandSideRight),
		)
	case ctx.GetTernaryOperator() != nil:
		return expressionChildExposesInOperator(
			ctx.GetCondition(),
			ternaryExpressionOperation(operandSideLeft),
		) || expressionChildExposesInOperator(
			ctx.GetOnFalse(),
			ternaryExpressionOperation(operandSideRight),
		)
	case ctx.Predicate() != nil:
		predicate, ok := ctx.Predicate().(*fql.PredicateContext)
		if !ok {
			return false
		}

		return predicateExposesInOperator(predicate, outer)
	default:
		precedence := expressionRootPrecedence(ctx)
		if precedence == precedenceNone {
			return false
		}

		return expressionChildExposesInOperator(
			ctx.GetLeft(),
			binaryExpressionOperation(precedence, operandSideLeft),
		) || expressionChildExposesInOperator(
			ctx.GetRight(),
			binaryExpressionOperation(precedence, operandSideRight),
		)
	}
}

func expressionChildExposesInOperator(ctx fql.IExpressionContext, outer expressionOperation) bool {
	expression, ok := ctx.(*fql.ExpressionContext)
	if !ok {
		return false
	}

	return expressionExposesInOperator(expression, outer)
}

func predicateExposesInOperator(ctx *fql.PredicateContext, outer expressionOperation) bool {
	if ctx == nil {
		return false
	}

	if ctx.InOperator() != nil {
		return true
	}

	if operator := ctx.ArrayOperator(); operator != nil && operator.InOperator() != nil {
		return true
	}

	if ctx.ExpressionAtom() != nil {
		atom, ok := ctx.ExpressionAtom().(*fql.ExpressionAtomContext)
		if !ok {
			return false
		}

		return expressionAtomExposesInOperator(atom, outer)
	}

	precedence := predicateRootPrecedence(ctx)
	if precedence == precedenceNone {
		return false
	}

	return predicateChildExposesInOperator(
		ctx.GetLeft(),
		binaryExpressionOperation(precedence, operandSideLeft),
	) || predicateChildExposesInOperator(
		ctx.GetRight(),
		binaryExpressionOperation(precedence, operandSideRight),
	)
}

func predicateChildExposesInOperator(ctx fql.IPredicateContext, outer expressionOperation) bool {
	predicate, ok := ctx.(*fql.PredicateContext)
	if !ok {
		return false
	}

	return predicateExposesInOperator(predicate, outer)
}

func expressionAtomExposesInOperator(ctx *fql.ExpressionAtomContext, outer expressionOperation) bool {
	if ctx == nil {
		return false
	}

	if ctx.OpenParen() != nil && ctx.ErrorOperator() == nil && ctx.RecoveryTails() == nil {
		inner, ok := ctx.Expression().(*fql.ExpressionContext)
		if !ok || !canRemoveExpressionParentheses(inner, outer) {
			return false
		}

		return expressionExposesInOperator(inner, outer)
	}

	precedence := expressionAtomRootPrecedence(ctx)
	if precedence == precedencePrimary {
		return false
	}

	return expressionAtomChildExposesInOperator(
		ctx.GetLeft(),
		binaryExpressionOperation(precedence, operandSideLeft),
	) || expressionAtomChildExposesInOperator(
		ctx.GetRight(),
		binaryExpressionOperation(precedence, operandSideRight),
	)
}

func expressionAtomChildExposesInOperator(ctx fql.IExpressionAtomContext, outer expressionOperation) bool {
	atom, ok := ctx.(*fql.ExpressionAtomContext)
	if !ok {
		return false
	}

	return expressionAtomExposesInOperator(atom, outer)
}
