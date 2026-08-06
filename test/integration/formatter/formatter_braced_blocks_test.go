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
		"UDF and MATCH braces": {
			input: `FUNC classify(x){RETURN MATCH x{1=>{value:1},_=>{value:0}}} RETURN classify(1)`,
			want:  "FUNC classify(x) {\n    RETURN MATCH x { 1 => { value: 1 }, _ => { value: 0 } }\n}\nRETURN classify(1)",
		},
		"unbraced FOR retained": {
			input: "FOR value IN [1]\nRETURN value",
			want:  "FOR value IN [1]\n    RETURN value",
		},
		"braced FOR retained": {
			input: `FOR value IN [1]{RETURN value}`,
			want:  "FOR value IN [1] {\n    RETURN value\n}",
		},
		"mixed nested FOR retained": {
			input: "FOR outer IN [1] {\nLET inner=(FOR value IN [outer]\nRETURN value)\nRETURN inner\n}",
			want:  "FOR outer IN [1] {\n    LET inner = (\n        FOR value IN [outer]\n            RETURN value\n    )\n    RETURN inner\n}",
		},
		"WHILE and DO WHILE braces": {
			input: "LET first=(FOR value WHILE value<1{RETURN value}) LET second=(FOR value DO WHILE false { RETURN value }) RETURN [first,second]",
			want:  "LET first = (\n    FOR value WHILE value < 1 {\n        RETURN value\n    }\n)\nLET second = (\n    FOR value DO WHILE FALSE {\n        RETURN value\n    }\n)\nRETURN [first, second]",
		},
		"braced FOR comments": {
			input: "FOR value IN [1] {\n// before return\nRETURN value\n// before close\n}",
			want:  "FOR value IN [1] {\n    // before return\n    RETURN value\n    // before close\n}",
		},
		"comments around opening braces": {
			input: "FUNC classify(value) // function header\n{\nRETURN MATCH value // match header\n{\n// before first arm\n1 => 1,\n_ => 0\n}\n}\nLET result = classify(1)\nFOR value IN [result] // loop header\n{\n// before return\nRETURN value\n}",
			want:  "FUNC classify(value) { // function header\n    RETURN MATCH value { // match header\n        // before first arm\n        1 => 1,\n        _ => 0,\n    }\n}\nLET result = classify(1)\nFOR value IN [result] { // loop header\n    // before return\n    RETURN value\n}",
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
