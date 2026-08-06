package vm

import (
	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/vm/internal/data"
)

func isExternalCapabilityReceiver(value runtime.Value) bool {
	if value == nil || value == runtime.None {
		return false
	}

	switch value.(type) {
	case runtime.Boolean,
		runtime.Int,
		runtime.Float,
		runtime.Duration,
		runtime.String,
		runtime.DateTime,
		runtime.Binary,
		*runtime.Array,
		*runtime.Object,
		*runtime.Range:
		return false
	default:
		return !data.IsVMOwnedValue(value)
	}
}

func comparisonNeedsSafepoint(left, right runtime.Value) bool {
	return !isCheapComparisonOperand(left) || !isCheapComparisonOperand(right)
}

func isCheapComparisonOperand(value runtime.Value) bool {
	if value == nil || value == runtime.None {
		return true
	}

	switch value.(type) {
	case runtime.Boolean,
		runtime.Int,
		runtime.Float,
		runtime.Duration,
		runtime.String,
		runtime.DateTime,
		runtime.Binary,
		*runtime.Range,
		*data.Regexp:
		return true
	default:
		return false
	}
}

func propertyReadUsesExternalCapability(src, prop runtime.Value) bool {
	switch prop.(type) {
	case runtime.Int, runtime.Float:
		_, ok := src.(runtime.IndexReadable)
		return ok && isExternalCapabilityReceiver(src)
	default:
		_, ok := src.(runtime.KeyReadable)
		return ok && isExternalCapabilityReceiver(src)
	}
}

func propertyWriteUsesExternalCapability(target, prop runtime.Value) bool {
	switch value := prop.(type) {
	case runtime.Int:
		if value < 0 {
			return false
		}

		_, ok := target.(runtime.IndexWritable)
		return ok && isExternalCapabilityReceiver(target)
	case runtime.Float:
		return false
	default:
		_, ok := target.(runtime.KeyWritable)
		return ok && isExternalCapabilityReceiver(target)
	}
}

func propertyDeleteUsesExternalCapability(target, prop runtime.Value) bool {
	switch value := prop.(type) {
	case runtime.Int:
		if value < 0 {
			return false
		}

		_, ok := target.(runtime.IndexRemovable)
		return ok && isExternalCapabilityReceiver(target)
	case runtime.Float:
		return false
	default:
		_, ok := target.(runtime.KeyRemovable)
		return ok && isExternalCapabilityReceiver(target)
	}
}

func keyDeleteUsesExternalCapability(target, key runtime.Value) bool {
	if _, ok := target.(runtime.KeyRemovable); ok {
		return isExternalCapabilityReceiver(target)
	}

	idx, ok := key.(runtime.Int)
	if !ok || idx < 0 {
		return false
	}

	_, ok = target.(runtime.IndexRemovable)
	return ok && isExternalCapabilityReceiver(target)
}
