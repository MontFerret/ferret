package core

import (
	"github.com/MontFerret/ferret/v2/pkg/runtime"
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
