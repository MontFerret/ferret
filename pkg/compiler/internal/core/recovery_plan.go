package core

import (
	"github.com/antlr4-go/antlr/v4"

	"github.com/MontFerret/ferret/v2/pkg/parser/fql"
)

type (
	RecoveryCondition int

	RecoveryActionKind int

	RecoveryPlan struct {
		OnError   *RecoveryHandler
		OnTimeout *RecoveryHandler
	}

	RecoveryHandler struct {
		ActionNode antlr.ParserRuleContext
		Expr       fql.IExpressionContext
		Retry      *RecoveryRetryPlan
		TailNode   antlr.ParserRuleContext
		ActionKind RecoveryActionKind
	}

	RecoveryRetryPlan struct {
		ActionNode      antlr.ParserRuleContext
		FinalActionNode antlr.ParserRuleContext
		FinalExpr       fql.IExpressionContext
		RetryPlan
		FinalActionKind RecoveryActionKind
		HasOr           bool
	}

	RecoveryPlanOptions struct {
		AllowTimeout bool
		HasTimeout   bool
	}

	RecoveryTailOwner interface {
		RecoveryTails() fql.IRecoveryTailsContext
	}
)

const SuppressRecoveryAction = "SUPPRESS"

const (
	RecoveryConditionUnknown RecoveryCondition = iota
	RecoveryConditionError
	RecoveryConditionTimeout
)

const (
	RecoveryActionUnknown RecoveryActionKind = iota
	RecoveryActionFail
	RecoveryActionReturn
	RecoveryActionRetry
)
