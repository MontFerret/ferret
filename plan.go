package ferret

import "github.com/MontFerret/ferret/pkg/vm"

type Plan struct {
	prog *vm.Program
}

func (p *Plan) NewSession() *Session {
	return &Session{
		vm: vm.New(p.prog),
	}
}
