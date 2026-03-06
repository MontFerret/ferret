package ferret

import (
	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/encoding"
	"github.com/MontFerret/ferret/v2/pkg/vm"
)

type Plan struct {
	prog     *bytecode.Program
	env      *vm.Environment
	encoding *encoding.Registry
}

func (p *Plan) NewSession(setters ...SessionOption) (*Session, error) {
	env, err := vm.ExtendEnvironment(p.env, setters)

	if err != nil {
		return nil, err
	}

	return &Session{
		// TODO: create a VM pool and get a VM from it instead of creating a new one for each session
		vm:       vm.New(p.prog),
		env:      env,
		encoding: p.encoding,
	}, nil
}
