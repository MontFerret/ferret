package compiler_test

import (
	"testing"

	parserd "github.com/MontFerret/ferret/v2/pkg/parser/diagnostics"
	"github.com/MontFerret/ferret/v2/test/spec"
	. "github.com/MontFerret/ferret/v2/test/spec/compile"
)

func TestImplicitCurrentShorthandErrors(t *testing.T) {
	RunSpecs(t, []spec.Spec{
		Failure(
			`
			RETURN .name
		`, E{
				Kind: parserd.SyntaxError,
			}, "Implicit current shorthand in return"),
		Failure(
			`
			RETURN ?.name
		`, E{
				Kind: parserd.SyntaxError,
			}, "Implicit current shorthand optional chaining in return"),
		Failure(
			`
			RETURN .[0]
		`, E{
				Kind: parserd.SyntaxError,
			}, "Implicit current shorthand computed access in return"),
		Failure(
			`
			RETURN .[*]
		`, E{
				Kind: parserd.SyntaxError,
			}, "Implicit current shorthand array expansion in return"),
		Failure(
			`
			RETURN .[**]
		`, E{
				Kind: parserd.SyntaxError,
			}, "Implicit current shorthand array contraction in return"),
		Failure(
			`
			RETURN .[? FILTER .active]
		`, E{
				Kind: parserd.SyntaxError,
			}, "Implicit current shorthand array question mark in return"),
		Failure(
			"RETURN .[~ css`x`]",
			E{
				Kind: parserd.SyntaxError,
			}, "Implicit current shorthand query apply in return"),
		Failure(
			`
			LET x = .name
			RETURN x
		`, E{
				Kind: parserd.SyntaxError,
			}, "Implicit current shorthand in LET assignment"),
		Failure(
			`
			RETURN [.name]
		`, E{
				Kind: parserd.SyntaxError,
			}, "Implicit current shorthand in array literal"),
		Failure(
			`
			RETURN {name: .name}
		`, E{
				Kind: parserd.SyntaxError,
			}, "Implicit current shorthand in object literal"),
		Failure(
			`
			RETURN UPPER(.name)
		`, E{
				Kind: parserd.SyntaxError,
			}, "Implicit current shorthand in function call"),
		Failure(
			`
			LET doc = {name: "x"}
			RETURN doc[.name]
		`, E{
				Kind: parserd.SyntaxError,
			}, "Implicit current shorthand in computed member expression"),
		Failure(
			`
			RETURN .name..10
		`, E{
				Kind: parserd.SyntaxError,
			}, "Implicit current shorthand in range expression"),
		Failure(
			`
			LET docs = [{active: true}]
			FOR d IN docs FILTER .active RETURN d
		`, E{
				Kind: parserd.SyntaxError,
			}, "Implicit current shorthand in FOR FILTER clause"),
		Failure(
			`
			LET docs = [{name: "a"}]
			FOR d IN docs SORT .name RETURN d
		`, E{
				Kind: parserd.SyntaxError,
			}, "Implicit current shorthand in FOR SORT clause"),
		Failure(
			`
			LET docs = [{limit: 1}]
			FOR d IN docs LIMIT .limit RETURN d
		`, E{
				Kind: parserd.SyntaxError,
			}, "Implicit current shorthand in FOR LIMIT clause"),
		Failure(
			`
			LET docs = [{group: 1}]
			FOR d IN docs COLLECT g = .group RETURN g
		`, E{
				Kind: parserd.SyntaxError,
			}, "Implicit current shorthand in FOR COLLECT clause"),
		Failure(
			`
			LET ok = WAITFOR EXISTS .name
			RETURN ok
		`, E{
				Kind: parserd.SyntaxError,
			}, "Implicit current shorthand in WAITFOR predicate"),
	})
}
