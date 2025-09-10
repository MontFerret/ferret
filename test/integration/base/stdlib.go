package base

import (
	"github.com/MontFerret/ferret/pkg/runtime"
	"github.com/MontFerret/ferret/pkg/stdlib"
)

func Stdlib() runtime.Functions {
	return stdlib.New()
}
