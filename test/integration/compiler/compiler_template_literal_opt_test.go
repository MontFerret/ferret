package compiler_test

import (
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/compiler"
	"github.com/MontFerret/ferret/v2/pkg/file"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/vm"
	"github.com/MontFerret/ferret/v2/test/integration/base"
)

func compileNoOpt(t *testing.T, expr string) *bytecode.Program {
	t.Helper()
	c := compiler.New(compiler.WithOptimizationLevel(compiler.O0))
	prog, err := c.Compile(file.NewSource("template-literal-opt", expr))
	if err != nil {
		t.Fatalf("compile failed: %v", err)
	}
	return prog
}

func execProgram(t *testing.T, prog *bytecode.Program) any {
	t.Helper()
	out, err := base.Exec(prog, false, vm.WithFunctions(base.Stdlib()))
	if err != nil {
		t.Fatalf("exec failed: %v", err)
	}
	return out
}

func execProgramWithOpts(t *testing.T, prog *bytecode.Program, opts ...vm.EnvironmentOption) any {
	t.Helper()

	allOpts := append([]vm.EnvironmentOption{vm.WithFunctions(base.Stdlib())}, opts...)
	out, err := base.Exec(prog, false, allOpts...)
	if err != nil {
		t.Fatalf("exec failed: %v", err)
	}
	return out
}

func countOpcode(prog *bytecode.Program, op bytecode.Opcode) int {
	if prog == nil {
		return 0
	}
	count := 0
	for _, inst := range prog.Bytecode {
		if inst.Opcode == op {
			count++
		}
	}
	return count
}

func assertOpcodeCount(t *testing.T, prog *bytecode.Program, op bytecode.Opcode, want int) {
	t.Helper()
	if got := countOpcode(prog, op); got != want {
		t.Fatalf("expected %d %s opcode(s), got %d", want, op.String(), got)
	}
}

func hasEmptyStringConstant(prog *bytecode.Program) bool {
	if prog == nil {
		return false
	}

	for _, c := range prog.Constants {
		s, ok := c.(runtime.String)
		if !ok {
			continue
		}

		if s == runtime.EmptyString {
			return true
		}
	}

	return false
}

func TestTemplateLiteral_ConstantFolding(t *testing.T) {
	testCases := []struct {
		name     string
		query    string
		opConcat int
		opAdd    int
		expected string
	}{
		{
			name:     "fully constant template",
			query:    "RETURN `foo-${1}-bar-${true}`",
			opConcat: 0,
			opAdd:    0,
			expected: "foo-1-bar-true",
		},
		{
			name:     "folds constant expressions into chunks",
			query:    "LET x = \"X\" RETURN `a-${1}-b-${x}-c-${true}-d`",
			opConcat: 1,
			opAdd:    0,
			expected: "a-1-b-X-c-true-d",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			prog := compileNoOpt(t, tc.query)
			assertOpcodeCount(t, prog, bytecode.OpConcat, tc.opConcat)
			assertOpcodeCount(t, prog, bytecode.OpAdd, tc.opAdd)
			out := execProgram(t, prog)
			str, ok := out.(string)
			if !ok {
				t.Fatalf("expected string output, got %T", out)
			}
			if str != tc.expected {
				t.Fatalf("expected %q, got %v", tc.expected, str)
			}
		})
	}
}

func TestTemplateLiteral_DynamicEdgeCases_NoEmptyConstArtifacts(t *testing.T) {
	testCases := []struct {
		name           string
		query          string
		opts           []vm.EnvironmentOption
		opConcat       int
		expected       string
		wantEmptyConst bool
	}{
		{
			name:           "single param interpolation",
			query:          "RETURN `${@foo}`",
			opts:           []vm.EnvironmentOption{vm.WithParam("foo", runtime.NewString("bar"))},
			opConcat:       1,
			expected:       "bar",
			wantEmptyConst: false,
		},
		{
			name:           "adjacent param interpolations",
			query:          "RETURN `${@a}${@b}`",
			opts:           []vm.EnvironmentOption{vm.WithParam("a", runtime.NewString("x")), vm.WithParam("b", runtime.NewString("y"))},
			opConcat:       1,
			expected:       "xy",
			wantEmptyConst: false,
		},
		{
			name:           "prefix literal and param",
			query:          "RETURN `pre-${@foo}`",
			opts:           []vm.EnvironmentOption{vm.WithParam("foo", runtime.NewString("bar"))},
			opConcat:       1,
			expected:       "pre-bar",
			wantEmptyConst: false,
		},
		{
			name:           "escaped interpolation marker constant folds",
			query:          "RETURN `cost=\\${1}`",
			opts:           nil,
			opConcat:       0,
			expected:       "cost=${1}",
			wantEmptyConst: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			prog := compileNoOpt(t, tc.query)
			assertOpcodeCount(t, prog, bytecode.OpConcat, tc.opConcat)

			if got := hasEmptyStringConstant(prog); got != tc.wantEmptyConst {
				t.Fatalf("unexpected empty-string constant presence: got %v, want %v", got, tc.wantEmptyConst)
			}

			out := execProgramWithOpts(t, prog, tc.opts...)
			str, ok := out.(string)
			if !ok {
				t.Fatalf("expected string output, got %T", out)
			}
			if str != tc.expected {
				t.Fatalf("expected %q, got %v", tc.expected, str)
			}
		})
	}
}
