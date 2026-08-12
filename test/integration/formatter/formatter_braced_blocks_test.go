package formatter_test

import (
	"strings"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/formatter"
	"github.com/MontFerret/ferret/v2/pkg/source"
)

func TestFormatterBracedBlocksAndForStyle(t *testing.T) {
	tests := map[string]struct {
		input string
		want  string
	}{
		"UDF and match braces": {
			input: `func classify(x){return match x{1=>{value:1},_=>{value:0}}} return classify(1)`,
			want:  "func classify(x) {\n    return match x { 1 => { value: 1 }, _ => { value: 0 } }\n}\nreturn classify(1)",
		},
		"empty UDF block": {
			input: `func noop(){}`,
			want:  "func noop() {\n}",
		},
		"comment-only UDF block": {
			input: "func noop() {\n// intentional no-op\n}",
			want:  "func noop() {\n    // intentional no-op\n}",
		},
		"comment-only script": {
			input: "// intentional no-op",
			want:  "// intentional no-op",
		},
		"use-only script": {
			input: "use FOO as F",
			want:  "use FOO as F",
		},
		"direct for return": {
			input: `return for value in [1]{return value}`,
			want:  "return for value in [1] {\n    return value\n}",
		},
		"direct unbraced for return": {
			input: "return for value in [1]\nRETURN value",
			want:  "return for value in [1]\n    return value",
		},
		"direct distinct for return in UDF": {
			input: `func values(){return distinct for value in [1,1]{return value}} return values()`,
			want:  "func values() {\n    return distinct for value in [1, 1] {\n        return value\n    }\n}\nreturn values()",
		},
		"unbraced for retained": {
			input: "for value in [1]\nRETURN value",
			want:  "for value in [1]\n    return value",
		},
		"braced for retained": {
			input: `for value in [1]{return value}`,
			want:  "for value in [1] {\n    return value\n}",
		},
		"mixed nested for retained": {
			input: "for outer in [1] {\nLET inner=(for value in [outer]\nRETURN value)\nRETURN inner\n}",
			want:  "for outer in [1] {\n    let inner = (\n        for value in [outer]\n            return value\n    )\n    return inner\n}",
		},
		"while and do while braces": {
			input: "let first=(for value while value<1{return value}) let second=(for value do while false { return value }) return [first,second]",
			want:  "let first = (\n    for value while value < 1 {\n        return value\n    }\n)\nlet second = (\n    for value do while false {\n        return value\n    }\n)\nreturn [first, second]",
		},
		"braced for comments": {
			input: "for value in [1] {\n// before return\nRETURN value\n// before close\n}",
			want:  "for value in [1] {\n    // before return\n    return value\n    // before close\n}",
		},
		"comments around opening braces": {
			input: "func classify(value) // function header\n{\nRETURN match value // match header\n{\n// before first arm\n1 => 1,\n_ => 0\n}\n}\nLET result = classify(1)\nFOR value in [result] // loop header\n{\n// before return\nRETURN value\n}",
			want:  "func classify(value) { // function header\n    return match value { // match header\n        // before first arm\n        1 => 1,\n        _ => 0,\n    }\n}\nlet result = classify(1)\nfor value in [result] { // loop header\n    // before return\n    return value\n}",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			first := formatSource(t, test.input)
			if first != test.want {
				t.Fatalf("formatted output mismatch\nwant:\n%q\ngot:\n%q", test.want, first)
			}

			second := formatSource(t, first)
			if second != first {
				t.Fatalf("formatter is not idempotent\nfirst:\n%q\nsecond:\n%q", first, second)
			}
		})
	}
}

func formatSource(t *testing.T, input string) string {
	t.Helper()

	var output strings.Builder
	if err := formatter.New().Format(&output, source.NewAnonymous(input)); err != nil {
		t.Fatalf("format source: %v", err)
	}

	return output.String()
}
