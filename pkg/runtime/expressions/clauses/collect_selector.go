package clauses

import "github.com/MontFerret/ferret/pkg/runtime/core"

type (
	CollectSelector struct {
		variable   string
		expression core.Expression
	}

	CollectAggregateSelector struct {
		variable    string
		aggregators []core.Expression
		reducer     core.Function
	}
)

func NewCollectSelector(variable string, exp core.Expression) (*CollectSelector, error) {
	if variable == "" {
		return nil, core.Error(core.ErrMissedArgument, "selector variable")
	}

	if exp == nil {
		return nil, core.Error(core.ErrMissedArgument, "selector reducer")
	}

	return &CollectSelector{variable, exp}, nil
}

func (selector *CollectSelector) Variable() string {
	return selector.variable
}

func (selector *CollectSelector) Expression() core.Expression {
	return selector.expression
}

func NewCollectAggregateSelector(variable string, aggr []core.Expression, reducer core.Function) (*CollectAggregateSelector, error) {
	if variable == "" {
		return nil, core.Error(core.ErrMissedArgument, "selector variable")
	}

	if reducer == nil {
		return nil, core.Error(core.ErrMissedArgument, "selector reducer")
	}

	if aggr == nil {
		return nil, core.Error(core.ErrMissedArgument, "selector aggregators")
	}

	return &CollectAggregateSelector{variable, aggr, reducer}, nil
}

func (selector *CollectAggregateSelector) Variable() string {
	return selector.variable
}

func (selector *CollectAggregateSelector) Expression() core.Function {
	return selector.reducer
}

func (selector *CollectAggregateSelector) Aggregators() []core.Expression {
	return selector.aggregators
}
