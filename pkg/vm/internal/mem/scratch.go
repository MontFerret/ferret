package mem

import "github.com/MontFerret/ferret/v2/pkg/runtime"

type Scratch struct {
	Params   []runtime.Value
	HostArgs []runtime.Value
}

func NewScratch(params int) *Scratch {
	s := &Scratch{}
	s.Init(params)

	return s
}

func (s *Scratch) Init(params int) {
	s.Params = makeNoneValues(params)
	s.HostArgs = nil
}

func (s *Scratch) ResizeParams(size int) {
	resizeNoneValues(&s.Params, size)
}

func (s *Scratch) ResizeHostArgs(size int) {
	resizeNoneValues(&s.HostArgs, size)
}

func (s *Scratch) Reset() {
	fillWithNoneAndClose(s.Params)
	fillWithNone(s.HostArgs)
}
