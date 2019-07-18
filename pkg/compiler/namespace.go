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

type NamespaceBuilder struct {
	funcs core.Functions
	name  string
}

func newRootNamespace() *NamespaceBuilder {
	ns := new(NamespaceBuilder)
	ns.funcs = make(core.Functions)

	return ns
}

func newNamespace(funcs core.Functions, name string) *NamespaceBuilder {
	return &NamespaceBuilder{funcs, name}
}

func (rs *NamespaceBuilder) Namespace(name string) core.Namespace {
	return newNamespace(rs.funcs, rs.makeFullName(name))
}

func (rs *NamespaceBuilder) RegisterFunction(name string, fun core.Function) error {
	nsName := rs.makeFullName(name)
	_, exists := rs.funcs[nsName]

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

	rs.funcs[strings.ToUpper(nsName)] = fun

	return nil
}

func (rs *NamespaceBuilder) RemoveFunction(name string) {
	delete(rs.funcs, strings.ToUpper(rs.makeFullName(name)))
}

func (rs *NamespaceBuilder) RegisterFunctions(funcs core.Functions) error {
	for name, fun := range funcs {
		if err := rs.RegisterFunction(name, fun); err != nil {
			return err
		}
	}

	return nil
}

func (rs *NamespaceBuilder) RegisteredFunctions() []string {
	res := make([]string, 0, len(rs.funcs))

	// root namespace, return all functions
	if rs.name == emptyNS {
		for k := range rs.funcs {
			res = append(res, k)
		}
	} else {
		nsPrefix := rs.name + separator
		for k := range rs.funcs {
			if strings.HasPrefix(k, nsPrefix) {
				res = append(res, k)
			}
		}
	}

	return res
}

func (rs *NamespaceBuilder) makeFullName(name string) string {
	if rs.name == emptyNS {
		return name
	}

	return rs.name + separator + name
}
