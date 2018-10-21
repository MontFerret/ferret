package collections

import (
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
)

type (
	Variables []values.String

	DataSet map[values.String]core.Value
)

var (
	DefaultValueVar = values.NewString("value")
	DefaultKeyVar   = values.NewString("key")
)

func (ds DataSet) Length() values.Int {
	return values.NewInt(len(ds))
}

func (ds DataSet) Keys() []string {
	keys := make([]string, 0, len(ds))

	for key := range ds {
		keys = append(keys, key.String())
	}

	return keys
}

func (ds DataSet) Get(key values.String) (core.Value, values.Boolean) {
	val, found := ds[key]

	if found {
		return val, values.True
	}

	return values.None, values.False
}

func (ds DataSet) GetOrDefault(key values.String) core.Value {
	val, found := ds[key]

	if found {
		return val
	}

	return values.None
}

func (ds DataSet) Set(key values.String, value core.Value) {
	ds[key] = value
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

			scope.SetVariable(variable.String(), value)
		}
	}

	return nil
}
