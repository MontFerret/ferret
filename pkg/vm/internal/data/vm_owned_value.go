package data

import "github.com/MontFerret/ferret/v2/pkg/runtime"

// IsVMOwnedValue reports whether value is implemented by the VM's internal
// data layer rather than by an external runtime capability.
func IsVMOwnedValue(value runtime.Value) bool {
	switch value.(type) {
	case *AggregateCollector,
		*GroupedAggregateCollector,
		*CounterCollector,
		*DataSet,
		*FastObject,
		*Iterator,
		*ClosableIterator,
		*KeyCollector,
		*KeyCounterCollector,
		*KeyGroupCollector,
		*MultiSorter,
		*Regexp,
		*Sorter,
		*StreamValue,
		*AggregateKey,
		*KV,
		*noopCollector:
		return true
	default:
		return false
	}
}
