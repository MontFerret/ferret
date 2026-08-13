package sdk

import (
	"fmt"
	"strings"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

type functionRegistrationKey struct {
	name  string
	arity byte
}

// RegisterFunctions validates and registers definitions in a namespace.
// Qualified function names, including namespace segments, are canonicalized to lowercase.
// Functions may share a name when their fixed arities differ or one is variadic.
// Validation is atomic: no definition is registered when any definition is invalid.
func RegisterFunctions(ns runtime.Namespace, definitions ...FunctionDef) error {
	if ns == nil {
		return fmt.Errorf("function namespace cannot be nil")
	}

	functions := ns.Function()
	seen := make(map[functionRegistrationKey]struct{}, len(definitions))

	for i := range definitions {
		name := strings.TrimSpace(definitions[i].name)
		if !runtime.HasTerminalFunctionName(name) {
			return fmt.Errorf("function name cannot be empty")
		}

		definitions[i].name = name
	}

	for _, definition := range definitions {
		arity, err := validateFunctionDefinition(definition)
		if err != nil {
			return err
		}

		key := functionRegistrationKey{name: runtime.NormalizeRegisteredName(definition.name), arity: arity}
		if _, exists := seen[key]; exists {
			return fmt.Errorf("function %q with the same arity is defined more than once", definition.name)
		}

		if hasFunctionDefinition(functions, definition.name, arity) {
			return fmt.Errorf("function %q with the same arity is already registered", definition.name)
		}

		seen[key] = struct{}{}
	}

	for _, definition := range definitions {
		registerFunctionDefinition(functions, definition)
	}

	return nil
}

func validateFunctionDefinition(definition FunctionDef) (byte, error) {
	switch fn := definition.function.(type) {
	case runtime.Function:
		if fn == nil {
			return 0, fmt.Errorf("function %q cannot be nil", definition.name)
		}
		return 0xff, nil
	case runtime.Function0:
		if fn == nil {
			return 0, fmt.Errorf("function %q cannot be nil", definition.name)
		}
		return 0, nil
	case runtime.Function1:
		if fn == nil {
			return 0, fmt.Errorf("function %q cannot be nil", definition.name)
		}
		return 1, nil
	case runtime.Function2:
		if fn == nil {
			return 0, fmt.Errorf("function %q cannot be nil", definition.name)
		}
		return 2, nil
	case runtime.Function3:
		if fn == nil {
			return 0, fmt.Errorf("function %q cannot be nil", definition.name)
		}
		return 3, nil
	case runtime.Function4:
		if fn == nil {
			return 0, fmt.Errorf("function %q cannot be nil", definition.name)
		}
		return 4, nil
	default:
		return 0, fmt.Errorf("function %q has unsupported type %T", definition.name, definition.function)
	}
}

func hasFunctionDefinition(functions runtime.FunctionDefs, name string, arity byte) bool {
	switch arity {
	case 0:
		return functions.A0().Has(name)
	case 1:
		return functions.A1().Has(name)
	case 2:
		return functions.A2().Has(name)
	case 3:
		return functions.A3().Has(name)
	case 4:
		return functions.A4().Has(name)
	case 0xff:
		return functions.Var().Has(name)
	default:
		return false
	}
}

func registerFunctionDefinition(functions runtime.FunctionDefs, definition FunctionDef) {
	switch fn := definition.function.(type) {
	case runtime.Function:
		functions.Var().Add(definition.name, fn)
	case runtime.Function0:
		functions.A0().Add(definition.name, fn)
	case runtime.Function1:
		functions.A1().Add(definition.name, fn)
	case runtime.Function2:
		functions.A2().Add(definition.name, fn)
	case runtime.Function3:
		functions.A3().Add(definition.name, fn)
	case runtime.Function4:
		functions.A4().Add(definition.name, fn)
	}
}
