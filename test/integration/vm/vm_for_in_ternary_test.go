package vm_test

import (
	"testing"
)

func TestForTernaryExpression(t *testing.T) {
	RunUseCases(t, []UseCase{
		CaseArray(`
			LET foo = FALSE
			RETURN foo ? TRUE : (FOR i IN 1..5 RETURN i*2)`,
			[]any{2, 4, 6, 8, 10}),
		CaseArray(`
			LET foo = FALSE
			RETURN foo ? TRUE : (FOR i IN 1..5 T::FAIL() RETURN i*2)?`,
			[]any{}),
		CaseArray(`
			LET foo = FALSE
			RETURN foo ? (FOR i IN 1..5 RETURN i) : (FOR i IN 1..5 RETURN i*2)`,
			[]any{2, 4, 6, 8, 10}),
		CaseArray(`
			LET foo = FALSE
			RETURN foo ? (FOR i IN 1..5 RETURN T::FAIL()) : (FOR i IN 1..5 RETURN T::FAIL())?`,
			[]any{}),
		CaseArray(`
			LET foo = TRUE
			RETURN foo ? (FOR i IN 1..5 RETURN T::FAIL())? : (FOR i IN 1..5 RETURN T::FAIL())`,
			[]any{}),
		CaseArray(`
			LET foo = FALSE
			LET res = foo ? TRUE : (FOR i IN 1..5 RETURN i*2) 
			RETURN res`,
			[]any{2, 4, 6, 8, 10}),
		Case(`
			LET foo = TRUE
			LET res = foo ? TRUE : (FOR i IN 1..5 RETURN i*2)
			RETURN res`,
			true),
		CaseArray(`
			LET foo = FALSE
			LET res = foo ? TRUE : (FOR i IN 1..5 RETURN i*2) 
			RETURN res`,
			[]any{2, 4, 6, 8, 10}),
		CaseArray(`
			LET foo = FALSE
			LET res = foo ? (FOR i IN 1..5 RETURN i) : (FOR i IN 1..5 RETURN i*2)
			RETURN res`,
			[]any{2, 4, 6, 8, 10}),
		CaseArray(`
			LET foo = TRUE
			LET res = foo ? (FOR i IN 1..5 RETURN i) : (FOR i IN 1..5 RETURN i*2)
			RETURN res`,
			[]any{1, 2, 3, 4, 5}),
		Case(`
			LET res = LENGTH((FOR i IN 1..5 RETURN T::FAIL())?) ? TRUE : FALSE
			RETURN res`,
			false),
		Case(`
			LET res = (FOR i IN 1..5 RETURN i)? ? TRUE : FALSE
			RETURN res
`,
			true),
	})
}
