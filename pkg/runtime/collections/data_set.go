package collections

import (
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
)

type DataSet map[string]core.Value

func NewDataSet() DataSet {
	return make(DataSet)
}

func (ds DataSet) Apply(scope *core.Scope, variables Variables) error {
	if err := ValidateDataSet(ds, variables); err != nil {
		return err
	}

	for _, variable := range variables {
		if variable != "" {
			value, found := ds[variable]

			if !found {
				return core.Errorf(core.ErrNotFound, "variable not found in a given data set: %s", variable)
			}

			scope.SetVariable(variable, value)
		}
	}

	return nil
}

func (ds DataSet) Set(key string, value core.Value) {
	ds[key] = value
}

func (ds DataSet) Get(key string) core.Value {
	val, found := ds[key]

	if found {
		return val
	}

	return values.None
}
