package mem

import "github.com/MontFerret/ferret/v2/pkg/runtime"

type Scratch struct {
	Params []runtime.Value
}

func NewScratch(params int) *Scratch {
	paramSlots := make([]runtime.Value, params)

	for i := 0; i < params; i++ {
		paramSlots[i] = runtime.None
	}

	return &Scratch{
		Params: paramSlots,
	}
}

func (s *Scratch) ResizeParams(size int) {
	if size < 0 || size == len(s.Params) {
		return
	}

	prevSize := len(s.Params)

	if size < prevSize {
		for i := size; i < prevSize; i++ {
			s.Params[i] = runtime.None
		}

		s.Params = s.Params[:size]
		return
	}

	if size > cap(s.Params) {
		params := make([]runtime.Value, size)
		copy(params, s.Params)
		s.Params = params

		for i := prevSize; i < size; i++ {
			s.Params[i] = runtime.None
		}

		return
	}

	s.Params = s.Params[:size]

	for i := prevSize; i < size; i++ {
		s.Params[i] = runtime.None
	}
}
