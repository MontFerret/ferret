package compiler

import (
	"strings"

	"github.com/pkg/errors"

	"github.com/MontFerret/ferret/pkg/runtime/core"
)

func copyFromNamespace(fns *core.Functions, namespace string) error {
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
