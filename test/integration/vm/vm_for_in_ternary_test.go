package vm_test

import (
	"testing"

	"github.com/MontFerret/ferret/v2/test/spec"
	. "github.com/MontFerret/ferret/v2/test/spec/exec"
)

func TestForTernaryExpression(t *testing.T) {
	RunSpecs(t, []spec.Spec{
		Array(`
			LET foo = FALSE
			RETURN foo ? TRUE : (FOR i IN 1..5 RETURN i*2)`,
			[]any{2, 4, 6, 8, 10}),
		Array(`
			LET foo = FALSE
			RETURN foo ? TRUE : (FOR i IN 1..5 T::FAIL() RETURN i*2)?`,
			[]any{}),
		Array(`
			LET foo = FALSE
			RETURN foo ? (FOR i IN 1..5 RETURN i) : (FOR i IN 1..5 RETURN i*2)`,
			[]any{2, 4, 6, 8, 10}),
		Array(`
			LET foo = FALSE
			RETURN foo ? (FOR i IN 1..5 RETURN T::FAIL()) : (FOR i IN 1..5 RETURN T::FAIL())?`,
			[]any{}),
		Array(`
			LET foo = TRUE
			RETURN foo ? (FOR i IN 1..5 RETURN T::FAIL())? : (FOR i IN 1..5 RETURN T::FAIL())`,
			[]any{}),
		Array(`
			LET foo = FALSE
			LET res = foo ? TRUE : (FOR i IN 1..5 RETURN i*2) 
			RETURN res`,
			[]any{2, 4, 6, 8, 10}),
		S(`
			LET foo = TRUE
			LET res = foo ? TRUE : (FOR i IN 1..5 RETURN i*2)
			RETURN res`,
			true),
		Array(`
			LET foo = FALSE
			LET res = foo ? TRUE : (FOR i IN 1..5 RETURN i*2) 
			RETURN res`,
			[]any{2, 4, 6, 8, 10}),
		Array(`
			LET foo = FALSE
			LET res = foo ? (FOR i IN 1..5 RETURN i) : (FOR i IN 1..5 RETURN i*2)
			RETURN res`,
			[]any{2, 4, 6, 8, 10}),
		Array(`
			LET foo = TRUE
			LET res = foo ? (FOR i IN 1..5 RETURN i) : (FOR i IN 1..5 RETURN i*2)
			RETURN res`,
			[]any{1, 2, 3, 4, 5}),
		S(`
			LET res = LENGTH((FOR i IN 1..5 RETURN T::FAIL())?) ? TRUE : FALSE
			RETURN res`,
			false),
		S(`
			LET res = (FOR i IN 1..5 RETURN i)? ? TRUE : FALSE
			RETURN res
`,
			true),
	})
}
