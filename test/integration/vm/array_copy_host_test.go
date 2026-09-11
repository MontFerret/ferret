package vm_test

import "github.com/MontFerret/ferret/v2/pkg/runtime"

type invalidArrayCopy struct {
	runtime.List
}

func (list *invalidArrayCopy) Copy() runtime.Value {
	return runtime.True
}
