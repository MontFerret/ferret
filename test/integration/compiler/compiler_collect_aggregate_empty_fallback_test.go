package compiler_test

import (
	"fmt"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/test/spec"
	. "github.com/MontFerret/ferret/v2/test/spec/compile"
)

func TestCollectAggregateGlobalEmptyFallbackLoadsNone(t *testing.T) {
	RunSpecs(t, []spec.Spec{
		ProgramCheck(`
LET users = []
FOR u IN users
	COLLECT AGGREGATE total = COUNT(u)
	RETURN total
`, func(prog *bytecode.Program) error {
			if !hasOpcode(prog.Bytecode, bytecode.OpJumpIfTrue) {
				return fmt.Errorf("expected OpJumpIfTrue for empty-aggregator branch")
			}

			if got := countOpcode(prog, bytecode.OpLoadNone); got < 1 {
				return fmt.Errorf("expected OpLoadNone in global aggregate empty fallback, got %d", got)
			}

			return nil
		}, "global aggregate empty fallback loads none"),
	})
}
