package sdk

import (
	"fmt"
	"strings"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

// RegisterFunctions validates and registers definitions in a namespace.
// Validation is atomic: no definition is registered when any definition is invalid.
func RegisterFunctions(ns runtime.Namespace, definitions ...FunctionDef) error {
	if ns == nil {
		return fmt.Errorf("function namespace cannot be nil")
	}

	functions := ns.Function()
	seen := make(map[string]struct{}, len(definitions))

	for i := range definitions {
		definitions[i].name = strings.TrimSpace(definitions[i].name)
		if definitions[i].name == "" {
			return fmt.Errorf("function name cannot be empty")
		}
	}

	for _, definition := range definitions {
		if _, exists := seen[definition.name]; exists {
			return fmt.Errorf("function %q is defined more than once", definition.name)
		}

		if functions.Has(definition.name) {
			return fmt.Errorf("function %q is already registered", definition.name)
		}

		if err := validateFunctionDefinition(definition); err != nil {
			return err
		}

		seen[definition.name] = struct{}{}
	}

	for _, definition := range definitions {
		registerFunctionDefinition(functions, definition)
	}

	return nil
}

func validateFunctionDefinition(definition FunctionDef) error {
	switch fn := definition.function.(type) {
	case runtime.Function:
		if fn == nil {
			return fmt.Errorf("function %q cannot be nil", definition.name)
		}
	case runtime.Function0:
		if fn == nil {
			return fmt.Errorf("function %q cannot be nil", definition.name)
		}
	case runtime.Function1:
		if fn == nil {
			return fmt.Errorf("function %q cannot be nil", definition.name)
		}
	case runtime.Function2:
		if fn == nil {
			return fmt.Errorf("function %q cannot be nil", definition.name)
		}
	case runtime.Function3:
		if fn == nil {
			return fmt.Errorf("function %q cannot be nil", definition.name)
		}
	case runtime.Function4:
		if fn == nil {
			return fmt.Errorf("function %q cannot be nil", definition.name)
		}
	default:
		return fmt.Errorf("function %q has unsupported type %T", definition.name, definition.function)
	}

	return nil
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
