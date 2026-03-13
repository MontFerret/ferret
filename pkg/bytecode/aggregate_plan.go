package bytecode

import (
	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

type AggregateKind int

const (
	AggregateCount AggregateKind = iota
	AggregateSum
	AggregateMin
	AggregateMax
	AggregateAverage
)

type AggregatePlan struct {
	Index map[string]int   `json:"index"`
	Keys  []runtime.String `json:"keys"`
	Kinds []AggregateKind  `json:"kinds"`
}

func NewAggregatePlan(keys []runtime.String, kinds []AggregateKind) AggregatePlan {
	idx := make(map[string]int, len(keys))

	for i, key := range keys {
		idx[key.String()] = i
	}

	return AggregatePlan{
		Keys:  keys,
		Kinds: kinds,
		Index: idx,
	}
}
