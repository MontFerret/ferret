package spec

import (
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
		suiteName := spec.SuiteName(r.Name)

		t.Run(suiteName, func(t *testing.T) {
			if spec.Base.SkipInfo.Active {
				t.Skip(spec.Base.SkipInfo.Reason)
			}

			var prog *bytecode.Program

			defer func() {
				if recovered := recover(); recovered != nil {
					PrintDebug(t, suiteName, prog)
					t.Fatalf("panic: %v", recovered)
				}
			}()

			prog, err := r.Compiler.Compile(file.NewSource(suiteName, spec.Base.Expression))
			if err != nil {
				if spec.Base.DebugOutput {
					PrintError(t, err)
				}

				if spec.Compile.Error.Defined() {
					spec.Compile.Error.Assert(t, err)
					return
				}

				t.Fatalf("unexpected compilation error:\n%s", diagnostics.Format(err))
			}

			if spec.Compile.Error.Defined() {
				t.Fatal("expected compilation error, got none")
			}

			if spec.Compile.Result.Defined() {
				spec.Compile.Result.Assert(t, prog)
			}

			if !spec.Exec.Result.Defined() && !spec.Exec.Error.Defined() {
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

			actual, err := Exec(prog, spec.Exec.RawOutput, options...)

			if err != nil {
				if spec.Exec.Error.Defined() {
					spec.Exec.Error.Assert(t, err)
					return
				}

				t.Fatalf("unexpected runtime error: %v", err)
			}

			if spec.Exec.Error.Defined() {
				t.Fatal("expected runtime error, got none")
			}

			if spec.Exec.Result.Defined() {
				spec.Exec.Result.Assert(t, actual)
			}

			if spec.Base.DebugOutput {
				PrintDebug(t, suiteName, prog)
			}
		})
	}
}
