package compiler_test

import (
	"testing"

	parserd "github.com/MontFerret/ferret/v2/pkg/parser/diagnostics"
	"github.com/MontFerret/ferret/v2/test/spec"
	. "github.com/MontFerret/ferret/v2/test/spec/compile"
)

func TestMathOperatorsRejectKnownNonNumericOperands(t *testing.T) {
	expected := func(operator string) E {
		return E{
			Kind:    parserd.SemanticError,
			Message: "Operator '" + operator + "' requires numeric operands",
			Hint:    "Use Int or Float values with this operator.",
		}
	}

	RunSpecs(t, []spec.Spec{
		Failure(`RETURN [120, 45, 300] * 2`, expected("*"), "array multiplication"),
		Failure(`RETURN "3" * 2`, expected("*"), "string multiplication"),
		Failure(`RETURN TRUE - 1`, expected("-"), "boolean subtraction"),
		Failure(`RETURN {} / 2`, expected("/"), "object division"),
		Failure(`RETURN -"x"`, expected("-"), "string unary negative"),
		Failure(`RETURN +"3"`, expected("+"), "string unary positive"),
	})
}
