package compiler_test

import (
	"fmt"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/compiler"
	"github.com/MontFerret/ferret/v2/test/spec"
	. "github.com/MontFerret/ferret/v2/test/spec/compile"
)

func TestDispatchShorthandCompiles(t *testing.T) {
	RunSpecs(t, []spec.Spec{
		Opcode(`
			@d <- "click"
			RETURN 1
		`, OpcodeExistence{
			Exists: []bytecode.Opcode{bytecode.OpDispatch},
		}, "Should compile shorthand dispatch as a statement"),
		Opcode(`
			RETURN @d <- "click"
		`, OpcodeExistence{
			Exists: []bytecode.Opcode{bytecode.OpDispatch},
		}, "Should compile shorthand dispatch as an expression"),
		Opcode(`
			LET target = @d
			RETURN target<-"click"
		`, OpcodeExistence{
			Exists: []bytecode.Opcode{bytecode.OpDispatch},
		}, "Should compile compact shorthand dispatch with string event"),
		Opcode(`
			LET target = @d
			LET eventName = "submit"
			RETURN target<-eventName
		`, OpcodeExistence{
			Exists: []bytecode.Opcode{bytecode.OpDispatch},
		}, "Should compile compact shorthand dispatch with variable event"),
		Opcode(`
			LET a = @d
			LET b = "click"
			RETURN a<-b
		`, OpcodeExistence{
			Exists: []bytecode.Opcode{bytecode.OpDispatch},
		}, "Should treat compact a<-b as shorthand dispatch"),
		Opcode(`
			LET a = 1
			LET b = 2
			RETURN a<(-b)
		`, OpcodeExistence{
			Exists:    []bytecode.Opcode{bytecode.OpLt},
			NotExists: []bytecode.Opcode{bytecode.OpDispatch},
		}, "Should compile parenthesized negative variable comparison without dispatch"),
		Opcode(`
			LET a = 1
			RETURN a<(-1)
		`, OpcodeExistence{
			Exists:    []bytecode.Opcode{bytecode.OpLt},
			NotExists: []bytecode.Opcode{bytecode.OpDispatch},
		}, "Should compile parenthesized negative literal comparison without dispatch"),
		Opcode(`
			LET button = @d
			button <- "click"
			LET form = @d
			form <- "submit"
			RETURN 1
		`, OpcodeExistence{
			Exists: []bytecode.Opcode{bytecode.OpDispatch},
		}, "Should compile variable-target shorthand dispatch"),
		Opcode("LET document = @d\n"+
			"document[~ css`button`][0] <- \"click\"\n"+
			"RETURN 1", OpcodeExistence{
			Exists: []bytecode.Opcode{bytecode.OpDispatch},
		}, "Should compile member-target shorthand dispatch"),
	})
}

func TestDispatchInForBodiesCompiles(t *testing.T) {
	RunSpecs(t, []spec.Spec{
		Opcode(`
			VAR i = 0
			FOR WHILE i < 1
				i += 1
				DISPATCH "click" IN @d
				RETURN i
		`, OpcodeExistence{
			Exists: []bytecode.Opcode{bytecode.OpDispatch},
		}, "Should compile long-form dispatch in FOR WHILE body"),
		Opcode(`
			VAR i = 0
			FOR WHILE i < 1
				i += 1
				@d <- "click"
				RETURN i
		`, OpcodeExistence{
			Exists: []bytecode.Opcode{bytecode.OpDispatch},
		}, "Should compile shorthand dispatch in FOR WHILE body"),
		Opcode(`
			FOR item IN [1]
				DISPATCH "click" IN @d
				RETURN item
		`, OpcodeExistence{
			Exists: []bytecode.Opcode{bytecode.OpDispatch},
		}, "Should compile long-form dispatch in FOR IN body"),
	})
}

func TestDispatchGroupedTargetsCompile(t *testing.T) {
	RunSpecsLevels(t, []spec.Spec{
		ProgramCheck(`
			LET page = WEB::HTML::OPEN("https://mockery.ferretlang.org/scenarios/forms/", {
				driver: "cdp"
			})

			DISPATCH "input" IN (QUERY ONE "#query" IN page USING css) WITH { value: "ferret" }
			DISPATCH "click" IN (QUERY ONE "#search-form button[type='submit']" IN page USING css)

			LET result = WAITFOR VALUE (QUERY ONE "#form-result" IN page USING css)
				WHEN .attributes.disabled == false
				TIMEOUT 10s

			RETURN result.textContent
		`, expectGroupedDispatchOrder, "Should compile grouped query dispatch targets in the reported query"),
		Opcode(`
			DISPATCH "input" IN (QUERY ONE "#query" IN @page USING css) WITH { value: "ferret" }
			RETURN 1
		`, OpcodeCount{Count: map[bytecode.Opcode]int{
			bytecode.OpQueryOne: 1,
			bytecode.OpDispatch: 1,
		}}, "Should compile a grouped query target with a payload"),
		Opcode(`
			DISPATCH "click" IN (QUERY ONE "#submit" IN @page USING css)
			RETURN 1
		`, OpcodeCount{Count: map[bytecode.Opcode]int{
			bytecode.OpQueryOne: 1,
			bytecode.OpDispatch: 1,
		}}, "Should compile a grouped query target without a payload"),
		Opcode(`
			DISPATCH "click" IN (MATCH @kind (
				"primary" => @primary,
				_ => @fallback,
			))
			RETURN 1
		`, OpcodeExistence{
			Exists: []bytecode.Opcode{bytecode.OpDispatch},
		}, "Should compile a non-query grouped expression target"),
		Opcode(`
			(QUERY ONE "#submit" IN @page USING css) <- "click"
			RETURN 1
		`, OpcodeCount{Count: map[bytecode.Opcode]int{
			bytecode.OpQueryOne: 1,
			bytecode.OpDispatch: 1,
		}}, "Should compile a grouped shorthand target"),
		Opcode(`
			LET obs = []
			RETURN WAITFOR EVENT "ready" IN obs
				TRIGGER DISPATCH "click" IN (QUERY ONE "#submit" IN @page USING css)
				TIMEOUT 1ms
				ON TIMEOUT RETURN NONE
		`, OpcodeCount{Count: map[bytecode.Opcode]int{
			bytecode.OpQueryOne: 1,
			bytecode.OpDispatch: 1,
		}}, "Should compile a grouped target in an inline WAITFOR trigger"),
	}, compiler.O0, compiler.O1)
}

func expectGroupedDispatchOrder(program *bytecode.Program) error {
	var queries []int
	var dispatches []int

	for i, instruction := range program.Bytecode {
		switch instruction.Opcode {
		case bytecode.OpQueryOne:
			queries = append(queries, i)
		case bytecode.OpDispatch:
			dispatches = append(dispatches, i)
		}
	}

	if len(queries) != 3 {
		return fmt.Errorf("expected 3 QUERY ONE instructions, got %d", len(queries))
	}

	if len(dispatches) != 2 {
		return fmt.Errorf("expected 2 DISPATCH instructions, got %d", len(dispatches))
	}

	if !(queries[0] < dispatches[0] && dispatches[0] < queries[1] && queries[1] < dispatches[1] && dispatches[1] < queries[2]) {
		return fmt.Errorf("unexpected QUERY ONE/DISPATCH order: queries=%v dispatches=%v", queries, dispatches)
	}

	return nil
}
