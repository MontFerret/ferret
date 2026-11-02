package compiler_test

import (
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/compiler"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/test/spec"
)

func TestCompiledStringValues(t *testing.T) {
	// These expectations describe FQL decoding, not Go/JSON/SQL escaping.
	tests := []struct {
		name    string
		literal string
		want    string
	}{
		{"empty", `''`, ""},
		{"recognized escapes", `"a\nb\tc"`, "a\nb\tc"},
		{"unknown escapes", `"\b\q\r\u0041\U0001F600"`, `\b\q\r\u0041\U0001F600`},
		{"ordinary backslashes", `"a\\b"`, `a\\b`},
		{"escaped double quote", `"a\"b"`, `a\"b`},
		{"escaped single quote", `'a\'b'`, `a\'b`},
		{"doubled double quote", `"a""b"`, `a""b`},
		{"doubled single quote", `'a''b'`, `a''b`},
		{"raw controls", "\"a\n\tb\rc\"", "a\n\tb\rc"},
		{"backslash before newline", "\"a\\\nb\"", "a\\\nb"},
		{"unicode escapes", `'é\中\n😀'`, "é\\中\n😀"},
		{"acute", `´abc´`, "abc"},
		{"acute escaped delimiter", `´a\´b´`, `a\´b`},
		{"acute trailing backslash", `´a\´`, `a\´`},
		{"acute trailing pair", `´a\\´`, `a\\`},
		{"template backslashes", "`a\\\\b`", `a\b`},
		{"template recognized escapes", "`a\\nb\\tc`", "a\nb\tc"},
		{"template escaped dollar", "`a\\${b}`", "a${b}"},
		{"template escaped backtick", "`a\\`b`", "a`b"},
		{"template ordinary interpolation", "`pre${'a\"b'}post`", `prea"bpost`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			program, err := spec.Compile("RETURN "+tt.literal, compiler.None)
			if err != nil {
				t.Fatalf("compile: %v", err)
			}

			if len(program.Constants) != 1 {
				t.Fatalf("constants = %v, want one string", program.Constants)
			}

			value, ok := program.Constants[0].(runtime.String)
			if !ok || string(value) != tt.want {
				t.Fatalf("constant = (%T) %q, want runtime.String %q", program.Constants[0], value, tt.want)
			}
		})
	}
}
