package mem

import "github.com/MontFerret/ferret/v2/pkg/runtime"

func fillWithNone(values []runtime.Value) {
	for i := range values {
		values[i] = runtime.None
	}
}

func makeNoneValues(size int) []runtime.Value {
	if size <= 0 {
		return nil
	}

	values := make([]runtime.Value, size)
	fillWithNone(values)

	return values
}

func resizeNoneValues(values *[]runtime.Value, size int) {
	current := *values
	if size < 0 || size == len(current) {
		return
	}

	prevSize := len(current)

	if size < prevSize {
		fillWithNone(current[size:prevSize])

		*values = current[:size]
		return
	}

	if size > cap(current) {
		resized := make([]runtime.Value, size)
		copy(resized, current)
		current = resized
	} else {
		current = current[:size]
	}

	fillWithNone(current[prevSize:size])

	*values = current
}
