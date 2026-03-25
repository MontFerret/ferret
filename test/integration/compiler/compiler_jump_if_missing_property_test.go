package compiler_test

import (
	"fmt"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/test/spec"
)

func TestMatchLoadPropertyConstEmission(t *testing.T) {
	RunSpecs(t, []spec.Spec{
		ProgramCheck(`
LET obj = { a: 1, b: 2 }
RETURN MATCH obj (
  { a: 1, b: v } => v,
  _ => 0,
)
`, func(prog *bytecode.Program) error {
			if !programHasOpcode(prog, bytecode.OpMatchLoadPropertyConst) {
				return fmt.Errorf("expected bytecode to contain %s", bytecode.OpMatchLoadPropertyConst)
			}

			if programHasOpcode(prog, bytecode.OpJumpIfMissingPropertyConst) {
				return fmt.Errorf("expected object-pattern lowering to avoid %s", bytecode.OpJumpIfMissingPropertyConst)
			}

			if programHasOpcode(prog, bytecode.OpLoadPropertyConst) {
				return fmt.Errorf("expected object-pattern lowering to avoid %s", bytecode.OpLoadPropertyConst)
			}

			if len(prog.Metadata.MatchFailTargets) != len(prog.Bytecode) {
				return fmt.Errorf("expected match fail targets metadata to align with bytecode: got %d, want %d", len(prog.Metadata.MatchFailTargets), len(prog.Bytecode))
			}

			return nil
		}, "object pattern uses match-specific property load"),
	})
}
