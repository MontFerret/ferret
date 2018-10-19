package collections

import "github.com/MontFerret/ferret/pkg/runtime/core"

type (
	Selector      func(set ResultSet) (core.Value, error)
	GroupSelector struct {
		key   Selector
		value Selector
	}
)

func NewGroupSelector(key, value Selector) *GroupSelector {
	return &GroupSelector{key, value}
}

func (selector *GroupSelector) Key(set ResultSet) (core.Value, error) {
	return selector.key(set)
}

func (selector *GroupSelector) Value(set ResultSet) (core.Value, error) {
	return selector.value(set)
}
