package mem

import "github.com/MontFerret/ferret/v2/pkg/runtime"

type Scratch struct {
	Params        []runtime.Value
	MissingParams []bool
	HostArgs      []runtime.Value
}

func NewScratch(params int) Scratch {
	return Scratch{
		Params:        makeNoneValues(params),
		MissingParams: make([]bool, params),
		HostArgs:      nil,
	}
}

func (s *Scratch) ResizeParams(size int) {
	resizeNoneValues(&s.Params, size)

	current := s.MissingParams
	if size < 0 || size == len(current) {
		return
	}

	prevSize := len(current)

	if size < prevSize {
		clear(current[size:prevSize])
		s.MissingParams = current[:size]
		return
	}

	if size > cap(current) {
		resized := make([]bool, size)
		copy(resized, current)
		current = resized
	} else {
		current = current[:size]
	}

	clear(current[prevSize:size])
	s.MissingParams = current
}

func (s *Scratch) ResizeHostArgs(size int) {
	resizeNoneValues(&s.HostArgs, size)
}

// Reset scrubs scratch slots to runtime.None. Scratch storage never closes
// values directly because params and staged args are borrowed.
func (s *Scratch) Reset() {
	fillWithNone(s.Params)
	clear(s.MissingParams)
	fillWithNone(s.HostArgs)
}
