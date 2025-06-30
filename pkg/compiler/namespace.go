package compiler

import (
	"regexp"
	"strings"

	"github.com/MontFerret/ferret/pkg/runtime"
)

var fnNameValidation = regexp.MustCompile("^[a-zA-Z]+[a-zA-Z0-9_]*(::[a-zA-Z]+[a-zA-Z0-9_]*)*$")

const emptyNS = ""
const separator = "::"

type NamespaceContainer struct {
	funcs runtime.Functions
	name  string
}

func NewRootNamespace() *NamespaceContainer {
	ns := new(NamespaceContainer)
	ns.funcs = runtime.NewFunctions()

	return ns
}

func NewNamespace(funcs runtime.Functions, name string) *NamespaceContainer {
	return &NamespaceContainer{funcs, strings.ToUpper(name)}
}

func (nc *NamespaceContainer) Namespace(name string) runtime.Namespace {
	return NewNamespace(nc.funcs, nc.makeFullName(name))
}

func (nc *NamespaceContainer) RemoveFunction(name string) {
	nc.funcs.Unset(nc.makeFullName(name))
}

func (nc *NamespaceContainer) RegisterFunctions(funcs runtime.Functions) error {
	nc.funcs.SetAll(funcs)

	return nil
}

func (nc *NamespaceContainer) RegisteredFunctions() []string {
	fnames := nc.funcs.Names()
	res := make([]string, 0, len(fnames))

	// root namespace, return all functions
	if nc.name == emptyNS {
		res = append(res, fnames...)
	} else {
		nsPrefix := nc.name + separator
		for _, k := range fnames {
			if strings.HasPrefix(k, nsPrefix) {
				res = append(res, k)
			}
		}
	}

	return res
}

func (nc *NamespaceContainer) Functions() runtime.Functions {
	return nc.funcs
}

func (nc *NamespaceContainer) makeFullName(name string) string {
	name = strings.ToUpper(name)

	if nc.name == emptyNS {
		return name
	}

	return nc.name + separator + name
}
