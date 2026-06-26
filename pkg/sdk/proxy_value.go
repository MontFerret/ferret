package sdk

import (
	"reflect"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

func proxyValueEqual[T any](expected runtime.Value, actual T, actualValue runtime.Value) bool {
	if expected == nil {
		expected = runtime.None
	}

	if runtimeValue, ok := any(actual).(runtime.Value); ok {
		if runtime.CompareValues(expected, runtimeValue) == 0 {
			return true
		}

		return reflect.DeepEqual(proxyUnwrap(expected), proxyUnwrap(runtimeValue))
	}

	if runtime.CompareValues(expected, actualValue) == 0 {
		return true
	}

	return reflect.DeepEqual(proxyUnwrap(expected), actual)
}

func proxyUnwrap(value runtime.Value) any {
	if unwrappable, ok := value.(runtime.Unwrappable); ok {
		return unwrappable.Unwrap()
	}

	return value
}
