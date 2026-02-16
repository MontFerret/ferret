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
	Keys  []runtime.String
	Kinds []AggregateKind
	Index map[string]int
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
