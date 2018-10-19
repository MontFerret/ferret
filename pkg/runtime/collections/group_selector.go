package collections

import "github.com/MontFerret/ferret/pkg/runtime/core"

type (
	GroupSelector interface {
		Key(set ResultSet) (core.Value, error)
		Value(set ResultSet) (core.Value, error)
	}

	GenericSelectorFn func(set ResultSet) (core.Value, error)

	GenericGroupSelector struct {
		key   GenericSelectorFn
		value GenericSelectorFn
	}
)

func NewGenericGroupSelector(key, value GenericSelectorFn) GroupSelector {
	return &GenericGroupSelector{key, value}
}

func (selector *GenericGroupSelector) Key(set ResultSet) (core.Value, error) {
	return selector.key(set)
}

func (selector *GenericGroupSelector) Value(set ResultSet) (core.Value, error) {
	return selector.value(set)
}
