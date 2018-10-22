package clauses

import "github.com/MontFerret/ferret/pkg/runtime/core"

type CollectSelector struct {
	variable   string
	expression core.Expression
}

func NewCollectSelector(variable string, exp core.Expression) (*CollectSelector, error) {
	if variable == "" {
		return nil, core.Error(core.ErrMissedArgument, "selector variable")
	}

	if exp == nil {
		return nil, core.Error(core.ErrMissedArgument, "selector expression")
	}

	return &CollectSelector{variable, exp}, nil
}

func (selector *CollectSelector) Variable() string {
	return selector.variable
}

func (selector *CollectSelector) Expression() core.Expression {
	return selector.expression
}
