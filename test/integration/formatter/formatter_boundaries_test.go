package formatter_test

import (
	"context"
	"errors"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/formatter"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/vm"
	"github.com/MontFerret/ferret/v2/test/spec"
)

type grammarBoundaryCase struct {
	name   string
	input  string
	want   string
	result string
}

func TestFormatterGrammarBoundaries(t *testing.T) {
	const fallback = `"Default value returned when no shipping option matches the destination address"`
	const function = "func f() => 1\n"

	tests := []grammarBoundaryCase{
		{
			name:   "current as left endpoint",
			input:  "return [1][* return . .. 10]",
			want:   "return [1][* return . .. 10]",
			result: "[[1,2,3,4,5,6,7,8,9,10]]",
		},
		{
			name:   "current as right endpoint",
			input:  "return [1][* return 1 .. .]",
			want:   "return [1][* return 1 .. .]",
			result: "[[1]]",
		},
		{
			name:   "current as both endpoints",
			input:  "return [1][* return . .. .]",
			want:   "return [1][* return . .. .]",
			result: "[[1]]",
		},
		{
			name:   "implicit member as right endpoint",
			input:  "let row={name:3} return [row][* return 1 .. .name]",
			want:   "let row = { name: 3 }\nreturn [row][* return 1 .. .name]",
			result: "[[1,2,3]]",
		},
		{
			name:   "numeric endpoints",
			input:  "return 1 .. 3",
			want:   "return 1..3",
			result: "[1,2,3]",
		},
		{
			name:   "variable endpoints",
			input:  "let lo=1 let hi=3 return lo .. hi",
			want:   "let lo = 1\nlet hi = 3\nreturn lo..hi",
			result: "[1,2,3]",
		},
		{
			name:   "recovery return before current endpoint",
			input:  function + "return [2][* return f() on error return (.) .. .]",
			want:   function + "return [2][* return f() on error return (.) .. .]",
			result: "[[1,2]]",
		},
		{
			name:   "grouped current recovery value before numeric endpoint",
			input:  function + "return [2][* return f() on error return (.) .. 3]",
			want:   function + "return [2][* return f() on error return (.)..3]",
			result: "[[1,2,3]]",
		},
		{
			name:   "numeric recovery value before outer range",
			input:  function + "return f() on error return (1) .. 3",
			want:   function + "return f() on error return (1)..3",
			result: "[1,2,3]",
		},
		{
			name:   "nested recovery grouping",
			input:  function + "return f() on error return (((1))) .. 3",
			want:   function + "return f() on error return (1)..3",
			result: "[1,2,3]",
		},
		{
			name:   "group ends a compound recovery value",
			input:  function + "return f() on error return (2) + (1) .. 3",
			want:   function + "return f() on error return 2 + (1)..3",
			result: "[1,2,3]",
		},
		{
			name:   "recovery on a member access before outer range",
			input:  "let o={value:1} return o.value on error return (2) .. 3",
			want:   "let o = { value: 1 }\nreturn o.value on error return (2)..3",
			result: "[1,2,3]",
		},
		{
			name:   "waitfor recovery value before outer range",
			input:  function + "return f() on error return (WAITFOR VALUE 1) .. 3",
			want:   function + "return f() on error return (waitfor value 1)..3",
			result: "[1,2,3]",
		},
		{
			name:   "recovery value at right endpoint",
			input:  function + "return 1 .. f() on error return (3)",
			want:   function + "return 1..f() on error return 3",
			result: "[1]",
		},
		{
			name:   "recovery without outer range",
			input:  function + "return f() on error return (1)",
			want:   function + "return f() on error return 1",
			result: "1",
		},
		{
			name:   "function return",
			input:  "func read(DISTINCT){return (DISTINCT.values)} return read({values:[1,2]})",
			want:   "func read(DISTINCT) {\n    return (DISTINCT.values)\n}\nreturn read({ values: [1, 2] })",
			result: "[1,2]",
		},
		{
			name:   "nested grouping in top-level return",
			input:  "let dIsTiNcT={values:[1,1,2]} return (((dIsTiNcT.values)))",
			want:   "let dIsTiNcT = { values: [1, 1, 2] }\nreturn (dIsTiNcT.values)",
			result: "[1,1,2]",
		},
		{
			name:   "comment inside required parentheses",
			input:  "let DISTINCT={values:2} return ( /* keep */ DISTINCT.values /* here */ )",
			want:   "let DISTINCT = { values: 2 }\nreturn (/* keep */ DISTINCT.values /* here */)",
			result: "2",
		},
		{
			name:   "nested member source",
			input:  "let DISTINCT={values:2} return ((DISTINCT)).values",
			want:   "let DISTINCT = { values: 2 }\nreturn (DISTINCT).values",
			result: "2",
		},
		{
			name:   "comment supplies the member source boundary",
			input:  "let DISTINCT={x:1} return ( /* keep */ (DISTINCT) ).x",
			want:   "let DISTINCT = { x: 1 }\nreturn (/* keep */ DISTINCT).x",
			result: "1",
		},
		{
			name:   "index expression could compile as distinct array literal",
			input:  "let DISTINCT=[10,20] return (DISTINCT[1])",
			want:   "let DISTINCT = [10, 20]\nreturn (DISTINCT[1])",
			result: "20",
		},
		{
			name:   "group at start of binary expression",
			input:  "let DISTINCT=2 return (DISTINCT)+1",
			want:   "let DISTINCT = 2\nreturn (DISTINCT) + 1",
			result: "3",
		},
		{
			name:   "existing precedence boundary",
			input:  "let DISTINCT=2 return (((DISTINCT) + 1) * 2)",
			want:   "let DISTINCT = 2\nreturn (DISTINCT + 1) * 2",
			result: "6",
		},
		{
			name:   "wrapped coalesce keeps its leading group after the inline layout attempt",
			input:  "let DISTINCT=2 return ((DISTINCT)) ?? " + fallback,
			want:   "let DISTINCT = 2\nreturn (DISTINCT)\n    ?? " + fallback,
			result: "2",
		},
		{
			name:   "ordinary return",
			input:  "let value=2 return ((value))",
			want:   "let value = 2\nreturn value",
			result: "2",
		},
		{
			name:   "distinct in nonleading operand",
			input:  "let DISTINCT=2 return 1+(DISTINCT)",
			want:   "let DISTINCT = 2\nreturn 1 + DISTINCT",
			result: "3",
		},
		{
			name:   "explicit return modifier",
			input:  "let DISTINCT={values:[1,1,2]} RETURN DISTINCT (DISTINCT.values)",
			want:   "let DISTINCT = { values: [1, 1, 2] }\nreturn distinct DISTINCT.values",
			result: "[1,2]",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.check(t, nil)
		})
	}

	t.Run("failing call uses the recovery range endpoint", func(t *testing.T) {
		calls := 0
		fail := vm.WithFunction("STEP", func(context.Context, ...runtime.Value) (runtime.Value, error) {
			calls++
			return runtime.None, errors.New("step failed")
		})
		tt := grammarBoundaryCase{
			input:  "return STEP() on error return (2) .. 4",
			want:   "return STEP() on error return (2)..4",
			result: "[2,3,4]",
		}
		tt.check(t, nil, fail)
		if calls != 2 {
			t.Fatalf("STEP ran %d times, want one call per execution", calls)
		}
	})
}

