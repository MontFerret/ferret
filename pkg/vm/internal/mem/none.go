package mem

import "github.com/MontFerret/ferret/v2/pkg/runtime"

func fillWithNone(values []runtime.Value) {
	for i := range values {
		values[i] = runtime.None
	}
}
