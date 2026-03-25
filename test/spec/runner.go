package spec

import (
	"fmt"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/compiler"
	"github.com/MontFerret/ferret/v2/pkg/diagnostics"
	"github.com/MontFerret/ferret/v2/pkg/file"
	"github.com/MontFerret/ferret/v2/pkg/vm"
)

type Runner struct {
	Name        string
	Compiler    *compiler.Compiler
	EnvOpts     []vm.EnvironmentOption
	SpecEnvOpts func(index int) []vm.EnvironmentOption
}

func (r *Runner) Run(t *testing.T, specs []Spec) {
	t.Helper()

	std := Stdlib()

	for index, spec := range specs {
		specName := spec.String()

		t.Run(fmt.Sprintf("%s/%s", r.Name, specName), func(t *testing.T) {
			t.Helper()

			if spec.ShouldSkip {
				t.Skip()
			}

			var prog *bytecode.Program

			defer func() {
				if recovered := recover(); recovered != nil {
					PrintDebug(t, specName, prog)
					t.Fatalf("panic: %v", recovered)
				}
			}()

			prog, err := r.Compiler.Compile(file.NewSource(specName, spec.Expression))
			if err != nil {
				if spec.DebugOutput {
					PrintError(t, err)
				}

				if spec.Compile.ErrorAssertion != nil {
					spec.Compile.ErrorAssertion(t, err, spec.Compile.Error)
					return
				}

				t.Fatalf("unexpected compilation error:\n%s", diagnostics.Format(err))
			}

			if spec.Compile.ErrorAssertion != nil || spec.Compile.Error != nil {
				t.Fatal("expected compilation error, got none")
			}

			if spec.Compile.ValueAssertion != nil {
				spec.Compile.ValueAssertion(t, prog, spec.Compile.Value)
			}

			exec := spec.Run.ErrorAssertion != nil || spec.Run.ValueAssertion != nil

			if !exec {
				return
			}

			options := []vm.EnvironmentOption{
				vm.WithNamespace(std),
			}

			if len(r.EnvOpts) > 0 {
				options = append(options, r.EnvOpts...)
			}

			if r.SpecEnvOpts != nil {
				specEnvOpts := r.SpecEnvOpts(index)
				options = append(options, specEnvOpts...)
			}

			actual, err := Exec(prog, spec.RawOutput, options...)

			if err != nil {
				if spec.Run.ErrorAssertion != nil {
					spec.Run.ErrorAssertion(t, err, spec.Run.Error)
					return
				}

				t.Fatalf("unexpected runtime error: %v", err)
			}

			if spec.Run.ErrorAssertion != nil || spec.Run.Error != nil {
				t.Fatal("expected runtime error, got none")
			}

			if spec.Run.ValueAssertion != nil {
				spec.Run.ValueAssertion(t, actual, spec.Run.Value)
			}

			if spec.DebugOutput {
				PrintDebug(t, specName, prog)
			}
		})
	}
}
