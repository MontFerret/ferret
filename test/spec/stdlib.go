package spec

import (
	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/stdlib"
)

func Stdlib() runtime.Namespace {
	return stdlib.New()
}
