package compiler_test

import (
	"testing"

	parserd "github.com/MontFerret/ferret/v2/pkg/parser/diagnostics"
)

func TestImplicitCurrentShorthandErrors(t *testing.T) {
	RunUseCases(t, []UseCase{
		ErrorCase(
			`
			RETURN .name
		`, E{
				Kind: parserd.SyntaxError,
			}, "Implicit current shorthand in return"),
		ErrorCase(
			`
			RETURN ?.name
		`, E{
				Kind: parserd.SyntaxError,
			}, "Implicit current shorthand optional chaining in return"),
		ErrorCase(
			`
			RETURN .[0]
		`, E{
				Kind: parserd.SyntaxError,
			}, "Implicit current shorthand computed access in return"),
		ErrorCase(
			`
			RETURN .[*]
		`, E{
				Kind: parserd.SyntaxError,
			}, "Implicit current shorthand array expansion in return"),
		ErrorCase(
			`
			RETURN .[**]
		`, E{
				Kind: parserd.SyntaxError,
			}, "Implicit current shorthand array contraction in return"),
		ErrorCase(
			`
			RETURN .[? FILTER .active]
		`, E{
				Kind: parserd.SyntaxError,
			}, "Implicit current shorthand array question mark in return"),
		ErrorCase(
			"RETURN .[~ css`x`]",
			E{
				Kind: parserd.SyntaxError,
			}, "Implicit current shorthand query apply in return"),
		ErrorCase(
			`
			LET x = .name
			RETURN x
		`, E{
				Kind: parserd.SyntaxError,
			}, "Implicit current shorthand in LET assignment"),
		ErrorCase(
			`
			RETURN [.name]
		`, E{
				Kind: parserd.SyntaxError,
			}, "Implicit current shorthand in array literal"),
		ErrorCase(
			`
			RETURN {name: .name}
		`, E{
				Kind: parserd.SyntaxError,
			}, "Implicit current shorthand in object literal"),
		ErrorCase(
			`
			RETURN UPPER(.name)
		`, E{
				Kind: parserd.SyntaxError,
			}, "Implicit current shorthand in function call"),
		ErrorCase(
			`
			LET doc = {name: "x"}
			RETURN doc[.name]
		`, E{
				Kind: parserd.SyntaxError,
			}, "Implicit current shorthand in computed member expression"),
		ErrorCase(
			`
			RETURN .name..10
		`, E{
				Kind: parserd.SyntaxError,
			}, "Implicit current shorthand in range expression"),
		ErrorCase(
			`
			LET docs = [{active: true}]
			FOR d IN docs FILTER .active RETURN d
		`, E{
				Kind: parserd.SyntaxError,
			}, "Implicit current shorthand in FOR FILTER clause"),
		ErrorCase(
			`
			LET docs = [{name: "a"}]
			FOR d IN docs SORT .name RETURN d
		`, E{
				Kind: parserd.SyntaxError,
			}, "Implicit current shorthand in FOR SORT clause"),
		ErrorCase(
			`
			LET docs = [{limit: 1}]
			FOR d IN docs LIMIT .limit RETURN d
		`, E{
				Kind: parserd.SyntaxError,
			}, "Implicit current shorthand in FOR LIMIT clause"),
		ErrorCase(
			`
			LET docs = [{group: 1}]
			FOR d IN docs COLLECT g = .group RETURN g
		`, E{
				Kind: parserd.SyntaxError,
			}, "Implicit current shorthand in FOR COLLECT clause"),
		ErrorCase(
			`
			LET ok = WAITFOR EXISTS .name
			RETURN ok
		`, E{
				Kind: parserd.SyntaxError,
			}, "Implicit current shorthand in WAITFOR predicate"),
	})
}
