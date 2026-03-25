package vm_test

import (
	. "github.com/MontFerret/ferret/v2/test/spec/exec"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/vm"
)

func TestMatchObjectPattern(t *testing.T) {
	RunSpecs(t, []Spec{
		S(
			`
LET obj = @obj
RETURN MATCH obj (
  { a: 1, b: v } => v,
  _ => 0,
)
`,
			2,
			"Should bind a matched object property",
		),
	}, vm.WithParams(map[string]runtime.Value{
		"obj": runtime.NewObjectWith(map[string]runtime.Value{
			"a": runtime.NewInt(1),
			"b": runtime.NewInt(2),
		}),
	}))

	RunSpecs(t, []Spec{
		S(
			`
LET obj = @obj
RETURN MATCH obj (
  { a: 1, b: v } => v,
  _ => 0,
)
`,
			0,
			"Should fall through when a required property is missing",
		),
	}, vm.WithParams(map[string]runtime.Value{
		"obj": runtime.NewObjectWith(map[string]runtime.Value{
			"a": runtime.NewInt(1),
		}),
	}))
}
