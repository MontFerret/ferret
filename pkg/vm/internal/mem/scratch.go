package mem

import "github.com/MontFerret/ferret/v2/pkg/runtime"

type Scratch struct {
	Params   []runtime.Value
	HostArgs []runtime.Value
}

func (s *Scratch) Init(params int) {
	s.Params = makeNoneValues(params)
	s.HostArgs = nil
}

func NewScratch(params int) *Scratch {
	s := &Scratch{}
	s.Init(params)
	return s
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
