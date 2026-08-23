package compiler_test

import (
	"fmt"
	"slices"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/compiler"
	"github.com/MontFerret/ferret/v2/pkg/source"
)

func TestBracedBlockSyntax(t *testing.T) {
	valid := map[string]string{
		"arrow UDF":                       `FUNC value() => 1 RETURN value()`,
		"block UDF":                       `FUNC value() { RETURN 1 } RETURN value()`,
		"returned UDF FOR":                `FUNC values() { RETURN FOR value IN [{ n: 1 }] { RETURN value.n } } RETURN values()`,
		"object MATCH pattern and result": `RETURN MATCH { kind: "ok" } { { kind: "ok" } => { value: 1 }, _ => {} }`,
		"guard MATCH":                     `LET value = 1 RETURN MATCH { WHEN value > 0 => { positive: true }, _ => {} }`,
		"braced IN":                       `FOR value IN { one: 1 } { RETURN value }`,
		"braced WHILE":                    `FOR value WHILE value < 1 { RETURN value }`,
		"braced DO WHILE":                 `FOR value DO WHILE false { RETURN value }`,
		"returnless braced IN":            `FOR value IN [1] { LET copy = value }`,
		"returnless braced WHILE":         `FOR WHILE false {}`,
		"returnless braced DO WHILE":      `FOR DO WHILE false {}`,
		"nested returnless statement":     `FOR outer IN [1] { FOR inner IN [outer] {} LET copy = outer }`,
		"nested mixed FOR":                `FOR outer IN [1] { LET inner = (FOR value IN [outer] RETURN value) RETURN inner }`,
	}

	for name, query := range valid {
		t.Run(name, func(t *testing.T) {
			if _, err := mustNewCompiler(t).Compile(source.NewAnonymous(query)); err != nil {
				t.Fatalf("compile valid syntax: %v", err)
			}
		})
	}
}

func TestRemovedBlockSyntaxIsRejected(t *testing.T) {
	invalid := map[string]string{
		"parenthesized UDF":   `FUNC value() ( RETURN 1 ) RETURN value()`,
		"parenthesized MATCH": `RETURN MATCH 1 (1 => 1, _ => 0)`,
		"unbraced MATCH":      `RETURN MATCH 1 1 => 1, _ => 0`,
		"unbraced returnless": `FOR value IN [1] LET copy = value`,
	}

	for name, query := range invalid {
		t.Run(name, func(t *testing.T) {
			if _, err := mustNewCompiler(t).Compile(source.NewAnonymous(query)); err == nil {
				t.Fatal("expected syntax error")
			}
		})
	}
}

func TestForBodyBracesDoNotChangeBytecode(t *testing.T) {
	const unbraced = `
RETURN FOR value IN [1, 2, 3]
  FILTER value > 1
  RETURN value * 2
`
	const braced = `
RETURN FOR value IN [1, 2, 3] {
  FILTER value > 1
  RETURN value * 2
}
`

	for _, level := range []compiler.OptimizationLevel{compiler.O0, compiler.O1} {
		t.Run(fmt.Sprintf("O%d", level), func(t *testing.T) {
			unbracedProgram := compileWithLevel(t, level, unbraced)
			bracedProgram := compileWithLevel(t, level, braced)

			if !slices.Equal(unbracedProgram.Bytecode, bracedProgram.Bytecode) {
				t.Fatalf("instruction/operand sequences differ:\nunbraced: %v\nbraced:   %v", unbracedProgram.Bytecode, bracedProgram.Bytecode)
			}
		})
	}
}
