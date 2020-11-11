package compiler

import (
	"regexp"
	"strings"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/pkg/errors"
)

var fnNameValidation = regexp.MustCompile("^[a-zA-Z]+[a-zA-Z0-9_]*(::[a-zA-Z]+[a-zA-Z0-9_]*)*$")

const emptyNS = ""
const separator = "::"

type NamespaceContainer struct {
	funcs *core.Functions
	name  string
}

func newRootNamespace() *NamespaceContainer {
	ns := new(NamespaceContainer)
	ns.funcs = core.NewFunctions()

	return ns
}

func newNamespace(funcs *core.Functions, name string) *NamespaceContainer {
	return &NamespaceContainer{funcs, strings.ToUpper(name)}
}

func (nc *NamespaceContainer) Namespace(name string) core.Namespace {
	return newNamespace(nc.funcs, nc.makeFullName(name))
}

func (nc *NamespaceContainer) MustRegisterFunction(name string, fun core.Function) {
	if err := nc.RegisterFunction(name, fun); err != nil {
		panic(err)
	}
}

func (nc *NamespaceContainer) RegisterFunction(name string, fun core.Function) error {
	nsName := nc.makeFullName(name)

	_, exists := nc.funcs.Get(nsName)

	if exists {
		return errors.Errorf("function already exists: %s", name)
	}

	// validation the name
	if strings.Contains(name, separator) {
		return errors.Errorf("invalid function name: %s", name)
	}

	if !fnNameValidation.MatchString(nsName) {
		return errors.Errorf("invalid function or namespace name: %s", nsName)
	}

	nc.funcs.Set(nsName, fun)

	return nil
}

func (nc *NamespaceContainer) RemoveFunction(name string) {
	nc.funcs.Unset(nc.makeFullName(name))
}

func (nc *NamespaceContainer) MustRegisterFunctions(funcs *core.Functions) {
	if err := nc.RegisterFunctions(funcs); err != nil {
		panic(err)
	}
}

func (nc *NamespaceContainer) RegisterFunctions(funcs *core.Functions) error {
	for _, name := range funcs.Names() {
		fun, _ := funcs.Get(name)

		if err := nc.RegisterFunction(name, fun); err != nil {
			return err
		}
	}

	return nil
}

func (nc *NamespaceContainer) RegisteredFunctions() []string {
	fnames := nc.funcs.Names()
	res := make([]string, 0, len(fnames))

	// root namespace, return all functions
	if nc.name == emptyNS {
		for _, k := range fnames {
			res = append(res, k)
		}
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

func (nc *NamespaceContainer) makeFullName(name string) string {
	if nc.name == emptyNS {
		return name
	}

	return nc.name + separator + name
}
