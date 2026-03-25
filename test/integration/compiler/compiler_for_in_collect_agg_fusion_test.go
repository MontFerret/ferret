package compiler_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	parserd "github.com/MontFerret/ferret/v2/pkg/parser/diagnostics"
	"github.com/MontFerret/ferret/v2/test/spec"
	"github.com/MontFerret/ferret/v2/test/spec/assert"
)

func TestCollectAggregateRequiresAtLeastOneArgument(t *testing.T) {
	RunSpecs(t, []spec.Spec{
		spec.New(`
		LET users = [1, 2, 3]
		FOR u IN users
			COLLECT AGGREGATE total = COUNT()
			RETURN total
	`).Expect().CompileError(assert.NewUnaryAssertion(func(actual any) error {
			err, ok := actual.(error)
			if !ok {
				return fmt.Errorf("expected error, got %T", actual)
			}

			ce := firstCompilationError(err)
			if ce == nil {
				return fmt.Errorf("expected compiler error type, got %T", err)
			}

			if ce.Kind != parserd.SemanticError {
				return fmt.Errorf("expected SemanticError, got %s", ce.Kind)
			}

			if !strings.Contains(ce.Message, "requires at least one argument") {
				return fmt.Errorf("expected arity message, got %q", ce.Message)
			}

			if !strings.Contains(ce.Message, "COUNT") {
				return fmt.Errorf("expected function name in message, got %q", ce.Message)
			}

			formatted := ce.Format()
			if strings.Contains(formatted, "goroutine ") || strings.Contains(formatted, "unhandled panic") {
				return fmt.Errorf("diagnostic should not include panic stack trace, got:\n%s", formatted)
			}

			return nil
		}), "aggregate requires at least one arg"),
	})
}

func TestCollectAggregateGroupedFusionSupportsScalarLiteralKeys(t *testing.T) {
	testCases := []struct {
		name string
		key  string
	}{
		{name: "integer", key: "1"},
		{name: "float", key: "1.5"},
		{name: "boolean", key: "true"},
		{name: "none", key: "NONE"},
	}

	specs := make([]spec.Spec, 0, len(testCases))
	for _, tc := range testCases {
		query := fmt.Sprintf(`
				LET users = [{ age: 1 }, { age: 2 }, { age: 3 }]
				FOR u IN users
					COLLECT g = %s
					AGGREGATE
						cnt = COUNT(u.age),
						sum = SUM(u.age),
						min = MIN(u.age)
					RETURN { g, cnt, sum, min }
			`, tc.key)
		specs = append(specs, ProgramCheck(query, func(program *bytecode.Program) error {
			if !hasAggregatePlan(program) {
				return fmt.Errorf("expected grouped fused aggregate plan for key expression %q", tc.key)
			}

			if hasFunctionCallOpcode(program) {
				return fmt.Errorf("expected fused grouped aggregation to avoid function call opcodes for key expression %q", tc.key)
			}

			return nil
		}, tc.name))
	}

	RunSpecs(t, specs)
}

func TestCollectAggregateGroupedFusionUsesAggregateGroupUpdateOpcode(t *testing.T) {
	RunSpecs(t, []spec.Spec{
		ProgramCheck(`
LET users = [{ age: 1 }, { age: 2 }, { age: 3 }]

FOR u IN users
	COLLECT g = u.age
	AGGREGATE
		cnt = COUNT(u.age),
		sum = SUM(u.age),
		min = MIN(u.age)
	RETURN { g, cnt, sum, min }
`, func(prog *bytecode.Program) error {
			if !hasAggregatePlan(prog) {
				return fmt.Errorf("expected grouped fused aggregate plan")
			}

			if !hasOpcode(prog.Bytecode, bytecode.OpAggregateGroupUpdate) {
				return fmt.Errorf("expected grouped fused aggregation to use OpAggregateGroupUpdate")
			}

			if hasOpcode(prog.Bytecode, bytecode.OpPushKV) {
				return fmt.Errorf("expected grouped fused aggregation without INTO to avoid raw PushKV group writes")
			}

			return nil
		}, "grouped fusion uses aggregate group update"),
	})
}

func TestCollectAggregateGroupedFusionSupportsComputedKeys(t *testing.T) {
	RunSpecs(t, []spec.Spec{
		ProgramCheck(`
LET users = [{ age: 1 }, { age: 2 }, { age: 3 }]

FOR u IN users
	COLLECT g = u.age % 2
	AGGREGATE
		cnt = COUNT(u.age),
		sum = SUM(u.age),
		min = MIN(u.age)
	RETURN { g, cnt, sum, min }
`, func(prog *bytecode.Program) error {
			if !hasAggregatePlan(prog) {
				return fmt.Errorf("expected grouped fused aggregate plan for computed group key")
			}

			if !hasOpcode(prog.Bytecode, bytecode.OpAggregateGroupUpdate) {
				return fmt.Errorf("expected computed-key grouped fusion to use OpAggregateGroupUpdate")
			}

			return nil
		}, "computed group key fusion"),
	})
}

func TestCollectAggregateGroupedFusionWithIntoKeepsGroupValueWrites(t *testing.T) {
	RunSpecs(t, []spec.Spec{
		ProgramCheck(`
LET users = [{ age: 1 }, { age: 2 }, { age: 3 }]

FOR u IN users
	COLLECT g = u.age
	AGGREGATE
		cnt = COUNT(u.age),
		sum = SUM(u.age),
		min = MIN(u.age)
	INTO groups
	RETURN { g, cnt, sum, min, groups }
`, func(prog *bytecode.Program) error {
			if !hasOpcode(prog.Bytecode, bytecode.OpAggregateGroupUpdate) {
				return fmt.Errorf("expected grouped aggregate INTO to use OpAggregateGroupUpdate")
			}

			if !hasOpcode(prog.Bytecode, bytecode.OpPushKV) {
				return fmt.Errorf("expected grouped aggregate INTO to keep raw group value writes")
			}

			return nil
		}, "grouped aggregate into keeps group writes"),
	})
}

func hasAggregatePlan(program *bytecode.Program) bool {
	return len(program.Metadata.AggregatePlans) > 0
}

func hasFunctionCallOpcode(program *bytecode.Program) bool {
	for _, instruction := range program.Bytecode {
		switch instruction.Opcode {
		case bytecode.OpHCall,
			bytecode.OpProtectedHCall,
			bytecode.OpCall,
			bytecode.OpProtectedCall,
			bytecode.OpTailCall:
			return true
		}
	}

	return false
}
