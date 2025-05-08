package stdlib

import (
	"github.com/MontFerret/ferret/pkg/runtime"
	"github.com/MontFerret/ferret/pkg/stdlib/arrays"
	"github.com/MontFerret/ferret/pkg/stdlib/collections"
	"github.com/MontFerret/ferret/pkg/stdlib/math"
	"github.com/MontFerret/ferret/pkg/stdlib/strings"
	"github.com/MontFerret/ferret/pkg/stdlib/testing"
)

func RegisterLib(ns runtime.Namespace) error {
	//if err := types.RegisterLib(ns); err != nil {
	//	return err
	//}

	if err := strings.RegisterLib(ns); err != nil {
		return err
	}

	if err := math.RegisterLib(ns); err != nil {
		return err
	}

	if err := collections.RegisterLib(ns); err != nil {
		return err
	}

	//if err := datetime.RegisterLib(ns); err != nil {
	//	return err
	//}
	//
	if err := arrays.RegisterLib(ns); err != nil {
		return err
	}

	//if err := objects.RegisterLib(ns); err != nil {
	//	return err
	//}
	//
	////if err := html.RegisterLib(ns); err != nil {
	////	return err
	////}
	//
	//if err := io.RegisterLib(ns); err != nil {
	//	return err
	//}
	//
	//if err := path.RegisterLib(ns); err != nil {
	//	return err
	//}
	//
	//if err := utils.RegisterLib(ns); err != nil {
	//	return err
	//}

	return testing.RegisterLib(ns)
}
