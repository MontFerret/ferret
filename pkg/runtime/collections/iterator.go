package collections

import (
	"context"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
)

type (
	Variables []string

	DataSet map[string]core.Value

	Iterator interface {
		HasNext() bool
		Next() (DataSet, error)
	}

	Iterable interface {
		Variables() Variables
		Iterate(ctx context.Context, scope *core.Scope) (Iterator, error)
	}
)

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

func (ds DataSet) Get(key string) core.Value {
	val, found := ds[key]

	if found {
		return val
	}

	return values.None
}

func ToSlice(iterator Iterator) ([]DataSet, error) {
	res := make([]DataSet, 0, 10)

	for iterator.HasNext() {
		ds, err := iterator.Next()

		if err != nil {
			return nil, err
		}

		res = append(res, ds)
	}

	return res, nil
}