func TestFormatterCommentLineBreaks(t *testing.T) {
	const item = "let item = { value: 3 }\n"
	tests := []grammarBoundaryCase{
		{
			name:   "array trailing line comment",
			input:  "return [({value:3} // keep\n).value]",
			want:   "return [\n    (\n        { value: 3 }\n        // keep\n    ).value\n]",
			result: "[3]",
		},
		{
			name:   "array leading line comment",
			input:  item + "return [(// keep\nitem).value]",
			want:   item + "return [\n    (\n        // keep\n        item\n    ).value\n]",
			result: "[3]",
		},
		{
			name:   "object trailing line comment",
			input:  "return {v: ({value:3} // keep\n).value}",
			want:   "return {\n    v: (\n        { value: 3 }\n        // keep\n    ).value\n}",
			result: `{"v":3}`,
		},
		{
			name:   "object leading line comment",
			input:  item + "return {v: (// keep\nitem).value}",
			want:   item + "return {\n    v: (\n        // keep\n        item\n    ).value\n}",
			result: `{"v":3}`,
		},
		{
			name:   "array multiline block comment",
			input:  "let item={value:3}\nreturn [(/* keep\n shape */ item).value]",
			want:   item + "return [\n    (\n        /* keep\n        shape */\n        item\n    ).value\n]",
			result: "[3]",
		},
		{
			name:   "blank lines inside block comment",
			input:  item + "return [(/* keep\n\n\n shape */ item).value]",
			want:   item + "return [\n    (\n        /* keep\n\n\n        shape */\n        item\n    ).value\n]",
			result: "[3]",
		},
		{
			name:   "object multiline block comment",
			input:  item + "return {v: (item /* keep\n shape */).value}",
			want:   item + "return {\n    v: (\n        item\n        /* keep\n        shape */\n    ).value\n}",
			result: `{"v":3}`,
		},
		{
			name:   "nested commented grouping",
			input:  item + "return [((item // keep\n)).value]",
			want:   item + "return [\n    ((\n        item\n        // keep\n    )).value\n]",
			result: "[3]",
		},
		{
			name:   "inline block comment containing slashes",
			input:  item + "return [(/* a//b */ item).value]",
			want:   item + "return [(/* a//b */ item).value]",
			result: "[3]",
		},
		{
			name:   "FQL line terminators inside parenthesis comments",
			input:  item + "return [(/* keep\r\n this\u2028 block */ item // then\r// close\u2029).value]",
			want:   item + "return [\n    (\n        /* keep\n        this\n        block */\n        item\n        // then\n        // close\n    ).value\n]",
			result: "[3]",
		},
		{
			name:   "optional breaks in nested arrays still fit inline",
			input:  "return {v:[[1+2]]}",
			want:   "return { v: [ [1 + 2] ] }",
			result: `{"v":[[3]]}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Run("default", func(t *testing.T) { tt.check(t, nil) })
			t.Run("width32", func(t *testing.T) { tt.check(t, []formatter.Option{formatter.WithPrintWidth(32)}) })
		})
	}
}

func (tt grammarBoundaryCase) check(t *testing.T, options []formatter.Option, env ...vm.EnvironmentOption) {
	t.Helper()

	original, err := spec.Compile(tt.input)
	if err != nil {
		t.Fatalf("compile original: %v", err)
	}

	formatted := formatSource(t, tt.input, options...)
	if formatted != tt.want {
		t.Errorf("formatted output:\n%s\nwant:\n%s", formatted, tt.want)
	}

	recompiled, err := spec.Compile(formatted)
	if err != nil {
		t.Fatalf("compile formatted: %v\n%s", err, formatted)
	}

	before, err := spec.Run(original, env...)
	if err != nil || string(before) != tt.result {
		t.Errorf("run original: got %s, want %s, error %v", before, tt.result, err)
	}

	after, err := spec.Run(recompiled, env...)
	if err != nil || string(after) != tt.result {
		t.Errorf("run formatted: got %s, want %s, error %v", after, tt.result, err)
	}

	if second := formatSource(t, formatted, options...); second != formatted {
		t.Errorf("second format changed output:\n%s\nfirst:\n%s", second, formatted)
	}
}
