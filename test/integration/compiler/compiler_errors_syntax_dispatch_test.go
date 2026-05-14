package compiler_test

import (
	"strings"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/compiler"
	parserd "github.com/MontFerret/ferret/v2/pkg/parser/diagnostics"
	"github.com/MontFerret/ferret/v2/pkg/source"
	"github.com/MontFerret/ferret/v2/test/spec"
	. "github.com/MontFerret/ferret/v2/test/spec/compile"
)

func TestSyntaxErrorsDispatch(t *testing.T) {
	RunSpecs(t, []spec.Spec{
		Failure(`
			LET obj = NONE
			LET ok = DISPATCH IN obj
			RETURN ok
		`, E{
			Kind: parserd.SyntaxError,
		}, "Missing DISPATCH event name"),
		Failure(`
			LET obj = NONE
			LET ok = DISPATCH "click" IN
			RETURN ok
		`, E{
			Kind: parserd.SyntaxError,
		}, "Missing DISPATCH target"),
		Failure(`
			LET obj = NONE
			LET ok = DISPATCH "click" IN obj WITH
			RETURN ok
		`, E{
			Kind: parserd.SyntaxError,
		}, "Missing DISPATCH payload"),
		Failure(`
			LET obj = NONE
			LET ok = DISPATCH "click" IN obj OPTIONS
			RETURN ok
		`, E{
			Kind: parserd.SyntaxError,
		}, "Missing DISPATCH options"),
		Failure(`
			LET obj = NONE
			LET ok = <- "click"
			RETURN ok
		`, E{
			Kind:    parserd.SyntaxError,
			Message: "Expected dispatch target before '<-'",
			Hint:    `Provide a dispatchable target, e.g. btn <- "click".`,
		}, "Missing shorthand dispatch target"),
		Failure(`
			RETURN <- "click"
		`, E{
			Kind:    parserd.SyntaxError,
			Message: "Expected dispatch target before '<-'",
			Hint:    `Provide a dispatchable target, e.g. btn <- "click".`,
		}, "Missing shorthand dispatch target after RETURN"),
		Failure(`
			LET ok = (<- "click")
			RETURN ok
		`, E{
			Kind:    parserd.SyntaxError,
			Message: "Expected dispatch target before '<-'",
			Hint:    `Provide a dispatchable target, e.g. btn <- "click".`,
		}, "Missing shorthand dispatch target in parenthesized expression"),
		Failure(`
			LET obj = NONE
			LET ok = obj <-
			RETURN ok
		`, E{
			Kind:    parserd.SyntaxError,
			Message: "Expected dispatch event after '<-'",
			Hint:    `Provide an event expression, e.g. btn <- "click".`,
		}, "Missing shorthand dispatch event"),
		Failure(`
			LET a = @d
			RETURN a<-1
		`, E{
			Kind:    parserd.SyntaxError,
			Message: "Expected dispatch event after '<-'",
			Hint:    `Provide an event expression, e.g. btn <- "click".`,
		}, "Numeric compact shorthand event should fail as dispatch syntax"),
		Failure(`
			LET obj = NONE
			LET ok = obj <- "input" WITH { value: "x" }
			RETURN ok
		`, E{
			Kind:    parserd.SyntaxError,
			Message: "Dispatch shorthand does not support WITH",
			Hint:    `Use the long form instead, e.g. DISPATCH "input" IN field WITH { value: "x" }.`,
		}, "Shorthand WITH should fail syntax checks"),
		Failure(`
			LET obj = NONE
			LET ok = obj <- "click" OPTIONS { bubbles: true }
			RETURN ok
		`, E{
			Kind:    parserd.SyntaxError,
			Message: "Dispatch shorthand does not support OPTIONS",
			Hint:    `Use the long form instead, e.g. DISPATCH "click" IN btn OPTIONS { bubbles: true }.`,
		}, "Shorthand OPTIONS should fail syntax checks"),
		Failure(`
			LET obj = NONE
			LET ok = "click" -> obj
			RETURN ok
		`, E{
			Kind: parserd.SyntaxError,
		}, "Old event-first shorthand should fail as plain syntax"),
	})
}

func TestDispatchSyntaxErrorsIgnoreCommentsAndStrings(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name string
		src  string
	}{
		{
			name: "string literal",
			src: `
				LET msg = "RETURN <-"
				LET x =
				RETURN x
			`,
		},
		{
			name: "single-line comment",
			src: `
				LET x = 1 // RETURN <-
				LET y =
				RETURN y
			`,
		},
		{
			name: "multi-line comment",
			src: `
				/* = <- */
				LET y =
				RETURN y
			`,
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			_, err := compiler.New(compiler.WithOptimizationLevel(compiler.O0)).Compile(source.New("dispatch_diag", tc.src))
			if err == nil {
				t.Fatal("expected compilation error")
			}

			diag := firstCompilationError(err)
			if diag == nil {
				t.Fatal("expected diagnostic")
			}

			if diag.Message == "Expected dispatch target before '<-'" {
				t.Fatalf("unexpected dispatch shorthand diagnostic for %s: %q", tc.name, diag.Message)
			}

			if strings.Contains(diag.Hint, `btn <- "click"`) {
				t.Fatalf("unexpected dispatch shorthand hint for %s: %q", tc.name, diag.Hint)
			}
		})
	}
}
