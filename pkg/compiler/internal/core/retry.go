package core

import (
	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/parser/fql"
	"github.com/antlr4-go/antlr/v4"
)

type (
	RetryBackoff int

	RetryPlan struct {
		Delay       DurationClause
		Backoff     RetryBackoff
		Count       int
		CountNode   antlr.ParserRuleContext
		DelayNode   antlr.ParserRuleContext
		BackoffNode antlr.ParserRuleContext
		HasDelay    bool
	}

	DurationClause interface {
		DurationLiteral() fql.IDurationLiteralContext
		IntegerLiteral() fql.IIntegerLiteralContext
		FloatLiteral() fql.IFloatLiteralContext
		Variable() fql.IVariableContext
		Param() fql.IParamContext
		MemberExpression() fql.IMemberExpressionContext
		FunctionCall() fql.IFunctionCallContext
	}

	RetryDelayState struct {
		BaseReg    bytecode.Operand
		CurrentReg bytecode.Operand
		ReadyReg   bytecode.Operand
	}
)

const (
	RetryBackoffNone RetryBackoff = iota
	RetryBackoffLinear
	RetryBackoffExponential
)
