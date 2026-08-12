package vm_test

import "github.com/MontFerret/ferret/v2/pkg/runtime"

type genericMap struct {
	runtime.Map
}

func (genericMap) Type() runtime.Type {
	return runtime.TypeMap
}
