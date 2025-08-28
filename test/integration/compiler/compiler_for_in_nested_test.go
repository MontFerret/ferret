package compiler_test

import (
	"testing"

	"github.com/MontFerret/ferret/pkg/vm"
)

func TestForNested(t *testing.T) {
	RunUseCases(t, []UseCase{
		SkipByteCodeCase(
			`FOR prop IN ["a"]
							FOR val IN [1, 2, 3]
								RETURN {[prop]: val}`,
			BC{
				I(vm.OpReturn, 0, 7),
			},
		),
		SkipByteCodeCase(
			`FOR val IN 1..3
							FOR prop IN ["a"]
								RETURN {[prop]: val}`,
			BC{
				I(vm.OpReturn, 0, 7),
			},
		),
		SkipByteCodeCase(
			`FOR prop IN ["a"]
							FOR val IN 1..3
								RETURN {[prop]: val}`,
			BC{
				I(vm.OpReturn, 0, 7),
			},
		),
		SkipByteCodeCase(
			`FOR prop IN ["a"]
							FOR val IN [1, 2, 3]
								FOR val2 IN [1, 2, 3]
									RETURN { [prop]: [val, val2] }`,
			BC{
				I(vm.OpReturn, 0, 7),
			},
		),
		SkipByteCodeCase(
			`FOR val IN [1, 2, 3]
							RETURN (
								FOR prop IN ["a", "b", "c"]
									RETURN { [prop]: val }
							)`,
			BC{
				I(vm.OpReturn, 0, 7),
			},
		),
		SkipByteCodeCase(
			`FOR val IN [1, 2, 3]
							LET sub = (
								FOR prop IN ["a", "b", "c"]
									RETURN { [prop]: val }
							)
		
							RETURN sub`,
			BC{
				I(vm.OpReturn, 0, 7),
			},
		),
		SkipByteCodeCase(`
LET strs = ["foo", "bar", "qaz", "abc"]

FOR s IN strs
	SORT s
	FOR n IN 0..1
		RETURN CONCAT(s, n)
`,
			BC{
				I(vm.OpReturn, 0, 7),
			},
		),
		SkipByteCodeCase(`
LET strs = ["foo", "bar", "qaz", "abc"]

FOR n IN 0..1
	FOR s IN strs
		SORT s
		RETURN CONCAT(s, n)
`,
			BC{
				I(vm.OpReturn, 0, 7),
			},
		),
	})
}
