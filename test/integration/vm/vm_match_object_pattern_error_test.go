package vm_test

import (
	"context"
	"errors"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/vm"
	"github.com/MontFerret/ferret/v2/test/spec"
	. "github.com/MontFerret/ferret/v2/test/spec/exec"
)

type errorMap struct {
	*runtime.Object
}

func (m *errorMap) ContainsKey(ctx context.Context, key runtime.Value) (runtime.Boolean, error) {
	return runtime.False, errors.New("boom")
}

func TestMatchObjectPatternContainsKeyError(t *testing.T) {
	RunSpecs(t, []spec.Spec{
		S(
			`
LET obj = @obj
RETURN MATCH obj (
  { a: 1 } => 1,
  _ => 0,
)
`,
			"Should surface ContainsKey errors",
		).Expect().ExecError(ShouldBeRuntimeError, &ExpectedRuntimeError{Contains: []string{"boom"}}),
	}, vm.WithParams(map[string]runtime.Value{
		"obj": &errorMap{
			Object: runtime.NewObjectWith(map[string]runtime.Value{
				"a": runtime.NewInt(1),
			}),
		},
	}))
}
