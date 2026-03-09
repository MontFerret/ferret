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

	if size > cap(s.Params) {
		s.Params = make([]runtime.Value, size)
		return
	}

	s.Params = s.Params[:size]
}
