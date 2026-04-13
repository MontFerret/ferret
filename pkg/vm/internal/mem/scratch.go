package mem

import "github.com/MontFerret/ferret/v2/pkg/runtime"

type Scratch struct {
	Params   []runtime.Value
	HostArgs []runtime.Value
}

func NewScratch(params int) Scratch {
	return Scratch{
		Params:   makeNoneValues(params),
		HostArgs: nil,
	}
}

func (s *Scratch) ResizeParams(size int) {
	resizeNoneValues(&s.Params, size)
}

func (s *Scratch) ResizeHostArgs(size int) {
	resizeNoneValues(&s.HostArgs, size)
}

// Reset scrubs scratch slots to runtime.None. Scratch storage never closes
// values directly because params and staged args are borrowed.
func (s *Scratch) Reset() {
	fillWithNone(s.Params)
	fillWithNone(s.HostArgs)
}
