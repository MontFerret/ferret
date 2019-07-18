package stdlib

import (
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/stdlib/arrays"
	"github.com/MontFerret/ferret/pkg/stdlib/collections"
	"github.com/MontFerret/ferret/pkg/stdlib/html"
	"github.com/MontFerret/ferret/pkg/stdlib/math"
	"github.com/MontFerret/ferret/pkg/stdlib/objects"
	"github.com/MontFerret/ferret/pkg/stdlib/strings"
	"github.com/MontFerret/ferret/pkg/stdlib/types"
	"github.com/MontFerret/ferret/pkg/stdlib/utils"
)

func NewLib(ns core.Namespace) error {
	if err := types.NewLib(ns); err != nil {
		return err
	}

	if err := strings.NewLib(ns); err != nil {
		return err
	}

	if err := math.NewLib(ns); err != nil {
		return err
	}

	if err := collections.NewLib(ns); err != nil {
		return err
	}

	if err := arrays.NewLib(ns); err != nil {
		return err
	}

	if err := objects.NewLib(ns); err != nil {
		return err
	}

	if err := html.NewLib(ns); err != nil {
		return err
	}

	return utils.NewLib(ns)
}
