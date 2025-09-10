package runtime

import (
	"strings"
)

const NamespaceSeparator = "::"
const emptyNS = ""

type (
	Namespace interface {
		Namespace(name string) Namespace
		Functions() FunctionsBuilder
		RegisterFunctions(funcs Functions) error
	}

	defaultNamespace struct {
		builder *defaultFunctionBuilder
		name    string
	}
)

func NewRootNamespace() Namespace {
	ns := new(defaultNamespace)
	ns.builder = newRootFunctionsBuilder()

	return ns
}

func NewNamespace(name string) Namespace {
	ns := new(defaultNamespace)
	ns.name = strings.ToUpper(name)
	ns.builder = newNamespaceFunctionsBuilder(ns.name)

	return ns
}

func (nc *defaultNamespace) Namespace(name string) Namespace {
	ns := new(defaultNamespace)
	ns.name = makeFunctionName(nc.name, name)
	ns.builder = newFunctionsBuilderInternalFrom(ns.name, nc.builder)

	return ns
}

func (nc *defaultNamespace) Functions() FunctionsBuilder {
	return nc.builder
}

func (nc *defaultNamespace) RegisterFunctions(funcs Functions) error {
	for _, fname := range funcs.Names() {
		if nc.builder.Has(fname) {
			return Errorf(ErrNotUnique, "function '%s' already exists", makeFunctionName(nc.name, fname))
		}
	}

	nc.builder.SetFrom(funcs)

	return nil
}
