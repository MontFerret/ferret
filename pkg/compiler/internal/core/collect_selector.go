package core

import (
	"github.com/MontFerret/ferret/pkg/runtime"
)

type CollectSelector struct {
	name runtime.String
}

func NewCollectSelector(name runtime.String) *CollectSelector {
	return &CollectSelector{
		name: name,
	}
}

func (s *CollectSelector) Name() runtime.String {
	return s.name
}
