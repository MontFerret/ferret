package compiler

import (
	"github.com/MontFerret/ferret/pkg/runtime"
	"regexp"
	"strings"

	"github.com/pkg/errors"
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

func (nc *NamespaceContainer) MustRegisterFunction(name string, fun runtime.Function) {
	if err := nc.RegisterFunction(name, fun); err != nil {
		panic(err)
	}
}

func (nc *NamespaceContainer) RegisterFunction(name string, fun runtime.Function) error {
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

func (nc *NamespaceContainer) MustRegisterFunctions(funcs runtime.Functions) {
	if err := nc.RegisterFunctions(funcs); err != nil {
		panic(err)
	}
}

func (nc *NamespaceContainer) RegisterFunctions(funcs runtime.Functions) error {
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
