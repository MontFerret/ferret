package datasource

import (
	"context"
	"github.com/MontFerret/ferret/pkg/runtime/collections"
	"github.com/MontFerret/ferret/pkg/runtime/core"
)

type (
	Variables  []string
	DataSource interface {
		Variables() Variables
		Iterate(ctx context.Context, scope *core.Scope) (collections.Iterator, error)
	}
)

func (vars Variables) Apply(scope *core.Scope, set collections.ResultSet) error {
	err := ValidateResultSet(vars, set)

	if err != nil {
		return err
	}

	for i, variable := range vars {
		if variable != "" {
			scope.SetVariable(variable, set[i])
		}
	}

	return nil
}
