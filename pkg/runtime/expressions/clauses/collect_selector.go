package clauses

import (
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
)

type CollectSelector struct {
	variable values.String
	exp      core.Expression
}

func NewCollectSelector(variable values.String, exp core.Expression) (*CollectSelector, error) {
	if variable == "" {
		return nil, core.Error(core.ErrMissedArgument, "selector variable")
	}

	if exp == nil {
		return nil, core.Error(core.ErrMissedArgument, "selector expression")
	}

	return &CollectSelector{variable, exp}, nil
}

func (selector *CollectSelector) Variable() values.String {
	return selector.variable
}

func (selector *CollectSelector) Expression() core.Expression {
	return selector.exp
}
