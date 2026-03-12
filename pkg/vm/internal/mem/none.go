package mem

import (
	"io"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

func fillWithNone(values []runtime.Value) {
	for i := range values {
		values[i] = runtime.None
	}
}

func fillWithNoneAndClose(values []runtime.Value) {
	for i := range values {
		closer, ok := values[i].(io.Closer)

		if ok {
			// TODO: Should we propagate the error here?
			_ = closer.Close()
		}

		values[i] = runtime.None
	}
}
