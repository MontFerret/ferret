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

func New() runtime.Functions {
	ns := runtime.NewRootNamespace()

	if err := RegisterLib(ns); err != nil {
		panic(err)
	}

	return ns.Functions().Build()
}

func RegisterLib(ns runtime.Namespace) error {
	if err := types.RegisterLib(ns); err != nil {
		return err
	}

	if err := strings.RegisterLib(ns); err != nil {
		return err
	}

	if err := math.RegisterLib(ns); err != nil {
		return err
	}

	if err := collections.RegisterLib(ns); err != nil {
		return err
	}

	if err := datetime.RegisterLib(ns); err != nil {
		return err
	}

	if err := arrays.RegisterLib(ns); err != nil {
		return err
	}

	if err := objects.RegisterLib(ns); err != nil {
		return err
	}

	////if err := html.RegisterLib(ns); err != nil {
	////	return err
	////}
	//
	if err := io.RegisterLib(ns); err != nil {
		return err
	}

	if err := path.RegisterLib(ns); err != nil {
		return err
	}

	if err := utils.RegisterLib(ns); err != nil {
		return err
	}

	return testing.RegisterLib(ns)
}
