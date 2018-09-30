package stdlib

import (
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/stdlib/arrays"
	"github.com/MontFerret/ferret/pkg/stdlib/collections"
	"github.com/MontFerret/ferret/pkg/stdlib/html"
	"github.com/MontFerret/ferret/pkg/stdlib/strings"
	"github.com/MontFerret/ferret/pkg/stdlib/types"
	"github.com/MontFerret/ferret/pkg/stdlib/utils"
)

func NewLib() map[string]core.Function {
	lib := make(map[string]core.Function)

	add := func(l map[string]core.Function) {
		for name, fn := range l {
			lib[name] = fn
		}
	}

	add(types.NewLib())
	add(strings.NewLib())
	add(collections.NewLib())
	add(arrays.NewLib())
	add(html.NewLib())
	add(utils.NewLib())

	return lib
}
