package runtime_test

import "github.com/MontFerret/ferret/v2/pkg/runtime"

type copyContractList struct {
	runtime.List
	copied runtime.Value
	calls  int
}

func (list *copyContractList) Copy() runtime.Value {
	list.calls++

	return list.copied
}
