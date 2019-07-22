package compiler

import (
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/pkg/errors"
	"regexp"
	"strings"
)

var fnNameValidation = regexp.MustCompile("^[a-zA-Z]+[a-zA-Z0-9_]*(::[a-zA-Z]+[a-zA-Z0-9_]*)*$")

const emptyNS = ""
const separator = "::"

type NamespaceContainer struct {
	funcs core.Functions
	name  string
}

func newRootNamespace() *NamespaceContainer {
	ns := new(NamespaceContainer)
	ns.funcs = make(core.Functions)

	return ns
}

func newNamespace(funcs core.Functions, name string) *NamespaceContainer {
	return &NamespaceContainer{funcs, strings.ToUpper(name)}
}

func (nc *NamespaceContainer) Namespace(name string) core.Namespace {
	return newNamespace(nc.funcs, nc.makeFullName(name))
}

func (nc *NamespaceContainer) RegisterFunction(name string, fun core.Function) error {
	nsName := nc.makeFullName(name)
	_, exists := nc.funcs[nsName]

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

	nc.funcs[strings.ToUpper(nsName)] = fun

	return nil
}

func (nc *NamespaceContainer) RemoveFunction(name string) {
	delete(nc.funcs, strings.ToUpper(nc.makeFullName(name)))
}

func (nc *NamespaceContainer) RegisterFunctions(funcs core.Functions) error {
	for name, fun := range funcs {
		if err := nc.RegisterFunction(name, fun); err != nil {
			return err
		}
	}

	return nil
}

func (nc *NamespaceContainer) RegisteredFunctions() []string {
	res := make([]string, 0, len(nc.funcs))

	// root namespace, return all functions
	if nc.name == emptyNS {
		for k := range nc.funcs {
			res = append(res, k)
		}
	} else {
		nsPrefix := nc.name + separator
		for k := range nc.funcs {
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
