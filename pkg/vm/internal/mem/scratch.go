package mem

import "github.com/MontFerret/ferret/v2/pkg/runtime"

type Scratch struct {
	Params   []runtime.Value
	HostArgs []runtime.Value
}

func NewScratch(params int) *Scratch {
	return &Scratch{
		Params: makeNoneValues(params),
	}
}

func (s *Scratch) ResizeParams(size int) {
	resizeNoneValues(&s.Params, size)
}

func (s *Scratch) ResizeHostArgs(size int) {
	resizeNoneValues(&s.HostArgs, size)
}

func makeNoneValues(size int) []runtime.Value {
	if size <= 0 {
		return nil
	}

	values := make([]runtime.Value, size)
	for i := 0; i < size; i++ {
		values[i] = runtime.None
	}

	return values
}

func resizeNoneValues(values *[]runtime.Value, size int) {
	current := *values
	if size < 0 || size == len(current) {
		return
	}

	prevSize := len(current)

	if size < prevSize {
		for i := size; i < prevSize; i++ {
			current[i] = runtime.None
		}

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

	for i := prevSize; i < size; i++ {
		current[i] = runtime.None
	}

	*values = current
}
