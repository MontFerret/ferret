package internal

import (
	"github.com/MontFerret/ferret/pkg/compiler/internal/core"
	"github.com/MontFerret/ferret/pkg/vm"
	"github.com/antlr4-go/antlr/v4"
	"strings"

	"github.com/MontFerret/ferret/pkg/runtime"

	"github.com/pkg/errors"
)

func loadConstant(ctx *FuncContext, value runtime.Value) vm.Operand {
	reg := ctx.Registers.Allocate(core.Temp)
	ctx.Emitter.EmitLoadConst(reg, ctx.Symbols.AddConstant(value))
	return reg
}

func sortDirection(dir antlr.TerminalNode) runtime.SortDirection {
	if dir == nil {
		return runtime.SortDirectionAsc
	}

	if strings.ToLower(dir.GetText()) == "desc" {
		return runtime.SortDirectionDesc
	}

	return runtime.SortDirectionAsc
}

func copyFromNamespace(fns *runtime.Functions, namespace string) error {
	// In the name of the function "A::B::C", the namespace is "A::B",
	// not "A::B::".
	//
	// So add "::" at the end.
	namespace += "::"

	// core.Functions cast every function name to upper case. Thus
	// namespace also should be in upper case.
	namespace = strings.ToUpper(namespace)

	for _, name := range fns.Names() {
		if !strings.HasPrefix(name, namespace) {
			continue
		}

		noprefix := strings.Replace(name, namespace, "", 1)

		if _, exists := fns.Get(noprefix); exists {
			return errors.Errorf(
				`collision occurred: "%s" already registered`,
				noprefix,
			)
		}

		fn, _ := fns.Get(name)
		fns.Set(noprefix, fn)
	}

	return nil
}
