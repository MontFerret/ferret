package vm_test

import "github.com/MontFerret/ferret/v2/pkg/runtime"

type spreadObjectLike struct {
	runtime.Map
}

func (*spreadObjectLike) ObjectLike() {}
