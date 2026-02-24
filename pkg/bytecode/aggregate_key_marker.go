package bytecode

import (
	"encoding/json"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

var (
	// AggregateKeyMarker is a sentinel value used to disambiguate internal aggregate keys
	// from user-provided group keys in fused grouped aggregation.
	AggregateKeyMarker = &aggregateKeyMarker{}

	// typeAggregateKeyMarker is the runtime type of AggregateKeyMarker.
	// It is used for encoding and decoding the marker in bytecode constants.
	typeAggregateKeyMarker = runtime.NewTypeFor[*aggregateKeyMarker]("bytecode", "__agg_key_marker__")
)

type aggregateKeyMarker struct{}

func (m *aggregateKeyMarker) Type() runtime.Type {
	return typeAggregateKeyMarker
}

func (m *aggregateKeyMarker) MarshalJSON() ([]byte, error) {
	return json.Marshal(typeAggregateKeyMarker.Name())
}

func (m *aggregateKeyMarker) String() string {
	return typeAggregateKeyMarker.Name()
}

func (m *aggregateKeyMarker) Hash() uint64 {
	return 0x9e3779b97f4a7c15
}

func (m *aggregateKeyMarker) Copy() runtime.Value {
	return m
}
