package compiler_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/compiler"
	"github.com/MontFerret/ferret/v2/pkg/diagnostics"
	"github.com/MontFerret/ferret/v2/pkg/file"
)

func TestCollectAggregateRequiresAtLeastOneArgument(t *testing.T) {
	query := `
		LET users = [1, 2, 3]
		FOR u IN users
			COLLECT AGGREGATE total = COUNT()
			RETURN total
	`

	c := compiler.New(compiler.WithOptimizationLevel(compiler.O0))
	_, err := c.Compile(file.NewSource("collect-aggregate-arity", query))

	if err == nil {
		t.Fatal("expected compilation error")
	}

	ce := firstCompilationError(err)
	if ce == nil {
		t.Fatalf("expected compiler error type, got %T", err)
	}

	if ce.Kind != diagnostics.Kind("SemanticError") {
		t.Fatalf("expected SemanticError, got %s", ce.Kind)
	}

	if !strings.Contains(ce.Message, "requires at least one argument") {
		t.Fatalf("expected arity message, got %q", ce.Message)
	}

	if !strings.Contains(ce.Message, "COUNT") {
		t.Fatalf("expected function name in message, got %q", ce.Message)
	}

	formatted := ce.Format()
	if strings.Contains(formatted, "goroutine ") || strings.Contains(formatted, "unhandled panic") {
		t.Fatalf("diagnostic should not include panic stack trace, got:\n%s", formatted)
	}
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

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
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

			c := compiler.New(compiler.WithOptimizationLevel(compiler.O0))
			program, err := c.Compile(file.NewSource("collect-aggregate-fused-literal", query))
			if err != nil {
				t.Fatalf("unexpected compilation error: %v", err)
			}

			if !hasAggregatePlan(program) {
				t.Fatalf("expected grouped fused aggregate plan for key expression %q", tc.key)
			}

			if hasFunctionCallOpcode(program) {
				t.Fatalf("expected fused grouped aggregation to avoid function call opcodes for key expression %q", tc.key)
			}
		})
	}
}

func firstCompilationError(err error) *diagnostics.Diagnostic {
	switch e := err.(type) {
	case *diagnostics.Diagnostic:
		return e
	case *diagnostics.DiagnosticSet:
		if e.Size() == 0 {
			return nil
		}

		return e.Errors()[0]
	default:
		return nil
	}
}

func hasAggregatePlan(program *bytecode.Program) bool {
	return len(program.Metadata.AggregatePlans) > 0
}

func hasFunctionCallOpcode(program *bytecode.Program) bool {
	for _, instruction := range program.Bytecode {
		switch instruction.Opcode {
		case bytecode.OpCall,
			bytecode.OpProtectedCall,
			bytecode.OpCall0,
			bytecode.OpProtectedCall0,
			bytecode.OpCall1,
			bytecode.OpProtectedCall1,
			bytecode.OpCall2,
			bytecode.OpProtectedCall2,
			bytecode.OpCall3,
			bytecode.OpProtectedCall3,
			bytecode.OpCall4,
			bytecode.OpProtectedCall4:
			return true
		}
	}

	return false
}
