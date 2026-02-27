package stdlib

import (
	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/stdlib/arrays"
	"github.com/MontFerret/ferret/v2/pkg/stdlib/collections"
	"github.com/MontFerret/ferret/v2/pkg/stdlib/datetime"
	"github.com/MontFerret/ferret/v2/pkg/stdlib/io"
	"github.com/MontFerret/ferret/v2/pkg/stdlib/math"
	"github.com/MontFerret/ferret/v2/pkg/stdlib/objects"
	"github.com/MontFerret/ferret/v2/pkg/stdlib/path"
	"github.com/MontFerret/ferret/v2/pkg/stdlib/strings"
	"github.com/MontFerret/ferret/v2/pkg/stdlib/testing"
	"github.com/MontFerret/ferret/v2/pkg/stdlib/types"
	"github.com/MontFerret/ferret/v2/pkg/stdlib/utils"
)

// New creates a new standard library.
// It registers all available functions and namespaces to the root namespace and returns it.
func New() runtime.Namespace {
	ns := runtime.NewRootNamespace()

	RegisterLib(ns)

	return ns
}

func RegisterLib(ns runtime.Namespace) {
	libs := []func(runtime.Namespace){
		types.RegisterLib,
		strings.RegisterLib,
		math.RegisterLib,
		collections.RegisterLib,
		datetime.RegisterLib,
		arrays.RegisterLib,
		objects.RegisterLib,
		io.RegisterLib,
		path.RegisterLib,
		utils.RegisterLib,
		testing.RegisterLib,
	}

	for _, lib := range libs {
		lib(ns)
	}
}
