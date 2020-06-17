package stdlib

import (
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/stdlib/arrays"
	"github.com/MontFerret/ferret/pkg/stdlib/collections"
	"github.com/MontFerret/ferret/pkg/stdlib/datetime"
	"github.com/MontFerret/ferret/pkg/stdlib/html"
	"github.com/MontFerret/ferret/pkg/stdlib/io"
	"github.com/MontFerret/ferret/pkg/stdlib/math"
	"github.com/MontFerret/ferret/pkg/stdlib/objects"
	"github.com/MontFerret/ferret/pkg/stdlib/path"
	"github.com/MontFerret/ferret/pkg/stdlib/strings"
	"github.com/MontFerret/ferret/pkg/stdlib/testing"
	"github.com/MontFerret/ferret/pkg/stdlib/types"
	"github.com/MontFerret/ferret/pkg/stdlib/utils"
)

func RegisterLib(ns core.Namespace) error {
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

	if err := html.RegisterLib(ns); err != nil {
		return err
	}

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
