package formatter_test

import (
	"testing"

	"github.com/MontFerret/ferret/v2/test/spec"
)

func TestFormatterGrammarBoundaries(t *testing.T) {
	const fallback = `"Default value returned when no shipping option matches the destination address"`

	tests := []struct {
		name   string
		input  string
		want   string
		result string
	}{
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
			original, err := spec.Compile(tt.input)
			if err != nil {
				t.Fatalf("compile original: %v", err)
			}

			formatted := formatSource(t, tt.input)
			if formatted != tt.want {
				t.Errorf("formatted output:\n%s\nwant:\n%s", formatted, tt.want)
			}

			recompiled, err := spec.Compile(formatted)
			if err != nil {
				t.Fatalf("compile formatted: %v\n%s", err, formatted)
			}

			if second := formatSource(t, formatted); second != formatted {
				t.Errorf("second format changed output:\n%s\nfirst:\n%s", second, formatted)
			}

			before, err := spec.Run(original)
			if err != nil || string(before) != tt.result {
				t.Errorf("run original: got %s, want %s, error %v", before, tt.result, err)
			}

			after, err := spec.Run(recompiled)
			if err != nil || string(after) != tt.result {
				t.Errorf("run formatted: got %s, want %s, error %v", after, tt.result, err)
			}
		})
	}
}
