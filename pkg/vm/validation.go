package vm

import (
	"strings"

	"github.com/MontFerret/ferret/pkg/runtime"
)

func validate(env *Environment, program *Program) error {
	if err := validateParams(env, program); err != nil {
		return err
	}

	if err := validateFunctions(env, program); err != nil {
		return err
	}

	return nil
}

func validateParams(env *Environment, program *Program) error {
	if len(program.Params) == 0 {
		return nil
	}

	// There might be no errors.
	// Thus, we allocate this slice lazily, on a first error.
	var missedParams []string

	for _, n := range program.Params {
		_, exists := env.params[n]

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

// TODO: Implement this function.
func validateFunctions(env *Environment, program *Program) error {
	//if len(program.Locations) == 0 {
	//	return nil
	//}
	//
	//// There might be no errors.
	//// Thus, we allocate this slice lazily, on a first error.
	//var missedFunctions []string
	//
	//for _, loc := range program.Locations {
	//	if loc.Function == "" {
	//		continue
	//	}
	//
	//	if _, exists := env.functions[loc.Function]; !exists {
	//		if missedFunctions == nil {
	//			missedFunctions = make([]string, 0, len(program.Locations))
	//		}
	//
	//		missedFunctions = append(missedFunctions, loc.Function)
	//	}
	//}
	//
	//if len(missedFunctions) > 0 {
	//	return runtime.Error(ErrFunctionNotFound, strings.Join(missedFunctions, ", "))
	//}
	//
	return nil
}
