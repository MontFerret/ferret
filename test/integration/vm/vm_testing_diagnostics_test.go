package vm_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/compiler"
	"github.com/MontFerret/ferret/v2/pkg/source"
	"github.com/MontFerret/ferret/v2/test/spec"
	. "github.com/MontFerret/ferret/v2/test/spec/exec"
)

func TestTestingAssertionDiagnostics(t *testing.T) {
	tests := []struct {
		name        string
		description string
	}{
		{name: "equal_scalar_context", description: "scalar equality includes custom context"},
		{name: "equal_structured_paths", description: "structured equality reports deterministic paths"},
		{name: "not_equal_structured", description: "negated equality reports the shared structured value"},
		{name: "has_missing_keys", description: "has reports every missing key in input order"},
		{name: "match_expression", description: "match reports the regular expression operand"},
	}

	specs := make([]spec.Spec, 0, len(tests))
	for _, test := range tests {
		specs = append(specs, testingDiagnosticSpec(t, test.name, test.description))
	}

	RunSpecs(t, specs)
}

func testingDiagnosticSpec(t *testing.T, name, description string) spec.Spec {
	t.Helper()

	directory := filepath.Join("testdata", "testing")
	scriptPath := filepath.Join(directory, name+".fql")
	goldenPath := filepath.Join(directory, name+".golden")

	scriptContent, err := os.ReadFile(scriptPath)
	if err != nil {
		t.Fatalf("read assertion diagnostic script %q: %v", scriptPath, err)
	}

	golden, err := os.ReadFile(goldenPath)
	if err != nil {
		t.Fatalf("read assertion diagnostic golden %q: %v", goldenPath, err)
	}

	scriptName := filepath.ToSlash(scriptPath)
	script := source.New(scriptName, string(scriptContent))
	input := spec.NewProgramSourceInput(spec.ProgramSource{
		Name: scriptName,
		Build: func(_ string, c *compiler.Compiler) (*bytecode.Program, error) {
			return c.Compile(script)
		},
	})

	return spec.NewSpecWith(input, description).Expect().ExecError(
		ShouldBeRuntimeError,
		&ExpectedRuntimeError{
			Message: "assertion error",
			Format:  string(golden),
		},
	)
}
