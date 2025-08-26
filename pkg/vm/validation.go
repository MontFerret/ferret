package vm

import (
	"fmt"
	"strings"

	"github.com/MontFerret/ferret/pkg/runtime"
)

func validate(env *Environment, program *Program) error {
	if err := validateParams(env, program); err != nil {
		return err
	}

	return validateFunctions(env, program)
}

func validateParams(env *Environment, program *Program) error {
	if len(program.Params) == 0 {
		return nil
	}

	// There might be no errors.
	// Thus, we allocate this slice lazily, on a first error.
	var missedParams []string

	for _, n := range program.Params {
		_, exists := env.Params[n]

		if !exists {
			if missedParams == nil {
				missedParams = make([]string, 0, len(program.Params))
			}

			missedParams = append(missedParams, "@"+n)
		}
	}

	if len(missedParams) > 0 {
		return runtime.Error(ErrMissedParam, strings.Join(missedParams, ", "))
	}

	return nil
}

func validateFunctions(env *Environment, program *Program) error {
	if len(program.Functions) == 0 {
		return nil
	}

	// There might be no errors.
	// Thus, we allocate this slice lazily, on a first error.
	var errors []string

	for name, args := range program.Functions {
		exists := env.Functions.Has(name)

		if !exists {
			if errors == nil {
				errors = make([]string, 0, len(program.Functions))
			}

			errors = append(errors, fmt.Sprintf("function `%s` not found", name))

			continue
		}

		// Check if the number of arguments matches.
		var matched bool

		switch args {
		case 4:
			matched = env.Functions.F4().Has(name)
		case 3:
			matched = env.Functions.F3().Has(name)
		case 2:
			matched = env.Functions.F2().Has(name)
		case 1:
			matched = env.Functions.F1().Has(name)
		case 0:
			matched = env.Functions.F0().Has(name)
		default:
			// Variable number of arguments.
			matched = env.Functions.F().Has(name)
		}

		// Check if the function is a variadic function.
		if !matched && args > -1 {
			matched = env.Functions.F().Has(name)
		}

		if !matched {
			// Tell the user that the function was not found with the specified number of arguments.
			errors = append(errors, fmt.Sprintf("function `%s` not found with %d arguments", name, args))
		}
	}

	if len(errors) > 0 {
		return runtime.Error(ErrFunctionNotFound, strings.Join(errors, ", "))
	}

	return nil
}
