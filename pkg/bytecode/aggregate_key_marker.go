package bytecode

import (
	"encoding/json"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

const aggregateKeyMarkerString = "__agg_key__"

// AggregateKeyMarker is a sentinel value used to disambiguate internal aggregate keys
// from user-provided group keys in fused grouped aggregation.
var AggregateKeyMarker = &aggregateKeyMarker{}

type aggregateKeyMarker struct{}

func (m *aggregateKeyMarker) MarshalJSON() ([]byte, error) {
	return json.Marshal(aggregateKeyMarkerString)
}

func (m *aggregateKeyMarker) String() string {
	return aggregateKeyMarkerString
}

func (m *aggregateKeyMarker) Hash() uint64 {
	return 0x9e3779b97f4a7c15
}

func (m *aggregateKeyMarker) Copy() runtime.Value {
	return m
}
