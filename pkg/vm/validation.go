package vm

import (
	"strings"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

func validate(env *Environment, program *bytecode.Program) error {
	if err := validateProgramVersion(program); err != nil {
		return err
	}

	return validateParams(env, program)
}

func validateProgramVersion(program *bytecode.Program) error {
	if program == nil {
		return runtime.Error(runtime.ErrInvalidOperation, "unsupported bytecode version; recompile query")
	}

	if program.ISAVersion != bytecode.Version {
		return runtime.Error(runtime.ErrInvalidOperation, "unsupported bytecode version; recompile query")
	}

	return nil
}

func validateParams(env *Environment, program *bytecode.Program) error {
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
