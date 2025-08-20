package runtime

import (
	"strings"
)

//var fnNameValidation = regexp.MustCompile("^[a-zA-Z]+[a-zA-Z0-9_]*(::[a-zA-Z]+[a-zA-Z0-9_]*)*$")

const NamespaceSeparator = "::"
const emptyNS = ""

type (
	Namespace interface {
		Functions() Functions
		Namespace(name string) Namespace
		RegisterFunctions(funs Functions) error
		RegisteredFunctions() []string
		RemoveFunction(name string)
	}

	nsContainer struct {
		funcs Functions
		name  string
	}
)

func NewRootNamespace() Namespace {
	ns := new(nsContainer)
	ns.funcs = NewFunctions()

	return ns
}

func NewNamespace(funcs Functions, name string) Namespace {
	return &nsContainer{funcs, strings.ToUpper(name)}
}

func (nc *nsContainer) Functions() Functions {
	return nc.funcs
}

func (nc *nsContainer) Namespace(name string) Namespace {
	return NewNamespace(nc.funcs, nc.makeFullName(name))
}

func (nc *nsContainer) RemoveFunction(name string) {
	nc.funcs.Unset(nc.makeFullName(name))
}

func (nc *nsContainer) RegisterFunctions(funcs Functions) error {
	var err error

	err = funcs.F().ForEach(func(function Function, s string) error {
		fullName := nc.makeFullName(s)

		if nc.funcs.F().Has(fullName) {
			return Errorf(ErrInvalidOperation, "function '%s' is already registered", fullName)
		}

		nc.funcs.F().Set(nc.makeFullName(s), function)

		return nil
	})

	if err != nil {
		return err
	}

	err = funcs.F0().ForEach(func(function Function0, s string) error {
		fullName := nc.makeFullName(s)

		if nc.funcs.F0().Has(fullName) {
			return Errorf(ErrInvalidOperation, "function '%s' is already registered", fullName)
		}

		nc.funcs.F0().Set(nc.makeFullName(s), function)

		return nil
	})

	if err != nil {
		return err
	}

	err = funcs.F1().ForEach(func(function Function1, s string) error {
		fullName := nc.makeFullName(s)

		if nc.funcs.F1().Has(fullName) {
			return Errorf(ErrInvalidOperation, "function '%s' is already registered", fullName)
		}

		nc.funcs.F1().Set(nc.makeFullName(s), function)

		return nil
	})

	if err != nil {
		return err
	}
	err = funcs.F2().ForEach(func(function Function2, s string) error {
		fullName := nc.makeFullName(s)

		if nc.funcs.F2().Has(fullName) {
			return Errorf(ErrInvalidOperation, "function '%s' is already registered", fullName)
		}

		nc.funcs.F2().Set(nc.makeFullName(s), function)

		return nil
	})

	if err != nil {
		return err
	}

	err = funcs.F3().ForEach(func(function Function3, s string) error {
		fullName := nc.makeFullName(s)

		if nc.funcs.F3().Has(fullName) {
			return Errorf(ErrInvalidOperation, "function '%s' is already registered", fullName)
		}

		nc.funcs.F3().Set(nc.makeFullName(s), function)

		return nil
	})

	if err != nil {
		return err
	}

	err = funcs.F4().ForEach(func(function Function4, s string) error {
		fullName := nc.makeFullName(s)

		if nc.funcs.F4().Has(fullName) {
			return Errorf(ErrInvalidOperation, "function '%s' is already registered", fullName)
		}

		nc.funcs.F4().Set(nc.makeFullName(s), function)

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

func (nc *nsContainer) RegisteredFunctions() []string {
	fnames := nc.funcs.Names()
	res := make([]string, 0, len(fnames))

	// root namespace, return all functions
	if nc.name == emptyNS {
		res = append(res, fnames...)
	} else {
		nsPrefix := nc.name + NamespaceSeparator
		for _, k := range fnames {
			if strings.HasPrefix(k, nsPrefix) {
				res = append(res, k)
			}
		}
	}

	return res
}

func (nc *nsContainer) makeFullName(name string) string {
	name = strings.ToUpper(name)

	if nc.name == emptyNS {
		return name
	}

	return nc.name + NamespaceSeparator + name
}
